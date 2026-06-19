package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	"monorepo-template/libs/go/config"
)

type Client struct {
	RDB    *redis.Client
	logger *zap.Logger
}

func New(cfg config.RedisConfig, logger *zap.Logger) (*Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	logger.Info("connected to Redis",
		zap.String("host", cfg.Host),
		zap.Int("port", cfg.Port),
	)

	return &Client{RDB: rdb, logger: logger}, nil
}

func (c *Client) Ping() error {
	return c.RDB.Ping(context.Background()).Err()
}

func (c *Client) Close() error {
	return c.RDB.Close()
}
