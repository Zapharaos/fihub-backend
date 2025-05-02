package mappers

import (
	"github.com/Zapharaos/fihub-backend/gen/go/securitypb"
	"github.com/Zapharaos/fihub-backend/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

// Test_RoleToProto tests the RoleToProto function
func Test_RoleToProto(t *testing.T) {
	// Create test role
	id := uuid.New()
	role := models.Role{
		Id:   id,
		Name: "admin",
	}

	// Convert to proto
	protoRole := RoleToProto(role)

	// Assert values were correctly converted
	assert.Equal(t, id.String(), protoRole.Id)
	assert.Equal(t, "admin", protoRole.Name)
}

// Test_RoleFromProto tests the RoleFromProto function
func Test_RoleFromProto(t *testing.T) {
	// Create test role
	id := uuid.New()
	protoRole := &securitypb.Role{
		Id:   id.String(),
		Name: "admin",
	}

	// Convert to model
	role := RoleFromProto(protoRole)

	// Assert values were correctly converted
	assert.Equal(t, id, role.Id)
	assert.Equal(t, "admin", role.Name)
}

// Test_RolesToProto tests the RolesToProto function
func Test_RolesToProto(t *testing.T) {
	// Create test roles
	id1 := uuid.New()
	id2 := uuid.New()
	roles := models.Roles{
		{
			Id:   id1,
			Name: "admin",
		},
		{
			Id:   id2,
			Name: "user",
		},
	}

	// Convert to proto
	protoRoles := RolesToProto(roles)

	// Assert values were correctly converted
	assert.Equal(t, 2, len(protoRoles))
	assert.Equal(t, id1.String(), protoRoles[0].Id)
	assert.Equal(t, "admin", protoRoles[0].Name)
	assert.Equal(t, id2.String(), protoRoles[1].Id)
	assert.Equal(t, "user", protoRoles[1].Name)
}

// Test_RolesFromProto tests the RolesFromProto function
func Test_RolesFromProto(t *testing.T) {
	// Create test roles
	id1 := uuid.New()
	id2 := uuid.New()
	protoRoles := []*securitypb.Role{
		{
			Id:   id1.String(),
			Name: "admin",
		},
		{
			Id:   id2.String(),
			Name: "user",
		},
	}

	// Convert to model
	roles := RolesFromProto(protoRoles)

	// Assert values were correctly converted
	assert.Equal(t, 2, len(roles))
	assert.Equal(t, id1, roles[0].Id)
	assert.Equal(t, "admin", roles[0].Name)
	assert.Equal(t, id2, roles[1].Id)
	assert.Equal(t, "user", roles[1].Name)
}

// Test_RoleWithPermissionsToProto tests the RoleWithPermissionsToProto function
func Test_RoleWithPermissionsToProto(t *testing.T) {
	// Create test role with permissions
	id := uuid.New()
	role := models.RoleWithPermissions{
		Role: models.Role{
			Id:   id,
			Name: "admin",
		},
		Permissions: models.Permissions{
			{
				Id:          uuid.New(),
				Value:       "read",
				Scope:       "admin",
				Description: "Test permission",
			},
		},
	}

	// Convert to proto
	protoRole := RoleWithPermissionsToProto(role)

	// Assert values were correctly converted
	assert.Equal(t, id.String(), protoRole.Role.Id)
	assert.Equal(t, "admin", protoRole.Role.Name)
	assert.Equal(t, 1, len(protoRole.Permissions))
	assert.Equal(t, "read", protoRole.Permissions[0].Value)
}

// Test_RoleWithPermissionsFromProto tests the RoleWithPermissionsFromProto function
func Test_RoleWithPermissionsFromProto(t *testing.T) {
	// Create test role with permissions
	id := uuid.New()
	protoRole := &securitypb.RoleWithPermissions{
		Role: &securitypb.Role{
			Id:   id.String(),
			Name: "admin",
		},
		Permissions: []*securitypb.Permission{
			{
				Id:          uuid.New().String(),
				Value:       "read",
				Scope:       "admin",
				Description: "Test permission",
			},
		},
	}

	// Convert to model
	role := RoleWithPermissionsFromProto(protoRole)

	// Assert values were correctly converted
	assert.Equal(t, id, role.Role.Id)
	assert.Equal(t, "admin", role.Role.Name)
	assert.Equal(t, 1, len(role.Permissions))
	assert.Equal(t, "read", role.Permissions[0].Value)
}

// Test_RolesWithPermissionsToProto tests the RolesWithPermissionsToProto function
func Test_RolesWithPermissionsToProto(t *testing.T) {
	// Create test roles with permissions
	id1 := uuid.New()
	id2 := uuid.New()
	roles := models.RolesWithPermissions{
		{
			Role: models.Role{
				Id:   id1,
				Name: "admin",
			},
			Permissions: models.Permissions{
				{
					Id:          uuid.New(),
					Value:       "read",
					Scope:       "admin",
					Description: "Test permission",
				},
			},
		},
		{
			Role: models.Role{
				Id:   id2,
				Name: "user",
			},
			Permissions: models.Permissions{
				{
					Id:          uuid.New(),
					Value:       "write",
					Scope:       "all",
					Description: "Test permission 2",
				},
			},
		},
	}

	// Convert to proto
	protoRoles := RolesWithPermissionsToProto(roles)

	// Assert values were correctly converted
	assert.Equal(t, 2, len(protoRoles))
	assert.Equal(t, id1.String(), protoRoles[0].Role.Id)
	assert.Equal(t, "admin", protoRoles[0].Role.Name)
	assert.Equal(t, 1, len(protoRoles[0].Permissions))
	assert.Equal(t, "read", protoRoles[0].Permissions[0].Value)
	assert.Equal(t, id2.String(), protoRoles[1].Role.Id)
	assert.Equal(t, "user", protoRoles[1].Role.Name)
	assert.Equal(t, 1, len(protoRoles[1].Permissions))
	assert.Equal(t, "write", protoRoles[1].Permissions[0].Value)
}

// Test_RolesWithPermissionsFromProto tests the RolesWithPermissionsFromProto function
func Test_RolesWithPermissionsFromProto(t *testing.T) {
	// Create test roles with permissions
	id1 := uuid.New()
	id2 := uuid.New()
	protoRoles := []*securitypb.RoleWithPermissions{
		{
			Role: &securitypb.Role{
				Id:   id1.String(),
				Name: "admin",
			},
			Permissions: []*securitypb.Permission{
				{
					Id:          uuid.New().String(),
					Value:       "read",
					Scope:       "admin",
					Description: "Test permission",
				},
			},
		},
		{
			Role: &securitypb.Role{
				Id:   id2.String(),
				Name: "user",
			},
			Permissions: []*securitypb.Permission{
				{
					Id:          uuid.New().String(),
					Value:       "write",
					Scope:       "all",
					Description: "Test permission 2",
				},
			},
		},
	}

	// Convert to model
	roles := RolesWithPermissionsFromProto(protoRoles)

	// Assert values were correctly converted
	assert.Equal(t, 2, len(roles))
	assert.Equal(t, id1, roles[0].Role.Id)
	assert.Equal(t, "admin", roles[0].Role.Name)
	assert.Equal(t, 1, len(roles[0].Permissions))
	assert.Equal(t, "read", roles[0].Permissions[0].Value)
	assert.Equal(t, id2, roles[1].Role.Id)
	assert.Equal(t, "user", roles[1].Role.Name)
	assert.Equal(t, 1, len(roles[1].Permissions))
	assert.Equal(t, "write", roles[1].Permissions[0].Value)
}
