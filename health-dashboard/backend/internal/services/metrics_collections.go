package services

import (
	"health-dashboard/internal/models"
	"sync"
	"time"
)

type MetricsCollector struct {
	metrics []models.Metric
	mu      sync.RWMutex
}

func NewMetricsCollector() *MetricsCollector {
	return &MetricsCollector{
		metrics: make([]models.Metric, 0),
	}
}

func (mc *MetricsCollector) AddMetric(m models.Metric) {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	cutoff := time.Now().Add(-24 * time.Hour)
	filtered := mc.metrics[:0]
	for _, metric := range mc.metrics {
		if metric.Timestamp.After(cutoff) {
			filtered = append(filtered, metric)
		}
	}

	mc.metrics = append(filtered, m)
}

func (mc *MetricsCollector) GetMetrics() []models.Metric {
	mc.mu.RLock()
	defer mc.mu.RUnlock()

	metrics := make([]models.Metric, len(mc.metrics))
	copy(metrics, mc.metrics)
	return metrics
}
