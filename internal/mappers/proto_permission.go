package mappers

import (
	"github.com/Zapharaos/fihub-backend/internal/models"
	"github.com/Zapharaos/fihub-backend/protogen"
	"github.com/google/uuid"
)

func PermissionToProto(permission models.Permission) *protogen.Permission {
	return &protogen.Permission{
		Id:          permission.Id.String(),
		Value:       permission.Value,
		Scope:       permission.Scope,
		Description: permission.Description,
	}
}

func PermissionFromProto(permission *protogen.Permission) models.Permission {
	return models.Permission{
		Id:          uuid.MustParse(permission.GetId()),
		Value:       permission.GetValue(),
		Scope:       permission.GetScope(),
		Description: permission.GetDescription(),
	}
}

func PermissionsToProto(permissions models.Permissions) []*protogen.Permission {
	protoPermissions := make([]*protogen.Permission, len(permissions))
	for i, permission := range permissions {
		protoPermissions[i] = PermissionToProto(permission)
	}
	return protoPermissions
}

func PermissionsFromProto(permissions []*protogen.Permission) models.Permissions {
	protoPermissions := make(models.Permissions, len(permissions))
	for i, permission := range permissions {
		protoPermissions[i] = PermissionFromProto(permission)
	}
	return protoPermissions
}
