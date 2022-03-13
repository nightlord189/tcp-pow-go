package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"time"
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
		go handleRequest(conn)
	}
}

func handleRequest(conn net.Conn) {
	fmt.Println("new client:", conn.RemoteAddr())
	defer conn.Close()

	reader := bufio.NewReader(conn)

	for {
		// Waiting for the client request
		req, err := reader.ReadString('\n')

		switch err {
		case nil:
			req := strings.TrimSpace(req)
			if req == ":QUIT" {
				fmt.Println("client requested server to close the connection so closing")
				return
			} else {
				fmt.Println("msg received:", req)
			}
		case io.EOF:
			fmt.Println("client closed the connection by terminating the process")
			return
		default:
			fmt.Printf("error: %v\n", err)
			return
		}
		time.Sleep(5 * time.Second)
		// Responding to the client request
		_, err = conn.Write([]byte("response\n"))
		if err != nil {
			fmt.Printf("failed to respond to client: %v\n", err)
		} else {
			fmt.Println("msg sent: response")
		}
	}
}
