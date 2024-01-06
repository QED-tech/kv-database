package storage

import (
	"database/internal/database/commands"
	"fmt"
)

type Storage interface {
	Get(key string) Result
	Set(key string, value string) Result
	Delete(key string) Result
}

type Result struct {
	Out string
}

type Engine struct {
	storage Storage
}

func NewEngine(storage Storage) *Engine {
	return &Engine{
		storage: storage,
	}
}

func (e *Engine) Execute(cmd commands.Command) (Result, error) {
	if err := cmd.Validate(); err != nil {
		return Result{}, err
	}

	switch cmd.Operation {
	case commands.SetOperation:
		return e.storage.Set(
			cmd.GetKey(),
			cmd.GetValue(),
		), nil
	case commands.GetOperation:
		return e.storage.Get(
			cmd.GetKey(),
		), nil
	case commands.DeleteOperation:
		return e.storage.Delete(
			cmd.GetKey(),
		), nil
	}

	return Result{}, fmt.Errorf("unsupported operation %s", cmd.Operation)
}
