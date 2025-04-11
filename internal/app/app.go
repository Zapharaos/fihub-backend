package app

import (
	"github.com/Zapharaos/fihub-backend/pkg/email"
	"github.com/Zapharaos/fihub-backend/pkg/env"
	"github.com/Zapharaos/fihub-backend/pkg/translation"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/text/language"
	"log"
)

var (
	// Version is the binary version + build number
	Version = ""
	// BuildDate is the date of build
	BuildDate = ""
)

// Init initialize all the app configuration and components
func Init() {

	// Load the .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
		return
	}

	// Setup Environment
	InitConfiguration()

	// Setup Logger
	InitLogger()

	zap.L().Info("Starting Fihub Backend", zap.String("version", Version), zap.String("build_date", BuildDate))

	// Setup Database
	initDatabase()

	// Setup Email
	email.ReplaceGlobals(email.NewSendgridService())

	// Setup Translations
	defaultLang := language.MustParse(env.GetString("DEFAULT_LANG", "en"))
	translation.ReplaceGlobals(translation.NewI18nService(defaultLang))
}

// InitLogger initializes the Zap logger.
func InitLogger() zap.Config {

	// Set environment config
	var zapConfig zap.Config
	if env.GetBool("LOGGER_PROD", true) {
		zapConfig = zap.NewProductionConfig()
	} else {
		zapConfig = zap.NewDevelopmentConfig()
	}

	zapConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	// Logger level
	switch env.GetString("LOGGER_LEVEL", "debug") {
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
