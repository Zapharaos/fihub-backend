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

func (r PostgresRepository) Create(transaction Transaction) (uuid.UUID, error) {

	// UUID
	transaction.ID = uuid.New()

	// Prepare query
	query := `INSERT INTO transactions (id, user_id, broker_id, date, transaction_type, asset, quantity, price, price_unit, fee)
			  VALUES (:id, :user_id, :broker_id, :date, :transaction_type, :asset, :quantity, :price, :price_unit, :fee)`
	params := map[string]interface{}{
		"id":               transaction.ID,
		"user_id":          transaction.UserID,
		"broker_id":        transaction.BrokerID,
		"date":             transaction.Date,
		"transaction_type": transaction.Type,
		"asset":            transaction.Asset,
		"quantity":         transaction.Quantity,
		"price":            transaction.Price,
		"price_unit":       transaction.PriceUnit,
		"fee":              transaction.Fee,
	}

	// Execute query
	_, err := r.conn.NamedQuery(query, params)
	if err != nil {
		return uuid.UUID{}, err
	}

	return transaction.ID, nil
}

func (r PostgresRepository) Get(transactionID uuid.UUID) (Transaction, bool, error) {

	// Prepare query
	query := `SELECT *
			  FROM transactions as t
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

func (r PostgresRepository) Delete(transaction Transaction) error {
	// Prepare query
	query := `DELETE FROM transactions as t WHERE t.id = $1`

	// Execute query
	_, err := r.conn.Exec(query, transaction.ID)
	if err != nil {
		return err
	}
	return nil
}

func (r PostgresRepository) Exists(transactionID uuid.UUID) (bool, error) {
	query := `SELECT COUNT(1) FROM transactions WHERE id = $1`

	var count int
	err := r.conn.Get(&count, query, transactionID)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r PostgresRepository) GetAll(userID uuid.UUID) ([]Transaction, error) {
	query := `SELECT *
			  FROM transactions
			  WHERE user_id = :user_id`
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
		&transaction.BrokerID,
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
