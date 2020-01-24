package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
	"grpc_server/env"
	"grpc_server/kafka"
	pb "grpc_server/proto"
	"grpc_server/server"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	env.Init()

	wg := sync.WaitGroup{}
	wg.Add(1)

	go setupKafka(&wg)
	//startServer(&wg)

	wg.Wait()
}

func setupKafka(topWg *sync.WaitGroup) {
	defer topWg.Done()

	log.Println("Setting up kafka...")

	brokers := []string{fmt.Sprintf("%v:%v", env.KafkaHost, env.KafkaPort)}
	topics := []string{"order-conf"}

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
	log.Println("Sarama consumer up and running!...")

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

func startServer(topWg *sync.WaitGroup) {
	defer topWg.Done()

	log.Println("Setting up gRPC server...")

	// GRPC Setup
	var grpcServer *grpc.Server

	if env.Debug {
		log.Printf("Network: %v\n", env.Network)
		log.Printf("Port: %v\n", env.Port)
		log.Printf("Secure: %v\n", env.Secure)
	}

	if env.Secure {
		creds, err := credentials.NewServerTLSFromFile(env.CertFile, env.KeyFile)
		if err != nil {
			panic(err)
		}

		grpcServer = grpc.NewServer(grpc.Creds(creds))
	} else {
		grpcServer = grpc.NewServer()
	}

	pb.RegisterExchangeServiceServer(grpcServer, &server.Server{})
	reflection.Register(grpcServer)

	listener, err := net.Listen(env.Network, ":"+env.Port)
	if err != nil {
		panic(err)
	}

	log.Println("Listening...")
	if err := grpcServer.Serve(listener); err != nil {
		panic(err)
	}
}
