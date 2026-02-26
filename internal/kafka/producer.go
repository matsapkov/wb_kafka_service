package kafka

import (
	"context"
	"github.com/segmentio/kafka-go"
	"log"
)

type KafkaProducer struct {
	writer *kafka.Writer
}

func NewKafkaProducer(brokers []string) (*KafkaProducer, error) {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Topic:    "orders",
		Balancer: &kafka.LeastBytes{},
	}

	return &KafkaProducer{
		writer: writer,
	}, nil
}

func (kp *KafkaProducer) SendMessage(message []byte) error {
	err := kp.writer.WriteMessages(context.Background(),
		kafka.Message{
			Value: message,
		},
	)
	if err != nil {
		log.Printf("Ошибка отправки сообщения в Kafka: %v", err)
		return err
	}

	log.Println("Сообщение успешно отправлено в Kafka")
	return nil
}
