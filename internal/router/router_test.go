package router

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewRouter(t *testing.T) {
	r := New()

	assert.NotNil(t, r, "Router should not be nil")

	// Test health check route
	req, _ := http.NewRequest("GET", "/api/v1/health", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "Health check route should return status OK")
}
