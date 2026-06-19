package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-telegram/bot"
	"go.uber.org/zap"

	"monorepo-template/apps/bot/internal/appconfig"
	"monorepo-template/apps/bot/internal/handler"
	"monorepo-template/apps/bot/internal/middleware"
	"monorepo-template/libs/go/config"
	"monorepo-template/libs/go/logger"
)

func main() {
	// 1. Config
	cfg, err := config.Load[appconfig.Config](config.Options{
		ConfigPath: "config/config.yml",
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config: %v\n", err)
		os.Exit(1)
	}

	// 2. Logger
	log, err := logger.New(logger.Config{
		Level:  cfg.Log.Level,
		Format: cfg.Log.Format,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to init logger: %v\n", err)
		os.Exit(1)
	}
	defer func() { _ = log.Sync() }()

	// 3. Bot
	// Middleware order matters: Recover wraps Logging, so panics inside
	// Logging are caught. Register Recover first.
	opts := []bot.Option{
		bot.WithDefaultHandler(handler.Default()),
		bot.WithMiddlewares(
			middleware.Recover(log),
			middleware.Logging(log),
		),
	}
	b, err := bot.New(cfg.Bot.Token, opts...)
	if err != nil {
		log.Fatal("failed to create bot", zap.Error(err))
	}

	// 4. Handlers
	b.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, handler.Start(b))
	b.RegisterHandler(bot.HandlerTypeMessageText, "/help", bot.MatchTypeExact, handler.Help(b))

	// 5. Graceful shutdown
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	log.Info("bot started")
	b.Start(ctx) // blocking — returns when ctx is cancelled

	// Give in-flight updates time to finish, then force exit.
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	<-shutdownCtx.Done()
	log.Info("bot stopped")
}
