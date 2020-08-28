package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"matcher/database"
	"matcher/engine"
	"matcher/env"
)

func main() {
	env.Init()

	// Initialize DB
	if _, err := database.Init("disable"); err != nil {
		log.Fatalln("Error setting up db:", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	// Initialize matching engine
	log.Println("main - setting up engine...")
	engine := engine.NewEngine()

	// Start matching engine
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go engine.Start(ctx, wg)

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