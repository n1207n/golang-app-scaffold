package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	AppEnv            string
	HTTPServerAddress string
	PostgresURL       string
	RedisURL          string
	RedisHost         string
	RedisPort         string
	RedisPassword     string
	SecretKey         string
}

// LoadConfig loads configuration from environment variables
func LoadConfig(path string) (*Config, error) {
	if os.Getenv("APP_ENV") != "production" { // Only load .env file if not in production
		err := godotenv.Load(path)
		if err != nil {
			log.Println("Warning: .env file not found or error loading .env file:", err)
		}
	}

	return &Config{
		AppEnv:            getEnv("APP_ENV", "development"),
		HTTPServerAddress: getEnv("HTTP_SERVER_ADDRESS", "0.0.0.0:8080"),
		PostgresURL:       getEnv("DATABASE_URL", "postgres://user:password@localhost:5432/mydatabase?sslmode=disable"),
		RedisURL:          getEnv("REDIS_URL", "redis://localhost:6379/0"),
		RedisHost:         getEnv("REDIS_HOST", "localhost"),
		RedisPort:         getEnv("REDIS_PORT", "6379"),
		RedisPassword:     getEnv("REDIS_PASSWORD", ""),
		SecretKey:         getEnv("SECRET_KEY", "supersecret"),
	}, nil
}

// Helper function to get an environment variable or return a default value
func getEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// Helper function to get an environment variable as int or return a default value
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}
