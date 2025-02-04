package database

import "github.com/jmoiron/sqlx"

// Databases holds all the database connections
type Databases struct {
	postgres *sqlx.DB
}

// NewDatabases creates a new Databases struct
func NewDatabases(sql *sqlx.DB) Databases {
	return Databases{
		postgres: sql,
	}
}

// Postgres returns the postgres database connection
func (db Databases) Postgres() *sqlx.DB {
	return db.postgres
}

var _databases Databases

// DB returns the global Databases struct
func DB() Databases {
	return _databases
}

// ReplaceGlobals replaces the global Databases struct with the provided one
func ReplaceGlobals(databases Databases) func() {
	prev := _databases
	_databases = databases
	return func() { _databases = prev }
}
