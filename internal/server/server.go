package server

import (
	"context"
	"fmt"
	"github.com/matsapkov/wb_kafka_service/internal/cache"
	"github.com/matsapkov/wb_kafka_service/internal/config/postgres"
	"github.com/matsapkov/wb_kafka_service/internal/kafka"
	"github.com/matsapkov/wb_kafka_service/internal/router"
	"log"
	"net/http"
	"os"
	"time"

	ordersHandler "github.com/matsapkov/wb_kafka_service/internal/handler/orders"

	ordersUsecase "github.com/matsapkov/wb_kafka_service/internal/usecase/orders/usecase"

	ordersRepo "github.com/matsapkov/wb_kafka_service/internal/repository/orders"
)

const limit = 1000

type Server struct {
	server http.Server
}

func NewServer() (*Server, error) {

	PGConfig := postgres.NewPostgresConfig()
	db, err := PGConfig.PGconnect()
	if err != nil {
		return nil, err
	}

	cache := cache.NewCache(limit)

	orderRepo := ordersRepo.NewPostgresOrders(db)

	orderUsecase := ordersUsecase.NewOrderUsecase(orderRepo, cache)

	orderHandler := ordersHandler.NewOrderHandler(orderUsecase)

	ctx := context.Background()
	orders, err := orderRepo.ListRecent(ctx, limit)
	if err != nil {
		log.Printf("Error listing recent orders: %v", err)
	} else {
		cache.WarmUp(orders)
	}

	mux := router.NewRouter(orderHandler)
	brokers := []string{os.Getenv("KAFKA_BROKERS")}
	topic := os.Getenv("KAFKA_TOPIC")

	consumer, err := kafka.NewKafkaConsumer(brokers, topic, orderRepo, cache)
	if err != nil {
		return nil, fmt.Errorf("failed to create kafka consumer: %w", err)
	}

	consumerCtx, cancel := context.WithCancel(context.Background())
	go consumer.StartConsuming(consumerCtx)

	go func() {
		<-consumerCtx.Done()
		log.Println("Shutting down kafka consumer...")
		cancel()
	}()

	server := http.Server{
		Addr:         ":8081",
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	return &Server{
		server: server,
	}, nil
}

func (s *Server) StartServer() {
	err := s.server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
