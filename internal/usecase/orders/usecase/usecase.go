package usecase

import (
	"github.com/matsapkov/wb_kafka_service/internal/cache"
	"github.com/matsapkov/wb_kafka_service/internal/repository/orders"
)

type Usecase struct {
	orderRepo    orders.OrdersRepository
	cacheStorage *cache.Cache
}

func NewOrderUsecase(orderRepo orders.OrdersRepository, cache *cache.Cache) *Usecase {
	return &Usecase{
		orderRepo:    orderRepo,
		cacheStorage: cache,
	}
}
