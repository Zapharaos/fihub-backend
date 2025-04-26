package app

import (
	userrepositories "github.com/Zapharaos/fihub-backend/cmd/user/app/repositories"
	"github.com/Zapharaos/fihub-backend/internal/database"
	"github.com/Zapharaos/fihub-backend/internal/password"
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
	userrepositories.ReplaceGlobals(userrepositories.NewPostgresRepository(dbClient))
	password.ReplaceGlobals(password.NewPostgresRepository(dbClient))
}
