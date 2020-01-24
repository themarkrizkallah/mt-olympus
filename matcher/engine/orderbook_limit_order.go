package engine

import (
	"time"

	"github.com/golang/protobuf/ptypes"

	pb "matcher/proto"
)

// Process a limit buy order
func (ob *OrderBook) processLimitBuy(order Order) (pb.OrderConf, []Trade) {
	trades := make([]Trade, 0, 1)
	numSells := len(ob.SellOrders)
	ts, _ := ptypes.TimestampProto(order.CreatedAt)
	orderConf := pb.OrderConf{
		UserId:               order.UserId,
		OrderId:              order.OrderId,
		Amount:               order.Amount,
		Price:                order.Price,
		Side:                 order.Side,
		Type:                 order.Type,
		Message:              "Confirmed",
		CreatedAt:            ts,
	}

	// Check if we have at least one matching order
	if numSells > 0 && ob.SellOrders[numSells-1].Price <= order.Price {
		// Traverse all orders that match
		for i := numSells - 1; i >= 0; i-- {
			sellOrder := ob.SellOrders[i]

			if sellOrder.Price > order.Price {
				break
			}

			// Fill the entire order
			if sellOrder.Amount >= order.Amount {
				trade := Trade{
					TakerId:    order.UserId,
					MakerId:    sellOrder.UserId,
					TakerOid:   order.OrderId,
					MakerOid:   sellOrder.OrderId,
					Amount:     order.Amount,
					Price:      sellOrder.Price,
					Base:       ob.Base,
					Quote:      ob.Quote,
					ExecutedAt: time.Now(),
				}
				trades = append(trades, trade)

				sellOrder.Amount -= order.Amount
				if sellOrder.Amount == 0 {
					ob.removeSellOrder(i)
				}

				orderConf.Message = "Filled"
				return orderConf, trades
			}
			// Fill a partial order and continue
			if sellOrder.Amount < order.Amount {
				trade := Trade{
					TakerId:    order.UserId,
					MakerId:    sellOrder.UserId,
					TakerOid:   order.OrderId,
					MakerOid:   sellOrder.OrderId,
					Amount:     sellOrder.Amount,
					Price:      sellOrder.Price,
					Base:       ob.Base,
					Quote:      ob.Quote,
					ExecutedAt: time.Now(),
				}
				trades = append(trades, trade)

				order.Amount -= sellOrder.Amount
				ob.removeSellOrder(i)
				orderConf.Message = "Partially filled"
			}
		}
	}
	// finally add the remaining order to the list
	ob.addBuyOrder(order)

	return orderConf, trades
}

// Process a limit sell order
func (ob *OrderBook) processLimitSell(order Order) (pb.OrderConf, []Trade) {
	trades := make([]Trade, 0, 1)
	numBuys := len(ob.BuyOrders)
	ts, _ := ptypes.TimestampProto(order.CreatedAt)
	orderConf := pb.OrderConf{
		UserId:               order.UserId,
		OrderId:              order.OrderId,
		Amount:               order.Amount,
		Price:                order.Price,
		Side:                 order.Side,
		Type:                 order.Type,
		Message:              "Confirmed",
		CreatedAt:            ts,
	}

	// Check if we have at least one matching order
	if numBuys > 0 && ob.BuyOrders[numBuys-1].Price >= order.Price {
		// Traverse all orders that match
		for i := numBuys - 1; i >= 0; i-- {
			buyOrder := ob.BuyOrders[i]

			if buyOrder.Price < order.Price {
				break
			}

			// Fill the entire order
			if buyOrder.Amount >= order.Amount {
				trade := Trade{
					TakerId:    order.UserId,
					MakerId:    buyOrder.UserId,
					TakerOid:   order.OrderId,
					MakerOid:   buyOrder.OrderId,
					Amount:     order.Amount,
					Price:      buyOrder.Price,
					Base:       ob.Base,
					Quote:      ob.Quote,
					ExecutedAt: time.Now(),
				}
				trades = append(trades, trade)

				buyOrder.Amount -= order.Amount
				if buyOrder.Amount == 0 {
					ob.removeBuyOrder(i)
				}

				orderConf.Message = "Filled"
				return orderConf, trades
			}

			// Fill a partial order and continue
			if buyOrder.Amount < order.Amount {
				trade := Trade{
					TakerId:    order.UserId,
					MakerId:    buyOrder.UserId,
					TakerOid:   order.OrderId,
					MakerOid:   buyOrder.OrderId,
					Amount:     buyOrder.Amount,
					Price:      buyOrder.Price,
					Base:       ob.Base,
					Quote:      ob.Quote,
					ExecutedAt: time.Now(),
				}

				trades = append(trades, trade)

				order.Amount -= buyOrder.Amount
				ob.removeBuyOrder(i)
				orderConf.Message = "Partially filled"
			}
		}
	}

	// Finally add the remaining order to the list
	ob.addSellOrder(order)

	return orderConf, trades
}
