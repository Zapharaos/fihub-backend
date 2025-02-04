package database

import (
	"github.com/Zapharaos/fihub-backend/pkg/env"
	"github.com/jmoiron/sqlx"
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
	dbClient.SetMaxOpenConns(env.GetInt("POSTGRES_MAX_OPEN_CONNS", 30))
	dbClient.SetMaxIdleConns(env.GetInt("POSTGRES_MAX_IDLE_CONNS", 30))

	return dbClient
}
