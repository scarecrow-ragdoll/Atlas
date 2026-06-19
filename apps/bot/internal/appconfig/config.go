package appconfig

import (
	"time"

	"monorepo-template/libs/go/config"
)

// Config is the bot service configuration.
type Config struct {
	Bot BotConfig        `mapstructure:"bot"`
	Log config.LogConfig `mapstructure:"log"`
}

// BotConfig holds Telegram bot settings.
type BotConfig struct {
	Token       string        `mapstructure:"token" validate:"required"`
	PollTimeout time.Duration `mapstructure:"poll_timeout"`
}
