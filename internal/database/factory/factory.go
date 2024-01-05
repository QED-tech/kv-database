package factory

import (
	"database/internal/database"
	"database/internal/database/compute"
	"database/internal/database/config"
	"database/internal/database/storage"
	in_mem "database/internal/database/storage/in-mem"
	"database/internal/shared/logger"
	"fmt"
)

func CreateDatabase(
	logger logger.Logger,
	conf *config.Config,
) (*database.Database, error) {
	if conf == nil {
		return nil, fmt.Errorf("config should be defined")
	}

	engine := storage.NewEngine(
		getStorage(conf),
	)

	return database.NewDatabase(
		logger,
		engine,
		compute.NewAnalyzer(),
		compute.NewParser(),
	), nil
}

func getStorage(conf *config.Config) storage.Storage {
	switch conf.Engine.Type {
	case "in_memory":
		return in_mem.NewInMemoryStorage()
	}

	return in_mem.NewInMemoryStorage()
}
