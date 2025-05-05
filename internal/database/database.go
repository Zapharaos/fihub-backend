package database

// Databases holds all the database connections
type Databases struct {
	postgres PostgresDB
	redis    RedisDB
}

// Postgres returns the postgres database connection
func (db Databases) Postgres() PostgresDB {
	return db.postgres
}

// Redis returns the redis database connection
func (db Databases) Redis() RedisDB {
	return db.redis
}

// SetPostgres sets the postgres database connection while preserving other connections
func (db Databases) SetPostgres(postgres PostgresDB) {
	db.postgres = postgres
}

// SetRedis sets the redis database connection while preserving other connections
func (db Databases) SetRedis(redis RedisDB) {
	db.redis = redis
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
