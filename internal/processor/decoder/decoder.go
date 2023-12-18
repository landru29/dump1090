// Package decoder is the default processor.
package decoder

import (
	"context"
	"time"

	"github.com/landru29/dump1090/internal/database"
	"github.com/landru29/dump1090/internal/modes"
)

// Process is the data processor.
type Process struct {
	ExtendedQuitters *database.Storage[modes.AircraftAddress, modes.ExtendedSquitter]
}

// New creates a data processor.
func New(ctx context.Context, dbLifeTime time.Duration) *Process {
	return &Process{
		ExtendedQuitters: database.New[modes.AircraftAddress, modes.ExtendedSquitter](ctx, dbLifeTime),
	}
}

// Process implements source.Processor the interface.
func (p Process) Process(data []byte) error {
	message := &modes.ExtendedSquitter{}
	if err := message.Unmarshal(data); err != nil {
		return err
	}

	p.ExtendedQuitters.Add(message.AircraftAddress, *message)

	return nil
}
