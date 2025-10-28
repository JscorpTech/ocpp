package main

import (
	"context"

	"github.com/JscorpTech/ocpp/internal/config"
	"github.com/JscorpTech/ocpp/internal/ocpp"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewProduction()
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
		DB:   0,
	})

	if err := rdb.Ping(ctx).Err(); err != nil {
		panic(err)
	}

	if err != nil {
		panic(err)
	}
	cfg := config.NewConfig()
	server := ocpp.NewServer(ctx, cfg, logger, rdb)
	go server.Run()
	select {}
}
