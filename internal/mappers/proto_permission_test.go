package mappers

import (
	"github.com/Zapharaos/fihub-backend/gen/go/securitypb"
	"github.com/Zapharaos/fihub-backend/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

// Test_PermissionToProto tests the PermissionToProto function
func Test_PermissionToProto(t *testing.T) {
	// Create test permission
	id := uuid.New()
	permission := models.Permission{
		Id:          id,
		Value:       "read",
		Scope:       "admin",
		Description: "Test permission",
	}

	// Convert to proto
	protoPermission := PermissionToProto(permission)

	// Assert values were correctly converted
	assert.Equal(t, id.String(), protoPermission.Id)
	assert.Equal(t, "read", protoPermission.Value)
	assert.Equal(t, "admin", protoPermission.Scope)
	assert.Equal(t, "Test permission", protoPermission.Description)
}

// Test_PermissionFromProto tests the PermissionFromProto function
func Test_PermissionFromProto(t *testing.T) {
	// Create test permission
	id := uuid.New()
	protoPermission := &securitypb.Permission{
		Id:          id.String(),
		Value:       "read",
		Scope:       "admin",
		Description: "Test permission",
	}

	// Convert to model
	permission := PermissionFromProto(protoPermission)

	// Assert values were correctly converted
	assert.Equal(t, id, permission.Id)
	assert.Equal(t, "read", permission.Value)
	assert.Equal(t, "admin", permission.Scope)
	assert.Equal(t, "Test permission", permission.Description)
}

// Test_PermissionsToProto tests the PermissionsToProto function
func Test_PermissionsToProto(t *testing.T) {
	// Create test permissions
	id1 := uuid.New()
	id2 := uuid.New()
	permissions := models.Permissions{
		{
			Id:          id1,
			Value:       "read",
			Scope:       "admin",
			Description: "Test permission 1",
		},
		{
			Id:          id2,
			Value:       "write",
			Scope:       "all",
			Description: "Test permission 2",
		},
	}

	// Convert to proto
	protoPermissions := PermissionsToProto(permissions)

	// Assert values were correctly converted
	assert.Equal(t, 2, len(protoPermissions))
	assert.Equal(t, id1.String(), protoPermissions[0].Id)
	assert.Equal(t, "read", protoPermissions[0].Value)
	assert.Equal(t, "admin", protoPermissions[0].Scope)
	assert.Equal(t, "Test permission 1", protoPermissions[0].Description)
	assert.Equal(t, id2.String(), protoPermissions[1].Id)
	assert.Equal(t, "write", protoPermissions[1].Value)
	assert.Equal(t, "all", protoPermissions[1].Scope)
	assert.Equal(t, "Test permission 2", protoPermissions[1].Description)
}

// Test_PermissionsFromProto tests the PermissionsFromProto function
func Test_PermissionsFromProto(t *testing.T) {
	// Create test permissions
	id1 := uuid.New()
	id2 := uuid.New()
	protoPermissions := []*securitypb.Permission{
		{
			Id:          id1.String(),
			Value:       "read",
			Scope:       "admin",
			Description: "Test permission 1",
		},
		{
			Id:          id2.String(),
			Value:       "write",
			Scope:       "all",
			Description: "Test permission 2",
		},
	}

	// Convert to model
	permissions := PermissionsFromProto(protoPermissions)

	// Assert values were correctly converted
	assert.Equal(t, 2, len(permissions))
	assert.Equal(t, id1, permissions[0].Id)
	assert.Equal(t, "read", permissions[0].Value)
	assert.Equal(t, "admin", permissions[0].Scope)
	assert.Equal(t, "Test permission 1", permissions[0].Description)
	assert.Equal(t, id2, permissions[1].Id)
	assert.Equal(t, "write", permissions[1].Value)
	assert.Equal(t, "all", permissions[1].Scope)
	assert.Equal(t, "Test permission 2", permissions[1].Description)
}
