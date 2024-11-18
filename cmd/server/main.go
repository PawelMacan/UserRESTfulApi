package main

import (
	"log"
	"os"

	"github.com/user-api/internal"
)

func main() {
	// Initialize router
	router := internal.NewRouter()
	router.SetupRoutes()

	// Get port from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start server
	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
