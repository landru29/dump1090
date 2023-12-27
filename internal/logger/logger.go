// Package logger manages the logger.
package logger

import (
	"context"
	"log/slog"

	"github.com/landru29/dump1090/internal/errors"
)

const (
	// ErrMissingLogger is when the logger is missing.
	ErrMissingLogger errors.Error = "missing logger"
)

type logContext struct{}

// WithLogger sets a logger in the context.
func WithLogger(ctx context.Context, log *slog.Logger) context.Context {
	return context.WithValue(ctx, logContext{}, log)
}

// Logger gets the logger from the context.
func Logger(ctx context.Context) (*slog.Logger, bool) {
	out, found := ctx.Value(logContext{}).(*slog.Logger)

	return out, found
}
