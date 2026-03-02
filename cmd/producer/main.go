package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/matsapkov/wb_kafka_service/internal/config"
	"github.com/matsapkov/wb_kafka_service/internal/kafka"
	"github.com/matsapkov/wb_kafka_service/internal/models"
)

func main() {
	count := flag.Int("count", 10, "number of messages")
	invalidPct := flag.Int("invalid-percent", 30, "percentage of invalid messages")
	interval := flag.Duration("interval", 200*time.Millisecond, "delay between messages")
	flag.Parse()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	producer, err := kafka.NewProducer(cfg.Kafka.Brokers, cfg.Kafka.Topic)
	if err != nil {
		log.Fatalf("create producer: %v", err)
	}
	defer func() {
		if err := producer.Close(); err != nil {
			log.Printf("producer close error: %v", err)
		}
	}()

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	ctx := context.Background()

	for i := 0; i < *count; i++ {
		valid := r.Intn(100) >= *invalidPct
		msg, err := generateMessage(r, valid)
		if err != nil {
			log.Printf("generate message error: %v", err)
			continue
		}

		if err := producer.SendMessage(ctx, msg); err != nil {
			log.Printf("send message error: %v", err)
			continue
		}
		log.Printf("sent message #%d valid=%t", i+1, valid)
		time.Sleep(*interval)
	}
}

func generateMessage(r *rand.Rand, valid bool) ([]byte, error) {
	if !valid {
		switch r.Intn(3) {
		case 0:
			return []byte("{not-json"), nil
		case 1:
			order, err := generateOrder(r)
			if err != nil {
				return nil, err
			}
			order.OrderUID = ""
			return json.Marshal(order)
		default:
			order, err := generateOrder(r)
			if err != nil {
				return nil, err
			}
			order.Payload = json.RawMessage(`{"delivery":{},"payment":{},"items":[]}`)
			return json.Marshal(order)
		}
	}

	order, err := generateOrder(r)
	if err != nil {
		return nil, err
	}
	return json.Marshal(order)
}

func generateOrder(r *rand.Rand) (models.Order, error) {
	now := time.Now().UTC()
	uid := fmt.Sprintf("order-%d", now.UnixNano()+int64(r.Intn(1000)))
	track := fmt.Sprintf("WB-%06d", r.Intn(1_000_000))

	payload := map[string]any{
		"delivery": map[string]any{
			"name":    fmt.Sprintf("User %d", r.Intn(1000)),
			"phone":   "+79991234567",
			"city":    "Moscow",
			"address": "Lenina 1",
			"email":   fmt.Sprintf("user%d@example.com", r.Intn(1000)),
		},
		"payment": map[string]any{
			"transaction": fmt.Sprintf("tx-%d", now.UnixNano()),
			"currency":    "RUB",
			"amount":      1000 + r.Intn(10000),
		},
		"items": []map[string]any{
			{
				"track_number": track,
				"name":         "Demo item",
				"price":        100 + r.Intn(1000),
			},
		},
	}

	payloadRaw, err := json.Marshal(payload)
	if err != nil {
		return models.Order{}, fmt.Errorf("marshal payload: %w", err)
	}

	return models.Order{
		OrderUID:    uid,
		TrackNumber: track,
		DateCreated: now,
		Payload:     payloadRaw,
		UpdatedAt:   now,
	}, nil
}
