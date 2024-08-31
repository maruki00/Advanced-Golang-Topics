package main

import (
	"flag"
	"fmt"
	"net"
)

func main() {
	name := flag.String("name", "anonymous", "This ur nikename")
	flag.Parse()

	// Connect to the server
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	defer conn.Close()

	for {
		_, err = conn.Write([]byte(*name))
		if err != nil {
			//fmt.Println("error sending data ", err.Error())
			panic(err.Error())
		}

		fmt.Print(*name)
		go handleServer(conn, *name)
	}

}

func handleServer(conn net.Conn, nName string) {
	defer conn.Close()
	// fmt.Println(nName)
	buffer := make([]byte, 1024)
	return
	for {
		// Read data from the client
		n, err := conn.Read(buffer)
		if err != nil {
			panic(err.Error())
		}
		// Process and use the data (here, we'll just print it)
		fmt.Printf("Received: %s\n", buffer[:n])
		var msg string
		_, err = fmt.Scan(msg)
		if err != nil {
			continue
		}
		conn.Write([]byte(msg))
		var m []byte
		_, _ = fmt.Scan(m)
		fmt.Println(m)
	}

	// _, _ = conn.Write([]byte("recived thanks ."))
}
