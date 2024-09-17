package main

import "fmt"

type Greeter struct {
	message Message
}

func NewGreeter(m Message) Greeter {
	return Greeter{message: m}
}

func (g Greeter) Greet() {
	fmt.Println(g.message.String())
}
