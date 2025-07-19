
package server

import (
	"fmt"
	"io"
	"net"
	"sync"
)

type Server struct {
	clients []net.Conn;
	mutex sync.Mutex;
	messageChannel chan string;
}

func NewServer() *Server {
	return &Server {
		messageChannel: make(chan string),
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

		s.messageChannel <- string(buf[:n]);
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

func (s *Server) BroadcastAllMessages() {

	for message := range s.messageChannel {
		fmt.Println("Someone said: ", message);

		s.mutex.Lock();
		fmt.Println("Connected clients: ", s.clients);
		s.mutex.Unlock();
	}
}

