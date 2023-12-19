package database

import "time"

// Configurator is the database configurator.
type Configurator[K comparable, T any] func(*Storage[K, T])

// WithLifetime specify the lifetime of the elements.
func WithLifetime[K comparable, T any](duration time.Duration) Configurator[K, T] {
	return func(s *Storage[K, T]) {
		s.lifetime = duration
	}
}

// WithCleanCycle specify the lifetime of the elements.
func WithCleanCycle[K comparable, T any](duration time.Duration) Configurator[K, T] {
	return func(s *Storage[K, T]) {
		s.cleanCycle = duration
	}
}
