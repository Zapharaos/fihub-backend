package postgres

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

// Credentials : Give All DB Information.
type Credentials struct {
	Host     string `json:"url,omitempty"`
	Port     string `json:"port,omitempty"`
	User     string `json:"user,omitempty"`
	Password string `json:"password,omitempty"`
	DbName   string `json:"dbname,omitempty"`
}

// DbConnection : init DB access.
func DbConnection(credentials Credentials) (*sqlx.DB, error) {

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		credentials.Host,
		credentials.Port,
		credentials.User,
		credentials.Password,
		credentials.DbName)

	// Connect
	db, err := sqlx.Open("postgres", psqlInfo)
	if err != nil {
		zap.L().Error("DbConnection.Open:", zap.Error(err))
		return nil, err
	}

	// Ping for verification
	if err = db.Ping(); err != nil {
		zap.L().Error("DbConnection.Ping:", zap.Error(err))
		return nil, err
	}

	return db, nil
}
