package usecase

import (
	"context"
	"fmt"

	"github.com/matsapkov/wb_kafka_service/internal/models"
)

func (u *Usecase) GetOrder(ctx context.Context, id string) (models.Order, error) {
	cached, ok := u.cacheStorage.Get(id)
	if ok {
		if u.metrics != nil {
			u.metrics.CacheHits.Inc()
		}
		return cached, nil
	}
	if u.metrics != nil {
		u.metrics.CacheMisses.Inc()
	}

	order, err := u.orderRepo.GetOrder(ctx, id)
	if err != nil {
		return models.Order{}, fmt.Errorf("usecase get order %s: %w", id, err)
	}

	u.cacheStorage.Set(id, order)
	return order, nil
}
