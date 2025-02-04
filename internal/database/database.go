package database

import "github.com/jmoiron/sqlx"

type Databases struct {
	postgres *sqlx.DB
}

func NewDatabases(sql *sqlx.DB) Databases {
	return Databases{
		postgres: sql,
	}
}

func (db Databases) Postgres() *sqlx.DB {
	return db.postgres
}

var _databases Databases

func DB() Databases {
	return _databases
}

func ReplaceGlobals(databases Databases) func() {
	prev := _databases
	_databases = databases
	return func() { _databases = prev }
}
