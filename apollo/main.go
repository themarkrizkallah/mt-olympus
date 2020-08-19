package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/gin-gonic/gin"

	"apollo/accounts"
	"apollo/database"
	"apollo/env"
	"apollo/kafka"
	"apollo/order"
	"apollo/product"
	"apollo/redis"
	"apollo/users"
)


func main() {
	env.Init()

	wg := sync.WaitGroup{}
	wg.Add(1)

	// Init kafka
	go initKafka(&wg)

	// Init DB
	if _, err := database.Init("disable"); err != nil {
		log.Fatalln("Error setting up db:", err)
	}

	// Init Redis
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

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	r.Use(gin.Recovery())

	// User related endpoints
	r.POST("/signup/", users.SignUp)
	r.POST("/login/", users.Login)

	// Product related endpoints
	r.GET("/products/", product.GetProducts)

	// Order related endpoints
	orderGroup := r.Group("/orders/")
	orderGroup.Use(AuthRequired()) // Require active user session
	{
		orderGroup.POST("/", order.PostOrder)
	}

	// Account related endpoints
	accountsGroup := r.Group("/accounts")
	accountsGroup.Use(AuthRequired()) // Require active user session
	{
		accountsGroup.GET("/", accounts.GetUserAccounts)
		accountsGroup.POST("/:account_id/deposit", accounts.Deposit)
		accountsGroup.POST("/:account_id/withdraw", accounts.Withdraw)
	}

	if err := r.Run(); err != nil {
		panic(err)
	}
}
