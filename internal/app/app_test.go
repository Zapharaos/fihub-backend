package app

import (
	"github.com/Zapharaos/fihub-backend/internal/brokers"
	"github.com/Zapharaos/fihub-backend/internal/database"
	"github.com/Zapharaos/fihub-backend/internal/transactions"
	"github.com/Zapharaos/fihub-backend/internal/users"
	"github.com/Zapharaos/fihub-backend/internal/users/password"
	"github.com/Zapharaos/fihub-backend/internal/users/permissions"
	"github.com/Zapharaos/fihub-backend/internal/users/roles"
	"github.com/Zapharaos/fihub-backend/pkg/email"
	"github.com/Zapharaos/fihub-backend/pkg/translation"
	"github.com/Zapharaos/fihub-backend/test"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	sqlmock "github.com/zhashkevych/go-sqlxmock"
	"go.uber.org/zap"
	"testing"
)

// TestInit tests the Init function to ensure that it correctly initializes the application.
// This test only verifies : env, email, translation.
// Any further function calls within Init are tested by their respective tests.
func TestInit(t *testing.T) {
	// Create a full test suite
	ts := test.TestSuite{}
	_ = ts.CreateConfigTranslationsFullTestSuite(t)
	defer ts.CleanTestSuite(t)

	// Mock viper configuration
	viper.Set("DEFAULT_LANG", "en")
	viper.Set("POSTGRES_HOST", "localhost")
	viper.Set("POSTGRES_PORT", "5432")
	viper.Set("POSTGRES_USER", "user")
	viper.Set("POSTGRES_PASSWORD", "password")
	viper.Set("POSTGRES_DB", "testdb")

	// Call Init function
	Init()

	// Assertions to verify initialization
	assert.NotNil(t, email.S(), "Email service should be initialized")
	assert.NotNil(t, translation.S(), "Translation service should be initialized")
}

// TestInitLogger tests the initLogger function to ensure that it correctly initializes the logger.
func TestInitLogger(t *testing.T) {
	// Create a full test suite
	ts := test.TestSuite{}
	_ = ts.CreateFullTestSuite(t)
	defer ts.CleanTestSuite(t)

	// Call initLogger function
	InitLogger()

	// Assertions to verify logger configuration
	assert.NotNil(t, zap.L(), "Logger should be initialized")
}

// TestInit tests the initDatabase function to ensure that it correctly initializes the database.
// This test only verifies the database.
// Any further function calls within initDatabase are tested by their respective tests.
func TestInitDatabase(t *testing.T) {
	// Create a full test suite
	ts := test.TestSuite{}
	_ = ts.CreateFullTestSuite(t)
	defer ts.CleanTestSuite(t)

	// Call initDatabase function
	initDatabase()

	// Assertions to verify Database initialization
	assert.NotNil(t, database.DB())
}

// TestInitPostgres tests the initPostgres function to ensure that it correctly initializes the repositories.
func TestInitPostgres(t *testing.T) {
	// Simulate a successful connection
	sqlxMock, _, err := sqlmock.Newx()
	assert.NoError(t, err)
	defer sqlxMock.Close()

	// Call initPostgres function with the mock connection
	initPostgres(sqlxMock)

	// Assertions to verify repositories initialization
	assert.NotNil(t, users.R(), "Users repository should be initialized")
	assert.NotNil(t, password.R(), "Password repository should be initialized")
	assert.NotNil(t, roles.R(), "Roles repository should be initialized")
	assert.NotNil(t, permissions.R(), "Permissions repository should be initialized")
	assert.NotNil(t, brokers.R(), "Brokers repository should be initialized")
	assert.NotNil(t, transactions.R(), "Transactions repository should be initialized")
}
