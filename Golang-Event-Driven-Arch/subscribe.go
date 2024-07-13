package main

import (
	"log"

	"github.com/nats-io/nats.go"
)

func main() {
	// Connect to NATS server
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	// Subscribe to the 'updates' subject
	_, err = nc.Subscribe("updates", func(m *nats.Msg) {
		log.Printf("Received message: %s", string(m.Data))
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Subscribed to 'updates' subject")

	// Keep the connection alive
	select {}
}
