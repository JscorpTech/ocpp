package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/JscorpTech/ocpp/internal/config"
	"github.com/JscorpTech/ocpp/internal/ocpp"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func main() {
	// Add panic recovery to ensure errors are logged before exit
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintf(os.Stderr, "FATAL: Application panicked: %v\n", r)
			os.Exit(1)
		}
	}()

	// Initialize logger with error handling
	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Fprintf(os.Stderr, "FATAL: Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	// Log application startup
	logger.Info("Starting OCPP server application")

	// Load environment variables (non-fatal if .env doesn't exist)
	if err := godotenv.Load(); err != nil {
		logger.Warn("Failed to load .env file, using environment variables", zap.Error(err))
	}

	// Load configuration with proper error handling
	cfg, err := config.NewConfig()
	if err != nil {
		logger.Fatal("Failed to load configuration", zap.Error(err))
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize Redis client with proper error handling
	logger.Info("Connecting to Redis", zap.String("addr", cfg.RedisAddr))
	rdb := redis.NewClient(&redis.Options{
		Addr: cfg.RedisAddr,
		DB:   0,
	})
	defer rdb.Close()

	// Test Redis connection with retry logic
	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		if err := rdb.Ping(ctx).Err(); err != nil {
			logger.Warn("Failed to connect to Redis, retrying...",
				zap.Error(err),
				zap.Int("attempt", i+1),
				zap.Int("max_retries", maxRetries))
			if i == maxRetries-1 {
				logger.Fatal("Failed to connect to Redis after max retries", zap.Error(err))
				os.Exit(1)
			}
			time.Sleep(time.Second * 2)
			continue
		}
		logger.Info("Successfully connected to Redis")
		break
	}

	// Initialize server
	server := ocpp.NewServer(ctx, cfg, logger, rdb)

	// Channel to capture server errors
	errChan := make(chan error, 1)

	// Run server in goroutine with panic recovery
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logger.Error("Server goroutine panicked", zap.Any("panic", r))
				errChan <- fmt.Errorf("server panic: %v", r)
			}
		}()

		logger.Info("Starting OCPP server", zap.String("addr", cfg.Addr))
		if err := server.Run(); err != nil {
			logger.Error("Server error", zap.Error(err))
			errChan <- err
		}
	}()

	// Setup signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	// Wait for shutdown signal or error
	select {
	case sig := <-sigChan:
		logger.Info("Received shutdown signal", zap.String("signal", sig.String()))
	case err := <-errChan:
		logger.Error("Server stopped with error", zap.Error(err))
		os.Exit(1)
	}

	// Graceful shutdown
	logger.Info("Shutting down gracefully...")
	cancel()

	// Give server time to cleanup
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	<-shutdownCtx.Done()
	logger.Info("Server shutdown complete")
}
