package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"os"
	"time"

	"github.com/matsapkov/wb_kafka_service/internal/kafka"
	"github.com/matsapkov/wb_kafka_service/internal/models"
)

func main() {
	filePath := flag.String("file", "", "Path to JSON file with order data")
	count := flag.Int("count", 1, "Number of test messages to send")
	flag.Parse()

	brokers := []string{"localhost:9092"}
	topic := "orders"

	producer, err := kafka.NewKafkaProducer(brokers, topic)
	if err != nil {
		log.Fatalf("Failed to create producer: %v", err)
	}
	defer producer.Close()
	log.Printf("Используемые brokers перед вызовом: %v", brokers)

	ctx := context.Background()

	if *filePath != "" {
		data, err := os.ReadFile(*filePath)
		if err != nil {
			log.Fatalf("Failed to read file: %v", err)
		}
		err = producer.SendMessage(ctx, data)
		if err != nil {
			log.Fatalf("Failed to send: %v", err)
		}
		log.Println("Sent message from file")
		return
	}

	for i := 1; i <= *count; i++ {
		order := generateTestOrder(i)
		data, err := json.Marshal(order)
		if err != nil {
			log.Printf("Marshal error: %v", err)
			continue
		}

		err = producer.SendMessage(ctx, data)
		if err != nil {
			log.Printf("Send error: %v", err)
			continue
		}

		time.Sleep(time.Millisecond * 500)
	}

	log.Printf("Sent %d test messages", *count)
}

func generateTestOrder(num int) models.Order {
	return models.Order{
		OrderUID:    "test-order-" + string(rune('A'+num-1)),
		TrackNumber: "TN-TEST-" + string(rune('0'+num)),
		DateCreated: time.Now().Add(-time.Hour * time.Duration(num)),
		UpdatedAt:   time.Now(),
		Payload: json.RawMessage(`{
			"delivery": {"name": "Test User ` + string(rune('A'+num-1)) + `", "phone": "+1234567890", "zip": "12345", "city": "Test City", "address": "Test Address", "region": "Test Region", "email": "test@example.com"},
			"payment": {"transaction": "trans-` + string(rune('0'+num)) + `", "currency": "USD", "provider": "test-provider", "amount": 100` + string(rune('0'+num)) + `, "payment_dt": 1234567890, "bank": "test-bank", "delivery_cost": 500, "goods_total": 9500, "custom_fee": 0},
			"items": [
				{"chrt_id": 12345` + string(rune('0'+num)) + `, "track_number": "TN-TEST-` + string(rune('0'+num)) + `", "price": 1000, "rid": "rid-` + string(rune('0'+num)) + `", "name": "Test Item ` + string(rune('A'+num-1)) + `", "sale": 10, "size": "M", "total_price": 900, "nm_id": 98765, "brand": "Test Brand", "status": 1}
			]
		}`),
	}
}
