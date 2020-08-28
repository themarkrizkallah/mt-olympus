package main

import (
	"context"
	"log"
	"sync"

	"github.com/Shopify/sarama"

	"hermes/database"
	"hermes/env"
	"hermes/proto"
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

type ChannelName = string

var (
	acceptedChanNames = []ChannelName{
		"heartbeat",
		"status",
		"ticker",
		"level2",
		"user",
		"matches",
		"full",
	}
	acceptedProductIDs []string
)

type Channel interface {
	init()
	subscribeClient(client *Client, productIDs []string)
	unsubscribeClient(client *Client, productIDs []string)
	unsubscribeClientFromAll(client *Client)
	broadcast(productID string, msg interface{})
}

func newChannel(chanName string) Channel {
	var (
		err     error
		channel Channel
	)

	// Need to fetch the product IDs from the database first
	// Cache product ids
	if len(acceptedProductIDs) == 0 {
		acceptedProductIDs, err = database.GetProductIDs()
		if err != nil {
			log.Fatalln("Failed setting up channel,", err)
		}
	}

	switch chanName {
	case "heartbeat":
		channel = &HeartbeatChannel{}
	case "status":
		channel = &StatusChannel{}
	case "ticker":
		channel = &TickerChannel{}
	case "level2":
		channel = &Level2Channel{}
	case "user":
		channel = &UserChannel{}
	case "matches":
		channel = &MatchesChannel{}
	case "full":
		channel = &FullChannel{}
	default:
		channel = &DefaultChannel{}
	}

	return channel
}

// Operations on the ChannelManager are not thread safe
type ChannelManager struct {
	// Channel to Product to Client map
	channels map[ChannelName]Channel

	// Kafka related fields
	consumer Consumer
	client   sarama.ConsumerGroup
	consumerWg *sync.WaitGroup

	// OrderRequest channel
	orderRequestChan chan proto.OrderRequest

	// OrderConf channel
	orderConfChan chan proto.OrderConf

	// TradeMessage channel
	tradeMsgChan chan proto.TradeMessage
}

func newChannelManager() *ChannelManager {
	log.Println("ChannelManager - setting up consumer group...")
	consumer, client := newConsumerGroup(env.GetKafkaBroker())

	cHub := &ChannelManager{
		make(map[ChannelName]Channel, len(acceptedChanNames)),
		consumer,
		client,
		&sync.WaitGroup{},
		make(chan proto.OrderRequest),
		make(chan proto.OrderConf),
		make(chan proto.TradeMessage),
	}

	// Setup channels
	for _, chanName := range acceptedChanNames {
		channel := newChannel(chanName)
		channel.init()
		cHub.channels[chanName] = channel
	}

	// Attach channels to consumer
	consumer.orderRequestChan = cHub.orderRequestChan
	consumer.orderConfChan = cHub.orderConfChan
	consumer.tradeMsgChan = cHub.tradeMsgChan

	return cHub
}

func (cm *ChannelManager) run(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	defer cm.cleanup()

	cm.consumerWg.Add(1)
	go cm.consumer.run([]string{"order.conf"}, cm.client, cm.consumerWg, ctx)

	// Await till the consumer has been set up
	<-cm.consumer.ready
	log.Println("ChannelManager - consumer up and running!")

	// TODO: Begin receiving messages
	for {
		select {
		case orderRequest := <- cm.orderRequestChan:
			log.Println("ChannelManager - received orderRequest", orderRequest)

		case orderConf := <- cm.orderConfChan:
			log.Println("ChannelManager - received orderConf", orderConf)

		case tradeMsg := <- cm.tradeMsgChan:
			log.Println("ChannelManager - received tradeMsg", tradeMsg)

		// Context cancelled
		case <-ctx.Done():
			log.Println("ChannelManager - context cancelled")
			break
		}
	}
}

func (cm *ChannelManager) cleanup(){
	cm.consumerWg.Wait()
	if err := cm.client.Close(); err != nil {
		log.Panicf("ChannelManager - error closing client: %s", err)
	}

	log.Println("ChannelManager - cleanup complete")
}

func (cm *ChannelManager) subscribeRequest(c *Client, subMsg SubscribeMessage) {
	for _, msg := range subMsg.ChannelMsgs {
		channel := cm.channels[msg.Name]

		if subMsg.MsgType == "subscribe" {
			log.Println("Hub - Subscribing client...")
			channel.subscribeClient(c, msg.ProductIDs)
		} else {
			log.Println("Hub - Unsubscribing client...")
			channel.unsubscribeClient(c, msg.ProductIDs)
		}
	}
}

func (cm *ChannelManager) unregisterClient(c *Client, chanName ChannelName) {
	// Unregister client from channel chanName
	if channel, ok := cm.channels[chanName]; ok {
		channel.unsubscribeClientFromAll(c)
	} else {
		log.Fatalln("ChannelManager - Fatal error, invalid channel name", chanName)
	}
}
