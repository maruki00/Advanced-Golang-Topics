package main

import (
	"fmt"
)

type App struct {
	rep Repository
	srv Service
}

func NewApp(repo Repository, srv Service) *App {
	fmt.Println("🧑🏼‍🎄")
	return &App{
		rep: repo,
		srv: srv,
	}
}
