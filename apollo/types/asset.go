package types

import "time"

type Asset struct {
	Id        string    `json:"id,omitempty"`
	Name      string    `json:"name,omitempty"`
	Tick      string    `json:"tick,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}
