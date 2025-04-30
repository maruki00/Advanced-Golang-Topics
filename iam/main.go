package main

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

//go:embed templates/*
var templatesFS embed.FS

//go:embed static/*
var staticFS embed.FS

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

// Role represents a role in the system
type Role struct {
	ID   uint   `json:"id" gorm:"primaryKey"`
	Name string `json:"name" gorm:"uniqueIndex"`
}

// CasbinRule is the structure for policy rules
type CasbinRule struct {
	ID    uint   `json:"id" gorm:"primaryKey"`
	Ptype string `json:"ptype"`
	V0    string `json:"v0"` // subject/role
	V1    string `json:"v1"` // object/resource
	V2    string `json:"v2"` // action
	V3    string `json:"v3"`
	V4    string `json:"v4"`
	V5    string `json:"v5"`
}

var db *gorm.DB
var enforcer *casbin.Enforcer
var jwtSecret []byte
var adapter *gormadapter.Adapter

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
	db.AutoMigrate(&User{}, &Role{})
}

// Setup Casbin enforcer
func setupCasbin() {
	var err error
	adapter, err = gormadapter.NewAdapterByDB(db)
	if err != nil {
		log.Fatalf("Failed to initialize Casbin adapter: %v", err)
	}

	// Load model from file
	m, err := model.NewModelFromFile("auth_model.conf")
	if err != nil {
		log.Fatalf("Failed to load model: %v", err)
	}

	enforcer, err = casbin.NewEnforcer(m, adapter)
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
		enforcer.AddPolicy("admin", "/api/roles", "*")
		enforcer.AddPolicy("admin", "/api/roles/*", "*")
		enforcer.AddPolicy("admin", "/api/policies", "*")
		enforcer.AddPolicy("admin", "/api/policies/*", "*")
		enforcer.AddPolicy("admin", "/api/permissions", "*")
		enforcer.AddPolicy("admin", "/api/permissions/*", "*")
		enforcer.AddPolicy("admin", "/admin/*", "*")

		// User policies
		enforcer.AddPolicy("user", "/api/users", "GET")
		enforcer.AddPolicy("user", "/api/users/:id", "GET")
		enforcer.AddPolicy("user", "/api/profile", "GET")
		enforcer.AddPolicy("user", "/api/profile", "PUT")

		// Guest (unauthenticated) policies
		enforcer.AddPolicy("guest", "/api/auth/login", "POST")
		enforcer.AddPolicy("guest", "/api/auth/register", "POST")
		enforcer.AddPolicy("guest", "/", "GET")
		enforcer.AddPolicy("guest", "/login", "GET")
		enforcer.AddPolicy("guest", "/register", "GET")
		enforcer.AddPolicy("guest", "/static/*", "GET")

		// Save policy changes
		enforcer.SavePolicy()
	}

	// Create roles table if it doesn't exist
	var roles []Role
	if result := db.Find(&roles); result.Error != nil || len(roles) == 0 {
		defaultRoles := []Role{
			{Name: "admin"},
			{Name: "user"},
			{Name: "guest"},
		}
		db.Create(&defaultRoles)
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
		// Skip authentication for public routes and static content
		if c.Request.URL.Path == "/api/auth/login" ||
			c.Request.URL.Path == "/api/auth/register" ||
			c.Request.URL.Path == "/" ||
			c.Request.URL.Path == "/login" ||
			c.Request.URL.Path == "/register" {
			c.Set("role", "guest")
			c.Next()
			return
		}

		// Check for session authentication
		session := sessions.Default(c)
		userID := session.Get("user_id")
		if userID != nil {
			var user User
			if result := db.First(&user, userID); result.Error == nil {
				c.Set("user", user)
				c.Set("role", user.Role)
				c.Next()
				return
			}
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
			// For API routes, return JSON error
			if c.GetHeader("Accept") == "application/json" || c.Request.URL.Path[:5] == "/api/" {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
					"error": "Forbidden",
				})
				return
			}

			// For web routes, redirect to login page
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		c.Next()
	}
}

// SetupTemplates configures the HTML templates
func SetupTemplates() *template.Template {
	templ := template.Must(template.New("").ParseFS(templatesFS, "templates/*.html"))
	return templ
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

	// Setup templates
	templ := SetupTemplates()
	r.SetHTMLTemplate(templ)

	// Setup sessions
	store := cookie.NewStore([]byte(getEnv("SESSION_SECRET", "secret")))
	r.Use(sessions.Sessions("auth-session", store))

	// Setup CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Serve static files
	staticFiles, _ := fs.Sub(staticFS, "static")
	r.StaticFS("/static", http.FS(staticFiles))

	// Web UI routes
	r.GET("/", func(c *gin.Context) {
		session := sessions.Default(c)
		userID := session.Get("user_id")

		if userID != nil {
			c.Redirect(http.StatusFound, "/admin/dashboard")
			return
		}

		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Auth System",
		})
	})

	r.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", gin.H{
			"title": "Login",
		})
	})

	r.GET("/register", func(c *gin.Context) {
		c.HTML(http.StatusOK, "register.html", gin.H{
			"title": "Register",
		})
	})

	// Admin dashboard routes
	adminRoutes := r.Group("/admin")
	adminRoutes.Use(Authentication(), Authorization())
	{
		adminRoutes.GET("/dashboard", func(c *gin.Context) {
			user, exists := c.Get("user")
			if !exists {
				c.Redirect(http.StatusFound, "/login")
				return
			}

			c.HTML(http.StatusOK, "dashboard.html", gin.H{
				"title": "Admin Dashboard",
				"user":  user.(User),
			})
		})

		adminRoutes.GET("/users", func(c *gin.Context) {
			var users []User
			db.Find(&users)

			c.HTML(http.StatusOK, "users.html", gin.H{
				"title": "User Management",
				"users": users,
			})
		})

		adminRoutes.GET("/roles", func(c *gin.Context) {
			var roles []Role
			db.Find(&roles)

			c.HTML(http.StatusOK, "roles.html", gin.H{
				"title": "Role Management",
				"roles": roles,
			})
		})

		adminRoutes.GET("/policies", func(c *gin.Context) {
			policies, _ := enforcer.GetPolicy()

			c.HTML(http.StatusOK, "policies.html", gin.H{
				"title":    "Policy Management",
				"policies": policies,
			})
		})
	}

	// API routes
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

				// Generate JWT token
				token, err := GenerateJWT(user)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
					return
				}

				// Set session
				session := sessions.Default(c)
				session.Set("user_id", user.ID)
				session.Save()

				c.JSON(http.StatusCreated, gin.H{
					"message": "User registered successfully",
					"token":   token,
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

				// Set session
				session := sessions.Default(c)
				session.Set("user_id", user.ID)
				session.Save()

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

			auth.POST("/logout", func(c *gin.Context) {
				session := sessions.Default(c)
				session.Clear()
				session.Save()

				c.JSON(http.StatusOK, gin.H{
					"message": "Logout successful",
				})
			})
		}

		// Role management API
		roles := api.Group("/roles")
		{
			roles.GET("", func(c *gin.Context) {
				var roles []Role
				db.Find(&roles)

				c.JSON(http.StatusOK, gin.H{
					"roles": roles,
				})
			})

			roles.POST("", func(c *gin.Context) {
				var role Role
				if err := c.ShouldBindJSON(&role); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}

				// Check if role already exists
				var existingRole Role
				if result := db.Where("name = ?", role.Name).First(&existingRole); result.Error == nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Role already exists"})
					return
				}

				if result := db.Create(&role); result.Error != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create role"})
					return
				}

				c.JSON(http.StatusCreated, gin.H{
					"message": "Role created successfully",
					"role":    role,
				})
			})

			roles.DELETE("/:id", func(c *gin.Context) {
				id := c.Param("id")
				var role Role

				// Check if role exists
				if result := db.First(&role, id); result.Error != nil {
					c.JSON(http.StatusNotFound, gin.H{"error": "Role not found"})
					return
				}

				// Check if role is in use
				var userCount int64
				db.Model(&User{}).Where("role = ?", role.Name).Count(&userCount)
				if userCount > 0 {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Role is in use by users"})
					return
				}

				// Delete role policies
				policies, _ := enforcer.GetFilteredPolicy(0, role.Name)
				for _, policy := range policies {
					enforcer.RemovePolicy(policy)
				}
				enforcer.SavePolicy()

				// Delete role
				db.Delete(&role)

				c.JSON(http.StatusOK, gin.H{
					"message": "Role deleted successfully",
				})
			})
		}

		// Policy management API
		policies := api.Group("/policies")
		{
			policies.GET("", func(c *gin.Context) {
				policies, _ := enforcer.GetPolicy()

				c.JSON(http.StatusOK, gin.H{
					"policies": policies,
				})
			})

			policies.POST("", func(c *gin.Context) {
				var policy []interface{}
				if err := c.ShouldBindJSON(&policy); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}

				if len(policy) < 3 {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Policy must have at least subject, object, and action"})
					return
				}

				added, err := enforcer.AddPolicy(policy...)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add policy"})
					return
				}

				if !added {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Policy already exists"})
					return
				}

				enforcer.SavePolicy()

				c.JSON(http.StatusCreated, gin.H{
					"message": "Policy added successfully",
					"policy":  policy,
				})
			})

			policies.DELETE("", func(c *gin.Context) {
				var policy []interface{}
				if err := c.ShouldBindJSON(&policy); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}

				if len(policy) < 3 {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Policy must have at least subject, object, and action"})
					return
				}

				// Remove policy
				removed, err := enforcer.RemovePolicy(policy...)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove policy"})
					return
				}

				if !removed {
					c.JSON(http.StatusNotFound, gin.H{"error": "Policy not found"})
					return
				}

				enforcer.SavePolicy()

				c.JSON(http.StatusOK, gin.H{
					"message": "Policy removed successfully",
				})
			})

			// Get policies by role
			policies.GET("/role/:role", func(c *gin.Context) {
				role := c.Param("role")
				policies, _ := enforcer.GetFilteredPolicy(0, role)

				c.JSON(http.StatusOK, gin.H{
					"role":     role,
					"policies": policies,
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

			// Get Casbin model
			admin.GET("/model", func(c *gin.Context) {
				model := enforcer.GetModel()
				modelAsText := ""

				for sec, t := range model {
					modelAsText += fmt.Sprintf("[%s]\n", sec)
					for k, v := range t {
						modelAsText += fmt.Sprintf("%s = %s\n", k, v)
					}
					modelAsText += "\n"
				}

				c.JSON(http.StatusOK, gin.H{
					"model": modelAsText,
				})
			})

			// Update user role
			admin.PUT("/users/:id/role", func(c *gin.Context) {
				id := c.Param("id")
				var user User

				if result := db.First(&user, id); result.Error != nil {
					c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
					return
				}

				var roleUpdate struct {
					Role string `json:"role" binding:"required"`
				}

				if err := c.ShouldBindJSON(&roleUpdate); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}

				// Check if role exists
				var role Role
				if result := db.Where("name = ?", roleUpdate.Role).First(&role); result.Error != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Role does not exist"})
					return
				}

				// Update user role
				user.Role = roleUpdate.Role
				db.Save(&user)

				c.JSON(http.StatusOK, gin.H{
					"message": "User role updated successfully",
					"user": gin.H{
						"id":       user.ID,
						"username": user.Username,
						"email":    user.Email,
						"role":     user.Role,
					},
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
					"id":       user.ID,
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
