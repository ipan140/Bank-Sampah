package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type CacheRepository interface {
	SetTokenBlacklist(ctx context.Context, token string, expiration time.Duration) error
	IsTokenBlacklisted(ctx context.Context, token string) (bool, error)
}

type cacheRepo struct {
	rdb *redis.Client
}

func NewCacheRepository(rdb *redis.Client) CacheRepository {
	return &cacheRepo{rdb: rdb}
}

func (c *cacheRepo) SetTokenBlacklist(ctx context.Context, token string, expiration time.Duration) error {
	if c.rdb == nil {
		return nil
	}
	key := fmt.Sprintf("blacklist:token:%s", token)
	return c.rdb.Set(ctx, key, "1", expiration).Err()
}

func (c *cacheRepo) IsTokenBlacklisted(ctx context.Context, token string) (bool, error) {
	if c.rdb == nil {
		return false, nil
	}
	key := fmt.Sprintf("blacklist:token:%s", token)
	exists, err := c.rdb.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return exists > 0, nil
}
