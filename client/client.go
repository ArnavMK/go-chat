package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

func setupCommunicationWithServer() {

	reader := bufio.NewReader(os.Stdin);
	username := getUsername(reader);

	conn, err := net.Dial("tcp", "localhost:8080");
	if err != nil {
		if err == io.EOF {
			fmt.Print("Lost connection with the server: ");
		}
		fmt.Println(err);
		return;
	}
	defer conn.Close();
	conn.Write([]byte("USERNAME:" + username));

	go handleIncomingMessages(conn);

	for {
		fmt.Printf("[%v]: ", username);
		input, _ := reader.ReadString('\n');
		input = strings.TrimSpace(input);
		conn.Write([]byte(input));
	}
}

func getUsername(reader *bufio.Reader) string{
	fmt.Print("Enter username: ");
	user, _ := reader.ReadString('\n');
	user = strings.TrimSpace(user);
	return user;
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
