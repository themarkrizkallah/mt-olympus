package order

import (
	pb "zeus/proto"
	"time"
)

type Order struct {
	UserId    string    `json:"user_id"`
	OrderId   string    `json:"order_id"`
	Amount    uint64    `json:"amount"`
	Price     uint64    `json:"price"`
	Side      pb.Side   `json:"side"`
	Type      pb.Type   `json:"type"`
	CreatedAt time.Time `json:"created_at"`
}
