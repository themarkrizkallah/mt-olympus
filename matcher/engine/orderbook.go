package engine

var orderbook OrderBook

// BuyOrders: sorted in ascending order
// SellOrders: sorted in descending order
type OrderBook struct {
	BuyOrders  []Order
	SellOrders []Order
}

func (ob *OrderBook) addBuyOrder(order Order) {
	n := len(ob.BuyOrders)
	i := n

	for i = n - 1; i >= 0; i-- {
		if ob.BuyOrders[i].Price < order.Price {
			break
		}
	}

	i++
	ob.BuyOrders = append(ob.BuyOrders, order)

	if i <= n - 1{
		copy(ob.BuyOrders[i+1:], ob.BuyOrders[i:])
		ob.BuyOrders[i] = order
	}
}

func (ob *OrderBook) addSellOrder(order Order) {
	n := len(ob.SellOrders)
	i := n

	for i = n - 1; i >= 0; i-- {
		if ob.SellOrders[i].Price > order.Price {
			break
		}
	}

	i++
	ob.SellOrders = append(ob.SellOrders, order)

	if i <= n - 1{
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

func InitializeOrderBook(capacity uint64){
	orderbook = OrderBook{
		BuyOrders:  make([]Order, 0, capacity),
		SellOrders: make([]Order, 0, capacity),
	}
}

func GetOrderBook() *OrderBook{
	return &orderbook
}