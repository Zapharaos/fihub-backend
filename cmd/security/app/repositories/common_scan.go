package repositories

import (
	"github.com/Zapharaos/fihub-backend/internal/models"
	"github.com/jmoiron/sqlx"
)

// ScanPermission scans the retrieved data from the database and returns a Permission
func ScanPermission(rows *sqlx.Rows) (models.Permission, error) {
	var permission models.Permission
	err := rows.Scan(
		&permission.Id,
		&permission.Value,
		&permission.Scope,
		&permission.Description,
	)
	if err != nil {
		return models.Permission{}, err
	}
	return permission, nil
}
