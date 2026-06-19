package handler

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

// Default returns a catch-all handler that does nothing.
func Default() bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {}
}
