// Package model ...
package model

// Squitter is the squitter message.
type Squitter interface {
	Name() string
	AircraftAddress() ICAOAddr
}

// Message is a received message from extended squitter.
type Message interface {
	Name() string
}

// UntypeArray removes types from all the elements of an array.
func UntypeArray[T any](data []T) []any {
	output := make([]any, len(data))
	for idx, elt := range data {
		output[idx] = (any)(elt)
	}

	return output
}
