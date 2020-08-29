package order

import (
	pb "apollo/proto"
	"time"

	"github.com/golang/protobuf/ptypes"
)

type ConfJSON struct {
	Id     string    `json:"order_id"`
	Amount int64     `json:"amount"`
	Price  int64     `json:"price"`
	Side   string    `json:"side"`
	Type   string    `json:"type"`
	Status string    `json:"status"`
	Time   time.Time `json:"time"`
}

func newConf(conf pb.OrderConf) ConfJSON {
	ts, _ := ptypes.Timestamp(conf.GetCreatedAt())

	return ConfJSON{
		Id:     conf.GetOrderId(),
		Amount: conf.GetAmount(),
		Price:  conf.GetPrice(),
		Side:   pb.Side_name[(int32)(conf.GetSide())],
		Type:   pb.Type_name[(int32)(conf.GetType())],
		Status: conf.GetStatus(),
		Time:   ts,
	}
}
