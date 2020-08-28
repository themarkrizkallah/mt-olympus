package main

import (
	"context"
	"log"
	"sync"
)

// Hub maintains the set of active clients and broadcasts messages to the clients.
type Hub struct {
	// Registered clients.
	clients ClientChannelMap

	// Channel manager
	chanManager *ChannelManager

	// Inbound messages from the clients.
	subscribe chan SubscribeRequest

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func newHub() *Hub {
	hub := &Hub{
		clients:     make(ClientChannelMap),
		chanManager: newChannelManager(),
		subscribe:   make(chan SubscribeRequest),
		register:    make(chan *Client),
		unregister:  make(chan *Client),
	}

	return hub
}

func (h *Hub) registerClient(c *Client) {
	// Check that client has not already been registered
	if _, ok := h.clients[c]; ok {
		log.Fatalln("Hub - Fatal error, already registered client")
	}

	h.clients[c] = make(map[string]bool)
}

func (h *Hub) unregisterClient(c *Client) {
	// Check that client has not already been unregistered
	if _, ok := h.clients[c]; ok {
		// Unregister client from channels its subscribed to
		for chanName := range h.clients[c] {
			h.chanManager.unregisterClient(c, chanName)
		}
		close(c.send)

	} else {
		log.Fatalln("Hub - Fatal error, already unregistered client")
	}
}

func (h *Hub) run(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	defer h.cleanup()

	log.Println("Hub - running...")

	wg.Add(1)
	go h.chanManager.run(ctx, wg)

	for {
		select {
		// New connection
		case client := <-h.register:
			log.Println("Hub - Registering client")
			h.registerClient(client)
			client.message(newConfirmationMessage("Successfully connected"))

		// Close connection
		case client := <-h.unregister:
			log.Println("Hub - Unregistering client")
			h.unregisterClient(client)

		// Subscribe request
		case subRequest := <-h.subscribe:
			h.chanManager.subscribeRequest(subRequest.Client, subRequest.SubMsg)

		// Context cancelled
		case <-ctx.Done():
			log.Println("Hub - context cancelled")
			break
		}
	}
}

func (h *Hub) cleanup() {
	close(h.subscribe)
	close(h.register)
	close(h.unregister)

	log.Println("Hub - cleanup complete")
}
