// Package application is the main application.
package application

import (
	"context"
	"time"

	"github.com/landru29/dump1090/internal/processor"
	"github.com/landru29/dump1090/internal/processor/decoder"
	"github.com/landru29/dump1090/internal/source"
	"github.com/landru29/dump1090/internal/source/rtl28xxx"
	"github.com/landru29/dump1090/internal/transport"
)

// Config is the application configuration.
type Config struct {
	FixturesFilename string
	FixtureLoop      bool
	DeviceIndex      uint32
	Frequency        uint32
	Gain             float64
	EnableAGC        bool
	DatabaseLifetime time.Duration
}

// App is the main application.
type App struct {
	starter source.Starter
}

// New creates a new application.
func New(ctx context.Context, cfg *Config, _ []transport.Transporter) (*App, error) {
	output := &App{}

	var processor processor.Processer = decoder.New(ctx, cfg.DatabaseLifetime)

	if cfg.FixturesFilename != "" {
		opts := []rtl28xxx.FileConfigurator{}
		if cfg.FixtureLoop {
			opts = append(opts, rtl28xxx.WithLoop())
		}

		output.starter = rtl28xxx.NewFile(cfg.FixturesFilename, processor, opts...)

		return output, nil
	}

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

	output.starter = rtl28xxx.New(opts...)

	return output, nil
}

// Start is the application entrypoint.
func (a *App) Start(ctx context.Context) error {
	return a.starter.Start(ctx)
}
