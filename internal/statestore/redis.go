package statestore

import (
	"context"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"go.uber.org/zap"
	"kurushimi/internal/config"
	"time"
)

type redisStore struct {
	healthCheckPool *redis.Pool
	redisPool       *redis.Pool
}

func NewRedis(cfg config.RedisConfig) StateStore {
	return &redisStore{
		healthCheckPool: getHealthCheckPool(cfg),
		redisPool:       getRedisPool(cfg),
	}
}

func (rs *redisStore) HealthCheck(ctx context.Context) error {
	redisConn, err := rs.healthCheckPool.GetContext(ctx)
	if err != nil {
		return err
	}
	defer handleConClose(redisConn)

	_, err = redisConn.Do("PING")
	if err != nil {
		return err
	}
	return nil
}

func getHealthCheckPool(cfg config.RedisConfig) *redis.Pool {
	return &redis.Pool{
		MaxIdle:      3,
		MaxActive:    0,
		IdleTimeout:  5 * time.Second,
		Wait:         true,
		TestOnBorrow: testOnBorrow,
		DialContext: func(ctx context.Context) (redis.Conn, error) {
			if ctx.Err() != nil {
				return nil, ctx.Err()
			}
			return redis.DialContext(ctx, "tcp", fmt.Sprintf("%s:%d", cfg.Host, cfg.Port))
		},
	}
}

func getRedisPool(cfg config.RedisConfig) *redis.Pool {
	return &redis.Pool{
		MaxIdle:      3,
		MaxActive:    0,
		IdleTimeout:  5 * time.Second,
		Wait:         true,
		TestOnBorrow: testOnBorrow,
		DialContext: func(ctx context.Context) (redis.Conn, error) {
			if ctx.Err() != nil {
				return nil, ctx.Err()
			}
			return redis.DialContext(ctx, "tcp", fmt.Sprintf("%s:%d", cfg.Host, cfg.Port))
		},
	}
}

func testOnBorrow(c redis.Conn, lastUsed time.Time) error {
	// Assume the connection is valid if it was used in 30 sec.
	if time.Since(lastUsed) < 15*time.Second {
		return nil
	}

	_, err := c.Do("PING")
	return err
}

func handleConClose(conn redis.Conn) {
	err := conn.Close()
	if err != nil {
		zap.S().Errorw("failed to close redis client connection.", err)
	}
}
