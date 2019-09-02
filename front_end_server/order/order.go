package order

import pb "front_end_server/proto"

type Order struct {
	ID        string `json:"id"`
	Amount    uint64 `json:"amount"`
	Price     uint64 `json:"price"`
	Side      bool   `json:"side"`
	CreatedAt string
}

func (order *Order) ToProtoObj() pb.OrderObj {
	return pb.OrderObj{
		ID:        order.ID,
		Amount:    order.Amount,
		Price:     order.Price,
		Side:      order.Side,
		CreatedAt: order.CreatedAt,
	}
}
