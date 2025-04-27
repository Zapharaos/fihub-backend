package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Zapharaos/fihub-backend/internal/models"
	"github.com/Zapharaos/fihub-backend/internal/utils"
	"github.com/lib/pq"
	"go.uber.org/zap"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// RolePostgresRepository is a repository containing the user roles data based on a PSQL database and
// implementing the repository interface
type RolePostgresRepository struct {
	conn *sqlx.DB
}

// NewRolePostgresRepository returns a new instance of RolePostgresRepository
func NewRolePostgresRepository(dbClient *sqlx.DB) RoleRepository {
	r := RolePostgresRepository{
		conn: dbClient,
	}
	var rr RoleRepository = &r
	return rr
}

// Create creates a new Role in the repository
func (r *RolePostgresRepository) Create(role models.Role, permissionUUIDs []uuid.UUID) (uuid.UUID, error) {

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
	_, err = tx.ExecContext(ctx, query, roleUUID, role.Name)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return uuid.Nil, fmt.Errorf("main error: %v, rollback error: %v", err, rollbackErr)
		}
		return uuid.Nil, err
	}

	// If no permissions are provided, we can commit and return
	if len(permissionUUIDs) == 0 {
		if err = tx.Commit(); err != nil {
			return uuid.Nil, err
		}
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
	result, err := tx.ExecContext(ctx, query, values...)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return uuid.Nil, fmt.Errorf("main error: %v, rollback error: %v", err, rollbackErr)
		}
		return uuid.Nil, err
	}

	// Check if all permissions were set
	if err = utils.CheckRowAffected(result, int64(len(permissionUUIDs))); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return uuid.Nil, fmt.Errorf("main error: %v, rollback error: %v", err, rollbackErr)
		}
		return uuid.Nil, err
	}

	if err = tx.Commit(); err != nil {
		return uuid.Nil, err
	}
	return roleUUID, nil
}

// Get search and returns a Role from the repository by its id
func (r *RolePostgresRepository) Get(roleUUID uuid.UUID) (models.Role, bool, error) {
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
		return models.Role{}, false, err
	}
	defer rows.Close()

	return utils.ScanFirst(rows, r.Scan)
}

// GetByName search and returns a Role from the repository by its name
func (r *RolePostgresRepository) GetByName(name string) (models.Role, bool, error) {
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
		return models.Role{}, false, err
	}
	defer rows.Close()

	return utils.ScanFirst(rows, r.Scan)
}

// GetWithPermissions returns a Role in the repository with its permissions
func (r *RolePostgresRepository) GetWithPermissions(uuid uuid.UUID) (models.RoleWithPermissions, bool, error) {
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
		return models.RoleWithPermissions{}, false, err
	}
	defer rows.Close()

	return utils.ScanFirst(rows, r.ScanWithPermissions)
}

// Update updates a Role in the repository
func (r *RolePostgresRepository) Update(role models.Role, permissionUUIDs []uuid.UUID) error {
	// Start transaction
	ctx := context.Background()
	tx, err := r.conn.BeginTx(ctx, nil)
	if err != nil {
		zap.L().Error("Cannot start transaction", zap.Error(err))
		return err
	}

	// Query to update the role
	query := `UPDATE roles as r SET name = $1 WHERE r.id = $2`
	_, err = tx.ExecContext(ctx, query, role.Name, role.Id)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return fmt.Errorf("main error: %v, rollback error: %v", err, rollbackErr)
		}
		return err
	}

	// Query to reset permissions
	query = `DELETE FROM role_permissions WHERE role_id = $1`
	_, err = tx.ExecContext(ctx, query, role.Id)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return fmt.Errorf("main error: %v, rollback error: %v", err, rollbackErr)
		}
		return err
	}

	// If no permissions are provided, we can commit and return
	if len(permissionUUIDs) == 0 {
		if err = tx.Commit(); err != nil {
			return err
		}
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
	result, err := tx.ExecContext(ctx, query, values...)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return fmt.Errorf("main error: %v, rollback error: %v", err, rollbackErr)
		}
		return err
	}

	// Check if all permissions were set
	if err = utils.CheckRowAffected(result, int64(len(permissionUUIDs))); err != nil {
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

// Delete deletes a Role in the repository
func (r *RolePostgresRepository) Delete(uuid uuid.UUID) error {
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

// List returns all Roles in the repository
func (r *RolePostgresRepository) List() ([]models.Role, error) {
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

// ListByUserId returns all the roles of a user in the repository
func (r *RolePostgresRepository) ListByUserId(userUUID uuid.UUID) ([]models.Role, error) {
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

// ListWithPermissions returns all Roles in the repository with their permissions
func (r *RolePostgresRepository) ListWithPermissions() (models.RolesWithPermissions, error) {
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

// SetForUser sets the roles of a User in the repository
func (r *RolePostgresRepository) SetForUser(userUUID uuid.UUID, roleUUIDs []uuid.UUID) error {

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

// AddToUsers adds a role to a list of User in the repository
func (r *RolePostgresRepository) AddToUsers(userUUIDs []uuid.UUID, roleUUID uuid.UUID) error {

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

// RemoveFromUsers removes a role from a list of User in the repository
func (r *RolePostgresRepository) RemoveFromUsers(userUUIDs []uuid.UUID, roleUUID uuid.UUID) error {
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

// SetPermissionsByRoleId sets the permissions of a Role in the repository
func (r *RolePostgresRepository) SetPermissionsByRoleId(roleUUID uuid.UUID, permissionUUIDs []uuid.UUID) error {

	// Start transaction
	ctx := context.Background()
	tx, err := r.conn.BeginTx(ctx, nil)
	if err != nil {
		zap.L().Error("Cannot start transaction", zap.Error(err))
		return err
	}

	// Query to reset permissions
	query := `DELETE FROM role_permissions WHERE role_id = $1`
	_, err = tx.ExecContext(ctx, query, roleUUID)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return fmt.Errorf("main error: %v, rollback error: %v", err, rollbackErr)
		}
		return err
	}

	// If no permissions are provided, we can commit and return
	if len(permissionUUIDs) == 0 {
		if err = tx.Commit(); err != nil {
			return err
		}
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
	result, err := tx.ExecContext(ctx, query, values...)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return fmt.Errorf("main error: %v, rollback error: %v", err, rollbackErr)
		}
		return err
	}

	// Check if all permissions were set
	if err = utils.CheckRowAffected(result, int64(len(permissionUUIDs))); err != nil {
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

// ListPermissionsByRoleId returns all Permissions for a given Role
func (r *RolePostgresRepository) ListPermissionsByRoleId(roleUUID uuid.UUID) (models.Permissions, error) {
	// Prepare query
	query := `SELECT p.id, p.value, p.scope, p.description
			  FROM permissions as p
			  INNER JOIN role_permissions as rp on p.id = rp.permission_id
			  WHERE rp.role_id = :id`
	params := map[string]interface{}{
		"id": roleUUID,
	}

	// Execute query
	rows, err := r.conn.NamedQuery(query, params)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return utils.ScanAll(rows, ScanPermission)
}

// ListPermissionsByUserId returns all Permissions for a given User
func (r *RolePostgresRepository) ListPermissionsByUserId(userUUID uuid.UUID) (models.Permissions, error) {
	// Prepare query
	query := `SELECT p.id, p.value, p.scope, p.description
			  FROM permissions as p
			  INNER JOIN role_permissions as rp on p.id = rp.permission_id
			  INNER JOIN user_roles as ur on rp.role_id = ur.role_id
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

	return utils.ScanAll(rows, ScanPermission)
}

// Scan scans the current row of the given rows and returns a Role
func (r *RolePostgresRepository) Scan(rows *sqlx.Rows) (models.Role, error) {
	var role models.Role
	err := rows.Scan(
		&role.Id,
		&role.Name,
	)
	if err != nil {
		return models.Role{}, err
	}
	return role, nil
}

// ScanWithPermissions scans the current row of the given rows and returns a RoleWithPermissions
func (r *RolePostgresRepository) ScanWithPermissions(rows *sqlx.Rows) (models.RoleWithPermissions, error) {

	var role models.RoleWithPermissions
	var permission models.Permission
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
		return models.RoleWithPermissions{}, err
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
func (r *RolePostgresRepository) ScanAllWithPermissions(rows *sqlx.Rows) (models.RolesWithPermissions, error) {
	var roles models.RolesWithPermissions
	rolesMap := make(map[uuid.UUID]int)

	for rows.Next() {
		// One row is a role with one single permission
		role, err := r.ScanWithPermissions(rows)
		if err != nil {
			return models.RolesWithPermissions{}, err
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
