package main

import "log"

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients ClientChannels

	// Channel to Product to Client map
	channels map[string] *Channel

	// Inbound messages from the clients.
	subscribe chan SubscribeRequest

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func (h *Hub) registerClient(c *Client) {
	// Check that client has not already been registered
	if _, ok := h.clients[c]; ok {
		log.Fatalln("Hub - Fatal error, already registered client")
	}

	h.clients[c] = make(map[string] bool)
}

func (h *Hub) unregisterClient(c *Client) {
	// Check that client has not already been unregistered
	if _, ok := h.clients[c]; ok {
		// Unregister client from channels its subscribed to
		for chanName, _ := range h.clients[c] {
			channel := h.channels[chanName]
			channel.unsubscribeClientFromAll(c)
		}
		close(c.send)

	} else {
		log.Fatalln("Hub - Fatal error, already unregistered client")
	}
}

func (h *Hub) subscribeClient(c *Client, subMsg SubscribeMessage) {
	for _, chanMsg := range subMsg.ChannelMsgs {
		channel := h.channels[chanMsg.Name]

		if subMsg.MsgType == "subscribe" {
			log.Println("Hub - Subscribing client...")
			channel.subscribeClient(c, chanMsg.ProductIDs)
		} else {
			log.Println("Hub - Unsubscribing client...")
			channel.unsubscribeClient(c, chanMsg.ProductIDs)
		}

	}
}

func newHub() *Hub {
	hub := Hub{
		subscribe:  make(chan SubscribeRequest),
		channels:   make(map[string] *Channel, len(acceptedChanNames)),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(ClientChannels),
	}

	// Initialize channels
	for _, chanName := range acceptedChanNames {
		hub.channels[chanName] = newChannel()
		log.Printf("Hub - Setup channel %s\n", chanName)
	}

	return &hub
}

func (h *Hub) run() {
	defer h.cleanup()

	for {
		select {
		// New connection
		case client := <-h.register:
			log.Println("Hub - Registering client")
			h.registerClient(client)
			client.message(newConfirmationMessage("Successfully connected"))

		// Closed connection
		case client := <-h.unregister:
			log.Println("Hub - Unregistering client")
			h.unregisterClient(client)

		case subRequest := <-h.subscribe:
			h.subscribeClient(subRequest.Client, subRequest.SubMsg)
		}
	}
}

func (h *Hub) cleanup() {
	close(h.subscribe)
	close(h.register)
	close(h.unregister)
}
