package factory

import (
	"database/internal/database"
	"database/internal/database/compute"
	"database/internal/database/config"
	"database/internal/database/storage"
	in_mem "database/internal/database/storage/in-mem"
	"database/internal/database/storage/wal"
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
		wal.NewWal(
			conf.Wal.FlushingBatchSize,
			conf.Wal.FlushingBatchTimeoutMS,
			logger,
			wal.NewWriter(
				parseSizeToBytes(conf.Wal.MaxSegmentSize),
				conf.Wal.DataDirectory,
			),
			wal.NewReader(conf.Wal.DataDirectory),
		),
		logger,
	)

	engine.TryRestore()
	engine.Run()

	return database.NewDatabase(
		logger,
		engine,
		compute.NewAnalyzer(),
		compute.NewParser(),
	), nil
}

func parseSizeToBytes(_ string) int64 {
	return 1024
}

func getStorage(conf *config.Config) storage.Storage {
	switch conf.Engine.Type {
	case config.DefaultEngineType:
		return in_mem.NewInMemoryStorage()
	}

	return in_mem.NewInMemoryStorage()
}
