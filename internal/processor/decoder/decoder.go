// Package decoder is the default processor.
package decoder

import (
	"context"
	"time"

	"github.com/landru29/dump1090/internal/database"
	"github.com/landru29/dump1090/internal/model"
)

// Process is the data processor.
type Process struct {
	ExtendedSquitters *database.Storage[model.ICAOAddr, model.ExtendedSquitter]
}

// New creates a data processor.
func New(ctx context.Context, dbLifeTime time.Duration) *Process {
	return &Process{
		ExtendedSquitters: database.New[model.ICAOAddr, model.ExtendedSquitter](
			ctx,
			database.WithLifetime[model.ICAOAddr, model.ExtendedSquitter](dbLifeTime),
		),
	}
}

// Process implements source.Processor the interface.
func (p Process) Process(data []byte) error {
	message := &model.ExtendedSquitter{}
	if err := message.UnmarshalModeS(data); err != nil {
		return err
	}

	p.ExtendedSquitters.Add(message.AircraftAddress, *message)

	return nil
}
