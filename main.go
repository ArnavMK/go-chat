package main 

import (
	"fmt"
	"net"
	"github.com/arnavmk/go-chat/server"
)

func main() {

	listener, err := net.Listen("tcp", ":8080");
	if err != nil {
		fmt.Println(err);
		return;
	}
	defer listener.Close();

	server := server.NewServer();
	go server.BroadcastAllMessages();

	for {
		connection, err := listener.Accept();

		if err != nil {
			fmt.Println(err);
			return;
		}

		server.AddClient(connection);
		go server.HandleConnections(connection);
	}
}
