package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

const (
	DefaultEngineType = "in_memory"

	DefaultAddress        = "127.0.0.1:8080"
	DefaultMaxConnections = 30

	DefaultLogLevel  = "production"
	LogLevelDev      = "dev"
	DefaultLogOutput = "stdout"

	DefaultWalFlushingBatchSize      = 100
	DefaultWalFlushingBatchTimeoutMS = 1000
	DefaultWalMaxSegmentSize         = "10MB"
	DefaultWalDataDirectory          = "/tmp/kv-database/wal"
)

type EngineConfig struct {
	Type string `yaml:"type"`
}

type NetworkConfig struct {
	MaxConnections int    `yaml:"max_connections"`
	Address        string `yaml:"address"`
}

type LoggingConfig struct {
	Level  string `yaml:"level"`
	Output string `yaml:"output"`
}

type WalConfig struct {
	FlushingBatchSize      int    `yaml:"flushing_batch_size"`
	FlushingBatchTimeoutMS int    `yaml:"flushing_batch_timeout_ms"`
	MaxSegmentSize         string `yaml:"max_segment_size"`
	DataDirectory          string `yaml:"data_directory"`
}

type Config struct {
	Engine  EngineConfig  `yaml:"engine"`
	Network NetworkConfig `yaml:"network"`
	Logging LoggingConfig `yaml:"logging"`
	Wal     WalConfig     `yaml:"wal"`
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
		Engine: EngineConfig{
			Type: DefaultEngineType,
		},
		Network: NetworkConfig{
			MaxConnections: DefaultMaxConnections,
			Address:        DefaultAddress,
		},
		Logging: LoggingConfig{
			Level:  DefaultLogLevel,
			Output: DefaultLogOutput,
		},
		Wal: WalConfig{
			FlushingBatchSize:      DefaultWalFlushingBatchSize,
			FlushingBatchTimeoutMS: DefaultWalFlushingBatchTimeoutMS,
			MaxSegmentSize:         DefaultWalMaxSegmentSize,
			DataDirectory:          DefaultWalDataDirectory,
		},
	}
}
