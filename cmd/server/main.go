package main

import (
	"database/internal/database/config"
	"database/internal/database/factory"
	"database/internal/network"
	"database/internal/shared/logger"
	"log"

	"go.uber.org/zap"
)

func main() {
	conf, err := config.ReadConfig()
	if err != nil {
		log.Fatalf("failed to load config, err: %v", err)
	}

	l, err := logger.NewLogger(conf)
	if err != nil {
		log.Fatalf("failed to create logger, err: %v", err)
	}

	database, err := factory.CreateDatabase(l, conf)
	if err != nil {
		log.Fatalf("failed to create database, err: %v", err)
	}

	l.Debug("config: ", zap.Reflect("config", *conf))
	s := network.NewTCPServer(database, l, conf)

	if err := s.Listen(); err != nil {
		log.Fatalf("failed to listen server, err: %v", err)
	}
}
