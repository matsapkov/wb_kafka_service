package cache

import (
	"sync"
	"time"

	"github.com/matsapkov/wb_kafka_service/internal/models"
)

type CacheStorage interface {
	Set(key string, value models.Order)
	Get(key string) (models.Order, bool)
	Delete(key string)
	WarmUp(orders []models.Order)
}

type entry struct {
	value     models.Order
	expiresAt time.Time
}

type Cache struct {
	storage map[string]entry
	mu      sync.RWMutex
	limit   int
	ttl     time.Duration
	stopCh  chan struct{}
}

func NewCache(limit int, ttl, cleanupInterval time.Duration) *Cache {
	c := &Cache{
		storage: make(map[string]entry, limit),
		limit:   limit,
		ttl:     ttl,
		stopCh:  make(chan struct{}),
	}

	if cleanupInterval > 0 {
		go c.cleanupExpired(cleanupInterval)
	}

	return c
}

func (c *Cache) Set(key string, value models.Order) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if len(c.storage) >= c.limit {
		for k := range c.storage {
			delete(c.storage, k)
			break
		}
	}

	c.storage[key] = entry{value: value, expiresAt: time.Now().Add(c.ttl)}
}

func (c *Cache) Get(key string) (models.Order, bool) {
	c.mu.RLock()
	item, ok := c.storage[key]
	c.mu.RUnlock()
	if !ok {
		return models.Order{}, false
	}

	if time.Now().After(item.expiresAt) {
		c.Delete(key)
		return models.Order{}, false
	}

	return item.value, true
}

func (c *Cache) Delete(key string) {
	c.mu.Lock()
	delete(c.storage, key)
	c.mu.Unlock()
}

func (c *Cache) WarmUp(orders []models.Order) {
	for _, value := range orders {
		c.Set(value.OrderUID, value)
	}
}

func (c *Cache) Close() {
	close(c.stopCh)
}

func (c *Cache) cleanupExpired(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-c.stopCh:
			return
		case <-ticker.C:
			now := time.Now()
			c.mu.Lock()
			for key, item := range c.storage {
				if now.After(item.expiresAt) {
					delete(c.storage, key)
				}
			}
			c.mu.Unlock()
		}
	}
}
