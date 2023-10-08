package application

import (
	"context"
	"log"

	"github.com/landru29/dump1090/internal/dump"
	"github.com/landru29/dump1090/internal/serialize"
	"github.com/landru29/dump1090/internal/transport"
)

type Config struct {
	FixturesFilename string
	DeviceIndex      uint32
	Frequency        uint32
	Gain             int
	EnableAGC        bool
}

type App struct {
	cfg         *Config
	formater    serialize.Serializer
	tranporters []transport.Transporter
}

func New(cfg *Config, formater serialize.Serializer, tranporters []transport.Transporter) (*App, error) {
	return &App{
		cfg:         cfg,
		formater:    formater,
		tranporters: tranporters,
	}, nil
}

func (a *App) Start(ctx context.Context) error {
	eventAircraft := make(chan *dump.Aircraft, 10)
	defer func() {
		close(eventAircraft)
	}()

	go func(acStream chan *dump.Aircraft, app *App) {
		for {
			ac := <-acStream

			for _, transporter := range a.tranporters {
				if err := transporter.Transport(ac); err != nil {
					log.Println(err)
				}
			}
		}
	}(eventAircraft, a)

	return dump.Start(
		ctx,
		a.cfg.DeviceIndex,
		a.cfg.Gain,
		a.cfg.Frequency,
		a.cfg.EnableAGC,
		a.cfg.FixturesFilename,
		nil,
		eventAircraft,
	)
}
