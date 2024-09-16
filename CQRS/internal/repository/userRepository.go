package repository

import (
	"CRQS-GO/internal/commands"
	"CRQS-GO/internal/models"
	"CRQS-GO/internal/queries"
	"fmt"
	"sync"
)

type Repository struct {
	mu    sync.RWMutex
	store map[int]models.User
}

func NewRepository() *Repository {
	return &Repository{
		store: make(map[int]models.User),
	}
}

func (r *Repository) HandleCreateUserCommand(cmd commands.CreateUserCommand) {
	r.mu.Lock()
	defer r.mu.Unlock()

	user := models.User{ID: cmd.ID, Name: cmd.Name}
	r.store[cmd.ID] = user
	fmt.Println("User created:", user)
}
func (r *Repository) HandleGetUserQuery(query queries.GetUserQuery) (models.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.store[query.ID]
	if !exists {
		return models.User{}, fmt.Errorf("user with ID %d not found", query.ID)
	}

	return user, nil
}
