package validation

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/matsapkov/wb_kafka_service/internal/models"
)

func TestValidator_ValidOrder(t *testing.T) {
	v := NewOrderValidator()
	order := models.Order{
		OrderUID:    "uid-1",
		TrackNumber: "WB-123",
		DateCreated: time.Now(),
		UpdatedAt:   time.Now(),
		Payload: json.RawMessage(`{
			"delivery": {"name": "John", "email": "john@example.com"},
			"payment": {"transaction": "tx-1", "currency": "RUB", "amount": 100},
			"items": [{"track_number": "WB-123", "name": "Item", "price": 100}]
		}`),
	}

	if err := v.Validate(order); err != nil {
		t.Fatalf("expected valid order, got error: %v", err)
	}
}

func TestValidator_InvalidOrder(t *testing.T) {
	v := NewOrderValidator()
	order := models.Order{OrderUID: "", TrackNumber: "WB-1", Payload: json.RawMessage(`{}`)}

	if err := v.Validate(order); err == nil {
		t.Fatalf("expected validation error")
	}
}
