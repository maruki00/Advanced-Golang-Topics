package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// User represents a user in our system
type User struct {
	gorm.Model
	Username     string `json:"username" gorm:"uniqueIndex"`
	PasswordHash string `json:"-"`
	Email        string `json:"email" gorm:"uniqueIndex"`
	Role         string `json:"role"`
}

// UserDTO for registration and login
type UserDTO struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Role     string `json:"role"`
}

// LoginDTO for login requests
type LoginDTO struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Claims defines the structure for JWT claims
type Claims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

var db *gorm.DB
var enforcer *casbin.Enforcer
var jwtSecret []byte

func init() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using default or environment variables")
	}

	// Get JWT secret from environment variable or use a default
	jwtSecret = []byte(getEnv("JWT_SECRET", "your-256-bit-secret"))
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

// Setup database
func setupDB() {
	var err error
	db, err = gorm.Open(sqlite.Open("auth.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto migrate the schema
	db.AutoMigrate(&User{})
}

// Setup Casbin enforcer
func setupCasbin() {
	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		log.Fatalf("Failed to initialize Casbin adapter: %v", err)
	}

	enforcer, err = casbin.NewEnforcer("auth_model.conf", adapter)
	if err != nil {
		log.Fatalf("Failed to create enforcer: %v", err)
	}

	// Load policies
	if err := enforcer.LoadPolicy(); err != nil {
		log.Fatalf("Failed to load policy: %v", err)
	}

	// Add policies if none exist
	policies, err := enforcer.GetPolicy()
	if err != nil || len(policies) == 0 {
		// Admin policies
		enforcer.AddPolicy("admin", "/api/admin/*", "*")
		enforcer.AddPolicy("admin", "/api/users", "*")
		enforcer.AddPolicy("admin", "/api/users/*", "*")

		// User policies
		enforcer.AddPolicy("user", "/api/users", "GET")
		enforcer.AddPolicy("user", "/api/users/:id", "GET")
		enforcer.AddPolicy("user", "/api/profile", "GET")
		enforcer.AddPolicy("user", "/api/profile", "PUT")

		// Guest (unauthenticated) policies
		enforcer.AddPolicy("guest", "/api/auth/login", "POST")
		enforcer.AddPolicy("guest", "/api/auth/register", "POST")

		// Save policy changes
		enforcer.SavePolicy()
	}
}

// HashPassword creates a bcrypt hash from a password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPasswordHash compares a password with a hash
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateJWT creates a new JWT token
func GenerateJWT(user User) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   user.Username,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	return tokenString, err
}

// Authentication middleware to validate JWT tokens
func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip authentication for public routes
		if c.Request.URL.Path == "/api/auth/login" || c.Request.URL.Path == "/api/auth/register" {
			c.Set("role", "guest")
			c.Next()
			return
		}

		// Get token from Authorization header
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.Set("role", "guest")
			c.Next()
			return
		}

		// Remove "Bearer " prefix if present
		if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
			tokenString = tokenString[7:]
		}

		// Parse and validate the token
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.Set("role", "guest")
			c.Next()
			return
		}

		// Set user information in context
		var user User
		if result := db.Where("username = ?", claims.Username).First(&user); result.Error != nil {
			c.Set("role", "guest")
			c.Next()
			return
		}

		c.Set("user", user)
		c.Set("role", user.Role)
		c.Next()
	}
}

// Authorization middleware using Casbin
func Authorization() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get role from context
		role, exists := c.Get("role")
		if !exists {
			role = "guest"
		}

		// Get request path and method
		path := c.Request.URL.Path
		method := c.Request.Method

		// Check permission
		allowed, err := enforcer.Enforce(role, path, method)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Authorization error",
			})
			return
		}

		if !allowed {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "Forbidden",
			})
			return
		}

		c.Next()
	}
}

func main() {
	// Setup database and Casbin
	setupDB()
	setupCasbin()

	// Create default admin if none exists
	var adminCount int64
	db.Model(&User{}).Where("role = ?", "admin").Count(&adminCount)
	if adminCount == 0 {
		hashedPassword, _ := HashPassword("admin123")
		admin := User{
			Username:     "admin",
			PasswordHash: hashedPassword,
			Email:        "admin@example.com",
			Role:         "admin",
		}
		db.Create(&admin)
	}

	// Initialize Gin
	r := gin.Default()

	// Public endpoints
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Welcome to the Auth API",
		})
	})

	// API routes with middleware
	api := r.Group("/api")
	api.Use(Authentication(), Authorization())
	{
		// Auth routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", func(c *gin.Context) {
				var userDTO UserDTO
				if err := c.ShouldBindJSON(&userDTO); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}

				// Check if username already exists
				var existingUser User
				if result := db.Where("username = ?", userDTO.Username).First(&existingUser); result.Error == nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Username already exists"})
					return
				}

				// Check if email already exists
				if result := db.Where("email = ?", userDTO.Email).First(&existingUser); result.Error == nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Email already exists"})
					return
				}

				// Hash the password
				hashedPassword, err := HashPassword(userDTO.Password)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
					return
				}

				// Set default role if not provided
				if userDTO.Role == "" {
					userDTO.Role = "user"
				}

				// Create user
				user := User{
					Username:     userDTO.Username,
					PasswordHash: hashedPassword,
					Email:        userDTO.Email,
					Role:         userDTO.Role,
				}

				if result := db.Create(&user); result.Error != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
					return
				}

				c.JSON(http.StatusCreated, gin.H{
					"message": "User registered successfully",
					"user": gin.H{
						"username": user.Username,
						"email":    user.Email,
						"role":     user.Role,
					},
				})
			})

			auth.POST("/login", func(c *gin.Context) {
				var loginDTO LoginDTO
				if err := c.ShouldBindJSON(&loginDTO); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}

				// Find user by username
				var user User
				if result := db.Where("username = ?", loginDTO.Username).First(&user); result.Error != nil {
					c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
					return
				}

				// Check password
				if !CheckPasswordHash(loginDTO.Password, user.PasswordHash) {
					c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
					return
				}

				// Generate JWT token
				token, err := GenerateJWT(user)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
					return
				}

				c.JSON(http.StatusOK, gin.H{
					"message": "Login successful",
					"token":   token,
					"user": gin.H{
						"username": user.Username,
						"email":    user.Email,
						"role":     user.Role,
					},
				})
			})
		}

		// Admin routes
		admin := api.Group("/admin")
		{
			admin.GET("/users", func(c *gin.Context) {
				var users []User
				db.Find(&users)

				c.JSON(http.StatusOK, gin.H{
					"users": users,
				})
			})
		}

		// User routes
		users := api.Group("/users")
		{
			users.GET("", func(c *gin.Context) {
				var users []struct {
					Username string
					Email    string
					Role     string
				}
				db.Model(&User{}).Select("username, email, role").Find(&users)

				c.JSON(http.StatusOK, gin.H{
					"users": users,
				})
			})

			users.GET("/:id", func(c *gin.Context) {
				id := c.Param("id")

				var user User
				if result := db.Where("username = ?", id).First(&user); result.Error != nil {
					c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
					return
				}

				c.JSON(http.StatusOK, gin.H{
					"user": gin.H{
						"username": user.Username,
						"email":    user.Email,
						"role":     user.Role,
					},
				})
			})
		}

		// Profile routes
		api.GET("/profile", func(c *gin.Context) {
			userInterface, exists := c.Get("user")
			if !exists {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
				return
			}

			user := userInterface.(User)
			c.JSON(http.StatusOK, gin.H{
				"user": gin.H{
					"username": user.Username,
					"email":    user.Email,
					"role":     user.Role,
				},
			})
		})

		api.PUT("/profile", func(c *gin.Context) {
			userInterface, exists := c.Get("user")
			if !exists {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
				return
			}

			user := userInterface.(User)

			var updateDTO struct {
				Email    string `json:"email"`
				Password string `json:"password"`
			}

			if err := c.ShouldBindJSON(&updateDTO); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			// Update user fields
			if updateDTO.Email != "" {
				user.Email = updateDTO.Email
			}

			if updateDTO.Password != "" {
				hashedPassword, err := HashPassword(updateDTO.Password)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
					return
				}
				user.PasswordHash = hashedPassword
			}

			// Save changes
			if result := db.Save(&user); result.Error != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"message": "Profile updated successfully",
				"user": gin.H{
					"username": user.Username,
					"email":    user.Email,
					"role":     user.Role,
				},
			})
		})
	}

	// Run the server
	port := getEnv("PORT", "8080")
	log.Printf("Server starting on port %s", port)
	r.Run(":" + port)
}
