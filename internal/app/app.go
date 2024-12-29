package app

import (
	"github.com/Zapharaos/fihub-backend/internal/auth/password"
	"github.com/Zapharaos/fihub-backend/internal/auth/users"
	"github.com/Zapharaos/fihub-backend/internal/brokers"
	"github.com/Zapharaos/fihub-backend/internal/database"
	"github.com/Zapharaos/fihub-backend/internal/transactions"
	"github.com/Zapharaos/fihub-backend/pkg/env"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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
	err := env.Load()
	if err != nil {
		log.Fatal(err)
		return
	}

	// Setup Logger
	initLogger()

	zap.L().Info("Starting Fihub Backend", zap.String("version", Version), zap.String("build_date", BuildDate))

	// Setup Database
	initPostgres()
	initRepositories()
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
		}
	}()

	zap.ReplaceGlobals(logger)
	return zapConfig
}

// initPostgres initializes the postgres connection.
func initPostgres() {

	zap.L().Info("Initializing Postgres")

	// Configure postgres
	credentials := postgres.Credentials{
		Host:     env.GetString("POSTGRES_HOST", "host"),
		Port:     env.GetString("POSTGRES_PORT", "port"),
		DbName:   env.GetString("POSTGRES_DB", "database_name"),
		User:     env.GetString("POSTGRES_USER", "user"),
		Password: env.GetString("POSTGRES_PASSWORD", "password"),
	}

	// Connect
	dbClient, err := postgres.DbConnection(credentials)
	if err != nil {
		zap.L().Fatal("main.DbConnection:", zap.Error(err))
	}

	zap.L().Info("Connected to Postgres")

	// Finish up configuration
	dbClient.SetMaxOpenConns(env.GetInt("POSTGRES_MAX_OPEN_CONNS", 30))
	dbClient.SetMaxIdleConns(env.GetInt("POSTGRES_MAX_IDLE_CONNS", 30))
	postgres.ReplaceGlobals(dbClient)
}

// initRepositories initializes the repositories. (postgres)
func initRepositories() {
	// Setup for postgres
	dbClient := postgres.DB()

	// Auth
	users.ReplaceGlobals(users.NewPostgresRepository(dbClient))
	password.ReplaceGlobals(password.NewPostgresRepository(dbClient))

	// Brokers
	brokerRepository := brokers.NewPostgresRepository(dbClient)
	userBrokerRepository := brokers.NewUserBrokerPostgresRepository(dbClient)
	imageBrokerRepository := brokers.NewImageBrokerPostgresRepository(dbClient)
	brokers.ReplaceGlobals(brokers.NewRepository(brokerRepository, userBrokerRepository, imageBrokerRepository))

	// Transactions
	transactions.ReplaceGlobals(transactions.NewPostgresRepository(dbClient))
}
