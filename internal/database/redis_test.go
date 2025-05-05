package database

import (
	"fmt"
	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestRedisDB_IsHealthy tests the IsHealthy method of the RedisDB struct.
func TestRedisDB_IsHealthy(t *testing.T) {
	t.Run("Redis is healthy", func(t *testing.T) {
		// Create a mocked Redis client
		mockClient, mock := redismock.NewClientMock()

		// Simulate a successful ping response
		mock.ExpectPing().SetVal("PONG")

		redisDB := RedisDB{Client: mockClient}

		// Test IsHealthy
		isHealthy := redisDB.IsHealthy()
		assert.True(t, isHealthy, "Redis connection should be healthy")

		// Ensure all expectations were met
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Redis is unhealthy", func(t *testing.T) {
		// Create a mocked Redis client
		mockClient, mock := redismock.NewClientMock()

		// Simulate a failed ping response
		mock.ExpectPing().SetErr(fmt.Errorf("ping error"))

		redisDB := RedisDB{Client: mockClient}

		// Test IsHealthy
		isHealthy := redisDB.IsHealthy()
		assert.False(t, isHealthy, "Redis connection should be unhealthy")

		// Ensure all expectations were met
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("No Redis connection", func(t *testing.T) {
		// Create a RedisDB instance with no client
		redisDB := RedisDB{Client: nil}

		// Test IsHealthy
		isHealthy := redisDB.IsHealthy()
		assert.False(t, isHealthy, "Redis connection should be unhealthy")
	})
}

// TestRedisDB_Close tests the Close method of the RedisDB struct.
func TestRedisDB_Close(t *testing.T) {
	t.Run("Close Redis connection successfully", func(t *testing.T) {
		// Create a mocked Redis client
		mockClient, _ := redismock.NewClientMock()

		redisDB := RedisDB{Client: mockClient}

		// Call Close
		redisDB.Close()

		// Call Close
		assert.NotPanics(t, func() { redisDB.Close() })
	})

	t.Run("Close Redis connection with nil client", func(t *testing.T) {
		// Create a RedisDB instance with no client
		redisDB := RedisDB{Client: nil}

		// Call Close
		assert.NotPanics(t, func() { redisDB.Close() })
	})
}
