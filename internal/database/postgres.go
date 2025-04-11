package database

import (
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// NewPostgresDB create a new Postgres DB.
func NewPostgresDB(db SqlDatabase) *sqlx.DB {
	zap.L().Info("Initializing Postgres")

	// Connect to Postgres
	dbClient, err := db.Connect()
	if err != nil {
		zap.L().Error("main.DbConnection:", zap.Error(err))
		return nil
	}

	zap.L().Info("Connected to Postgres")

	// Finish up configuration
	dbClient.SetMaxOpenConns(viper.GetInt("POSTGRES_MAX_OPEN_CONNS"))
	dbClient.SetMaxIdleConns(viper.GetInt("POSTGRES_MAX_IDLE_CONNS"))

	return dbClient
}
