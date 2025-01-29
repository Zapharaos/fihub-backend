// Package email provides functionality for sending emails.
//
// This package defines an interface for email services, allowing for different implementations.
// The default implementation uses the SendGrid service to send emails.
//
// In your main.go or application initialization file, you can initialize the email service like this:
//
//	package main
//
//	import (
//	    "github.com/Zapharaos/fihub-backend/pkg/email"
//	)
//
//	func main() {
//	    // Initialize the email service
//	    emailService := email.NewSendgridService()
//
//	    // Replace the global email service instance
//	    email.ReplaceGlobals(emailService)
//	}
//
// To send an email:
//
//	err := service.Send("recipient@example.com", "Subject", "Plain text content", "HTML content")
//
// For more information, see the documentation for the SendGrid library.
package email
