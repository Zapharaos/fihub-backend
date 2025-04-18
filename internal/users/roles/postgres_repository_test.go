package roles_test

import (
	"errors"
	"github.com/Zapharaos/fihub-backend/internal/models"
	"github.com/Zapharaos/fihub-backend/internal/users/roles"
	"github.com/Zapharaos/fihub-backend/test"
	"github.com/google/uuid"
	sqlxmock "github.com/zhashkevych/go-sqlxmock"
	"testing"
)

// TestPostgresRepository_Get test the Get method
func TestPostgresRepository_Get(t *testing.T) {
	var sqlxMock test.Sqlx
	sqlxMock.CreateFullTestSqlx(t)
	defer sqlxMock.CleanTestSqlx()

	roles.ReplaceGlobals(roles.NewPostgresRepository(sqlxMock.DB))

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
			_, found, err := roles.R().Get(tt.roleID)
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

// TestPostgresRepository_GetByName test the GetByName method
func TestPostgresRepository_GetByName(t *testing.T) {
	var sqlxMock test.Sqlx
	sqlxMock.CreateFullTestSqlx(t)
	defer sqlxMock.CleanTestSqlx()

	roles.ReplaceGlobals(roles.NewPostgresRepository(sqlxMock.DB))

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
			_, found, err := roles.R().GetByName(tt.roleName)
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

// TestPostgresRepository_Create test the Create method
func TestPostgresRepository_Create(t *testing.T) {
	var sqlxMock test.Sqlx
	sqlxMock.CreateFullTestSqlx(t)
	defer sqlxMock.CleanTestSqlx()

	roles.ReplaceGlobals(roles.NewPostgresRepository(sqlxMock.DB))

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
			_, err := roles.R().Create(tt.role, tt.permissionIDs)
			if (err != nil) != tt.expectErr {
				t.Errorf("Create() error = %v, expectErr %v", err, tt.expectErr)
			}
		})
	}
}

// TestPostgresRepository_Update test the Update method
func TestPostgresRepository_Update(t *testing.T) {
	var sqlxMock test.Sqlx
	sqlxMock.CreateFullTestSqlx(t)
	defer sqlxMock.CleanTestSqlx()

	roles.ReplaceGlobals(roles.NewPostgresRepository(sqlxMock.DB))

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
			err := roles.R().Update(tt.role, tt.permissionIDs)
			if (err != nil) != tt.expectErr {
				t.Errorf("Update() error = %v, expectErr %v", err, tt.expectErr)
			}
		})
	}
}

// TestPostgresRepository_Delete test the Delete method
func TestPostgresRepository_Delete(t *testing.T) {
	var sqlxMock test.Sqlx
	sqlxMock.CreateFullTestSqlx(t)
	defer sqlxMock.CleanTestSqlx()

	roles.ReplaceGlobals(roles.NewPostgresRepository(sqlxMock.DB))

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
			err := roles.R().Delete(tt.roleID)
			if (err != nil) != tt.expectErr {
				t.Errorf("Delete() error = %v, expectErr %v", err, tt.expectErr)
			}
		})
	}
}

// TestPostgresRepository_GetAll test the GetAll method
func TestPostgresRepository_GetAll(t *testing.T) {
	var sqlxMock test.Sqlx
	sqlxMock.CreateFullTestSqlx(t)
	defer sqlxMock.CleanTestSqlx()

	roles.ReplaceGlobals(roles.NewPostgresRepository(sqlxMock.DB))

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
			roles, err := roles.R().GetAll()
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

// TestPostgresRepository_GetWithPermissions test the GetWithPermissions method
func TestPostgresRepository_GetWithPermissions(t *testing.T) {
	var sqlxMock test.Sqlx
	sqlxMock.CreateFullTestSqlx(t)
	defer sqlxMock.CleanTestSqlx()

	roles.ReplaceGlobals(roles.NewPostgresRepository(sqlxMock.DB))

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
			_, found, err := roles.R().GetWithPermissions(tt.roleID)
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

// TestPostgresRepository_GetAllWithPermissions test the GetAllWithPermissions method
func TestPostgresRepository_GetAllWithPermissions(t *testing.T) {
	var sqlxMock test.Sqlx
	sqlxMock.CreateFullTestSqlx(t)
	defer sqlxMock.CleanTestSqlx()

	roles.ReplaceGlobals(roles.NewPostgresRepository(sqlxMock.DB))

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
			roles, err := roles.R().GetAllWithPermissions()
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

// TestPostgresRepository_GetRolesByUserId test the GetRolesByUserId method
func TestPostgresRepository_GetRolesByUserId(t *testing.T) {
	var sqlxMock test.Sqlx
	sqlxMock.CreateFullTestSqlx(t)
	defer sqlxMock.CleanTestSqlx()

	roles.ReplaceGlobals(roles.NewPostgresRepository(sqlxMock.DB))

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
			roles, err := roles.R().GetRolesByUserId(tt.userID)
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

// TestPostgresRepository_SetRolePermissions test the SetRolePermissions method
func TestPostgresRepository_SetRolePermissions(t *testing.T) {
	var sqlxMock test.Sqlx
	sqlxMock.CreateFullTestSqlx(t)
	defer sqlxMock.CleanTestSqlx()

	roles.ReplaceGlobals(roles.NewPostgresRepository(sqlxMock.DB))

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
			err := roles.R().SetRolePermissions(tt.roleID, tt.permissionIDs)
			if (err != nil) != tt.expectErr {
				t.Errorf("SetRolePermissions() error = %v, expectErr %v", err, tt.expectErr)
			}
		})
	}
}
