package wal

import (
	"bufio"
	"database/internal/database/commands"
	"strings"
)

type Reader struct {
	dataDirectory string
}

func NewReader(dataDirectory string) *Reader {
	return &Reader{dataDirectory: dataDirectory}
}

func (r *Reader) Read() <-chan commands.Command {
	out := make(chan commands.Command)

	segments, err := findSegments(r.dataDirectory)
	if err != nil {
		return out
	}

	go func() {
		defer func() {
			close(out)
		}()

		for _, segment := range segments.src {
			for log := range segment.read() {
				out <- mapLogToCommand(log)
			}
		}
	}()

	return out
}

func mapLogToCommand(log string) commands.Command {
	scan := bufio.NewScanner(strings.NewReader(log))
	scan.Split(bufio.ScanWords)

	var (
		index     int
		operation commands.Operation
		arguments []string
	)

	for scan.Scan() {
		switch index {
		case 0:
			operation = commands.Operation(scan.Text())
		default:
			arguments = append(arguments, scan.Text())
		}

		index++
	}

	return commands.Command{
		Operation: operation,
		Arguments: append([]string(nil), arguments...),
	}
}
