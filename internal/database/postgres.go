package postgres

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
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
	log.Println(psqlInfo)
	db, err := sqlx.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal("DbConnection.Open:", err)
		return nil, err
	}
	if err = db.Ping(); err != nil {
		log.Fatal("DbConnection.Ping:", err)
		return nil, err
	}
	return db, nil
}
