package main

import (
	"context"
	"fmt"
	"github.com/nightlord189/tcp-pow-go/internal/pkg/cache"
	"github.com/nightlord189/tcp-pow-go/internal/pkg/clock"
	"github.com/nightlord189/tcp-pow-go/internal/pkg/config"
	"github.com/nightlord189/tcp-pow-go/internal/server"
	"math/rand"
	"time"
)

func main() {
	fmt.Println("start server")

	// loading config from file and env
	configInst, err := config.Load("config/config.json")
	if err != nil {
		fmt.Println("error load config:", err)
		return
	}

	// init context to pass config down
	ctx := context.Background()
	ctx = context.WithValue(ctx, "config", configInst)
	ctx = context.WithValue(ctx, "clock", clock.SystemClock{})

	cacheInst, err := cache.InitRedisCache(ctx, configInst.CacheHost, configInst.CachePort)
	if err != nil {
		fmt.Println("error init cache:", err)
		return
	}
	ctx = context.WithValue(ctx, "cache", cacheInst)

	// seed random generator to randomize order of quotes
	rand.Seed(time.Now().UnixNano())

	// run server
	serverAddress := fmt.Sprintf("%s:%d", configInst.ServerHost, configInst.ServerPort)
	err = server.Run(ctx, serverAddress)
	if err != nil {
		fmt.Println("server error:", err)
	}
}
