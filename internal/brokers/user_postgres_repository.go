package brokers

import (
	"github.com/Zapharaos/fihub-backend/internal/models"
	"github.com/Zapharaos/fihub-backend/internal/utils"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// UserPostgresRepository is a postgres interface for UserRepository
type UserPostgresRepository struct {
	conn *sqlx.DB
}

// NewUserPostgresRepository returns a new instance of UserPostgresRepository
func NewUserPostgresRepository(dbClient *sqlx.DB) UserRepository {
	r := UserPostgresRepository{
		conn: dbClient,
	}
	var repo UserRepository = &r
	return repo
}

// Create use to create a BrokerUser
func (r *UserPostgresRepository) Create(userBroker models.BrokerUser) error {
	// Prepare query
	query := `INSERT INTO user_brokers (user_id, broker_id)
				VALUES (:user_id, :broker_id)`
	params := map[string]interface{}{
		"user_id":   userBroker.UserID,
		"broker_id": userBroker.Broker.ID,
	}

	// Execute query
	_, err := r.conn.NamedQuery(query, params)
	if err != nil {
		return err
	}

	return nil
}

// Delete use to delete a BrokerUser
func (r *UserPostgresRepository) Delete(userBroker models.BrokerUser) error {
	// Prepare query
	query := `DELETE FROM user_brokers as ub
			  WHERE ub.user_id = :user_id AND ub.broker_id = :broker_id`
	params := map[string]interface{}{
		"user_id":   userBroker.UserID,
		"broker_id": userBroker.Broker.ID,
	}

	// Execute query
	result, err := r.conn.NamedExec(query, params)
	if err != nil {
		return err
	}

	return utils.CheckRowAffected(result, 1)
}

// Exists use to check if a BrokerUser exists
func (r *UserPostgresRepository) Exists(userBroker models.BrokerUser) (bool, error) {
	// Prepare query
	query := `SELECT *
			  FROM user_brokers as ub
			  WHERE ub.user_id = :user_id AND ub.broker_id = :broker_id`
	params := map[string]interface{}{
		"user_id":   userBroker.UserID,
		"broker_id": userBroker.Broker.ID,
	}

	// Execute query
	rows, err := r.conn.NamedQuery(query, params)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	return rows.Next(), nil
}

// GetAll use to get all Broker for a BrokerUser
func (r *UserPostgresRepository) GetAll(userID uuid.UUID) ([]models.BrokerUser, error) {
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

	return utils.ScanAll(rows, r.Scan)
}

func (r *UserPostgresRepository) Scan(rows *sqlx.Rows) (models.BrokerUser, error) {
	var userBroker models.BrokerUser
	err := rows.Scan(
		&userBroker.Broker.ID,
		&userBroker.Broker.Name,
		&userBroker.Broker.ImageID,
	)
	if err != nil {
		return models.BrokerUser{}, err
	}

	return userBroker, nil
}
