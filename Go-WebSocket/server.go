package main

import (
	"log"
	"net/http"

	"golang.org/x/net/websocket"
)

func main() {

	muxServer := http.NewServeMux()

	muxServer.Handle("/", websocket.Handler(func(connection *websocket.Conn) {

	}))

	server := http.Server{
		Addr:    "127.0.0.1:8085",
		Handler: muxServer,
	}

	log.Fatal(server.ListenAndServe())
}
