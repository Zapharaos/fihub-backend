package translation

import (
	"fmt"
	"golang.org/x/text/language"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// translationFile represents a discovered translation file
type translationFile struct {
	path   string
	locale language.Tag
}

// Supported translation file extensions
var supportedExtensions = []string{".toml", ".json", ".yaml", ".yml"}

// Regular expression for validating translation filename format
// Matches: prefix.locale.ext where prefix contains only alphanumeric chars, hyphens, underscores
// and locale follows BCP 47 format (letters, numbers, hyphens)
var translationFilenameRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]+\.[a-zA-Z0-9-]+\.[a-zA-Z0-9]+$`)

// Maximum file size for translation files (1MB should be more than enough)
const maxTranslationFileSize = 1024 * 1024

// discoverTranslationFiles scans the given directory for translation files.
// If filePrefixes is empty, all valid translation files are returned.
// If filePrefixes is provided, only files with those prefixes are returned.
func discoverTranslationFiles(translationsPath string, filePrefixes ...string) ([]translationFile, error) {
	// Check if the path exists and is a directory
	pathInfo, err := os.Stat(translationsPath)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("translations path does not exist: %s", translationsPath)
	}
	if err != nil {
		return nil, fmt.Errorf("cannot access translations path %s: %w", translationsPath, err)
	}
	if !pathInfo.IsDir() {
		return nil, fmt.Errorf("translations path is not a directory: %s", translationsPath)
	}

	// Read directory contents
	entries, err := os.ReadDir(translationsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read translations directory: %w", err)
	}

	var translationFiles []translationFile
	var invalidFiles []string

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		fileName := entry.Name()
		fullPath := filepath.Join(translationsPath, fileName)

		// Check if file has supported extension
		if !hasSupportedExtension(fileName) {
			continue
		}

		// Check if the file has the correct prefix (if specified)
		if len(filePrefixes) > 0 {
			hasValidPrefix := false
			for _, prefix := range filePrefixes {
				if strings.HasPrefix(fileName, prefix+".") {
					hasValidPrefix = true
					break
				}
			}
			if !hasValidPrefix {
				continue
			}
		}

		// Validate the file is actually a valid file
		if !isValidFile(fullPath) {
			invalidFiles = append(invalidFiles, fileName)
			continue
		}

		// Extract and validate locale from filename
		locale, err := extractAndValidateLocaleFromFilename(fileName)
		if err != nil {
			invalidFiles = append(invalidFiles, fileName)
			continue
		}

		translationFiles = append(translationFiles, translationFile{
			path:   fullPath,
			locale: locale,
		})
	}

	// Report invalid files if any were found
	if len(invalidFiles) > 0 {
		prefixMsg := ""
		if len(filePrefixes) > 0 {
			prefixMsg = fmt.Sprintf(" with prefixes %v", filePrefixes)
		}
		supportedExtsStr := strings.Join(supportedExtensions, ", ")
		return translationFiles, fmt.Errorf("found %d invalid translation files%s: %v (files must follow format 'prefix.{locale}.{ext}' where ext is one of: %s)", len(invalidFiles), prefixMsg, invalidFiles, supportedExtsStr)
	}

	return translationFiles, nil
}

// hasSupportedExtension checks if the filename has a supported translation file extension
func hasSupportedExtension(filename string) bool {
	lower := strings.ToLower(filename)
	for _, ext := range supportedExtensions {
		if strings.HasSuffix(lower, ext) {
			return true
		}
	}
	return false
}

// isValidFile performs comprehensive validation of a translation file
func isValidFile(filePath string) bool {
	// Check file size to prevent loading extremely large files
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return false
	}

	if fileInfo.Size() > maxTranslationFileSize {
		return false
	}

	// Check if file is readable
	file, err := os.Open(filePath)
	if err != nil {
		return false
	}
	defer file.Close()

	return true
}

// extractLocaleFromFilename extracts the language tag from a translation filename
// Expected format: "prefix.{locale}.ext" (e.g., "active.en.toml", "active.fr.json")
func extractLocaleFromFilename(filename string) (language.Tag, error) {
	// Remove any supported extension
	nameWithoutExt := filename
	for _, ext := range supportedExtensions {
		if strings.HasSuffix(strings.ToLower(filename), ext) {
			nameWithoutExt = strings.TrimSuffix(filename, ext)
			break
		}
	}

	// Split by dots
	parts := strings.Split(nameWithoutExt, ".")
	if len(parts) < 2 {
		return language.Und, fmt.Errorf("invalid translation filename format: %s", filename)
	}

	// The locale should be the last part before the extension
	localeStr := parts[len(parts)-1]

	// Parse the language tag
	locale, err := language.Parse(localeStr)
	if err != nil {
		return language.Und, fmt.Errorf("invalid locale in filename %s: %w", filename, err)
	}

	return locale, nil
}

// extractAndValidateLocaleFromFilename extracts and validates the locale from the filename
func extractAndValidateLocaleFromFilename(filename string) (language.Tag, error) {
	// First, validate the filename format using regex
	if !translationFilenameRegex.MatchString(filename) {
		supportedExtsStr := strings.Join(supportedExtensions, ", ")
		return language.Und, fmt.Errorf("filename '%s' does not match expected format 'prefix.locale.{ext}' where ext is one of: %s (only alphanumeric, hyphens, underscores allowed)", filename, supportedExtsStr)
	}

	// Sanitize filename to prevent path traversal attacks
	cleanFilename := filepath.Base(filename)
	if cleanFilename != filename {
		return language.Und, fmt.Errorf("filename '%s' contains invalid path characters", filename)
	}

	// Extract locale using the existing function
	locale, err := extractLocaleFromFilename(filename)
	if err != nil {
		return language.Und, err
	}

	// Additional BCP 47 validation
	if err := validateBCP47Locale(locale); err != nil {
		return language.Und, fmt.Errorf("invalid BCP 47 locale in filename '%s': %w", filename, err)
	}

	return locale, nil
}

// validateBCP47Locale performs additional validation on the parsed language tag
func validateBCP47Locale(tag language.Tag) error {
	// Check if the tag is valid and not undefined
	if tag == language.Und {
		return fmt.Errorf("undefined language tag")
	}

	// Get the base language
	base, _ := tag.Base()
	if base.String() == "" {
		return fmt.Errorf("invalid base language")
	}

	// Ensure the tag string representation is reasonable (not too long)
	tagStr := tag.String()
	if len(tagStr) > 35 { // BCP 47 recommends max 35 characters
		return fmt.Errorf("language tag too long: %s", tagStr)
	}

	// Ensure no suspicious characters in the tag
	if strings.ContainsAny(tagStr, "/\\<>:\"|?*") {
		return fmt.Errorf("language tag contains invalid characters: %s", tagStr)
	}

	return nil
}
