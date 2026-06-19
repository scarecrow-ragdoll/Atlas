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
	"monorepo-template/libs/go/logger"
)

func TestLogging_LogsUpdateFields(t *testing.T) {
	core, logs := observer.New(zap.DebugLevel)
	log := zap.New(core)

	inner := func(ctx context.Context, b *bot.Bot, update *models.Update) {
		// verify logger is in context
		l := logger.FromContext(ctx)
		assert.NotNil(t, l)
	}

	wrapped := middleware.Logging(log)(inner)
	wrapped(context.Background(), nil, &models.Update{
		ID: 42,
		Message: &models.Message{
			Chat: models.Chat{ID: 123},
		},
	})

	require.Equal(t, 1, logs.Len())
	entry := logs.All()[0]
	assert.Equal(t, "update processed", entry.Message)

	fields := entry.ContextMap()
	assert.Equal(t, "middleware.Logging", fields["op"])
	assert.Equal(t, int64(42), fields["update_id"])
	assert.Equal(t, int64(123), fields["chat_id"])
	assert.Contains(t, fields, "duration")
}

func TestLogging_NilMessage_NoPanic(t *testing.T) {
	core, logs := observer.New(zap.DebugLevel)
	log := zap.New(core)

	inner := func(ctx context.Context, b *bot.Bot, update *models.Update) {}

	wrapped := middleware.Logging(log)(inner)

	// non-message update (e.g., callback query) — Message is nil
	assert.NotPanics(t, func() {
		wrapped(context.Background(), nil, &models.Update{ID: 99})
	})

	require.Equal(t, 1, logs.Len())
	fields := logs.All()[0].ContextMap()
	assert.Equal(t, int64(0), fields["chat_id"])
}
