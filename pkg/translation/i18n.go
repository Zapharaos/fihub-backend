package translation

import (
	"encoding/json"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
)

// I18nService implements the Service interface using i18n
type I18nService struct {
	bundle      *i18n.Bundle
	localizers  map[language.Tag]*i18n.Localizer
	defaultLang language.Tag
}

// NewI18nService returns a new instance of I18nService with a custom file prefix
// defaultLang: the default language to use when a requested language is not available
// translationsPath: path to the directory containing translation files
// filePrefixes: the prefixes that translation files can have (e.g., "active" for "active.en.toml")
func NewI18nService(defaultLang language.Tag, translationsPath string, filePrefixes ...string) (Service, error) {
	// Create a new bundle
	bundle := i18n.NewBundle(defaultLang)

	// Register unmarshal functions for all supported file formats
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	bundle.RegisterUnmarshalFunc("yaml", yaml.Unmarshal)
	bundle.RegisterUnmarshalFunc("yml", yaml.Unmarshal)

	// Discover and load translation files from the given path
	translationFiles, err := discoverTranslationFiles(translationsPath, filePrefixes...)
	if err != nil {
		return nil, fmt.Errorf("failed to discover translation files: %w", err)
	}

	if len(translationFiles) == 0 {
		return nil, fmt.Errorf("no corresponding translation files were found in path: %s", translationsPath)
	}

	// Load all discovered translation files
	availableLocales := make([]language.Tag, 0, len(translationFiles))
	for _, file := range translationFiles {
		_, err := bundle.LoadMessageFile(file.path)
		if err != nil {
			return nil, fmt.Errorf("failed to load translation file %s: %w", file.path, err)
		}
		availableLocales = append(availableLocales, file.locale)
	}

	// Verify that the default language is available
	defaultFound := false
	for _, locale := range availableLocales {
		if locale == defaultLang {
			defaultFound = true
			break
		}
	}
	if !defaultFound {
		return nil, fmt.Errorf("default language %s not found in available translations files", defaultLang)
	}

	// Create localizers for each available language
	localizers := make(map[language.Tag]*i18n.Localizer, len(availableLocales))
	for _, locale := range availableLocales {
		localizers[locale] = i18n.NewLocalizer(bundle, locale.String())
	}

	// Create the service
	s := I18nService{
		bundle:      bundle,
		localizers:  localizers,
		defaultLang: defaultLang,
	}
	return &s, nil
}

// Localizer returns the requested localizer and a boolean indicating if the localizer was found
// If the requested localizer is not found, returns the default language localizer
func (t *I18nService) Localizer(language language.Tag) (interface{}, bool, error) {
	localizer, found := t.localizers[language]
	if !found {
		// Return the default language localizer
		defaultLocalizer := t.localizers[t.defaultLang]
		if defaultLocalizer == nil {
			return nil, false, fmt.Errorf("default localizer %s not found, please check your translations configuration", t.defaultLang)
		}
		return defaultLocalizer, false, nil
	}
	return localizer, true, nil
}

// Translate returns a localized message for the given localizer and message
// Returns the translated message, a boolean indicating success, and an error if something went wrong
func (t *I18nService) Translate(localizer interface{}, message *Message) (string, bool, error) {
	// Verify that the localizer is of the correct type
	loc, ok := localizer.(*i18n.Localizer)
	if !ok {
		return "", false, fmt.Errorf("invalid localizer type: expected *i18n.Localizer, got %T", localizer)
	}

	// Validate that message is not nil
	if message == nil {
		return "", false, fmt.Errorf("message cannot be nil")
	}

	// Validate that message ID is not empty
	if message.ID == "" {
		return "", false, fmt.Errorf("message ID cannot be empty")
	}

	// Map Message to i18n.LocalizeConfig
	localizeConfig := &i18n.LocalizeConfig{
		MessageID:    message.ID,
		TemplateData: message.Data,
		PluralCount:  message.PluralCount,
	}

	// Localize the message
	result, err := loc.Localize(localizeConfig)
	if err != nil {
		return "", false, fmt.Errorf("failed to localize message '%s': %w", message.ID, err)
	}

	return result, true, nil
}

// MustTranslate returns a localized message, panicking on error
// This is useful when you're confident the translation should always work
func (t *I18nService) MustTranslate(localizer interface{}, message *Message) string {
	result, _, err := t.Translate(localizer, message)
	if err != nil {
		panic(fmt.Sprintf("translation failed: %v", err))
	}
	return result
}
