package compute

import (
	"database/internal/database/commands"
	"errors"
	"fmt"
)

var (
	ErrMinTokensLen = errors.New(
		fmt.Sprintf("[database] minimum tokens length for execute command - %d", minimumTokensLen),
	)
	ErrMaxTokensLen = errors.New(
		fmt.Sprintf("[database] maximum tokens length for execute command - %d", maxTokensLen),
	)
)

type Analyzer struct {
}

func NewAnalyzer() *Analyzer {
	return &Analyzer{}
}

type Tokens []string

func (t Tokens) Len() int { return len(t) }

const minimumTokensLen = 2
const maxTokensLen = 3

func (t Tokens) Validate() error {
	if t.Len() < minimumTokensLen {
		return ErrMinTokensLen
	}

	if t.Len() > maxTokensLen {
		return ErrMaxTokensLen
	}

	return nil
}

func (a *Analyzer) Analyze(tokens Tokens) (commands.Command, error) {
	if err := tokens.Validate(); err != nil {
		return commands.Command{}, err
	}

	cmd := commands.NewCommand(commands.Operation(tokens[0]), tokens[1:])

	if err := cmd.Validate(); err != nil {
		return commands.Command{}, err
	}

	return cmd, nil
}
