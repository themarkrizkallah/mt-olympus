package main

import (
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"

	"grpc_server/env"
	"grpc_server/kafka"
	pb "grpc_server/proto"
	"grpc_server/server"
)

func startServer() {
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

	// Setup Kafka Producer
	kafka.CreateProducer()
	//kafka.CreateConsumer()

	listener, err := net.Listen(env.Network, ":" + env.Port)
	if err != nil {
		panic(err)
	}

	log.Println("Listening...")

	if err := grpcServer.Serve(listener); err != nil {
		panic(err)
	}
}

func main() {
	env.Init()
	startServer()
}
