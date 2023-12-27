package database

import "time"

// ChainedElement is a database element.
type ChainedElement[T any] struct {
	date     time.Time
	data     T
	next     *ChainedElement[T]
	previous *ChainedElement[T]
}

// Last is the last Element of the chain.
func (e *ChainedElement[T]) Last() *ChainedElement[T] {
	output := e
	for output.next != nil {
		output = output.next
	}

	return output
}

// Root is the first Element of the chain.
func (e ChainedElement[T]) Root() *ChainedElement[T] {
	output := &e
	for e.previous != nil {
		output = e.previous
	}

	return output
}

// NewElement creates a new element.
// func NewChainedElement[T any](data T) ChainedElement[T] {
// 	return ChainedElement[T]{
// 		date: time.Now(),
// 		data: data,
// 	}
// }

// Data is the element data.
func (e ChainedElement[T]) Data() T { //nolint: ireturn
	return e.data
}

func (e ChainedElement[T]) expired(lifetime time.Duration) bool {
	return time.Now().After(e.date.Add(lifetime))
}
