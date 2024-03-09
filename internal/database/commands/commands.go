package commands

import (
	"errors"
	"fmt"
	"strings"
)

type Operation string

const (
	DeleteOperation Operation = "DEL"
	GetOperation    Operation = "GET"
	SetOperation    Operation = "SET"
)

func (o Operation) String() string {
	return string(o)
}

func (o Operation) IsValid() bool {
	switch o {
	case DeleteOperation, GetOperation, SetOperation:
		return true
	default:
		return false
	}
}

func (o Operation) GetArgumentsCount() int {
	switch o {
	case DeleteOperation:
		return 1
	case GetOperation:
		return 1
	case SetOperation:
		return 2
	default:
		return 0
	}
}

type Command struct {
	Operation Operation
	Arguments []string
}

func (c Command) String() string {
	return fmt.Sprintf("%s %s", c.Operation, strings.Join(c.Arguments, " "))
}

func (c Command) GetKey() string {
	const keyIndex = 0

	return c.Arguments[keyIndex]
}

func (c Command) GetValue() string {
	const valIndex = 1

	return c.Arguments[valIndex]
}

var (
	ErrInvalidCommand        = errors.New("invalid command")
	ErrInvalidArgumentsCount = errors.New("invalid arguments count")
)

func NewCommand(operation Operation, arguments []string) Command {
	return Command{Operation: operation, Arguments: arguments}
}

func (c Command) Validate() error {
	if !c.Operation.IsValid() {
		return ErrInvalidCommand
	}

	if c.Operation.GetArgumentsCount() != len(c.Arguments) {
		return ErrInvalidArgumentsCount
	}

	return nil
}
