package main

import (
	"fmt"
	"net"
)

type Pair struct {
	Nikname string
	Conn net.Conn
	Addr net.Addr
}

type Server struct {
	onlines map[string]bool
	pairs map[string][2]Pair
}


func NewServer() *Server{

}

func NewPair(nikename string, conn net.Conn, addr net.Addr) *Pair{
	return &Pair{
		Nikname : nikename,
		Conn : con,
		Addr : addr,
	}
}


func main() {
	onlines := map[string]bool,0)
	// Listen for incoming connections
	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Server is listening on port 8080")

	for {
		// Accept incoming connections
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		// Handle client connection in a goroutine
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
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

}
