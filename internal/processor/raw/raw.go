// Package raw is the raw decoder.
package raw

import (
	"encoding/hex"
	"log/slog"
	"strings"
)

// New creates a raw processor.
func New(log *slog.Logger) *Processor {
	return &Processor{log: log.With("processor", "raw")}
}

// Processor is an empty processor.
type Processor struct {
	log *slog.Logger
}

// Process implements the Processer interface.
func (e Processor) Process(data []byte) error {
	e.log.Info("New message", "data", strings.ToUpper(hex.EncodeToString(data)))
	return nil
}
