package logging

import (
	"context"
	"log/slog"
)

type contextKey struct{}

var loggerKey = contextKey{}

// WithLogger returns a new context with the provided logger.
func WithLogger(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

// Ctx returns the logger from the context.
// If no logger is found, it returns the default slog logger.
func Ctx(ctx context.Context) *slog.Logger {
	if logger, ok := ctx.Value(loggerKey).(*slog.Logger); ok {
		return logger
	}
	return slog.Default()
}

// Err returns a slog attribute for an error with a consistent key.
func Err(err error) slog.Attr {
	return slog.Any("error", err)
}
