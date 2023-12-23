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
	ExtendedSquitters *database.Storage[model.ICAOAddr, model.Squitter]
}

// New creates a data processor.
func New(ctx context.Context, dbLifeTime time.Duration) *Process {
	return &Process{
		ExtendedSquitters: database.New[model.ICAOAddr, model.Squitter](
			ctx,
			database.WithLifetime[model.ICAOAddr, model.Squitter](dbLifeTime),
		),
	}
}

// Process implements source.Processor the interface.
func (p Process) Process(data []byte) error {
	modes := model.ModeS(data)

	if err := modes.CheckSum(); err != nil {
		return err
	}

	squitter, err := modes.Squitter()
	if err != nil {
		return err
	}

	p.ExtendedSquitters.Add(squitter.AircraftAddress(), squitter)

	return nil
}
