package app

/*func TestInit(t *testing.T) {
	// Create a full test suite
	ts := test.TestSuite{}
	_ = ts.CreateFullTestSuite(t)
	defer ts.CleanTestSuite(t)

	// Call Init function
	Init()

	// Assertions to verify initialization
	assert.NotNil(t, email.S(), "Email service should be initialized")
	assert.NotNil(t, translation.S(), "Translation service should be initialized")
}*/

/*func TestInitLogger(t *testing.T) {
	// Create a full test suite
	ts := test.TestSuite{}
	_ = ts.CreateFullTestSuite(t)
	defer ts.CleanTestSuite(t)

	// Call initLogger function
	zapConfig := initLogger()

	// Assertions to verify logger configuration
	assert.NotNil(t, zap.L(), "Logger should be initialized")
	assert.Equal(t, zapcore.ISO8601TimeEncoder, zapConfig.EncoderConfig.EncodeTime, "Logger time encoder should be ISO8601")
}

func TestInitPostgres(t *testing.T) {
	// Create a full test suite
	ts := test.TestSuite{}
	_ = ts.CreateFullTestSuite(t)
	defer ts.CleanTestSuite(t)

	// Call initPostgres function
	initPostgres()

	// Assertions to verify Postgres initialization
	assert.NotNil(t, database.DB(), "Postgres client should be initialized")
}

func TestInitRepositories(t *testing.T) {
	// Create a full test suite
	ts := test.TestSuite{}
	_ = ts.CreateFullTestSuite(t)
	defer ts.CleanTestSuite(t)

	// Call initRepositories function
	initRepositories()

	// Assertions to verify repositories initialization
	assert.NotNil(t, users.R(), "Users repository should be initialized")
	assert.NotNil(t, password.R(), "Password repository should be initialized")
	assert.NotNil(t, roles.R(), "Roles repository should be initialized")
	assert.NotNil(t, permissions.R(), "Permissions repository should be initialized")
	assert.NotNil(t, brokers.R(), "Brokers repository should be initialized")
	assert.NotNil(t, transactions.R(), "Transactions repository should be initialized")
}
*/
