package main

import (
	"log"

	"github.com/nats-io/nats.go"
)

func main() {

	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	err = nc.Publish("Updates", []byte("Hello NAT's"))
	if err != nil {
		log.Println("Could not pusblish the message!")
	}

	log.Println("Connected to NATS server")

}
