// main.go
package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"

	"health-dashboard/backend/internal/config"
	"health-dashboard/backend/internal/server"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

func main() {

	// Add debugging support
	if os.Getenv("DEBUG") == "true" {
		gin.SetMode(gin.DebugMode)
		// Add delve debugging port
		log.Println("Debugger is enabled on port 2345")
	}
	// Load configuration from environment variables
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Get Azure credentials
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("Failed to get credentials: %v", err)
	}

	// Create server with configuration
	srv, err := server.NewServer(
		cred,
		cfg.SubscriptionID,
		cfg.ResourceGroup,
		cfg.Region,
		cfg.Role,
		cfg.ProfileName,
	)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	if err := srv.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// package main

// import (
// 	"fmt"
// 	"log"
// 	"net/http"
// 	"time"

// 	"health-dashboard/backend/internal/models"

// 	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
// 	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/trafficmanager/armtrafficmanager"
// 	"github.com/gin-gonic/gin"
// )

// type Server struct {
// 	cred           *azidentity.DefaultAzureCredential
// 	profilesClient *armtrafficmanager.ProfilesClient
// }

// func NewServer() (*Server, error) {
// 	cred, err := azidentity.NewDefaultAzureCredential(nil)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to create credential: %v", err)
// 	}

// 	clientFactory, err := armtrafficmanager.NewClientFactory("<subscription-id>", cred, nil)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to create client factory: %v", err)
// 	}

// 	return &Server{
// 		cred:           cred,
// 		profilesClient: clientFactory.NewProfilesClient(),
// 	}, nil
// }

// func (s *Server) setupRoutes() *gin.Engine {
// 	r := gin.Default()

// 	// API routes group
// 	api := r.Group("/api")
// 	{
// 		api.GET("/system", func(c *gin.Context) {
// 			info := models.SystemInfo{
// 				Region:           "West US",
// 				Hostname:         "webapp-1",
// 				ContainerVersion: "1.0.0",
// 				Timestamp:        time.Now(),
// 			}
// 			c.JSON(http.StatusOK, info)
// 		})

// 		api.GET("/health/status", func(c *gin.Context) {
// 			status := models.HealthStatus{
// 				TrafficManager:    "active",
// 				KeyVault:          "connected",
// 				ContainerRegistry: "connected",
// 				RegionStatus:      "healthy",
// 				Role:              "primary",
// 				LastChecked:       time.Now(),
// 			}
// 			c.JSON(http.StatusOK, status)
// 		})

// 		api.GET("/health/live", func(c *gin.Context) {
// 			c.String(http.StatusOK, "200 OK")
// 		})

// 		api.POST("/failover/trigger", func(c *gin.Context) {
// 			var request struct {
// 				TargetRegion string `json:"targetRegion"`
// 			}
// 			if err := c.BindJSON(&request); err != nil {
// 				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 				return
// 			}

// 			response := gin.H{
// 				"success":        true,
// 				"timestamp":      time.Now(),
// 				"previousRegion": "West US",
// 				"newRegion":      request.TargetRegion,
// 			}
// 			c.JSON(http.StatusOK, response)
// 		})

// 		api.GET("/failover/history", func(c *gin.Context) {
// 			history := models.FailoverHistory{
// 				LastFailover:   time.Now(),
// 				CurrentPrimary: "West US",
// 				FailoverCount:  1,
// 			}
// 			c.JSON(http.StatusOK, history)
// 		})
// 	}

// 	// Serve static files - note the nested static directory
// 	r.Static("/static", "/app/static/static")

// 	// Serve index.html for root path and SPA routes
// 	r.GET("/", serveIndex)
// 	r.NoRoute(serveIndex)

// 	return r
// }

// func serveIndex(c *gin.Context) {
// 	c.File("/app/static/index.html")
// }

// func main() {
// 	server, err := NewServer()
// 	if err != nil {
// 		log.Fatalf("Failed to create server: %v", err)
// 	}

// 	r := server.setupRoutes()
// 	if err := r.Run(":8080"); err != nil {
// 		log.Fatalf("Failed to start server: %v", err)
// 	}
// }
