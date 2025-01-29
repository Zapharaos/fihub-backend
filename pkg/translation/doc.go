// Package translation provides functionality for internationalization and localization.
//
// This package defines an interface for translation services, allowing for different implementations.
// The default implementation uses the go-i18n library to manage translations and localize messages.
// It supports multiple languages and allows for easy addition of new translations.
//
// In your main.go or application initialization file, you can initialize the translation service like this:
//
//	package main
//
//	import (
//	    "github.com/Zapharaos/fihub-backend/pkg/translation"
//	    "golang.org/x/text/language"
//	)
//
//	func main() {
//	    // Initialize the translation service
//	    translationService := translation.NewI18nService(language.English)
//
//	    // Replace the global translation service instance
//	    translation.ReplaceGlobals(translationService)
//	}
//
// To get a localized message:
//
//	localizer, _ := translation.S().Localizer(language.French)
//	message := translation.S().Message(localizer, &translation.Message{ID: "HelloWorld"})
//
// For more information, see the documentation for the go-i18n library.
package translation
