package middleware

import (
	"context"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"go.uber.org/zap"

	"monorepo-template/libs/go/logger"
)

// Logging logs each update with metadata and injects an enriched logger into ctx.
// Guards against nil update.Message (non-message updates like callbacks).
func Logging(log *zap.Logger) bot.Middleware {
	return func(next bot.HandlerFunc) bot.HandlerFunc {
		return func(ctx context.Context, b *bot.Bot, update *models.Update) {
			const op = "middleware.Logging"
			start := time.Now()

			var chatID int64
			if update.Message != nil {
				chatID = update.Message.Chat.ID
			}

			l := log.With(
				zap.String("op", op),
				zap.Int64("update_id", update.ID),
				zap.Int64("chat_id", chatID),
			)
			ctx = logger.WithContext(ctx, l)

			next(ctx, b, update)

			l.Info("update processed",
				zap.Duration("duration", time.Since(start)),
			)
		}
	}
}
