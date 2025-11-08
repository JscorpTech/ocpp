package main

import (
	"context"
	"log"

	"github.com/JscorpTech/ocpp/internal/config"
	"github.com/JscorpTech/ocpp/internal/ocpp"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	_ = godotenv.Load()
	cfg := config.NewConfig()
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr: cfg.RedisAddr,
		DB:   0,
	})
	if err := rdb.Ping(ctx).Err(); err != nil {
		panic(err)
	}
	server := ocpp.NewServer(ctx, cfg, logger, rdb)
	go func() {
		if err := server.Run(); err != nil {
			log.Panic(err)
		}
	}()
	select {}
}
