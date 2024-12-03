package app

import (
	"github.com/Zapharaos/fihub-backend/internal/auth/users"
	"github.com/Zapharaos/fihub-backend/internal/database"
	"github.com/Zapharaos/fihub-backend/pkg/env"
	"log"
)

// Init initialize all the app configuration and components
func Init() {

	// Load the .env file
	err := env.Load()
	if err != nil {
		log.Fatal(err)
		return
	}

	initPostgres()
	initRepositories()
}

func initPostgres() {
	credentials := postgres.Credentials{
		Host:     env.GetString("POSTGRES_HOST", "host"),
		Port:     env.GetString("POSTGRES_PORT", "port"),
		DbName:   env.GetString("POSTGRES_DB", "database_name"),
		User:     env.GetString("POSTGRES_USER", "user"),
		Password: env.GetString("POSTGRES_PASSWORD", "password"),
	}
	dbClient, err := postgres.DbConnection(credentials)
	if err != nil {
		log.Fatal("main.DbConnection:", err)
	}
	dbClient.SetMaxOpenConns(env.GetInt("POSTGRES_MAX_OPEN_CONNS", 30))
	dbClient.SetMaxIdleConns(env.GetInt("POSTGRES_MAX_IDLE_CONNS", 30))
	postgres.ReplaceGlobals(dbClient)
}

func initRepositories() {
	dbClient := postgres.DB()
	users.ReplaceGlobals(users.NewPostgresRepository(dbClient))
}
