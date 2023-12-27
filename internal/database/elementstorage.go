package database

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/landru29/dump1090/cmd/logger"
)

// ElementStorage is the database storage.
type ElementStorage[K comparable, T any] struct {
	data       map[K]*Element[T]
	timer      *time.Ticker
	mutex      sync.Mutex
	lifetime   time.Duration
	cleanCycle time.Duration
	stop       chan struct{}
	waitgroup  *sync.WaitGroup
	closed     bool
}

// NewElementStorage initialize a storage.
func NewElementStorage[K comparable, T any]( //nolint: dupl
	ctx context.Context,
	opts ...ElementConfigurator[K, T],
) *ElementStorage[K, T] {
	out := &ElementStorage[K, T]{
		data:       map[K]*Element[T]{},
		lifetime:   defaultLifetime,
		cleanCycle: defaultCleanCycle,
		stop:       make(chan struct{}),
		waitgroup:  &sync.WaitGroup{},
	}

	for _, opt := range opts {
		opt(out)
	}

	out.timer = time.NewTicker(out.cleanCycle)

	out.waitgroup.Add(1)

	go func() {
		log, found := logger.Logger(ctx)
		if found {
			log = log.With("database", "chained", "type", fmt.Sprintf("%T", new(T)))

			log.Info("starting cleaner")
		}

		defer func() {
			out.timer.Stop()

			out.closed = true

			out.waitgroup.Done()

			if found {
				log.Info("stopping cleaner")
			}
		}()

		for {
			select {
			case <-out.timer.C:
				out.Clean()
			case <-ctx.Done():
				return
			case <-out.stop:
				return
			}
		}
	}()

	return out
}

// Clean removes outdated elements.
func (s *ElementStorage[K, T]) Clean() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for key, element := range s.data {
		if element.expired(s.lifetime) {
			delete(s.data, key)
		}
	}
}

// Add stores a new element.
func (s *ElementStorage[K, T]) Add(key K, element T) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.data[key] = &Element[T]{
		data: element,
		date: time.Now(),
	}
}

// Keys list the available keys.
func (s *ElementStorage[K, T]) Keys() []K {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	out := make([]K, len(s.data))

	idx := 0

	for key := range s.data {
		out[idx] = key

		idx++
	}

	return out
}

// Element is the list of data for a specified key.
func (s *ElementStorage[K, T]) Element(key K) *T {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if element, found := s.data[key]; found {
		return &element.data
	}

	return nil
}

// Close implements the io.Closer interface.
func (s *ElementStorage[K, T]) Close() error {
	if !s.closed {
		s.stop <- struct{}{}

		s.waitgroup.Wait()
	}

	return nil
}
