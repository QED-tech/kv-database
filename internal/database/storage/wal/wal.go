package wal

import (
	"database/internal/database/commands"
	"database/internal/shared/logger"
	"database/internal/tools"
	"sync"
	"time"
)

type LogsWriter interface {
	Write([]commands.Command) error
}

type LogsReader interface {
	Read() <-chan commands.Command
}

type Wal struct {
	logger     logger.Logger
	logsWriter LogsWriter
	logsReader LogsReader

	flushingBatchSize    int
	flushingBatchTimeout time.Duration
	batchCh              chan []commands.Command
	flushingErrCh        chan error

	mu        sync.Mutex
	logsBatch []commands.Command
}

func (w *Wal) Run() {
	w.logger.Debug("[wal] Is started")

	ticker := time.NewTicker(w.flushingBatchTimeout)
	go func() {
		defer ticker.Stop()

		for {
			select {
			case batch := <-w.batchCh:
				if err := w.logsWriter.Write(batch); err != nil {
					w.logger.Errorf("failed write log: %v", err)
					w.flushingErrCh <- err
				}

				w.flushingErrCh <- nil
			case <-ticker.C:
				tools.WithLock(&w.mu, func() {
					if err := w.logsWriter.Write(w.logsBatch); err != nil {
						w.logger.Errorf("failed write log: %v", err)
					}

					w.logsBatch = nil
				})
			}
		}
	}()
}

func NewWal(
	flushingBatchSize,
	flushingBatchTimeoutMS int,
	logger logger.Logger,
	writer LogsWriter,
	reader LogsReader,
) *Wal {
	wal := &Wal{
		flushingBatchSize:    flushingBatchSize,
		flushingBatchTimeout: time.Duration(flushingBatchTimeoutMS) * time.Millisecond,
		logger:               logger,
		batchCh:              make(chan []commands.Command),
		flushingErrCh:        make(chan error),
		logsWriter:           writer,
		logsReader:           reader,
	}

	return wal
}

func (w *Wal) ReadLogs() <-chan commands.Command {
	return w.logsReader.Read()
}

func (w *Wal) WriteLog(cmd commands.Command) error {
	if cmd.Operation == commands.GetOperation {
		return nil
	}

	var err error

	tools.WithLock(&w.mu, func() {
		w.logsBatch = append(w.logsBatch, cmd)

		if len(w.logsBatch) >= w.flushingBatchSize {
			w.logger.Debug("[wal] logs overflow")

			w.batchCh <- w.logsBatch
			err = <-w.flushingErrCh

			w.logsBatch = nil
		}
	})

	return err
}
