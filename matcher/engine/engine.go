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
)

const chanSize = 10

type Engine struct {
	orderBook *OrderBook

	// Order channel
	requests <-chan pb.OrderRequest

	// Kafka consumer related fields
	consumer   Consumer
	client     sarama.ConsumerGroup
	consumerWg *sync.WaitGroup

	// Kafka producer related fields
	producer Producer
}

func NewEngine() *Engine {
	orderBook := newOrderBook(env.Base, env.Quote)
	orderChan := make(chan pb.OrderRequest, chanSize)

	log.Println("Engine - setting up consumer group")
	consumer, client := newConsumerGroup(env.GetKafkaBroker(), getConsumptionTopic(orderBook.productId))
	consumer.requests = orderChan

	log.Println("Engine - setting up producer")
	producer := Producer{newAsyncProducer(env.GetKafkaBroker())}

	return &Engine{
		orderBook:  orderBook,
		requests:   orderChan,
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
		case request := <-e.requests:
			log.Println("Engine - request received")
			conf, trades := e.receiveOrder(request)
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

func (e *Engine) receiveOrder(request pb.OrderRequest) (pb.OrderConf, []pb.TradeMessage) {
	conf, trades := e.orderBook.Process(request) // Process the request
	if len(trades) > 0 {
		log.Printf("Engine - completed trade(s): %+v\n", trades)
	}
	return conf, trades
}

func (e *Engine) confirmOrder(conf pb.OrderConf) {
	data, err := proto.Marshal(&conf)
	if err != nil {
		log.Panicf("Engine - error marshalling data: %s", err)
	}

	log.Printf("Engine - confirming order on %s.%s", confTopic, e.orderBook.productId)
	e.producer.sendMessage(fmt.Sprintf("%s.%s", confTopic, e.orderBook.productId), data)
}

func (e *Engine) broadcastTrades(trades []pb.TradeMessage) {
	for _, trade := range trades {
		data, err := proto.Marshal(&trade)
		if err != nil {
			log.Fatalln("Error marshalling trade:", err)
		}

		log.Printf("Engine - broadcasting trades on %s.%s", tradesTopic, e.orderBook.productId)
		e.producer.sendMessage(fmt.Sprintf("%s.%s", tradesTopic, e.orderBook.productId), data)
	}
}