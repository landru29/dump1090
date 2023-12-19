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
	ExtendedSquitters *database.Storage[modes.AircraftAddress, modes.ExtendedSquitter]
}

// New creates a data processor.
func New(ctx context.Context, dbLifeTime time.Duration) *Process {
	return &Process{
		ExtendedSquitters: database.New[modes.AircraftAddress, modes.ExtendedSquitter](
			ctx,
			database.WithLifetime[modes.AircraftAddress, modes.ExtendedSquitter](dbLifeTime),
		),
	}
}

// Process implements source.Processor the interface.
func (p Process) Process(data []byte) error {
	message := &modes.ExtendedSquitter{}
	if err := message.Unmarshal(data); err != nil {
		return err
	}

	p.ExtendedSquitters.Add(message.AircraftAddress, *message)

	return nil
}
