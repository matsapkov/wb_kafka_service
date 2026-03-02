package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	HTTPAddr             string
	DBDSN                string
	CacheLimit           int
	CacheTTL             time.Duration
	CacheCleanupInterval time.Duration
	Kafka                KafkaConfig
}

type KafkaConfig struct {
	Brokers  []string
	Topic    string
	DLQTopic string
	GroupID  string
	MinBytes int
	MaxBytes int
	MaxWait  time.Duration
}

func Load() (Config, error) {
	cfg := Config{
		HTTPAddr:             getEnv("HTTP_ADDR", ":8081"),
		DBDSN:                os.Getenv("DB_DSN"),
		CacheLimit:           getEnvAsInt("CACHE_LIMIT", 1000),
		CacheTTL:             getEnvAsDuration("CACHE_TTL", 5*time.Minute),
		CacheCleanupInterval: getEnvAsDuration("CACHE_CLEANUP_INTERVAL", time.Minute),
		Kafka: KafkaConfig{
			Brokers:  splitCSV(getEnv("KAFKA_BROKERS", "localhost:9092")),
			Topic:    getEnv("KAFKA_TOPIC", "orders"),
			DLQTopic: getEnv("KAFKA_DLQ_TOPIC", "orders.dlq"),
			GroupID:  getEnv("KAFKA_GROUP", "order-service-group"),
			MinBytes: getEnvAsInt("KAFKA_MIN_BYTES", 10_000),
			MaxBytes: getEnvAsInt("KAFKA_MAX_BYTES", 10_000_000),
			MaxWait:  getEnvAsDuration("KAFKA_MAX_WAIT", time.Second),
		},
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return fallback
	}
	return v
}

func getEnvAsInt(key string, fallback int) int {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(v)
	if err != nil {
		return fallback
	}
	return parsed
}

func getEnvAsDuration(key string, fallback time.Duration) time.Duration {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return fallback
	}
	parsed, err := time.ParseDuration(v)
	if err != nil {
		return fallback
	}
	return parsed
}

func splitCSV(value string) []string {
	parts := strings.Split(value, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			out = append(out, trimmed)
		}
	}
	return out
}
