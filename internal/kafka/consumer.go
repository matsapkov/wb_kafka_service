package kafka

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/matsapkov/wb_kafka_service/internal/cache"
	"github.com/matsapkov/wb_kafka_service/internal/models"
	"github.com/matsapkov/wb_kafka_service/internal/repository/orders"

	"github.com/segmentio/kafka-go"
)

type KafkaConsumer struct {
	reader *kafka.Reader
	repo   orders.OrdersRepository
	cache  cache.CacheStorage
}

func NewKafkaConsumer(brokers []string, topic string, repo orders.OrdersRepository, cache cache.CacheStorage) (*KafkaConsumer, error) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		Topic:    topic,
		GroupID:  "order-service-group",
		MinBytes: 10e3,
		MaxBytes: 10e6,
		MaxWait:  time.Second * 1,
	})

	return &KafkaConsumer{
		reader: reader,
		repo:   repo,
		cache:  cache,
	}, nil
}

func (kc *KafkaConsumer) StartConsuming(ctx context.Context) {
	defer kc.reader.Close()

	for {
		select {
		case <-ctx.Done():
			log.Println("Kafka consumer stopped gracefully")
			return
		default:
			msg, err := kc.reader.FetchMessage(ctx)
			if err != nil {
				log.Printf("Fetch error: %v (will retry after 5s)", err)
				time.Sleep(5 * time.Second)
				continue
			}

			log.Printf("Received message at offset %d: %s", msg.Offset, string(msg.Value))

			var order models.Order
			if err := json.Unmarshal(msg.Value, &order); err != nil {
				log.Printf("Unmarshal error: %v → skip", err)
				_ = kc.reader.CommitMessages(ctx, msg)
				continue
			}

			err = kc.repo.SaveOrder(ctx, order.OrderUID, order.TrackNumber,
				order.DateCreated, order.UpdatedAt, order.Payload)
			if err != nil {
				log.Printf("Save error: %v → will retry on next read", err)
				continue
			}

			kc.cache.Set(order.OrderUID, order)

			if err := kc.reader.CommitMessages(ctx, msg); err != nil {
				log.Printf("Commit error: %v", err)
			} else {
				log.Printf("Processed and committed order %s", order.OrderUID)
			}
		}
	}
}
