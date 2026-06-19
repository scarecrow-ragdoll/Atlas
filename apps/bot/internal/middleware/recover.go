package middleware

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"go.uber.org/zap"
)

// Recover catches panics in downstream handlers and logs them.
func Recover(log *zap.Logger) bot.Middleware {
	return func(next bot.HandlerFunc) bot.HandlerFunc {
		return func(ctx context.Context, b *bot.Bot, update *models.Update) {
			const op = "middleware.Recover"
			defer func() {
				if r := recover(); r != nil {
					log.Error("panic recovered",
						zap.String("op", op),
						zap.Any("panic", r),
						zap.Int64("update_id", update.ID),
					)
				}
			}()
			next(ctx, b, update)
		}
	}
}
