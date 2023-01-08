package statestore

import (
	"context"
	"fmt"
	rs "github.com/go-redsync/redsync/v4"
	"github.com/gomodule/redigo/redis"
	"go.uber.org/zap"
	"os"
	"time"
)

var (
	logger, _ = zap.NewProduction()

	redsync *rs.Redsync

	healthCheckPool = getHealthCheckPool()
	redisPool       = getRedisPool()
	mutex           *rs.Mutex
)

func HealthCheck() error {
	redisConn, err := healthCheckPool.GetContext(context.Background())
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

func getHealthCheckPool() *redis.Pool {
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
			return redis.DialContext(ctx, "tcp", getRedisAddress())
		},
	}
}

func getRedisPool() *redis.Pool {
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
			return redis.DialContext(ctx, "tcp", getRedisAddress())
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
		logger.Error("failed to close redis client connection.", zap.Error(err))
	}
}

func getRedisAddress() string {
	host := "localhost"
	if eHost := os.Getenv("REDIS_HOST"); eHost != "" {
		host = eHost
	}

	return fmt.Sprintf("%s:6379", host)
}
