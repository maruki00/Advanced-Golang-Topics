package socket

import (
	"fmt"
	"log"

	"golang.org/x/net/websocket"
)

type Message struct {
	Subject string `json: "subject"`
}

type Config struct {
	Clients        map[string]*websocket.Conn
	RegisterClient chan *websocket.Conn
	RemoveClient   chan *websocket.Conn
	Message        chan Message
}

func NewConfig() *Config {
	return &Config{
		Clients:        make(map[string]*websocket.Conn),
		RegisterClient: make(chan *websocket.Conn),
		RemoveClient:   make(chan *websocket.Conn),
		MessageData:        make(chan Message),
	}
}

func (config *Config) RegisterClient(client *websocket.Conn) {
	config.Clients[client.RemoteAddr().String()] = client
	fmt.Println("clients : ", config.Clients)

}

func (config *Config) RemoveClient(client *websocket.Conn) {
	delete(config.Clients, client.RemoteAddr().String())
	fmt.Println("clients : ", config.Clients)
}

func (Config *Config) MessageData(message Message) {
	for _, client := range Config.Clients {
		err := websocket.JSON.Send(client, message)
		if err != nil {
			log.Fatal("[+] Error Sending : ", err.Error())
		}
	}
}

func (config *Config) RunSocket() {
	for {
		select {
		case registerClient := <-config.RegisterClient:
			config.RegisterClient(registerClient)

		
		case removeClient := <- config.RemoveClient
			config.RemoveClient(removeClient)
		
		case messageData := <- config.MessageData:
			config.MessageData(messageData)




		}
	}
}
