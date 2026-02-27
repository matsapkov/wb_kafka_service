package kafka

import (
	"context"
	"fmt"
	"log"

	"github.com/segmentio/kafka-go"
)

type KafkaProducer struct {
	writer *kafka.Writer
}

func NewKafkaProducer(brokers []string, topic string) (*KafkaProducer, error) {
	if len(brokers) == 0 {
		return nil, fmt.Errorf("brokers list cannot be empty")
	}

	log.Printf("Создаётся producer с brokers: %v, topic: %s", brokers, topic)

	writer := &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		Async:        false,
		RequiredAcks: kafka.RequireOne,
	}

	if err := writer.WriteMessages(context.Background(),
		kafka.Message{Value: []byte("healthcheck")},
	); err != nil {
		writer.Close()
		return nil, fmt.Errorf("cannot connect to kafka: %w", err)
	}

	return &KafkaProducer{writer: writer}, nil
}

func (kp *KafkaProducer) SendMessage(ctx context.Context, message []byte) error {
	err := kp.writer.WriteMessages(ctx,
		kafka.Message{
			Value: message,
		},
	)
	if err != nil {
		log.Printf("Ошибка отправки сообщения в Kafka: %v", err)
		return err
	}

	log.Printf("Сообщение успешно отправлено в Kafka: %s", string(message))
	return nil
}

func (kp *KafkaProducer) Close() error {
	return kp.writer.Close()
}
