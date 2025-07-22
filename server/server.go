package main

import (
	"fmt"
	"io"
	"net"
	"sync"
)

type Server struct {
	clients map[net.Conn]Client;
	mutex sync.Mutex;
	messageChannel chan map[Client]string;
}

type Client struct {	
	conn net.Conn
	username string
}

func NewServer() *Server {
	return &Server {
		messageChannel: make(chan map[Client]string),
		mutex: sync.Mutex{},
		clients: make(map[net.Conn]Client, 0),	
	}
}

func (s *Server) HandleConnections(conn net.Conn, firstMessageCounter int) {

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
	
		if firstMessageCounter < 1 {
			s.AddClient(conn, string(buf[:n])); firstMessageCounter++;
		} else {
			client := s.GetClient(conn);
			s.messageChannel <- map[Client]string{client: string(buf[:n])};
		}
	}
}

func (s *Server) GetClient(conn net.Conn) Client{
	s.mutex.Lock(); defer s.mutex.Unlock();
	return s.clients[conn];

}

func (s *Server) AddClient(connection net.Conn, username string) {
		
	client := Client {
		conn: connection,
		username: username,
	}
	s.mutex.Lock();
	s.clients[connection] = client;
	s.mutex.Unlock();
}



func (s *Server) RemoveConnection(connection net.Conn) {
	fmt.Println("Removing the connections");
	s.mutex.Lock();
	delete(s.clients, connection);
	s.mutex.Unlock();
}

func (s *Server) RelayMessagesToClients() {

	for message := range s.messageChannel {
		for sender, text := range message {
			fmt.Printf("Sender: %v    Message: %v\n", sender.conn.RemoteAddr(), text);
			for _, client := range s.clients {
				if sender.conn != client.conn {
					text := "[" + sender.username + "]" + ": " + text;
					client.conn.Write([]byte(text));
				}
			}
		}
	}
}

