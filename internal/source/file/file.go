// Package file is file data source.
package file

import (
	"context"
	"io"
	"os"

	"github.com/landru29/dump1090/internal/source"
)

// Source is the data source.
type Source struct {
	filename string
	loop     bool

	processer source.Processer
}

// Configurator is the Source configurator.
type Configurator func(*Source)

// New creates a new data source process.
func New(filename string, processer source.Processer, opts ...Configurator) *Source {
	output := &Source{
		filename: filename,
	}

	return output
}

// WithLoop is the data loop configurator.
func WithLoop() Configurator {
	return func(s *Source) {
		s.loop = true
	}
}

// Start implements the source.Starter interface.
func (s *Source) Start(ctx context.Context) error {
	if s.loop {
		for {
			select {
			case <-ctx.Done():
				return nil
			default:
				if err := s.start(ctx); err != nil {
					return err
				}
			}
		}
	}

	return s.start(ctx)
}

func (s *Source) start(ctx context.Context) error {
	fd, err := os.Open(s.filename)
	if err != nil {
		return err
	}

	defer func(closer io.Closer) {
		_ = closer.Close()
	}(fd)

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			data := make([]byte, 1024)
			size, err := fd.Read(data)

			if size > 0 {
				s.processer.Process(data[:size])
			}

			if err == io.EOF {
				return nil
			}

			if err != nil {
				return err
			}
		}
	}
}
