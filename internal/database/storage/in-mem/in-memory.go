package in_mem

import (
	"database/internal/database/storage"
	"fmt"
	"sync"
)

type InMemoryStorage struct {
	mu  sync.Mutex
	src map[string]string
}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		src: make(map[string]string),
	}
}

func (s *InMemoryStorage) Get(key string) storage.Result {
	s.mu.Lock()
	val, ok := s.src[key]
	s.mu.Unlock()

	if !ok {
		return storage.Result{
			Out: fmt.Sprintf("[database] value by key '%s' not found", key),
		}
	}

	return storage.Result{
		Out: val,
	}
}
func (s *InMemoryStorage) Set(key string, value string) storage.Result {
	s.mu.Lock()
	s.src[key] = value
	s.mu.Unlock()

	return storage.Result{
		Out: "[database] OK",
	}
}
func (s *InMemoryStorage) Delete(key string) storage.Result {
	s.mu.Lock()
	delete(s.src, key)
	s.mu.Unlock()

	return storage.Result{
		Out: "[database] OK",
	}
}
