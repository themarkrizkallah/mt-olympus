package types

import "time"

type Product struct {
	Id        string    `json:"id"`
	BaseId    string    `json:"base_id"`
	QuoteId   string    `json:"quote_id"`
	BaseTick  string    `json:"base_tick"`
	QuoteTick string    `json:"quote_tick"`
	CreatedAt time.Time `json:"created_at"`
}
