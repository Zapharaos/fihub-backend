package users

import (
	"errors"
	"github.com/Zapharaos/fihub-backend/internal/utils"
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
func (r *PostgresRepository) Create(user UserWithPassword) (uuid.UUID, error) {

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
func (r *PostgresRepository) Get(userID uuid.UUID) (User, bool, error) {

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
		return User{}, false, err
	}
	defer rows.Close()

	return utils.ScanFirst(rows, scanUser)
}

// GetByEmail use to retrieve a user by email
func (r *PostgresRepository) GetByEmail(email string) (User, bool, error) {

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
		return User{}, false, err
	}
	defer rows.Close()

	return utils.ScanFirst(rows, scanUser)
}

// Exists checks if a User with requested email exists in the repository
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
func (r *PostgresRepository) Authenticate(email string, password string) (User, bool, error) {
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
		return User{}, false, err
	}
	defer rows.Close()

	// Retrieve user
	var userWithPassword UserWithPassword
	if rows.Next() {
		userWithPassword, err = scanUserWithPassword(rows)
		if err != nil {
			return User{}, false, err
		}

		err = bcrypt.CompareHashAndPassword([]byte(userWithPassword.Password), []byte(password))
		if err == nil {
			return userWithPassword.ToUser(), true, nil
		}
	}

	return User{}, false, errors.New("no user found, invalid credentials")
}

// Update method used to update a user
func (r *PostgresRepository) Update(user User) error {

	// Prepare query
	query := `UPDATE users as u
			  SET email = :email, updated_at = :updated_at
			  WHERE u.id = :id`
	params := map[string]interface{}{
		"id":         user.ID,
		"email":      user.Email,
		"updated_at": time.Now().Truncate(1 * time.Millisecond).UTC(),
	}

	// Execute query
	result, err := r.conn.NamedExec(query, params)
	if err != nil {
		return err
	}

	return utils.CheckRowAffected(result, 1)
}

// UpdateWithPassword method used to update a user with password
func (r *PostgresRepository) UpdateWithPassword(user UserWithPassword) error {

	// Hash password before saving
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Prepare query
	query := `UPDATE users as u
			  SET password = :password, updated_at = :updated_at
			  WHERE u.id = :id`
	params := map[string]interface{}{
		"id":         user.ID,
		"password":   hashedPassword,
		"updated_at": time.Now().Truncate(1 * time.Millisecond).UTC(),
	}

	// Execute query
	result, err := r.conn.NamedExec(query, params)
	if err != nil {
		return err
	}

	return utils.CheckRowAffected(result, 1)
}

// Delete method used to delete a user
func (r *PostgresRepository) Delete(userID uuid.UUID) error {

	// Prepare query
	query := `DELETE FROM users as u
			  WHERE u.id = :id`
	params := map[string]interface{}{
		"id": userID,
	}

	// Execute query
	result, err := r.conn.NamedExec(query, params)
	if err != nil {
		return err
	}

	return utils.CheckRowAffected(result, 1)
}

func scanUser(rows *sqlx.Rows) (User, error) {

	userWithPassword, err := scanUserWithPassword(rows)
	if err != nil {
		return User{}, err
	}

	return userWithPassword.ToUser(), nil
}

func scanUserWithPassword(rows *sqlx.Rows) (UserWithPassword, error) {
	var userWithPassword UserWithPassword
	err := rows.Scan(
		&userWithPassword.ID,
		&userWithPassword.Email,
		&userWithPassword.Password,
		&userWithPassword.CreatedAt,
		&userWithPassword.UpdatedAt,
	)
	if err != nil {
		return UserWithPassword{}, err
	}

	return userWithPassword, nil
}
