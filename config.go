package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Rule struct {
	Name       string `yaml:"name"`
	Pattern    string `yaml:"pattern"`
	TargetDir  string `yaml:"target_dir"`
	TargetName string `yaml:"target_name"`
}

type Config struct {
	Rules []Rule `yaml:"rules"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}
	return &cfg, nil
}
