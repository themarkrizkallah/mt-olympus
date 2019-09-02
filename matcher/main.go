package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"matcher/engine"
	"matcher/env"
	"matcher/kafka"
)

func startMatcher(){
	kafka.CreateProducer()
	kafka.CreateConsumer()

	client := kafka.GetConsumerClient()
	consumer := kafka.GetConsumer()
	ctx, cancel := context.WithCancel(context.Background())
	topics := kafka.GetConsumerTopics()

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			if err := (*client).Consume(ctx, topics, consumer); err != nil {
				log.Panicf("Error from consumer: %v", err)
			}

			// check if context was cancelled, signaling that the consumer should stop
			if ctx.Err() != nil {
				return
			}
			consumer.Ready = make(chan bool)
		}
	}()

	<-consumer.Ready // Await till the consumer has been set up
	log.Println("Sarama consumer up and running!...")

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-ctx.Done():
		log.Println("terminating: context cancelled")
	case <-sigterm:
		log.Println("terminating: via signal")
	}

	cancel()
	wg.Wait()

	if err := (*client).Close(); err != nil {
		log.Panicf("Error closing client: %v", err)
	}
}


func main() {
	env.Init()
	engine.InitializeOrderBook(env.OrderBookCap)
	startMatcher()
}
