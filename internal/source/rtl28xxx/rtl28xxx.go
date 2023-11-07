// Package rtlsdr is the RTL28xxx data source.
package rtl28xxx

import (
	"context"

	"github.com/landru29/dump1090/internal/errors"
	"github.com/landru29/dump1090/internal/source"
)

const (
	ErrNoDeviceFound errors.Error = "no device found"

	modeSfrequency = 1090000000
	sampleRate     = 2000000

	asyncBufNumber = 12
	dataLen        = (16 * 16384) /* 256k */
)

// Configurator is the Source configurator.
type Configurator func(*Source)

// Source is the data source process.
type Source struct {
	deviceIndex uint32
	frequency   uint32
	gain        float64
	enableAGC   bool

	processer source.Processer
	dev       *Device
}

// New creates a new data source process.
func New(processer source.Processer, opts ...Configurator) *Source {
	output := &Source{
		deviceIndex: 0,
		frequency:   modeSfrequency,
		gain:        0,
		enableAGC:   false,
	}

	for _, opt := range opts {
		opt(output)
	}

	return output
}

// Start implements the source.Starter interface.
func (s *Source) Start(ctx context.Context) error {
	deviceCount := DeviceCount()
	if deviceCount == 0 {
		return ErrNoDeviceFound
	}

	deviceIndex := uint32(0)

	if s.deviceIndex < deviceCount {
		deviceIndex = s.deviceIndex
	}

	device, err := OpenDevice(uint32(deviceIndex), s.processer)
	if err != nil {
		return err
	}

	s.dev = device

	if err := s.dev.SetCenterFreq(modeSfrequency); err != nil {
		return err
	}

	if err := s.dev.SetSampleRate(sampleRate); err != nil {
		return err
	}

	if err := s.dev.SetAgcMode(s.enableAGC); err != nil {
		return err
	}

	if err := s.dev.SetTunerGainMode(s.gain > 0); err != nil {
		return err
	}

	if s.gain > 0 {
		if err := s.dev.SetTunerGain(s.gain); err != nil {
			return err
		}
	}

	return s.dev.ReadAsync(ctx, asyncBufNumber, dataLen)
}

// WithDeviceIndex configures the device index.
func WithDeviceIndex(index int) Configurator {
	return func(s *Source) {
		if index > 0 {
			s.deviceIndex = uint32(index)
		}
	}
}

// WithFrequencyv configures the frequency.
func WithFrequency(frequency uint32) Configurator {
	return func(s *Source) {
		s.frequency = frequency
	}
}

// WithGain configures the gain.
func WithGain(gain float64) Configurator {
	return func(s *Source) {
		s.gain = gain
	}
}

// WithAGC enables AGC.
func WithAGC() Configurator {
	return func(s *Source) {
		s.enableAGC = true
	}
}
