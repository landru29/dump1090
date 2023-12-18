// Package source is the data source.
package source

import "context"

// Starter is a process starter.
type Starter interface {
	Start(ctx context.Context) error
}
