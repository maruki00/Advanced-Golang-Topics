package main

import (
	"fmt"
	"sync"
)

// User is a simple model representing a user



// Command Handlers (Write-side)

// CreateUserCommand represents a command to create a new user


// HandleCreateUserCommand handles the logic for creating a user


// Query Handlers (Read-side)

// GetUserQuery represents a query to get user data by ID


// HandleGetUserQuery handles the logic for fetching a user by ID


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
