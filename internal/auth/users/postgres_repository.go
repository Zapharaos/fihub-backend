package users

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Zapharaos/fihub-backend/internal/auth/roles"
	"github.com/Zapharaos/fihub-backend/internal/utils"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"go.uber.org/zap"
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

// GetWithRoles returns a User with its roles in the repository
func (r *PostgresRepository) GetWithRoles(userID uuid.UUID) (UserWithRoles, error) {
	// Prepare query
	query := `SELECT u.id, u.email, u.created_at, u.updated_at, r.id, r.name
			  FROM users as u
			  LEFT JOIN user_roles as ur on u.id = ur.user_id
			  LEFT JOIN roles as r on ur.role_id = r.id
			  WHERE u.id = :id`
	params := map[string]interface{}{
		"id": userID,
	}

	// Execute query
	rows, err := r.conn.NamedQuery(query, params)
	if err != nil {
		return UserWithRoles{}, err
	}
	defer rows.Close()

	results, err := scanUsersWithRoles(rows)
	if err != nil {
		return UserWithRoles{}, err
	}

	return results[0], nil
}

// GetAllWithRoles returns a User with its roles in the repository
func (r *PostgresRepository) GetAllWithRoles() ([]UserWithRoles, error) {
	// Prepare query
	query := `SELECT u.id, u.email, u.created_at, u.updated_at, r.id, r.name
			  FROM users as u
			  LEFT JOIN user_roles as ur on u.id = ur.user_id
			  LEFT JOIN roles as r on ur.role_id = r.id`
	params := map[string]interface{}{}

	// Execute query
	rows, err := r.conn.NamedQuery(query, params)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanUsersWithRoles(rows)
}

// GetUsersByRoleID returns all Users for a role in the repository
func (r *PostgresRepository) GetUsersByRoleID(roleUUID uuid.UUID) ([]User, error) {
	// Prepare query
	query := `SELECT u.id, u.email, u.password, u.created_at, u.updated_at
			  FROM users as u
			  INNER JOIN user_roles as ur on u.id = ur.user_id
			  WHERE ur.role_id = :id`
	params := map[string]interface{}{
		"id": roleUUID,
	}

	// Execute query
	rows, err := r.conn.NamedQuery(query, params)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return utils.ScanAll(rows, scanUser)
}

// UpdateWithRoles updates a user with its roles in the repository
func (r *PostgresRepository) UpdateWithRoles(user UserWithRoles, roleUUIDs []uuid.UUID) error {
	// Start transaction
	ctx := context.Background()
	tx, err := r.conn.BeginTx(ctx, nil)
	if err != nil {
		zap.L().Error("Cannot start transaction", zap.Error(err))
		return err
	}

	// Prepare query
	query := `UPDATE users as u SET updated_at = $1 WHERE u.id = $2`
	result, err := tx.ExecContext(ctx, query, time.Now().Truncate(1*time.Millisecond).UTC(), user.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// If no roles are provided, we can commit and return
	if len(roleUUIDs) == 0 {
		tx.Commit()
		return nil
	}

	// Execute query to reset roles
	query = `DELETE FROM user_roles as ur WHERE ur.user_id = $1`
	result, err = tx.ExecContext(ctx, query, user.ID)
	if err != nil {
		tx.Rollback()
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
	result, err = tx.ExecContext(ctx, query, values...)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Check if all roles were set
	if err = utils.CheckRowAffected(result, int64(len(roleUUIDs))); err != nil {
		tx.Rollback()
		return err
	}

	// Commit transaction
	tx.Commit()
	return nil
}

// SetUserRoles sets the roles of a user in the repository
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
	result, err := tx.ExecContext(ctx, query, userUUID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// If no roles are provided, we can commit and return
	if len(roleUUIDs) == 0 {
		tx.Commit()
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
	result, err = tx.ExecContext(ctx, query, values...)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Check if all roles were set
	if err = utils.CheckRowAffected(result, int64(len(roleUUIDs))); err != nil {
		tx.Rollback()
		return err
	}

	// Commit transaction
	tx.Commit()
	return nil
}

// AddUsersRole adds a role to a list of users in the repository
func (r *PostgresRepository) AddUsersRole(userUUIDs []uuid.UUID, roleUUID uuid.UUID) error {

	// Start transaction
	ctx := context.Background()
	tx, err := r.conn.BeginTx(ctx, nil)
	if err != nil {
		zap.L().Error("Cannot start transaction", zap.Error(err))
		return err
	}

	// Prepare query to add new role to users
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
		tx.Rollback()
		return err
	}

	// Check if all roles were set
	if err = utils.CheckRowAffected(result, int64(len(userUUIDs))); err != nil {
		tx.Rollback()
		return err
	}

	// Commit transaction
	tx.Commit()
	return nil
}

// RemoveUsersRole removes a role from a list of users in the repository
func (r *PostgresRepository) RemoveUsersRole(userUUIDs []uuid.UUID, roleUUID uuid.UUID) error {
	// Start transaction
	ctx := context.Background()
	tx, err := r.conn.BeginTx(ctx, nil)
	if err != nil {
		zap.L().Error("Cannot start transaction", zap.Error(err))
		return err
	}

	// Query to remove role from users
	query := `DELETE FROM user_roles WHERE user_id = ANY(?) AND role_id = ?`
	result, err := tx.ExecContext(ctx, query, pq.Array(userUUIDs), roleUUID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Check if all roles were set
	if err = utils.CheckRowAffected(result, int64(len(userUUIDs))); err != nil {
		tx.Rollback()
		return err
	}

	// Commit transaction
	tx.Commit()
	return nil
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

func scanUserWithRoles(rows *sqlx.Rows) (UserWithRoles, error) {

	var userWithRoles UserWithRoles
	var role roles.RoleWithPermissions
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
		return UserWithRoles{}, err
	}

	// If role exists, add it to the user
	if role.Id != uuid.Nil {
		role.Name = roleName.String
		userWithRoles.Roles = append(userWithRoles.Roles, role)
	}

	return userWithRoles, nil
}

func scanUsersWithRoles(rows *sqlx.Rows) ([]UserWithRoles, error) {
	var usersWithRoles []UserWithRoles
	userMap := make(map[uuid.UUID]int)

	for rows.Next() {
		// One row is a user with one single role
		userWithRoles, err := scanUserWithRoles(rows)
		if err != nil {
			return []UserWithRoles{}, err
		}

		// Retrieve user from map if exists
		index, exists := userMap[userWithRoles.ID]

		// If user does not exist, add it to the map and the list
		if !exists {
			userMap[userWithRoles.ID] = len(usersWithRoles)
			usersWithRoles = append(usersWithRoles, userWithRoles)
			continue
		}

		// If user already exists but has no bew roles, skip
		if len(userWithRoles.Roles) == 0 {
			continue
		}

		// Add the role to the user
		usersWithRoles[index].Roles = append(usersWithRoles[index].Roles, userWithRoles.Roles[0])
	}

	return usersWithRoles, nil
}
