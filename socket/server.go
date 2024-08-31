package main

import (
	"fmt"
	"net"
	"time"
)

type Pair struct {
	Nikname  string
	Conn     net.Conn
	isOnline bool
}

type Server struct {
	pairs map[string]Pair
	rooms map[string]Pair
}

func NewServer() *Server {
	return &Server{
		pairs: make(map[string]Pair, 0),
		rooms: make(map[string]Pair, 0),
	}
}

func NewPair(nikename string, conn net.Conn) *Pair {
	return &Pair{
		Nikname:  nikename,
		Conn:     conn,
		isOnline: true,
	}
}

func (s *Server) JoinChat(pair1 Pair, pair2 Pair) error {
	if pair, ok := s.pairs[pair1.Nikname]; !ok || !pair.isOnline {
		return fmt.Errorf("the pair %s offline for now, or doesnt exists", pair1.Nikname)
	}

	if pair, ok := s.pairs[pair2.Nikname]; !ok || !pair.isOnline {
		return fmt.Errorf("the pair %s offline for now, or doesnt exists", pair2.Nikname)
	}

	_, ok := s.rooms[pair1.Nikname]
	if ok {
		return fmt.Errorf("%s already connected ", pair1.Nikname)
	}
	_, ok = s.rooms[pair2.Nikname]
	if ok {
		return fmt.Errorf("%s already connected ", pair2.Nikname)
	}

	s.rooms[pair1.Nikname] = pair2
	s.rooms[pair2.Nikname] = pair1

	return nil
}

func (s *Server) Logout(pair1 Pair) {
	if pair, ok := s.pairs[pair1.Nikname]; ok {
		pair.isOnline = false
	}
}

func (s *Server) Login(pair1 Pair) {

	if pair, ok := s.pairs[pair1.Nikname]; ok {
		pair.isOnline = true
	}
}

func (s *Server) SendMessage(from, to Pair, message string) error {

	if pair, ok := s.rooms[from.Nikname]; pair != to || !ok {
		return fmt.Errorf("you re not connected to that pair")
	}

	if _, ok := s.pairs[to.Nikname]; !ok {
		delete(s.rooms, from.Nikname)
		return fmt.Errorf("%s are offline ", to.Nikname)
	}

	msg := fmt.Sprintf("from [%s] : %s", from.Nikname, message)
	to.Conn.Write([]byte(msg))
	return nil
}

func (s *Server) Register(pair Pair) error {

	if _, ok := s.pairs[pair.Nikname]; ok {
		return fmt.Errorf("pair %s alredy exists", pair.Nikname)
	}

	s.pairs[pair.Nikname] = pair
	return nil
}

func main() {

	// commands := []string{
	// 	"join",
	// 	"register",
	// 	"logout",
	// 	"send",
	// }
	s := NewServer()
	// Listen for incoming connections
	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Server is listening on port 8080")

	for {

		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}
		var nikname []byte
		conn.Read(nikname)
		conn.Write([]byte("Welcome " + string(nikname)))
		go s.handleClient(conn)
	}
}

func (s *Server) handleClient(conn net.Conn) {
	defer conn.Close()

	buffer := make([]byte, 1024)

	for {
		// Read data from the client
		n, err := conn.Read(buffer)
		if err != nil {
			return
		}

		var pair2 chan Pair

		nikname := string(buffer[:n])
		pair, ok := s.pairs[nikname]
		if !ok {
			p := NewPair(nikname, conn)
			s.Register(*p)
		}
		go s.GetOnlinePair(pair2)
		s.JoinChat(pair, <-pair2)
	}

}

func (s *Server) GetOnlinePair(pair2 chan Pair) {

	for {
		for _, pair := range s.pairs {
			if pair.isOnline {
				pair2 <- pair
				return

			}
		}
		time.Sleep(time.Microsecond * 100)

	}

}
