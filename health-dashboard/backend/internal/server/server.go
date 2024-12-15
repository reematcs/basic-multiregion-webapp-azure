package server

import (
	"fmt"
	"log"

	"health-dashboard/backend/internal/handlers"
	"health-dashboard/backend/internal/services"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/gin-gonic/gin"
)

type Server struct {
	cred             *azidentity.DefaultAzureCredential
	healthChecker    *services.HealthChecker
	metricsCollector *services.MetricsCollector
}

func NewServer(cred *azidentity.DefaultAzureCredential, subscriptionID, resourceGroup, region, role, profileName string) (*Server, error) {
	healthChecker, err := services.NewHealthChecker(cred, subscriptionID, resourceGroup, region, role, profileName)
	if err != nil {
		return nil, fmt.Errorf("failed to create health checker: %v", err)
	}

	metricsCollector := services.NewMetricsCollector()

	return &Server{
		cred:             cred,
		healthChecker:    healthChecker,
		metricsCollector: metricsCollector,
	}, nil
}

func (s *Server) setupRoutes() (*gin.Engine, error) {
	r := gin.Default()

	healthHandler := handlers.NewHealthHandler(s.healthChecker, s.metricsCollector)

	api := r.Group("/api")
	{
		api.GET("/system", healthHandler.HandleSystemInfo)
		api.GET("/health/status", healthHandler.HandleHealthStatus)
		api.GET("/health/live", healthHandler.HandleLiveCheck)
		api.POST("/failover/trigger", healthHandler.HandleFailoverTrigger)
		api.GET("/failover/history", healthHandler.HandleFailoverHistory)
	}

	// Serve static files from the /static directory
	r.Static("/static", "/static")
	r.StaticFile("/", "/static/index.html")
	r.NoRoute(func(c *gin.Context) {
		c.File("/static/index.html")
	})

	return r, nil
}

func (s *Server) Run(addr string) error {
	r, err := s.setupRoutes()
	if err != nil {
		return fmt.Errorf("failed to setup routes: %v", err)
	}

	log.Printf("Server is running on %s", addr)
	return r.Run(addr)
}
