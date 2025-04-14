package handlers

import (
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestHealthCheckHandler tests the HealthCheckHandler function
func TestHealthCheckHandler(t *testing.T) {
	apiBasePath := viper.GetString("API_BASE_PATH")
	w := httptest.NewRecorder()
	r := httptest.NewRequest("PUT", apiBasePath+"/health", nil)

	HealthCheckHandler(w, r)
	response := w.Result()
	defer response.Body.Close()

	assert.Equal(t, http.StatusOK, response.StatusCode)
}
