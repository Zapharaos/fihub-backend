//go:build example
// +build example

package main

import (
	"fmt"
	"github.com/Zapharaos/fihub-backend/pkg/translation"
	"golang.org/x/text/language"
	"log"
)

func main() {
	fmt.Println("=== Translation Package Demo ===")

	// Step 1: Initialize the translation service
	fmt.Println("1. Initializing translation service...")
	translationService, err := translation.NewI18nService(
		language.English, // default language
		"config/",        // translations directory
		"messages",       // file prefix
	)
	if err != nil {
		log.Fatalf("Failed to initialize translation service: %v", err)
	}

	// Replace the global translation service instance
	translation.ReplaceGlobals(translationService)
	fmt.Println("✓ Translation service initialized successfully")

	// Step 2: Demonstrate basic translation usage
	fmt.Println("2. Basic Translation Examples:")
	demonstrateBasicTranslation()

	// Step 3: Demonstrate MustTranslate usage
	fmt.Println("\n3. MustTranslate Examples:")
	demonstrateMustTranslate()

	// Step 4: Demonstrate template variables
	fmt.Println("\n4. Template Variables Examples:")
	demonstrateTemplateVariables()

	// Step 5: Demonstrate pluralization
	fmt.Println("\n5. Pluralization Examples:")
	demonstratePluralization()

	// Step 6: Demonstrate fallback behavior
	fmt.Println("\n6. Fallback Behavior:")
	demonstrateFallback()
}

func demonstrateBasicTranslation() {
	languages := []language.Tag{language.English, language.French, language.Spanish, language.German}

	for _, lang := range languages {
		localizer, found, err := translation.S().Localizer(lang)
		if err != nil {
			fmt.Printf("Error getting localizer for %s: %v\n", lang, err)
			continue
		}

		message := &translation.Message{ID: "hello_world"}
		result, success, err := translation.S().Translate(localizer, message)

		foundStr := "✓"
		if !found {
			foundStr = "⚠ (fallback)"
		}

		if err != nil {
			fmt.Printf("  %s: Error - %v\n", lang, err)
		} else if success {
			fmt.Printf("  %s %s: %s\n", foundStr, lang, result)
		} else {
			fmt.Printf("  %s %s: %s (translation failed)\n", foundStr, lang, result)
		}
	}
}

func demonstrateMustTranslate() {
	languages := []language.Tag{language.English, language.French}

	for _, lang := range languages {
		localizer, _, err := translation.S().Localizer(lang)
		if err != nil {
			fmt.Printf("Error getting localizer for %s: %v\n", lang, err)
			continue
		}

		// Safe usage of MustTranslate - we know this key exists
		result := translation.S().MustTranslate(localizer, &translation.Message{
			ID: "goodbye",
		})
		fmt.Printf("  %s: %s\n", lang, result)
	}

	// Example of error handling with MustTranslate
	fmt.Println("\n  Demonstrating MustTranslate panic (commented out for safety):")
	fmt.Println("  // This would panic: translation.S().MustTranslate(localizer, &translation.Message{ID: \"nonexistent_key\"})")
}

func demonstrateTemplateVariables() {
	localizer, _, _ := translation.S().Localizer(language.English)

	// Simple template variable
	message := &translation.Message{
		ID: "welcome_user",
		Data: map[string]interface{}{
			"Name": "John Doe",
		},
	}
	result := translation.S().MustTranslate(localizer, message)
	fmt.Printf("  English (with name): %s\n", result)

	// Multiple variables
	message = &translation.Message{
		ID: "user_stats",
		Data: map[string]interface{}{
			"Name":  "Alice",
			"Posts": 42,
			"Likes": 128,
		},
	}
	result = translation.S().MustTranslate(localizer, message)
	fmt.Printf("  English (with stats): %s\n", result)

	// French version
	localizerFr, _, _ := translation.S().Localizer(language.French)
	result = translation.S().MustTranslate(localizerFr, message)
	fmt.Printf("  French (with stats): %s\n", result)
}

func demonstratePluralization() {
	localizer, _, _ := translation.S().Localizer(language.English)

	pluralCounts := []int{0, 1, 2, 5}

	for _, count := range pluralCounts {
		var message *translation.Message

		// Handle zero case separately since English CLDR doesn't support "zero" category
		if count == 0 {
			message = &translation.Message{
				ID: "no_items",
			}
		} else {
			message = &translation.Message{
				ID:          "item_count",
				PluralCount: count,
				Data: map[string]interface{}{
					"Count": count,
				},
			}
		}

		result := translation.S().MustTranslate(localizer, message)
		fmt.Printf("  %d items: %s\n", count, result)
	}

	// French pluralization (different rules)
	fmt.Println("\n  French pluralization:")
	localizerFr, _, _ := translation.S().Localizer(language.French)
	for _, count := range pluralCounts {
		var message *translation.Message

		// Handle zero case for French too
		if count == 0 {
			message = &translation.Message{
				ID: "no_items",
			}
		} else {
			message = &translation.Message{
				ID:          "item_count",
				PluralCount: count,
				Data: map[string]interface{}{
					"Count": count,
				},
			}
		}

		result := translation.S().MustTranslate(localizerFr, message)
		fmt.Printf("  %d items: %s\n", count, result)
	}
}

func demonstrateFallback() {
	// Try to get a localizer for a language that doesn't exist
	unsupportedLang := language.MustParse("ja") // Japanese - not in our example files

	localizer, found, err := translation.S().Localizer(unsupportedLang)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if !found {
		fmt.Printf("Japanese localizer not found, falling back to default (English)\n")
	}

	message := &translation.Message{ID: "hello_world"}
	result := translation.S().MustTranslate(localizer, message)
	fmt.Printf("  Result: %s\n", result)

	// Try to translate a key that doesn't exist
	fmt.Println("\n  When translating a missing key:")
	localizerEn, _, _ := translation.S().Localizer(language.English)
	message = &translation.Message{ID: "nonexistent_key"}
	result, success, err := translation.S().Translate(localizerEn, message)

	if err != nil {
		fmt.Printf("  Error translating missing key: %v\n", err)
	} else if !success {
		fmt.Printf("  Translation failed, returned: %s\n", result)
	}
}
