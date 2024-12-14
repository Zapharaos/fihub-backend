package handlers

import (
	"github.com/Zapharaos/fihub-backend/internal/handlers/render"
	"go.uber.org/zap"
	"net/http"
)

func CreateTransaction(w http.ResponseWriter, r *http.Request) {
	userCtx, found := getUserFromContext(r)
	if !found {
		zap.L().Debug("No context user provided")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	render.JSON(w, r, userCtx)
}

func GetTransaction(w http.ResponseWriter, r *http.Request) {
	userCtx, found := getUserFromContext(r)
	if !found {
		zap.L().Debug("No context user provided")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	render.JSON(w, r, userCtx)
}

func UpdateTransaction(w http.ResponseWriter, r *http.Request) {
	userCtx, found := getUserFromContext(r)
	if !found {
		zap.L().Debug("No context user provided")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	render.JSON(w, r, userCtx)
}

func DeleteTransaction(w http.ResponseWriter, r *http.Request) {
	userCtx, found := getUserFromContext(r)
	if !found {
		zap.L().Debug("No context user provided")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	render.JSON(w, r, userCtx)
}

func GetAll(w http.ResponseWriter, r *http.Request) {
	userCtx, found := getUserFromContext(r)
	if !found {
		zap.L().Debug("No context user provided")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	render.JSON(w, r, userCtx)
}
