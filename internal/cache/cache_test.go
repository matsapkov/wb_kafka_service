package cache

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/matsapkov/wb_kafka_service/internal/models"
)

func TestCacheExpiresEntries(t *testing.T) {
	c := NewCache(10, 20*time.Millisecond, 5*time.Millisecond)
	t.Cleanup(c.Close)

	c.Set("order-1", models.Order{OrderUID: "order-1", Payload: json.RawMessage(`{"ok":true}`)})

	if _, ok := c.Get("order-1"); !ok {
		t.Fatalf("expected value in cache before ttl")
	}

	time.Sleep(30 * time.Millisecond)

	if _, ok := c.Get("order-1"); ok {
		t.Fatalf("expected value to be expired")
	}
}
