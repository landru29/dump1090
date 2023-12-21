// Package database is a simple database with a cleaner when elements are too old.
package database

import (
	"context"
	"sync"
	"time"
)

const (
	defaultLifetime   time.Duration = time.Minute
	defaultCleanCycle time.Duration = time.Minute / 3
)

// Storage is the database storage.
type Storage[K comparable, T any] struct {
	data       map[K][]Element[T]
	timer      *time.Ticker
	mutex      sync.Mutex
	lifetime   time.Duration
	cleanCycle time.Duration
}

// New initialize a storage.
func New[K comparable, T any](ctx context.Context, opts ...Configurator[K, T]) *Storage[K, T] {
	out := &Storage[K, T]{
		data:       map[K][]Element[T]{},
		lifetime:   defaultLifetime,
		cleanCycle: defaultCleanCycle,
	}

	for _, opt := range opts {
		opt(out)
	}

	out.timer = time.NewTicker(out.cleanCycle)

	go func() {
		defer out.timer.Stop()

		for {
			select {
			case <-out.timer.C:
				out.Clean()
			case <-ctx.Done():
				return
			}
		}
	}()

	return out
}

// Clean removes outdated elements.
func (s *Storage[K, T]) Clean() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for key, list := range s.data {
		newLst := []Element[T]{}

		for _, elt := range list {
			if !elt.expired(s.lifetime) {
				newLst = append(newLst, elt)
			}
		}

		s.data[key] = newLst

		if len(s.data[key]) == 0 {
			delete(s.data, key)
		}
	}
}

// Add stores a new element.
func (s *Storage[K, T]) Add(key K, element T) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	list, found := s.data[key]
	if !found {
		list = []Element[T]{}
	}

	list = append(list, NewElement[T](element))

	s.data[key] = list
}

// Keys list the available keys.
func (s *Storage[K, T]) Keys() []K {
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
func (s *Storage[K, T]) Elements(key K) []T {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if elements, found := s.data[key]; found {
		out := make([]T, len(elements))

		for idx, elt := range elements {
			out[idx] = elt.Data()
		}

		return out
	}

	return nil
}
