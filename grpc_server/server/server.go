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

func (s *Server) PlaceOrder(c context.Context, orderRequest *pb.OrderRequest) (*pb.PlaceOrderResponse, error) {
	log.Printf("Processing: %+v\n", orderRequest)

	data, err := proto.Marshal(orderRequest)
	if err != nil {
		log.Println("Marshaling error: ", err)
		return &pb.PlaceOrderResponse{Body: ""}, err
	}

	msg := &sarama.ProducerMessage{
		Topic: "orders",
		Value: sarama.ByteEncoder(data),
	}

	_, _, err = kafka.Producer.SendMessage(msg)
	if err != nil {
		return &pb.PlaceOrderResponse{Body: ""}, err
	}

	return &pb.PlaceOrderResponse{Body: "Order successfully placed."}, nil
}
