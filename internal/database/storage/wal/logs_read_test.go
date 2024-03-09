package wal

import (
	"database/internal/database/commands"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	testWALDataDir      = "./testdata"
	testEmptyWALDataDir = "./testdata-empty"
)

func TestReaderRead(t *testing.T) {
	t.Run("test reading wal logs", func(t *testing.T) {
		r := &Reader{
			dataDirectory: testWALDataDir,
		}

		results := make([]commands.Command, 0)
		expected := []commands.Command{
			{
				Operation: commands.SetOperation,
				Arguments: []string{"key", "value"},
			},
			{
				Operation: commands.DeleteOperation,
				Arguments: []string{"key"},
			},
			{
				Operation: commands.SetOperation,
				Arguments: []string{"key1", "val"},
			},
			{
				Operation: commands.DeleteOperation,
				Arguments: []string{"key1"},
			},
		}

		is := require.New(t)

		got := r.Read()

		for log := range got {
			results = append(results, log)
		}

		is.Equal(expected, results)
	})
}

func TestReaderReadWhenEmptyLogs(t *testing.T) {
	t.Run("test reading wal logs", func(t *testing.T) {
		r := &Reader{
			dataDirectory: testEmptyWALDataDir,
		}

		results := make([]commands.Command, 0)
		expected := make([]commands.Command, 0)

		is := require.New(t)

		got := r.Read()

		for log := range got {
			results = append(results, log)
		}

		is.Equal(expected, results)
	})
}
