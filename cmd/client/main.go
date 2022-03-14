package main

import (
	"context"
	"fmt"
	"github.com/nightlord189/tcp-pow-go/internal/client"
	"github.com/nightlord189/tcp-pow-go/internal/pkg/config"
)

func main() {
	fmt.Println("start client")

	// loading config from file and env
	configInst, err := config.Load("config/config.json")
	if err != nil {
		fmt.Println("error load config:", err)
		return
	}

	// init context to pass config down
	ctx := context.Background()
	ctx = context.WithValue(ctx, "config", configInst)

	address := fmt.Sprintf("%s:%d", configInst.ServerHost, configInst.ServerPort)

	// run client
	err = client.Run(ctx, address)
	if err != nil {
		fmt.Println("client error:", err)
	}
}
