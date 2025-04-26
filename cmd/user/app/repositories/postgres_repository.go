package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Zapharaos/fihub-backend/internal/models"
	"github.com/Zapharaos/fihub-backend/internal/utils"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"go.uber.org/zap"
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

// GetWithRoles returns a User with its roles in the repository
func (r *PostgresRepository) GetWithRoles(userID uuid.UUID) (models.UserWithRoles, error) {
	// Prepare query
	query := `SELECT u.ID, u.email, u.created_at, u.updated_at, r.ID, r.name
			  FROM Users as u
			  LEFT JOIN user_roles as ur on u.ID = ur.user_id
			  LEFT JOIN roles as r on ur.role_id = r.ID
			  WHERE u.ID = :ID`
	params := map[string]interface{}{
		"ID": userID,
	}

	// Execute query
	rows, err := r.conn.NamedQuery(query, params)
	if err != nil {
		return models.UserWithRoles{}, err
	}
	defer rows.Close()

	results, err := r.ScanMultiplesWithRoles(rows)
	if err != nil {
		return models.UserWithRoles{}, err
	}

	return results[0], nil
}

// GetAllWithRoles returns a User with its roles in the repository
func (r *PostgresRepository) GetAllWithRoles() ([]models.UserWithRoles, error) {
	// Prepare query
	query := `SELECT u.ID, u.email, u.created_at, u.updated_at, r.ID, r.name
			  FROM Users as u
			  LEFT JOIN user_roles as ur on u.ID = ur.user_id
			  LEFT JOIN roles as r on ur.role_id = r.ID`
	params := map[string]interface{}{}

	// Execute query
	rows, err := r.conn.NamedQuery(query, params)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.ScanMultiplesWithRoles(rows)
}

// GetUsersByRoleID returns all User for a role in the repository
func (r *PostgresRepository) GetUsersByRoleID(roleUUID uuid.UUID) ([]models.User, error) {
	// Prepare query
	query := `SELECT u.ID, u.email, u.password, u.created_at, u.updated_at
			  FROM Users as u
			  INNER JOIN user_roles as ur on u.ID = ur.user_id
			  WHERE ur.role_id = :ID`
	params := map[string]interface{}{
		"ID": roleUUID,
	}

	// Execute query
	rows, err := r.conn.NamedQuery(query, params)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return utils.ScanAll(rows, r.Scan)
}

// UpdateWithRoles updates a User with its roles in the repository
func (r *PostgresRepository) UpdateWithRoles(user models.UserWithRoles, roleUUIDs []uuid.UUID) error {
	// Start transaction
	ctx := context.Background()
	tx, err := r.conn.BeginTx(ctx, nil)
	if err != nil {
		zap.L().Error("Cannot start transaction", zap.Error(err))
		return err
	}

	// Prepare query
	query := `UPDATE Users as u SET updated_at = $1 WHERE u.ID = $2`
	_, err = tx.ExecContext(ctx, query, time.Now().Truncate(1*time.Millisecond).UTC(), user.ID)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return fmt.Errorf("main error: %v, rollback error: %v", err, rollbackErr)
		}
		return err
	}

	// If no roles are provided, we can commit and return
	if len(roleUUIDs) == 0 {
		if err = tx.Commit(); err != nil {
			return err
		}
		return nil
	}

	// Execute query to reset roles
	query = `DELETE FROM user_roles as ur WHERE ur.user_id = $1`
	_, err = tx.ExecContext(ctx, query, user.ID)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return fmt.Errorf("main error: %v, rollback error: %v", err, rollbackErr)
		}
		return err
	}

	// Prepare query to set new roles
	query = `INSERT INTO user_roles (user_id, role_id) VALUES `
	var values []interface{}
	for i, roleUUID := range roleUUIDs {
		query += fmt.Sprintf("($%d, $%d),", i*2+1, i*2+2)
		values = append(values, user.ID, roleUUID)
	}
	query = query[:len(query)-1] // Remove the trailing comma

	// Execute query
	result, err := tx.ExecContext(ctx, query, values...)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return fmt.Errorf("main error: %v, rollback error: %v", err, rollbackErr)
		}
		return err
	}

	// Check if all roles were set
	if err = utils.CheckRowAffected(result, int64(len(roleUUIDs))); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return fmt.Errorf("main error: %v, rollback error: %v", err, rollbackErr)
		}
		return err
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

// SetUserRoles sets the roles of a User in the repository
func (r *PostgresRepository) SetUserRoles(userUUID uuid.UUID, roleUUIDs []uuid.UUID) error {

	// Start transaction
	ctx := context.Background()
	tx, err := r.conn.BeginTx(ctx, nil)
	if err != nil {
		zap.L().Error("Cannot start transaction", zap.Error(err))
		return err
	}

	// Execute query to reset roles
	query := `DELETE FROM user_roles as ur WHERE ur.user_id = $1`
	_, err = tx.ExecContext(ctx, query, userUUID)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return fmt.Errorf("main error: %v, rollback error: %v", err, rollbackErr)
		}
		return err
	}

	// If no roles are provided, we can commit and return
	if len(roleUUIDs) == 0 {
		if err = tx.Commit(); err != nil {
			return err
		}
		return nil
	}

	// Prepare query to set new roles
	query = `INSERT INTO user_roles (user_id, role_id) VALUES `
	var values []interface{}
	for i, roleUUID := range roleUUIDs {
		query += fmt.Sprintf("($%d, $%d),", i*2+1, i*2+2)
		values = append(values, userUUID, roleUUID)
	}
	query = query[:len(query)-1] // Remove the trailing comma

	// Execute query
	result, err := tx.ExecContext(ctx, query, values...)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return fmt.Errorf("main error: %v, rollback error: %v", err, rollbackErr)
		}
		return err
	}

	// Check if all roles were set
	if err = utils.CheckRowAffected(result, int64(len(roleUUIDs))); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return fmt.Errorf("main error: %v, rollback error: %v", err, rollbackErr)
		}
		return err
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

// AddUsersRole adds a role to a list of User in the repository
func (r *PostgresRepository) AddUsersRole(userUUIDs []uuid.UUID, roleUUID uuid.UUID) error {

	// Start transaction
	ctx := context.Background()
	tx, err := r.conn.BeginTx(ctx, nil)
	if err != nil {
		zap.L().Error("Cannot start transaction", zap.Error(err))
		return err
	}

	// Prepare query to add new role to Users
	query := `INSERT INTO user_roles (user_id, role_id) VALUES `
	var values []interface{}
	for i, userUUID := range userUUIDs {
		query += fmt.Sprintf("($%d, $%d),", i*2+1, i*2+2)
		values = append(values, userUUID, roleUUID)
	}
	query = query[:len(query)-1] // Remove the trailing comma

	// Execute query
	result, err := tx.ExecContext(ctx, query, values...)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return fmt.Errorf("main error: %v, rollback error: %v", err, rollbackErr)
		}
		return err
	}

	// Check if all roles were set
	if err = utils.CheckRowAffected(result, int64(len(userUUIDs))); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return fmt.Errorf("main error: %v, rollback error: %v", err, rollbackErr)
		}
		return err
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

// RemoveUsersRole removes a role from a list of User in the repository
func (r *PostgresRepository) RemoveUsersRole(userUUIDs []uuid.UUID, roleUUID uuid.UUID) error {
	// Start transaction
	ctx := context.Background()
	tx, err := r.conn.BeginTx(ctx, nil)
	if err != nil {
		zap.L().Error("Cannot start transaction", zap.Error(err))
		return err
	}

	// Query to remove role from Users
	query := `DELETE FROM user_roles WHERE user_id = ANY(?) AND role_id = ?`
	result, err := tx.ExecContext(ctx, query, pq.Array(userUUIDs), roleUUID)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return fmt.Errorf("main error: %v, rollback error: %v", err, rollbackErr)
		}
		return err
	}

	// Check if all roles were set
	if err = utils.CheckRowAffected(result, int64(len(userUUIDs))); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return fmt.Errorf("main error: %v, rollback error: %v", err, rollbackErr)
		}
		return err
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
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

func (r *PostgresRepository) ScanWithRoles(rows *sqlx.Rows) (models.UserWithRoles, error) {

	var userWithRoles models.UserWithRoles
	var role models.RoleWithPermissions
	var roleName sql.NullString
	err := rows.Scan(
		&userWithRoles.ID,
		&userWithRoles.Email,
		&userWithRoles.CreatedAt,
		&userWithRoles.UpdatedAt,
		&role.Id,
		&roleName,
	)

	// Check if there is an error
	if err != nil {
		return models.UserWithRoles{}, err
	}

	// If role exists, add it to the User
	if role.Id != uuid.Nil {
		role.Name = roleName.String
		userWithRoles.Roles = append(userWithRoles.Roles, role)
	}

	return userWithRoles, nil
}

func (r *PostgresRepository) ScanMultiplesWithRoles(rows *sqlx.Rows) ([]models.UserWithRoles, error) {
	var usersWithRoles []models.UserWithRoles
	userMap := make(map[uuid.UUID]int)

	for rows.Next() {
		// One row is a User with one single role
		userWithRoles, err := r.ScanWithRoles(rows)
		if err != nil {
			return []models.UserWithRoles{}, err
		}

		// Retrieve User from map if exists
		index, exists := userMap[userWithRoles.ID]

		// If User does not exist, add it to the map and the list
		if !exists {
			userMap[userWithRoles.ID] = len(usersWithRoles)
			usersWithRoles = append(usersWithRoles, userWithRoles)
			continue
		}

		// If User already exists but has no bew roles, skip
		if len(userWithRoles.Roles) == 0 {
			continue
		}

		// Add the role to the User
		usersWithRoles[index].Roles = append(usersWithRoles[index].Roles, userWithRoles.Roles[0])
	}

	return usersWithRoles, nil
}
