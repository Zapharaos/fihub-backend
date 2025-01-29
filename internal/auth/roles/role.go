package roles

import (
	"fmt"
	"github.com/Zapharaos/fihub-backend/internal/auth/permissions"
	"github.com/google/uuid"
)

type Role struct {
	Id   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

// IsValid checks if a role is valid and has no missing mandatory fields
func (r Role) IsValid() (bool, error) {
	if r.Name == "" {
		return false, fmt.Errorf("missing name")
	}
	if len(r.Name) < 3 {
		return false, fmt.Errorf("name is too short (less than 3 characters)")
	}
	return true, nil
}

type RoleWithPermissions struct {
	Role
	Permissions []permissions.Permission `json:"permissions"`
}
