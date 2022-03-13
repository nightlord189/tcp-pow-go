package main

import (
	"fmt"
	"github.com/nightlord189/tcp-pow-go/internal/client"
)

func main() {
	fmt.Println("start client")

	err := client.Run("localhost:3333")
	if err != nil {
		fmt.Println("client error:", err)
	}
}
