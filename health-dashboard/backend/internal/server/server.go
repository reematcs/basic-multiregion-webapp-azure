package server

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"health-dashboard/backend/internal/config"
	"health-dashboard/backend/internal/handlers"
	"health-dashboard/backend/internal/services"
)

type Server struct {
	config        *config.Config
	healthHandler *handlers.HealthHandler
}

func NewServer(cfg *config.Config) (*Server, error) {
	// Initialize Azure credentials
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create Azure credentials: %v", err)
	}

	// Initialize services
	healthChecker, err := services.NewHealthChecker(
		cred,
		cfg.SubscriptionID,
		cfg.ResourceGroup,
		cfg.Region,
		cfg.Role,
		cfg.ProfileName,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create health checker: %v", err)
	}

	metricsCollector := services.NewMetricsCollector()

	// Initialize handlers
	healthHandler := handlers.NewHealthHandler(healthChecker, metricsCollector)

	return &Server{
		config:        cfg,
		healthHandler: healthHandler,
	}, nil
}

func (s *Server) Start() error {
	log.Printf("Initializing server...")

	// Switch to debug mode for development
	gin.SetMode(gin.DebugMode)

	r := gin.Default()

	// Enable CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://frontend-dev:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	log.Printf("Setting up routes...")
	// API routes group
	api := r.Group("/api")
	{
		api.GET("/system", s.healthHandler.HandleSystemInfo)
		api.GET("/health/status", s.healthHandler.HandleHealthStatus)
		api.GET("/health/live", s.healthHandler.HandleLiveCheck)
		api.GET("/metrics", s.healthHandler.HandleMetrics)
		api.POST("/failover/trigger", s.healthHandler.HandleFailoverTrigger)
		api.GET("/failover/history", s.healthHandler.HandleFailoverHistory)
	}

	log.Printf("Server initialization complete, starting to listen on port %s", s.config.Port)
	return r.Run(":" + s.config.Port)
}
