package engine

import (
	"log"

	"github.com/golang/protobuf/proto"

	pb "matcher/proto"
)

type Trade struct {
	TakerID string
	MakerID string
	Amount  uint64
	Price   uint64
}

func (trade *Trade) FromProto(msg []byte) {
	tradeMessage := &pb.TradeObj{}
	err := proto.Unmarshal(msg, tradeMessage)

	if err != nil {
		log.Fatalln("Error unmarshalling:", err)
	}

	trade.TakerID = tradeMessage.TakerID
	trade.MakerID = tradeMessage.MakerID
	trade.Amount = tradeMessage.Amount
	trade.Price = tradeMessage.Price
}

func (trade *Trade) ToProto() []byte {
	tradeMessage := &pb.TradeObj{
		TakerID: trade.TakerID,
		MakerID: trade.MakerID,
		Amount:  trade.Amount,
		Price:   trade.Price,
	}

	protoMsg, err := proto.Marshal(tradeMessage)

	if err != nil {
		log.Fatalln("Error marshalling:", err)
	}

	return protoMsg
}
