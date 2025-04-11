package app

import (
	"github.com/Zapharaos/fihub-backend/internal/auth/password"
	"github.com/Zapharaos/fihub-backend/internal/auth/permissions"
	"github.com/Zapharaos/fihub-backend/internal/auth/roles"
	"github.com/Zapharaos/fihub-backend/internal/auth/users"
	"github.com/Zapharaos/fihub-backend/internal/brokers"
	"github.com/Zapharaos/fihub-backend/internal/database"
	"github.com/Zapharaos/fihub-backend/internal/transactions"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
)

// initDatabase initializes the database connections.
func initDatabase() {
	postgres := database.NewPostgresDB(database.NewSqlDatabase(database.SqlCredentials{
		Host:     viper.GetString("POSTGRES_HOST"),
		Port:     viper.GetString("POSTGRES_PORT"),
		User:     viper.GetString("POSTGRES_USER"),
		Password: viper.GetString("POSTGRES_PASSWORD"),
		DbName:   viper.GetString("POSTGRES_DB"),
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
