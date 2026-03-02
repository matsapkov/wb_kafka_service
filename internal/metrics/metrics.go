package metrics

import "github.com/prometheus/client_golang/prometheus"

type Metrics struct {
	ProcessedMessages prometheus.Counter
	InvalidMessages   prometheus.Counter
	DLQPublished      prometheus.Counter
	CacheHits         prometheus.Counter
	CacheMisses       prometheus.Counter
}

func New() *Metrics {
	m := &Metrics{
		ProcessedMessages: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "orders_processed_total",
			Help: "Total number of processed Kafka order messages.",
		}),
		InvalidMessages: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "orders_invalid_total",
			Help: "Total number of invalid Kafka order messages.",
		}),
		DLQPublished: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "orders_dlq_total",
			Help: "Total number of messages published to DLQ.",
		}),
		CacheHits: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "orders_cache_hits_total",
			Help: "Total number of cache hits for order requests.",
		}),
		CacheMisses: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "orders_cache_misses_total",
			Help: "Total number of cache misses for order requests.",
		}),
	}

	prometheus.MustRegister(
		m.ProcessedMessages,
		m.InvalidMessages,
		m.DLQPublished,
		m.CacheHits,
		m.CacheMisses,
	)

	return m
}
