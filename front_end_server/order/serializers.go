package order

import "time"

type Payload struct {
	ID     string `json:"id"`
	Amount uint64 `json:"amount"`
	Price  uint64 `json:"price"`
	Side   bool   `json:"side"`
}

// Parse Parses a UserPayload into a User
func (payload *Payload) Parse() Order {
	return Order{
		ID:        payload.ID,
		Amount:    payload.Amount,
		Price:     payload.Price,
		Side:      payload.Side,
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
}
