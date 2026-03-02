package kafka

import (
	"context"
	"errors"
	"fmt"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(brokers []string, topic string) (*Producer, error) {
	if len(brokers) == 0 {
		return nil, errors.New("brokers list cannot be empty")
	}
	if topic == "" {
		return nil, errors.New("topic cannot be empty")
	}

	return &Producer{
		writer: &kafka.Writer{
			Addr:         kafka.TCP(brokers...),
			Topic:        topic,
			Balancer:     &kafka.LeastBytes{},
			RequiredAcks: kafka.RequireOne,
		},
	}, nil
}

func (p *Producer) SendMessage(ctx context.Context, message []byte) error {
	if len(message) == 0 {
		return fmt.Errorf("empty message")
	}

	return p.writer.WriteMessages(ctx, kafka.Message{Value: message})
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
