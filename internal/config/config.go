package config

type Config struct {
	Addr string
}

func NewConfig() *Config {
	return &Config{
		Addr: ":10800",
	}
}
