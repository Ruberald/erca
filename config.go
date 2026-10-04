package main

import (
	"encoding/json"
	"os"
)

type Config struct {
	Port     int      `json:"port"`
	Backends []string `json:"backends"`
}

func loadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config Config

	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
