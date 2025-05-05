package app

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
)

// InitLogger initializes the Zap logger.
func InitLogger() zap.Config {

	// Set environment config
	var zapConfig zap.Config
	if viper.GetString("APP_ENV") != "production" {
		zapConfig = zap.NewDevelopmentConfig()
	} else {
		zapConfig = zap.NewProductionConfig()
	}

	zapConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	// Logger level
	switch viper.GetString("LOGGER_LEVEL") {
	case "debug":
		zapConfig.Level.SetLevel(zap.DebugLevel)
	case "info":
		zapConfig.Level.SetLevel(zap.InfoLevel)
	case "warn":
		zapConfig.Level.SetLevel(zap.WarnLevel)
	case "error":
		zapConfig.Level.SetLevel(zap.ErrorLevel)
	case "dpanic":
		zapConfig.Level.SetLevel(zap.DPanicLevel)
	case "panic":
		zapConfig.Level.SetLevel(zap.PanicLevel)
	case "fatal":
		zapConfig.Level.SetLevel(zap.FatalLevel)
	default:
		zapConfig.Level.SetLevel(zap.InfoLevel)
	}

	// Constructs logger
	logger, err := zapConfig.Build(zap.AddStacktrace(zap.ErrorLevel))
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}

	defer func() {
		if err = logger.Sync(); err != nil {
			log.Printf("can't sync logger: %v\n", err)
		}
	}()

	zap.ReplaceGlobals(logger)
	return zapConfig
}
