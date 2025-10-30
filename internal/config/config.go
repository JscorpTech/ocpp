package config

import (
	"os"
)

type Config struct {
	Addr    string
	BaseUrl string
}

func NewConfig() *Config {
	baseUrl := os.Getenv("BASE_URL")
	if baseUrl == "" {
		panic("Base url is required")
	}
	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = ":10800"
	}
	return &Config{
		BaseUrl: baseUrl,
		Addr:    addr,
	}
}
