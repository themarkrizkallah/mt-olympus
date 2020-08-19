package main

import (
	"hermes/database"
	"log"
)

/*
Server supported channels are:
	- "heartbeat"
    - "status"
    - "ticket"
    - "level2"
    - "user"
	- "matches"
	- "full"
*/
var (
	acceptedChanNames  = []string{"heartbeat", "status", "ticket", "level2", "user", "matches", "full"}
	acceptedProductIDs []string
)


type Channel struct {
	// Maps product_ids to set of subscribed clients
	cMap map[string] ClientSet
}

func newChannel() *Channel{
	var err error

	// Need to fetch the product IDs from the database first
	// Cache product ids
	if len(acceptedProductIDs) == 0 {
		acceptedProductIDs, err = database.GetProductIDs()
		if err != nil {
			log.Fatalln("Failed setting up channel,", err)
		}
	}

	channel := Channel{cMap: make(map[string] ClientSet)}
	for _, id := range acceptedProductIDs {
		channel.cMap[id] = make(ClientSet)
	}

	return &channel
}

func (c *Channel) subscribeClient(client *Client, productIDs []string){
	for _, id := range productIDs {
		c.cMap[id][client] = true
	}
}

func (c *Channel) unsubscribeClientFromAll(client *Client) {
	for id, _ := range c.cMap {
		delete(c.cMap[id], client)
	}
}

func (c *Channel) unsubscribeClient(client *Client, productIDs []string) {
	for _, id := range productIDs {
		delete(c.cMap[id], client)
	}
}
