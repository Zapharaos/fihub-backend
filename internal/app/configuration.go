package app

import (
	"github.com/spf13/viper"
	"log"
)

// ConfigPath is the toml configuration file path
var ConfigPath = "config"

// ConfigName is the toml configuration file name
var ConfigName = "fihub-backend"

// EnvPrefix is the standard environment variable prefix
var EnvPrefix = "FIHUB"

// InitConfiguration initializes the application configuration
func InitConfiguration() {
	// Set up Viper to read the main configuration file
	viper.SetConfigName(ConfigName)
	viper.AddConfigPath(ConfigPath)
	viper.SetEnvPrefix(EnvPrefix)
	viper.AutomaticEnv()

	// Read the configuration file
	err := viper.ReadInConfig()
	if err != nil {
		log.Printf("Failed to read configuration file: %v", err)
	}

	// Custom plugins config
	v := viper.New()
	v.SetConfigName("services")
	v.AddConfigPath(ConfigPath)
	err = v.ReadInConfig()
	if err != nil {
		log.Printf("No plugins configuration found: %v", err)
		return
	}

	// Merge plugin configuration into the main configuration
	err = viper.MergeConfigMap(v.AllSettings())
	if err != nil {
		log.Printf("Failed to merge plugins configuration: %v", err)
	}
}
