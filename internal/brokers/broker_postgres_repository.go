package brokers

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// BrokerPostgresRepository is a repository containing the Issue definition based on a PSQL database and
// implementing the repository interface
type BrokerPostgresRepository struct {
	conn *sqlx.DB
}

// NewPostgresRepository returns a new instance of PostgresRepository
func NewPostgresRepository(dbClient *sqlx.DB) BrokerRepository {
	r := BrokerPostgresRepository{
		conn: dbClient,
	}
	var repo BrokerRepository = &r
	return repo
}

func (r *BrokerPostgresRepository) Exists(id uuid.UUID) (bool, error) {
	// Prepare query
	query := `SELECT *
			  FROM brokers as b
			  WHERE b.id = :id`
	params := map[string]interface{}{
		"id": id,
	}

	// Execute query
	rows, err := r.conn.NamedQuery(query, params)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	return rows.Next(), nil
}

// GetAll use to retrieve all brokers
func (r *BrokerPostgresRepository) GetAll() ([]Broker, error) {

	// Prepare query
	query := `SELECT *
			  FROM brokers`

	// Execute query
	rows, err := r.conn.Queryx(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Retrieve all brokers
	var brokers []Broker
	for rows.Next() {
		broker, err := scanBroker(rows)
		if err != nil {
			return nil, err
		}
		brokers = append(brokers, broker)
	}

	// Handle error
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return brokers, nil
}

func scanBroker(rows *sqlx.Rows) (Broker, error) {
	var broker Broker
	err := rows.Scan(
		&broker.ID,
		&broker.Name,
	)
	if err != nil {
		return Broker{}, err
	}

	return broker, nil
}
