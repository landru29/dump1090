// Package decoder is the default processor.
package decoder

import (
	"context"
	"log/slog"
	"time"

	"github.com/landru29/dump1090/internal/database"
	"github.com/landru29/dump1090/internal/model"
	"github.com/landru29/dump1090/internal/transport"
)

// Configurator is the Process configurator.
type Configurator func(*Process)

// Process is the data processor.
type Process struct {
	ExtendedSquitters *database.ChainedStorage[model.ICAOAddr, model.Squitter]
	log               *slog.Logger
	checkCRC          bool
	dbLifeTime        time.Duration
	transporters      []transport.Transporter
}

// New creates a data processor.
func New(ctx context.Context, log *slog.Logger, opts ...Configurator) *Process {
	process := &Process{
		log:          log,
		transporters: []transport.Transporter{},
	}

	for _, opt := range opts {
		opt(process)
	}

	process.ExtendedSquitters = database.NewChainedStorage[model.ICAOAddr, model.Squitter](
		ctx,
		database.ChainedWithLifetime[model.ICAOAddr, model.Squitter](process.dbLifeTime),
	)

	return process
}

// WithChecksumCheck requests the CRC check.
func WithChecksumCheck() Configurator {
	return func(process *Process) {
		process.checkCRC = true
	}
}

// WithDatabaseLifetime sets the lifetime of database elements.
func WithDatabaseLifetime(dbLifeTime time.Duration) Configurator {
	return func(process *Process) {
		process.dbLifeTime = dbLifeTime
	}
}

// WithTransporter add a new transporter.
func WithTransporter(transporter transport.Transporter) Configurator {
	return func(process *Process) {
		process.transporters = append(process.transporters, transporter)
	}
}

// Process implements source.Processor the interface.
func (p Process) Process(data []byte) error {
	modes := model.ModeS(data)

	log := p.log.With("message", modes.String())

	if p.checkCRC {
		if err := modes.CheckSum(); err != nil {
			log.Error("checksum", "msg", err)

			return err
		}
	}

	squitter, err := modes.Squitter()
	if err != nil {
		log.Error("squitter", "msg", err)

		return err
	}

	log.Info("processing message")

	p.ExtendedSquitters.Add(squitter.AircraftAddress(), squitter)

	for _, transporter := range p.transporters {
		if err := transporter.Transport(nil); err != nil {
			log.Error("transport", "msg", err)
		}
	}

	return nil
}
