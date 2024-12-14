package handlers

import (
	"health-dashboard/internal/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func HandleSystemInfo(c *gin.Context) {
	info := models.SystemInfo{
		Region:           "East US",
		Hostname:         "webapp-1",
		ContainerVersion: "1.0.0",
		Timestamp:        time.Now(),
	}
	c.JSON(http.StatusOK, info)
}

func HandleHealthStatus(c *gin.Context) {
	status := models.HealthStatus{
		TrafficManager:    "active",
		KeyVault:          "connected",
		ContainerRegistry: "connected",
		RegionStatus:      "healthy",
		Role:              "primary",
		LastChecked:       time.Now(),
	}
	c.JSON(http.StatusOK, status)
}

func HandleLiveCheck(c *gin.Context) {
	c.String(http.StatusOK, "200 OK")
}

func HandleFailoverTrigger(c *gin.Context) {
	var request models.FailoverRequest
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response := gin.H{
		"success":        true,
		"timestamp":      time.Now(),
		"previousRegion": "East US",
		"newRegion":      request.TargetRegion,
	}
	c.JSON(http.StatusOK, response)
}

func HandleFailoverHistory(c *gin.Context) {
	history := models.FailoverHistory{
		LastFailover:   time.Now(),
		CurrentPrimary: "East US",
		FailoverCount:  1,
	}
	c.JSON(http.StatusOK, history)
}
