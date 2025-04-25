package service

import (
	"errors"
	"github.com/Zapharaos/fihub-backend/cmd/user/app/repositories"
	"github.com/Zapharaos/fihub-backend/internal/models"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

var (
	ErrorUserNotFound = errors.New("user not found")
)

// LoadFullUser loads the roles from the database
func LoadFullUser(userId uuid.UUID) (*models.UserWithRoles, bool, error) {
	user, ok, err := repositories.R().U().Get(userId)
	if err != nil {
		return nil, false, err
	}
	if !ok {
		return nil, false, ErrorUserNotFound
	}

	fullRoles, err := LoadUserRoles(user.ID)
	if err != nil {
		return nil, false, err
	}

	return &models.UserWithRoles{
		User:  user,
		Roles: fullRoles,
	}, true, nil
}

// LoadUserRoles loads the roles from the database
func LoadUserRoles(userId uuid.UUID) (models.RolesWithPermissions, error) {
	userRoles, err := repositories.R().R().GetRolesByUserId(userId)
	if err != nil {
		zap.L().Error("GetUserRoles.GetRolesByUserId", zap.Error(err))
		return nil, err
	}

	userRolesWithPermissions := make(models.RolesWithPermissions, 0)

	for _, role := range userRoles {
		perms, err := repositories.R().P().GetAllByRoleId(role.Id)
		if err != nil {
			zap.L().Error("GetUserRoles.GetAllByRoleId", zap.Error(err))
			return nil, err
		}
		userRolesWithPermissions = append(userRolesWithPermissions, models.RoleWithPermissions{
			Role:        role,
			Permissions: perms,
		})
	}

	return userRolesWithPermissions, nil
}
