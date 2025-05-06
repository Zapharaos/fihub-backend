package database

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"time"
)

type RedisDB struct {
	Client *redis.Client
}

// NewRedisDB creates a new Redis client.
func NewRedisDB() RedisDB {
	zap.L().Info("Connecting to Redis...")

	// Connect to Redis
	host := viper.GetString("REDIS_HOST")
	port := viper.GetString("REDIS_PORT")
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: viper.GetString("REDIS_PASSWORD"),
		DB:       viper.GetInt("REDIS_DB"),
		PoolSize: viper.GetInt("REDIS_POOL_SIZE"),
	})

	// Verify connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := client.Ping(ctx).Result()
	if err != nil {
		zap.L().Error("Failed to connect to Redis", zap.Error(err))
		return RedisDB{
			Client: nil,
		}
	}

	zap.L().Info("Connected to Redis")

	return RedisDB{
		Client: client,
	}
}

// IsHealthy checks if the Redis connection is healthy by running a ping command.
func (r RedisDB) IsHealthy() bool {
	// No Redis connection
	if r.Client == nil {
		return false
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Attempt to ping Redis
	_, err := r.Client.Ping(ctx).Result()
	return err == nil
}

// Close closes the Redis connection.
func (r RedisDB) Close() {
	if r.Client != nil {
		err := r.Client.Close()
		if err != nil {
			zap.L().Error("Failed to close Redis connection", zap.Error(err))
		}
	}
}
