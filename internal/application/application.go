// Package application is the main application.
package application

import (
	"context"
	"log/slog"
	"time"

	"github.com/landru29/dump1090/internal/database"
	"github.com/landru29/dump1090/internal/model"
	"github.com/landru29/dump1090/internal/processor"
	"github.com/landru29/dump1090/internal/source"
	"github.com/landru29/dump1090/internal/source/rtl28xxx"
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
	starter    source.Starter
	log        *slog.Logger
	aircraftDB *database.ElementStorage[model.ICAOAddr, model.Aircraft]
	messageDB  *database.ChainedStorage[model.ICAOAddr, model.Squitter]
}

// New creates a new application.
func New(
	log *slog.Logger,
	cfg *Config,
	processors []processor.Processer,
	aircraftDB *database.ElementStorage[model.ICAOAddr, model.Aircraft],
	messageDB *database.ChainedStorage[model.ICAOAddr, model.Squitter],
) (*App, error) {
	output := &App{
		log:        log,
		aircraftDB: aircraftDB,
		messageDB:  messageDB,
	}

	switch {
	case cfg.FixturesFilename != "":
		// Source is a file
		opts := []rtl28xxx.FileConfigurator{}
		if cfg.FixtureLoop {
			opts = append(opts, rtl28xxx.WithLoop())
		}

		output.starter = rtl28xxx.NewFile(cfg.FixturesFilename, processors, opts...)

		return output, nil
	default:
		// Source is the device
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
}

// Start is the application entrypoint.
func (a *App) Start(ctx context.Context) error {
	a.log.Info("Starting application")

	return a.starter.Start(ctx)
}
