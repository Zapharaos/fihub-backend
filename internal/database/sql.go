package database

// _ "github.com/lib/pq" // is used to register the PostgreSQL driver with the sqlx package.
import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// SqlDatabase is an interface for a SQL database.
type SqlDatabase interface {
	Connect() (*sqlx.DB, error)
}

// SqlCredentials holds the credentials for a SQL database.
type SqlCredentials struct {
	Host     string `json:"url,omitempty"`
	Port     string `json:"port,omitempty"`
	User     string `json:"user,omitempty"`
	Password string `json:"password,omitempty"`
	DbName   string `json:"dbname,omitempty"`
}

// Sql is a struct that holds the credentials for a SQL database.
type Sql struct {
	SqlCredentials
}

// NewSqlDatabase creates a new SQL database.
func NewSqlDatabase(credentials SqlCredentials) SqlDatabase {
	d := Sql{
		SqlCredentials: credentials,
	}
	var db SqlDatabase = &d
	return db
}

// Connect creates a new SQL connection.
func (s Sql) Connect() (*sqlx.DB, error) {

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		s.Host, s.Port, s.User, s.Password, s.DbName)

	// Connect
	db, err := sqlx.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	// Ping for verification
	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
