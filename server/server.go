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

func (s *Server) BroadcastAllMessages() {
    
    for message := range s.messageChannel {
        fmt.Println("Someone said: ", message);
    }
}

