// Package model gather all business models.
package model

// UntypeArray removes types from all the elements of an array.
func UntypeArray[T any](data []T) []any {
	output := make([]any, len(data))
	for idx, elt := range data {
		output[idx] = (any)(elt)
	}

	return output
}
