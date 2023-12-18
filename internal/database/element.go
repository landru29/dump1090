package database

import "time"

// Element is a database element.
type Element[T any] struct {
	date time.Time
	data T
}

// NewElement creates a new element.
func NewElement[T any](data T) Element[T] {
	return Element[T]{
		date: time.Now(),
		data: data,
	}
}

// Data is the element data.
func (e Element[T]) Data() T { //nolint: ireturn
	return e.data
}

func (e Element[T]) expired(lifetime time.Duration) bool {
	return time.Now().After(e.date.Add(lifetime))
}
