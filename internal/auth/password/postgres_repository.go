package password

import (
	"database/sql"
	"errors"
	"github.com/Zapharaos/fihub-backend/internal/utils"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"time"
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

func (p PostgresRepository) Create(request Request) (Request, error) {
	// Prepare query
	query := `INSERT INTO password_reset_tokens (user_id, token, expires_at)
				VALUES (:user_id, :token, :expires_at)
				RETURNING id, user_id, token, expires_at, created_at`
	params := map[string]interface{}{
		"user_id":    request.UserID,
		"token":      request.Token,
		"expires_at": request.ExpiresAt,
	}

	// Execute query
	rows, err := p.conn.NamedQuery(query, params)
	if err != nil {
		return Request{}, err
	}
	defer rows.Close()

	// Scan result
	result, _, err := utils.ScanFirst(rows, p.Scan)

	return result, err
}

func (p PostgresRepository) GetRequestID(userID uuid.UUID, token string) (uuid.UUID, error) {
	// Prepare query
	query := `SELECT id
			  FROM password_reset_tokens as p
			  WHERE p.user_id = $1 AND p.token = $2 AND p.expires_at > NOW()
			  LIMIT 1`

	// Execute query
	var requestID uuid.UUID
	err := p.conn.Get(&requestID, query, userID, token)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return uuid.Nil, nil
		}
		return uuid.Nil, err
	}

	return requestID, nil
}

// GetExpiresAt retrieves the expiration time for the existing password reset request for a given user.
func (p PostgresRepository) GetExpiresAt(userID uuid.UUID) (time.Time, error) {
	// Prepare query
	query := `SELECT expires_at
              FROM password_reset_tokens
              WHERE user_id = $1 AND expires_at > NOW()
              LIMIT 1`

	// Execute query
	var expiresAt time.Time
	err := p.conn.Get(&expiresAt, query, userID)
	if err != nil {
		return time.Time{}, err
	}

	return expiresAt, nil
}

func (p PostgresRepository) Delete(requestID uuid.UUID) error {
	// Prepare query
	query := `DELETE FROM password_reset_tokens as p
			  WHERE p.id = :id`
	params := map[string]interface{}{
		"id": requestID,
	}

	// Execute query
	result, err := p.conn.NamedExec(query, params)
	if err != nil {
		return err
	}

	return utils.CheckRowAffected(result, 1)
}

func (p PostgresRepository) Valid(userID uuid.UUID, requestID uuid.UUID) (bool, error) {

	// Prepare query
	query := `SELECT *
			  FROM password_reset_tokens as p
			  WHERE p.id = :id AND p.user_id = :user_id AND p.expires_at > NOW()`
	params := map[string]interface{}{
		"id":      requestID,
		"user_id": userID,
	}

	// Execute query
	rows, err := p.conn.NamedQuery(query, params)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	return rows.Next(), nil
}

func (p PostgresRepository) ValidForUser(userID uuid.UUID) (bool, error) {
	// Prepare query
	query := `SELECT *
			  FROM password_reset_tokens as p
			  WHERE p.user_id = :user_id AND p.expires_at > NOW()`
	params := map[string]interface{}{
		"user_id": userID,
	}

	// Execute query
	rows, err := p.conn.NamedQuery(query, params)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	return rows.Next(), nil
}

func (p PostgresRepository) Scan(rows *sqlx.Rows) (Request, error) {
	var request Request
	err := rows.Scan(
		&request.ID,
		&request.UserID,
		&request.Token,
		&request.ExpiresAt,
		&request.CreatedAt,
	)
	if err != nil {
		return Request{}, err
	}

	return request, nil
}
