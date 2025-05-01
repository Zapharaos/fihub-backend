package mappers

import (
	"github.com/Zapharaos/fihub-backend/gen/go/securitypb"
	"github.com/Zapharaos/fihub-backend/internal/models"
	"github.com/google/uuid"
)

// RoleToProto converts a models.Role to a securitypb.Role
func RoleToProto(role models.Role) *securitypb.Role {
	return &securitypb.Role{
		Id:   role.Id.String(),
		Name: role.Name,
	}
}

// RoleFromProto converts a securitypb.Role to a models.Role
func RoleFromProto(role *securitypb.Role) models.Role {
	return models.Role{
		Id:   uuid.MustParse(role.GetId()),
		Name: role.GetName(),
	}
}

// RolesToProto converts a models.Roles to a slice of securitypb.Role
func RolesToProto(roles models.Roles) []*securitypb.Role {
	protoRoles := make([]*securitypb.Role, len(roles))
	for i, role := range roles {
		protoRoles[i] = RoleToProto(role)
	}
	return protoRoles
}

// RolesFromProto converts a slice of securitypb.Role to a models.Roles
func RolesFromProto(roles []*securitypb.Role) models.Roles {
	protoRoles := make(models.Roles, len(roles))
	for i, role := range roles {
		protoRoles[i] = RoleFromProto(role)
	}
	return protoRoles
}

// RoleWithPermissionsToProto converts a models.RoleWithPermissions to a securitypb.RoleWithPermissions
func RoleWithPermissionsToProto(role models.RoleWithPermissions) *securitypb.RoleWithPermissions {
	return &securitypb.RoleWithPermissions{
		Role:        RoleToProto(role.Role),
		Permissions: PermissionsToProto(role.Permissions),
	}
}

// RoleWithPermissionsFromProto converts a securitypb.RoleWithPermissions to a models.RoleWithPermissions
func RoleWithPermissionsFromProto(role *securitypb.RoleWithPermissions) models.RoleWithPermissions {
	return models.RoleWithPermissions{
		Role:        RoleFromProto(role.GetRole()),
		Permissions: PermissionsFromProto(role.GetPermissions()),
	}
}

// RolesWithPermissionsToProto converts a models.RolesWithPermissions to a slice of securitypb.RoleWithPermissions
func RolesWithPermissionsToProto(roles models.RolesWithPermissions) []*securitypb.RoleWithPermissions {
	protoRoles := make([]*securitypb.RoleWithPermissions, len(roles))
	for i, role := range roles {
		protoRoles[i] = RoleWithPermissionsToProto(role)
	}
	return protoRoles
}

// RolesWithPermissionsFromProto converts a slice of securitypb.RoleWithPermissions to a models.RolesWithPermissions
func RolesWithPermissionsFromProto(roles []*securitypb.RoleWithPermissions) models.RolesWithPermissions {
	protoRoles := make(models.RolesWithPermissions, len(roles))
	for i, role := range roles {
		protoRoles[i] = RoleWithPermissionsFromProto(role)
	}
	return protoRoles
}
