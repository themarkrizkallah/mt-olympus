package server

import (
	"context"
	"log"

	"github.com/Shopify/sarama"
	"github.com/golang/protobuf/proto"

	"grpc_server/kafka"
	pb "grpc_server/proto"
)

type Server struct{}

func (s *Server) CreateOrder(c context.Context, order *pb.OrderObj) (*pb.CreateOrderResponse, error) {
	log.Printf("Processing: %+v", order)

	data, err := proto.Marshal(order)
	if err != nil {
		log.Fatal("Marshaling error: ", err)
	}

	kafka.Producer.Input() <- &sarama.ProducerMessage{
		Topic: "orders",
		Value: sarama.ByteEncoder(data),
	}

	return &pb.CreateOrderResponse{Response: "Cool"}, nil
}
