package main;

import (
	"fmt"
	"io"
	"net"
	"sync"
)

type Server struct {
	clients []net.Conn;
	mutex sync.Mutex;
	messageChannel chan map[net.Conn]string;
}

func NewServer() *Server {
	return &Server {
		messageChannel: make(chan map[net.Conn]string),
		mutex: sync.Mutex{},
		clients: make([]net.Conn, 0),	
	}
}

func (s *Server) HandleConnections(conn net.Conn) {

	defer conn.Close();
	buf := make([]byte, 1024);

	for {
		n, err := conn.Read(buf);

		if err == io.EOF {
			fmt.Println("Client broke the connection");
			s.RemoveConnection(conn);
			return;
		}

		if err != nil {
			fmt.Println(err);
			return;
		}

		s.messageChannel <- map[net.Conn]string{conn: string(buf[:n])};
	}
}

func (s *Server) AddClient(connection net.Conn) {
	s.mutex.Lock(); 
	s.clients = append(s.clients, connection);
	s.mutex.Unlock();
}



func (s *Server) RemoveConnection(connection net.Conn) {
	fmt.Println("Removing the connections");

	s.mutex.Lock();
	for i, c := range s.clients {
		if c == connection {
			s.clients = append(s.clients[:i], s.clients[i+1:]...);
			break;
		}
	} 
	s.mutex.Unlock();
}

func (s *Server) RelayMessagesToClients() {

	for message := range s.messageChannel {
		for sender, text := range message {
			fmt.Printf("Sender: %v    Message: %v\n", sender.RemoteAddr(), text);
			for _, client := range s.clients {
				if sender != client {
					client.Write([]byte(text));
				}
			}
		}
	}
}

