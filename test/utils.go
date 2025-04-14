package test

import (
	"os"
	"testing"
)

// checkErrDeleteFolder Check if an error is not nil, delete the folder and fail the test
func checkErrDeleteFolder(t *testing.T, err error, dname string) {
	if err == nil {
		return
	}
	t.Error(err)
	err = os.RemoveAll(dname)
	if err != nil {
		t.Error(err)
	}
	t.FailNow()
}

type TestSuite struct {
	dname   string // temporary directory name (and new working directory during the test)
	oldWd   string // old working directory (to go back to it after the test)
	created bool   // whether the temporary directory has been created
}

// CreateFullTestSuite Create a full test suite
// Please clean test suite after use (defer ts.CleanTestSuite())
func (ts *TestSuite) CreateFullTestSuite(t *testing.T) (directoryPath string) {
	// Create temporary directory
	dname, err := os.MkdirTemp("", "fihub-backend")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	ts.created = true
	ts.dname = dname

	// Save the old working directory
	oldWd, err := os.Getwd()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	ts.oldWd = oldWd

	// Chdir to a temporary directory
	err = os.Chdir(dname)
	checkErrDeleteFolder(t, err, dname)

	// .env file
	file, err := os.Create(".env")
	checkErrDeleteFolder(t, err, dname)
	err = file.Close()
	checkErrDeleteFolder(t, err, dname)

	// Setup config directory
	ts.createConfigFile(t)

	return dname
}

func (ts *TestSuite) createConfigFile(t *testing.T) {
	// Create the config directory
	err := os.MkdirAll("config", os.ModePerm)
	checkErrDeleteFolder(t, err, ts.dname)

	// Create the config file
	file, err := os.Create("config/fihub-test.toml")
	checkErrDeleteFolder(t, err, ts.dname)

	// Write the config content
	_, err = file.WriteString(`
		APP_ENV = "test"
	`)
	checkErrDeleteFolder(t, err, ts.dname)

	// Close the file
	err = file.Close()
	checkErrDeleteFolder(t, err, ts.dname)
}

// CreateConfigTranslationsFullTestSuite Create a full test suite with a config/translations directory
func (ts *TestSuite) CreateConfigTranslationsFullTestSuite(t *testing.T) string {
	// Create temporary directory
	dname := ts.CreateFullTestSuite(t)

	// Create config translations directory
	err := os.MkdirAll("config/translations", os.ModePerm)
	checkErrDeleteFolder(t, err, dname)

	// Create the default english translation file
	file, err := os.Create("config/translations/active.en.toml")
	checkErrDeleteFolder(t, err, dname)

	// Write the translation content
	_, err = file.WriteString(`
		[hello]
		other = "Hello, {{.name}}!"
	`)
	checkErrDeleteFolder(t, err, dname)

	// Close the file
	err = file.Close()
	checkErrDeleteFolder(t, err, dname)

	// French translation file
	file, err = os.Create("config/translations/active.fr.toml")
	checkErrDeleteFolder(t, err, dname)
	err = file.Close()
	checkErrDeleteFolder(t, err, dname)

	return dname
}

// CleanTestSuite Clean a test suite
func (ts *TestSuite) CleanTestSuite(t *testing.T) {
	// go back to the old working directory
	err := os.Chdir(ts.oldWd)
	if err != nil {
		t.Error(err)
	}

	// remove the temporary directory
	err = os.RemoveAll(ts.dname)
	if err != nil {
		t.Error(err)
	}
}
