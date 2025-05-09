package repositories

import (
	"github.com/Zapharaos/fihub-backend/internal/models"
	"github.com/Zapharaos/fihub-backend/internal/utils"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// PostgresRepository is a postgres interface for Repository
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

// Create use to create a Transaction
func (r PostgresRepository) Create(transactionInput models.TransactionInput) (uuid.UUID, error) {

	// Prepare query
	query := `INSERT INTO transactions (id, user_id, broker_id, date, transaction_type, asset, quantity, price, price_unit, fee)
			  VALUES (:id, :user_id, :broker_id, :date, :transaction_type, :asset, :quantity, :price, :price_unit, :fee)
			  RETURNING id`

	// Create parameter map
	params := map[string]interface{}{
		"id":               uuid.New(),
		"user_id":          transactionInput.UserID,
		"broker_id":        transactionInput.BrokerID,
		"date":             transactionInput.Date,
		"transaction_type": transactionInput.Type,
		"asset":            transactionInput.Asset,
		"quantity":         transactionInput.Quantity,
		"price":            transactionInput.Price,
		"price_unit":       transactionInput.PriceUnit,
		"fee":              transactionInput.Fee,
	}

	// Execute query
	rows, err := r.conn.NamedQuery(query, params)
	if err != nil {
		return uuid.Nil, err
	}
	defer rows.Close()

	// Retrieve the created transaction ID
	var id uuid.UUID
	if rows.Next() {
		if err := rows.Scan(&id); err != nil {
			return uuid.Nil, err
		}
		return id, nil
	}

	return id, nil
}

// Get use to retrieve a Transaction by its id
func (r PostgresRepository) Get(transactionID uuid.UUID) (models.Transaction, bool, error) {

	// Prepare query
	query := `SELECT b.id AS "broker.id", b.name AS "broker.name", b.image_id AS "broker.image_id",
       			t.id, t.user_id, t.date, t.transaction_type, t.asset, t.quantity, t.price, t.price_unit, t.fee
			  FROM transactions as t
			  JOIN brokers as b ON t.broker_id = b.id
			  WHERE t.id = :id`
	params := map[string]interface{}{
		"id": transactionID,
	}

	// Execute query
	rows, err := r.conn.NamedQuery(query, params)
	if err != nil {
		return models.Transaction{}, false, err
	}
	defer rows.Close()

	return utils.ScanFirstStruct[models.Transaction](rows)
}

// Update use to update a Transaction
func (r PostgresRepository) Update(transactionInput models.TransactionInput) error {
	// Prepare query
	query := `UPDATE transactions
			  SET broker_id = :broker_id,
			      date = :date,
				  transaction_type = :transaction_type,
				  asset = :asset,
				  quantity = :quantity,
				  price = :price,
				  price_unit = :price_unit,
				  fee = :fee
			  WHERE id = :id`
	params := map[string]interface{}{
		"id":               transactionInput.ID,
		"broker_id":        transactionInput.BrokerID,
		"date":             transactionInput.Date,
		"transaction_type": transactionInput.Type,
		"asset":            transactionInput.Asset,
		"quantity":         transactionInput.Quantity,
		"price":            transactionInput.Price,
		"price_unit":       transactionInput.PriceUnit,
		"fee":              transactionInput.Fee,
	}

	// Execute query
	result, err := r.conn.NamedExec(query, params)
	if err != nil {
		return err
	}

	return utils.CheckRowAffected(result, 1)
}

// Delete use to delete a Transaction
func (r PostgresRepository) Delete(transaction models.Transaction) error {
	// Prepare query
	query := `DELETE FROM transactions as t WHERE t.id = :id`
	params := map[string]interface{}{
		"id": transaction.ID,
	}

	// Execute query
	result, err := r.conn.NamedExec(query, params)
	if err != nil {
		return err
	}

	return utils.CheckRowAffected(result, 1)
}

// DeleteByBroker use to delete a Transaction
func (r PostgresRepository) DeleteByBroker(transaction models.Transaction) error {
	// Prepare query
	query := `DELETE FROM transactions as t WHERE t.broker_id = :broker_id`
	params := map[string]interface{}{
		"broker_id": transaction.Broker.ID,
	}

	// Execute query
	_, err := r.conn.NamedExec(query, params)
	if err != nil {
		return err
	}

	return nil
}

// Exists use to check if a Transaction exists
func (r PostgresRepository) Exists(transactionID uuid.UUID, userID uuid.UUID) (bool, error) {
	query := `SELECT id
			  FROM transactions
			  WHERE id = :id AND user_id = :user_id`
	params := map[string]interface{}{
		"id":      transactionID,
		"user_id": userID,
	}

	// Execute query
	rows, err := r.conn.NamedQuery(query, params)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	return rows.Next(), nil
}

// GetAll use to retrieve all Transactions
func (r PostgresRepository) GetAll(userID uuid.UUID) ([]models.Transaction, error) {

	// Prepare query
	query := `SELECT b.id AS "broker.id", b.name AS "broker.name", b.image_id AS "broker.image_id",
       			t.id, t.user_id, t.date, t.transaction_type, t.asset, t.quantity, t.price, t.price_unit, t.fee
			  FROM transactions as t
			  JOIN brokers as b ON t.broker_id = b.id
			  WHERE t.user_id = :user_id`
	params := map[string]interface{}{
		"user_id": userID,
	}

	// Execute query
	rows, err := r.conn.NamedQuery(query, params)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return utils.ScanAllStruct[models.Transaction](rows)
}
