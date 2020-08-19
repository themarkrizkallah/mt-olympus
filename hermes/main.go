package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"hermes/database"
	"hermes/env"
)

func main() {
	env.Init()

	log.Println("Setting up db...")

	// Init DB
	if _, err := database.Init("disable"); err != nil {
		log.Fatalln("Error setting up db:", err)
	}

	// Initialize messaging hub
	log.Println("Setting up hub...")
	hub := newHub()
	go hub.run()

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
}
