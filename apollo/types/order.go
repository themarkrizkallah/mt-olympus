package types

import (
	"time"

	"github.com/golang/protobuf/ptypes"

	pb "apollo/proto"
)

type Order struct {
	UserId  string  `json:"user_id"`
	OrderId string  `json:"order_id"`
	Amount  int64   `json:"amount"`
	Price   int64   `json:"price"`
	Side    pb.Side `json:"side"`
	Type    pb.Type `json:"type"`
}

type OrderConf struct {
	OrderId   string    `json:"order_id"`
	Amount    int64     `json:"amount"`
	Price     int64     `json:"price"`
	Side      pb.Side   `json:"side"`
	Type      pb.Type   `json:"type"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

func FromProto(oc pb.OrderConf) OrderConf {
	ts, _ := ptypes.Timestamp(oc.CreatedAt)

	return OrderConf{
		OrderId:   oc.OrderId,
		Amount:    oc.Amount,
		Price:     oc.Price,
		Side:      oc.Side,
		Type:      oc.Type,
		Status:    oc.Status,
		CreatedAt: ts,
	}
}
