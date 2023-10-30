// Package application is the main application.
package application

import (
	"context"
	"log"

	"github.com/landru29/dump1090/internal/dump"
	"github.com/landru29/dump1090/internal/transport"
)

const (
	acChannelSize = 10
)

// Config is the application configuration.
type Config struct {
	FixturesFilename string
	DeviceIndex      uint32
	Frequency        uint32
	Gain             int
	EnableAGC        bool
}

// App is the main application.
type App struct {
	cfg         *Config
	tranporters []transport.Transporter
}

// New creates a new application.
func New(cfg *Config, tranporters []transport.Transporter) (*App, error) {
	return &App{
		cfg:         cfg,
		tranporters: tranporters,
	}, nil
}

// Start is the application entrypoint.
func (a *App) Start(ctx context.Context, loop bool) error {
	eventAircraft := make(chan *dump.Aircraft, acChannelSize)
	defer func() {
		close(eventAircraft)
	}()

	go func(acStream chan *dump.Aircraft) {
		for {
			ac := <-acStream

			for _, transporter := range a.tranporters {
				if err := transporter.Transport(ac); err != nil {
					log.Println(err)
				}
			}
		}
	}(eventAircraft)

	return dump.Start(
		ctx,
		a.cfg.DeviceIndex,
		a.cfg.Gain,
		a.cfg.Frequency,
		a.cfg.EnableAGC,
		a.cfg.FixturesFilename,
		nil,
		eventAircraft,
		loop,
	)
}
