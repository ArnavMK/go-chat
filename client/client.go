package main

import (
	"fmt"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080");
	if err != nil {
		fmt.Println(err);
		return;
	}
	defer conn.Close();

	go handleIncomingMessages(conn);

	for {
		var input string;
		fmt.Scan(&input);
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
