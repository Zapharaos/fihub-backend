package app

import (
	"github.com/Zapharaos/fihub-backend/internal/database"
	"github.com/spf13/viper"
)

// InitPostgres init the postgres database.
func InitPostgres() bool {
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

// InitRedis init the redis database.
func InitRedis() bool {
	// Connect to Redis
	redis := database.NewRedisDB()

	// Verify connection
	if redis.IsHealthy() {
		database.DB().SetRedis(redis)
		return true
	}
	return false
}
