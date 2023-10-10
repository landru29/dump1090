// Package errors manages the text errors.
package errors

// Error is a text error.
type Error string

// Error implements the error interface.
func (e Error) Error() string {
	return string(e)
}
