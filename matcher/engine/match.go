package engine

import (
	"github.com/golang/protobuf/ptypes"
	"time"

	pb "matcher/proto"
)

type Match struct {
	size, price int64
	takerId     string
	makerId     string
	takerOid    string
	makerOid    string
	executedAt  time.Time
}

func matchesToTrades(matches []Match, side pb.Side, productId string) []pb.TradeMessage {
	var bidUserId, askUserId string

	trades := make([]pb.TradeMessage, 0, len(matches))

	for _, match := range matches {
		ts, _ := ptypes.TimestampProto(match.executedAt)

		if side == pb.Side_BUY {
			bidUserId = match.takerId
			askUserId = match.makerId
		} else {
			bidUserId = match.makerId
			askUserId = match.takerId
		}

		trade := pb.TradeMessage{
			BuyerId:    bidUserId,
			SellerId:   askUserId,
			TakerId:    match.takerId,
			MakerId:    match.makerId,
			TakerOid:   match.takerOid,
			MakerOid:   match.makerOid,
			Amount:     match.size,
			Price:      match.price,
			Side:       side,
			ProductId:  productId,
			ExecutedAt: ts,
		}

		trades = append(trades, trade)
	}

	return trades
}