package database

import "time"

// Element is a database element.
type Element[T any] struct {
	date time.Time
	data T
}

// Data is the element data.
func (e Element[T]) Data() T { //nolint: ireturn
	return e.data
}

func (e Element[T]) expired(lifetime time.Duration) bool {
	return time.Now().After(e.date.Add(lifetime))
}
