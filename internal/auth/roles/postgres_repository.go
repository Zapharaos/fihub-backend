package roles

import (
	"github.com/Zapharaos/fihub-backend/internal/utils"

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
func (r *PostgresRepository) Create(role Role) (uuid.UUID, error) {

	newUUID := uuid.New()

	// Prepare query
	query := `INSERT INTO roles (id, name)
				VALUES (:id, :name)`
	params := map[string]interface{}{
		"id":   newUUID,
		"name": role.Name,
	}

	// Execute query
	_, err := r.conn.NamedQuery(query, params)
	if err != nil {
		return uuid.UUID{}, err
	}

	return newUUID, nil

}

// Update updates a Role in the repository
func (r *PostgresRepository) Update(role Role) error {
	// Prepare query
	query := `UPDATE roles as r
			  SET name = :name,
			  WHERE p.id = :id`
	params := map[string]interface{}{
		"id":   role.Id,
		"name": role.Name,
	}

	// Execute query
	result, err := r.conn.NamedExec(query, params)
	if err != nil {
		return err
	}

	return utils.CheckRowAffected(result, 1)
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

// GetRolesByUserId returns all the roles of a user in the repository
func (r *PostgresRepository) GetRolesByUserId(userUUID uuid.UUID) ([]Role, error) {
	// Prepare query
	query := `SELECT *
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
