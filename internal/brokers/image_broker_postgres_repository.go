package brokers

import (
	"github.com/Zapharaos/fihub-backend/internal/utils"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// ImageBrokerPostgresRepository is a repository containing the Issue definition based on a PSQL database and
// implementing the repository interface
type ImageBrokerPostgresRepository struct {
	conn *sqlx.DB
}

// NewImageBrokerPostgresRepository returns a new instance of ImageBrokerPostgresRepository
func NewImageBrokerPostgresRepository(dbClient *sqlx.DB) ImageBrokerRepository {
	r := ImageBrokerPostgresRepository{
		conn: dbClient,
	}
	var repo ImageBrokerRepository = &r
	return repo
}

// Create creates a new Image in the repository
func (r *ImageBrokerPostgresRepository) Create(brokerImage BrokerImage) error {

	// Prepare query
	query := `INSERT INTO broker_image (id, broker_id, name, data)
				VALUES (:id, :broker_id, :name, :data)`
	params := map[string]interface{}{
		"id":        brokerImage.ID,
		"broker_id": brokerImage.BrokerID,
		"name":      brokerImage.Name,
		"data":      brokerImage.Data,
	}

	// Execute query
	_, err := r.conn.NamedQuery(query, params)
	if err != nil {
		return err
	}

	return nil
}

// Get searches and returns an Image from the repository by its id
func (r *ImageBrokerPostgresRepository) Get(brokerImageID uuid.UUID) (BrokerImage, bool, error) {
	// Prepare query
	query := `SELECT *
			  FROM broker_image as bi
			  WHERE bi.id = :id`
	params := map[string]interface{}{
		"id": brokerImageID,
	}

	// Execute query
	rows, err := r.conn.NamedQuery(query, params)
	if err != nil {
		return BrokerImage{}, false, err
	}
	defer rows.Close()

	return utils.ScanFirst(rows, scanBrokerImage)
}

// Update updates an Image in the repository
func (r *ImageBrokerPostgresRepository) Update(brokerImage BrokerImage) error {
	// Prepare query
	query := `UPDATE broker_image as bi
			  SET broker_id = :broker_id, name = :name, data = :data
			  WHERE bi.id = :id`
	params := map[string]interface{}{
		"id":        brokerImage.ID,
		"broker_id": brokerImage.BrokerID,
		"name":      brokerImage.Name,
		"data":      brokerImage.Data,
	}

	// Execute query
	result, err := r.conn.NamedExec(query, params)
	if err != nil {
		return err
	}

	return utils.CheckRowAffected(result, 1)
}

// Delete deletes an Image from the repository
func (r *ImageBrokerPostgresRepository) Delete(brokerImageID uuid.UUID) error {
	// Prepare query
	query := `DELETE FROM broker_image as bi
			  WHERE bi.id = :id`
	params := map[string]interface{}{
		"id": brokerImageID,
	}

	// Execute query
	result, err := r.conn.NamedExec(query, params)
	if err != nil {
		return err
	}

	return utils.CheckRowAffected(result, 1)
}

// Exists checks if an Image exists in the repository
func (r *ImageBrokerPostgresRepository) Exists(brokerID uuid.UUID, brokerImageID uuid.UUID) (bool, error) {
	// Prepare query
	query := `SELECT *
			  FROM broker_image as bi
			  WHERE bi.id = :id AND bi.broker_id = :broker_id`
	params := map[string]interface{}{
		"id":        brokerImageID,
		"broker_id": brokerID,
	}

	// Execute query
	rows, err := r.conn.NamedQuery(query, params)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	return rows.Next(), nil
}

func scanBrokerImage(rows utils.RowScanner) (BrokerImage, error) {
	var brokerImage BrokerImage
	err := rows.Scan(
		&brokerImage.ID,
		&brokerImage.BrokerID,
		&brokerImage.Name,
		&brokerImage.Data,
	)
	if err != nil {
		return BrokerImage{}, err
	}

	return brokerImage, nil
}
