package repositories

import (
	"errors"
	"github.com/Zapharaos/fihub-backend/internal/models"
	"github.com/Zapharaos/fihub-backend/internal/utils"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
	"time"
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

// Create method used to create a User
func (r *PostgresRepository) Create(user models.UserWithPassword) (uuid.UUID, error) {

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
	query := `INSERT INTO Users (ID, email, password, created_at, updated_at)
				VALUES (:ID, :email, :password, :created_at, :updated_at)`
	params := map[string]interface{}{
		"ID":         userID,
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

// Get use to retrieve a User by ID
func (r *PostgresRepository) Get(userID uuid.UUID) (models.User, bool, error) {

	// Prepare query
	query := `SELECT *
			  FROM Users as u
			  WHERE u.ID = :ID`
	params := map[string]interface{}{
		"ID": userID,
	}

	// Execute query
	rows, err := r.conn.NamedQuery(query, params)
	if err != nil {
		return models.User{}, false, err
	}
	defer rows.Close()

	return utils.ScanFirst(rows, r.Scan)
}

// GetByEmail use to retrieve a User by email
func (r *PostgresRepository) GetByEmail(email string) (models.User, bool, error) {

	// Prepare query
	query := `SELECT *
			  FROM Users as u
			  WHERE u.email = :email`
	params := map[string]interface{}{
		"email": email,
	}

	// Execute query
	rows, err := r.conn.NamedQuery(query, params)
	if err != nil {
		return models.User{}, false, err
	}
	defer rows.Close()

	return utils.ScanFirst(rows, r.Scan)
}

// Exists checks if a User with requested email exists in the repository
func (r *PostgresRepository) Exists(email string) (bool, error) {
	// Prepare query
	query := `SELECT *
			  FROM Users as u
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
func (r *PostgresRepository) Authenticate(email string, password string) (models.User, bool, error) {
	// Prepare query
	query := `SELECT *
			  FROM Users as u
			  WHERE u.email = :email`
	params := map[string]interface{}{
		"email": email,
	}

	// Execute query
	rows, err := r.conn.NamedQuery(query, params)
	if err != nil {
		return models.User{}, false, err
	}
	defer rows.Close()

	// Retrieve User
	var userWithPassword models.UserWithPassword
	if rows.Next() {
		userWithPassword, err = r.ScanWithPassword(rows)
		if err != nil {
			return models.User{}, false, err
		}

		err = bcrypt.CompareHashAndPassword([]byte(userWithPassword.Password), []byte(password))
		if err == nil {
			return userWithPassword.User, true, nil
		}
	}

	return models.User{}, false, errors.New("no User Found, invalid credentials")
}

// Update method used to update a User
func (r *PostgresRepository) Update(user models.User) error {

	// Prepare query
	query := `UPDATE Users as u
			  SET email = :email, updated_at = :updated_at
			  WHERE u.ID = :ID`
	params := map[string]interface{}{
		"ID":         user.ID,
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

// UpdateWithPassword method used to update a User with password
func (r *PostgresRepository) UpdateWithPassword(user models.UserWithPassword) error {

	// Hash password before saving
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Prepare query
	query := `UPDATE Users as u
			  SET password = :password, updated_at = :updated_at
			  WHERE u.ID = :ID`
	params := map[string]interface{}{
		"ID":         user.ID,
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

// Delete method used to delete a User
func (r *PostgresRepository) Delete(userID uuid.UUID) error {

	// Prepare query
	query := `DELETE FROM Users as u
			  WHERE u.ID = :ID`
	params := map[string]interface{}{
		"ID": userID,
	}

	// Execute query
	result, err := r.conn.NamedExec(query, params)
	if err != nil {
		return err
	}

	return utils.CheckRowAffected(result, 1)
}

func (r *PostgresRepository) Scan(rows *sqlx.Rows) (models.User, error) {

	userWithPassword, err := r.ScanWithPassword(rows)
	if err != nil {
		return models.User{}, err
	}

	return userWithPassword.User, nil
}

func (r *PostgresRepository) ScanWithPassword(rows *sqlx.Rows) (models.UserWithPassword, error) {
	var userWithPassword models.UserWithPassword
	err := rows.Scan(
		&userWithPassword.ID,
		&userWithPassword.Email,
		&userWithPassword.Password,
		&userWithPassword.CreatedAt,
		&userWithPassword.UpdatedAt,
	)
	if err != nil {
		return models.UserWithPassword{}, err
	}

	return userWithPassword, nil
}
