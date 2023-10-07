package application

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/landru29/dump1090/internal/dump"
)

const (
	cleanDelay time.Duration = time.Second * 10

	outOfDateAC time.Duration = time.Minute
)

type Config struct {
	FixturesFilename string
	DeviceIndex      uint32
	Frequency        uint32
	Gain             int
	EnableAGC        bool
}

type App struct {
	cfg          *Config
	aircraftPool map[uint32]*dump.Aircraft
	mutex        sync.Mutex
}

func New(cfg *Config) (*App, error) {
	return &App{
		cfg:          cfg,
		aircraftPool: map[uint32]*dump.Aircraft{},
	}, nil
}

func (a *App) Start(ctx context.Context) error {
	eventAircraft := make(chan *dump.Aircraft, 10)
	defer func() {
		close(eventAircraft)
	}()

	go func(app *App) {
		app.acCleaner(ctx)
	}(a)

	go func(acStream chan *dump.Aircraft, app *App) {
		for {
			ac := <-acStream

			app.mutex.Lock()
			app.aircraftPool[ac.Addr] = ac
			app.mutex.Unlock()

			fmt.Println(ac)
		}
	}(eventAircraft, a)

	return dump.Start(
		a.cfg.DeviceIndex,
		a.cfg.Gain,
		a.cfg.Frequency,
		a.cfg.EnableAGC,
		a.cfg.FixturesFilename,
		nil,
		eventAircraft,
	)
}

func (a *App) acCleaner(ctx context.Context) {
	ticker := time.NewTicker(cleanDelay)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			a.mutex.Lock()
			for idx, ac := range a.aircraftPool {
				if time.Since(ac.Seen) > outOfDateAC {
					delete(a.aircraftPool, idx)
				}
			}
			a.mutex.Unlock()
		}
	}
}
