package brokers_test

import (
	"errors"
	"github.com/Zapharaos/fihub-backend/internal/brokers"
	"github.com/Zapharaos/fihub-backend/internal/models"
	"github.com/Zapharaos/fihub-backend/test"
	"github.com/google/uuid"
	sqlxmock "github.com/zhashkevych/go-sqlxmock"
	"testing"
)

// TestUserPostgresRepository_Create tests the Create method
func TestUserPostgresRepository_Create(t *testing.T) {
	var sqlxMock test.Sqlx
	sqlxMock.CreateFullTestSqlx(t)
	defer sqlxMock.CleanTestSqlx()

	brokers.ReplaceGlobals(brokers.NewRepository(nil, brokers.NewUserPostgresRepository(sqlxMock.DB), nil))

	tests := []struct {
		name      string
		mockSetup func()
		expectErr bool
	}{
		{
			name: "Fail user broker creation",
			mockSetup: func() {
				sqlxMock.Mock.ExpectQuery("INSERT INTO user_brokers").WillReturnError(errors.New("error"))
			},
			expectErr: true,
		},
		{
			name: "Create user broker",
			mockSetup: func() {

				rows := sqlxmock.NewRows([]string{"user_id", "broker_id"})
				sqlxMock.Mock.ExpectQuery("INSERT INTO user_brokers").WillReturnRows(rows)
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			err := brokers.R().U().Create(models.BrokerUser{})
			if (err != nil) != tt.expectErr {
				t.Errorf("Create() error = %v, expectErr %v", err, tt.expectErr)
			}
		})
	}
}

// TestUserPostgresRepository_Delete tests the Delete method
func TestUserPostgresRepository_Delete(t *testing.T) {
	var sqlxMock test.Sqlx
	sqlxMock.CreateFullTestSqlx(t)
	defer sqlxMock.CleanTestSqlx()

	brokers.ReplaceGlobals(brokers.NewRepository(nil, brokers.NewUserPostgresRepository(sqlxMock.DB), nil))

	tests := []struct {
		name      string
		mockSetup func()
		expectErr bool
	}{
		{
			name: "Fail user broker delete",
			mockSetup: func() {
				sqlxMock.Mock.ExpectExec("DELETE FROM user_brokers").WillReturnError(errors.New("error"))
			},
			expectErr: true,
		},
		{
			name: "Delete user broker",
			mockSetup: func() {
				sqlxMock.Mock.ExpectExec("DELETE FROM user_brokers").WillReturnResult(sqlxmock.NewResult(1, 1))
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			err := brokers.R().U().Delete(models.BrokerUser{})
			if (err != nil) != tt.expectErr {
				t.Errorf("Delete() error = %v, expectErr %v", err, tt.expectErr)
			}
		})
	}
}

// TestUserPostgresRepository_Exists tests the Exists method
func TestUserPostgresRepository_Exists(t *testing.T) {
	var sqlxMock test.Sqlx
	sqlxMock.CreateFullTestSqlx(t)
	defer sqlxMock.CleanTestSqlx()

	brokers.ReplaceGlobals(brokers.NewRepository(nil, brokers.NewUserPostgresRepository(sqlxMock.DB), nil))

	tests := []struct {
		name         string
		user         models.BrokerUser
		mockSetup    func()
		expectErr    bool
		expectExists bool
	}{
		{
			name: "Fail user broker exists check",
			mockSetup: func() {
				sqlxMock.Mock.ExpectQuery("SELECT").WillReturnError(errors.New("error"))
			},
			expectErr:    true,
			expectExists: false,
		},
		{
			name: "BrokerUser broker exists",
			mockSetup: func() {
				rows := sqlxmock.NewRows([]string{"user_id", "broker_id"}).AddRow(uuid.New(), uuid.New())
				sqlxMock.Mock.ExpectQuery("SELECT").WillReturnRows(rows)
			},
			expectErr:    false,
			expectExists: true,
		},
		{
			name: "BrokerUser broker does not exist",
			mockSetup: func() {
				rows := sqlxmock.NewRows([]string{"user_id", "broker_id"})
				sqlxMock.Mock.ExpectQuery("SELECT").WillReturnRows(rows)
			},
			expectErr:    false,
			expectExists: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			exists, err := brokers.R().U().Exists(models.BrokerUser{})
			if (err != nil) != tt.expectErr {
				t.Errorf("Exists() error = %v, expectErr %v", err, tt.expectErr)
			}
			if exists != tt.expectExists {
				t.Errorf("Exists() exists = %v, expectExists %v", exists, tt.expectExists)
			}
		})
	}
}

// TestUserPostgresRepository_GetAll tests the GetAll method
func TestUserPostgresRepository_GetAll(t *testing.T) {
	var sqlxMock test.Sqlx
	sqlxMock.CreateFullTestSqlx(t)
	defer sqlxMock.CleanTestSqlx()

	brokers.ReplaceGlobals(brokers.NewRepository(nil, brokers.NewUserPostgresRepository(sqlxMock.DB), nil))

	tests := []struct {
		name        string
		userID      uuid.UUID
		mockSetup   func()
		expectErr   bool
		expectCount int
	}{
		{
			name:   "Fail user broker retrieval",
			userID: uuid.New(),
			mockSetup: func() {
				sqlxMock.Mock.ExpectQuery("SELECT").WillReturnError(errors.New("error"))
			},
			expectErr:   true,
			expectCount: 0,
		},
		{
			name:   "Retrieve user brokers",
			userID: uuid.New(),
			mockSetup: func() {
				rows := sqlxmock.NewRows([]string{"id", "name", "image_id"}).
					AddRow(uuid.New(), "broker_name", uuid.New()).
					AddRow(uuid.New(), "broker_name", uuid.New())
				sqlxMock.Mock.ExpectQuery("SELECT").WillReturnRows(rows)
			},
			expectErr:   false,
			expectCount: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			userBrokers, err := brokers.R().U().GetAll(tt.userID)
			if (err != nil) != tt.expectErr {
				t.Errorf("GetAll() error = %v, expectErr %v", err, tt.expectErr)
				return
			}
			if len(userBrokers) != tt.expectCount {
				t.Errorf("GetAll() count = %v, expectCount %v", len(userBrokers), tt.expectCount)
			}
		})
	}
}
