package accounts

import "time"

type Account struct {
	Id        string    `json:"account_id"`
	UserId    string    `json:"user_id,omitempty"`
	AssetId   string    `json:"asset_id"`
	AssetTick string    `json:"tick,omitempty"`
	Balance   int64     `json:"balance"`
	Holds     int64     `json:"holds"`
	CreatedAt time.Time `json:"created_at"`
}
