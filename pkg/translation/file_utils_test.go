package translation

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Zapharaos/fihub-backend/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"
)

// TestDiscoverTranslationFiles tests the main discovery function
func TestDiscoverTranslationFiles(t *testing.T) {
	t.Run("Nonexistent path", func(t *testing.T) {
		files, err := discoverTranslationFiles("/nonexistent/path")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "translations path does not exist")
		assert.Nil(t, files)
	})

	t.Run("Path is not a directory", func(t *testing.T) {
		// Create a temporary file instead of directory
		ts := test.TestSuite{}
		dir := ts.CreateFullTestSuite(t)
		defer ts.CleanTestSuite(t)

		// Create a file instead of directory
		filePath := filepath.Join(dir, "notadir.txt")
		file, err := os.Create(filePath)
		require.NoError(t, err)
		file.Close()

		files, err := discoverTranslationFiles(filePath)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "translations path is not a directory")
		assert.Nil(t, files)
	})

	t.Run("Empty directory", func(t *testing.T) {
		ts := test.TestSuite{}
		dir := ts.CreateFullTestSuite(t)
		defer ts.CleanTestSuite(t)

		// Create empty directory
		translationsDir := filepath.Join(dir, "translations")
		err := os.MkdirAll(translationsDir, os.ModePerm)
		require.NoError(t, err)

		files, err := discoverTranslationFiles(translationsDir)
		assert.NoError(t, err)
		assert.Empty(t, files)
	})

	t.Run("Valid translation files without prefix filter", func(t *testing.T) {
		ts := test.TestSuite{}
		dir := ts.CreateFullTestSuite(t)
		defer ts.CleanTestSuite(t)

		// Create translations directory
		translationsDir := filepath.Join(dir, "translations")
		err := os.MkdirAll(translationsDir, os.ModePerm)
		require.NoError(t, err)

		// Create valid translation files
		validFiles := []string{
			"active.en.toml",
			"active.fr.json",
			"messages.de.yaml",
			"common.es.yml",
		}

		for _, fileName := range validFiles {
			filePath := filepath.Join(translationsDir, fileName)
			file, err := os.Create(filePath)
			require.NoError(t, err)
			_, err = file.WriteString("test content")
			require.NoError(t, err)
			file.Close()
		}

		files, err := discoverTranslationFiles(translationsDir)
		assert.NoError(t, err)
		assert.Len(t, files, 4)

		// Check that all files were discovered
		fileNames := make(map[string]bool)
		for _, file := range files {
			fileNames[filepath.Base(file.path)] = true
		}
		for _, expectedFile := range validFiles {
			assert.True(t, fileNames[expectedFile], "Expected file %s not found", expectedFile)
		}
	})

	t.Run("Valid translation files with prefix filter", func(t *testing.T) {
		ts := test.TestSuite{}
		dir := ts.CreateFullTestSuite(t)
		defer ts.CleanTestSuite(t)

		// Create translations directory
		translationsDir := filepath.Join(dir, "translations")
		err := os.MkdirAll(translationsDir, os.ModePerm)
		require.NoError(t, err)

		// Create files with different prefixes
		testFiles := []string{
			"active.en.toml",
			"active.fr.json",
			"messages.en.yaml",
			"common.es.yml",
		}

		for _, fileName := range testFiles {
			filePath := filepath.Join(translationsDir, fileName)
			file, err := os.Create(filePath)
			require.NoError(t, err)
			_, err = file.WriteString("test content")
			require.NoError(t, err)
			file.Close()
		}

		// Filter for "active" prefix only
		files, err := discoverTranslationFiles(translationsDir, "active")
		assert.NoError(t, err)
		assert.Len(t, files, 2)

		// Check that only active files were found
		for _, file := range files {
			fileName := filepath.Base(file.path)
			assert.True(t, strings.HasPrefix(fileName, "active."), "File %s should have 'active' prefix", fileName)
		}
	})

	t.Run("Files with unsupported extensions are ignored", func(t *testing.T) {
		ts := test.TestSuite{}
		dir := ts.CreateFullTestSuite(t)
		defer ts.CleanTestSuite(t)

		// Create translations directory
		translationsDir := filepath.Join(dir, "translations")
		err := os.MkdirAll(translationsDir, os.ModePerm)
		require.NoError(t, err)

		// Create files with various extensions
		testFiles := []string{
			"active.en.toml",   // valid
			"active.fr.txt",    // invalid extension
			"active.de.xml",    // invalid extension
			"messages.es.json", // valid
		}

		for _, fileName := range testFiles {
			filePath := filepath.Join(translationsDir, fileName)
			file, err := os.Create(filePath)
			require.NoError(t, err)
			_, err = file.WriteString("test content")
			require.NoError(t, err)
			file.Close()
		}

		files, err := discoverTranslationFiles(translationsDir)
		assert.NoError(t, err)
		assert.Len(t, files, 2) // Only the .toml and .json files
	})

	t.Run("Invalid filename format returns error", func(t *testing.T) {
		ts := test.TestSuite{}
		dir := ts.CreateFullTestSuite(t)
		defer ts.CleanTestSuite(t)

		// Create translations directory
		translationsDir := filepath.Join(dir, "translations")
		err := os.MkdirAll(translationsDir, os.ModePerm)
		require.NoError(t, err)

		// Create files with invalid names
		invalidFiles := []string{
			"active.en.toml",              // valid (for comparison)
			"invalid-name.toml",           // missing locale
			"active..toml",                // empty locale
			"active.invalid-locale!.toml", // invalid characters in locale
		}

		for _, fileName := range invalidFiles {
			filePath := filepath.Join(translationsDir, fileName)
			file, err := os.Create(filePath)
			require.NoError(t, err)
			_, err = file.WriteString("test content")
			require.NoError(t, err)
			file.Close()
		}

		files, err := discoverTranslationFiles(translationsDir)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid translation files")
		assert.Len(t, files, 1) // Only the valid file should be included
	})

	t.Run("Large file is considered invalid", func(t *testing.T) {
		ts := test.TestSuite{}
		dir := ts.CreateFullTestSuite(t)
		defer ts.CleanTestSuite(t)

		// Create translations directory
		translationsDir := filepath.Join(dir, "translations")
		err := os.MkdirAll(translationsDir, os.ModePerm)
		require.NoError(t, err)

		// Create a file that's too large
		largeFilePath := filepath.Join(translationsDir, "large.en.toml")
		file, err := os.Create(largeFilePath)
		require.NoError(t, err)

		// Write more than maxTranslationFileSize (1MB)
		largeContent := strings.Repeat("a", maxTranslationFileSize+1)
		_, err = file.WriteString(largeContent)
		require.NoError(t, err)
		file.Close()

		_, err = discoverTranslationFiles(translationsDir)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid translation files")
	})

	t.Run("Directories are ignored", func(t *testing.T) {
		ts := test.TestSuite{}
		dir := ts.CreateFullTestSuite(t)
		defer ts.CleanTestSuite(t)

		// Create translations directory
		translationsDir := filepath.Join(dir, "translations")
		err := os.MkdirAll(translationsDir, os.ModePerm)
		require.NoError(t, err)

		// Create a subdirectory that looks like a translation file
		subDir := filepath.Join(translationsDir, "active.en.toml")
		err = os.MkdirAll(subDir, os.ModePerm)
		require.NoError(t, err)

		// Create a valid file
		validFile := filepath.Join(translationsDir, "messages.fr.json")
		file, err := os.Create(validFile)
		require.NoError(t, err)
		_, err = file.WriteString("test content")
		require.NoError(t, err)
		file.Close()

		files, err := discoverTranslationFiles(translationsDir)
		assert.NoError(t, err)
		assert.Len(t, files, 1) // Only the valid file, directory ignored
		assert.Equal(t, "messages.fr.json", filepath.Base(files[0].path))
	})
}

// TestHasSupportedExtension tests the extension validation function
func TestHasSupportedExtension(t *testing.T) {
	testCases := []struct {
		filename string
		expected bool
	}{
		{"active.en.toml", true},
		{"active.en.json", true},
		{"active.en.yaml", true},
		{"active.en.yml", true},
		{"active.en.TOML", true}, // case insensitive
		{"active.en.JSON", true}, // case insensitive
		{"active.en.txt", false},
		{"active.en.xml", false},
		{"active.en", false},
		{"", false},
		{"active.en.toml.backup", false}, // extension must be at the end
	}

	for _, tc := range testCases {
		t.Run(tc.filename, func(t *testing.T) {
			result := hasSupportedExtension(tc.filename)
			assert.Equal(t, tc.expected, result, "Expected %v for filename %s", tc.expected, tc.filename)
		})
	}
}

// TestIsValidFile tests the file validation function
func TestIsValidFile(t *testing.T) {
	ts := test.TestSuite{}
	dir := ts.CreateFullTestSuite(t)
	defer ts.CleanTestSuite(t)

	t.Run("Valid file", func(t *testing.T) {
		filePath := filepath.Join(dir, "valid.toml")
		file, err := os.Create(filePath)
		require.NoError(t, err)
		_, err = file.WriteString("test content")
		require.NoError(t, err)
		file.Close()

		assert.True(t, isValidFile(filePath))
	})

	t.Run("Nonexistent file", func(t *testing.T) {
		nonexistentPath := filepath.Join(dir, "nonexistent.toml")
		assert.False(t, isValidFile(nonexistentPath))
	})

	t.Run("File too large", func(t *testing.T) {
		largeFilePath := filepath.Join(dir, "large.toml")
		file, err := os.Create(largeFilePath)
		require.NoError(t, err)

		// Write more than maxTranslationFileSize
		largeContent := strings.Repeat("a", maxTranslationFileSize+1)
		_, err = file.WriteString(largeContent)
		require.NoError(t, err)
		file.Close()

		assert.False(t, isValidFile(largeFilePath))
	})
}

// TestExtractLocaleFromFilename tests locale extraction from filenames
func TestExtractLocaleFromFilename(t *testing.T) {
	testCases := []struct {
		filename      string
		expectedTag   language.Tag
		expectedError bool
	}{
		{"active.en.toml", language.English, false},
		{"active.fr.json", language.French, false},
		{"messages.de.yaml", language.German, false},
		{"common.es.yml", language.Spanish, false},
		{"app.zh-CN.toml", language.MustParse("zh-CN"), false},
		{"test.pt-BR.json", language.MustParse("pt-BR"), false},
		{"invalid.toml", language.Und, true},                // missing locale
		{"noextension", language.Und, true},                 // no extension
		{"active..toml", language.Und, true},                // empty locale
		{"active.invalid-locale!.toml", language.Und, true}, // invalid language code
	}

	for _, tc := range testCases {
		t.Run(tc.filename, func(t *testing.T) {
			tag, err := extractLocaleFromFilename(tc.filename)
			if tc.expectedError {
				assert.Error(t, err)
				assert.Equal(t, language.Und, tag)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedTag, tag)
			}
		})
	}
}

// TestExtractAndValidateLocaleFromFilename tests the comprehensive validation function
func TestExtractAndValidateLocaleFromFilename(t *testing.T) {
	testCases := []struct {
		filename      string
		expectedTag   language.Tag
		expectedError bool
		errorContains string
	}{
		{"active.en.toml", language.English, false, ""},
		{"active.fr.json", language.French, false, ""},
		{"messages_v2.de.yaml", language.German, false, ""},
		{"common-ui.es.yml", language.Spanish, false, ""},

		// Invalid format cases
		{"invalid.toml", language.Und, true, "does not match expected format"},
		{"active..toml", language.Und, true, "does not match expected format"},
		{"active.en", language.Und, true, "does not match expected format"},
		{"", language.Und, true, "does not match expected format"},

		// Path traversal attempts (caught by regex validation)
		{"../active.en.toml", language.Und, true, "does not match expected format"},
		{"active.en.toml/../", language.Und, true, "does not match expected format"},

		// Invalid characters in prefix (caught by regex validation)
		{"active@.en.toml", language.Und, true, "does not match expected format"},
		{"active .en.toml", language.Und, true, "does not match expected format"},

		// Path that would pass regex but fail path traversal check
		{"active.en.toml", language.English, false, ""}, // This should actually pass

		// Invalid locale
		{"active.invalidlang.toml", language.Und, true, "invalid locale"},
	}

	for _, tc := range testCases {
		t.Run(tc.filename, func(t *testing.T) {
			tag, err := extractAndValidateLocaleFromFilename(tc.filename)
			if tc.expectedError {
				assert.Error(t, err)
				if tc.errorContains != "" {
					assert.Contains(t, err.Error(), tc.errorContains)
				}
				assert.Equal(t, language.Und, tag)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedTag, tag)
			}
		})
	}
}

// TestValidateBCP47Locale tests BCP 47 locale validation
func TestValidateBCP47Locale(t *testing.T) {
	testCases := []struct {
		tag           language.Tag
		expectedError bool
		errorContains string
	}{
		{language.English, false, ""},
		{language.French, false, ""},
		{language.MustParse("zh-CN"), false, ""},
		{language.MustParse("pt-BR"), false, ""},
		{language.MustParse("en-US"), false, ""},

		// Invalid cases
		{language.Und, true, "undefined language tag"},

		// Test would need a way to create invalid tags, which is difficult
		// since language.Parse() handles most validation
	}

	for _, tc := range testCases {
		t.Run(tc.tag.String(), func(t *testing.T) {
			err := validateBCP47Locale(tc.tag)
			if tc.expectedError {
				assert.Error(t, err)
				if tc.errorContains != "" {
					assert.Contains(t, err.Error(), tc.errorContains)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestTranslationFileStruct tests the translationFile struct
func TestTranslationFileStruct(t *testing.T) {
	// Test that the struct holds the expected data
	file := translationFile{
		path:   "/path/to/active.en.toml",
		locale: language.English,
	}

	assert.Equal(t, "/path/to/active.en.toml", file.path)
	assert.Equal(t, language.English, file.locale)
}

// TestSupportedExtensions tests that all expected extensions are supported
func TestSupportedExtensions(t *testing.T) {
	expectedExtensions := []string{".toml", ".json", ".yaml", ".yml"}

	assert.Equal(t, expectedExtensions, supportedExtensions)
	assert.Len(t, supportedExtensions, 4)
}

// TestMaxTranslationFileSize tests the file size constant
func TestMaxTranslationFileSize(t *testing.T) {
	expectedSize := 1024 * 1024 // 1MB
	assert.Equal(t, expectedSize, maxTranslationFileSize)
}

// TestTranslationFilenameRegex tests the regex pattern
func TestTranslationFilenameRegex(t *testing.T) {
	validNames := []string{
		"active.en.toml",
		"messages.fr.json",
		"common-ui.de.yaml",
		"app_v2.es.yml",
		"test123.zh-CN.toml",
	}

	invalidNames := []string{
		"active..toml",
		"active.en",
		"active.en.toml.backup",
		"active@.en.toml",
		"active .en.toml",
		"",
		"active",
		".en.toml",
	}

	for _, name := range validNames {
		t.Run("valid_"+name, func(t *testing.T) {
			assert.True(t, translationFilenameRegex.MatchString(name), "Expected %s to match regex", name)
		})
	}

	for _, name := range invalidNames {
		t.Run("invalid_"+name, func(t *testing.T) {
			assert.False(t, translationFilenameRegex.MatchString(name), "Expected %s to NOT match regex", name)
		})
	}
}
