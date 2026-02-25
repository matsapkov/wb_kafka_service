package dto

import (
	"encoding/json"
	"time"
)

type OutputDto struct {
	OrderUID    string          `json:"order_uid"`
	TrackNumber string          `json:"track_number"`
	DateCreated time.Time       `json:"date_created"`
	Payload     json.RawMessage `json:"payload"`
	UpdatedAt   time.Time       `json:"updated_at"`
}
