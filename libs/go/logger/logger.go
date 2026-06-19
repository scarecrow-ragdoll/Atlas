package logger

import (
	"fmt"

	"go.uber.org/zap"
)

// Config holds logger configuration.
type Config struct {
	Level  string // "debug", "info", "warn", "error"; default "info"
	Format string // "json" or "console"; default "json"
}

var buildZapLogger = func(cfg zap.Config) (*zap.Logger, error) {
	return cfg.Build()
}

// New creates a configured *zap.Logger.
// "json" format uses zap.NewProductionConfig(), "console" uses zap.NewDevelopmentConfig().
// Returns error if level string is invalid.
func New(cfg Config) (*zap.Logger, error) {
	if cfg.Level == "" {
		cfg.Level = "info"
	}

	level, err := zap.ParseAtomicLevel(cfg.Level)
	if err != nil {
		return nil, fmt.Errorf("parse log level %q: %w", cfg.Level, err)
	}

	var zapCfg zap.Config
	if cfg.Format == "console" {
		zapCfg = zap.NewDevelopmentConfig()
	} else {
		zapCfg = zap.NewProductionConfig()
	}
	zapCfg.Level = level

	l, err := buildZapLogger(zapCfg)
	if err != nil {
		return nil, fmt.Errorf("build logger: %w", err)
	}

	return l, nil
}
