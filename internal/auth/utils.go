package auth

import (
	"fmt"
	"github.com/Zapharaos/fihub-backend/internal/auth/permissions"
	"github.com/Zapharaos/fihub-backend/internal/auth/roles"
	"github.com/Zapharaos/fihub-backend/internal/auth/users"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// LoadFullUser loads the roles from the database
func LoadFullUser(userId uuid.UUID) (*users.UserWithRoles, bool, error) {
	if users.R() == nil {
		zap.L().Fatal("users repository is not initialized")
		return nil, false, nil
	}
	user, ok, err := users.R().Get(userId)
	if err != nil {
		return nil, false, err
	}
	if !ok {
		return nil, false, fmt.Errorf("user not found")
	}

	fullRoles, err := LoadUserRoles(user.ID)
	if err != nil {
		return nil, false, err
	}

	return &users.UserWithRoles{
		User:  user,
		Roles: fullRoles,
	}, true, nil
}

// LoadUserRoles loads the roles from the database
func LoadUserRoles(userId uuid.UUID) ([]roles.RoleWithPermissions, error) {
	userRoles, err := roles.R().GetRolesByUserId(userId)
	if err != nil {
		zap.L().Error("GetUserRoles.GetRolesByUserId", zap.Error(err))
		return nil, err
	}

	userRolesWithPermissions := make([]roles.RoleWithPermissions, 0)

	for _, role := range userRoles {
		perms, err := permissions.R().GetAllByRoleId(role.Id)
		if err != nil {
			zap.L().Error("GetUserRoles.GetAllByRoleId", zap.Error(err))
			return nil, err
		}
		userRolesWithPermissions = append(userRolesWithPermissions, roles.RoleWithPermissions{
			Role:        role,
			Permissions: perms,
		})
	}

	return userRolesWithPermissions, nil
}
