package router

import (
	"github.com/gorilla/mux"

	"github.com/matsapkov/wb_kafka_service/internal/handler/orders"
)

func NewRouter(orderHandler *orders.Handler) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/order/{order_id}", orderHandler.GetOrderHandler).Methods("GET")

	return router
}
