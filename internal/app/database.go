package app

import (
	"github.com/Zapharaos/fihub-backend/internal/database"
	"github.com/Zapharaos/fihub-backend/internal/users"
	"github.com/Zapharaos/fihub-backend/internal/users/password"
	"github.com/Zapharaos/fihub-backend/internal/users/permissions"
	"github.com/Zapharaos/fihub-backend/internal/users/roles"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
)

// InitDatabase initializes the database connections.
func InitDatabase() {
	postgres := database.NewPostgresDB(database.NewSqlDatabase(database.SqlCredentials{
		Host:     viper.GetString("POSTGRES_HOST"),
		Port:     viper.GetString("POSTGRES_PORT"),
		User:     viper.GetString("POSTGRES_USER"),
		Password: viper.GetString("POSTGRES_PASSWORD"),
		DbName:   viper.GetString("POSTGRES_DB"),
	}))
	database.ReplaceGlobals(database.NewDatabases(postgres))
}

// InitPostgres initializes the postgres repositories.
func InitPostgres(dbClient *sqlx.DB) {
	// Auth
	users.ReplaceGlobals(users.NewPostgresRepository(dbClient))
	password.ReplaceGlobals(password.NewPostgresRepository(dbClient))

	// Roles
	roles.ReplaceGlobals(roles.NewPostgresRepository(dbClient))

	// Permissions
	permissions.ReplaceGlobals(permissions.NewPostgresRepository(dbClient))
}
