package router

import (
	"github.com/Zapharaos/fihub-backend/cmd/api/app/clients"
	"github.com/Zapharaos/fihub-backend/cmd/api/app/server"
	"github.com/Zapharaos/fihub-backend/gen/go/healthpb"
	"github.com/Zapharaos/fihub-backend/test/mocks"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestNewRouter tests the New function to ensure it creates a new router instance.
// Verifies that the health check route is working correctly.
func TestNewRouter(t *testing.T) {
	viper.Set("API_BASE_PATH", "/api/v1")
	r := New(server.Config{
		CORS:        true,
		Security:    true,
		GatewayMode: true,
		AllowOrigin: "https://*,http://*",
	})

	// Assert that the router instance is not nil
	assert.NotNil(t, r, "Router should not be nil")

	// Set up a mock controller and a mock HealthServiceClient
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	hc := mocks.NewMockHealthServiceClient(ctrl)
	hc.EXPECT().CheckHealth(gomock.Any(), gomock.Any()).Return(&healthpb.HealthResponse{
		IsHealthy: true,
	}, nil)
	clients.ReplaceGlobals(clients.NewClients(
		clients.WithHealthClient(hc),
	))

	// Test health check route
	apiBasePath := viper.GetString("API_BASE_PATH")
	req, _ := http.NewRequest("GET", apiBasePath+"/health", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "Health check route should return status OK")
}
