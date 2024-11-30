package main

import "github.com/google/wire"



func InintApp(db map[int]any) *App {
	wire.Build(
		NewRepository,
		NewService,
		NewApp,
	)	
	return &App{}
}
