package app

import (
	"github.com/Zapharaos/fihub-backend/internal/database"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"time"
)

// ConnectPostgres connects to the postgres database.
func ConnectPostgres() bool {
	// Connect to Postgres
	postgres := database.NewPostgresDB(database.NewSqlDatabase(database.SqlCredentials{
		Host:     viper.GetString("POSTGRES_HOST"),
		Port:     viper.GetString("POSTGRES_PORT"),
		User:     viper.GetString("POSTGRES_USER"),
		Password: viper.GetString("POSTGRES_PASSWORD"),
		DbName:   viper.GetString("POSTGRES_DB"),
	}))

	// Verify connection
	if postgres.IsHealthy() {
		database.DB().SetPostgres(postgres)
		return true
	}
	return false
}

// ConnectRedis connects to the redis database.
func ConnectRedis() bool {
	// Connect to Redis
	redis := database.NewRedisDB()

	// Verify connection
	if redis.IsHealthy() {
		database.DB().SetRedis(redis)
		return true
	}
	return false
}

// StartDatabaseHealthChecks starts a goroutine that periodically checks
// the health of the databases and attempts to reconnect if necessary.
// interval: the time between health checks
// onSuccessfulConnection: optional callback function executed after successful reconnection
// Returns a function that can be called to stop the health check.
func StartPostgresHealthCheck(interval time.Duration, onSuccessfulConnection func()) func() {
	ticker := time.NewTicker(interval)
	done := make(chan bool)

	go func() {
		for {
			select {
			case <-done:
				ticker.Stop()
				return
			case <-ticker.C:
				if !database.DB().Postgres().IsHealthy() {
					if ConnectPostgres() {
						zap.L().Info("Reconnected to Postgres")
						if onSuccessfulConnection != nil {
							onSuccessfulConnection()
						}
					}
				}
			}
		}
	}()

	// Return function to stop the health check
	return func() {
		done <- true
	}
}
