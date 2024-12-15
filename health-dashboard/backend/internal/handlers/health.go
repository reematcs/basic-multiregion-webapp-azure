package handlers

import (
	"context"
	"net/http"
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
	status, err := h.healthChecker.CheckHealth(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

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
