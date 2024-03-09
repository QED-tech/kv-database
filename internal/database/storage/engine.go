package storage

import (
	"database/internal/database/commands"
	"database/internal/shared/logger"
	"fmt"
	"go.uber.org/zap"
)

type Storage interface {
	Get(key string) Result
	Set(key string, value string) Result
	Delete(key string) Result
}

type Wal interface {
	Run()
	WriteLog(cmd commands.Command) error
	ReadLogs() <-chan commands.Command
}

type Engine struct {
	storage Storage
	wal     Wal
	logger  logger.Logger
}

func NewEngine(storage Storage, wal Wal, log logger.Logger) *Engine {
	return &Engine{
		storage: storage,
		wal:     wal,
		logger:  log,
	}
}

type Result struct {
	Out string
}

func (e *Engine) Run() {
	e.logger.Debug("[engine] Is started")

	if e.wal != nil {
		e.wal.Run()
	}
}

func (e *Engine) Execute(cmd commands.Command, writeToWAL bool) (Result, error) {
	if err := cmd.Validate(); err != nil {
		return Result{}, err
	}

	if writeToWAL && e.wal != nil {
		if err := e.wal.WriteLog(cmd); err != nil {
			return Result{}, err
		}
	}

	switch cmd.Operation {
	case commands.SetOperation:
		return e.set(cmd)
	case commands.GetOperation:
		return e.get(cmd)
	case commands.DeleteOperation:
		return e.delete(cmd)
	}

	return Result{}, fmt.Errorf("unsupported operation %s", cmd.Operation)
}

func (e *Engine) TryRestore() {
	if e.wal == nil {
		return
	}

	e.logger.Debug("[engine] start trying to restore from wal")

	for cmd := range e.wal.ReadLogs() {
		result, err := e.Execute(cmd, false)
		if err != nil {
			e.logger.Warnf("[engine] failed to execute command from wal, err: %s", err.Error())
		}

		e.logger.Debug("[engine] execute command from wal", zap.String("out", result.Out))
	}
}

func (e *Engine) delete(cmd commands.Command) (Result, error) {
	return e.storage.Delete(
		cmd.GetKey(),
	), nil
}

func (e *Engine) set(cmd commands.Command) (Result, error) {
	return e.storage.Set(
		cmd.GetKey(),
		cmd.GetValue(),
	), nil
}

func (e *Engine) get(cmd commands.Command) (Result, error) {
	return e.storage.Get(
		cmd.GetKey(),
	), nil
}
