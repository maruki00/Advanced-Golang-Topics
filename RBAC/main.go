package main

import (
	"fmt"
)

// Role Base Access Control
type Role struct {
	Name        string
	Permissions []string
}

type User struct {
	Name string
	Role Role
}

func (u User) HasPermission(permission string) bool {
	for _, p := range u.Role.Permissions {
		if p == permission {
			return true
		}
	}
	return false
}

func main() {

	adminRole := Role{"Admin", []string{"read", "write", "delete"}}
	editorRole := Role{"Editor", []string{"read", "write"}}
	viewerRole := Role{"Viewer", []string{"read"}}

	users := []User{
		{"Alice", editorRole},
		{"Bob", viewerRole},
		{"Charlie", adminRole},
	}

	for _, user := range users {
		fmt.Printf("%s can delete: %v\n", user.Name, user.HasPermission("delete"))
	}
}
