package database

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"time"
)

type PostgresDB struct {
	DB *sqlx.DB
}

// NewPostgresDB create a new Postgres DB.
func NewPostgresDB(db SqlDatabase) PostgresDB {
	zap.L().Info("Connecting to Postgres...")

	// Connect to Postgres
	dbClient, err := db.Connect()
	if err != nil {
		zap.L().Error("main.DbConnection:", zap.Error(err))
		return PostgresDB{
			DB: nil,
		}
	}

	zap.L().Info("Connected to Postgres")

	// Finish up configuration
	dbClient.SetMaxOpenConns(viper.GetInt("POSTGRES_MAX_OPEN_CONNS"))
	dbClient.SetMaxIdleConns(viper.GetInt("POSTGRES_MAX_IDLE_CONNS"))

	return PostgresDB{
		DB: dbClient,
	}
}

// IsHealthy checks if the database connection is healthy by running a simple query.
func (p PostgresDB) IsHealthy() bool {
	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// No database connection
	if p.DB == nil {
		return false
	}

	// Attempt to ping the database
	err := p.DB.PingContext(ctx)
	if err != nil {
		return false
	}

	return true
}
