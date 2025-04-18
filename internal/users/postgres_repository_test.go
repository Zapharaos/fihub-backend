package users_test

import (
	"errors"
	"github.com/Zapharaos/fihub-backend/internal/models"
	"github.com/Zapharaos/fihub-backend/internal/users"
	"github.com/Zapharaos/fihub-backend/test"
	"github.com/google/uuid"
	sqlxmock "github.com/zhashkevych/go-sqlxmock"
	"golang.org/x/crypto/bcrypt"
	"testing"
	"time"
)

// TestPostgresRepository_Create test the Create method
func TestPostgresRepository_Create(t *testing.T) {
	var sqlxMock test.Sqlx
	sqlxMock.CreateFullTestSqlx(t)
	defer sqlxMock.CleanTestSqlx()

	users.ReplaceGlobals(users.NewPostgresRepository(sqlxMock.DB))

	tests := []struct {
		name      string
		user      models.UserWithPassword
		mockSetup func()
		expectErr bool
	}{
		{
			name: "Fail user creation",
			user: models.UserWithPassword{User: models.User{Email: "test@example.com"}, Password: "password"},
			mockSetup: func() {
				sqlxMock.Mock.ExpectQuery("INSERT INTO Users").WillReturnError(errors.New("error"))
			},
			expectErr: true,
		},
		{
			name: "Create user",
			user: models.UserWithPassword{User: models.User{Email: "test@example.com"}, Password: "password"},
			mockSetup: func() {
				rows := sqlxmock.NewRows([]string{"id"}).AddRow(uuid.New())
				sqlxMock.Mock.ExpectQuery("INSERT INTO Users").WillReturnRows(rows)
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			_, err := users.R().Create(tt.user)
			if (err != nil) != tt.expectErr {
				t.Errorf("Create() error = %v, expectErr %v", err, tt.expectErr)
			}
		})
	}
}

// TestPostgresRepository_Get test the Get method
func TestPostgresRepository_Get(t *testing.T) {
	var sqlxMock test.Sqlx
	sqlxMock.CreateFullTestSqlx(t)
	defer sqlxMock.CleanTestSqlx()

	users.ReplaceGlobals(users.NewPostgresRepository(sqlxMock.DB))

	tests := []struct {
		name        string
		userID      uuid.UUID
		mockSetup   func()
		expectErr   bool
		expectFound bool
	}{
		{
			name:   "Fail user retrieval",
			userID: uuid.New(),
			mockSetup: func() {
				sqlxMock.Mock.ExpectQuery("SELECT").WillReturnError(errors.New("error"))
			},
			expectErr:   true,
			expectFound: false,
		},
		{
			name:   "Retrieve user",
			userID: uuid.New(),
			mockSetup: func() {
				rows := sqlxmock.NewRows([]string{"id", "email", "password", "created_at", "updated_at"}).
					AddRow(uuid.New(), "test@example.com", "password", time.Now(), time.Now())
				sqlxMock.Mock.ExpectQuery("SELECT").WillReturnRows(rows)
			},
			expectErr:   false,
			expectFound: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			_, found, err := users.R().Get(tt.userID)
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

// TestPostgresRepository_GetByEmail test the GetByEmail method
func TestPostgresRepository_GetByEmail(t *testing.T) {
	var sqlxMock test.Sqlx
	sqlxMock.CreateFullTestSqlx(t)
	defer sqlxMock.CleanTestSqlx()

	users.ReplaceGlobals(users.NewPostgresRepository(sqlxMock.DB))

	tests := []struct {
		name        string
		email       string
		mockSetup   func()
		expectErr   bool
		expectFound bool
	}{
		{
			name:  "Fail user retrieval by email",
			email: "test@example.com",
			mockSetup: func() {
				sqlxMock.Mock.ExpectQuery("SELECT").WillReturnError(errors.New("error"))
			},
			expectErr:   true,
			expectFound: false,
		},
		{
			name:  "Retrieve user by email",
			email: "test@example.com",
			mockSetup: func() {
				rows := sqlxmock.NewRows([]string{"id", "email", "password", "created_at", "updated_at"}).
					AddRow(uuid.New(), "test@example.com", "password", time.Now(), time.Now())
				sqlxMock.Mock.ExpectQuery("SELECT").WillReturnRows(rows)
			},
			expectErr:   false,
			expectFound: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			_, found, err := users.R().GetByEmail(tt.email)
			if (err != nil) != tt.expectErr {
				t.Errorf("GetByEmail() error = %v, expectErr %v", err, tt.expectErr)
				return
			}
			if found != tt.expectFound {
				t.Errorf("GetByEmail() found = %v, expectFound %v", found, tt.expectFound)
			}
		})
	}
}

// TestPostgresRepository_Exists test the Exists method
func TestPostgresRepository_Exists(t *testing.T) {
	var sqlxMock test.Sqlx
	sqlxMock.CreateFullTestSqlx(t)
	defer sqlxMock.CleanTestSqlx()

	users.ReplaceGlobals(users.NewPostgresRepository(sqlxMock.DB))

	tests := []struct {
		name         string
		mockSetup    func()
		expectErr    bool
		expectExists bool
	}{
		{
			name: "Fail user exists check",
			mockSetup: func() {
				sqlxMock.Mock.ExpectQuery("SELECT").WillReturnError(errors.New("error"))
			},
			expectErr:    true,
			expectExists: false,
		},
		{
			name: "User exists",
			mockSetup: func() {
				rows := sqlxmock.NewRows([]string{"id"}).AddRow(1)
				sqlxMock.Mock.ExpectQuery("SELECT").WillReturnRows(rows)
			},
			expectErr:    false,
			expectExists: true,
		},
		{
			name: "User does not exist",
			mockSetup: func() {
				rows := sqlxmock.NewRows([]string{"id"})
				sqlxMock.Mock.ExpectQuery("SELECT").WillReturnRows(rows)
			},
			expectErr:    false,
			expectExists: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			exists, err := users.R().Exists("")
			if (err != nil) != tt.expectErr {
				t.Errorf("Exists() error = %v, expectErr %v", err, tt.expectErr)
			}
			if exists != tt.expectExists {
				t.Errorf("Exists() exists = %v, expectExists %v", exists, tt.expectExists)
			}
		})
	}
}

// TestPostgresRepository_Authenticate test the Authenticate method
func TestPostgresRepository_Authenticate(t *testing.T) {
	var sqlxMock test.Sqlx
	sqlxMock.CreateFullTestSqlx(t)
	defer sqlxMock.CleanTestSqlx()

	users.ReplaceGlobals(users.NewPostgresRepository(sqlxMock.DB))

	tests := []struct {
		name       string
		mockSetup  func()
		expectErr  bool
		expectAuth bool
	}{
		{
			name: "Fail user authentication",
			mockSetup: func() {
				sqlxMock.Mock.ExpectQuery("SELECT").WillReturnError(errors.New("error"))
			},
			expectErr:  true,
			expectAuth: false,
		},
		{
			name: "Authenticate user",
			mockSetup: func() {
				hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
				rows := sqlxmock.NewRows([]string{"id", "email", "password", "created_at", "updated_at"}).
					AddRow(uuid.New(), "", string(hashedPassword), time.Now(), time.Now())
				sqlxMock.Mock.ExpectQuery("SELECT").WillReturnRows(rows)
			},
			expectErr:  false,
			expectAuth: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			_, auth, err := users.R().Authenticate("", "password")
			if (err != nil) != tt.expectErr {
				t.Errorf("Authenticate() error = %v, expectErr %v", err, tt.expectErr)
			}
			if auth != tt.expectAuth {
				t.Errorf("Authenticate() users = %v, expectAuth %v", auth, tt.expectAuth)
			}
		})
	}
}

// TestPostgresRepository_Update tests the Update method
func TestPostgresRepository_Update(t *testing.T) {
	var sqlxMock test.Sqlx
	sqlxMock.CreateFullTestSqlx(t)
	defer sqlxMock.CleanTestSqlx()

	users.ReplaceGlobals(users.NewPostgresRepository(sqlxMock.DB))

	tests := []struct {
		name      string
		user      models.User
		mockSetup func()
		expectErr bool
	}{
		{
			name: "Fail user update",
			user: models.User{ID: uuid.New(), Email: "test@example.com"},
			mockSetup: func() {
				sqlxMock.Mock.ExpectExec("UPDATE Users").WillReturnError(errors.New("error"))
			},
			expectErr: true,
		},
		{
			name: "Update user",
			user: models.User{ID: uuid.New(), Email: "test@example.com"},
			mockSetup: func() {
				sqlxMock.Mock.ExpectExec("UPDATE Users").WillReturnResult(sqlxmock.NewResult(1, 1))
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			err := users.R().Update(tt.user)
			if (err != nil) != tt.expectErr {
				t.Errorf("Update() error = %v, expectErr %v", err, tt.expectErr)
			}
		})
	}
}

// TestPostgresRepository_UpdateWithPassword tests the UpdateWithPassword method
func TestPostgresRepository_UpdateWithPassword(t *testing.T) {
	var sqlxMock test.Sqlx
	sqlxMock.CreateFullTestSqlx(t)
	defer sqlxMock.CleanTestSqlx()

	users.ReplaceGlobals(users.NewPostgresRepository(sqlxMock.DB))

	tests := []struct {
		name      string
		user      models.UserWithPassword
		mockSetup func()
		expectErr bool
	}{
		{
			name: "Fail user update with password",
			user: models.UserWithPassword{User: models.User{ID: uuid.New(), Email: "test@example.com"}, Password: "newpassword"},
			mockSetup: func() {
				sqlxMock.Mock.ExpectExec("UPDATE Users").WillReturnError(errors.New("error"))
			},
			expectErr: true,
		},
		{
			name: "Update user with password",
			user: models.UserWithPassword{User: models.User{ID: uuid.New(), Email: "test@example.com"}, Password: "newpassword"},
			mockSetup: func() {
				sqlxMock.Mock.ExpectExec("UPDATE Users").WillReturnResult(sqlxmock.NewResult(1, 1))
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			err := users.R().UpdateWithPassword(tt.user)
			if (err != nil) != tt.expectErr {
				t.Errorf("UpdateWithPassword() error = %v, expectErr %v", err, tt.expectErr)
			}
		})
	}
}

// TestPostgresRepository_Delete tests the Delete method
func TestPostgresRepository_Delete(t *testing.T) {
	var sqlxMock test.Sqlx
	sqlxMock.CreateFullTestSqlx(t)
	defer sqlxMock.CleanTestSqlx()

	users.ReplaceGlobals(users.NewPostgresRepository(sqlxMock.DB))

	tests := []struct {
		name      string
		userID    uuid.UUID
		mockSetup func()
		expectErr bool
	}{
		{
			name:   "Fail user delete",
			userID: uuid.New(),
			mockSetup: func() {
				sqlxMock.Mock.ExpectExec("DELETE FROM Users").WillReturnError(errors.New("error"))
			},
			expectErr: true,
		},
		{
			name:   "Delete user",
			userID: uuid.New(),
			mockSetup: func() {
				sqlxMock.Mock.ExpectExec("DELETE FROM Users").WillReturnResult(sqlxmock.NewResult(1, 1))
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			err := users.R().Delete(tt.userID)
			if (err != nil) != tt.expectErr {
				t.Errorf("Delete() error = %v, expectErr %v", err, tt.expectErr)
			}
		})
	}
}

// TestPostgresRepository_GetWithRoles tests the GetWithRoles method
func TestPostgresRepository_GetWithRoles(t *testing.T) {
	var sqlxMock test.Sqlx
	sqlxMock.CreateFullTestSqlx(t)
	defer sqlxMock.CleanTestSqlx()

	users.ReplaceGlobals(users.NewPostgresRepository(sqlxMock.DB))

	tests := []struct {
		name      string
		userID    uuid.UUID
		mockSetup func()
		expectErr bool
	}{
		{
			name:   "Fail user retrieval with roles",
			userID: uuid.New(),
			mockSetup: func() {
				sqlxMock.Mock.ExpectQuery("SELECT").WillReturnError(errors.New("error"))
			},
			expectErr: true,
		},
		{
			name:   "Retrieve user with roles",
			userID: uuid.New(),
			mockSetup: func() {
				rows := sqlxmock.NewRows([]string{"ID", "email", "created_at", "updated_at", "role_id", "role_name"}).
					AddRow(uuid.New(), "test@example.com", time.Now(), time.Now(), uuid.New(), "role_name")
				sqlxMock.Mock.ExpectQuery("SELECT").WillReturnRows(rows)
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			_, err := users.R().GetWithRoles(tt.userID)
			if (err != nil) != tt.expectErr {
				t.Errorf("GetWithRoles() error = %v, expectErr %v", err, tt.expectErr)
			}
		})
	}
}

// TestPostgresRepository_GetAllWithRoles tests the GetAllWithRoles method
func TestPostgresRepository_GetAllWithRoles(t *testing.T) {
	var sqlxMock test.Sqlx
	sqlxMock.CreateFullTestSqlx(t)
	defer sqlxMock.CleanTestSqlx()

	users.ReplaceGlobals(users.NewPostgresRepository(sqlxMock.DB))

	tests := []struct {
		name      string
		mockSetup func()
		expectErr bool
	}{
		{
			name: "Fail users retrieval with roles",
			mockSetup: func() {
				sqlxMock.Mock.ExpectQuery("SELECT").WillReturnError(errors.New("error"))
			},
			expectErr: true,
		},
		{
			name: "Retrieve users with roles",
			mockSetup: func() {
				rows := sqlxmock.NewRows([]string{"ID", "email", "created_at", "updated_at", "role_id", "role_name"}).
					AddRow(uuid.New(), "test@example.com", time.Now(), time.Now(), uuid.New(), "role_name").
					AddRow(uuid.New(), "test2@example.com", time.Now(), time.Now(), uuid.New(), "role_name2")
				sqlxMock.Mock.ExpectQuery("SELECT").WillReturnRows(rows)
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			_, err := users.R().GetAllWithRoles()
			if (err != nil) != tt.expectErr {
				t.Errorf("GetAllWithRoles() error = %v, expectErr %v", err, tt.expectErr)
			}
		})
	}
}

// TestPostgresRepository_GetUsersByRoleID tests the GetUsersByRoleID method
func TestPostgresRepository_GetUsersByRoleID(t *testing.T) {
	var sqlxMock test.Sqlx
	sqlxMock.CreateFullTestSqlx(t)
	defer sqlxMock.CleanTestSqlx()

	users.ReplaceGlobals(users.NewPostgresRepository(sqlxMock.DB))

	tests := []struct {
		name      string
		roleID    uuid.UUID
		mockSetup func()
		expectErr bool
	}{
		{
			name:   "Fail users retrieval by role ID",
			roleID: uuid.New(),
			mockSetup: func() {
				sqlxMock.Mock.ExpectQuery("SELECT").WillReturnError(errors.New("error"))
			},
			expectErr: true,
		},
		{
			name:   "Retrieve users by role ID",
			roleID: uuid.New(),
			mockSetup: func() {
				rows := sqlxmock.NewRows([]string{"ID", "email", "password", "created_at", "updated_at"}).
					AddRow(uuid.New(), "test@example.com", "password", time.Now(), time.Now()).
					AddRow(uuid.New(), "test2@example.com", "password", time.Now(), time.Now())
				sqlxMock.Mock.ExpectQuery("SELECT").WillReturnRows(rows)
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			_, err := users.R().GetUsersByRoleID(tt.roleID)
			if (err != nil) != tt.expectErr {
				t.Errorf("GetUsersByRoleID() error = %v, expectErr %v", err, tt.expectErr)
			}
		})
	}
}

// TestPostgresRepository_UpdateWithRoles tests the UpdateWithRoles method
func TestPostgresRepository_UpdateWithRoles(t *testing.T) {
	var sqlxMock test.Sqlx
	sqlxMock.CreateFullTestSqlx(t)
	defer sqlxMock.CleanTestSqlx()

	users.ReplaceGlobals(users.NewPostgresRepository(sqlxMock.DB))

	tests := []struct {
		name      string
		user      models.UserWithRoles
		roleUUIDs []uuid.UUID
		mockSetup func()
		expectErr bool
	}{
		{
			name:      "Fail user update with roles",
			user:      models.UserWithRoles{User: models.User{ID: uuid.New(), Email: "test@example.com"}},
			roleUUIDs: []uuid.UUID{uuid.New()},
			mockSetup: func() {
				sqlxMock.Mock.ExpectBegin()
				sqlxMock.Mock.ExpectExec("UPDATE Users").WillReturnError(errors.New("error"))
				sqlxMock.Mock.ExpectRollback()
			},
			expectErr: true,
		},
		{
			name:      "Update user with roles",
			user:      models.UserWithRoles{User: models.User{ID: uuid.New(), Email: "test@example.com"}},
			roleUUIDs: []uuid.UUID{uuid.New()},
			mockSetup: func() {
				sqlxMock.Mock.ExpectBegin()
				sqlxMock.Mock.ExpectExec("UPDATE Users").WillReturnResult(sqlxmock.NewResult(1, 1))
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
			err := users.R().UpdateWithRoles(tt.user, tt.roleUUIDs)
			if (err != nil) != tt.expectErr {
				t.Errorf("UpdateWithRoles() error = %v, expectErr %v", err, tt.expectErr)
			}
		})
	}
}

// TestPostgresRepository_SetUserRoles tests the SetUserRoles method
func TestPostgresRepository_SetUserRoles(t *testing.T) {
	var sqlxMock test.Sqlx
	sqlxMock.CreateFullTestSqlx(t)
	defer sqlxMock.CleanTestSqlx()

	users.ReplaceGlobals(users.NewPostgresRepository(sqlxMock.DB))

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
			err := users.R().SetUserRoles(tt.userUUID, tt.roleUUIDs)
			if (err != nil) != tt.expectErr {
				t.Errorf("SetUserRoles() error = %v, expectErr %v", err, tt.expectErr)
			}
		})
	}
}

// TestPostgresRepository_AddUsersRole tests the AddUsersRole method
func TestPostgresRepository_AddUsersRole(t *testing.T) {
	var sqlxMock test.Sqlx
	sqlxMock.CreateFullTestSqlx(t)
	defer sqlxMock.CleanTestSqlx()

	users.ReplaceGlobals(users.NewPostgresRepository(sqlxMock.DB))

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
			err := users.R().AddUsersRole(tt.userUUIDs, tt.roleUUID)
			if (err != nil) != tt.expectErr {
				t.Errorf("AddUsersRole() error = %v, expectErr %v", err, tt.expectErr)
			}
		})
	}
}

// TestPostgresRepository_RemoveUsersRole tests the RemoveUsersRole method
func TestPostgresRepository_RemoveUsersRole(t *testing.T) {
	var sqlxMock test.Sqlx
	sqlxMock.CreateFullTestSqlx(t)
	defer sqlxMock.CleanTestSqlx()

	users.ReplaceGlobals(users.NewPostgresRepository(sqlxMock.DB))

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
			err := users.R().RemoveUsersRole(tt.userUUIDs, tt.roleUUID)
			if (err != nil) != tt.expectErr {
				t.Errorf("RemoveUsersRole() error = %v, expectErr %v", err, tt.expectErr)
			}
		})
	}
}
