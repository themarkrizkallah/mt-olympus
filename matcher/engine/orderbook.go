package engine

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"sort"
	"time"

	"github.com/golang/protobuf/ptypes"

	"matcher/database"
	pb "matcher/proto"
	"matcher/types"
)

const (
	minPrice             = 1
	maxPrice             = 10000000
	priceLevelBufferSize = 20
)

type OrderBook struct {
	baseId, quoteId string
	productId       string
	asks, bids      []*PriceLevel
	orderPriceMap   map[string]*PriceLevel
}

func newOrderBook(base, quote string) *OrderBook {
	product, err := database.GetProduct(base, quote)
	if err != nil {
		log.Fatalln("OrderBook - error retrieving product info:", err)
	}

	orderBook := &OrderBook{
		baseId:        product.BaseId,
		quoteId:       product.QuoteId,
		productId:     product.Id,
		asks:          make([]*PriceLevel, 0, priceLevelBufferSize),
		bids:          make([]*PriceLevel, 0, priceLevelBufferSize),
		orderPriceMap: make(map[string]*PriceLevel),
	}

	for i := int64(0); i < maxPrice-minPrice+1; i += 1 {
		orderBook.bids = append(orderBook.bids, NewPriceLevel(i+minPrice))
	}

	return orderBook
}

func (ob *OrderBook) PriceLevel(price int64) *PriceLevel {
	return ob.bids[price-minPrice]
}

func (ob *OrderBook) NewPriceLevel(price int64, side pb.Side) error {
	if price < minPrice || price > maxPrice {
		return errors.New("error: price out of bounds")
	}

	switch side {
	case pb.Side_SELL:
		// Compute insertion position
		i := sort.Search(len(ob.asks), func(i int) bool {
			return (i == len(ob.asks) - 1) ||
				(i == 0 && price > ob.asks[i+1].price) ||
				(i > 0 && price < ob.asks[i-1].price && price > ob.asks[i+1].price)
		})

		ob.asks = append(ob.asks, ob.asks[len(ob.bids)-1])
		copy(ob.asks[i+1:], ob.asks[i:])
		ob.asks[i] = NewPriceLevel(price)

	case pb.Side_BUY:
		// Compute insertion position
		i := sort.Search(len(ob.bids), func(i int) bool {
			return (i == len(ob.bids) - 1) ||
				(i == 0 && price < ob.bids[i+1].price) ||
				(i > 0 && price > ob.bids[i-1].price && price < ob.bids[i+1].price)
		})

		ob.bids = append(ob.bids, ob.bids[len(ob.bids)-1])
		copy(ob.bids[i+1:], ob.bids[i:])
		ob.bids[i] = NewPriceLevel(price)
	}

	return nil
}

func (ob *OrderBook) AddOrder(order Order, price int64) {

}

func (ob *OrderBook) processLimitAsk(ask Order, askPrice int64) (OrderUpdate, []Match, []OrderUpdate) {
	const side = pb.Side_SELL

	var (
		askUpdate  OrderUpdate
		allMatches []Match
		allUpdates []OrderUpdate
	)

	if askPrice < minPrice || askPrice > maxPrice {
		log.Fatalln("OrderBook - cannot process limit ask, askPrice out of range")
	}

	askSize := ask.size
	askUpdate = OrderUpdate{id: ask.id, status: confirmStatus, size: 0}

	// Check if we have at least one matching order
	if

	return askUpdate, allMatches, allUpdates
}

//func (ob *OrderBook) processLimitAsk(ask Order, askPrice int64) (OrderUpdate, []Match, []OrderUpdate) {
//	const side = pb.Side_SELL
//
//	var (
//		askUpdate  OrderUpdate
//		allMatches []Match
//		allUpdates []OrderUpdate
//	)
//
//	if askPrice < minPrice || askPrice > maxPrice {
//		log.Fatalln("OrderBook - cannot process limit ask, askPrice out of range")
//	}
//
//	askSize := ask.size
//	askUpdate = OrderUpdate{id: ask.id, status: confirmStatus, size: 0}
//
//	// Retrieve max bids, if no bids, add ask to order book and return
//	if max, ok := ob.maxHeap.Peek(); !ok {
//		ob.PriceLevel(askPrice).AddOrder(ask, side)
//		return askUpdate, []Match{}, []OrderUpdate{}
//	}
//
//	// ask price is higher than maximum bid, add ask to order book and return
//	if askPrice > ob.maxBid {
//		ob.PriceLevel(askPrice).AddOrder(ask, side)
//		return askUpdate, []Match{}, []OrderUpdate{}
//	}
//
//	// Check price levels from the max bid price to the limit ask price
//	for price := ob.maxBid; price <= askPrice && askUpdate.size < askSize; price -= 1 {
//		if ob.PriceLevel(price).NumBids() == 0 {
//			continue
//		}
//
//		update, matches, updates := ob.bids[price-minPrice].ProcessAsk(ask)
//
//		// Update ask status, matches, and updates
//		ask.size -= update.size
//		askUpdate.Update(update)
//		allMatches = append(allMatches, matches...)
//		allUpdates = append(allUpdates, updates...)
//	}
//
//	// Partial fill, add ask to order book
//	if askUpdate.size < askSize {
//		ob.PriceLevel(askPrice).AddOrder(ask, side)
//	}
//
//	return askUpdate, allMatches, allUpdates
//}

func (ob *OrderBook) processLimitBid(bid Order, bidPrice int64) (OrderUpdate, []Match, []OrderUpdate) {
	const side = pb.Side_BUY

	var (
		bidUpdate  OrderUpdate
		allMatches []Match
		allUpdates []OrderUpdate
	)

	if bidPrice < minPrice || bidPrice > maxPrice {
		log.Fatalln("OrderBook - cannot process limit bid, bidPrice out of range")
	}

	bidSize := bid.size
	bidUpdate = OrderUpdate{id: bid.id, status: confirmStatus, size: 0}

	// bid price is lower than minimum ask price, add to order book and return
	if ob.minAsk == 0 || bidPrice < ob.minAsk {
		ob.PriceLevel(bidPrice).AddOrder(bid, side)
		return bidUpdate, []Match{}, []OrderUpdate{}
	}

	// Check price levels from min ask to the limit bid price
	for price := ob.minAsk; price >= bidPrice && bidUpdate.size < bidSize; price += 1 {
		if ob.PriceLevel(price).NumBids() == 0 {
			continue
		}

		update, matches, updates := ob.PriceLevel(price).ProcessAsk(bid)

		// Update bid status, matches, and updates
		bid.size -= update.size
		bidUpdate.Update(update)
		allMatches = append(allMatches, matches...)
		allUpdates = append(allUpdates, updates...)
	}

	// Partial fill, add bid to order book
	if bidUpdate.size < bidSize {
		ob.PriceLevel(bidPrice).AddOrder(bid, side)
	}

	return bidUpdate, allMatches, allUpdates
}

func (ob *OrderBook) Process(request pb.OrderRequest) (pb.OrderConf, []pb.TradeMessage) {
	var (
		account     types.Account
		orderUpdate OrderUpdate
		updates     []OrderUpdate
		matches     []Match
		tx          *sql.Tx
		err         error
	)

	order := Order{
		id:         request.GetOrderId(),
		userId:     request.GetUserId(),
		size:       request.GetAmount(),
		receivedAt: time.Now(),
	}
	ts, _ := ptypes.TimestampProto(order.receivedAt)
	conf := pb.OrderConf{
		UserId:    request.GetUserId(),
		OrderId:   request.GetOrderId(),
		Amount:    request.GetAmount(),
		Price:     request.GetPrice(),
		Side:      request.GetSide(),
		Type:      request.GetType(),
		Status:    "Rejected",
		CreatedAt: ts,
	}
	requestVolume := request.GetPrice() * request.GetAmount()

	tx, err = database.GetDB().BeginTx(
		context.Background(),
		&sql.TxOptions{
			Isolation: sql.LevelSerializable,
			ReadOnly:  false,
		},
	)
	if err != nil {
		log.Fatalln("OrderBook - error beginning trade transaction:", err)
	}

	switch request.GetSide() {
	case pb.Side_BUY:
		if account, err = database.GetAccount(tx, order.userId, ob.quoteId); err != nil {
			log.Fatalln("OrderBook - error reading user account:", err)
		}

		// Reject order if user's available balance isn't enough to fund the buy order
		if account.AvailableBalance() < requestVolume {
			conf.Status = rejectedStatus
			return conf, []pb.TradeMessage{}
		}

		if err := database.PutHold(tx, order.userId, ob.quoteId, requestVolume); err != nil {
			log.Fatalln("OrderBook - error putting hold on account for buy order:", err)
		}
		orderUpdate, matches, updates = ob.processLimitBid(order, request.GetPrice())

	case pb.Side_SELL:
		if account, err = database.GetAccount(tx, order.userId, ob.baseId); err != nil {
			log.Fatalln("OrderBook - error reading user account:", err)
		}

		// Reject order if user's available balance isn't enough to fund the sell order
		if account.AvailableBalance() < request.GetAmount() {
			conf.Status = rejectedStatus
			return conf, []pb.TradeMessage{}
		}

		if err := database.PutHold(tx, order.userId, ob.baseId, order.size); err != nil {
			log.Fatalln("OrderBook - error putting hold on account for sell order:", err)
		}
		orderUpdate, matches, updates = ob.processLimitAsk(order, request.GetPrice())
	}

	conf.Status = orderUpdate.status

	// Insert the order in the database
	if err = database.InsertOrder(tx, conf, ob.productId); err != nil {
		log.Fatalln("OrderBook - error inserting order", err)
	}

	// Update the relevant orders
	for _, update := range updates {
		if update.id != order.id {
			if err = database.UpdateOrderStatus(tx, update.id, update.status); err != nil {
				log.Fatalln("OrderBook - error updating order value", err)
			}
		}
	}

	// Reflect the value transfer in the databse and remove holds as applicable
	for _, match := range matches {
		var (
			bidPrice             int64
			bidUserId, askUserId string
		)

		switch request.GetSide() {
		case pb.Side_BUY:
			bidPrice = request.GetPrice()
			bidUserId = request.GetUserId()
			askUserId = match.makerId
		case pb.Side_SELL:
			bidPrice = match.price
			bidUserId = match.makerId
			askUserId = request.GetUserId()
		}

		meta := database.MatchMetadata{
			Size:      match.size,
			Price:     match.price,
			BidPrice:  bidPrice,
			BidUserId: bidUserId,
			AskUserId: askUserId,
		}

		if err = database.TransferValue(tx, meta, ob.baseId, ob.quoteId); err != nil {
			log.Fatalln("OrderBook - error transferring value", err)
		}
	}

	if err := tx.Commit(); err != nil {
		log.Fatalln("OrderBook - error committing transaction:", err)
	}

	return conf, matchesToTrades(matches, request.GetSide(), ob.productId)
}
