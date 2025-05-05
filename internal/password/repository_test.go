package password_test

/*// TestReplaceGlobals tests the ReplaceGlobals function
// It verifies that the global repository can be replaced and restored correctly.
func TestReplaceGlobals(t *testing.T) {
	// Replace the global repository with a mocks repository
	mockRepository := &mocks.Repositor{}
	restore := password.ReplaceGlobals(mockRepository)

	// Verify that the global repository instance has been replaced
	assert.Equal(t, mockRepository, password.R())

	// Restore the global repository instance
	restore()

	// Verify that the global repository instance has been restored
	assert.NotEqual(t, mockRepository, password.R())
}

// TestRepository tests the R function
// It verifies that the global repository can be accessed correctly.
func TestRepository(t *testing.T) {
	// Replace the global repository with a mocks repository
	mockRepository := &mocks.PasswordRepository{}
	restore := password.ReplaceGlobals(mockRepository)
	defer restore()

	// Access the global repository
	assert.Equal(t, mockRepository, password.R())
}*/
