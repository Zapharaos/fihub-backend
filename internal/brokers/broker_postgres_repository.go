package brokers

import (
	"github.com/Zapharaos/fihub-backend/internal/utils"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// BrokerPostgresRepository is a postgres interface for BrokerRepository
type BrokerPostgresRepository struct {
	conn *sqlx.DB
}

// NewPostgresRepository returns a new instance of BrokerRepository
func NewPostgresRepository(dbClient *sqlx.DB) BrokerRepository {
	r := BrokerPostgresRepository{
		conn: dbClient,
	}
	var repo BrokerRepository = &r
	return repo
}

// Create use to create a Broker
func (r *BrokerPostgresRepository) Create(broker Broker) (uuid.UUID, error) {

	// UUID
	broker.ID = uuid.New()

	// Prepare query
	query := `INSERT INTO brokers (id, name, disabled)
			  VALUES (:id, :name, :disabled)`
	params := map[string]interface{}{
		"id":       broker.ID,
		"name":     broker.Name,
		"disabled": broker.Disabled,
	}

	// Execute query
	_, err := r.conn.NamedQuery(query, params)
	if err != nil {
		return uuid.Nil, err
	}

	return broker.ID, nil
}

// Get use to retrieve a Broker by its id
func (r *BrokerPostgresRepository) Get(id uuid.UUID) (Broker, bool, error) {

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
		return Broker{}, false, err
	}
	defer rows.Close()

	return utils.ScanFirst(rows, scanBroker)
}

// Update use to update a Broker
func (r *BrokerPostgresRepository) Update(broker Broker) error {

	// Prepare query
	query := `UPDATE brokers
			  SET name = :name, disabled = :disabled
			  WHERE id = :id`
	params := map[string]interface{}{
		"id":       broker.ID,
		"name":     broker.Name,
		"disabled": broker.Disabled,
	}

	// Execute query
	result, err := r.conn.NamedExec(query, params)
	if err != nil {
		return err
	}

	return utils.CheckRowAffected(result, 1)
}

// Delete use to delete a Broker
func (r *BrokerPostgresRepository) Delete(id uuid.UUID) error {

	// Prepare query
	query := `DELETE FROM brokers
			  WHERE id = :id`
	params := map[string]interface{}{
		"id": id,
	}

	// Execute query
	result, err := r.conn.NamedExec(query, params)
	if err != nil {
		return err
	}

	return utils.CheckRowAffected(result, 1)
}

// Exists use to check if a Broker exists
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

// ExistsByName use to check if a Broker exists with a given name
func (r *BrokerPostgresRepository) ExistsByName(name string) (bool, error) {
	// Prepare query
	query := `SELECT *
			  FROM brokers as b
			  WHERE b.name = :name`
	params := map[string]interface{}{
		"name": name,
	}

	// Execute query
	rows, err := r.conn.NamedQuery(query, params)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	return rows.Next(), nil
}

// GetAll use to retrieve all existing Broker
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

	return utils.ScanAll(rows, scanBroker)
}

// GetAllEnabled use to retrieve all enabled Broker
func (r *BrokerPostgresRepository) GetAllEnabled() ([]Broker, error) {

	// Prepare query
	query := `SELECT *
			  FROM brokers
			  WHERE disabled = false`

	// Execute query
	rows, err := r.conn.Queryx(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return utils.ScanAll(rows, scanBroker)
}

// SetImage use to set an Image to a Broker
func (r *BrokerPostgresRepository) SetImage(id uuid.UUID, imageID uuid.UUID) error {
	// Prepare query
	query := `UPDATE brokers
			  SET image_id = :image_id
			  WHERE id = :id`
	params := map[string]interface{}{
		"id":       id,
		"image_id": imageID,
	}

	// Execute query
	result, err := r.conn.NamedExec(query, params)
	if err != nil {
		return err
	}

	return utils.CheckRowAffected(result, 1)
}

// HasImage use to check if a Broker has an Image
func (r *BrokerPostgresRepository) HasImage(id uuid.UUID) (bool, error) {
	// Prepare query
	query := `SELECT *
			  FROM brokers as b
			  WHERE b.id = :id AND b.image_id IS NOT NULL`
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

// DeleteImage use to delete an Image from a Broker
func (r *BrokerPostgresRepository) DeleteImage(id uuid.UUID) error {
	// Prepare query
	query := `UPDATE brokers
			  SET image_id = NULL
			  WHERE id = :id`
	params := map[string]interface{}{
		"id": id,
	}

	// Execute query
	result, err := r.conn.NamedExec(query, params)
	if err != nil {
		return err
	}

	return utils.CheckRowAffected(result, 1)
}

func scanBroker(rows *sqlx.Rows) (Broker, error) {
	var broker Broker
	err := rows.Scan(
		&broker.ID,
		&broker.Name,
		&broker.ImageID,
		&broker.Disabled,
	)
	if err != nil {
		return Broker{}, err
	}

	return broker, nil
}
