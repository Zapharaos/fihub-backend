package repositories_test

import (
	"errors"
	"github.com/Zapharaos/fihub-backend/cmd/broker/app/repositories"
	"github.com/Zapharaos/fihub-backend/internal/models"
	"github.com/Zapharaos/fihub-backend/test"
	"github.com/google/uuid"
	sqlxmock "github.com/zhashkevych/go-sqlxmock"
	"testing"
)

// TestImagePostgresRepository_Create tests the Create method
func TestImagePostgresRepository_Create(t *testing.T) {
	var sqlxMock test.Sqlx
	sqlxMock.CreateFullTestSqlx(t)
	defer sqlxMock.CleanTestSqlx()

	repositories.ReplaceGlobals(repositories.NewRepository(nil, nil, repositories.NewImagePostgresRepository(sqlxMock.DB)))

	tests := []struct {
		name      string
		mockSetup func()
		expectErr bool
	}{
		{
			name: "Fail image creation",
			mockSetup: func() {
				sqlxMock.Mock.ExpectQuery("INSERT INTO broker_image").WillReturnError(errors.New("error"))
			},
			expectErr: true,
		},
		{
			name: "Create image",
			mockSetup: func() {
				rows := sqlxmock.NewRows([]string{"id", "broker_id", "name", "data"})
				sqlxMock.Mock.ExpectQuery("INSERT INTO broker_image").WillReturnRows(rows)
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			err := repositories.R().I().Create(models.BrokerImage{})
			if (err != nil) != tt.expectErr {
				t.Errorf("Create() error = %v, expectErr %v", err, tt.expectErr)
			}
		})
	}
}

// TestImagePostgresRepository_Get tests the Get method
func TestImagePostgresRepository_Get(t *testing.T) {
	var sqlxMock test.Sqlx
	sqlxMock.CreateFullTestSqlx(t)
	defer sqlxMock.CleanTestSqlx()

	repositories.ReplaceGlobals(repositories.NewRepository(nil, nil, repositories.NewImagePostgresRepository(sqlxMock.DB)))

	tests := []struct {
		name        string
		mockSetup   func()
		expectErr   bool
		expectFound bool
	}{
		{
			name: "Fail image retrieval",
			mockSetup: func() {
				sqlxMock.Mock.ExpectQuery("SELECT").WillReturnError(errors.New("error"))
			},
			expectErr:   true,
			expectFound: false,
		},
		{
			name: "Retrieve image",
			mockSetup: func() {
				rows := sqlxmock.NewRows([]string{"id", "broker_id", "name", "data"}).
					AddRow(uuid.New(), uuid.New(), "test", []byte("data"))
				sqlxMock.Mock.ExpectQuery("SELECT").WillReturnRows(rows)
			},
			expectErr:   false,
			expectFound: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			_, found, err := repositories.R().I().Get(uuid.New())
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

// TestImagePostgresRepository_Update tests the Update method
func TestImagePostgresRepository_Update(t *testing.T) {
	var sqlxMock test.Sqlx
	sqlxMock.CreateFullTestSqlx(t)
	defer sqlxMock.CleanTestSqlx()

	repositories.ReplaceGlobals(repositories.NewRepository(nil, nil, repositories.NewImagePostgresRepository(sqlxMock.DB)))

	tests := []struct {
		name      string
		mockSetup func()
		expectErr bool
	}{
		{
			name: "Fail image update",
			mockSetup: func() {
				sqlxMock.Mock.ExpectExec("UPDATE broker_image").WillReturnError(errors.New("error"))
			},
			expectErr: true,
		},
		{
			name: "Update image",
			mockSetup: func() {
				sqlxMock.Mock.ExpectExec("UPDATE broker_image").WillReturnResult(sqlxmock.NewResult(1, 1))
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			err := repositories.R().I().Update(models.BrokerImage{})
			if (err != nil) != tt.expectErr {
				t.Errorf("Update() error = %v, expectErr %v", err, tt.expectErr)
			}
		})
	}
}

// TestImagePostgresRepository_Delete tests the Delete method
func TestImagePostgresRepository_Delete(t *testing.T) {
	var sqlxMock test.Sqlx
	sqlxMock.CreateFullTestSqlx(t)
	defer sqlxMock.CleanTestSqlx()

	repositories.ReplaceGlobals(repositories.NewRepository(nil, nil, repositories.NewImagePostgresRepository(sqlxMock.DB)))

	tests := []struct {
		name      string
		mockSetup func()
		expectErr bool
	}{
		{
			name: "Fail image delete",
			mockSetup: func() {
				sqlxMock.Mock.ExpectExec("DELETE FROM broker_image").WillReturnError(errors.New("error"))
			},
			expectErr: true,
		},
		{
			name: "Delete image",
			mockSetup: func() {
				sqlxMock.Mock.ExpectExec("DELETE FROM broker_image").WillReturnResult(sqlxmock.NewResult(1, 1))
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			err := repositories.R().I().Delete(uuid.New())
			if (err != nil) != tt.expectErr {
				t.Errorf("Delete() error = %v, expectErr %v", err, tt.expectErr)
			}
		})
	}
}

// TestImagePostgresRepository_Exists tests the Exists method
func TestImagePostgresRepository_Exists(t *testing.T) {
	var sqlxMock test.Sqlx
	sqlxMock.CreateFullTestSqlx(t)
	defer sqlxMock.CleanTestSqlx()

	repositories.ReplaceGlobals(repositories.NewRepository(nil, nil, repositories.NewImagePostgresRepository(sqlxMock.DB)))

	tests := []struct {
		name         string
		brokerID     uuid.UUID
		imageID      uuid.UUID
		mockSetup    func()
		expectErr    bool
		expectExists bool
	}{
		{
			name:     "Fail image exists check",
			brokerID: uuid.New(),
			imageID:  uuid.New(),
			mockSetup: func() {
				sqlxMock.Mock.ExpectQuery("SELECT").WillReturnError(errors.New("error"))
			},
			expectErr:    true,
			expectExists: false,
		},
		{
			name:     "BrokerImage exists",
			brokerID: uuid.New(),
			imageID:  uuid.New(),
			mockSetup: func() {
				rows := sqlxmock.NewRows([]string{"id"}).AddRow(1)
				sqlxMock.Mock.ExpectQuery("SELECT").WillReturnRows(rows)
			},
			expectErr:    false,
			expectExists: true,
		},
		{
			name:     "BrokerImage does not exist",
			brokerID: uuid.New(),
			imageID:  uuid.New(),
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
			exists, err := repositories.R().I().Exists(tt.brokerID, tt.imageID)
			if (err != nil) != tt.expectErr {
				t.Errorf("Exists() error = %v, expectErr %v", err, tt.expectErr)
			}
			if exists != tt.expectExists {
				t.Errorf("Exists() exists = %v, expectExists %v", exists, tt.expectExists)
			}
		})
	}
}
