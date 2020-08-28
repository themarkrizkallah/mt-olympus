package engine

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/golang/protobuf/ptypes"

	"matcher/database"
	pb "matcher/proto"
	"matcher/types"
)

const capacity = 100

// BuyOrders: sorted in ascending order
// SellOrders: sorted in descending order
type OrderBook struct {
	Base, Quote     string
	BaseId, QuoteId string
	ProductId       string
	BuyOrders       []types.Order
	SellOrders      []types.Order
}

// Process an order and return the trades generated before adding the remaining amount to the market
func (ob *OrderBook) Process(order types.Order) (pb.OrderConf, []pb.TradeMessage) {
	var (
		account      types.Account
		orderConf    pb.OrderConf
		orderUpdates []types.OrderUpdate
		trades       []types.Trade
		tx           *sql.Tx
		err          error
	)
	order.CreatedAt = time.Now()

	tx, err = database.GetDB().BeginTx(
		context.Background(),
		&sql.TxOptions{
			Isolation: sql.LevelSerializable,
			ReadOnly:  false,
		},
	)
	if err != nil {
		log.Fatalln("Error beginning trade transaction:", err)
	}

	if order.Side == pb.Side_BUY {
		if account, err = database.GetAccount(tx, order.UserId, ob.QuoteId); err != nil {
			log.Fatalln("Error reading user account:", err)
		}

		// Reject order if user's available balance isn't enough to fund the buy order
		if account.AvailableBalance() < order.Amount*order.Price {
			ts, _ := ptypes.TimestampProto(order.CreatedAt)
			orderConf = pb.OrderConf{
				UserId:    order.UserId,
				OrderId:   order.OrderId,
				Amount:    order.Amount,
				Price:     order.Price,
				Side:      order.Side,
				Type:      order.Type,
				Status:    "Rejected",
				CreatedAt: ts,
			}

			return orderConf, []pb.TradeMessage{}
		}

		if err := database.PutHold(tx, order.UserId, ob.QuoteId, order.Amount*order.Price); err != nil {
			log.Fatalln("Error putting hold on account for buy order:", err)
		}
		orderConf, trades, orderUpdates = ob.processLimitBuy(order)
	} else {
		if account, err = database.GetAccount(tx, order.UserId, ob.BaseId); err != nil {
			log.Fatalln("Error reading user account:", err)
		}

		// Reject order if user's available balance isn't enough to fund the sell order
		if account.AvailableBalance() < order.Amount {
			ts, _ := ptypes.TimestampProto(order.CreatedAt)
			orderConf = pb.OrderConf{
				UserId:    order.UserId,
				OrderId:   order.OrderId,
				Amount:    order.Amount,
				Price:     order.Price,
				Side:      order.Side,
				Type:      order.Type,
				Status:    "Rejected",
				CreatedAt: ts,
			}
			return orderConf, []pb.TradeMessage{}
		}

		if err := database.PutHold(tx, order.UserId, ob.BaseId, order.Amount); err != nil {
			log.Fatalln("Error putting hold on account for sell order:", err)
		}
		orderConf, trades, orderUpdates = ob.processLimitSell(order)
	}

	// Insert the order in the database
	if err = database.InsertOrder(tx, order, orderConf.GetStatus(), ob.ProductId); err != nil {
		log.Fatalln("Error inserting order", err)
	}

	// Update the relevant orderChan
	for _, orderUpdate := range orderUpdates {
		if err = database.UpdateOrderStatus(tx, orderUpdate); err != nil {
			log.Fatalln("Error updating order value", err)
		}
	}

	// Reflect the value transfer in the DB and remove holds (if appropriate)
	for _, trade := range trades {
		if err = database.TransferValue(tx, &trade, ob.BaseId, ob.QuoteId); err != nil {
			log.Fatalln("Error transferring value", err)
		}
	}

	if err := tx.Commit(); err != nil {
		log.Fatalln("Error committing transaction:", err)
	}

	return orderConf, types.TradesToMessages(trades)
}

func (ob *OrderBook) addBuyOrder(order types.Order) {
	n := len(ob.BuyOrders)
	i := n

	for i = n - 1; i >= 0; i-- {
		if ob.BuyOrders[i].Price < order.Price {
			break
		}
	}

	i++
	ob.BuyOrders = append(ob.BuyOrders, order)

	if i <= n-1 {
		copy(ob.BuyOrders[i+1:], ob.BuyOrders[i:])
		ob.BuyOrders[i] = order
	}
}

func (ob *OrderBook) addSellOrder(order types.Order) {
	n := len(ob.SellOrders)
	i := n

	for i = n - 1; i >= 0; i-- {
		if ob.SellOrders[i].Price > order.Price {
			break
		}
	}

	i++
	ob.SellOrders = append(ob.SellOrders, order)

	if i <= n-1 {
		copy(ob.SellOrders[i+1:], ob.SellOrders[i:])
		ob.SellOrders[i] = order
	}
}

func (ob *OrderBook) removeBuyOrder(i int) {
	ob.BuyOrders = append(ob.BuyOrders[:i], ob.BuyOrders[i+1:]...)
}

func (ob *OrderBook) removeSellOrder(i int) {
	ob.SellOrders = append(ob.SellOrders[:i], ob.SellOrders[i+1:]...)
}

func newOrderBook(base, quote string) *OrderBook {
	product, err := database.GetProduct(base, quote)
	if err != nil {
		log.Fatalln("Error retrieving product info:", err)
	}

	return &OrderBook{
		Base:       base,
		Quote:      quote,
		BaseId:     product.BaseId,
		QuoteId:    product.QuoteId,
		ProductId:  product.Id,
		BuyOrders:  make([]types.Order, 0, capacity),
		SellOrders: make([]types.Order, 0, capacity),
	}
}
