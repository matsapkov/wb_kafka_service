package usecase

import (
	"github.com/matsapkov/wb_kafka_service/internal/cache"
	"github.com/matsapkov/wb_kafka_service/internal/metrics"
	"github.com/matsapkov/wb_kafka_service/internal/repository/orders"
)

type Usecase struct {
	orderRepo    orders.OrdersRepository
	cacheStorage cache.CacheStorage
	metrics      *metrics.Metrics
}

func NewOrderUsecase(orderRepo orders.OrdersRepository, cacheStorage cache.CacheStorage, m *metrics.Metrics) *Usecase {
	return &Usecase{
		orderRepo:    orderRepo,
		cacheStorage: cacheStorage,
		metrics:      m,
	}
}
