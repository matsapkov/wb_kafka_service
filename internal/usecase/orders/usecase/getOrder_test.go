package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/matsapkov/wb_kafka_service/internal/models"
)

type repoMock struct {
	order models.Order
	err   error
}

func (m repoMock) GetOrder(_ context.Context, _ string) (models.Order, error) {
	if m.err != nil {
		return models.Order{}, m.err
	}
	return m.order, nil
}

func (m repoMock) ListRecent(_ context.Context, _ int) ([]models.Order, error) {
	return nil, nil
}

func (m repoMock) SaveOrder(_ context.Context, _, _ string, _, _ time.Time, _ json.RawMessage) error {
	return nil
}

type cacheMock struct {
	items map[string]models.Order
}

func (c *cacheMock) Set(key string, value models.Order) { c.items[key] = value }
func (c *cacheMock) Get(key string) (models.Order, bool) {
	v, ok := c.items[key]
	return v, ok
}
func (c *cacheMock) Delete(key string)       { delete(c.items, key) }
func (c *cacheMock) WarmUp(_ []models.Order) {}

func TestGetOrder_CacheHit(t *testing.T) {
	cm := &cacheMock{items: map[string]models.Order{"o1": {OrderUID: "o1"}}}
	u := NewOrderUsecase(repoMock{}, cm, nil)

	order, err := u.GetOrder(context.Background(), "o1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if order.OrderUID != "o1" {
		t.Fatalf("unexpected order uid: %s", order.OrderUID)
	}
}

func TestGetOrder_RepoError(t *testing.T) {
	cm := &cacheMock{items: map[string]models.Order{}}
	u := NewOrderUsecase(repoMock{err: errors.New("db down")}, cm, nil)

	_, err := u.GetOrder(context.Background(), "o1")
	if err == nil {
		t.Fatalf("expected error")
	}
}
