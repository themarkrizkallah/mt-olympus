package order

import (
	"time"

	"github.com/golang/protobuf/ptypes"

	pb "front_end_server/proto"
)

type Payload struct {
	UserId  string  `json:"user_id"`
	OrderId string  `json:"order_id"`
	Amount  uint64  `json:"amount"`
	Price   uint64  `json:"price"`
	Side    pb.Side `json:"side"`
	Type    pb.Type `json:"type"`
}

// ToOrderRequest converts an Order to an OrderRequest
func (o *Order) ToOrderRequest() pb.OrderRequest {
	ts, err := ptypes.TimestampProto(o.CreatedAt)
	if err != nil {
		panic(err)
	}

	return pb.OrderRequest{
		UserId:    o.UserId,
		OrderId:   o.OrderId,
		Amount:    o.Amount,
		Price:     o.Price,
		Side:      o.Side,
		Type:      o.Type,
		CreatedAt: ts,
	}
}

// Parse parses a Payload into an Order
func (p *Payload) Parse() Order {
	return Order{
		UserId:    p.UserId,
		OrderId:   p.OrderId,
		Amount:    p.Amount,
		Price:     p.Price,
		Side:      p.Side,
		Type:      p.Type,
		CreatedAt: time.Now(),
	}
}
