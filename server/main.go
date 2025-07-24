package main

import (
	"fmt"
	"net"
)
func main() {

	listener, err := net.Listen("tcp", ":8080");
	if err != nil {
		fmt.Println(err);
		return;
	}
	defer listener.Close();

	server := NewServer();
	go server.RelayMessagesToClients();

	for {
		connection, err := listener.Accept();
		if err != nil {
			fmt.Println(err);
			return;
		}
	
		go server.HandleConnections(connection, false);
	}
}
