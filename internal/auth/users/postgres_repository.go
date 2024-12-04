package users

import (
	"errors"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
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

// Create method used to create a user
func (r *PostgresRepository) Create(user *UserWithPassword) (uuid.UUID, error) {

	// UUID
	userID := uuid.New()

	// Hash password before saving
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return uuid.UUID{}, err
	}

	// Get timestamps
	creationTS := time.Now().Truncate(1 * time.Millisecond).UTC()
	updateTS := creationTS

	// Prepare query
	query := `INSERT INTO users (id, email, password, created_at, updated_at)
				VALUES (:id, :email, :password, :created_at, :updated_at)`
	params := map[string]interface{}{
		"id":         userID,
		"email":      user.Email,
		"password":   hashedPassword,
		"created_at": creationTS,
		"updated_at": updateTS,
	}

	// Execute query
	_, err = r.conn.NamedQuery(query, params)
	if err != nil {
		return uuid.UUID{}, err
	}

	return userID, nil
}

// Get use to retrieve a user by id
func (r *PostgresRepository) Get(userID uuid.UUID) (*User, error) {

	// Prepare query
	query := `SELECT *
			  FROM users as u
			  WHERE u.id = :id`
	params := map[string]interface{}{
		"id": userID,
	}

	// Execute query
	rows, err := r.conn.NamedQuery(query, params)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Retrieve user
	var user *User
	if rows.Next() {
		user, err = scanUser(rows)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, nil
	}

	return user, nil
}

// Exists checks if a User exists in the repository
func (r *PostgresRepository) Exists(email string) (bool, error) {
	// Prepare query
	query := `SELECT *
			  FROM users as u
			  WHERE u.email = :email`
	params := map[string]interface{}{
		"email": email,
	}

	// Execute query
	rows, err := r.conn.NamedQuery(query, params)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	return rows.Next(), nil
}

// Authenticate returns a User from the repository by its login and password
func (r *PostgresRepository) Authenticate(email string, password string) (*User, error) {
	// Prepare query
	query := `SELECT *
			  FROM users as u
			  WHERE u.email = :email`
	params := map[string]interface{}{
		"email": email,
	}

	// Execute query
	rows, err := r.conn.NamedQuery(query, params)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Retrieve user
	var userWithPassword *UserWithPassword
	if rows.Next() {
		userWithPassword, err = scanUserWithPassword(rows)
		if err != nil {
			return nil, err
		}

		err = bcrypt.CompareHashAndPassword([]byte(userWithPassword.Password), []byte(password))
		if err == nil {
			return userWithPassword.ToUser(), nil
		}
	}

	return nil, errors.New("no user found, invalid credentials")
}

func scanUser(rows *sqlx.Rows) (*User, error) {

	userWithPassword, err := scanUserWithPassword(rows)
	if err != nil {
		return nil, err
	}

	return userWithPassword.ToUser(), nil
}

func scanUserWithPassword(rows *sqlx.Rows) (*UserWithPassword, error) {
	var userWithPassword UserWithPassword
	err := rows.Scan(
		&userWithPassword.ID,
		&userWithPassword.Email,
		&userWithPassword.Password,
		&userWithPassword.CreatedAt,
		&userWithPassword.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &userWithPassword, nil
}
