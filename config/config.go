package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port int
}

func Load() *Config {
	port := 8080
	if p := os.Getenv("PORT"); p != "" {
		if v, err := strconv.Atoi(p); err == nil {
			port = v
		}
	}
	return &Config{Port: port}
}
