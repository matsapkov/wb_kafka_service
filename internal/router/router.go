package router

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/matsapkov/wb_kafka_service/internal/handler/orders"
	"github.com/matsapkov/wb_kafka_service/internal/metrics"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func NewRouter(orderHandler *orders.Handler, _ *metrics.Metrics) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", serveIndex).Methods(http.MethodGet)
	r.HandleFunc("/order/{order_id}", orderHandler.GetOrderHandler).Methods(http.MethodGet)
	r.Handle("/metrics", promhttp.Handler()).Methods(http.MethodGet)
	return r
}

func serveIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/index.html")
}
