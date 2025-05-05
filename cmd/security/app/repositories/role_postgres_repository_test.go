package repositories_test

import (
	"errors"
	"github.com/Zapharaos/fihub-backend/cmd/security/app/repositories"
	"github.com/Zapharaos/fihub-backend/internal/models"
	"github.com/Zapharaos/fihub-backend/test"
	"github.com/google/uuid"
	sqlxmock "github.com/zhashkevych/go-sqlxmock"
	"testing"
)

// TestRolePostgresRepository_Create test the RolePostgresRepository.Create method
func TestRolePostgresRepository_Create(t *testing.T) {
	var sqlxMock test.Sqlx
	sqlxMock.CreateFullTestSqlx(t)
	defer sqlxMock.CleanTestSqlx()

	repositories.ReplaceGlobals(repositories.NewRepository(repositories.NewRolePostgresRepository(sqlxMock.DB), nil))

	tests := []struct {
		name          string
		role          models.Role
		permissionIDs []uuid.UUID
		mockSetup     func()
		expectErr     bool
	}{
		{
			name: "Fail role creation",
			role: models.Role{Name: "role_name"},
			mockSetup: func() {
				sqlxMock.Mock.ExpectBegin()
				sqlxMock.Mock.ExpectExec("INSERT INTO roles").WillReturnError(errors.New("error"))
				sqlxMock.Mock.ExpectRollback()
			},
			expectErr: true,
		},
		{
			name: "Create role without permissions",
			role: models.Role{Name: "role_name"},
			mockSetup: func() {
				sqlxMock.Mock.ExpectBegin()
				sqlxMock.Mock.ExpectExec("INSERT INTO roles").WillReturnResult(sqlxmock.NewResult(1, 1))
				sqlxMock.Mock.ExpectCommit()
			},
			expectErr: false,
		},
		{
			name:          "Create role with permissions",
			role:          models.Role{Name: "role_name"},
			permissionIDs: []uuid.UUID{uuid.New(), uuid.New()},
			mockSetup: func() {
				sqlxMock.Mock.ExpectBegin()
				sqlxMock.Mock.ExpectExec("INSERT INTO roles").WillReturnResult(sqlxmock.NewResult(1, 1))
				sqlxMock.Mock.ExpectExec("INSERT INTO role_permissions").WillReturnResult(sqlxmock.NewResult(1, 2))
				sqlxMock.Mock.ExpectCommit()
			},
			expectErr: false,
		},
		{
			name:          "Fail to set role permissions",
			role:          models.Role{Name: "role_name"},
			permissionIDs: []uuid.UUID{uuid.New(), uuid.New()},
			mockSetup: func() {
				sqlxMock.Mock.ExpectBegin()
				sqlxMock.Mock.ExpectExec("INSERT INTO roles").WillReturnResult(sqlxmock.NewResult(1, 1))
				sqlxMock.Mock.ExpectExec("INSERT INTO role_permissions").WillReturnError(errors.New("error"))
				sqlxMock.Mock.ExpectRollback()
			},
			expectErr: true,
		},
		{
			name:          "Fail to commit transaction",
			role:          models.Role{Name: "role_name"},
			permissionIDs: []uuid.UUID{uuid.New(), uuid.New()},
			mockSetup: func() {
				sqlxMock.Mock.ExpectBegin()
				sqlxMock.Mock.ExpectExec("INSERT INTO roles").WillReturnResult(sqlxmock.NewResult(1, 1))
				sqlxMock.Mock.ExpectExec("INSERT INTO role_permissions").WillReturnResult(sqlxmock.NewResult(1, 2))
				sqlxMock.Mock.ExpectCommit().WillReturnError(errors.New("error"))
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			_, err := repositories.R().R().Create(tt.role, tt.permissionIDs)
			if (err != nil) != tt.expectErr {
				t.Errorf("Create() error = %v, expectErr %v", err, tt.expectErr)
			}
		})
	}
}

// TestRolePostgresRepository_Get test the RolePostgresRepository.Get method
func TestRolePostgresRepository_Get(t *testing.T) {
	var sqlxMock test.Sqlx
	sqlxMock.CreateFullTestSqlx(t)
	defer sqlxMock.CleanTestSqlx()

	repositories.ReplaceGlobals(repositories.NewRepository(repositories.NewRolePostgresRepository(sqlxMock.DB), nil))

	tests := []struct {
		name        string
		roleID      uuid.UUID
		mockSetup   func()
		expectErr   bool
		expectFound bool
	}{
		{
			name:   "Fail role retrieval",
			roleID: uuid.New(),
			mockSetup: func() {
				sqlxMock.Mock.ExpectQuery("SELECT").WillReturnError(errors.New("error"))
			},
			expectErr:   true,
			expectFound: false,
		},
		{
			name:   "Retrieve role",
			roleID: uuid.New(),
			mockSetup: func() {
				rows := sqlxmock.NewRows([]string{"id", "name"}).
					AddRow(uuid.New(), "role_name")
				sqlxMock.Mock.ExpectQuery("SELECT").WillReturnRows(rows)
			},
			expectErr:   false,
			expectFound: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			_, found, err := repositories.R().R().Get(tt.roleID)
			if (err != nil) != tt.expectErr {
				t.Errorf("Get() error = %v, expectErr %v", err, tt.expectErr)
				return
			}
			if found != tt.expectFound {
				t.Errorf("Get() found = %v, expectFound %v", found, tt.expectFound)
			}
		})
	}
}

// TestRolePostgresRepository_GetByName test the RolePostgresRepository.GetByName method
func TestRolePostgresRepository_GetByName(t *testing.T) {
	var sqlxMock test.Sqlx
	sqlxMock.CreateFullTestSqlx(t)
	defer sqlxMock.CleanTestSqlx()

	repositories.ReplaceGlobals(repositories.NewRepository(repositories.NewRolePostgresRepository(sqlxMock.DB), nil))

	tests := []struct {
		name        string
		roleName    string
		mockSetup   func()
		expectErr   bool
		expectFound bool
	}{
		{
			name:     "Fail role retrieval by name",
			roleName: "role_name",
			mockSetup: func() {
				sqlxMock.Mock.ExpectQuery("SELECT").WillReturnError(errors.New("error"))
			},
			expectErr:   true,
			expectFound: false,
		},
		{
			name:     "Retrieve role by name",
			roleName: "role_name",
			mockSetup: func() {
				rows := sqlxmock.NewRows([]string{"id", "name"}).
					AddRow(uuid.New(), "role_name")
				sqlxMock.Mock.ExpectQuery("SELECT").WillReturnRows(rows)
			},
			expectErr:   false,
			expectFound: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			_, found, err := repositories.R().R().GetByName(tt.roleName)
			if (err != nil) != tt.expectErr {
				t.Errorf("GetByName() error = %v, expectErr %v", err, tt.expectErr)
				return
			}
			if found != tt.expectFound {
				t.Errorf("GetByName() found = %v, expectFound %v", found, tt.expectFound)
			}
		})
	}
}

// TestRolePostgresRepository_GetWithPermissions test the RolePostgresRepository.GetWithPermissions method
func TestRolePostgresRepository_GetWithPermissions(t *testing.T) {
	var sqlxMock test.Sqlx
	sqlxMock.CreateFullTestSqlx(t)
	defer sqlxMock.CleanTestSqlx()

	repositories.ReplaceGlobals(repositories.NewRepository(repositories.NewRolePostgresRepository(sqlxMock.DB), nil))

	tests := []struct {
		name        string
		roleID      uuid.UUID
		mockSetup   func()
		expectErr   bool
		expectFound bool
	}{
		{
			name:   "Fail role retrieval with permissions",
			roleID: uuid.New(),
			mockSetup: func() {
				sqlxMock.Mock.ExpectQuery("SELECT").WillReturnError(errors.New("error"))
			},
			expectErr:   true,
			expectFound: false,
		},
		{
			name:   "Retrieve role with permissions",
			roleID: uuid.New(),
			mockSetup: func() {
				rows := sqlxmock.NewRows([]string{"id", "name", "permission_id", "value", "scope", "description"}).
					AddRow(uuid.New(), "role_name", uuid.New(), "value", "scope", "description")
				sqlxMock.Mock.ExpectQuery("SELECT").WillReturnRows(rows)
			},
			expectErr:   false,
			expectFound: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			_, found, err := repositories.R().R().GetWithPermissions(tt.roleID)
			if (err != nil) != tt.expectErr {
				t.Errorf("GetWithPermissions() error = %v, expectErr %v", err, tt.expectErr)
				return
			}
			if found != tt.expectFound {
				t.Errorf("GetWithPermissions() found = %v, expectFound %v", found, tt.expectFound)
			}
		})
	}
}

// TestRolePostgresRepository_Update test the RolePostgresRepository.Update method
func TestRolePostgresRepository_Update(t *testing.T) {
	var sqlxMock test.Sqlx
	sqlxMock.CreateFullTestSqlx(t)
	defer sqlxMock.CleanTestSqlx()

	repositories.ReplaceGlobals(repositories.NewRepository(repositories.NewRolePostgresRepository(sqlxMock.DB), nil))

	tests := []struct {
		name          string
		role          models.Role
		permissionIDs []uuid.UUID
		mockSetup     func()
		expectErr     bool
	}{
		{
			name: "Fail role update",
			role: models.Role{Id: uuid.New(), Name: "role_name"},
			mockSetup: func() {
				sqlxMock.Mock.ExpectBegin()
				sqlxMock.Mock.ExpectExec("UPDATE roles").WillReturnError(errors.New("error"))
				sqlxMock.Mock.ExpectRollback()
			},
			expectErr: true,
		},
		{
			name: "Update role without permissions",
			role: models.Role{Id: uuid.New(), Name: "role_name"},
			mockSetup: func() {
				sqlxMock.Mock.ExpectBegin()
				sqlxMock.Mock.ExpectExec("UPDATE roles").WillReturnResult(sqlxmock.NewResult(1, 1))
				sqlxMock.Mock.ExpectExec("DELETE FROM role_permissions").WillReturnResult(sqlxmock.NewResult(1, 1))
				sqlxMock.Mock.ExpectCommit()
			},
			expectErr: false,
		},
		{
			name:          "Update role with permissions",
			role:          models.Role{Id: uuid.New(), Name: "role_name"},
			permissionIDs: []uuid.UUID{uuid.New(), uuid.New()},
			mockSetup: func() {
				sqlxMock.Mock.ExpectBegin()
				sqlxMock.Mock.ExpectExec("UPDATE roles").WillReturnResult(sqlxmock.NewResult(1, 1))
				sqlxMock.Mock.ExpectExec("DELETE FROM role_permissions").WillReturnResult(sqlxmock.NewResult(1, 1))
				sqlxMock.Mock.ExpectExec("INSERT INTO role_permissions").WillReturnResult(sqlxmock.NewResult(1, 2))
				sqlxMock.Mock.ExpectCommit()
			},
			expectErr: false,
		},
		{
			name:          "Fail to set role permissions",
			role:          models.Role{Id: uuid.New(), Name: "role_name"},
			permissionIDs: []uuid.UUID{uuid.New(), uuid.New()},
			mockSetup: func() {
				sqlxMock.Mock.ExpectBegin()
				sqlxMock.Mock.ExpectExec("UPDATE roles").WillReturnResult(sqlxmock.NewResult(1, 1))
				sqlxMock.Mock.ExpectExec("DELETE FROM role_permissions").WillReturnResult(sqlxmock.NewResult(1, 1))
				sqlxMock.Mock.ExpectExec("INSERT INTO role_permissions").WillReturnError(errors.New("error"))
				sqlxMock.Mock.ExpectRollback()
			},
			expectErr: true,
		},
		{
			name:          "Fail to commit transaction",
			role:          models.Role{Id: uuid.New(), Name: "role_name"},
			permissionIDs: []uuid.UUID{uuid.New(), uuid.New()},
			mockSetup: func() {
				sqlxMock.Mock.ExpectBegin()
				sqlxMock.Mock.ExpectExec("UPDATE roles").WillReturnResult(sqlxmock.NewResult(1, 1))
				sqlxMock.Mock.ExpectExec("DELETE FROM role_permissions").WillReturnResult(sqlxmock.NewResult(1, 1))
				sqlxMock.Mock.ExpectExec("INSERT INTO role_permissions").WillReturnResult(sqlxmock.NewResult(1, 2))
				sqlxMock.Mock.ExpectCommit().WillReturnError(errors.New("error"))
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			err := repositories.R().R().Update(tt.role, tt.permissionIDs)
			if (err != nil) != tt.expectErr {
				t.Errorf("Update() error = %v, expectErr %v", err, tt.expectErr)
			}
		})
	}
}

// TestRolePostgresRepository_Delete test the RolePostgresRepository.Delete method
func TestRolePostgresRepository_Delete(t *testing.T) {
	var sqlxMock test.Sqlx
	sqlxMock.CreateFullTestSqlx(t)
	defer sqlxMock.CleanTestSqlx()

	repositories.ReplaceGlobals(repositories.NewRepository(repositories.NewRolePostgresRepository(sqlxMock.DB), nil))

	tests := []struct {
		name      string
		roleID    uuid.UUID
		mockSetup func()
		expectErr bool
	}{
		{
			name:   "Fail role deletion",
			roleID: uuid.New(),
			mockSetup: func() {
				sqlxMock.Mock.ExpectExec("DELETE FROM roles").WillReturnError(errors.New("error"))
			},
			expectErr: true,
		},
		{
			name:   "Delete role",
			roleID: uuid.New(),
			mockSetup: func() {
				sqlxMock.Mock.ExpectExec("DELETE FROM roles").WillReturnResult(sqlxmock.NewResult(1, 1))
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			err := repositories.R().R().Delete(tt.roleID)
			if (err != nil) != tt.expectErr {
				t.Errorf("Delete() error = %v, expectErr %v", err, tt.expectErr)
			}
		})
	}
}

// TestRolePostgresRepository_List test the RolePostgresRepository.List method
func TestRolePostgresRepository_List(t *testing.T) {
	var sqlxMock test.Sqlx
	sqlxMock.CreateFullTestSqlx(t)
	defer sqlxMock.CleanTestSqlx()

	repositories.ReplaceGlobals(repositories.NewRepository(repositories.NewRolePostgresRepository(sqlxMock.DB), nil))

	tests := []struct {
		name        string
		mockSetup   func()
		expectErr   bool
		expectCount int
	}{
		{
			name: "Fail role retrieval",
			mockSetup: func() {
				sqlxMock.Mock.ExpectQuery("SELECT").WillReturnError(errors.New("error"))
			},
			expectErr:   true,
			expectCount: 0,
		},
		{
			name: "Retrieve roles",
			mockSetup: func() {
				rows := sqlxmock.NewRows([]string{"id", "name"}).
					AddRow(uuid.New(), "role_name1").
					AddRow(uuid.New(), "role_name2")
				sqlxMock.Mock.ExpectQuery("SELECT").WillReturnRows(rows)
			},
			expectErr:   false,
			expectCount: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			roles, err := repositories.R().R().List()
			if (err != nil) != tt.expectErr {
				t.Errorf("GetAll() error = %v, expectErr %v", err, tt.expectErr)
				return
			}
			if len(roles) != tt.expectCount {
				t.Errorf("GetAll() count = %v, expectCount %v", len(roles), tt.expectCount)
			}
		})
	}
}

// TestRolePostgresRepository_ListByUserId test the RolePostgresRepository.ListByUserId method
func TestRolePostgresRepository_ListByUserId(t *testing.T) {
	var sqlxMock test.Sqlx
	sqlxMock.CreateFullTestSqlx(t)
	defer sqlxMock.CleanTestSqlx()

	repositories.ReplaceGlobals(repositories.NewRepository(repositories.NewRolePostgresRepository(sqlxMock.DB), nil))

	tests := []struct {
		name        string
		userID      uuid.UUID
		mockSetup   func()
		expectErr   bool
		expectCount int
	}{
		{
			name:   "Fail roles retrieval by user ID",
			userID: uuid.New(),
			mockSetup: func() {
				sqlxMock.Mock.ExpectQuery("SELECT").WillReturnError(errors.New("error"))
			},
			expectErr:   true,
			expectCount: 0,
		},
		{
			name:   "Retrieve roles by user ID",
			userID: uuid.New(),
			mockSetup: func() {
				rows := sqlxmock.NewRows([]string{"id", "name"}).
					AddRow(uuid.New(), "role_name1").
					AddRow(uuid.New(), "role_name2")
				sqlxMock.Mock.ExpectQuery("SELECT").WillReturnRows(rows)
			},
			expectErr:   false,
			expectCount: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			roles, err := repositories.R().R().ListByUserId(tt.userID)
			if (err != nil) != tt.expectErr {
				t.Errorf("GetRolesByUserId() error = %v, expectErr %v", err, tt.expectErr)
				return
			}
			if len(roles) != tt.expectCount {
				t.Errorf("GetRolesByUserId() count = %v, expectCount %v", len(roles), tt.expectCount)
			}
		})
	}
}

// TestRolePostgresRepository_ListWithPermissions test the RolePostgresRepository.ListWithPermissions method
func TestRolePostgresRepository_ListWithPermissions(t *testing.T) {
	var sqlxMock test.Sqlx
	sqlxMock.CreateFullTestSqlx(t)
	defer sqlxMock.CleanTestSqlx()

	repositories.ReplaceGlobals(repositories.NewRepository(repositories.NewRolePostgresRepository(sqlxMock.DB), nil))

	tests := []struct {
		name        string
		mockSetup   func()
		expectErr   bool
		expectCount int
	}{
		{
			name: "Fail roles retrieval with permissions",
			mockSetup: func() {
				sqlxMock.Mock.ExpectQuery("SELECT").WillReturnError(errors.New("error"))
			},
			expectErr:   true,
			expectCount: 0,
		},
		{
			name: "Retrieve roles with permissions",
			mockSetup: func() {
				rows := sqlxmock.NewRows([]string{"id", "name", "permission_id", "value", "scope", "description"}).
					AddRow(uuid.New(), "role_name1", uuid.New(), "value1", "scope1", "description1").
					AddRow(uuid.New(), "role_name2", uuid.New(), "value2", "scope2", "description2")
				sqlxMock.Mock.ExpectQuery("SELECT").WillReturnRows(rows)
			},
			expectErr:   false,
			expectCount: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			roles, err := repositories.R().R().ListWithPermissions()
			if (err != nil) != tt.expectErr {
				t.Errorf("GetAllWithPermissions() error = %v, expectErr %v", err, tt.expectErr)
				return
			}
			if len(roles) != tt.expectCount {
				t.Errorf("GetAllWithPermissions() count = %v, expectCount %v", len(roles), tt.expectCount)
			}
		})
	}
}

// TestRolePostgresRepository_ListWithPermissionsByUserId test the RolePostgresRepository.ListWithPermissionsByUserId method
func TestRolePostgresRepository_ListWithPermissionsByUserId(t *testing.T) {
	var sqlxMock test.Sqlx
	sqlxMock.CreateFullTestSqlx(t)
	defer sqlxMock.CleanTestSqlx()

	// Mock data
	roleID1 := uuid.New()
	roleID2 := uuid.New()
	permissionID1 := uuid.New()
	permissionID2 := uuid.New()

	repositories.ReplaceGlobals(repositories.NewRepository(repositories.NewRolePostgresRepository(sqlxMock.DB), nil))

	tests := []struct {
		name         string
		userID       uuid.UUID
		mockSetup    func()
		expectErr    bool
		expectResult models.RolesWithPermissions
	}{
		{
			name:   "Fail roles and permissions retrieval for user",
			userID: uuid.New(),
			mockSetup: func() {
				sqlxMock.Mock.ExpectQuery("SELECT").WillReturnError(errors.New("error"))
			},
			expectErr:    true,
			expectResult: models.RolesWithPermissions{},
		},
		{
			name:   "Retrieve role with permissions for user",
			userID: uuid.New(),
			mockSetup: func() {
				rows := sqlxmock.NewRows([]string{"r.id", "r.name", "p.id", "p.value", "p.scope", "p.description"}).
					AddRow(roleID1, "name1", permissionID1, "value1", "scope1", "description1").
					AddRow(roleID1, "name1", permissionID2, "value2", "scope2", "description2").
					AddRow(roleID2, "name2", permissionID2, "value2", "scope2", "description2")
				sqlxMock.Mock.ExpectQuery("SELECT").WillReturnRows(rows)
			},
			expectErr: false,
			expectResult: models.RolesWithPermissions{
				{
					Role:        models.Role{Id: roleID1},
					Permissions: []models.Permission{{Id: permissionID1}, {Id: permissionID2}},
				},
				{
					Role:        models.Role{Id: roleID2},
					Permissions: []models.Permission{{Id: permissionID2}},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			rp, err := repositories.R().R().ListWithPermissionsByUserId(tt.userID)
			if (err != nil) != tt.expectErr {
				t.Errorf("ListWithPermissionsByUserId() error = %v, expectErr %v", err, tt.expectErr)
				return
			}
			if len(rp) != len(tt.expectResult) {
				t.Errorf("ListWithPermissionsByUserId() length = %v, expectResultLength %v", len(rp), len(tt.expectResult))
			}
			for i, role := range rp {
				if role.Role.Id != tt.expectResult[i].Role.Id {
					t.Errorf("ListWithPermissionsByUserId() role ID = %v, expectResultRoleID %v", role.Role.Id, tt.expectResult[i].Role.Id)
				}
				if len(role.Permissions) != len(tt.expectResult[i].Permissions) {
					t.Errorf("ListWithPermissionsByUserId() permissions length = %v, expectResultPermissionsLength %v", len(role.Permissions), len(tt.expectResult[i].Permissions))
				}
				for j, permission := range role.Permissions {
					if permission.Id != tt.expectResult[i].Permissions[j].Id {
						t.Errorf("ListWithPermissionsByUserId() permission ID = %v, expectResultPermissionID %v", permission.Id, tt.expectResult[i].Permissions[j].Id)
					}
				}
			}
		})
	}
}

// TestRolePostgresRepository_SetForUser tests the RolePostgresRepository.SetForUser method
func TestRolePostgresRepository_SetForUser(t *testing.T) {
	var sqlxMock test.Sqlx
	sqlxMock.CreateFullTestSqlx(t)
	defer sqlxMock.CleanTestSqlx()

	repositories.ReplaceGlobals(repositories.NewRepository(repositories.NewRolePostgresRepository(sqlxMock.DB), nil))

	tests := []struct {
		name      string
		userUUID  uuid.UUID
		roleUUIDs []uuid.UUID
		mockSetup func()
		expectErr bool
	}{
		{
			name:      "Fail set user roles",
			userUUID:  uuid.New(),
			roleUUIDs: []uuid.UUID{uuid.New()},
			mockSetup: func() {
				sqlxMock.Mock.ExpectBegin()
				sqlxMock.Mock.ExpectExec("DELETE FROM user_roles").WillReturnError(errors.New("error"))
				sqlxMock.Mock.ExpectRollback()
			},
			expectErr: true,
		},
		{
			name:      "Set user roles",
			userUUID:  uuid.New(),
			roleUUIDs: []uuid.UUID{uuid.New()},
			mockSetup: func() {
				sqlxMock.Mock.ExpectBegin()
				sqlxMock.Mock.ExpectExec("DELETE FROM user_roles").WillReturnResult(sqlxmock.NewResult(1, 1))
				sqlxMock.Mock.ExpectExec("INSERT INTO user_roles").WillReturnResult(sqlxmock.NewResult(1, 1))
				sqlxMock.Mock.ExpectCommit()
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			err := repositories.R().R().SetForUser(tt.userUUID, tt.roleUUIDs)
			if (err != nil) != tt.expectErr {
				t.Errorf("SetUserRoles() error = %v, expectErr %v", err, tt.expectErr)
			}
		})
	}
}

// TestRolePostgresRepository_AddUsersRole tests the RolePostgresRepository.AddToUsers method
func TestRolePostgresRepository_AddUsersRole(t *testing.T) {
	var sqlxMock test.Sqlx
	sqlxMock.CreateFullTestSqlx(t)
	defer sqlxMock.CleanTestSqlx()

	repositories.ReplaceGlobals(repositories.NewRepository(repositories.NewRolePostgresRepository(sqlxMock.DB), nil))

	tests := []struct {
		name      string
		userUUIDs []uuid.UUID
		roleUUID  uuid.UUID
		mockSetup func()
		expectErr bool
	}{
		{
			name:      "Fail add users role",
			userUUIDs: []uuid.UUID{uuid.New()},
			roleUUID:  uuid.New(),
			mockSetup: func() {
				sqlxMock.Mock.ExpectBegin()
				sqlxMock.Mock.ExpectExec("INSERT INTO user_roles").WillReturnError(errors.New("error"))
				sqlxMock.Mock.ExpectRollback()
			},
			expectErr: true,
		},
		{
			name:      "Add users role",
			userUUIDs: []uuid.UUID{uuid.New()},
			roleUUID:  uuid.New(),
			mockSetup: func() {
				sqlxMock.Mock.ExpectBegin()
				sqlxMock.Mock.ExpectExec("INSERT INTO user_roles").WillReturnResult(sqlxmock.NewResult(1, 1))
				sqlxMock.Mock.ExpectCommit()
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			err := repositories.R().R().AddToUsers(tt.userUUIDs, tt.roleUUID)
			if (err != nil) != tt.expectErr {
				t.Errorf("AddUsersRole() error = %v, expectErr %v", err, tt.expectErr)
			}
		})
	}
}

// TestRolePostgresRepository_RemoveUsersRole tests the RolePostgresRepository.RemoveFromUsers method
func TestRolePostgresRepository_RemoveUsersRole(t *testing.T) {
	var sqlxMock test.Sqlx
	sqlxMock.CreateFullTestSqlx(t)
	defer sqlxMock.CleanTestSqlx()

	repositories.ReplaceGlobals(repositories.NewRepository(repositories.NewRolePostgresRepository(sqlxMock.DB), nil))

	tests := []struct {
		name      string
		userUUIDs []uuid.UUID
		roleUUID  uuid.UUID
		mockSetup func()
		expectErr bool
	}{
		{
			name:      "Fail remove users role",
			userUUIDs: []uuid.UUID{uuid.New()},
			roleUUID:  uuid.New(),
			mockSetup: func() {
				sqlxMock.Mock.ExpectBegin()
				sqlxMock.Mock.ExpectExec("DELETE FROM user_roles").WillReturnError(errors.New("error"))
				sqlxMock.Mock.ExpectRollback()
			},
			expectErr: true,
		},
		{
			name:      "Remove users role",
			userUUIDs: []uuid.UUID{uuid.New()},
			roleUUID:  uuid.New(),
			mockSetup: func() {
				sqlxMock.Mock.ExpectBegin()
				sqlxMock.Mock.ExpectExec("DELETE FROM user_roles").WillReturnResult(sqlxmock.NewResult(1, 1))
				sqlxMock.Mock.ExpectCommit()
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			err := repositories.R().R().RemoveFromUsers(tt.userUUIDs, tt.roleUUID)
			if (err != nil) != tt.expectErr {
				t.Errorf("RemoveUsersRole() error = %v, expectErr %v", err, tt.expectErr)
			}
		})
	}
}

// TestRolePostgresRepository_ListUsersByRoleId test the RolePostgresRepository.ListUsersByRoleId method
func TestRolePostgresRepository_ListUsersByRoleId(t *testing.T) {
	var sqlxMock test.Sqlx
	sqlxMock.CreateFullTestSqlx(t)
	defer sqlxMock.CleanTestSqlx()

	repositories.ReplaceGlobals(repositories.NewRepository(repositories.NewRolePostgresRepository(sqlxMock.DB), nil))

	tests := []struct {
		name        string
		roleID      uuid.UUID
		mockSetup   func()
		expectErr   bool
		expectCount int
	}{
		{
			name:   "Fail permissions retrieval by role ID",
			roleID: uuid.New(),
			mockSetup: func() {
				sqlxMock.Mock.ExpectQuery("SELECT").WillReturnError(errors.New("error"))
			},
			expectErr:   true,
			expectCount: 0,
		},
		{
			name:   "Retrieve permissions by role ID",
			roleID: uuid.New(),
			mockSetup: func() {
				rows := sqlxmock.NewRows([]string{"id"}).
					AddRow(uuid.New()).
					AddRow(uuid.New())
				sqlxMock.Mock.ExpectQuery("SELECT").WillReturnRows(rows)
			},
			expectErr:   false,
			expectCount: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			p, err := repositories.R().R().ListUsersByRoleId(tt.roleID)
			if (err != nil) != tt.expectErr {
				t.Errorf("GetAllByRoleId() error = %v, expectErr %v", err, tt.expectErr)
				return
			}
			if len(p) != tt.expectCount {
				t.Errorf("GetAllByRoleId() count = %v, expectCount %v", len(p), tt.expectCount)
			}
		})
	}
}

// TestRolePostgresRepository_ListUsers test the RolePostgresRepository.ListUsers method
func TestRolePostgresRepository_ListUsers(t *testing.T) {
	var sqlxMock test.Sqlx
	sqlxMock.CreateFullTestSqlx(t)
	defer sqlxMock.CleanTestSqlx()

	repositories.ReplaceGlobals(repositories.NewRepository(repositories.NewRolePostgresRepository(sqlxMock.DB), nil))

	tests := []struct {
		name        string
		roleID      uuid.UUID
		mockSetup   func()
		expectErr   bool
		expectCount int
	}{
		{
			name:   "Fail users retrieval",
			roleID: uuid.New(),
			mockSetup: func() {
				sqlxMock.Mock.ExpectQuery("SELECT").WillReturnError(errors.New("error"))
			},
			expectErr:   true,
			expectCount: 0,
		},
		{
			name:   "Retrieve users",
			roleID: uuid.New(),
			mockSetup: func() {
				rows := sqlxmock.NewRows([]string{"id"}).
					AddRow(uuid.New()).
					AddRow(uuid.New())
				sqlxMock.Mock.ExpectQuery("SELECT").WillReturnRows(rows)
			},
			expectErr:   false,
			expectCount: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			p, err := repositories.R().R().ListUsers()
			if (err != nil) != tt.expectErr {
				t.Errorf("GetAllByRoleId() error = %v, expectErr %v", err, tt.expectErr)
				return
			}
			if len(p) != tt.expectCount {
				t.Errorf("GetAllByRoleId() count = %v, expectCount %v", len(p), tt.expectCount)
			}
		})
	}
}

// TestRolePostgresRepository_SetPermissionsByRoleId test the RolePostgresRepository.SetPermissionsByRoleId method
func TestRolePostgresRepository_SetPermissionsByRoleId(t *testing.T) {
	var sqlxMock test.Sqlx
	sqlxMock.CreateFullTestSqlx(t)
	defer sqlxMock.CleanTestSqlx()

	repositories.ReplaceGlobals(repositories.NewRepository(repositories.NewRolePostgresRepository(sqlxMock.DB), nil))

	tests := []struct {
		name          string
		roleID        uuid.UUID
		permissionIDs []uuid.UUID
		mockSetup     func()
		expectErr     bool
	}{
		{
			name:          "Fail to set role permissions",
			roleID:        uuid.New(),
			permissionIDs: []uuid.UUID{uuid.New(), uuid.New()},
			mockSetup: func() {
				sqlxMock.Mock.ExpectBegin()
				sqlxMock.Mock.ExpectExec("DELETE FROM role_permissions").WillReturnError(errors.New("error"))
				sqlxMock.Mock.ExpectRollback()
			},
			expectErr: true,
		},
		{
			name:          "Set role permissions",
			roleID:        uuid.New(),
			permissionIDs: []uuid.UUID{uuid.New(), uuid.New()},
			mockSetup: func() {
				sqlxMock.Mock.ExpectBegin()
				sqlxMock.Mock.ExpectExec("DELETE FROM role_permissions").WillReturnResult(sqlxmock.NewResult(1, 1))
				sqlxMock.Mock.ExpectExec("INSERT INTO role_permissions").WillReturnResult(sqlxmock.NewResult(1, 2))
				sqlxMock.Mock.ExpectCommit()
			},
			expectErr: false,
		},
		{
			name:          "Set role permissions with empty permissions",
			roleID:        uuid.New(),
			permissionIDs: []uuid.UUID{},
			mockSetup: func() {
				sqlxMock.Mock.ExpectBegin()
				sqlxMock.Mock.ExpectExec("DELETE FROM role_permissions").WillReturnResult(sqlxmock.NewResult(1, 1))
				sqlxMock.Mock.ExpectCommit()
			},
			expectErr: false,
		},
		{
			name:          "Fail to delete existing permissions",
			roleID:        uuid.New(),
			permissionIDs: []uuid.UUID{uuid.New(), uuid.New()},
			mockSetup: func() {
				sqlxMock.Mock.ExpectBegin()
				sqlxMock.Mock.ExpectExec("DELETE FROM role_permissions").WillReturnError(errors.New("error"))
				sqlxMock.Mock.ExpectRollback()
			},
			expectErr: true,
		},
		{
			name:          "Fail to insert new permissions",
			roleID:        uuid.New(),
			permissionIDs: []uuid.UUID{uuid.New(), uuid.New()},
			mockSetup: func() {
				sqlxMock.Mock.ExpectBegin()
				sqlxMock.Mock.ExpectExec("DELETE FROM role_permissions").WillReturnResult(sqlxmock.NewResult(1, 1))
				sqlxMock.Mock.ExpectExec("INSERT INTO role_permissions").WillReturnError(errors.New("error"))
				sqlxMock.Mock.ExpectRollback()
			},
			expectErr: true,
		},
		{
			name:          "Fail to commit transaction",
			roleID:        uuid.New(),
			permissionIDs: []uuid.UUID{uuid.New(), uuid.New()},
			mockSetup: func() {
				sqlxMock.Mock.ExpectBegin()
				sqlxMock.Mock.ExpectExec("DELETE FROM role_permissions").WillReturnResult(sqlxmock.NewResult(1, 1))
				sqlxMock.Mock.ExpectExec("INSERT INTO role_permissions").WillReturnResult(sqlxmock.NewResult(1, 2))
				sqlxMock.Mock.ExpectCommit().WillReturnError(errors.New("error"))
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			err := repositories.R().R().SetPermissionsByRoleId(tt.roleID, tt.permissionIDs)
			if (err != nil) != tt.expectErr {
				t.Errorf("SetRolePermissions() error = %v, expectErr %v", err, tt.expectErr)
			}
		})
	}
}

// TestRolePostgresRepository_ListPermissionsByRoleId test the RolePostgresRepository.ListPermissionsByRoleId method
func TestRolePostgresRepository_ListPermissionsByRoleId(t *testing.T) {
	var sqlxMock test.Sqlx
	sqlxMock.CreateFullTestSqlx(t)
	defer sqlxMock.CleanTestSqlx()

	repositories.ReplaceGlobals(repositories.NewRepository(repositories.NewRolePostgresRepository(sqlxMock.DB), nil))

	tests := []struct {
		name        string
		mockSetup   func()
		expectErr   bool
		expectCount int
	}{
		{
			name: "Fail permissions retrieval",
			mockSetup: func() {
				sqlxMock.Mock.ExpectQuery("SELECT").WillReturnError(errors.New("error"))
			},
			expectErr:   true,
			expectCount: 0,
		},
		{
			name: "Retrieve permissions",
			mockSetup: func() {
				rows := sqlxmock.NewRows([]string{"id", "value", "scope", "description"}).
					AddRow(uuid.New(), "value1", "scope1", "description1").
					AddRow(uuid.New(), "value2", "scope2", "description2")
				sqlxMock.Mock.ExpectQuery("SELECT").WillReturnRows(rows)
			},
			expectErr:   false,
			expectCount: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			p, err := repositories.R().R().ListPermissionsByRoleId(uuid.New())
			if (err != nil) != tt.expectErr {
				t.Errorf("GetAll() error = %v, expectErr %v", err, tt.expectErr)
				return
			}
			if len(p) != tt.expectCount {
				t.Errorf("GetAll() count = %v, expectCount %v", len(p), tt.expectCount)
			}
		})
	}
}

// TestRolePostgresRepository_ListPermissionsByUserId test the RolePostgresRepository.ListPermissionsByUserId method
func TestRolePostgresRepository_ListPermissionsByUserId(t *testing.T) {
	var sqlxMock test.Sqlx
	sqlxMock.CreateFullTestSqlx(t)
	defer sqlxMock.CleanTestSqlx()

	repositories.ReplaceGlobals(repositories.NewRepository(repositories.NewRolePostgresRepository(sqlxMock.DB), nil))

	tests := []struct {
		name        string
		mockSetup   func()
		expectErr   bool
		expectCount int
	}{
		{
			name: "Fail permissions retrieval",
			mockSetup: func() {
				sqlxMock.Mock.ExpectQuery("SELECT").WillReturnError(errors.New("error"))
			},
			expectErr:   true,
			expectCount: 0,
		},
		{
			name: "Retrieve permissions",
			mockSetup: func() {
				rows := sqlxmock.NewRows([]string{"id", "value", "scope", "description"}).
					AddRow(uuid.New(), "value1", "scope1", "description1").
					AddRow(uuid.New(), "value2", "scope2", "description2")
				sqlxMock.Mock.ExpectQuery("SELECT").WillReturnRows(rows)
			},
			expectErr:   false,
			expectCount: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			p, err := repositories.R().R().ListPermissionsByUserId(uuid.New())
			if (err != nil) != tt.expectErr {
				t.Errorf("GetAll() error = %v, expectErr %v", err, tt.expectErr)
				return
			}
			if len(p) != tt.expectCount {
				t.Errorf("GetAll() count = %v, expectCount %v", len(p), tt.expectCount)
			}
		})
	}
}
