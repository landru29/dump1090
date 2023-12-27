package database

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/landru29/dump1090/cmd/logger"
)

// ChainedStorage is the database storage.
type ChainedStorage[K comparable, T any] struct {
	data       map[K]*ChainedElement[T]
	timer      *time.Ticker
	mutex      sync.Mutex
	lifetime   time.Duration
	cleanCycle time.Duration
	stop       chan struct{}
	waitgroup  *sync.WaitGroup
	closed     bool
}

// NewChainedStorage initialize a storage.
func NewChainedStorage[K comparable, T any]( //nolint: dupl
	ctx context.Context,
	opts ...ChainedConfigurator[K, T],
) *ChainedStorage[K, T] {
	out := &ChainedStorage[K, T]{
		data:       map[K]*ChainedElement[T]{},
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
func (s *ChainedStorage[K, T]) Clean() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for key, chain := range s.data {
		current := chain.Last()

		for current != nil && !current.expired(s.lifetime) {
			current = current.previous
		}

		if current == nil {
			continue
		}

		if current.next != nil {
			current.next.previous = nil
			s.data[key] = current.next

			continue
		}

		delete(s.data, key)
	}
}

// Add stores a new element.
func (s *ChainedStorage[K, T]) Add(key K, element T) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	chain, found := s.data[key]
	if found {
		last := chain.Last()
		last.next = &ChainedElement[T]{
			previous: last,
			data:     element,
			date:     time.Now(),
		}

		return
	}

	s.data[key] = &ChainedElement[T]{
		data: element,
		date: time.Now(),
	}
}

// Keys list the available keys.
func (s *ChainedStorage[K, T]) Keys() []K {
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

// Elements is the list of data for a specified key.
func (s *ChainedStorage[K, T]) Elements(key K) []T {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if elements, found := s.data[key]; found {
		out := []T{elements.data}

		for elements.next != nil {
			out = append(out, elements.data)

			elements = elements.next
		}

		return out
	}

	return nil
}

// Close implements the io.Closer interface.
func (s *ChainedStorage[K, T]) Close() error {
	if !s.closed {
		s.stop <- struct{}{}

		s.waitgroup.Wait()
	}

	return nil
}
