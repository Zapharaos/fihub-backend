package handlers

import (
	"go.uber.org/zap"
	"net/http"
)

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write([]byte("OK")); err != nil {
		zap.L().Error("Unable to write response")
		w.WriteHeader(http.StatusInternalServerError)
	}
}
