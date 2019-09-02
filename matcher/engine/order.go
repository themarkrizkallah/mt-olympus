package engine

import (
	"log"

	"github.com/golang/protobuf/proto"

	pb "matcher/proto"
)

type Order struct {
	ID        string
	Amount    uint64
	Price     uint64
	Side      bool
	CreatedAt string
}

func (order *Order) FromProto(msg []byte) {
	orderMessage := &pb.OrderObj{}
	err := proto.Unmarshal(msg, orderMessage)

	if err != nil {
		log.Fatalln("Error unmarshalling:", err)
	}

	order.ID = orderMessage.ID
	order.Amount = orderMessage.Amount
	order.Price = orderMessage.Price
	order.Side = orderMessage.Side
	order.CreatedAt = orderMessage.CreatedAt
}

func (order *Order) ToProto() []byte {
	orderMessage := &pb.OrderObj{
		ID:        order.ID,
		Amount:    order.Amount,
		Price:     order.Price,
		Side:      order.Side,
		CreatedAt: order.CreatedAt,
	}

	protoMsg, err := proto.Marshal(orderMessage)

	if err != nil {
		log.Fatalln("Error marshalling:", err)
	}

	return protoMsg
}
