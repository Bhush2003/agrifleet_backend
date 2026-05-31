package database

import (
	"context"
	"fmt"

	"github.com/agrifleet/backend/internal/config"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

// NewRedis creates and returns a Redis client.
func NewRedis(cfg *config.Config) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	ctx := context.Background()
	if _, err := client.Ping(ctx).Result(); err != nil {
		return nil, fmt.Errorf("database.NewRedis: %w", err)
	}

	log.Info().Msg("Redis connected successfully")
	return client, nil
}
