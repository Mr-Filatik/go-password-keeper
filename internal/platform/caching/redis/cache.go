// Package redis provides a generic implementation for caching via redis.
package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/mr-filatik/go-password-keeper/internal/platform/caching/redis/adapter"
	"github.com/mr-filatik/go-password-keeper/internal/platform/logging"
	"github.com/redis/go-redis/v9"
)

// Cacher describes the Cacher structure for communicating with the simple redis service.
type Cacher struct {
	logger logging.Logger
	client *redis.Client
	config CacherConfig
}

// CacherConfig describes the configuration for Cacher.
type CacherConfig struct {
	ClientName string
	Address    string // Address
	DBNumber   int    // Database number
	Username   string
	Password   string
}

// NewCacher creates a new *Cacher instance.
//
// Parameters:
//   - conf CacherConfig: config;
//   - logger logging.Logger: logger.
func NewCacher(conf CacherConfig, logger logging.Logger) *Cacher {
	logger.Info("Cacher creating...")

	redis.SetLogger(adapter.NewLoggerAdapter(logger))

	chr := &Cacher{
		logger: logger,
		client: nil,
		config: conf,
	}

	logger.Info("Cacher create is successful")

	return chr
}

// Start - starting the cacher.
//
// Implements the server.IServer interface.
func (c *Cacher) Start(_ context.Context) error {
	c.logger.Info(
		"Cacher starting...",
		"address", c.config.Address,
		"database", c.config.DBNumber,
	)

	//nolint:exhaustruct // other options use the default value
	redisOptions := &redis.Options{
		Addr:       c.config.Address,
		ClientName: c.config.ClientName,
		DB:         c.config.DBNumber, // 0 - default (for Sentinel database and cluster the count is 1)
		Username:   c.config.Username,
		Password:   c.config.Password,
	}

	c.client = redis.NewClient(redisOptions)

	c.logger.Info("Cacher start is successful")

	return nil
}

// SetValue stores the value as a string by key.
func (c *Cacher) SetValue(
	ctx context.Context,
	key string,
	value string,
	expiration time.Duration,
) error {
	err := c.client.Set(ctx, key, value, expiration).Err()
	if err != nil {
		return fmt.Errorf("set value: %w", err)
	}

	return nil
}

// GetValue gets the value as a string by key.
func (c *Cacher) GetValue(ctx context.Context, key string) (string, error) {
	value, err := c.client.Get(ctx, key).Result()
	if err != nil {
		return "", fmt.Errorf("get value: %w", err)
	}

	return value, nil
}
