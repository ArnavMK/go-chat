package main

import (
	"fmt"
	"io"
	"net"
)


func main() {

    ln, err := net.Listen("tcp", ":8080");
    fmt.Println("Started server at localhost :8080 tcp");

    messages := make(chan string);

    if err != nil {
        fmt.Println(err);
        return;
    }

    defer ln.Close();

    go handleCentralMessageChannel(messages);

    for {
        connection, err := ln.Accept();
        if err != nil {
            fmt.Println(err);
            return;
        }

        go handleConnections(connection, messages);
    }
}

func handleConnections(conn net.Conn, messages chan string) {

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
        messages <- string(buf[:n]);
    }
}

func handleCentralMessageChannel(messages chan string) {
    
    for message := range messages {
        fmt.Println("Someone said: ", message);
    }
}

