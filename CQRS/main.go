package main

import (
	"fmt"
	"sync"
)

// User is a simple model representing a user
type User struct {
	ID   int
	Name string
}

// Repository simulates a database using an in-memory store
type Repository struct {
	mu    sync.RWMutex
	store map[int]User
}

// NewRepository creates a new instance of the repository
func NewRepository() *Repository {
	return &Repository{
		store: make(map[int]User),
	}
}

// Command Handlers (Write-side)

// CreateUserCommand represents a command to create a new user
type CreateUserCommand struct {
	ID   int
	Name string
}

// HandleCreateUserCommand handles the logic for creating a user
func (r *Repository) HandleCreateUserCommand(cmd CreateUserCommand) {
	r.mu.Lock()
	defer r.mu.Unlock()

	user := User{ID: cmd.ID, Name: cmd.Name}
	r.store[cmd.ID] = user
	fmt.Println("User created:", user)
}

// Query Handlers (Read-side)

// GetUserQuery represents a query to get user data by ID
type GetUserQuery struct {
	ID int
}

// HandleGetUserQuery handles the logic for fetching a user by ID
func (r *Repository) HandleGetUserQuery(query GetUserQuery) (User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.store[query.ID]
	if !exists {
		return User{}, fmt.Errorf("user with ID %d not found", query.ID)
	}

	return user, nil
}

func main() {
	repo := NewRepository()

	// Command: Create a user
	createCmd := CreateUserCommand{ID: 1, Name: "John Doe"}
	repo.HandleCreateUserCommand(createCmd)

	// Query: Fetch the user by ID
	query := GetUserQuery{ID: 1}
	user, err := repo.HandleGetUserQuery(query)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("User found:", user)
	}
}
