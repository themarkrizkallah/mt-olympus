package types

import (
	"time"

	"github.com/golang/protobuf/proto"

	pb "matcher/proto"
)

type Order struct {
	UserId    string    `json:"user_id"`
	OrderId   string    `json:"order_id"`
	Amount    int64     `json:"amount"`
	Price     int64     `json:"price"`
	Side      pb.Side   `json:"side"`
	Type      pb.Type   `json:"type"`
	CreatedAt time.Time `json:"created_at"`
}

type OrderUpdate struct {
	OrderId string
	Status  string
}

func ProtoToOrder(msg []byte) (Order, error) {
	orderRequest := &pb.OrderRequest{}

	err := proto.Unmarshal(msg, orderRequest)
	if err != nil {
		return Order{}, err
	}

	order := Order{
		UserId:  orderRequest.GetUserId(),
		OrderId: orderRequest.GetOrderId(),
		Amount:  orderRequest.GetAmount(),
		Price:   orderRequest.GetPrice(),
		Side:    orderRequest.GetSide(),
		Type:    orderRequest.GetType(),
	}

	return order, nil
}
