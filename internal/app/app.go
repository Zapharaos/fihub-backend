package app

import (
	"github.com/Zapharaos/fihub-backend/internal/auth/password"
	"github.com/Zapharaos/fihub-backend/internal/auth/permissions"
	"github.com/Zapharaos/fihub-backend/internal/auth/roles"
	"github.com/Zapharaos/fihub-backend/internal/auth/users"
	"github.com/Zapharaos/fihub-backend/internal/brokers"
	"github.com/Zapharaos/fihub-backend/internal/database"
	"github.com/Zapharaos/fihub-backend/internal/transactions"
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

	// Setup Logger
	initLogger()

	zap.L().Info("Starting Fihub Backend", zap.String("version", Version), zap.String("build_date", BuildDate))

	// Setup Database
	initDatabase()

	// Setup Email
	email.ReplaceGlobals(email.NewSendgridService())

	// Setup Translations
	defaultLang := language.MustParse(env.GetString("DEFAULT_LANG", "en"))
	translation.ReplaceGlobals(translation.NewI18nService(defaultLang))
}

// initPostgres initializes the Zap logger.
func initLogger() zap.Config {

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

// initDatabase initializes the database connections.
func initDatabase() {
	postgres := database.NewPostgresDB(database.NewSqlDatabase(database.SqlCredentials{
		Host:     env.GetString("POSTGRES_HOST", "localhost"),
		Port:     env.GetString("POSTGRES_PORT", "5432"),
		User:     env.GetString("POSTGRES_USER", "postgres"),
		Password: env.GetString("POSTGRES_PASSWORD", "password"),
		DbName:   env.GetString("POSTGRES_DB", "postgres"),
	}))
	database.ReplaceGlobals(database.NewDatabases(postgres))

	// Initialize the postgres repositories
	initPostgres()
}

// initPostgres initializes the postgres repositories.
func initPostgres() {
	// Setup for postgres
	dbClient := database.DB().Postgres()

	// Auth
	users.ReplaceGlobals(users.NewPostgresRepository(dbClient))
	password.ReplaceGlobals(password.NewPostgresRepository(dbClient))

	// Roles
	roles.ReplaceGlobals(roles.NewPostgresRepository(dbClient))

	// Permissions
	permissions.ReplaceGlobals(permissions.NewPostgresRepository(dbClient))

	// Brokers
	brokerRepository := brokers.NewPostgresRepository(dbClient)
	userBrokerRepository := brokers.NewUserBrokerPostgresRepository(dbClient)
	imageBrokerRepository := brokers.NewImageBrokerPostgresRepository(dbClient)
	brokers.ReplaceGlobals(brokers.NewRepository(brokerRepository, userBrokerRepository, imageBrokerRepository))

	// Transactions
	transactions.ReplaceGlobals(transactions.NewPostgresRepository(dbClient))
}
