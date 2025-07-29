package otp

import (
	"fmt"
	"github.com/Zapharaos/fihub-backend/pkg/email/templates"
	"github.com/Zapharaos/fihub-backend/pkg/translation"
	"go.uber.org/zap"
	"golang.org/x/text/language"
	"time"
)

// BuildOtpEmailContents prepares the email contents for OTP verification
func BuildOtpEmailContents(userLanguage language.Tag, otp string, timeLimit time.Duration) (subject, plainTextContent, htmlContent string, err error) {
	// Get localizer
	loc, _, err := translation.S().Localizer(userLanguage)
	if err != nil {
		zap.L().Error("Failed to get localizer", zap.Error(err))
		return
	}

	// Prepare email data with translations
	subject = translation.S().MustTranslate(loc, &translation.Message{ID: "EmailOtpTitle"})
	plainTextContent = translation.S().MustTranslate(loc, &translation.Message{
		ID: "EmailOtpPlainTextContent",
		Data: map[string]interface{}{
			"Otp": otp,
		},
	})

	// Prepare email html template
	htmlContentTemplate := templates.NewOtpTemplate(templates.OtpData{
		OTP:      otp,
		Greeting: translation.S().MustTranslate(loc, &translation.Message{ID: "EmailGreeting"}),
		MainContent: translation.S().MustTranslate(loc, &translation.Message{
			ID: "EmailOtpContentForgotPassword",
			Data: map[string]interface{}{
				"Duration": fmt.Sprintf("%d", int(timeLimit.Minutes())),
			},
		}),
		DoNotShare: translation.S().MustTranslate(loc, &translation.Message{ID: "EmailOtpDoNotShare"}),
	})

	// Prepare email layout labels
	labels := templates.LayoutLabels{
		Help: translation.S().MustTranslate(loc, &translation.Message{ID: "EmailFooterHelp"}),
		Copyrights: translation.S().MustTranslate(loc, &translation.Message{
			ID: "EmailFooterCopyrights",
			Data: map[string]interface{}{
				"Year": time.Now().Year(),
			},
		}),
	}

	// Render email html content
	htmlContent, err = htmlContentTemplate.Build(labels)
	if err != nil {
		// Log error and use plain text content instead of HTML
		zap.L().Error("Render email content", zap.Error(err))
		htmlContent = plainTextContent
	}

	return
}
