package main

import (
	"fmt"
)

// Access Control List
const (
	Read  = "read"
	Write = "write"
	Exec  = "execute"
)

type ACL map[string][]string

func (acl ACL) HasPermission(user, permission string) bool {
	for _, p := range acl[user] {
		if p == permission {
			return true
		}
	}
	return false
}

func main() {
	fileACL := ACL{
		"Alice": {"read", "write"},
		"Bob":   {"read", "execute"},
	}

	fmt.Println("Alice can write:", fileACL.HasPermission("Alice", Write))
	fmt.Println("Bob can write:", fileACL.HasPermission("Bob", Write))
}
