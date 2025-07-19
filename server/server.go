package main

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


func (s *Server) handleConnections(conn net.Conn) {

    defer conn.Close();
    buf := make([]byte, 1024);

    for {
        n, err := conn.Read(buf);

        if err == io.EOF {
            fmt.Println("Client broke the connection");
            return;
        }

        if err != nil {
            fmt.Println(err);
            return;
        }
        
        // each client go-routine sends the message to the central channel.
		s.messageChannel <- string(buf[:n]);
    }
}

func (s *Server) broadcastAllMessages() {
    
    for message := range s.messageChannel {
        fmt.Println("Someone said: ", message);
    }
}

func main() {

	listener, err := net.Listen("tcp", ":8080");
	if err != nil {
		fmt.Println(err);
		return;
	}

	defer listener.Close();

	server := &Server {
		clients: make([]net.Conn, 0),
		mutex: sync.Mutex{},
		messageChannel: make(chan string),
	}

	go server.broadcastAllMessages();

	for {
		connection, err := listener.Accept();
		if err != nil {
			if err == io.EOF {
				fmt.Println("Client has disconnected");
				return;
			}
			return;
		}

		go server.handleConnections(connection);
	}
}
