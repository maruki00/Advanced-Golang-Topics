package main

import (
	"fmt"
	"net"
)

func main() {
	// Connect to the server
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer conn.Close()

	for {
		_, err = conn.Write([]byte("hello world"))
		if err != nil {
			fmt.Println("error sending data")
		}

		handleServer(conn)
	}

}

func handleServer(conn net.Conn) {
	defer conn.Close()
	buffer := make([]byte, 1024)

	for {
		// Read data from the client
		n, err := conn.Read(buffer)
		if err != nil {
			return
		}

		// Process and use the data (here, we'll just print it)
		fmt.Printf("Received: %s\n", buffer[:n])
	}

	_, _ = conn.Write([]byte("recived thanks ."))
}
