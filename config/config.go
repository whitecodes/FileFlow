package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port   int
	DBPath string
}

func Load() *Config {
	port := 8080
	if p := os.Getenv("PORT"); p != "" {
		if v, err := strconv.Atoi(p); err == nil {
			port = v
		}
	}
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "fileflow.db"
	}
	return &Config{Port: port, DBPath: dbPath}
}
