package database

import (
	"context"
	"fmt"
	"github.com/Zapharaos/fihub-backend/pkg/env"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"time"
)

func New() (*sqlx.DB, error) {

	// Configure connection
	dataSourceName := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		env.GetString("POSTGRES_USER", "user"),
		env.GetString("POSTGRES_PASSWORD", "password"),
		env.GetString("POSTGRES_HOST", "host"),
		env.GetString("POSTGRES_PORT", "port"),
		env.GetString("POSTGRES_DB", "database_name"))

	// Open connection to database
	db, err := sqlx.Open("postgres", dataSourceName)
	if err != nil {
		return nil, err
	}

	// Connection settings
	db.SetMaxOpenConns(env.GetInt("POSTGRES_MAX_OPEN_CONNS", 30))
	db.SetMaxIdleConns(env.GetInt("POSTGRES_MAX_IDLE_CONNS", 30))

	duration, err := time.ParseDuration(env.GetString("POSTGRES_MAX_IDLE_TIME", "user"))
	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(duration)

	log.Println("Testing connection...")

	// Checks connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = db.PingContext(ctx); err != nil {
		return nil, err
	}

	log.Println("Database connection established")
	return db, nil
}

func Stop(db *sqlx.DB) {
	db.Close()

	log.Println("Database connection closed")
}
