// Package templates provides functionality for managing email templates.
//
// This package provides a main HTML email layout. Multiple content templates can be created
// to be used as the main item inside the initial email layout.
//
// Usage:
//
// To define the main HTML content and set up its structure:
//
//	htmlContentTemplate := templates.NewCustomTemplate(templates.CustomData{
//	    Field1: "Hello",
//	    Field2: "$123456789",
//	    Field3: "29/01/2025",
//	})
//
// To set the variables for the initial email layout:
//
//	labels := templates.LayoutLabels{
//	    Help:       "Help",
//	    Copyrights: "Copyrights",
//	}
//
// To build the full email HTML content with the specified data:
//
//	htmlContent, err := htmlContentTemplate.Build(labels)
//	if err != nil {
//	    // Handle error
//	}
//
// For more information, see the documentation for the text/template and html/template packages.
package templates
