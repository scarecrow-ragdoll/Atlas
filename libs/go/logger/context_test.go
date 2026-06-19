package logger_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"

	"monorepo-template/libs/go/logger"
)

func TestWithContext_FromContext_RoundTrip(t *testing.T) {
	core, _ := observer.New(zap.DebugLevel)
	l := zap.New(core)

	ctx := logger.WithContext(context.Background(), l)
	got := logger.FromContext(ctx)

	assert.Equal(t, l, got)
}

func TestFromContext_EmptyCtx_ReturnsNop(t *testing.T) {
	got := logger.FromContext(context.Background())
	assert.NotNil(t, got)
	// nop logger should not panic on usage
	got.Info("should not panic")
}
