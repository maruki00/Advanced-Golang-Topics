package main

import (
	"time"

	"github.com/google/wire"
)

// InitializeGreeter creates a Greeter instance with all dependencies injected
func InitializeGreeter() Greeter {
	wire.Build(NewGreeter, GetCurrentTimeMessage, time.Now)
	return Greeter{}
}
