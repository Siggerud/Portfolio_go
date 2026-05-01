package core

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

// Sectors represents the overall configuration
type Sectors struct {
	Sectors map[string]Sector `yaml:"Sectors"`
}

type Pattern struct {
	StockFilenamePattern string `yaml:"stocks"`
	FundFilenamepattern  string `yaml:"funds"`
}

type Config struct {
	Patterns Pattern `yaml:"Patterns"`
}

func ReadConfigYaml(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config YAML: %w", err)
	}

	return &config, nil
}

// ReadSectorYaml reads the YAML configuration file into a Config struct
func ReadSectorYaml(filename string) (*Sectors, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var sectors Sectors
	err = yaml.Unmarshal(data, &sectors)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal YAML: %w", err)
	}

	return &sectors, nil
}
