package middleware_test

import (
	"context"
	"testing"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"

	"monorepo-template/apps/bot/internal/middleware"
)

func TestRecover_CatchesPanic(t *testing.T) {
	core, logs := observer.New(zap.DebugLevel)
	log := zap.New(core)

	panicking := func(ctx context.Context, b *bot.Bot, update *models.Update) {
		panic("test panic")
	}

	wrapped := middleware.Recover(log)(panicking)

	// passing nil for *bot.Bot is safe here because the handler panics before using it
	assert.NotPanics(t, func() {
		wrapped(context.Background(), nil, &models.Update{ID: 1})
	})

	require.Equal(t, 1, logs.Len())
	assert.Equal(t, "panic recovered", logs.All()[0].Message)

	fields := logs.All()[0].ContextMap()
	assert.Equal(t, "middleware.Recover", fields["op"])
	assert.Equal(t, "test panic", fields["panic"])
	assert.Equal(t, int64(1), fields["update_id"])
}

func TestRecover_NoRecoverWhenNoPanic(t *testing.T) {
	core, logs := observer.New(zap.DebugLevel)
	log := zap.New(core)

	called := false
	normal := func(ctx context.Context, b *bot.Bot, update *models.Update) {
		called = true
	}

	wrapped := middleware.Recover(log)(normal)
	wrapped(context.Background(), nil, &models.Update{ID: 2})

	assert.True(t, called)
	assert.Equal(t, 0, logs.Len())
}
