package server

import (
	"context"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/golang/protobuf/proto"
	"grpc_server/kafka"
	pb "grpc_server/proto"
	"log"
	"sync"
)

var ChannelMap sync.Map

type Server struct{}

type OrderResponse struct{}

func (s *Server) PlaceOrder(c context.Context, orderRequest *pb.OrderRequest) (*pb.PlaceOrderResponse, error) {
	var response * pb.PlaceOrderResponse

	log.Printf("Processing: %+v\n", orderRequest)

	data, err := proto.Marshal(orderRequest)
	if err != nil {
		log.Println("Marshaling error: ", err)
		return &pb.PlaceOrderResponse{Body: ""}, err
	}

	msg := &sarama.ProducerMessage{
		Topic: "order.request",
		Value: sarama.ByteEncoder(data),
	}

	channel := make(chan OrderResponse)
	ChannelMap.Store(orderRequest.OrderId, channel)

	_, _, err = kafka.SyncProducer.SendMessage(msg)
	if err != nil {
		response = &pb.PlaceOrderResponse{Body: err.Error()}
	} else {
		ch, _ := ChannelMap.Load(orderRequest.OrderId)

		response = &pb.PlaceOrderResponse{Body: fmt.Sprintf("+%v", <- ch.(chan OrderResponse))}
	}

	ChannelMap.Delete(orderRequest.OrderId)
	return response, err
}
