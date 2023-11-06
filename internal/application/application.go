// Package application is the main application.
package application

import (
	"context"

	"github.com/landru29/dump1090/internal/modes"
	"github.com/landru29/dump1090/internal/source"
	"github.com/landru29/dump1090/internal/source/file"
	"github.com/landru29/dump1090/internal/source/rtl28xxx"
	"github.com/landru29/dump1090/internal/transport"
)

const (
	acChannelSize = 10
)

// Config is the application configuration.
type Config struct {
	FixturesFilename string
	FixtureLoop      bool
	DeviceIndex      uint32
	Frequency        uint32
	Gain             float64
	EnableAGC        bool
}

// App is the main application.
type App struct {
	starter source.Starter
}

// New creates a new application.
func New(cfg *Config, tranporters []transport.Transporter) (*App, error) {
	output := &App{}

	if cfg.FixturesFilename != "" {
		opts := []file.Configurator{}
		if cfg.FixtureLoop {
			opts = append(opts, file.WithLoop())
		}

		output.starter = file.New(cfg.FixturesFilename, modes.New(tranporters), opts...)
	} else {
		rtl28xxx.InitTables()

		opts := []rtl28xxx.Configurator{}

		if cfg.DeviceIndex > 0 {
			opts = append(opts, rtl28xxx.WithDeviceIndex(int(cfg.DeviceIndex)))
		}

		if cfg.EnableAGC {
			opts = append(opts, rtl28xxx.WithAGC())
		}

		if cfg.Frequency > 0 {
			opts = append(opts, rtl28xxx.WithFrequency(cfg.Frequency))
		}

		if cfg.Gain > 0 {
			opts = append(opts, rtl28xxx.WithGain(cfg.Gain))
		}

		output.starter = rtl28xxx.New(modes.New(tranporters), opts...)
	}

	return output, nil
}

// Start is the application entrypoint.
func (a *App) Start(ctx context.Context, loop bool) error {
	return a.starter.Start(ctx)
}
