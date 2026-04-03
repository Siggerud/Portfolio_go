package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Sector represents a sector in the portfolio configuration
type Sector struct {
	Weight float64  `yaml:"weight"`
	Stocks []string `yaml:"stocks"`
	Funds  []string `yaml:"funds"`
}

// Config represents the overall configuration
type Config struct {
	Sectors map[string]Sector `yaml:"Sectors"`
}

// ReadConfig reads the YAML configuration file into a Config struct
func ReadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal YAML: %w", err)
	}

	return &config, nil
}
