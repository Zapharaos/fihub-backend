package app

import (
	"github.com/Zapharaos/fihub-backend/internal/database"
	"go.uber.org/zap"
)

var (
	// Version is the binary version + build number
	Version = ""
	// BuildDate is the date of build
	BuildDate = ""
)

// RecoverPanic is a function that recovers from a panic and logs the error.
func RecoverPanic() {
	if r := recover(); r != nil {
		zap.L().Error("Unhandled panic occurred", zap.Any("panic", r))
		panic(r)
	}
}

// CleanResources is a function that closes all app resources including databases.
func CleanResources() {
	database.DB().CloseAll()
}
