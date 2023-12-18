// Package empty is an empty processor.
package empty

// New creates a ampty processor.
func New() *Processor {
	return &Processor{}
}

// Processor is an empty processor.
type Processor struct{}

// Process implements the Processer interface.
func (e Processor) Process(_ []byte) {}
