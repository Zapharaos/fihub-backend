package postgres

import (
	"github.com/jmoiron/sqlx"
	"sync"
)

var (
	_globalMu sync.RWMutex
	_globalDB *sqlx.DB
)

// DB is the main accessor to the global postgresql client singleton
func DB() *sqlx.DB {
	_globalMu.RLock()
	db := _globalDB
	_globalMu.RUnlock()
	return db
}

// ReplaceGlobals replace the global postgresql client singleton with the provided one
func ReplaceGlobals(dbClient *sqlx.DB) func() {
	_globalMu.Lock()
	prev := _globalDB
	_globalDB = dbClient
	_globalMu.Unlock()
	return func() { ReplaceGlobals(prev) }
}
