package permissions_test

import (
	"errors"
	"github.com/Zapharaos/fihub-backend/internal/auth/permissions"
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

	permissions.ReplaceGlobals(permissions.NewPostgresRepository(sqlxMock.DB))

	tests := []struct {
		name         string
		permissionID uuid.UUID
		mockSetup    func()
		expectErr    bool
		expectFound  bool
	}{
		{
			name:         "Fail permission retrieval",
			permissionID: uuid.New(),
			mockSetup: func() {
				sqlxMock.Mock.ExpectQuery("SELECT").WillReturnError(errors.New("error"))
			},
			expectErr:   true,
			expectFound: false,
		},
		{
			name:         "Retrieve permission",
			permissionID: uuid.New(),
			mockSetup: func() {
				rows := sqlxmock.NewRows([]string{"id", "value", "scope", "description"}).
					AddRow(uuid.New(), "value", "scope", "description")
				sqlxMock.Mock.ExpectQuery("SELECT").WillReturnRows(rows)
			},
			expectErr:   false,
			expectFound: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			_, found, err := permissions.R().Get(tt.permissionID)
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

// TestPostgresRepository_Create test the Create method
func TestPostgresRepository_Create(t *testing.T) {
	var sqlxMock test.Sqlx
	sqlxMock.CreateFullTestSqlx(t)
	defer sqlxMock.CleanTestSqlx()

	permissions.ReplaceGlobals(permissions.NewPostgresRepository(sqlxMock.DB))

	tests := []struct {
		name       string
		permission permissions.Permission
		mockSetup  func()
		expectErr  bool
	}{
		{
			name:       "Fail permission creation",
			permission: permissions.Permission{},
			mockSetup: func() {
				sqlxMock.Mock.ExpectQuery("INSERT INTO permissions").WillReturnError(errors.New("error"))
			},
			expectErr: true,
		},
		{
			name:       "Create permission",
			permission: permissions.Permission{},
			mockSetup: func() {
				sqlxMock.Mock.ExpectQuery("INSERT INTO permissions").WillReturnRows(sqlxmock.NewRows([]string{"id"}).AddRow(uuid.New()))
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			_, err := permissions.R().Create(tt.permission)
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

	permissions.ReplaceGlobals(permissions.NewPostgresRepository(sqlxMock.DB))

	tests := []struct {
		name       string
		permission permissions.Permission
		mockSetup  func()
		expectErr  bool
	}{
		{
			name: "Fail permission update",
			mockSetup: func() {
				sqlxMock.Mock.ExpectExec("UPDATE permissions").WillReturnError(errors.New("error"))
			},
			expectErr: true,
		},
		{
			name: "Update permission",
			mockSetup: func() {
				sqlxMock.Mock.ExpectExec("UPDATE permissions").WillReturnResult(sqlxmock.NewResult(1, 1))
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			err := permissions.R().Update(tt.permission)
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

	permissions.ReplaceGlobals(permissions.NewPostgresRepository(sqlxMock.DB))

	tests := []struct {
		name         string
		permissionID uuid.UUID
		mockSetup    func()
		expectErr    bool
	}{
		{
			name:         "Fail permission delete",
			permissionID: uuid.New(),
			mockSetup: func() {
				sqlxMock.Mock.ExpectExec("DELETE FROM permissions").WillReturnError(errors.New("error"))
			},
			expectErr: true,
		},
		{
			name:         "Delete permission",
			permissionID: uuid.New(),
			mockSetup: func() {
				sqlxMock.Mock.ExpectExec("DELETE FROM permissions").WillReturnResult(sqlxmock.NewResult(1, 1))
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			err := permissions.R().Delete(tt.permissionID)
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

	permissions.ReplaceGlobals(permissions.NewPostgresRepository(sqlxMock.DB))

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
			p, err := permissions.R().GetAll()
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

// TestPostgresRepository_GetAllByRoleId test the GetAllByRoleId method
func TestPostgresRepository_GetAllByRoleId(t *testing.T) {
	var sqlxMock test.Sqlx
	sqlxMock.CreateFullTestSqlx(t)
	defer sqlxMock.CleanTestSqlx()

	permissions.ReplaceGlobals(permissions.NewPostgresRepository(sqlxMock.DB))

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
			p, err := permissions.R().GetAllByRoleId(tt.roleID)
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

// TestPostgresRepository_GetAllForUser test the GetAllForUser method
func TestPostgresRepository_GetAllForUser(t *testing.T) {
	var sqlxMock test.Sqlx
	sqlxMock.CreateFullTestSqlx(t)
	defer sqlxMock.CleanTestSqlx()

	permissions.ReplaceGlobals(permissions.NewPostgresRepository(sqlxMock.DB))

	tests := []struct {
		name        string
		userID      uuid.UUID
		mockSetup   func()
		expectErr   bool
		expectCount int
	}{
		{
			name:   "Fail permissions retrieval for user",
			userID: uuid.New(),
			mockSetup: func() {
				sqlxMock.Mock.ExpectQuery("SELECT").WillReturnError(errors.New("error"))
			},
			expectErr:   true,
			expectCount: 0,
		},
		{
			name:   "Retrieve permissions for user",
			userID: uuid.New(),
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
			p, err := permissions.R().GetAllForUser(tt.userID)
			if (err != nil) != tt.expectErr {
				t.Errorf("GetAllForUser() error = %v, expectErr %v", err, tt.expectErr)
				return
			}
			if len(p) != tt.expectCount {
				t.Errorf("GetAllForUser() count = %v, expectCount %v", len(p), tt.expectCount)
			}
		})
	}
}
