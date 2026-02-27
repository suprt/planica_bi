package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/suprt/planica_bi/backend/internal/config"
	"github.com/suprt/planica_bi/backend/internal/logger"
	"go.uber.org/zap"
)

// Cache wraps Redis client for caching operations
type Cache struct {
	client *redis.Client
	ctx    context.Context
}

// NewCache creates a new cache instance
func NewCache(cfg *config.Config) (*Cache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisHost + ":" + cfg.RedisPort,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	// Test connection
	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &Cache{
		client: client,
		ctx:    ctx,
	}, nil
}

// Close closes the Redis connection
func (c *Cache) Close() error {
	return c.client.Close()
}

// Get retrieves a value from cache
func (c *Cache) Get(key string, dest interface{}) error {
	val, err := c.client.Get(c.ctx, key).Result()
	if err == redis.Nil {
		return ErrCacheMiss
	}
	if err != nil {
		return fmt.Errorf("cache get error: %w", err)
	}

	if err := json.Unmarshal([]byte(val), dest); err != nil {
		return fmt.Errorf("cache unmarshal error: %w", err)
	}

	return nil
}

// Set stores a value in cache with TTL
func (c *Cache) Set(key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("cache marshal error: %w", err)
	}

	if err := c.client.Set(c.ctx, key, data, ttl).Err(); err != nil {
		return fmt.Errorf("cache set error: %w", err)
	}

	return nil
}

// Delete removes a key from cache
func (c *Cache) Delete(key string) error {
	if err := c.client.Del(c.ctx, key).Err(); err != nil {
		return fmt.Errorf("cache delete error: %w", err)
	}
	return nil
}

// InvalidatePattern removes all keys matching the pattern
func (c *Cache) InvalidatePattern(pattern string) error {
	iter := c.client.Scan(c.ctx, 0, pattern, 0).Iterator()
	var keys []string

	for iter.Next(c.ctx) {
		keys = append(keys, iter.Val())
	}

	if err := iter.Err(); err != nil {
		return fmt.Errorf("cache scan error: %w", err)
	}

	if len(keys) > 0 {
		if err := c.client.Del(c.ctx, keys...).Err(); err != nil {
			return fmt.Errorf("cache delete pattern error: %w", err)
		}

		if logger.Log != nil {
			logger.Log.Info("Cache invalidated",
				zap.String("pattern", pattern),
				zap.Int("keys_count", len(keys)),
			)
		}
	}

	return nil
}

// ErrCacheMiss is returned when a key is not found in cache
var ErrCacheMiss = &CacheMissError{}

// CacheMissError represents a cache miss
type CacheMissError struct{}

func (e *CacheMissError) Error() string {
	return "cache miss"
}

// Cache key prefixes
// All these are reference data that change only on manual admin actions, not during sync
const (
	KeyPrefixCounters        = "counters:project:"
	KeyPrefixGoals           = "goals:counter:"
	KeyPrefixDirectAccounts  = "direct:accounts:project:"
	KeyPrefixDirectCampaigns = "direct:campaigns:account:"
)

// BuildKey constructs a cache key from prefix and ID
func BuildKey(prefix string, id uint) string {
	return fmt.Sprintf("%s%d", prefix, id)
}
