package client

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
	"zeus/env"
	pb "zeus/proto"
)

var (
	exchangeServiceClient pb.ExchangeServiceClient
	exchangeServiceConn   *grpc.ClientConn
)

func InitExchangeService() {
	var (
		conn *grpc.ClientConn
		host = env.GrpcHostName + ":" + env.GrpcPort
		err  error
	)

	if env.Debug {
		log.Printf("Grpc Hostname: %v\n", env.GrpcHostName)
		log.Printf("Port: %v\n", env.Port)
		log.Printf("Host: %v\n", host)
		log.Printf("Secure: %v\n", env.Secure)
	}

	if env.Secure {
		creds, err := credentials.NewClientTLSFromFile(env.CertFile, "")
		if err != nil {
			log.Fatalln("Error getting creds:", err)
		}

		conn, err = grpc.Dial(host, grpc.WithTransportCredentials(creds))
	} else {
		conn, err = grpc.Dial(host, grpc.WithInsecure())
	}

	if err != nil {
		panic(err)
	}

	exchangeServiceConn = conn
	exchangeServiceClient = pb.NewExchangeServiceClient(conn)
}

func GetExchangeServiceClient() *pb.ExchangeServiceClient {
	return &exchangeServiceClient
}

func Cleanup() {
	err := exchangeServiceConn.Close()
	if err != nil {
		log.Fatalln("Could not close connection:", err)
	}
}
