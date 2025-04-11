package app

import (
	"github.com/Zapharaos/fihub-backend/internal/auth/password"
	"github.com/Zapharaos/fihub-backend/internal/auth/permissions"
	"github.com/Zapharaos/fihub-backend/internal/auth/roles"
	"github.com/Zapharaos/fihub-backend/internal/auth/users"
	"github.com/Zapharaos/fihub-backend/internal/brokers"
	"github.com/Zapharaos/fihub-backend/internal/database"
	"github.com/Zapharaos/fihub-backend/internal/transactions"
	"github.com/Zapharaos/fihub-backend/pkg/env"
	"github.com/jmoiron/sqlx"
)

// initDatabase initializes the database connections.
func initDatabase() {
	postgres := database.NewPostgresDB(database.NewSqlDatabase(database.SqlCredentials{
		Host:     env.GetString("POSTGRES_HOST", "localhost"),
		Port:     env.GetString("POSTGRES_PORT", "5432"),
		User:     env.GetString("POSTGRES_USER", "postgres"),
		Password: env.GetString("POSTGRES_PASSWORD", "password"),
		DbName:   env.GetString("POSTGRES_DB", "postgres"),
	}))
	database.ReplaceGlobals(database.NewDatabases(postgres))

	// Initialize the postgres repositories
	initPostgres(database.DB().Postgres())
}

// initPostgres initializes the postgres repositories.
func initPostgres(dbClient *sqlx.DB) {
	// Auth
	users.ReplaceGlobals(users.NewPostgresRepository(dbClient))
	password.ReplaceGlobals(password.NewPostgresRepository(dbClient))

	// Roles
	roles.ReplaceGlobals(roles.NewPostgresRepository(dbClient))

	// Permissions
	permissions.ReplaceGlobals(permissions.NewPostgresRepository(dbClient))

	// Brokers
	brokerRepository := brokers.NewPostgresRepository(dbClient)
	userBrokerRepository := brokers.NewUserPostgresRepository(dbClient)
	imageBrokerRepository := brokers.NewImagePostgresRepository(dbClient)
	brokers.ReplaceGlobals(brokers.NewRepository(brokerRepository, userBrokerRepository, imageBrokerRepository))

	// Transactions
	transactions.ReplaceGlobals(transactions.NewPostgresRepository(dbClient))
}
