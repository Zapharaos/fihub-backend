// Package env provides utilities for loading and retrieving environment variables.
//
// This package offers functions to load environment variables from a `.env` file and retrieve them
// as different types such as string, integer, boolean, and duration.
//
// Usage:
//
// To load environment variables from a `.env` file:
//
//	err := env.Load()
//	if err != nil {
//	    log.Fatal("Error loading .env file")
//	}
//
// To retrieve an environment variable as a string:
//
//	value := env.GetString("KEY", "default")
//
// To retrieve an environment variable as an integer:
//
//	value := env.GetInt("KEY", 42)
//
// To retrieve an environment variable as a boolean:
//
//	value := env.GetBool("KEY", true)
//
// To retrieve an environment variable as a duration:
//
//	value := env.GetDuration("KEY", time.Second)
//
// For more information, see the documentation for the `github.com/joho/godotenv` library.
package env
