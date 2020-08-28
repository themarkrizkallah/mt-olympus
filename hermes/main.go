package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/gin-gonic/gin"

	"hermes/database"
	"hermes/env"
)

func main() {
	env.Init()

	log.Println("main - setting up db...")

	// Initialize DB
	if _, err := database.Init("disable"); err != nil {
		log.Fatalln("main - error setting up db:", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	// Initialize messaging hub
	log.Println("main - setting up hub...")
	hub := newHub()

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go hub.run(ctx, wg)

	// Setup webserver
	r := gin.Default()

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	r.Use(gin.Recovery())

	// WebSocket endpoint
	r.GET("/ws", func(c *gin.Context) {
		serveWs(hub, c)
	})

	if err := r.Run(); err != nil {
		panic(err)
	}

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-ctx.Done():
		log.Println("main - terminating: context cancelled")
	case <-sigterm:
		log.Println("main - terminating: via signal")
	}
	cancel()

	wg.Wait()
}
