package rtl28xxx

import (
	"context"
	"io"
	"os"
	"syscall"

	"github.com/landru29/dump1090/internal/source"
)

// SourceFile is the data file source.
type SourceFile struct {
	filename string
	loop     bool

	processor source.Processer
}

// Configurator is the Source configurator.
type FileConfigurator func(*SourceFile)

// New creates a new data source process.
func NewFile(filename string, processor source.Processer, opts ...FileConfigurator) *SourceFile {
	output := &SourceFile{
		filename:  filename,
		processor: processor,
	}

	for _, opt := range opts {
		opt(output)
	}

	return output
}

// WithLoop is the data loop configurator.
func WithLoop() FileConfigurator {
	return func(s *SourceFile) {
		s.loop = true
	}
}

// Start implements the source.Starter interface.
func (s *SourceFile) Start(ctx context.Context) error {
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

	err := s.start(ctx)

	syscall.Kill(syscall.Getpid(), syscall.SIGTERM)

	return err
}

func (s *SourceFile) start(ctx context.Context) error {
	fd, err := os.Open(s.filename)
	if err != nil {
		return err
	}

	defer func(closer io.Closer) {
		_ = closer.Close()
	}(fd)

	return NewReader(fd, s.processor).Start(ctx)
}
