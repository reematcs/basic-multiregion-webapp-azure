package handlers

import (
	"context"
	"net/http"
	"sort"
	"strconv"
	"time"

	"health-dashboard/backend/internal/models"
	"health-dashboard/backend/internal/services"

	"github.com/gin-gonic/gin"
)

type HealthHandler struct {
	healthChecker    *services.HealthChecker
	metricsCollector *services.MetricsCollector
}

func NewHealthHandler(healthChecker *services.HealthChecker, metricsCollector *services.MetricsCollector) *HealthHandler {
	return &HealthHandler{
		healthChecker:    healthChecker,
		metricsCollector: metricsCollector,
	}
}

func (h *HealthHandler) HandleSystemInfo(c *gin.Context) {
	info := models.SystemInfo{
		Region:           "West US",
		Hostname:         "webapp-1",
		ContainerVersion: "1.0.0",
		Timestamp:        time.Now(),
	}
	c.JSON(http.StatusOK, info)
}

func (h *HealthHandler) HandleHealthStatus(c *gin.Context) {
	ctx := context.Background()
	startTime := time.Now()

	status, err := h.healthChecker.CheckHealth(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Record latency metric
	latency := time.Since(startTime).Milliseconds()
	h.metricsCollector.AddMetric(models.Metric{
		Name:      "latency",
		Value:     strconv.FormatInt(latency, 10),
		Timestamp: time.Now(),
	})

	// Record request count metric
	h.metricsCollector.AddMetric(models.Metric{
		Name:      "request_count",
		Value:     "1",
		Timestamp: time.Now(),
	})

	// Record health status metric
	h.metricsCollector.AddMetric(models.Metric{
		Name:      "health_status",
		Value:     status.RegionStatus,
		Timestamp: time.Now(),
	})

	c.JSON(http.StatusOK, status)
}

func (h *HealthHandler) HandleLiveCheck(c *gin.Context) {
	c.String(http.StatusOK, "200 OK")
}

func (h *HealthHandler) HandleFailoverTrigger(c *gin.Context) {
	var request struct {
		TargetRegion string `json:"targetRegion"`
	}
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := context.Background()
	err := h.healthChecker.TriggerFailover(ctx, request.TargetRegion)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := gin.H{
		"success":        true,
		"timestamp":      time.Now(),
		"previousRegion": "West US",
		"newRegion":      request.TargetRegion,
	}
	c.JSON(http.StatusOK, response)
}

func (h *HealthHandler) HandleFailoverHistory(c *gin.Context) {
	ctx := context.Background()
	history, err := h.healthChecker.GetFailoverHistory(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, history)
}
func (h *HealthHandler) HandleMetrics(c *gin.Context) {
	// Transform metrics to the format expected by frontend
	type FrontendMetric struct {
		Timestamp    time.Time `json:"timestamp"`
		Latency      float64   `json:"latency"`
		RequestCount int64     `json:"requestCount"`
	}

	// Get metrics for the last hour
	oneHourAgo := time.Now().Add(-1 * time.Hour)
	latencyMetrics := h.metricsCollector.GetMetricsByName("latency")
	requestMetrics := h.metricsCollector.GetMetricsByName("request_count")

	// Create a map to hold the combined metrics
	metricsByMinute := make(map[time.Time]FrontendMetric)

	// Process latency metrics
	for _, metric := range latencyMetrics {
		if metric.Timestamp.After(oneHourAgo) {
			// Parse latency value from string
			latencyValue, err := strconv.ParseFloat(metric.Value, 64)
			if err != nil {
				continue // Skip invalid values
			}

			// Round to nearest minute for grouping
			minute := metric.Timestamp.Truncate(time.Minute)
			m := metricsByMinute[minute]
			m.Timestamp = minute
			m.Latency = latencyValue
			metricsByMinute[minute] = m
		}
	}

	// Process request count metrics
	for _, metric := range requestMetrics {
		if metric.Timestamp.After(oneHourAgo) {
			// Parse request count value from string
			countValue, err := strconv.ParseInt(metric.Value, 10, 64)
			if err != nil {
				continue // Skip invalid values
			}

			minute := metric.Timestamp.Truncate(time.Minute)
			m := metricsByMinute[minute]
			m.Timestamp = minute
			m.RequestCount = countValue
			metricsByMinute[minute] = m
		}
	}

	// Convert map to sorted slice
	var result []FrontendMetric
	for _, metric := range metricsByMinute {
		result = append(result, metric)
	}

	// Sort by timestamp
	sort.Slice(result, func(i, j int) bool {
		return result[i].Timestamp.Before(result[j].Timestamp)
	})

	// If no metrics are available, provide sample data
	if len(result) == 0 {
		now := time.Now()
		result = []FrontendMetric{
			{
				Timestamp:    now.Add(-2 * time.Minute),
				Latency:      120,
				RequestCount: 45,
			},
			{
				Timestamp:    now.Add(-1 * time.Minute),
				Latency:      115,
				RequestCount: 52,
			},
			{
				Timestamp:    now,
				Latency:      118,
				RequestCount: 48,
			},
		}
	}

	c.JSON(http.StatusOK, result)
}
