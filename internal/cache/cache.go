package cache

import (
	"github.com/matsapkov/wb_kafka_service/internal/models"

	"sync"
)

type Cache struct {
	Storage map[string]models.Order
	mu      sync.RWMutex
}

func NewCache(limit int) *Cache {
	return &Cache{
		Storage: make(map[string]models.Order, limit),
	}
}

type CacheStorage interface {
	Set(key string, value models.Order)
	Get(key string) (models.Order, bool)
	WarmUp(orders []models.Order)
}

func (c *Cache) Set(key string, value models.Order) {
	c.mu.Lock()
	c.Storage[key] = value
	c.mu.Unlock()
}

func (c *Cache) Get(key string) (models.Order, bool) {
	c.mu.RLock()
	value, ok := c.Storage[key]
	c.mu.RUnlock()
	return value, ok
}

func (c *Cache) WarmUp(orders []models.Order) {
	c.mu.Lock()
	for _, value := range orders {
		c.Storage[value.OrderUID] = value
	}
	c.mu.Unlock()
}
