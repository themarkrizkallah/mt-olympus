package engine

import (
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"

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
	ReceivedAt time.Time
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
	o.CreatedAt, _ = ptypes.Timestamp(orderRequest.CreatedAt)
	o.ReceivedAt = time.Now()

	return nil
}

func ProtoToOrder(msg []byte) (Order, error) {
	orderRequest := &pb.OrderRequest{}

	err := proto.Unmarshal(msg, orderRequest)
	if err != nil {
		return Order{}, err
	}

	ts, _ := ptypes.Timestamp(orderRequest.CreatedAt)
	order := Order{
		UserId:     orderRequest.UserId,
		OrderId:    orderRequest.OrderId,
		Amount:     orderRequest.Amount,
		Price:      orderRequest.Price,
		Side:       orderRequest.Side,
		Type:       orderRequest.Type,
		CreatedAt:  ts,
		ReceivedAt: time.Now(),
	}

	return order, nil
}
