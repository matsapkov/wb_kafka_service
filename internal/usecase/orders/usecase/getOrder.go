package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/matsapkov/wb_kafka_service/internal/models"
)

func (u *Usecase) GetOrder(ctx context.Context, Id uuid.UUID) (models.Order, error) {
	order, err := u.orderRepo.GetOrder(ctx, Id)
	if err != nil {

	}

	return order, nil
}
