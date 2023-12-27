package database

import "time"

// ChainedConfigurator is the database configurator for chained storage.
type ChainedConfigurator[K comparable, T any] func(*ChainedStorage[K, T])

// ChainedWithLifetime specify the lifetime of the elements.
func ChainedWithLifetime[K comparable, T any](duration time.Duration) ChainedConfigurator[K, T] {
	return func(s *ChainedStorage[K, T]) {
		s.lifetime = duration
	}
}

// ChainedWithCleanCycle specify the lifetime of the elements.
func ChainedWithCleanCycle[K comparable, T any](duration time.Duration) ChainedConfigurator[K, T] {
	return func(s *ChainedStorage[K, T]) {
		s.cleanCycle = duration
	}
}

// ElementConfigurator is the database configurator for chained storage.
type ElementConfigurator[K comparable, T any] func(*ElementStorage[K, T])

// ElementWithLifetime specify the lifetime of the elements.
func ElementWithLifetime[K comparable, T any](duration time.Duration) ElementConfigurator[K, T] {
	return func(s *ElementStorage[K, T]) {
		s.lifetime = duration
	}
}

// ElementWithCleanCycle specify the lifetime of the elements.
func ElementWithCleanCycle[K comparable, T any](duration time.Duration) ElementConfigurator[K, T] {
	return func(s *ElementStorage[K, T]) {
		s.cleanCycle = duration
	}
}
