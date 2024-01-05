package tools

import (
	"sync"
)

type Semaphore struct {
	cond    *sync.Cond
	current int
	maxSize int
}

func NewSemaphore(maxSize int) *Semaphore {
	return &Semaphore{
		cond:    sync.NewCond(&sync.Mutex{}),
		maxSize: maxSize,
	}
}

func (s *Semaphore) Acquire() {
	s.cond.L.Lock()
	defer s.cond.L.Unlock()

	for s.current >= s.maxSize {
		s.cond.Wait()
	}

	s.current++
}

func (s *Semaphore) Release() {
	s.cond.L.Lock()
	defer s.cond.L.Unlock()

	s.current--
	s.cond.Broadcast()
}
