package main

import (
	"context"
	"fmt"
	"github.com/nightlord189/tcp-pow-go/internal/pkg/config"
	"github.com/nightlord189/tcp-pow-go/internal/server"
)

func main() {
	fmt.Println("start server")

	configInst, err := config.Load("config/config.json")
	if err != nil {
		fmt.Println("error load config:", err)
		return
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, "config", configInst)
	address := fmt.Sprintf("%s:%d", configInst.ServerHost, configInst.ServerPort)

	err = server.Run(ctx, address)
	if err != nil {
		fmt.Println("server error:", err)
	}
}
