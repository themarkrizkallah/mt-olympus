package engine

import (
	"github.com/golang/protobuf/ptypes"
	pb "matcher/proto"
	"matcher/types"
)

const (
	filled = "Filled"
	partiallyFilled = "Partially Filled"
)

// Process a limit buy order
func (ob *OrderBook) processLimitBuy(order types.Order) (pb.OrderConf, []types.Trade, []types.OrderUpdate) {
	trades := make([]types.Trade, 0, 1)
	orderUpdates := make([]types.OrderUpdate, 0, 1)

	numSells := len(ob.SellOrders)
	ts, _ := ptypes.TimestampProto(order.CreatedAt)

	orderConf := pb.OrderConf{
		UserId:    order.UserId,
		OrderId:   order.OrderId,
		Amount:    order.Amount,
		Price:     order.Price,
		Side:      order.Side,
		Type:      order.Type,
		Status:    "Confirmed",
		CreatedAt: ts,
	}

	// Check if we have at least one matching order
	if numSells > 0 && ob.SellOrders[numSells-1].Price <= order.Price {
		// Traverse all orders that match
		for i := numSells - 1; i >= 0; i-- {
			sellOrder := ob.SellOrders[i]

			if ob.SellOrders[i].Price > order.Price {
				break
			}

			// Begin with the assumption that the buy order will be filled
			trade := types.Trade{
				TradeMsg: pb.TradeMessage{
					BuyerId:    order.UserId,
					SellerId:   sellOrder.UserId,
					TakerId:    order.UserId,
					MakerId:    sellOrder.UserId,
					TakerOid:   order.OrderId,
					MakerOid:   sellOrder.OrderId,
					Amount:     order.Amount,
					Price:      sellOrder.Price,
					Side:       pb.Side_BUY,
					ProductId:  ob.ProductId,
					ExecutedAt: ptypes.TimestampNow(),
				},
				Buy:  order,
				Sell: sellOrder,
			}
			orderUpdate := types.OrderUpdate{
				OrderId: sellOrder.OrderId,
				Status:  filled,
			}

			// Fill the entire order
			if sellOrder.Amount >= order.Amount {
				sellOrder.Amount -= order.Amount

				if sellOrder.Amount == 0 {
					ob.removeSellOrder(i)
				} else {
					orderUpdate.Status = partiallyFilled
					ob.SellOrders[i] = sellOrder
				}

				orderConf.Status = filled
				trades = append(trades, trade)
				orderUpdates = append(orderUpdates, orderUpdate)
				break
			}

			// Fill a partial order and continue
			trade.TradeMsg.Amount = sellOrder.Amount
			order.Amount -= sellOrder.Amount
			ob.removeSellOrder(i)
			orderConf.Status = partiallyFilled

			trades = append(trades, trade)
			orderUpdates = append(orderUpdates, orderUpdate)
		}
	}

	// Add the remaining order to the book if it isn't filled
	if order.Amount > 0 {
		ob.addBuyOrder(order)
	}
	return orderConf, trades, orderUpdates
}

// Process a limit sell order
func (ob *OrderBook) processLimitSell(order types.Order) (pb.OrderConf, []types.Trade, []types.OrderUpdate) {
	trades := make([]types.Trade, 0, 1)
	orderUpdates := make([]types.OrderUpdate, 0, 1)

	numBuys := len(ob.BuyOrders)
	ts, _ := ptypes.TimestampProto(order.CreatedAt)

	orderConf := pb.OrderConf{
		UserId:    order.UserId,
		OrderId:   order.OrderId,
		Amount:    order.Amount,
		Price:     order.Price,
		Side:      order.Side,
		Type:      order.Type,
		Status:    "Confirmed",
		CreatedAt: ts,
	}

	// Check if we have at least one matching order
	if numBuys > 0 && ob.BuyOrders[numBuys-1].Price >= order.Price {
		// Traverse all orders that match
		for i := numBuys - 1; i >= 0; i-- {
			buyOrder := ob.BuyOrders[i]

			if ob.BuyOrders[i].Price < order.Price {
				break
			}

			// Begin with the assumption that the sell order will be filled
			trade := types.Trade{
				TradeMsg: pb.TradeMessage{
					BuyerId:    buyOrder.UserId,
					SellerId:   order.UserId,
					TakerId:    order.UserId,
					MakerId:    buyOrder.UserId,
					TakerOid:   order.OrderId,
					MakerOid:   buyOrder.OrderId,
					Amount:     order.Amount,
					Price:      buyOrder.Price,
					Side:       pb.Side_SELL,
					ProductId:  ob.ProductId,
					ExecutedAt: ptypes.TimestampNow(),
				},
				Buy:  buyOrder,
				Sell: order,
			}
			orderUpdate := types.OrderUpdate{
				OrderId: buyOrder.OrderId,
				Status:  filled,
			}

			// Fill the entire order
			if buyOrder.Amount >= order.Amount {
				buyOrder.Amount -= order.Amount
				if buyOrder.Amount == 0 {
					ob.removeBuyOrder(i)
				} else {
					orderUpdate.Status = partiallyFilled
					ob.BuyOrders[i] = buyOrder
				}

				orderConf.Status = filled
				trades = append(trades, trade)
				orderUpdates = append(orderUpdates, orderUpdate)
				break

			}
			// Fill a partial order and continue
			trade.TradeMsg.Amount = buyOrder.Amount
			order.Amount -= buyOrder.Amount
			ob.removeBuyOrder(i)
			orderConf.Status = partiallyFilled
			trades = append(trades, trade)
			orderUpdates = append(orderUpdates, orderUpdate)
		}
	}

	// Add the remaining order to the book if it isn't filled
	if order.Amount > 0 {
		ob.addSellOrder(order)
	}
	return orderConf, trades, orderUpdates
}
