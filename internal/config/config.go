package config

import (
	"errors"
	"os"
)

type Config struct {
	Addr      string
	BaseUrl   string
	RedisAddr string
}

func NewConfig() (*Config, error) {
	baseUrl := os.Getenv("BASE_URL")
	if baseUrl == "" {
		return nil, errors.New("BASE_URL environment variable is required")
	}
	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = ":10800"
	}

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "127.0.0.1:6379"
	}
	return &Config{
		BaseUrl:   baseUrl,
		Addr:      addr,
		RedisAddr: redisAddr,
	}, nil
}
