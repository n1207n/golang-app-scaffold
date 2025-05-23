package main

import (
	"github.com/yourusername/yourprojectname/config" // TODO: Replace yourprojectname with your actual module path
	"log"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig(".env") // Load .env file from the root directory
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	log.Printf("Configuration loaded successfully. App Env: %s", cfg.AppEnv)
	log.Printf("Server starting on %s", cfg.HTTPServerAddress)

	// TODO: Initialize DB connection (Postgres with pgx)
	// TODO: Initialize Redis connection
	// TODO: Initialize Repositories (DB and Cache)
	// TODO: Initialize Services
	// TODO: Initialize Gin router and Handlers
	// TODO: Start HTTP server

	// Placeholder for server start
	// router := gin.Default()
	// router.GET("/ping", func(c *gin.Context) {
	// 	c.JSON(200, gin.H{
	// 		"message": "pong",
	// 	})
	// })
	// log.Fatal(router.Run(cfg.HTTPServerAddress))
}
