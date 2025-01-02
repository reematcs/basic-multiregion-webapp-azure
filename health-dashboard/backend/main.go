package main

import (
	"log"
	"os"

	"health-dashboard/backend/internal/config"
	"health-dashboard/backend/internal/server"

	"github.com/gin-gonic/gin"
)

func main() {
	log.Printf("Starting application...")
	log.Printf("Environment variables:")
	log.Printf("REGION: %s", os.Getenv("REGION"))
	log.Printf("ROLE: %s", os.Getenv("ROLE"))
	log.Printf("LOCAL_MODE: %s", os.Getenv("LOCAL_MODE"))

	// Add debugging support
	if os.Getenv("DEBUG") == "true" {
		gin.SetMode(gin.DebugMode)
		log.Println("Debug mode enabled")
	}

	// Load configuration from environment variables
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	log.Printf("Configuration loaded successfully")

	// Create server
	srv, err := server.NewServer(cfg)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}
	log.Printf("Server created successfully")

	// Start server
	log.Printf("Starting server...")
	if err := srv.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
