package wal

import (
	"bufio"
	"database/internal/database/commands"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"sort"
	"strconv"
)

type segment struct {
	ID   int
	path string
}

func newSegment(ID int, dataDirectory string) *segment {
	return &segment{ID: ID, path: createSegmentPath(ID, dataDirectory)}
}

var ErrSegmentGreaterThanMaxSize = errors.New("segment greater than max size")

func (s *segment) createNext(dataDirectory string) *segment {
	return &segment{
		ID:   s.ID + 1,
		path: createSegmentPath(s.ID+1, dataDirectory),
	}
}

func (s *segment) read() <-chan string {
	out := make(chan string)

	file, err := os.Open(s.path)
	if err != nil {
		return out
	}

	scan := bufio.NewScanner(file)

	go func() {
		defer func() {
			close(out)
			file.Close()
		}()

		for scan.Scan() {
			out <- scan.Text()
		}
	}()

	return out
}

func (s *segment) write(log commands.Command, maxSize int64) error {
	flags := os.O_APPEND | os.O_CREATE | os.O_WRONLY

	file, err := os.OpenFile(s.path, flags, fs.FileMode(0644))
	if err != nil {
		return fmt.Errorf("failed to open segment file: %w", err)
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return err
	}

	if stat.Size() > maxSize {
		return ErrSegmentGreaterThanMaxSize
	}

	if _, err = file.WriteString(log.String() + "\n"); err != nil {
		return err
	}

	if err = file.Sync(); err != nil {
		return err
	}

	return nil
}

type Segments struct {
	src []segment
}

func (s *Segments) isEmpty() bool {
	return len(s.src) == 0
}

func (s *Segments) sortASC() *Segments {
	copySrc := append([]segment(nil), s.src...)

	sort.Slice(copySrc, func(prev, next int) bool {
		return copySrc[prev].ID < copySrc[next].ID
	})

	return &Segments{
		src: copySrc,
	}
}

func (s *Segments) last() *segment {
	if s.isEmpty() {
		return nil
	}

	return &s.src[len(s.src)-1]
}

func createSegmentPath(segmentID int, dataDirectory string) string {
	return path.Join(dataDirectory, fmt.Sprintf("wal-%d.txt", segmentID))
}

func findSegments(dataDirectory string) (*Segments, error) {
	entries, err := os.ReadDir(dataDirectory)
	if err != nil {
		return nil, fmt.Errorf("failed to read wal directory")
	}

	segments := make([]segment, 0)
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		parseSegmentID := onlyDigitsRegexp.FindString(entry.Name())
		segID, err := strconv.Atoi(parseSegmentID)
		if err != nil {
			return nil, fmt.Errorf("failed to cast segmentID to int: %w", err)
		}

		seg := segment{
			ID:   segID,
			path: path.Join(dataDirectory, entry.Name()),
		}

		segments = append(segments, seg)
	}

	return &Segments{
		src: segments,
	}, nil
}
