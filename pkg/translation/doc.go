// Package translation provides functionality for internationalization and localization.
//
// This package defines an interface for translation services, allowing for different implementations.
// The default implementation uses the go-i18n library to manage translations and localize messages.
// It supports multiple languages and allows for easy addition of new translations.
//
// # Supported File Extensions
//
// The translation service supports the following file formats:
//   - .toml (TOML format)
//   - .json (JSON format)
//   - .yaml (YAML format)
//   - .yml (YAML format)
//
// # File Naming Convention
//
// Translation files must follow the naming pattern: `{prefix}.{locale}.{extension}`
//
// Examples:
//   - active.en.toml
//   - active.fr.json
//   - messages.es.yaml
//   - app.de.yml
//
// Where:
//   - prefix: alphanumeric characters, hyphens, and underscores only
//   - locale: valid BCP 47 language tag (e.g., en, fr, es-ES, zh-CN)
//   - extension: one of the supported file extensions listed above
//
// # File Prefix Usage
//
// File prefixes allow you to organize translation files by context or module.
// When initializing the translation service, you can specify which prefixes to load:
//
//	// Load only files with "active" prefix (active.en.toml, active.fr.json, etc.)
//	service, err := translation.NewI18nService(language.English, "config/translations", "active")
//
//	// Load files with multiple prefixes
//	service, err := translation.NewI18nService(language.English, "config/translations", "active", "messages")
//
//	// Load all valid translation files (no prefix filter)
//	service, err := translation.NewI18nService(language.English, "config/translations")
//
// # Multiple Files for Same Locale
//
// When multiple translation files exist for the same locale with different extensions,
// the go-i18n library will merge their contents. Files are processed in the order they
// are discovered. If duplicate message IDs exist across files for the same locale,
// the last loaded file will take precedence.
//
// Example: If both `active.fr.toml` and `active.fr.json` exist, both will be loaded
// and their translations merged for the French locale.
//
// # File Size and Validation
//
// Translation files are subject to the following constraints:
//   - Maximum file size: 1MB
//   - Files must be readable and properly formatted
//   - Locale codes must follow BCP 47 standards
//   - Filenames cannot contain path traversal characters
//
// When a requested language is not available, the service automatically falls back
// to the default language specified during initialization.
//
// For more information, see the documentation for the go-i18n library.
package translation
