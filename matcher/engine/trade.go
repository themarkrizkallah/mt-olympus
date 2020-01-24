package engine

import (
	"log"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"

	pb "matcher/proto"
)

type Trade struct {
	TakerId    string
	MakerId    string
	TakerOid   string
	MakerOid   string
	Amount     uint64
	Price      uint64
	Base       string
	Quote      string
	ExecutedAt time.Time
}

func (t *Trade) ToProto() []byte {
	ts, _ := ptypes.TimestampProto(t.ExecutedAt)

	tradeMessage := &pb.TradeMessage{
		TakerId:    t.TakerId,
		MakerId:    t.MakerId,
		TakerOid:   t.TakerOid,
		MakerOid:   t.MakerOid,
		Amount:     t.Amount,
		Price:      t.Price,
		Base:       t.Base,
		Quote:      t.Quote,
		ExecutedAt: ts,
	}

	data, err := proto.Marshal(tradeMessage)
	if err != nil {
		log.Panicf("Marshaling error: %v\n", err)
	}

	return data
}
