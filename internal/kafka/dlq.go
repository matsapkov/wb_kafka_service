package kafka

import (
	"context"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
)

type DLQWriter struct {
	writer *kafka.Writer
}

func NewDLQWriter(brokers []string, topic string) (*DLQWriter, error) {
	if len(brokers) == 0 {
		return nil, fmt.Errorf("dlq brokers are required")
	}
	if topic == "" {
		return nil, fmt.Errorf("dlq topic is required")
	}

	return &DLQWriter{
		writer: &kafka.Writer{
			Addr:         kafka.TCP(brokers...),
			Topic:        topic,
			Balancer:     &kafka.LeastBytes{},
			RequiredAcks: kafka.RequireOne,
		},
	}, nil
}

func (d *DLQWriter) Publish(ctx context.Context, msg kafka.Message, reason string) error {
	headers := append(msg.Headers,
		kafka.Header{Key: "x-error-reason", Value: []byte(reason)},
		kafka.Header{Key: "x-original-topic", Value: []byte(msg.Topic)},
		kafka.Header{Key: "x-failed-at", Value: []byte(time.Now().UTC().Format(time.RFC3339Nano))},
	)

	return d.writer.WriteMessages(ctx, kafka.Message{
		Key:     msg.Key,
		Value:   msg.Value,
		Headers: headers,
		Time:    time.Now(),
	})
}

func (d *DLQWriter) Close() error {
	return d.writer.Close()
}
