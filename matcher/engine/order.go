package engine

import (
	"time"

	"github.com/golang/protobuf/proto"

	pb "matcher/proto"
)

type Order struct {
	UserId     string
	OrderId    string
	Amount     uint64
	Price      uint64
	Side       pb.Side
	Type       pb.Type
	CreatedAt  time.Time
}

func ProtoToOrder(msg []byte) (Order, error) {
	orderRequest := &pb.OrderRequest{}

	err := proto.Unmarshal(msg, orderRequest)
	if err != nil {
		return Order{}, err
	}

	order := Order{
		UserId:     orderRequest.GetUserId(),
		OrderId:    orderRequest.GetOrderId(),
		Amount:     orderRequest.GetAmount(),
		Price:      orderRequest.GetPrice(),
		Side:       orderRequest.GetSide(),
		Type:       orderRequest.GetType(),
	}

	return order, nil
}
