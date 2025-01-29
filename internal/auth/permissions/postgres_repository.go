package permissions

import (
	"github.com/Zapharaos/fihub-backend/internal/utils"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// PostgresRepository is a repository containing the user permissions data based on a PSQL database and
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

// Get search and returns a Permission from the repository by its id
func (r *PostgresRepository) Get(permissionUUID uuid.UUID) (Permission, bool, error) {
	// Prepare query
	query := `SELECT *
			  FROM permissions as p
			  WHERE p.id = :id`
	params := map[string]interface{}{
		"id": permissionUUID,
	}

	// Execute query
	rows, err := r.conn.NamedQuery(query, params)
	if err != nil {
		return Permission{}, false, err
	}
	defer rows.Close()

	return utils.ScanFirst(rows, r.Scan)
}

// Create creates a new Permission in the repository
func (r *PostgresRepository) Create(permission Permission) (uuid.UUID, error) {

	newUUID := uuid.New()

	// Prepare query
	query := `INSERT INTO permissions (id, value, scope, description)
				VALUES (:id, :value, :scope, :description)`
	params := map[string]interface{}{
		"id":          newUUID,
		"value":       permission.Value,
		"scope":       permission.Scope,
		"description": permission.Description,
	}

	// Execute query
	_, err := r.conn.NamedQuery(query, params)
	if err != nil {
		return uuid.UUID{}, err
	}

	return newUUID, nil
}

// Update updates a Permission in the repository
func (r *PostgresRepository) Update(permission Permission) error {
	// Prepare query
	query := `UPDATE permissions as p
			  SET value = :value, scope = :scope, description = :description
			  WHERE p.id = :id`
	params := map[string]interface{}{
		"id":          permission.Id,
		"value":       permission.Value,
		"scope":       permission.Scope,
		"description": permission.Description,
	}

	// Execute query
	result, err := r.conn.NamedExec(query, params)
	if err != nil {
		return err
	}

	return utils.CheckRowAffected(result, 1)
}

// Delete deletes a Permission in the repository
func (r *PostgresRepository) Delete(uuid uuid.UUID) error {
	// Prepare query
	query := `DELETE FROM permissions as p
			  WHERE p.id = :id`
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

// GetAll returns all User Permissions in the repository
func (r *PostgresRepository) GetAll() ([]Permission, error) {
	// Prepare query
	query := `SELECT *
			  FROM permissions`

	// Execute query
	rows, err := r.conn.Queryx(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return utils.ScanAll(rows, r.Scan)
}

// GetAllByRoleId returns all Permissions for a given Role
func (r *PostgresRepository) GetAllByRoleId(roleUUID uuid.UUID) ([]Permission, error) {
	// Prepare query
	query := `SELECT *
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

	return utils.ScanAll(rows, r.Scan)
}

// GetAllByRoleIds returns all Permissions for a given list of Roles
func (r *PostgresRepository) GetAllByRoleIds(roleUUID []uuid.UUID) ([]Permission, error) {
	// Prepare query
	query := `SELECT *
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

	return utils.ScanAll(rows, r.Scan)
}

// Scan scans the retrieved data from the database and returns a Permission
func (r *PostgresRepository) Scan(rows *sqlx.Rows) (Permission, error) {
	var permission Permission
	err := rows.Scan(
		&permission.Id,
		&permission.Value,
		&permission.Scope,
		&permission.Description,
	)
	if err != nil {
		return Permission{}, err
	}
	return permission, nil
}
