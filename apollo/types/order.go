package types

import (
	"time"

	pb "apollo/proto"
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

type Conf struct {
	UserId      string    `json:"user_id"`
	OrderId     string    `json:"order_id"`
	CreatedAt   time.Time `json:"created_at"`
	ConfirmedAt time.Time `json:"confirmed_at"`
}
