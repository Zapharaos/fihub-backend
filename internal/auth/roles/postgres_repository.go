package roles

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Zapharaos/fihub-backend/internal/auth/permissions"
	"github.com/Zapharaos/fihub-backend/internal/utils"
	"go.uber.org/zap"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// PostgresRepository is a repository containing the user roles data based on a PSQL database and
// implementing the repository interface
type PostgresRepository struct {
	conn *sqlx.DB
}

// NewPostgresRepository returns a new instance of PostgresRepository
func NewPostgresRepository(dbClient *sqlx.DB) Repository {
	r := PostgresRepository{
		conn: dbClient,
	}
	var ifm Repository = &r
	return ifm
}

// Get search and returns a Role from the repository by its id
func (r *PostgresRepository) Get(roleUUID uuid.UUID) (Role, bool, error) {
	// Prepare query
	query := `SELECT *
			  FROM roles as r
			  WHERE r.id = :id`
	params := map[string]interface{}{
		"id": roleUUID,
	}

	// Execute query
	rows, err := r.conn.NamedQuery(query, params)
	if err != nil {
		return Role{}, false, err
	}
	defer rows.Close()

	return utils.ScanFirst(rows, r.Scan)
}

// GetByName search and returns a Role from the repository by its name
func (r *PostgresRepository) GetByName(name string) (Role, bool, error) {
	// Prepare query
	query := `SELECT *
			  FROM roles as r
			  WHERE r.name = :name`
	params := map[string]interface{}{
		"name": name,
	}

	// Execute query
	rows, err := r.conn.NamedQuery(query, params)
	if err != nil {
		return Role{}, false, err
	}
	defer rows.Close()

	return utils.ScanFirst(rows, r.Scan)
}

// Create creates a new Role in the repository
func (r *PostgresRepository) Create(role Role, permissionUUIDs []uuid.UUID) (uuid.UUID, error) {

	roleUUID := uuid.New()

	// Start transaction
	ctx := context.Background()
	tx, err := r.conn.BeginTx(ctx, nil)
	if err != nil {
		zap.L().Error("Cannot start transaction", zap.Error(err))
		return uuid.Nil, err
	}

	// Query to create the role
	query := `INSERT INTO roles (id, name) VALUES ($1, $2)`
	result, err := tx.ExecContext(ctx, query, roleUUID, role.Name)
	if err != nil {
		tx.Rollback()
		return uuid.Nil, err
	}

	// If no permissions are provided, we can commit and return
	if len(permissionUUIDs) == 0 {
		tx.Commit()
		return roleUUID, nil
	}

	// Prepare query to set new permissions
	query = `INSERT INTO role_permissions (role_id, permission_id) VALUES `
	var values []interface{}
	for i, permissionUUID := range permissionUUIDs {
		query += fmt.Sprintf("($%d, $%d),", i*2+1, i*2+2)
		values = append(values, roleUUID, permissionUUID)
	}
	query = query[:len(query)-1] // Remove the trailing comma

	// Execute query
	result, err = tx.ExecContext(ctx, query, values...)
	if err != nil {
		tx.Rollback()
		return uuid.Nil, err
	}

	// Check if all permissions were set
	if err = utils.CheckRowAffected(result, int64(len(permissionUUIDs))); err != nil {
		tx.Rollback()
		return uuid.Nil, err
	}

	tx.Commit()
	return roleUUID, nil
}

// Update updates a Role in the repository
func (r *PostgresRepository) Update(role Role, permissionUUIDs []uuid.UUID) error {
	// Start transaction
	ctx := context.Background()
	tx, err := r.conn.BeginTx(ctx, nil)
	if err != nil {
		zap.L().Error("Cannot start transaction", zap.Error(err))
		return err
	}

	// Query to update the role
	query := `UPDATE roles as r SET name = $1 WHERE r.id = $2`
	result, err := tx.ExecContext(ctx, query, role.Name, role.Id)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Query to reset permissions
	query = `DELETE FROM role_permissions WHERE role_id = $1`
	result, err = tx.ExecContext(ctx, query, role.Id)
	if err != nil {
		tx.Rollback()
		return err
	}

	// If no permissions are provided, we can commit and return
	if len(permissionUUIDs) == 0 {
		tx.Commit()
		return nil
	}

	// Prepare query to set new permissions
	query = `INSERT INTO role_permissions (role_id, permission_id) VALUES `
	var values []interface{}
	for i, permissionUUID := range permissionUUIDs {
		query += fmt.Sprintf("($%d, $%d),", i*2+1, i*2+2)
		values = append(values, role.Id, permissionUUID)
	}
	query = query[:len(query)-1] // Remove the trailing comma

	// Execute query
	result, err = tx.ExecContext(ctx, query, values...)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Check if all permissions were set
	if err = utils.CheckRowAffected(result, int64(len(permissionUUIDs))); err != nil {
		tx.Rollback()
		return err
	}

	// Commit transaction
	tx.Commit()
	return nil
}

// Delete deletes a Role in the repository
func (r *PostgresRepository) Delete(uuid uuid.UUID) error {
	// Prepare query
	query := `DELETE FROM roles as r
			  WHERE r.id = :id`
	params := map[string]interface{}{
		"id": uuid,
	}

	// Execute query
	result, err := r.conn.NamedExec(query, params)
	if err != nil {
		return err
	}

	return utils.CheckRowAffected(result, 1)
}

// GetAll returns all Roles in the repository
func (r *PostgresRepository) GetAll() ([]Role, error) {
	// Prepare query
	query := `SELECT *
			  FROM roles`

	// Execute query
	rows, err := r.conn.Queryx(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return utils.ScanAll(rows, r.Scan)
}

// GetWithPermissions returns a Role in the repository with its permissions
func (r *PostgresRepository) GetWithPermissions(uuid uuid.UUID) (RoleWithPermissions, bool, error) {
	// Prepare query
	query := `SELECT r.id, r.name, p.id, p.value, p.scope, p.description
			  FROM roles as r
			  LEFT JOIN role_permissions as rp on r.id = rp.role_id
			  LEFT JOIN permissions as p on rp.permission_id = p.id
			  WHERE r.id = :id`
	params := map[string]interface{}{
		"id": uuid,
	}

	// Execute query
	rows, err := r.conn.NamedQuery(query, params)
	if err != nil {
		return RoleWithPermissions{}, false, err
	}
	defer rows.Close()

	return utils.ScanFirst(rows, r.ScanWithPermissions)
}

// GetAllWithPermissions returns all Roles in the repository with their permissions
func (r *PostgresRepository) GetAllWithPermissions() (RolesWithPermissions, error) {
	// Prepare query
	query := `SELECT r.id, r.name, p.id, p.value, p.scope, p.description
			  FROM roles as r
			  LEFT JOIN role_permissions as rp on r.id = rp.role_id
			  LEFT JOIN permissions as p on rp.permission_id = p.id`
	params := map[string]interface{}{}

	// Execute query
	rows, err := r.conn.NamedQuery(query, params)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.ScanAllWithPermissions(rows)
}

// GetRolesByUserId returns all the roles of a user in the repository
func (r *PostgresRepository) GetRolesByUserId(userUUID uuid.UUID) ([]Role, error) {
	// Prepare query
	query := `SELECT r.id, r.name
			  FROM roles as r
			  INNER JOIN user_roles as ur on r.id = ur.role_id
			  WHERE ur.user_id = :id`
	params := map[string]interface{}{
		"id": userUUID,
	}

	// Execute query
	rows, err := r.conn.NamedQuery(query, params)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return utils.ScanAll(rows, r.Scan)
}

func (r *PostgresRepository) SetRolePermissions(roleUUID uuid.UUID, permissionUUIDs []uuid.UUID) error {

	// Start transaction
	ctx := context.Background()
	tx, err := r.conn.BeginTx(ctx, nil)
	if err != nil {
		zap.L().Error("Cannot start transaction", zap.Error(err))
		return err
	}

	// Query to reset permissions
	query := `DELETE FROM role_permissions WHERE role_id = $1`
	result, err := tx.ExecContext(ctx, query, roleUUID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// If no permissions are provided, we can commit and return
	if len(permissionUUIDs) == 0 {
		tx.Commit()
		return nil
	}

	// Prepare query to set new permissions
	query = `INSERT INTO role_permissions (role_id, permission_id) VALUES `
	var values []interface{}
	for i, permissionUUID := range permissionUUIDs {
		query += fmt.Sprintf("($%d, $%d),", i*2+1, i*2+2)
		values = append(values, roleUUID, permissionUUID)
	}
	query = query[:len(query)-1] // Remove the trailing comma

	// Execute query
	result, err = tx.ExecContext(ctx, query, values...)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Check if all permissions were set
	if err = utils.CheckRowAffected(result, int64(len(permissionUUIDs))); err != nil {
		tx.Rollback()
		return err
	}

	// Commit transaction
	tx.Commit()
	return nil
}

// Scan scans the current row of the given rows and returns a Role
func (r *PostgresRepository) Scan(rows *sqlx.Rows) (Role, error) {
	var role Role
	err := rows.Scan(
		&role.Id,
		&role.Name,
	)
	if err != nil {
		return Role{}, err
	}
	return role, nil
}

// ScanWithPermissions scans the current row of the given rows and returns a RoleWithPermissions
func (r *PostgresRepository) ScanWithPermissions(rows *sqlx.Rows) (RoleWithPermissions, error) {

	var role RoleWithPermissions
	var permission permissions.Permission
	var pValue sql.NullString
	var pScope sql.NullString
	var pDescription sql.NullString

	err := rows.Scan(
		&role.Id,
		&role.Name,
		&permission.Id,
		&pValue,
		&pScope,
		&pDescription,
	)

	// Check if there is an error
	if err != nil {
		return RoleWithPermissions{}, err
	}

	// If role exists, add it to the user
	if permission.Id != uuid.Nil {
		permission.Value = pValue.String
		permission.Scope = pScope.String
		permission.Description = pDescription.String
		role.Permissions = append(role.Permissions, permission)
	}

	return role, nil
}

// ScanAllWithPermissions scans all rows of the given rows and returns a list of RoleWithPermissions
func (r *PostgresRepository) ScanAllWithPermissions(rows *sqlx.Rows) (RolesWithPermissions, error) {
	var roles RolesWithPermissions
	rolesMap := make(map[uuid.UUID]int)

	for rows.Next() {
		// One row is a role with one single permission
		role, err := r.ScanWithPermissions(rows)
		if err != nil {
			return RolesWithPermissions{}, err
		}

		// Retrieve role from map if exists
		id := role.Id
		index, exists := rolesMap[id]

		// If role does not exist, add it to the map and the list
		if !exists {
			rolesMap[id] = len(roles)
			roles = append(roles, role)
			continue
		}

		// If role already exists but has no new permission, skip
		if len(role.Permissions) == 0 {
			continue
		}

		// Add the permission to the role
		roles[index].Permissions = append(roles[index].Permissions, role.Permissions[0])
	}

	return roles, nil
}
