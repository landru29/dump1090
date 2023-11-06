// Package source is the data source.
package source

import "context"

// Starter is a process starter.
type Starter interface {
	Start(ctx context.Context) error
}

// Processer is a data processor.
type Processer interface {
	Process(data []byte)
}
