package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/Shopify/sarama"

	"matcher/database"
	"matcher/engine"
	"matcher/env"
	"matcher/kafka"
)

func main() {
	env.Init()

	// Init DB
	if _, err := database.Init("disable"); err != nil {
		log.Fatalln("Error setting up db:", err)
	}

	engine.InitializeOrderBook(env.OrderBookCap, env.Base, env.Quote)
	startMatcher()
}

func startMatcher() {
	var (
		producer *sarama.AsyncProducer
		consumer kafka.Consumer
		client sarama.ConsumerGroup
		err error
	)

	brokers := []string{fmt.Sprintf("%v:%v", env.KafkaHost, env.KafkaPort)}
	topics := []string{"order.request"}

	if 	producer, err = kafka.SetupProducer(brokers); err != nil {
		log.Fatalln("Unable to connect async producer to kafka server:", err)
	}
	log.Println("Async producer up and running!...")
	defer (*producer).Close()

	if consumer, client, err = kafka.SetupConsumer(brokers); err != nil {
		log.Fatal("Unable to connect consumer group to kafka server:", err)
	}
	log.Println("Connected to consumer group!...")
	defer client.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			if err := client.Consume(ctx, topics, &consumer); err != nil {
				log.Panicf("Error from consumer: %v", err)
			}
			// check if context was cancelled, signaling that the consumer should stop
			if ctx.Err() != nil {
				return
			}
			consumer.Ready = make(chan bool, 1)
		}
	}()

	<-consumer.Ready // Await till the consumer has been set up
	log.Println("Consumer up and running!...")

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-ctx.Done():
		log.Println("terminating: context cancelled")
	case <-sig:
		log.Println("terminating: via signal")
	}

	wg.Wait()
	if err = client.Close(); err != nil {
		log.Panicf("Error closing client: %v", err)
	}
}
