package handler_test

import (
	"context"
	"errors"
	"testing"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"monorepo-template/apps/bot/internal/handler"
)

type mockSender struct {
	lastParams *bot.SendMessageParams
	err        error
}

func (m *mockSender) SendMessage(ctx context.Context, params *bot.SendMessageParams) (*models.Message, error) {
	m.lastParams = params
	return &models.Message{}, m.err
}

func TestStart_SendsWelcome(t *testing.T) {
	s := &mockSender{}
	h := handler.Start(s)

	update := &models.Update{
		Message: &models.Message{
			Chat: models.Chat{ID: 123},
		},
	}

	h(context.Background(), nil, update)

	require.NotNil(t, s.lastParams)
	assert.Equal(t, int64(123), s.lastParams.ChatID)
	assert.Contains(t, s.lastParams.Text, "Welcome")
}

func TestStart_LogsSendError(t *testing.T) {
	s := &mockSender{err: errors.New("telegram failed")}
	h := handler.Start(s)

	assert.NotPanics(t, func() {
		h(context.Background(), nil, &models.Update{
			Message: &models.Message{Chat: models.Chat{ID: 123}},
		})
	})
	require.NotNil(t, s.lastParams)
}

func TestHelp_SendsHelpText(t *testing.T) {
	s := &mockSender{}
	h := handler.Help(s)

	update := &models.Update{
		Message: &models.Message{
			Chat: models.Chat{ID: 456},
		},
	}

	h(context.Background(), nil, update)

	require.NotNil(t, s.lastParams)
	assert.Equal(t, int64(456), s.lastParams.ChatID)
	assert.Contains(t, s.lastParams.Text, "/start")
	assert.Contains(t, s.lastParams.Text, "/help")
}

func TestHelp_LogsSendError(t *testing.T) {
	s := &mockSender{err: errors.New("telegram failed")}
	h := handler.Help(s)

	assert.NotPanics(t, func() {
		h(context.Background(), nil, &models.Update{
			Message: &models.Message{Chat: models.Chat{ID: 456}},
		})
	})
	require.NotNil(t, s.lastParams)
}

func TestDefault_DoesNotPanic(t *testing.T) {
	h := handler.Default()

	assert.NotPanics(t, func() {
		h(context.Background(), nil, &models.Update{ID: 1})
	})
}
