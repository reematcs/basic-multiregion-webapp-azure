// main.go
package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/trafficmanager/armtrafficmanager"
	"github.com/gin-gonic/gin"
)

// SystemInfo represents the basic region information
type SystemInfo struct {
	Region           string    `json:"region"`
	Hostname         string    `json:"hostname"`
	ContainerVersion string    `json:"containerVersion"`
	Timestamp        time.Time `json:"timestamp"`
}

// HealthStatus represents the health check response
type HealthStatus struct {
	TrafficManager    string    `json:"trafficManager"`
	KeyVault          string    `json:"keyVault"`
	ContainerRegistry string    `json:"containerRegistry"`
	RegionStatus      string    `json:"regionStatus"`
	Role              string    `json:"role"`
	LastChecked       time.Time `json:"lastChecked"`
}

// FailoverHistory represents the failover status
type FailoverHistory struct {
	LastFailover   time.Time `json:"lastFailover"`
	CurrentPrimary string    `json:"currentPrimary"`
	FailoverCount  int       `json:"failoverCount"`
}
type Server struct {
	cred           *azidentity.DefaultAzureCredential
	profilesClient *armtrafficmanager.ProfilesClient
}

func NewServer() (*Server, error) {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create credential: %v", err)
	}

	clientFactory, err := armtrafficmanager.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create client factory: %v", err)
	}

	// Create the ProfilesClient
	profilesClient := clientFactory.NewProfilesClient()

	return &Server{
		cred:           cred,
		profilesClient: profilesClient,
	}, nil
}

func (s *Server) setupRoutes() *gin.Engine {
	r := gin.Default()

	// System information endpoint
	r.GET("/api/system", func(c *gin.Context) {
		info := SystemInfo{
			Region:           "East US",
			Hostname:         "webapp-1",
			ContainerVersion: "1.0.0",
			Timestamp:        time.Now(),
		}
		c.JSON(http.StatusOK, info)
	})

	// Health status endpoint
	r.GET("/api/health/status", func(c *gin.Context) {
		status := HealthStatus{
			TrafficManager:    "active",
			KeyVault:          "connected",
			ContainerRegistry: "connected",
			RegionStatus:      "healthy",
			Role:              "primary",
			LastChecked:       time.Now(),
		}
		c.JSON(http.StatusOK, status)
	})

	// Live health check endpoint
	r.GET("/api/health/live", func(c *gin.Context) {
		c.String(http.StatusOK, "200 OK")
	})

	// Failover trigger endpoint
	r.POST("/api/failover/trigger", func(c *gin.Context) {
		var request struct {
			TargetRegion string `json:"targetRegion"`
		}
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
	})

	// Failover history endpoint
	r.GET("/api/failover/history", func(c *gin.Context) {
		history := FailoverHistory{
			LastFailover:   time.Now(),
			CurrentPrimary: "East US",
			FailoverCount:  1,
		}
		c.JSON(http.StatusOK, history)
	})

	return r
}

func main() {
	server, err := NewServer()
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	r := server.setupRoutes()
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
