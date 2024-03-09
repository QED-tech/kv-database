package database

import (
	"database/internal/database/commands"
	"database/internal/database/compute"
	"database/internal/database/storage"
	"database/internal/shared/logger"
)

//go:generate go run go.uber.org/mock/mockgen -package database -destination mock.go -source database.go Engine Analyzer Parser

type Parser interface {
	Parse(query string) ([]string, error)
}

type Analyzer interface {
	Analyze(compute.Tokens) (commands.Command, error)
}

type Engine interface {
	Execute(command commands.Command, writeToWAL bool) (storage.Result, error)
}

type Database struct {
	logger   logger.Logger
	engine   Engine
	analyzer Analyzer
	parser   Parser
}

func NewDatabase(
	logger logger.Logger,
	engine Engine,
	analyzer Analyzer,
	parser Parser,
) *Database {
	return &Database{
		logger:   logger,
		engine:   engine,
		analyzer: analyzer,
		parser:   parser,
	}
}

func (db *Database) Handle(input string) string {
	tokens, err := db.parser.Parse(input)
	if err != nil {
		db.logger.Infof("[database] error parsing query: %v", err)

		return err.Error()
	}

	db.logger.Infof(
		"[database] parsed tokens %v", tokens,
	)

	command, err := db.analyzer.Analyze(tokens)
	if err != nil {
		db.logger.Infof("[database] error analyze query: %v", err)

		return err.Error()
	}

	result, err := db.engine.Execute(command, true)
	if err != nil {
		db.logger.Errorf("[database] failed to execute query: %v", err)

		return err.Error()
	}

	return result.Out
}
