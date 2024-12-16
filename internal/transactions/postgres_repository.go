package transactions

import (
	"github.com/Zapharaos/fihub-backend/internal/utils"
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

func (r PostgresRepository) Create(transactionInput TransactionInput) (uuid.UUID, error) {

	// UUID
	transactionInput.ID = uuid.New()

	// Prepare query
	query := `INSERT INTO transactions (id, user_id, broker_id, date, transaction_type, asset, quantity, price, price_unit, fee)
			  VALUES (:id, :user_id, :broker_id, :date, :transaction_type, :asset, :quantity, :price, :price_unit, :fee)`
	params := map[string]interface{}{
		"id":               transactionInput.ID,
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
	_, err := r.conn.NamedQuery(query, params)
	if err != nil {
		return uuid.UUID{}, err
	}

	return transactionInput.ID, nil
}

func (r PostgresRepository) Get(transactionID uuid.UUID) (Transaction, bool, error) {

	// Prepare query
	query := `SELECT t.id, t.user_id, b.id, b.name, b.image_id, t.date, t.transaction_type, t.asset, t.quantity, t.price, t.price_unit, t.fee
			  FROM transactions as t
			  JOIN brokers as b ON t.broker_id = b.id
			  WHERE t.id = :id`
	params := map[string]interface{}{
		"id": transactionID,
	}

	// Execute query
	rows, err := r.conn.NamedQuery(query, params)
	if err != nil {
		return Transaction{}, false, err
	}
	defer rows.Close()

	// Retrieve user
	var transaction Transaction
	if rows.Next() {
		transaction, err = scanTransaction(rows)
		if err != nil {
			return Transaction{}, false, err
		}
	} else {
		return Transaction{}, false, nil
	}

	return transaction, true, nil
}

func (r PostgresRepository) Update(transactionInput TransactionInput) error {
	// Prepare query
	query := `UPDATE transactions
			  SET date = :date,
				  transaction_type = :transaction_type,
				  asset = :asset,
				  quantity = :quantity,
				  price = :price,
				  price_unit = :price_unit,
				  fee = :fee
			  WHERE id = :id`
	params := map[string]interface{}{
		"id":               transactionInput.ID,
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

func (r PostgresRepository) Delete(transaction Transaction) error {
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

func (r PostgresRepository) Exists(transactionID uuid.UUID, userID uuid.UUID) (bool, error) {
	query := `SELECT COUNT(1)
			  FROM transactions
			  WHERE id = :id AND user_id = :user_id`
	params := map[string]interface{}{
		"id":      transactionID,
		"user_id": userID,
	}

	var count int
	err := r.conn.Get(&count, query, params)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r PostgresRepository) GetAll(userID uuid.UUID) ([]Transaction, error) {

	// Prepare query
	query := `SELECT t.id, t.user_id, b.id, b.name, b.image_id, t.date, t.transaction_type, t.asset, t.quantity, t.price, t.price_unit, t.fee
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

	// Retrieve transactions
	var transactions []Transaction
	for rows.Next() {
		transaction, err := scanTransaction(rows)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}

	// Handle error
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return transactions, nil
}

func scanTransaction(rows *sqlx.Rows) (Transaction, error) {
	var transaction Transaction
	err := rows.Scan(
		&transaction.ID,
		&transaction.UserID,
		&transaction.Broker.ID,
		&transaction.Broker.Name,
		&transaction.Broker.ImageID,
		&transaction.Date,
		&transaction.Type,
		&transaction.Asset,
		&transaction.Quantity,
		&transaction.Price,
		&transaction.PriceUnit,
		&transaction.Fee,
	)
	if err != nil {
		return Transaction{}, err
	}

	return transaction, nil
}
