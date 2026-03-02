package router

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/matsapkov/wb_kafka_service/internal/handler/orders"
	"github.com/matsapkov/wb_kafka_service/internal/models"
)

type usecaseMock struct {
	order models.Order
	err   error
}

func (m usecaseMock) GetOrder(_ context.Context, _ string) (models.Order, error) {
	if m.err != nil {
		return models.Order{}, m.err
	}
	return m.order, nil
}

func TestOrderRoute(t *testing.T) {
	payload := json.RawMessage(`{"delivery":{"name":"A","email":"a@b.c"},"payment":{"transaction":"tx","amount":1},"items":[{"name":"i","price":1}]}`)
	h := orders.NewOrderHandler(usecaseMock{order: models.Order{OrderUID: "uid", TrackNumber: "WB", DateCreated: time.Now(), UpdatedAt: time.Now(), Payload: payload}})
	r := NewRouter(h, nil)

	req := httptest.NewRequest(http.MethodGet, "/order/uid", nil)
	res := httptest.NewRecorder()
	r.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", res.Code)
	}
}
