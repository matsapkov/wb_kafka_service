package usecase

import (
	"context"
	"fmt"

	"github.com/matsapkov/wb_kafka_service/internal/models"
)

func (u *Usecase) GetOrder(ctx context.Context, Id string) (models.Order, error) {
	cached, ok := u.cacheStorage.Get(Id)
	if ok {
		return cached, nil
	}

	order, err := u.orderRepo.GetOrder(ctx, Id)
	if err != nil {
		return models.Order{}, fmt.Errorf("usecase: get order %s: %w", Id, err)
	}

	u.cacheStorage.Set(Id, order)
	return order, nil
}
