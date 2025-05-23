package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yourusername/yourprojectname/config" // Ensure this matches your go.mod module name
	// TODO: Import repository, service, handler packages once created
	// "github.com/yourusername/yourprojectname/internal/handler"
	// "github.com/yourusername/yourprojectname/internal/repository"
	// "github.com/yourusername/yourprojectname/internal/service"
	// "github.com/yourusername/yourprojectname/db/sqlc" // SQLC generated code
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig(".env")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	log.Printf("Configuration loaded successfully. App Env: %s, Server: %s", cfg.AppEnv, cfg.HTTPServerAddress)

	// Initialize Database connection
	dbPool, err := initDB(cfg.PostgresURL)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer dbPool.Close()
	log.Println("Database connection pool established.")

	// Initialize Redis client
	rdb, err := initRedis(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize Redis: %v", err)
	}
	defer rdb.Close()
	log.Println("Redis client initialized.")

	// TODO: Initialize SQLC store (example, adjust as per your sqlc setup)
	// store := db.NewStore(dbPool) // Assuming NewStore takes *pgxpool.Pool

	// TODO: Initialize Repositories
	// userRepository := repository.NewUserRepository(store, rdb) // Example

	// TODO: Initialize Services
	// userService := service.NewUserService(userRepository) // Example

	// Initialize Gin router
	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.Default()

	// TODO: Setup routes using Handlers
	// userHandler := handler.NewUserHandler(userService)
	// v1 := router.Group("/api/v1")
	// {
	// 	userRoutes := v1.Group("/users")
	// 	{
	// 		userRoutes.POST("", userHandler.CreateUser)
	// 		userRoutes.GET("/:id", userHandler.GetUser)
	// 	}
	// }

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	// Start HTTP server
	srv := &http.Server{
		Addr:    cfg.HTTPServerAddress,
		Handler: router,
	}

	go func() {
		log.Printf("Server listening on %s", cfg.HTTPServerAddress)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}

func initDB(databaseURL string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("unable to parse database URL: %w", err)
	}

	// You can configure pool settings here, e.g.,
	// config.MaxConns = 10
	// config.MinConns = 2
	// config.MaxConnLifetime = time.Hour
	// config.MaxConnIdleTime = 30 * time.Minute
	// config.HealthCheckPeriod = time.Minute
	// config.ConnConfig.ConnectTimeout = 5 * time.Second

	dbPool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	// Ping the database to verify connection
	if err = dbPool.Ping(context.Background()); err != nil {
		dbPool.Close() // Close the pool if ping fails
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	return dbPool, nil
}

func initRedis(cfg *config.Config) (*redis.Client, error) {
	var rdb *redis.Client
	if cfg.RedisURL != "" { // Prefer RedisURL if available
		opt, err := redis.ParseURL(cfg.RedisURL)
		if err != nil {
			return nil, fmt.Errorf("could not parse Redis URL: %w", err)
		}
		rdb = redis.NewClient(opt)
	} else { // Fallback to host/port if RedisURL is not set
		rdb = redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort),
			Password: cfg.RedisPassword, // no password set if empty
			DB:       0,                 // use default DB
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		rdb.Close() // Close the client if ping fails
		return nil, fmt.Errorf("could not ping Redis: %w", err)
	}
	return rdb, nil
}
