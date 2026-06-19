package logger

import (
	"context"

	"go.uber.org/zap"
)

type ctxKey struct{}

// WithContext returns a new context with the given logger stored as a value.
func WithContext(ctx context.Context, l *zap.Logger) context.Context {
	return context.WithValue(ctx, ctxKey{}, l)
}

// FromContext extracts the logger from context.
// Returns zap.NewNop() if no logger is present.
func FromContext(ctx context.Context) *zap.Logger {
	if l, ok := ctx.Value(ctxKey{}).(*zap.Logger); ok {
		return l
	}
	return zap.NewNop()
}
