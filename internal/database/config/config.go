package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

type EngineConfig struct {
	Type string `yaml:"type,omitempty"`
}

type NetworkConfig struct {
	MaxConnections int    `yaml:"max_connections,omitempty"`
	Address        string `yaml:"address,omitempty"`
}

type LoggingConfig struct {
	Level  string `yaml:"level,omitempty"`
	Output string `yaml:"output,omitempty"`
}

type Config struct {
	Engine  *EngineConfig  `yaml:"engine,omitempty"`
	Network *NetworkConfig `yaml:"network,omitempty"`
	Logging *LoggingConfig `yaml:"logging,omitempty"`
}

func ReadConfig() (*Config, error) {
	conf, err := os.ReadFile("./internal/database/config/config.yaml")
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	c := defaultConfig()

	err = yaml.Unmarshal(conf, &c)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &c, nil
}

func defaultConfig() Config {
	return Config{
		Engine:  &EngineConfig{Type: "in_memory"},
		Network: &NetworkConfig{MaxConnections: 30, Address: "127.0.0.1:8080"},
		Logging: &LoggingConfig{Level: "production", Output: "stdout"},
	}
}
