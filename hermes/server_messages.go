package main

import (
	"time"

	"github.com/golang/protobuf/ptypes"

	pb "hermes/proto"
)

/*
Server supported message types are:
	- "confirmation"
	- "error"
*/
const (
	confMessageType   = "confirmation"
	tickerMessageType = "ticker"
	errorMessageType  = "error"
)

type ConfirmationMessage struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

func newConfirmationMessage(msg string) ConfirmationMessage {
	return ConfirmationMessage{confMessageType, msg}
}

/*
{
    "type": "ticker",
    "trade_id": 20153558,
    "sequence": 3262786978,
    "time": "2017-09-02T17:05:49.250000Z",
    "product_id": "BTC-USD",
    "price": "4388.01000000",
    "side": "buy", // Taker side
    "last_size": "0.03000000",
    "best_bid": "4388",
    "best_ask": "4388.01"
}
*/

type TickerMessage struct {
	Type      string    `json:"type"`
	Id        string    `json:"trade_id"`
	Time      time.Time `json:"time"`
	ProductId string    `json:"product_id"`
	Price     int64     `json:"price"`
	Side      string    `json:"side"`
}

// TODO
func newTickerMessage(productId string, tradeMsg pb.TradeMessage) TickerMessage {
	ts, _ := ptypes.Timestamp(tradeMsg.GetExecutedAt())

	return TickerMessage{
		Type:      tickerMessageType,
		Id:        "TRADE-ID-TODO",
		Time:      ts,
		ProductId: productId,
		Price:     tradeMsg.GetPrice(),
		Side:      pb.Side_name[(int32)(tradeMsg.GetSide())],
	}
}
