package password

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

func (p PostgresRepository) Create(request Request) error {
	// Prepare query
	query := `INSERT INTO password_reset_tokens (user_id, token, expires_at)
				VALUES (:user_id, :token, :expires_at)`
	params := map[string]interface{}{
		"user_id":    request.UserID,
		"token":      request.Token,
		"expires_at": request.ExpiresAt,
	}

	// Execute query
	_, err := p.conn.NamedQuery(query, params)

	return err
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
		return uuid.Nil, err
	}

	return requestID, nil
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
