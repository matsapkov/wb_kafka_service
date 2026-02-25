package models

import (
	"encoding/json"
	"time"
)

type Order struct {
	OrderUID    string          `db:"order_uid"`
	TrackNumber string          `db:"track_number"`
	DateCreated time.Time       `db:"date_created"`
	Payload     json.RawMessage `db:"payload"`
	UpdatedAt   time.Time       `db:"updated_at"`
}
