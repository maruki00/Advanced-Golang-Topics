package main

import "time"

type Message struct {
	content string
}

func NewMessage() Message {
	return Message{content: "Hello, World!"}
}

func GetCurrentTimeMessage(t time.Time) Message {
	if t.Hour() < 12 {
		return Message{content: "Good Morning!"}
	}
	return Message{content: "Good Afternoon!"}
}

func (m Message) String() string {
	return m.content
}
