package main

import (
	"apollo/database"
	"apollo/redis"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/gin-gonic/gin"

	"apollo/env"
	"apollo/kafka"
	"apollo/order"
	"apollo/users"
)

func main() {
	env.Init()

	wg := sync.WaitGroup{}
	wg.Add(1)

	go initKafka(&wg)
	if _, err := database.Init("disable"); err != nil {
		log.Fatalln("Error setting up db:", err)
	}

	if _, err := redis.Init(); err != nil {
		log.Fatalln("Error initializing Redis client:", err)
	}

	startServer(&wg)

	wg.Wait()
}

func initKafka(topWg *sync.WaitGroup) {
	defer topWg.Done()

	log.Println("Setting up kafka...")

	brokers := []string{fmt.Sprintf("%v:%v", env.KafkaHost, env.KafkaPort)}
	topics := []string{"order.conf"}

	// Setup Kafka Producer
	if _, err := kafka.CreateAsyncProducer(brokers); err != nil {
		log.Fatalln("Error setting up Kafka:", err)
	}

	consumer, client, err := kafka.SetupConsumer(brokers)
	if err != nil {
		log.Fatalln("Unable to connect consumer group to kafka server:", err)
	}
	defer client.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wg := sync.WaitGroup{}
	wg.Add(2)

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
	log.Println("Sarama consumer up and running!...")

	go kafka.PipelineRequests(&wg)

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

func startServer(wg *sync.WaitGroup) {
	defer wg.Done()

	r := gin.Default()

	r.POST("/signup/", users.SignUp)
	r.POST("/login/", users.Login)

	// Endpoints that require authentication
	authorized := r.Group("/")
	authorized.Use(users.AuthRequired())
	{
		authorized.POST("/orders/", order.CreateOrder)
	}

	if err := r.Run(); err != nil {
		panic(err)
	}
}
