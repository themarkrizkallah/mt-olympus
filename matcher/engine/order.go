package engine

import "time"

type Order struct {
	id         string
	userId     string
	size       int64
	receivedAt time.Time
}

const (
	confirmStatus     = "confirmed"
	fillStatus        = "filled"
	partialFillStatus = "partially filled"
	rejectedStatus    = "rejected"
)

type OrderUpdate struct {
	id     string
	status string
	size   int64
}

func (oa *OrderUpdate) Update(update OrderUpdate) {
	oa.size += update.size
	oa.status = update.status
}
