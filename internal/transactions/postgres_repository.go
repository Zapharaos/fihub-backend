package transactions

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// PostgresRepository is a repository containing the Issue definition based on a PSQL database and
// implementing the repository interface
type PostgresRepository struct {
	conn *sqlx.DB
}

// NewPostgresRepository returns a new instance of PostgresRepository
func NewPostgresRepository(dbClient *sqlx.DB) Repository {
	r := PostgresRepository{
		conn: dbClient,
	}
	var repo Repository = &r
	return repo
}

func (p PostgresRepository) Create(transaction Transaction) error {
	//TODO implement me
	panic("implement me")
}

func (p PostgresRepository) Get(transaction Transaction) error {
	//TODO implement me
	panic("implement me")
}

func (p PostgresRepository) Update(transaction Transaction) error {
	//TODO implement me
	panic("implement me")
}

func (p PostgresRepository) Delete(transaction Transaction) error {
	//TODO implement me
	panic("implement me")
}

func (p PostgresRepository) Exists(transaction Transaction) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (p PostgresRepository) GetAll(userID uuid.UUID) ([]Transaction, error) {
	//TODO implement me
	panic("implement me")
}
