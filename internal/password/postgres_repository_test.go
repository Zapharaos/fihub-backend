package password_test

import (
	"database/sql"
	"errors"
	"github.com/Zapharaos/fihub-backend/internal/models"
	"github.com/Zapharaos/fihub-backend/internal/password"
	"github.com/Zapharaos/fihub-backend/test"
	"github.com/google/uuid"
	sqlxmock "github.com/zhashkevych/go-sqlxmock"
	"testing"
	"time"
)

// TestPostgresRepository_Create test the Create method
func TestPostgresRepository_Create(t *testing.T) {
	var sqlxMock test.Sqlx
	sqlxMock.CreateFullTestSqlx(t)
	defer sqlxMock.CleanTestSqlx()

	password.ReplaceGlobals(password.NewPostgresRepository(sqlxMock.DB))

	tests := []struct {
		name      string
		request   models.PasswordRequest
		mockSetup func()
		expectErr bool
	}{
		{
			name:    "Fail request creation",
			request: models.PasswordRequest{},
			mockSetup: func() {
				sqlxMock.Mock.ExpectQuery("INSERT INTO password_reset_tokens").WillReturnError(errors.New("error"))
			},
			expectErr: true,
		},
		{
			name:    "Create request",
			request: models.PasswordRequest{},
			mockSetup: func() {
				rows := sqlxmock.NewRows([]string{"id", "user_id", "token", "expires_at", "created_at"}).
					AddRow(uuid.New(), uuid.New(), "token", time.Now().Add(1*time.Hour), time.Now())
				sqlxMock.Mock.ExpectQuery("INSERT INTO password_reset_tokens").WillReturnRows(rows)
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			_, err := password.R().Create(tt.request)
			if (err != nil) != tt.expectErr {
				t.Errorf("Create() error = %v, expectErr %v", err, tt.expectErr)
			}
		})
	}
}

// TestPostgresRepository_GetRequestID tests the GetRequestID method
func TestPostgresRepository_GetRequestID(t *testing.T) {
	var sqlxMock test.Sqlx
	sqlxMock.CreateFullTestSqlx(t)
	defer sqlxMock.CleanTestSqlx()

	password.ReplaceGlobals(password.NewPostgresRepository(sqlxMock.DB))

	tests := []struct {
		name      string
		userID    uuid.UUID
		token     string
		mockSetup func()
		expectErr bool
	}{
		{
			name:   "Fail to get request ID",
			userID: uuid.New(),
			token:  "invalid_token",
			mockSetup: func() {
				sqlxMock.Mock.ExpectQuery("SELECT id FROM password_reset_tokens").WillReturnError(sql.ErrNoRows)
			},
			expectErr: false,
		},
		{
			name:   "Get request ID",
			userID: uuid.New(),
			token:  "valid_token",
			mockSetup: func() {
				rows := sqlxmock.NewRows([]string{"id"}).AddRow(uuid.New())
				sqlxMock.Mock.ExpectQuery("SELECT id FROM password_reset_tokens").WillReturnRows(rows)
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			_, err := password.R().GetRequestID(tt.userID, tt.token)
			if (err != nil) != tt.expectErr {
				t.Errorf("GetRequestID() error = %v, expectErr %v", err, tt.expectErr)
				return
			}
		})
	}
}

// TestPostgresRepository_GetExpiresAt tests the GetExpiresAt method
func TestPostgresRepository_GetExpiresAt(t *testing.T) {
	var sqlxMock test.Sqlx
	sqlxMock.CreateFullTestSqlx(t)
	defer sqlxMock.CleanTestSqlx()

	password.ReplaceGlobals(password.NewPostgresRepository(sqlxMock.DB))

	tests := []struct {
		name      string
		userID    uuid.UUID
		mockSetup func()
		expectErr bool
	}{
		{
			name:   "Fail to get expires at",
			userID: uuid.New(),
			mockSetup: func() {
				sqlxMock.Mock.ExpectQuery("SELECT expires_at FROM password_reset_tokens").WillReturnError(sql.ErrNoRows)
			},
			expectErr: true,
		},
		{
			name:   "Get expires at",
			userID: uuid.New(),
			mockSetup: func() {
				expiresAt := time.Now().Add(1 * time.Hour)
				rows := sqlxmock.NewRows([]string{"expires_at"}).AddRow(expiresAt)
				sqlxMock.Mock.ExpectQuery("SELECT expires_at FROM password_reset_tokens").WillReturnRows(rows)
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			_, err := password.R().GetExpiresAt(tt.userID)
			if (err != nil) != tt.expectErr {
				t.Errorf("GetExpiresAt() error = %v, expectErr %v", err, tt.expectErr)
				return
			}
		})
	}
}

// TestPostgresRepository_Delete test the Delete method
func TestPostgresRepository_Delete(t *testing.T) {
	var sqlxMock test.Sqlx
	sqlxMock.CreateFullTestSqlx(t)
	defer sqlxMock.CleanTestSqlx()

	password.ReplaceGlobals(password.NewPostgresRepository(sqlxMock.DB))

	tests := []struct {
		name      string
		requestID uuid.UUID
		mockSetup func()
		expectErr bool
	}{
		{
			name: "Fail request delete",
			mockSetup: func() {
				sqlxMock.Mock.ExpectExec("DELETE FROM password_reset_tokens").WillReturnError(errors.New("error"))
			},
			expectErr: true,
		},
		{
			name: "Delete request",
			mockSetup: func() {
				sqlxMock.Mock.ExpectExec("DELETE FROM password_reset_tokens").WillReturnResult(sqlxmock.NewResult(1, 1))
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			err := password.R().Delete(tt.requestID)
			if (err != nil) != tt.expectErr {
				t.Errorf("Delete() error = %v, expectErr %v", err, tt.expectErr)
			}
		})
	}
}

// TestPostgresRepository_Valid tests the Valid method
func TestPostgresRepository_Valid(t *testing.T) {
	var sqlxMock test.Sqlx
	sqlxMock.CreateFullTestSqlx(t)
	defer sqlxMock.CleanTestSqlx()

	password.ReplaceGlobals(password.NewPostgresRepository(sqlxMock.DB))

	tests := []struct {
		name        string
		userID      uuid.UUID
		requestID   uuid.UUID
		mockSetup   func()
		expectErr   bool
		expectValid bool
	}{
		{
			name:      "Fail to validate request",
			userID:    uuid.New(),
			requestID: uuid.New(),
			mockSetup: func() {
				sqlxMock.Mock.ExpectQuery("SELECT \\* FROM password_reset_tokens").WillReturnError(errors.New("error"))
			},
			expectErr:   true,
			expectValid: false,
		},
		{
			name:      "Valid request",
			userID:    uuid.New(),
			requestID: uuid.New(),
			mockSetup: func() {
				rows := sqlxmock.NewRows([]string{"id"}).AddRow(uuid.New())
				sqlxMock.Mock.ExpectQuery("SELECT \\* FROM password_reset_tokens").WillReturnRows(rows)
			},
			expectErr:   false,
			expectValid: true,
		},
		{
			name:      "Invalid request",
			userID:    uuid.New(),
			requestID: uuid.New(),
			mockSetup: func() {
				rows := sqlxmock.NewRows([]string{"id"})
				sqlxMock.Mock.ExpectQuery("SELECT \\* FROM password_reset_tokens").WillReturnRows(rows)
			},
			expectErr:   false,
			expectValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			valid, err := password.R().Valid(tt.userID, tt.requestID)
			if (err != nil) != tt.expectErr {
				t.Errorf("Valid() error = %v, expectErr %v", err, tt.expectErr)
				return
			}
			if valid != tt.expectValid {
				t.Errorf("Valid() valid = %v, expectValid %v", valid, tt.expectValid)
			}
		})
	}
}

// TestPostgresRepository_ValidForUser tests the ValidForUser method
func TestPostgresRepository_ValidForUser(t *testing.T) {
	var sqlxMock test.Sqlx
	sqlxMock.CreateFullTestSqlx(t)
	defer sqlxMock.CleanTestSqlx()

	password.ReplaceGlobals(password.NewPostgresRepository(sqlxMock.DB))

	tests := []struct {
		name        string
		userID      uuid.UUID
		mockSetup   func()
		expectErr   bool
		expectValid bool
	}{
		{
			name:   "Fail to validate request for user",
			userID: uuid.New(),
			mockSetup: func() {
				sqlxMock.Mock.ExpectQuery("SELECT \\* FROM password_reset_tokens").WillReturnError(errors.New("error"))
			},
			expectErr:   true,
			expectValid: false,
		},
		{
			name:   "Valid request for user",
			userID: uuid.New(),
			mockSetup: func() {
				rows := sqlxmock.NewRows([]string{"id"}).AddRow(uuid.New())
				sqlxMock.Mock.ExpectQuery("SELECT \\* FROM password_reset_tokens").WillReturnRows(rows)
			},
			expectErr:   false,
			expectValid: true,
		},
		{
			name:   "Invalid request for user",
			userID: uuid.New(),
			mockSetup: func() {
				rows := sqlxmock.NewRows([]string{"id"})
				sqlxMock.Mock.ExpectQuery("SELECT \\* FROM password_reset_tokens").WillReturnRows(rows)
			},
			expectErr:   false,
			expectValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			valid, err := password.R().ValidForUser(tt.userID)
			if (err != nil) != tt.expectErr {
				t.Errorf("ValidForUser() error = %v, expectErr %v", err, tt.expectErr)
				return
			}
			if valid != tt.expectValid {
				t.Errorf("ValidForUser() valid = %v, expectValid %v", valid, tt.expectValid)
			}
		})
	}
}
