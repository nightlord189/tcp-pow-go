package main

import (
	"fmt"
	"github.com/nightlord189/tcp-pow-go/internal/server"
	"net"
	"os"
)

func main() {
	fmt.Println("start server")

	listener, err := net.Listen("tcp", "localhost:3333")
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer listener.Close()
	fmt.Println("listening...", listener.Addr())
	for {
		// Listen for an incoming connection.
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		// Handle connections in a new goroutine.
		go server.HandleConnection(conn)
	}
}
