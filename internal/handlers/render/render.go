package render

import (
	"encoding/json"
	"go.uber.org/zap"
	"net/http"
)

type ErrorResponse struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

type CountResponse struct {
	Count int64 `json:"count"`
}

// Error returns an HTTP status 500 with a specific error message
func Error(w http.ResponseWriter, r *http.Request, err error, message string) {
	resp := ErrorResponse{Title: "Internal Server Error - Please contact the administrator."}
	if err != nil {
		if message != "" {
			zap.L().Error(message, zap.Error(err))
		} else {
			resp.Message = err.Error()
		}
	}
	w.WriteHeader(http.StatusInternalServerError)
	JSON(w, r, resp)
}

// BadRequest returns an HTTP status 400 with a specific error message
func BadRequest(w http.ResponseWriter, r *http.Request, err error) {
	resp := ErrorResponse{Title: "Bad Request - Please check your request."}
	if err != nil {
		zap.L().Debug("Bad Request", zap.Error(err))
		resp.Message = err.Error()
	}
	w.WriteHeader(http.StatusBadRequest)
	JSON(w, r, resp)
}

// OK returns an HTTP status 200 with an empty body
func OK(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

// NotImplemented returns an HTTP status 501
func NotImplemented(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	JSON(w, r, map[string]interface{}{"message": "Not Implemented"})
}

// JSON try to encode an interface and returns it in a specific ResponseWriter (or returns an internal server error)
func JSON(w http.ResponseWriter, r *http.Request, data interface{}) {
	OK(w, r)

	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		zap.L().Error("Render JSON encode", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// NotFound returns an HTTP status 404 with a specific error message
func NotFound(w http.ResponseWriter, r *http.Request, err error) {
	resp := ErrorResponse{Title: "Not Found"}
	if err != nil {
		zap.L().Debug("Not Found", zap.Error(err))
		resp.Message = err.Error()
	}
	w.WriteHeader(http.StatusNotFound)
	JSON(w, r, resp)
}

// Count returns an HTTP status 200 with a JSON object containing the count (CountResponse)
func Count(w http.ResponseWriter, r *http.Request, count int64) {
	JSON(w, r, CountResponse{Count: count})
}
