// Package database is a simple database with a cleaner when elements are too old.
package database

import (
	"context"
	"sync"
	"time"
)

// Storage is the database storage.
type Storage[K comparable, T any] struct {
	data     map[K][]Element[T]
	timer    *time.Ticker
	mutex    sync.Mutex
	lifetime time.Duration
}

// New initialize a storage.
func New[K comparable, T any](ctx context.Context, lifetime time.Duration) *Storage[K, T] {
	out := &Storage[K, T]{
		data:     map[K][]Element[T]{},
		lifetime: lifetime,
	}

	out.timer = time.NewTicker(time.Second)

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
