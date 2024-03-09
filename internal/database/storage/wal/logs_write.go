package wal

import (
	"database/internal/database/commands"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"regexp"
)

type Writer struct {
	maxSegmentSizeBytes int64
	dataDirectory       string
}

func NewWriter(maxSegmentSizeBytes int64, dataDirectory string) *Writer {
	return &Writer{
		maxSegmentSizeBytes: maxSegmentSizeBytes,
		dataDirectory:       dataDirectory,
	}
}

func (w *Writer) Write(logs []commands.Command) error {
	if err := w.createSegmentsDir(); err != nil {
		return fmt.Errorf("failed to create directory for wal: %w", err)
	}

	seg, err := w.getLastSegment()
	if err != nil {
		return fmt.Errorf("failed to create segment file for wal: %w", err)
	}

	for _, log := range logs {
		if err = seg.write(log, w.maxSegmentSizeBytes); err != nil {
			if errors.Is(err, ErrSegmentGreaterThanMaxSize) {
				if err = seg.createNext(w.dataDirectory).write(log, w.maxSegmentSizeBytes); err != nil {
					return err
				}
				continue
			}

			return err
		}
	}

	return nil
}

func (w *Writer) createSegmentsDir() error {
	return os.MkdirAll(w.dataDirectory, fs.FileMode(0755))
}

var onlyDigitsRegexp = regexp.MustCompile("[0-9]+")

const (
	firstSegmentID = 1
)

func (w *Writer) getLastSegment() (*segment, error) {
	segments, err := findSegments(w.dataDirectory)
	if err != nil {
		return nil, err
	}

	var seg *segment

	seg = segments.sortASC().last()
	if seg == nil {
		return newSegment(firstSegmentID, w.dataDirectory), nil
	}

	return seg, nil
}
