package app

import (
	"github.com/Zapharaos/fihub-backend/test"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"testing"
)

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
