package engine

var orderbook OrderBook

// BuyOrders: sorted in ascending order
// SellOrders: sorted in descending order
type OrderBook struct {
	Base       string
	Quote      string
	BuyOrders  []Order
	SellOrders []Order
}

type OneSidedError struct {}

func (e OneSidedError) Error() string {
	return "Orderbook is one sided"
}

/*
func(ob *OrderBook) GetSpread() (uint64, error) {
	var err error

	numBuys := len(ob.BuyOrders)
	numSells := len(ob.SellOrders)
	spread := uint64(0)

	if numBuys > 0 && numSells > 0{
		spread = ob.BuyOrders[0].Price - ob.SellOrders[0].Price

		//if spread < 0 {
		//	spread *= -1
		//}

		err = nil
	} else {
		err = OneSidedError{}
	}

	return 0, err
}

 */

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

	if i <= n-1 {
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

func InitializeOrderBook(capacity uint64) {
	orderbook = OrderBook{
		BuyOrders:  make([]Order, 0, capacity),
		SellOrders: make([]Order, 0, capacity),
	}
}

func GetOrderBook() *OrderBook {
	return &orderbook
}
