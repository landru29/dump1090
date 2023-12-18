// Package processor defines the way to process input data.
package processor

// Processer is a data processor.
type Processer interface {
	Process(data []byte) error
}
