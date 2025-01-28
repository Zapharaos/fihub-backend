package translation

import (
	"errors"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"go.uber.org/zap"
	"golang.org/x/text/language"
)

// I18nService implements the Service interface using i18n
type I18nService struct {
	bundle     *i18n.Bundle
	localizers map[language.Tag]*i18n.Localizer
}

// NewI18nService returns a new instance of I18nService
func NewI18nService(defaultLang language.Tag) Service {

	// Create a new bundle
	bundle := i18n.NewBundle(defaultLang)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	// Load the translations
	// No need to load active.en.toml since we are providing default translations.
	bundle.MustLoadMessageFile("config/translations/active.en.toml")
	bundle.MustLoadMessageFile("config/translations/active.fr.toml")

	// Create localizers for each language
	localizers := map[language.Tag]*i18n.Localizer{
		defaultLang:     i18n.NewLocalizer(bundle, defaultLang.String()),
		language.French: i18n.NewLocalizer(bundle, language.French.String()),
	}

	// Create the service
	s := I18nService{
		bundle:     bundle,
		localizers: localizers,
	}
	var service Service = &s
	return service
}

// Localizer returns the requested localizer and an error if any
func (t *I18nService) Localizer(language language.Tag) (interface{}, error) {
	localizer, found := t.localizers[language]
	if !found {
		return nil, errors.New(fmt.Sprintf("localizer %q not found", language))
	}
	return localizer, nil
}

// Message returns a localized message for the given localizer and message
func (t *I18nService) Message(localizer interface{}, message *Message) string {
	// Verify that the localizer is of the correct type
	loc, ok := localizer.(*i18n.Localizer)
	if !ok {
		return ""
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
		zap.L().Error("Failed to localize", zap.Error(err))
		return ""
	}
	return result
}
