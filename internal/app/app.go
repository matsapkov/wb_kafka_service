package app

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/lib/pq"
	"github.com/matsapkov/wb_kafka_service/internal/cache"
	"github.com/matsapkov/wb_kafka_service/internal/config"
	"github.com/matsapkov/wb_kafka_service/internal/handler/orders"
	"github.com/matsapkov/wb_kafka_service/internal/kafka"
	"github.com/matsapkov/wb_kafka_service/internal/metrics"
	orderRepo "github.com/matsapkov/wb_kafka_service/internal/repository/orders"
	"github.com/matsapkov/wb_kafka_service/internal/router"
	orderUsecase "github.com/matsapkov/wb_kafka_service/internal/usecase/orders/usecase"
	"github.com/matsapkov/wb_kafka_service/internal/validation"
)

type App struct {
	httpServer *http.Server
	consumer   *kafka.Consumer
	cache      *cache.Cache
	db         *sql.DB
}

func New(cfg config.Config) (*App, error) {
	if cfg.DBDSN == "" {
		return nil, fmt.Errorf("DB_DSN is required")
	}
	if len(cfg.Kafka.Brokers) == 0 {
		return nil, fmt.Errorf("KAFKA_BROKERS is required")
	}

	db, err := sql.Open("postgres", cfg.DBDSN)
	if err != nil {
		return nil, fmt.Errorf("open postgres: %w", err)
	}

	db.SetMaxOpenConns(10)
	db.SetConnMaxIdleTime(3 * time.Minute)
	db.SetConnMaxLifetime(15 * time.Minute)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping postgres: %w", err)
	}

	ordersRepository := orderRepo.NewPostgresOrders(db)
	ordersCache := cache.NewCache(cfg.CacheLimit, cfg.CacheTTL, cfg.CacheCleanupInterval)

	if err := warmupCache(context.Background(), ordersRepository, ordersCache, cfg.CacheLimit); err != nil {
		log.Printf("cache warmup warning: %v", err)
	}

	m := metrics.New()
	uc := orderUsecase.NewOrderUsecase(ordersRepository, ordersCache, m)
	h := orders.NewOrderHandler(uc)
	r := router.NewRouter(h, m)

	consumer, err := kafka.NewConsumer(cfg.Kafka, ordersRepository, ordersCache, validation.NewOrderValidator(), m)
	if err != nil {
		return nil, fmt.Errorf("create consumer: %w", err)
	}

	httpServer := &http.Server{
		Addr:         cfg.HTTPAddr,
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return &App{
		httpServer: httpServer,
		consumer:   consumer,
		cache:      ordersCache,
		db:         db,
	}, nil
}

func (a *App) Start(ctx context.Context) error {
	go a.consumer.Start(ctx)

	errCh := make(chan error, 1)
	go func() {
		if err := a.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	select {
	case <-ctx.Done():
		return a.Shutdown(context.Background())
	case err := <-errCh:
		return err
	}
}

func (a *App) Shutdown(ctx context.Context) error {
	a.consumer.Close()
	a.cache.Close()

	shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := a.httpServer.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("http shutdown: %w", err)
	}

	if err := a.db.Close(); err != nil {
		return fmt.Errorf("db close: %w", err)
	}

	return nil
}

func warmupCache(ctx context.Context, repo orderRepo.OrdersRepository, c *cache.Cache, limit int) error {
	orders, err := repo.ListRecent(ctx, limit)
	if err != nil {
		return err
	}
	c.WarmUp(orders)
	return nil
}
