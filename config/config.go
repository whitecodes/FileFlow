package config

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Port       int
	DBPath     string
	SearchDirs []string
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
	searchDirs := []string{"."}
	if d := os.Getenv("SEARCH_DIRS"); d != "" {
		searchDirs = strings.Split(d, ",")
	}
	return &Config{Port: port, DBPath: dbPath, SearchDirs: searchDirs}
}
