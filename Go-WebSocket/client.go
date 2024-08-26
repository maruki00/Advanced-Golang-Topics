package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"golang.org/x/net/websocket"
	"golang.org/x/text/message"

	

)

func main() {
	con, err := websocket.Dial("ws://127.0.0.1:8085", "", createIP())


	if err != nil {
		log.Fatal(err.Error())
	}
	defer con.Close()

	go recieveMessage(con)
	//sendMessage(con)
}

func createIP() string {
	var ip [4]int

	for i := 0; i < len(ip); i++ {
		rand.Seed(time.Now().UnixNano())
		ip[i] = rand.Intn(256)
	}

	return fmt.Sprintf("http://%d.%d.%d.%d", ip[0], ip[1], ip[2], ip[3])
}

func recieveMessage(con *websocket.Conn) {
	for {

		var message  socket_pkg.Message
		err := websocket.JSON.Receive(con, message)
		if err != nil {
			log.Fatal("Error recieving Data : ", err)
			continue
		}
		fmt.Println("Message recieved : ", message.Subject)

		select {}
	}
}


func sendMessage(con *websocket.Conn) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()

		message := socket.Message{
			sub
		}
	}
}