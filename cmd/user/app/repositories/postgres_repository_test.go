package repositories_test

import (
	"errors"
	"github.com/Zapharaos/fihub-backend/cmd/user/app/repositories"
	"github.com/Zapharaos/fihub-backend/internal/models"
	"github.com/Zapharaos/fihub-backend/test"
	"github.com/google/uuid"
	sqlxmock "github.com/zhashkevych/go-sqlxmock"
	"golang.org/x/crypto/bcrypt"
	"testing"
	"time"
)

// TestUserPostgresRepository_Create test the RolePostgresRepository.Create method
func TestUserPostgresRepository_Create(t *testing.T) {
	var sqlxMock test.Sqlx
	sqlxMock.CreateFullTestSqlx(t)
	defer sqlxMock.CleanTestSqlx()

	repositories.ReplaceGlobals(repositories.NewPostgresRepository(sqlxMock.DB))

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
			_, err := repositories.R().Create(tt.user)
			if (err != nil) != tt.expectErr {
				t.Errorf("Create() error = %v, expectErr %v", err, tt.expectErr)
			}
		})
	}
}

// TestUserPostgresRepository_Get test the RolePostgresRepository.Get method
func TestUserPostgresRepository_Get(t *testing.T) {
	var sqlxMock test.Sqlx
	sqlxMock.CreateFullTestSqlx(t)
	defer sqlxMock.CleanTestSqlx()

	repositories.ReplaceGlobals(repositories.NewPostgresRepository(sqlxMock.DB))

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
				rows := sqlxmock.NewRows([]string{"id", "email", "created_at", "updated_at"}).
					AddRow(uuid.New(), "test@example.com", time.Now(), time.Now())
				sqlxMock.Mock.ExpectQuery("SELECT").WillReturnRows(rows)
			},
			expectErr:   false,
			expectFound: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			_, found, err := repositories.R().Get(tt.userID)
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

// TestUserPostgresRepository_GetByEmail test the RolePostgresRepository.GetByEmail method
func TestUserPostgresRepository_GetByEmail(t *testing.T) {
	var sqlxMock test.Sqlx
	sqlxMock.CreateFullTestSqlx(t)
	defer sqlxMock.CleanTestSqlx()

	repositories.ReplaceGlobals(repositories.NewPostgresRepository(sqlxMock.DB))

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
				rows := sqlxmock.NewRows([]string{"id", "email", "created_at", "updated_at"}).
					AddRow(uuid.New(), "test@example.com", time.Now(), time.Now())
				sqlxMock.Mock.ExpectQuery("SELECT").WillReturnRows(rows)
			},
			expectErr:   false,
			expectFound: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			_, found, err := repositories.R().GetByEmail(tt.email)
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

// TestUserPostgresRepository_Exists test the RolePostgresRepository.Exists method
func TestUserPostgresRepository_Exists(t *testing.T) {
	var sqlxMock test.Sqlx
	sqlxMock.CreateFullTestSqlx(t)
	defer sqlxMock.CleanTestSqlx()

	repositories.ReplaceGlobals(repositories.NewPostgresRepository(sqlxMock.DB))

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
			exists, err := repositories.R().Exists("")
			if (err != nil) != tt.expectErr {
				t.Errorf("Exists() error = %v, expectErr %v", err, tt.expectErr)
			}
			if exists != tt.expectExists {
				t.Errorf("Exists() exists = %v, expectExists %v", exists, tt.expectExists)
			}
		})
	}
}

// TestUserPostgresRepository_Authenticate test the RolePostgresRepository.Authenticate method
func TestUserPostgresRepository_Authenticate(t *testing.T) {
	var sqlxMock test.Sqlx
	sqlxMock.CreateFullTestSqlx(t)
	defer sqlxMock.CleanTestSqlx()

	repositories.ReplaceGlobals(repositories.NewPostgresRepository(sqlxMock.DB))

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
			_, auth, err := repositories.R().Authenticate("", "password")
			if (err != nil) != tt.expectErr {
				t.Errorf("Authenticate() error = %v, expectErr %v", err, tt.expectErr)
			}
			if auth != tt.expectAuth {
				t.Errorf("Authenticate() users = %v, expectAuth %v", auth, tt.expectAuth)
			}
		})
	}
}

// TestUserPostgresRepository_Update tests the RolePostgresRepository.Update method
func TestUserPostgresRepository_Update(t *testing.T) {
	var sqlxMock test.Sqlx
	sqlxMock.CreateFullTestSqlx(t)
	defer sqlxMock.CleanTestSqlx()

	repositories.ReplaceGlobals(repositories.NewPostgresRepository(sqlxMock.DB))

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
			err := repositories.R().Update(tt.user)
			if (err != nil) != tt.expectErr {
				t.Errorf("Update() error = %v, expectErr %v", err, tt.expectErr)
			}
		})
	}
}

// TestUserPostgresRepository_UpdateWithPassword tests the RolePostgresRepository.UpdateWithPassword method
func TestUserPostgresRepository_UpdateWithPassword(t *testing.T) {
	var sqlxMock test.Sqlx
	sqlxMock.CreateFullTestSqlx(t)
	defer sqlxMock.CleanTestSqlx()

	repositories.ReplaceGlobals(repositories.NewPostgresRepository(sqlxMock.DB))

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
			err := repositories.R().UpdateWithPassword(tt.user)
			if (err != nil) != tt.expectErr {
				t.Errorf("UpdateWithPassword() error = %v, expectErr %v", err, tt.expectErr)
			}
		})
	}
}

// TestUserPostgresRepository_Delete tests the RolePostgresRepository.Delete method
func TestUserPostgresRepository_Delete(t *testing.T) {
	var sqlxMock test.Sqlx
	sqlxMock.CreateFullTestSqlx(t)
	defer sqlxMock.CleanTestSqlx()

	repositories.ReplaceGlobals(repositories.NewPostgresRepository(sqlxMock.DB))

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
			err := repositories.R().Delete(tt.userID)
			if (err != nil) != tt.expectErr {
				t.Errorf("Delete() error = %v, expectErr %v", err, tt.expectErr)
			}
		})
	}
}

// TestUserPostgresRepository_List test the RolePostgresRepository.List method
func TestUserPostgresRepository_List(t *testing.T) {
	var sqlxMock test.Sqlx
	sqlxMock.CreateFullTestSqlx(t)
	defer sqlxMock.CleanTestSqlx()

	repositories.ReplaceGlobals(repositories.NewPostgresRepository(sqlxMock.DB))

	tests := []struct {
		name         string
		userID       uuid.UUID
		mockSetup    func()
		expectErr    bool
		expectLength int
	}{
		{
			name:   "Fail user retrieval",
			userID: uuid.New(),
			mockSetup: func() {
				sqlxMock.Mock.ExpectQuery("SELECT").WillReturnError(errors.New("error"))
			},
			expectErr:    true,
			expectLength: 0,
		},
		{
			name:   "Retrieve user",
			userID: uuid.New(),
			mockSetup: func() {
				rows := sqlxmock.NewRows([]string{"id", "email", "created_at", "updated_at"}).
					AddRow(uuid.New(), "test@example.com", time.Now(), time.Now()).
					AddRow(uuid.New(), "test@example.com", time.Now(), time.Now())
				sqlxMock.Mock.ExpectQuery("SELECT").WillReturnRows(rows)
			},
			expectErr:    false,
			expectLength: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			users, err := repositories.R().List()
			if (err != nil) != tt.expectErr {
				t.Errorf("List() error = %v, expectErr %v", err, tt.expectErr)
				return
			}
			if len(users) != tt.expectLength {
				t.Errorf("List() length = %v, expectFound %v", len(users), tt.expectLength)
			}
		})
	}
}
