package repositories

import (
	"github.com/Zapharaos/fihub-backend/internal/models"
	"github.com/Zapharaos/fihub-backend/internal/utils"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// ImagePostgresRepository is a postgres interface for ImageRepository
type ImagePostgresRepository struct {
	conn *sqlx.DB
}

// NewImagePostgresRepository returns a new instance of ImagePostgresRepository
func NewImagePostgresRepository(dbClient *sqlx.DB) ImageRepository {
	r := ImagePostgresRepository{
		conn: dbClient,
	}
	var repo ImageRepository = &r
	return repo
}

// Create creates a new BrokerImage in the repository
func (r *ImagePostgresRepository) Create(brokerImage models.BrokerImage) error {

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

// Get searches and returns an BrokerImage from the repository by its id
func (r *ImagePostgresRepository) Get(brokerImageID uuid.UUID) (models.BrokerImage, bool, error) {
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
		return models.BrokerImage{}, false, err
	}
	defer rows.Close()

	return utils.ScanFirst(rows, r.Scan)
}

// Update updates an BrokerImage in the repository
func (r *ImagePostgresRepository) Update(brokerImage models.BrokerImage) error {
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

// Delete deletes an BrokerImage from the repository
func (r *ImagePostgresRepository) Delete(brokerImageID uuid.UUID) error {
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

// Exists checks if an BrokerImage exists in the repository
func (r *ImagePostgresRepository) Exists(brokerID uuid.UUID, brokerImageID uuid.UUID) (bool, error) {
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

func (r *ImagePostgresRepository) Scan(rows *sqlx.Rows) (models.BrokerImage, error) {
	var brokerImage models.BrokerImage
	err := rows.Scan(
		&brokerImage.ID,
		&brokerImage.BrokerID,
		&brokerImage.Name,
		&brokerImage.Data,
	)
	if err != nil {
		return models.BrokerImage{}, err
	}

	return brokerImage, nil
}
