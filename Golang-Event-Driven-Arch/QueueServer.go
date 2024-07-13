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

	// Subscribe to the 'updates' subject with a queue group
	_, err = nc.QueueSubscribe("updates", "workers", func(m *nats.Msg) {
		log.Printf("Worker received message: %s", string(m.Data))
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Subscribed to 'updates' subject with queue group 'workers'")

	// Keep the connection alive
	select {}
}
