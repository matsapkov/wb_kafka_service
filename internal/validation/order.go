package validation

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/matsapkov/wb_kafka_service/internal/models"
)

type OrderValidator interface {
	Validate(order models.Order) error
}

type Validator struct{}

type orderPayload struct {
	Delivery deliveryPayload `json:"delivery"`
	Payment  paymentPayload  `json:"payment"`
	Items    []itemPayload   `json:"items"`
}

type deliveryPayload struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	City    string `json:"city"`
	Address string `json:"address"`
	Email   string `json:"email"`
}

type paymentPayload struct {
	Transaction string `json:"transaction"`
	Currency    string `json:"currency"`
	Amount      int    `json:"amount"`
}

type itemPayload struct {
	TrackNumber string `json:"track_number"`
	Name        string `json:"name"`
	Price       int    `json:"price"`
}

func NewOrderValidator() *Validator {
	return &Validator{}
}

func (v *Validator) Validate(order models.Order) error {
	if strings.TrimSpace(order.OrderUID) == "" {
		return errors.New("order_uid is required")
	}
	if strings.TrimSpace(order.TrackNumber) == "" {
		return errors.New("track_number is required")
	}
	if len(order.Payload) == 0 {
		return errors.New("payload is required")
	}

	var payload orderPayload
	if err := json.Unmarshal(order.Payload, &payload); err != nil {
		return fmt.Errorf("payload parse: %w", err)
	}

	if strings.TrimSpace(payload.Delivery.Name) == "" || strings.TrimSpace(payload.Delivery.Email) == "" {
		return errors.New("delivery.name and delivery.email are required")
	}
	if strings.TrimSpace(payload.Payment.Transaction) == "" || payload.Payment.Amount <= 0 {
		return errors.New("payment.transaction and payment.amount must be valid")
	}
	if len(payload.Items) == 0 {
		return errors.New("at least one item is required")
	}
	for i, item := range payload.Items {
		if strings.TrimSpace(item.Name) == "" || item.Price <= 0 {
			return fmt.Errorf("item[%d] must contain valid name and price", i)
		}
	}

	return nil
}
