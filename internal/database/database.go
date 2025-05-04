package database

// Databases holds all the database connections
type Databases struct {
	postgres PostgresDB
}

// NewDatabases creates a new Databases struct
func NewDatabases(postgres PostgresDB) Databases {
	return Databases{
		postgres: postgres,
	}
}

// Postgres returns the postgres database connection
func (db Databases) Postgres() PostgresDB {
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
