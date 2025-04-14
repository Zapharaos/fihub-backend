package app

import (
	"github.com/spf13/viper"
	"log"
)

// ConfigPath is the toml configuration file path
var ConfigPath = "config"

// EnvPrefix is the standard environment variable prefix
var EnvPrefix = "FIHUB"

// InitConfiguration initializes the application configuration
func InitConfiguration(name string) {
	// Set up Viper to read the main configuration file
	viper.SetConfigName("fihub-" + name)
	viper.AddConfigPath(ConfigPath)
	viper.SetEnvPrefix(EnvPrefix)
	viper.AutomaticEnv()

	// Read the configuration file
	err := viper.ReadInConfig()
	if err != nil {
		log.Printf("Failed to read configuration file: %v", err)
	}
}
