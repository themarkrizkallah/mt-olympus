package types

import pb "matcher/proto"

type Trade struct {
	TradeMsg pb.TradeMessage
	Buy      Order
	Sell     Order
}

func TradesToMessages(trades []Trade) []pb.TradeMessage {
	tradeMsgs := make([]pb.TradeMessage, 0, len(trades))

	for _, trade := range trades {
		tradeMsgs = append(tradeMsgs, trade.TradeMsg)
	}

	return tradeMsgs
}
