package main

import (
	"CRQS-GO/internal/commands"
	"CRQS-GO/internal/queries"
	"CRQS-GO/internal/repository"
	"fmt"
)

func main() {
	repo := repository.NewRepository()

	// Command: Create a user
	createCmd := commands.CreateUserCommand{ID: 1, Name: "John Doe"}
	repo.HandleCreateUserCommand(createCmd)

	// Query: Fetch the user by ID
	query := queries.GetUserQuery{ID: 1}
	user, err := repo.HandleGetUserQuery(query)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("User found:", user)
	}
}
