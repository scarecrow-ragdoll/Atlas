package handler

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"go.uber.org/zap"

	"monorepo-template/apps/bot/internal/botapi"
	"monorepo-template/libs/go/logger"
)

// Help returns a handler for the /help command.
func Help(s botapi.Sender) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		const op = "handler.Help"
		log := logger.FromContext(ctx).With(zap.String("op", op))
		log.Debug("handling /help")

		if _, err := s.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Available commands:\n/start — Start the bot\n/help — Show this message",
		}); err != nil {
			log.Error("failed to send message", zap.Error(err))
		}
	}
}
