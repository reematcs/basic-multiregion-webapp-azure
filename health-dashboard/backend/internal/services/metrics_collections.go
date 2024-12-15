// services/metrics_collector.go
package services

import (
	"health-dashboard/backend/internal/models"
	"sort"
	"sync"
	"time"
)

type MetricsCollector struct {
	metrics      []models.Metric
	mu           sync.RWMutex
	retentionDur time.Duration
}

func NewMetricsCollector() *MetricsCollector {
	return &MetricsCollector{
		metrics:      make([]models.Metric, 0),
		retentionDur: 24 * time.Hour, // Default 24-hour retention
	}
}

func (mc *MetricsCollector) AddMetric(m models.Metric) {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	// Clean up old metrics
	cutoff := time.Now().Add(-mc.retentionDur)
	filtered := mc.metrics[:0]
	for _, metric := range mc.metrics {
		if metric.Timestamp.After(cutoff) {
			filtered = append(filtered, metric)
		}
	}

	// Add new metric
	filtered = append(filtered, m)

	// Sort by timestamp
	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].Timestamp.Before(filtered[j].Timestamp)
	})

	mc.metrics = filtered
}

func (mc *MetricsCollector) GetMetrics() []models.Metric {
	mc.mu.RLock()
	defer mc.mu.RUnlock()

	// Return a copy to prevent external modifications
	metrics := make([]models.Metric, len(mc.metrics))
	copy(metrics, mc.metrics)
	return metrics
}

func (mc *MetricsCollector) GetMetricsByName(name string) []models.Metric {
	mc.mu.RLock()
	defer mc.mu.RUnlock()

	var filtered []models.Metric
	for _, metric := range mc.metrics {
		if metric.Name == name {
			filtered = append(filtered, metric)
		}
	}
	return filtered
}

func (mc *MetricsCollector) GetMetricsInTimeRange(start, end time.Time) []models.Metric {
	mc.mu.RLock()
	defer mc.mu.RUnlock()

	var filtered []models.Metric
	for _, metric := range mc.metrics {
		if (metric.Timestamp.Equal(start) || metric.Timestamp.After(start)) &&
			(metric.Timestamp.Equal(end) || metric.Timestamp.Before(end)) {
			filtered = append(filtered, metric)
		}
	}
	return filtered
}

func (mc *MetricsCollector) SetRetentionDuration(duration time.Duration) {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	mc.retentionDur = duration
}

func (mc *MetricsCollector) ClearMetrics() {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	mc.metrics = mc.metrics[:0]
}
