package orders

import (
	"context"

	"github.com/matsapkov/wb_kafka_service/internal/models"
)

type Usecase interface {
	GetOrder(ctx context.Context, id string) (models.Order, error)
}
