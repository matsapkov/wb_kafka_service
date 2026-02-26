package kafka

import (
	"context"
	"encoding/json"
	"github.com/matsapkov/wb_kafka_service/internal/cache"
	"github.com/matsapkov/wb_kafka_service/internal/models"
	"github.com/matsapkov/wb_kafka_service/internal/repository/orders"
	"log"

	"github.com/segmentio/kafka-go"
)

type KafkaConsumer struct {
	reader *kafka.Reader
	topic  string
	repo   orders.OrdersRepository
	cache  cache.Cache
}

func NewKafkaConsumer(brokers []string, topic string, repo orders.OrdersRepository) (*KafkaConsumer, error) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   topic,
		GroupID: "order-service-group",
	})

	return &KafkaConsumer{
		reader: reader,
		topic:  topic,
		repo:   repo,
	}, nil
}

func (kc *KafkaConsumer) StartConsuming() {
	for {
		msg, err := kc.reader.ReadMessage(context.Background())
		if err != nil {
			log.Fatalf("Ошибка получения сообщения: %v", err)
		}

		log.Printf("Получено сообщение: %s", string(msg.Value))

		var order models.Order
		err = json.Unmarshal(msg.Value, &order)
		if err != nil {
			log.Printf("Ошибка парсинга сообщения: %v", err)
			continue
		}

		err = kc.repo.SaveOrder(context.Background(), order.OrderUID, order.TrackNumber, order.DateCreated, order.UpdatedAt, order.Payload)
		if err != nil {
			log.Printf("Ошибка сохранения заказа в базе данных: %v", err)
			continue
		}

		kc.cache.Set(order.OrderUID, order)
	}
}
