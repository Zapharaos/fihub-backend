package repositories_test

import (
	"errors"
	"github.com/Zapharaos/fihub-backend/cmd/transaction/app/repositories"
	"github.com/Zapharaos/fihub-backend/internal/models"
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

	repositories.ReplaceGlobals(repositories.NewPostgresRepository(sqlxMock.DB))

	tests := []struct {
		name        string
		transaction models.TransactionInput
		mockSetup   func()
		expectErr   bool
	}{

		{
			name:        "Fail transaction creation",
			transaction: models.TransactionInput{},
			mockSetup: func() {
				sqlxMock.Mock.ExpectQuery("INSERT INTO transactions").WillReturnError(errors.New("error"))
			},
			expectErr: true,
		},
		{
			name:        "Create transaction",
			transaction: models.TransactionInput{},
			mockSetup: func() {
				rows := sqlxmock.NewRows([]string{"id"}).AddRow(uuid.New())
				sqlxMock.Mock.ExpectQuery("INSERT INTO transactions").WillReturnRows(rows)
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			_, err := repositories.R().Create(tt.transaction)
			if (err != nil) != tt.expectErr {
				t.Errorf("Get() error new = %v, expectErr %v", err, tt.expectErr)
				return
			}
		})
	}
}

// TestPostgresRepository_Get test the Get method
func TestPostgresRepository_Get(t *testing.T) {
	var sqlxMock test.Sqlx
	sqlxMock.CreateFullTestSqlx(t)
	defer sqlxMock.CleanTestSqlx()

	repositories.ReplaceGlobals(repositories.NewPostgresRepository(sqlxMock.DB))

	tests := []struct {
		name          string
		transactionID uuid.UUID
		mockSetup     func()
		expectErr     bool
		expectFound   bool
	}{
		{
			name:          "Fail transaction retrieval",
			transactionID: uuid.New(),
			mockSetup: func() {
				sqlxMock.Mock.ExpectQuery("SELECT").WillReturnError(errors.New("error"))
			},
			expectErr:   true,
			expectFound: false,
		},
		{
			name:          "Retrieve transaction",
			transactionID: uuid.New(),
			mockSetup: func() {
				rows := sqlxmock.NewRows([]string{"broker.id", "broker.name", "broker.image_id", "id", "user_id", "date", "transaction_type", "asset", "quantity", "price", "price_unit", "fee"}).
					AddRow(uuid.New(), "broker_name", uuid.New(), uuid.New(), uuid.New(), time.Now(), "type", "asset", 0, 0.0, 0.0, 0.0)
				sqlxMock.Mock.ExpectQuery("SELECT").WillReturnRows(rows)
			},
			expectErr:   false,
			expectFound: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			_, found, err := repositories.R().Get(tt.transactionID)
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

// TestPostgresRepository_Update test the Update method
func TestPostgresRepository_Update(t *testing.T) {
	var sqlxMock test.Sqlx
	sqlxMock.CreateFullTestSqlx(t)
	defer sqlxMock.CleanTestSqlx()

	repositories.ReplaceGlobals(repositories.NewPostgresRepository(sqlxMock.DB))

	tests := []struct {
		name        string
		transaction models.TransactionInput
		mockSetup   func()
		expectErr   bool
	}{
		{
			name: "Fail transaction update",
			mockSetup: func() {
				sqlxMock.Mock.ExpectExec("UPDATE transactions").WillReturnError(errors.New("error"))
			},
			expectErr: true,
		},
		{
			name: "Update transaction",
			mockSetup: func() {
				sqlxMock.Mock.ExpectExec("UPDATE transactions").WillReturnResult(sqlxmock.NewResult(1, 1))
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			err := repositories.R().Update(models.TransactionInput{})
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

	repositories.ReplaceGlobals(repositories.NewPostgresRepository(sqlxMock.DB))

	tests := []struct {
		name        string
		transaction models.Transaction
		mockSetup   func()
		expectErr   bool
	}{
		{
			name: "Fail transaction delete",
			mockSetup: func() {
				sqlxMock.Mock.ExpectExec("DELETE FROM transactions").WillReturnError(errors.New("error"))
			},
			expectErr: true,
		},
		{
			name: "Delete transaction",
			mockSetup: func() {
				sqlxMock.Mock.ExpectExec("DELETE FROM transactions").WillReturnResult(sqlxmock.NewResult(1, 1))
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			err := repositories.R().Delete(models.Transaction{})
			if (err != nil) != tt.expectErr {
				t.Errorf("Delete() error = %v, expectErr %v", err, tt.expectErr)
			}
		})
	}
}

// TestPostgresRepository_DeleteByBroker test the DeleteByBroker method
func TestPostgresRepository_DeleteByBroker(t *testing.T) {
	var sqlxMock test.Sqlx
	sqlxMock.CreateFullTestSqlx(t)
	defer sqlxMock.CleanTestSqlx()

	repositories.ReplaceGlobals(repositories.NewPostgresRepository(sqlxMock.DB))

	tests := []struct {
		name        string
		transaction models.Transaction
		mockSetup   func()
		expectErr   bool
	}{
		{
			name: "Fail transaction delete",
			mockSetup: func() {
				sqlxMock.Mock.ExpectExec("DELETE FROM transactions").WillReturnError(errors.New("error"))
			},
			expectErr: true,
		},
		{
			name: "Delete transaction",
			mockSetup: func() {
				sqlxMock.Mock.ExpectExec("DELETE FROM transactions").WillReturnResult(sqlxmock.NewResult(1, 10))
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			err := repositories.R().DeleteByBroker(models.Transaction{})
			if (err != nil) != tt.expectErr {
				t.Errorf("Delete() error = %v, expectErr %v", err, tt.expectErr)
			}
		})
	}
}

// TestPostgresRepository_Exists test the Exists method
func TestPostgresRepository_Exists(t *testing.T) {
	var sqlxMock test.Sqlx
	sqlxMock.CreateFullTestSqlx(t)
	defer sqlxMock.CleanTestSqlx()

	repositories.ReplaceGlobals(repositories.NewPostgresRepository(sqlxMock.DB))

	tests := []struct {
		name          string
		transactionID uuid.UUID
		userID        uuid.UUID
		mockSetup     func()
		expectErr     bool
		expectExists  bool
	}{
		{
			name:          "Fail transaction exists check",
			transactionID: uuid.New(),
			userID:        uuid.New(),
			mockSetup: func() {
				sqlxMock.Mock.ExpectQuery("SELECT").WillReturnError(errors.New("error"))
			},
			expectErr:    true,
			expectExists: false,
		},
		{
			name:          "Transaction exists",
			transactionID: uuid.New(),
			userID:        uuid.New(),
			mockSetup: func() {
				rows := sqlxmock.NewRows([]string{"id"}).AddRow(1)
				sqlxMock.Mock.ExpectQuery("SELECT").WillReturnRows(rows)
			},
			expectErr:    false,
			expectExists: true,
		},
		{
			name:          "Transaction does not exist",
			transactionID: uuid.New(),
			userID:        uuid.New(),
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
			exists, err := repositories.R().Exists(tt.transactionID, tt.userID)
			if (err != nil) != tt.expectErr {
				t.Errorf("Exists() error = %v, expectErr %v", err, tt.expectErr)
			}
			if exists != tt.expectExists {
				t.Errorf("Exists() exists = %v, expectExists %v", exists, tt.expectExists)
			}
		})
	}
}

// TestPostgresRepository_GetAll test the GetAll method
func TestPostgresRepository_GetAll(t *testing.T) {
	var sqlxMock test.Sqlx
	sqlxMock.CreateFullTestSqlx(t)
	defer sqlxMock.CleanTestSqlx()

	repositories.ReplaceGlobals(repositories.NewPostgresRepository(sqlxMock.DB))

	tests := []struct {
		name        string
		userID      uuid.UUID
		mockSetup   func()
		expectErr   bool
		expectCount int
	}{
		{
			name:   "Fail transaction retrieval",
			userID: uuid.New(),
			mockSetup: func() {
				sqlxMock.Mock.ExpectQuery("SELECT").WillReturnError(errors.New("error"))
			},
			expectErr:   true,
			expectCount: 0,
		},
		{
			name:   "Retrieve transactions",
			userID: uuid.New(),
			mockSetup: func() {
				rows := sqlxmock.NewRows([]string{"broker.id", "broker.name", "broker.image_id", "id", "user_id", "date", "transaction_type", "asset", "quantity", "price", "price_unit", "fee"}).
					AddRow(uuid.New(), "broker_name", uuid.New(), uuid.New(), uuid.New(), time.Now(), "type", "asset", 0, 0.0, 0.0, 0.0).
					AddRow(uuid.New(), "broker_name", uuid.New(), uuid.New(), uuid.New(), time.Now(), "type", "asset", 0, 0.0, 0.0, 0.0)
				sqlxMock.Mock.ExpectQuery("SELECT").WillReturnRows(rows)
			},
			expectErr:   false,
			expectCount: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			transactions, err := repositories.R().GetAll(tt.userID)
			if (err != nil) != tt.expectErr {
				t.Errorf("GetAll() error = %v, expectErr %v", err, tt.expectErr)
				return
			}
			if len(transactions) != tt.expectCount {
				t.Errorf("GetAll() count = %v, expectCount %v", len(transactions), tt.expectCount)
			}
		})
	}
}
