package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yourusername/yourprojectname/config"

	// "github.com/yourusername/yourprojectname/internal/handler"
	"github.com/yourusername/yourprojectname/internal/repository"
	// "github.com/yourusername/yourprojectname/internal/service"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	log.Printf("Configuration loaded successfully. App Env: %s, Server: %d", cfg.AppEnv, cfg.AppPort)

	// Initialize Database connection
	dbPool, err := initDB(cfg.DbURL)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer dbPool.Close()
	log.Println("Database connection pool established.")

	// Initialize Redis client
	rdb, err := initRedis(cfg.RedisURL)
	if err != nil {
		log.Fatalf("Failed to initialize Redis: %v", err)
	}
	defer rdb.Close()
	log.Println("Redis client initialized.")

	// Initialize Repositories
	userRepo := repository.NewDBUserRepository(sqlcQuerier)
	log.Println("User repository initialized.")

	// TODO: Initialize Services
	// userService := service.NewUserService(userRepo) // Example
	// log.Println("User service initialized.")

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
		Addr:    strconv.Itoa(cfg.AppPort),
		Handler: router,
	}

	go func() {
		log.Printf("Server listening on %d", cfg.AppPort)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
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
	pgxpoolCfg, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("unable to parse database URL: %w", err)
	}

	// You can configure pool settings here, e.g.,
	// pgxpool_cfg.MaxConns = 10
	// pgxpool_cfg.MinConns = 2
	// pgxpool_cfg.MaxConnLifetime = time.Hour
	// pgxpool_cfg.MaxConnIdleTime = 30 * time.Minute
	// pgxpool_cfg.HealthCheckPeriod = time.Minute
	// pgxpool_cfg.ConnConfig.ConnectTimeout = 5 * time.Second

	dbPool, err := pgxpool.NewWithConfig(context.Background(), pgxpoolCfg)
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

func initRedis(redisURL string) (*redis.Client, error) {
	var rdb *redis.Client

	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("could not parse Redis URL: %w", err)
	}
	rdb = redis.NewClient(opt)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = rdb.Ping(ctx).Result()
	if err != nil {
		rdb.Close() // Close the client if ping fails
		return nil, fmt.Errorf("could not ping Redis: %w", err)
	}
	return rdb, nil
}
