package engine

import (
	"container/list"
	pb "matcher/proto"
	"time"
)

type PriceLevel struct {
	price      int64
	asks       *list.List
	bids       *list.List
	orderSet   map[string]*list.Element
	sellVolume int64
	buyVolume  int64
}

func NewPriceLevel(price int64) *PriceLevel {
	return &PriceLevel{
		price:      price,
		asks:       list.New(),
		bids:       list.New(),
		orderSet:   make(map[string]*list.Element),
		sellVolume: 0,
		buyVolume:  0,
	}
}

func (pl *PriceLevel) NumAsks() int {
	return pl.bids.Len()
}

func (pl *PriceLevel) NumBids() int {
	return pl.asks.Len()
}

func (pl *PriceLevel) BuyVolume() int64 {
	return pl.price * pl.buyVolume
}

func (pl *PriceLevel) SellVolume() int64 {
	return pl.price * pl.sellVolume
}

func (pl *PriceLevel) AddOrder(order Order, side pb.Side) *list.Element {
	var e *list.Element

	switch side {
	case pb.Side_BUY:
		e = pl.bids.PushBack(order)
		pl.buyVolume += order.size
	case pb.Side_SELL:
		e = pl.asks.PushBack(order)
		pl.sellVolume += order.size
	}

	pl.orderSet[order.id] = e
	return e
}

func (pl *PriceLevel) RemoveById(id string, side pb.Side) (Order, bool) {
	var order Order

	if e, ok := pl.orderSet[id]; ok {
		switch side{
		case pb.Side_BUY:
			order = pl.bids.Remove(e).(Order)
			pl.buyVolume -= order.size
		case pb.Side_SELL:
			order = pl.asks.Remove(e).(Order)
			pl.sellVolume -= order.size
		}

		delete(pl.orderSet, order.id)
		return order, true
	}

	return order, false
}

func (pl *PriceLevel) ProcessAsk(ask Order) (OrderUpdate, []Match, []OrderUpdate) {
	matches := make([]Match, 0, 1)
	updates := make([]OrderUpdate, 0, 1)
	askUpdate := OrderUpdate{id: ask.id, status: confirmStatus, size: 0}

	for e := pl.bids.Front(); e != nil; e = e.Next() {
		bid := e.Value.(Order)

		// Fill ask if possible
		if bid.size >= ask.size {
			// Decrement bid and mark ask as filled
			bid.size -= ask.size
			pl.buyVolume -= ask.size
			askUpdate.status = fillStatus
			askUpdate.size = ask.size

			// Record match
			match := Match{
				price:      pl.price,
				size:       ask.size,
				takerId:    ask.userId,
				makerId:    bid.userId,
				takerOid:   ask.id,
				makerOid:   bid.id,
				executedAt: time.Now(),
			}
			matches = append(matches, match)

			// Remove bid if filled
			if bid.size == 0 {
				updates = append(updates, OrderUpdate{id: bid.id, status: fillStatus})
			} else {
				updates = append(updates, OrderUpdate{id: bid.id, status: partialFillStatus})
			}
			break

		// Partially fill ask and continue
		} else {
			// Decrement ask, fill bid, and mark ask as partially filled
			ask.size -= bid.size
			updates = append(updates, OrderUpdate{id: bid.id, status: fillStatus})
			askUpdate.status = partialFillStatus
			askUpdate.size += bid.size

			// Record match
			match := Match{
				price:      pl.price,
				size:       bid.size,
				takerId:    ask.userId,
				makerId:    bid.userId,
				takerOid:   ask.id,
				makerOid:   bid.id,
				executedAt: time.Now(),
			}
			matches = append(matches, match)
		}
	}

	// Remove filled asks from price level
	for _, update := range updates {
		if update.status == fillStatus {
			pl.RemoveById(update.id, pb.Side_BUY)
		}
	}

	return askUpdate, matches, updates
}

func (pl *PriceLevel) ProcessBid(bid Order) (OrderUpdate, []Match, []OrderUpdate) {
	matches := make([]Match, 0, 1)
	updates := make([]OrderUpdate, 0, 1)
	bidUpdate := OrderUpdate{id: bid.id, status: confirmStatus, size: 0}

	for e := pl.asks.Front(); e != nil; e = e.Next() {
		ask := e.Value.(Order)

		// Fill bid if possible
		if ask.size >= bid.size {
			// Decrement ask and mark bid as filled
			ask.size -= bid.size
			pl.sellVolume -= bid.size
			bidUpdate.status = fillStatus
			bidUpdate.size = ask.size

			// Record match
			match := Match{
				price:      pl.price,
				size:       bid.size,
				takerId:    bid.userId,
				makerId:    ask.userId,
				takerOid:   bid.id,
				makerOid:   ask.id,
				executedAt: time.Now(),
			}
			matches = append(matches, match)

			// Remove ask if filled
			if ask.size == 0 {
				updates = append(updates, OrderUpdate{id: ask.id, status: fillStatus})
			} else {
				updates = append(updates, OrderUpdate{id: ask.id, status: partialFillStatus})
			}
			break

			// Partially fill bid and continue
		} else {
			// Decrement bid, fill ask, and mark bid as partially filled
			bid.size -= ask.size
			updates = append(updates, OrderUpdate{id: ask.id, status: fillStatus})
			bidUpdate.status = partialFillStatus
			bidUpdate.size += ask.size

			// Record match
			match := Match{
				price:      pl.price,
				size:       ask.size,
				takerId:    bid.userId,
				makerId:    ask.userId,
				takerOid:   bid.id,
				makerOid:   ask.id,
				executedAt: time.Now(),
			}
			matches = append(matches, match)
		}
	}

	// Remove filled asks from price level
	for _, update := range updates {
		if update.status == fillStatus {
			pl.RemoveById(update.id, pb.Side_SELL)
		}
	}

	return bidUpdate, matches, updates
}
