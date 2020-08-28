package engine

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/Shopify/sarama"
	"github.com/golang/protobuf/proto"

	"matcher/env"
	pb "matcher/proto"
	"matcher/types"
)

type Engine struct {
	orderbook *OrderBook

	// Order channel
	orderChan <-chan types.Order

	// Kafka consumer related fields
	consumer   Consumer
	client     sarama.ConsumerGroup
	consumerWg *sync.WaitGroup

	// Kafka producer related fields
	producer Producer
}

func NewEngine() *Engine {
	orderbook := newOrderBook(env.Base, env.Quote)
	orderChan := make(chan types.Order)

	log.Println("Engine - setting up consumer group")
	consumer, client := newConsumerGroup(env.GetKafkaBroker(), getConsumptionTopic(orderbook.ProductId))
	consumer.orderChan = orderChan

	log.Println("Engine - setting up producer")
	producer := Producer{newAsyncProducer(env.GetKafkaBroker())}

	return &Engine{
		orderbook:  orderbook,
		orderChan:  orderChan,
		consumer:   consumer,
		client:     client,
		consumerWg: &sync.WaitGroup{},
		producer:   producer,
	}
}

func (e *Engine) Start(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	defer e.cleanup()

	e.consumerWg.Add(1)
	go e.consumer.run(ctx, e.client, e.consumerWg)

	// Await till the consumer has been set up
	<-e.consumer.ready
	log.Println("Engine - consumer up and running!")

	for {
		select {
		case order := <-e.orderChan:
			log.Println("Engine - order received")
			conf, trades := e.receiveOrder(order)
			e.confirmOrder(conf)
			e.broadcastTrades(trades)

		case <-ctx.Done():
			break
		}
	}
}

func (e *Engine) cleanup() {
	e.consumerWg.Wait()
	if err := e.client.Close(); err != nil {
		log.Panicf("Engine - error closing client: %s", err)
	}

	log.Println("Engine - cleanup complete")
}

func (e *Engine) receiveOrder(order types.Order) (pb.OrderConf, []pb.TradeMessage) {
	conf, trades := e.orderbook.Process(order) // Process the order
	if len(trades) > 0 {
		log.Printf("Engine - completed trade(s): %+v\n", trades)
	}
	return conf, trades
}

func (e *Engine) confirmOrder(conf pb.OrderConf) {
	log.Printf("Engine - confirming order on %s", confTopic)

	data, err := proto.Marshal(&conf)
	if err != nil {
		log.Panicf("Engine - error marshalling data: %s", err)
	}

	e.producer.sendMessage(fmt.Sprintf("%s.%s", confTopic, e.orderbook.ProductId), data)
}

func (e *Engine) broadcastTrades(trades []pb.TradeMessage) {
	for _, trade := range trades {
		data, err := proto.Marshal(&trade)
		if err != nil {
			log.Fatalln("Error marshalling trade:", err)
		}
		e.producer.sendMessage(fmt.Sprintf("%s.%s", tradesTopic, e.orderbook.ProductId), data)
	}
}