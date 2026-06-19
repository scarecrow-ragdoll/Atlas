package handler

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"go.uber.org/zap"

	"monorepo-template/apps/bot/internal/botapi"
	"monorepo-template/libs/go/logger"
)

// Start returns a handler for the /start command.
func Start(s botapi.Sender) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		const op = "handler.Start"
		log := logger.FromContext(ctx).With(zap.String("op", op))
		log.Debug("handling /start")

		if _, err := s.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Welcome! Use /help to see available commands.",
		}); err != nil {
			log.Error("failed to send message", zap.Error(err))
		}
	}
}
