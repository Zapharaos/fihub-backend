package brokers

import (
	"github.com/Zapharaos/fihub-backend/internal/utils"
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

// Create use to create a broker
func (r *BrokerPostgresRepository) Create(broker Broker) (uuid.UUID, error) {

	// UUID
	broker.ID = uuid.New()

	// Prepare query
	query := `INSERT INTO brokers (id, name)
			  VALUES (:id, :name)`
	params := map[string]interface{}{
		"id":   broker.ID,
		"name": broker.Name,
	}

	// Execute query
	_, err := r.conn.NamedQuery(query, params)
	if err != nil {
		return uuid.Nil, err
	}

	return broker.ID, nil
}

// Get use to retrieve a broker
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

	// Retrieve broker
	if rows.Next() {
		broker, err := scanBroker(rows)
		if err != nil {
			return Broker{}, false, err
		}
		return broker, true, nil
	}

	return Broker{}, false, nil
}

// Update use to update a broker
func (r *BrokerPostgresRepository) Update(broker Broker) error {

	// Prepare query
	query := `UPDATE brokers
			  SET name = :name
			  WHERE id = :id`
	params := map[string]interface{}{
		"id":   broker.ID,
		"name": broker.Name,
	}

	// Execute query
	result, err := r.conn.NamedExec(query, params)
	if err != nil {
		return err
	}

	return utils.CheckRowAffected(result, 1)
}

// Delete use to delete a broker
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

// Exists use to check if a broker exists
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

// ExistsByName use to check if a broker exists with a given name
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

// SetImage use to set an image to a broker
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

// HasImage use to check if a broker has an image
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

// DeleteImage use to delete an image from a broker
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
	)
	if err != nil {
		return Broker{}, err
	}

	return broker, nil
}
