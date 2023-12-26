package database

import "time"

// Element is a database element.
type Element[T any] struct {
	date     time.Time
	data     T
	next     *Element[T]
	previous *Element[T]
}

// Last is the last Element of the chain.
func (e *Element[T]) Last() *Element[T] {
	output := e
	for output.next != nil {
		output = output.next
	}

	return output
}

// Root is the first Element of the chain.
func (e Element[T]) Root() *Element[T] {
	output := &e
	for e.previous != nil {
		output = e.previous
	}

	return output
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
