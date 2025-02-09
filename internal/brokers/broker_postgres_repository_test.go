package brokers_test

import (
	"errors"
	"github.com/Zapharaos/fihub-backend/internal/brokers"
	"github.com/Zapharaos/fihub-backend/test"
	"github.com/google/uuid"
	sqlxmock "github.com/zhashkevych/go-sqlxmock"
	"testing"
)

// TestBrokerPostgresRepository_Create tests the Create method
func TestBrokerPostgresRepository_Create(t *testing.T) {
	var sqlxMock test.Sqlx
	sqlxMock.CreateFullTestSqlx(t)
	defer sqlxMock.CleanTestSqlx()

	brokers.ReplaceGlobals(brokers.NewRepository(brokers.NewPostgresRepository(sqlxMock.DB), nil, nil))

	tests := []struct {
		name      string
		mockSetup func()
		expectErr bool
	}{
		{
			name: "Fail broker creation",
			mockSetup: func() {
				sqlxMock.Mock.ExpectQuery("INSERT INTO brokers").WillReturnError(errors.New("error"))
			},
			expectErr: true,
		},
		{
			name: "Create broker",
			mockSetup: func() {
				rows := sqlxmock.NewRows([]string{"id"}).AddRow(uuid.New())
				sqlxMock.Mock.ExpectQuery("INSERT INTO brokers").WillReturnRows(rows)
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			_, err := brokers.R().B().Create(brokers.Broker{})
			if (err != nil) != tt.expectErr {
				t.Errorf("Create() error = %v, expectErr %v", err, tt.expectErr)
			}
		})
	}
}

// TestBrokerPostgresRepository_Get tests the Get method
func TestBrokerPostgresRepository_Get(t *testing.T) {
	var sqlxMock test.Sqlx
	sqlxMock.CreateFullTestSqlx(t)
	defer sqlxMock.CleanTestSqlx()

	brokers.ReplaceGlobals(brokers.NewRepository(brokers.NewPostgresRepository(sqlxMock.DB), nil, nil))

	tests := []struct {
		name        string
		id          uuid.UUID
		mockSetup   func()
		expectErr   bool
		expectFound bool
	}{
		{
			name: "Fail broker retrieval",
			id:   uuid.New(),
			mockSetup: func() {
				sqlxMock.Mock.ExpectQuery("SELECT").WillReturnError(errors.New("error"))
			},
			expectErr:   true,
			expectFound: false,
		},
		{
			name: "Broker not found",
			id:   uuid.New(),
			mockSetup: func() {
				rows := sqlxmock.NewRows([]string{"id", "name", "image_id", "disabled"})
				sqlxMock.Mock.ExpectQuery("SELECT").WillReturnRows(rows)
			},
			expectErr:   false,
			expectFound: false,
		},
		{
			name: "Broker found",
			id:   uuid.New(),
			mockSetup: func() {
				rows := sqlxmock.NewRows([]string{"id", "name", "image_id", "disabled"}).
					AddRow(uuid.New(), "broker_name", uuid.New(), false)
				sqlxMock.Mock.ExpectQuery("SELECT").WillReturnRows(rows)
			},
			expectErr:   false,
			expectFound: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			_, found, err := brokers.R().B().Get(tt.id)
			if (err != nil) != tt.expectErr {
				t.Errorf("Get() error = %v, expectErr %v", err, tt.expectErr)
			}
			if found != tt.expectFound {
				t.Errorf("Get() found = %v, expectFound %v", found, tt.expectFound)
			}
		})
	}
}

// TestBrokerPostgresRepository_Update tests the Update method
func TestBrokerPostgresRepository_Update(t *testing.T) {
	var sqlxMock test.Sqlx
	sqlxMock.CreateFullTestSqlx(t)
	defer sqlxMock.CleanTestSqlx()

	brokers.ReplaceGlobals(brokers.NewRepository(brokers.NewPostgresRepository(sqlxMock.DB), nil, nil))

	tests := []struct {
		name      string
		broker    brokers.Broker
		mockSetup func()
		expectErr bool
	}{
		{
			name: "Fail broker update",
			broker: brokers.Broker{
				ID:       uuid.New(),
				Name:     "broker_name",
				Disabled: false,
			},
			mockSetup: func() {
				sqlxMock.Mock.ExpectExec("UPDATE brokers").WillReturnError(errors.New("error"))
			},
			expectErr: true,
		},
		{
			name: "Update broker",
			broker: brokers.Broker{
				ID:       uuid.New(),
				Name:     "broker_name",
				Disabled: false,
			},
			mockSetup: func() {
				sqlxMock.Mock.ExpectExec("UPDATE brokers").WillReturnResult(sqlxmock.NewResult(1, 1))
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			err := brokers.R().B().Update(tt.broker)
			if (err != nil) != tt.expectErr {
				t.Errorf("Update() error = %v, expectErr %v", err, tt.expectErr)
			}
		})
	}
}

// TestBrokerPostgresRepository_Delete tests the Delete method
func TestBrokerPostgresRepository_Delete(t *testing.T) {
	var sqlxMock test.Sqlx
	sqlxMock.CreateFullTestSqlx(t)
	defer sqlxMock.CleanTestSqlx()

	brokers.ReplaceGlobals(brokers.NewRepository(brokers.NewPostgresRepository(sqlxMock.DB), nil, nil))

	tests := []struct {
		name      string
		id        uuid.UUID
		mockSetup func()
		expectErr bool
	}{
		{
			name: "Fail broker deletion",
			id:   uuid.New(),
			mockSetup: func() {
				sqlxMock.Mock.ExpectExec("DELETE FROM brokers").WillReturnError(errors.New("error"))
			},
			expectErr: true,
		},
		{
			name: "Delete broker",
			id:   uuid.New(),
			mockSetup: func() {
				sqlxMock.Mock.ExpectExec("DELETE FROM brokers").WillReturnResult(sqlxmock.NewResult(1, 1))
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			err := brokers.R().B().Delete(tt.id)
			if (err != nil) != tt.expectErr {
				t.Errorf("Delete() error = %v, expectErr %v", err, tt.expectErr)
			}
		})
	}
}

// TestBrokerPostgresRepository_Exists tests the Exists method
func TestBrokerPostgresRepository_Exists(t *testing.T) {
	var sqlxMock test.Sqlx
	sqlxMock.CreateFullTestSqlx(t)
	defer sqlxMock.CleanTestSqlx()

	brokers.ReplaceGlobals(brokers.NewRepository(brokers.NewPostgresRepository(sqlxMock.DB), nil, nil))

	tests := []struct {
		name        string
		id          uuid.UUID
		mockSetup   func()
		expectErr   bool
		expectExist bool
	}{
		{
			name: "Fail broker exists check",
			id:   uuid.New(),
			mockSetup: func() {
				sqlxMock.Mock.ExpectQuery("SELECT").WillReturnError(errors.New("error"))
			},
			expectErr:   true,
			expectExist: false,
		},
		{
			name: "Broker does not exist",
			id:   uuid.New(),
			mockSetup: func() {
				rows := sqlxmock.NewRows([]string{"id", "name", "image_id", "disabled"})
				sqlxMock.Mock.ExpectQuery("SELECT").WillReturnRows(rows)
			},
			expectErr:   false,
			expectExist: false,
		},
		{
			name: "Broker exists",
			id:   uuid.New(),
			mockSetup: func() {
				rows := sqlxmock.NewRows([]string{"id", "name", "image_id", "disabled"}).
					AddRow(uuid.New(), "broker_name", uuid.New(), false)
				sqlxMock.Mock.ExpectQuery("SELECT").WillReturnRows(rows)
			},
			expectErr:   false,
			expectExist: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			exists, err := brokers.R().B().Exists(tt.id)
			if (err != nil) != tt.expectErr {
				t.Errorf("Exists() error = %v, expectErr %v", err, tt.expectErr)
			}
			if exists != tt.expectExist {
				t.Errorf("Exists() exists = %v, expectExist %v", exists, tt.expectExist)
			}
		})
	}
}

// TestBrokerPostgresRepository_ExistsByName tests the ExistsByName method
func TestBrokerPostgresRepository_ExistsByName(t *testing.T) {
	var sqlxMock test.Sqlx
	sqlxMock.CreateFullTestSqlx(t)
	defer sqlxMock.CleanTestSqlx()

	brokers.ReplaceGlobals(brokers.NewRepository(brokers.NewPostgresRepository(sqlxMock.DB), nil, nil))

	tests := []struct {
		name        string
		brokerName  string
		mockSetup   func()
		expectErr   bool
		expectExist bool
	}{
		{
			name:       "Fail broker exists by name check",
			brokerName: "broker_name",
			mockSetup: func() {
				sqlxMock.Mock.ExpectQuery("SELECT").WillReturnError(errors.New("error"))
			},
			expectErr:   true,
			expectExist: false,
		},
		{
			name:       "Broker does not exist by name",
			brokerName: "broker_name",
			mockSetup: func() {
				rows := sqlxmock.NewRows([]string{"id", "name", "image_id", "disabled"})
				sqlxMock.Mock.ExpectQuery("SELECT").WillReturnRows(rows)
			},
			expectErr:   false,
			expectExist: false,
		},
		{
			name:       "Broker exists by name",
			brokerName: "broker_name",
			mockSetup: func() {
				rows := sqlxmock.NewRows([]string{"id", "name", "image_id", "disabled"}).
					AddRow(uuid.New(), "broker_name", uuid.New(), false)
				sqlxMock.Mock.ExpectQuery("SELECT").WillReturnRows(rows)
			},
			expectErr:   false,
			expectExist: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			exists, err := brokers.R().B().ExistsByName(tt.brokerName)
			if (err != nil) != tt.expectErr {
				t.Errorf("ExistsByName() error = %v, expectErr %v", err, tt.expectErr)
			}
			if exists != tt.expectExist {
				t.Errorf("ExistsByName() exists = %v, expectExist %v", exists, tt.expectExist)
			}
		})
	}
}

// TestBrokerPostgresRepository_GetAll tests the GetAll method
func TestBrokerPostgresRepository_GetAll(t *testing.T) {
	var sqlxMock test.Sqlx
	sqlxMock.CreateFullTestSqlx(t)
	defer sqlxMock.CleanTestSqlx()

	brokers.ReplaceGlobals(brokers.NewRepository(brokers.NewPostgresRepository(sqlxMock.DB), nil, nil))

	tests := []struct {
		name      string
		mockSetup func()
		expectErr bool
		expectLen int
	}{
		{
			name: "Fail broker retrieval",
			mockSetup: func() {
				sqlxMock.Mock.ExpectQuery("SELECT").WillReturnError(errors.New("error"))
			},
			expectErr: true,
			expectLen: 0,
		},
		{
			name: "No brokers found",
			mockSetup: func() {
				rows := sqlxmock.NewRows([]string{"id", "name", "image_id", "disabled"})
				sqlxMock.Mock.ExpectQuery("SELECT").WillReturnRows(rows)
			},
			expectErr: false,
			expectLen: 0,
		},
		{
			name: "Brokers found",
			mockSetup: func() {
				rows := sqlxmock.NewRows([]string{"id", "name", "image_id", "disabled"}).
					AddRow(uuid.New(), "broker_name1", uuid.New(), false).
					AddRow(uuid.New(), "broker_name2", uuid.New(), true)
				sqlxMock.Mock.ExpectQuery("SELECT").WillReturnRows(rows)
			},
			expectErr: false,
			expectLen: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			b, err := brokers.R().B().GetAll()
			if (err != nil) != tt.expectErr {
				t.Errorf("GetAll() error = %v, expectErr %v", err, tt.expectErr)
			}
			if len(b) != tt.expectLen {
				t.Errorf("GetAll() len = %v, expectLen %v", len(b), tt.expectLen)
			}
		})
	}
}

// TestBrokerPostgresRepository_GetAllEnabled tests the GetAllEnabled method
func TestBrokerPostgresRepository_GetAllEnabled(t *testing.T) {
	var sqlxMock test.Sqlx
	sqlxMock.CreateFullTestSqlx(t)
	defer sqlxMock.CleanTestSqlx()

	brokers.ReplaceGlobals(brokers.NewRepository(brokers.NewPostgresRepository(sqlxMock.DB), nil, nil))

	tests := []struct {
		name      string
		mockSetup func()
		expectErr bool
		expectLen int
	}{
		{
			name: "Fail enabled broker retrieval",
			mockSetup: func() {
				sqlxMock.Mock.ExpectQuery("SELECT").WillReturnError(errors.New("error"))
			},
			expectErr: true,
			expectLen: 0,
		},
		{
			name: "No enabled brokers found",
			mockSetup: func() {
				rows := sqlxmock.NewRows([]string{"id", "name", "image_id", "disabled"})
				sqlxMock.Mock.ExpectQuery("SELECT").WillReturnRows(rows)
			},
			expectErr: false,
			expectLen: 0,
		},
		{
			name: "Enabled brokers found",
			mockSetup: func() {
				rows := sqlxmock.NewRows([]string{"id", "name", "image_id", "disabled"}).
					AddRow(uuid.New(), "broker_name1", uuid.New(), false).
					AddRow(uuid.New(), "broker_name2", uuid.New(), false)
				sqlxMock.Mock.ExpectQuery("SELECT").WillReturnRows(rows)
			},
			expectErr: false,
			expectLen: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			b, err := brokers.R().B().GetAllEnabled()
			if (err != nil) != tt.expectErr {
				t.Errorf("GetAllEnabled() error = %v, expectErr %v", err, tt.expectErr)
			}
			if len(b) != tt.expectLen {
				t.Errorf("GetAllEnabled() len = %v, expectLen %v", len(b), tt.expectLen)
			}
		})
	}
}

// TestBrokerPostgresRepository_SetImage tests the SetImage method
func TestBrokerPostgresRepository_SetImage(t *testing.T) {
	var sqlxMock test.Sqlx
	sqlxMock.CreateFullTestSqlx(t)
	defer sqlxMock.CleanTestSqlx()

	brokers.ReplaceGlobals(brokers.NewRepository(brokers.NewPostgresRepository(sqlxMock.DB), nil, nil))

	tests := []struct {
		name      string
		id        uuid.UUID
		imageID   uuid.UUID
		mockSetup func()
		expectErr bool
	}{
		{
			name:    "Fail setting image",
			id:      uuid.New(),
			imageID: uuid.New(),
			mockSetup: func() {
				sqlxMock.Mock.ExpectExec("UPDATE brokers").WillReturnError(errors.New("error"))
			},
			expectErr: true,
		},
		{
			name:    "Set image",
			id:      uuid.New(),
			imageID: uuid.New(),
			mockSetup: func() {
				sqlxMock.Mock.ExpectExec("UPDATE brokers").WillReturnResult(sqlxmock.NewResult(1, 1))
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			err := brokers.R().B().SetImage(tt.id, tt.imageID)
			if (err != nil) != tt.expectErr {
				t.Errorf("SetImage() error = %v, expectErr %v", err, tt.expectErr)
			}
		})
	}
}

// TestBrokerPostgresRepository_HasImage tests the HasImage method
func TestBrokerPostgresRepository_HasImage(t *testing.T) {
	var sqlxMock test.Sqlx
	sqlxMock.CreateFullTestSqlx(t)
	defer sqlxMock.CleanTestSqlx()

	brokers.ReplaceGlobals(brokers.NewRepository(brokers.NewPostgresRepository(sqlxMock.DB), nil, nil))

	tests := []struct {
		name        string
		id          uuid.UUID
		mockSetup   func()
		expectErr   bool
		expectExist bool
	}{
		{
			name: "Fail broker has image check",
			id:   uuid.New(),
			mockSetup: func() {
				sqlxMock.Mock.ExpectQuery("SELECT").WillReturnError(errors.New("error"))
			},
			expectErr:   true,
			expectExist: false,
		},
		{
			name: "Broker does not have image",
			id:   uuid.New(),
			mockSetup: func() {
				rows := sqlxmock.NewRows([]string{"id", "name", "image_id", "disabled"})
				sqlxMock.Mock.ExpectQuery("SELECT").WillReturnRows(rows)
			},
			expectErr:   false,
			expectExist: false,
		},
		{
			name: "Broker has image",
			id:   uuid.New(),
			mockSetup: func() {
				rows := sqlxmock.NewRows([]string{"id", "name", "image_id", "disabled"}).
					AddRow(uuid.New(), "broker_name", uuid.New(), false)
				sqlxMock.Mock.ExpectQuery("SELECT").WillReturnRows(rows)
			},
			expectErr:   false,
			expectExist: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			exists, err := brokers.R().B().HasImage(tt.id)
			if (err != nil) != tt.expectErr {
				t.Errorf("HasImage() error = %v, expectErr %v", err, tt.expectErr)
			}
			if exists != tt.expectExist {
				t.Errorf("HasImage() exists = %v, expectExist %v", exists, tt.expectExist)
			}
		})
	}
}

// TestBrokerPostgresRepository_DeleteImage tests the DeleteImage method
func TestBrokerPostgresRepository_DeleteImage(t *testing.T) {
	var sqlxMock test.Sqlx
	sqlxMock.CreateFullTestSqlx(t)
	defer sqlxMock.CleanTestSqlx()

	brokers.ReplaceGlobals(brokers.NewRepository(brokers.NewPostgresRepository(sqlxMock.DB), nil, nil))

	tests := []struct {
		name      string
		id        uuid.UUID
		mockSetup func()
		expectErr bool
	}{
		{
			name: "Fail deleting image",
			id:   uuid.New(),
			mockSetup: func() {
				sqlxMock.Mock.ExpectExec("UPDATE brokers").WillReturnError(errors.New("error"))
			},
			expectErr: true,
		},
		{
			name: "Delete image",
			id:   uuid.New(),
			mockSetup: func() {
				sqlxMock.Mock.ExpectExec("UPDATE brokers").WillReturnResult(sqlxmock.NewResult(1, 1))
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			err := brokers.R().B().DeleteImage(tt.id)
			if (err != nil) != tt.expectErr {
				t.Errorf("DeleteImage() error = %v, expectErr %v", err, tt.expectErr)
			}
		})
	}
}
