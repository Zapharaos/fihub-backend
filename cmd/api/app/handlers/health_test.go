package handlers

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestHealthCheckHandler tests the HealthCheckHandler function
func TestHealthCheckHandler(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("PUT", "/api/v1/health", nil)

	HealthCheckHandler(w, r)
	response := w.Result()
	defer response.Body.Close()

	assert.Equal(t, http.StatusOK, response.StatusCode)
}
