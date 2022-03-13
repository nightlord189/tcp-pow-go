package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strings"
	"time"
)

func main() {
	fmt.Println("start client")

	conn, err := net.Dial("tcp", "localhost:3333")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("connected")

	defer conn.Close()

	reader := bufio.NewReader(conn)

	for {
		if _, err = conn.Write([]byte("message1\n")); err != nil {
			fmt.Printf("failed to send the client request: %v\n", err)
		} else {
			fmt.Println("msg sent: message1")
		}

		// Waiting for the server response
		fmt.Println("waiting for new msg")
		response, err := reader.ReadString('\n')

		switch err {
		case nil:
			fmt.Println("msg received:", strings.TrimSpace(response))
		case io.EOF:
			fmt.Println("server closed the connection")
			return
		default:
			fmt.Printf("server error: %v\n", err)
			return
		}
		time.Sleep(1 * time.Second)
	}
}
