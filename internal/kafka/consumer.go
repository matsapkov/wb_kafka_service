package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/matsapkov/wb_kafka_service/internal/cache"
	"github.com/matsapkov/wb_kafka_service/internal/config"
	"github.com/matsapkov/wb_kafka_service/internal/metrics"
	"github.com/matsapkov/wb_kafka_service/internal/models"
	"github.com/matsapkov/wb_kafka_service/internal/repository/orders"
	"github.com/matsapkov/wb_kafka_service/internal/validation"
	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader    *kafka.Reader
	repo      orders.OrdersRepository
	cache     cache.CacheStorage
	validator validation.OrderValidator
	dlq       *DLQWriter
	metrics   *metrics.Metrics
}

func NewConsumer(cfg config.KafkaConfig, repo orders.OrdersRepository, cacheStorage cache.CacheStorage, validator validation.OrderValidator, m *metrics.Metrics) (*Consumer, error) {
	if len(cfg.Brokers) == 0 {
		return nil, errors.New("kafka brokers are required")
	}

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  cfg.Brokers,
		Topic:    cfg.Topic,
		GroupID:  cfg.GroupID,
		MinBytes: cfg.MinBytes,
		MaxBytes: cfg.MaxBytes,
		MaxWait:  cfg.MaxWait,
	})

	dlq, err := NewDLQWriter(cfg.Brokers, cfg.DLQTopic)
	if err != nil {
		_ = reader.Close()
		return nil, err
	}

	return &Consumer{
		reader:    reader,
		repo:      repo,
		cache:     cacheStorage,
		validator: validator,
		dlq:       dlq,
		metrics:   m,
	}, nil
}

func (c *Consumer) Start(ctx context.Context) {
	for {
		msg, err := c.reader.FetchMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return
			}
			log.Printf("kafka fetch error: %v", err)
			time.Sleep(time.Second)
			continue
		}

		if err := c.handleMessage(ctx, msg); err != nil {
			log.Printf("kafka message processing failed: %v", err)
			time.Sleep(time.Second)
		}
	}
}

func (c *Consumer) Close() {
	if err := c.reader.Close(); err != nil {
		log.Printf("kafka reader close error: %v", err)
	}
	if err := c.dlq.Close(); err != nil {
		log.Printf("kafka dlq close error: %v", err)
	}
}

func (c *Consumer) handleMessage(ctx context.Context, msg kafka.Message) error {
	var order models.Order
	if err := json.Unmarshal(msg.Value, &order); err != nil {
		return c.rejectAndCommit(ctx, msg, fmt.Sprintf("unmarshal: %v", err))
	}

	if err := c.validator.Validate(order); err != nil {
		return c.rejectAndCommit(ctx, msg, fmt.Sprintf("validation: %v", err))
	}

	if err := c.repo.SaveOrder(ctx, order.OrderUID, order.TrackNumber, order.DateCreated, order.UpdatedAt, order.Payload); err != nil {
		return c.rejectAndCommit(ctx, msg, fmt.Sprintf("db save: %v", err))
	}

	c.cache.Set(order.OrderUID, order)
	if c.metrics != nil {
		c.metrics.ProcessedMessages.Inc()
	}

	if err := c.reader.CommitMessages(ctx, msg); err != nil {
		return fmt.Errorf("commit processed message: %w", err)
	}

	return nil
}

func (c *Consumer) rejectAndCommit(ctx context.Context, msg kafka.Message, reason string) error {
	if c.metrics != nil {
		c.metrics.InvalidMessages.Inc()
	}

	if err := c.dlq.Publish(ctx, msg, reason); err != nil {
		return fmt.Errorf("publish to dlq: %w", err)
	}
	if c.metrics != nil {
		c.metrics.DLQPublished.Inc()
	}

	if err := c.reader.CommitMessages(ctx, msg); err != nil {
		return fmt.Errorf("commit rejected message: %w", err)
	}

	return nil
}
