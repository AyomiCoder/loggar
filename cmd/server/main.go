package main

import (
	"log"
	"os"

	"github.com/AyomiCoder/loggar/api"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Get configuration from environment
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	// Initialize database
	if err := api.InitDB(databaseURL); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Start server
	log.Printf("Starting server on port %s...", port)
	if err := api.Run(port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
