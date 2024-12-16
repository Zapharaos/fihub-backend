package brokers

import (
	"github.com/Zapharaos/fihub-backend/internal/utils"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// UserBrokerPostgresRepository is a repository containing the Issue definition based on a PSQL database and
// implementing the repository interface
type UserBrokerPostgresRepository struct {
	conn *sqlx.DB
}

// NewUserBrokerPostgresRepository returns a new instance of UserBrokerPostgresRepository
func NewUserBrokerPostgresRepository(dbClient *sqlx.DB) UserBrokerRepository {
	r := UserBrokerPostgresRepository{
		conn: dbClient,
	}
	var repo UserBrokerRepository = &r
	return repo
}

func (r *UserBrokerPostgresRepository) Create(userBroker UserBroker) error {
	// Prepare query
	query := `INSERT INTO user_brokers (user_id, broker_id)
				VALUES (:user_id, :broker_id)`
	params := map[string]interface{}{
		"user_id":   userBroker.UserID,
		"broker_id": userBroker.BrokerID,
	}

	// Execute query
	_, err := r.conn.NamedQuery(query, params)
	if err != nil {
		return err
	}

	return nil
}

func (r *UserBrokerPostgresRepository) Delete(userBroker UserBroker) error {
	// Prepare query
	query := `DELETE FROM user_brokers as ub
			  WHERE ub.user_id = :user_id AND ub.broker_id = :broker_id`
	params := map[string]interface{}{
		"user_id":   userBroker.UserID,
		"broker_id": userBroker.BrokerID,
	}

	// Execute query
	result, err := r.conn.NamedExec(query, params)
	if err != nil {
		return err
	}

	return utils.CheckRowAffected(result, 1)
}

func (r *UserBrokerPostgresRepository) Exists(userBroker UserBroker) (bool, error) {
	// Prepare query
	query := `SELECT *
			  FROM user_brokers as ub
			  WHERE ub.user_id = :user_id AND ub.broker_id = :broker_id`
	params := map[string]interface{}{
		"user_id":   userBroker.UserID,
		"broker_id": userBroker.BrokerID,
	}

	// Execute query
	rows, err := r.conn.NamedQuery(query, params)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	return rows.Next(), nil
}

func (r *UserBrokerPostgresRepository) GetAll(userID uuid.UUID) ([]Broker, error) {
	// Prepare query
	query := `SELECT b.id, b.name, b.image_id
			  FROM user_brokers as ub
			  JOIN brokers AS b ON ub.broker_id = b.id
			  WHERE ub.user_id = :user_id`
	params := map[string]interface{}{
		"user_id": userID,
	}

	// Execute query
	rows, err := r.conn.NamedQuery(query, params)
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
