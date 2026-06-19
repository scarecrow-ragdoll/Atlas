package config

import (
	"fmt"
	"time"
)

// PostgresConfig holds PostgreSQL connection settings.
type PostgresConfig struct {
	Host            string        `mapstructure:"host"              validate:"required"`
	Port            int           `mapstructure:"port"              validate:"required,gt=0"`
	User            string        `mapstructure:"user"              validate:"required"`
	Password        string        `mapstructure:"password"`
	DB              string        `mapstructure:"db"                validate:"required"`
	SSLMode         string        `mapstructure:"sslmode"`
	MaxConns        int32         `mapstructure:"max_conns"`
	MinConns        int32         `mapstructure:"min_conns"`
	MaxConnIdleTime time.Duration `mapstructure:"max_conn_idle_time"`
}

// DSN returns a PostgreSQL connection string.
func (c PostgresConfig) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.User, c.Password, c.Host, c.Port, c.DB, c.SSLMode,
	)
}

// RedisConfig holds Redis connection settings.
type RedisConfig struct {
	Host     string `mapstructure:"host"     validate:"required"`
	Port     int    `mapstructure:"port"     validate:"required,gt=0"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

// LogConfig holds logging settings.
type LogConfig struct {
	Level  string `mapstructure:"level"  validate:"omitempty,oneof=debug info warn error"`
	Format string `mapstructure:"format" validate:"omitempty,oneof=json text"`
}

// ServerConfig holds HTTP server settings.
type ServerConfig struct {
	Port            int           `mapstructure:"port"             validate:"required,gt=0"`
	ReadTimeout     time.Duration `mapstructure:"read_timeout"`
	WriteTimeout    time.Duration `mapstructure:"write_timeout"`
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"`
	Env             string        `mapstructure:"env"              validate:"omitempty,oneof=development staging production"`
	CORSOrigins     []string      `mapstructure:"cors_origins"`
}
