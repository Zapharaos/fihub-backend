package mappers

import (
	"github.com/Zapharaos/fihub-backend/internal/models"
	"github.com/Zapharaos/fihub-backend/protogen"
	"github.com/google/uuid"
)

// RoleToProto converts a models.Role to a protogen.Role
func RoleToProto(role models.Role) *protogen.Role {
	return &protogen.Role{
		Id:   role.Id.String(),
		Name: role.Name,
	}
}

// RoleFromProto converts a protogen.Role to a models.Role
func RoleFromProto(role *protogen.Role) models.Role {
	return models.Role{
		Id:   uuid.MustParse(role.GetId()),
		Name: role.GetName(),
	}
}

// RolesToProto converts a models.Roles to a slice of protogen.Role
func RolesToProto(roles models.Roles) []*protogen.Role {
	protoRoles := make([]*protogen.Role, len(roles))
	for i, role := range roles {
		protoRoles[i] = RoleToProto(role)
	}
	return protoRoles
}

// RolesFromProto converts a slice of protogen.Role to a models.Roles
func RolesFromProto(roles []*protogen.Role) models.Roles {
	protoRoles := make(models.Roles, len(roles))
	for i, role := range roles {
		protoRoles[i] = RoleFromProto(role)
	}
	return protoRoles
}

// RoleWithPermissionsToProto converts a models.RoleWithPermissions to a protogen.RoleWithPermissions
func RoleWithPermissionsToProto(role models.RoleWithPermissions) *protogen.RoleWithPermissions {
	return &protogen.RoleWithPermissions{
		Role:        RoleToProto(role.Role),
		Permissions: PermissionsToProto(role.Permissions),
	}
}

// RoleWithPermissionsFromProto converts a protogen.RoleWithPermissions to a models.RoleWithPermissions
func RoleWithPermissionsFromProto(role *protogen.RoleWithPermissions) models.RoleWithPermissions {
	return models.RoleWithPermissions{
		Role:        RoleFromProto(role.GetRole()),
		Permissions: PermissionsFromProto(role.GetPermissions()),
	}
}

// RolesWithPermissionsToProto converts a models.RolesWithPermissions to a slice of protogen.RoleWithPermissions
func RolesWithPermissionsToProto(roles models.RolesWithPermissions) []*protogen.RoleWithPermissions {
	protoRoles := make([]*protogen.RoleWithPermissions, len(roles))
	for i, role := range roles {
		protoRoles[i] = RoleWithPermissionsToProto(role)
	}
	return protoRoles
}

// RolesWithPermissionsFromProto converts a slice of protogen.RoleWithPermissions to a models.RolesWithPermissions
func RolesWithPermissionsFromProto(roles []*protogen.RoleWithPermissions) models.RolesWithPermissions {
	protoRoles := make(models.RolesWithPermissions, len(roles))
	for i, role := range roles {
		protoRoles[i] = RoleWithPermissionsFromProto(role)
	}
	return protoRoles
}
