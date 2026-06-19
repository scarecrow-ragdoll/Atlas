package config

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

var (
	// ErrEnvFileLoad indicates .env file could not be loaded.
	ErrEnvFileLoad = errors.New("failed to load .env file")

	// ErrConfigLoad indicates config file could not be read or parsed.
	ErrConfigLoad = errors.New("failed to load config file")

	// ErrValidation indicates config validation failed.
	ErrValidation = errors.New("config validation failed")
)

// Options configures the loading behavior.
type Options struct {
	// ConfigPath is the path to the YAML config file (required).
	ConfigPath string

	// EnvFile is the path to a .env file (optional, "" = skip).
	EnvFile string

	// EnvPrefix is an optional prefix for env var lookup (e.g. "API" → API_POSTGRES_HOST).
	EnvPrefix string
}

// Load reads configuration from YAML + .env + env vars into T, then validates.
// Priority: real env vars > .env file > YAML defaults.
func Load[T any](opts Options) (T, error) {
	var zero T

	// Step 1: Load .env file (if specified).
	// godotenv.Load does NOT overwrite pre-existing env vars.
	if opts.EnvFile != "" {
		if err := loadDotEnv(opts.EnvFile); err != nil {
			return zero, err
		}
	}

	// Step 2: Read YAML config file.
	v := viper.New()
	v.SetConfigFile(opts.ConfigPath)

	if err := v.ReadInConfig(); err != nil {
		return zero, fmt.Errorf("%w: %v", ErrConfigLoad, err)
	}

	// Step 3: Enable automatic env var resolution.
	// Replaces "." with "_" so postgres.host → POSTGRES_HOST.
	if opts.EnvPrefix != "" {
		v.SetEnvPrefix(opts.EnvPrefix)
	}
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Step 4: Unmarshal into T.
	var cfg T
	if err := v.Unmarshal(&cfg); err != nil {
		return zero, fmt.Errorf("%w: %v", ErrConfigLoad, err)
	}

	// Step 5: Validate struct tags.
	if err := Validate(cfg); err != nil {
		return zero, err
	}

	return cfg, nil
}
