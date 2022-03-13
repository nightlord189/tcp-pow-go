package main

import (
	"fmt"
	"github.com/nightlord189/tcp-pow-go/internal/server"
)

func main() {
	fmt.Println("start server")

	err := server.Run("localhost:3333")
	if err != nil {
		fmt.Println("server error:", err)
	}
}
