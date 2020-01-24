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
		UserId:     orderRequest.UserId,
		OrderId:    orderRequest.OrderId,
		Amount:     orderRequest.Amount,
		Price:      orderRequest.Price,
		Side:       orderRequest.Side,
		Type:       orderRequest.Type,
	}

	return order, nil
}

func (o *Order) ToProto(msg []byte) error {
	orderRequest := &pb.OrderRequest{}

	err := proto.Unmarshal(msg, orderRequest)
	if err != nil {
		return err
	}

	o.UserId = orderRequest.UserId
	o.OrderId = orderRequest.OrderId
	o.Amount = orderRequest.Amount
	o.Price = orderRequest.Price
	o.Side = orderRequest.Side
	o.Type = orderRequest.Type

	return nil
}
