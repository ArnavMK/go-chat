package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	setupCommunicationWithServer();
}

func setupCommunicationWithServer() {

	reader := bufio.NewReader(os.Stdin);
	fmt.Print("Enter username: ");
	user, _ := reader.ReadString('\n');
	user = strings.TrimSpace(user);
	fmt.Println("You will be called ", user);

	conn, err := net.Dial("tcp", "localhost:8080");
	if err != nil {
		fmt.Println(err);
		return;
	}
	defer conn.Close();
	conn.Write([]byte(user));

	go handleIncomingMessages(conn);

	for {
		input, _ := reader.ReadString('\n');
		input = strings.TrimSpace(input);
		conn.Write([]byte(input));
	}
}

func handleIncomingMessages(conn net.Conn) {
	
	buf := make([]byte, 1024);
	for {
		n, err := conn.Read(buf);
		if err != nil {
			fmt.Println(err);
			return;
		}
		
		fmt.Println(string(buf[:n]));
	}
}
