package handlers

import (
	"github.com/Zapharaos/fihub-backend/cmd/api/app/clients"
	"github.com/Zapharaos/fihub-backend/protogen"
	"github.com/Zapharaos/fihub-backend/test/mocks"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestHealthCheckHandler tests the HealthCheckHandler function
func TestHealthCheckHandler(t *testing.T) {
	apiBasePath := viper.GetString("API_BASE_PATH")
	w := httptest.NewRecorder()
	r := httptest.NewRequest("PUT", apiBasePath+"/health", nil)

	// Set up a mock controller and a mock HealthServiceClient
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	hc := mocks.NewHealthServiceClient(ctrl)
	hc.EXPECT().CheckHealth(gomock.Any(), gomock.Any()).Return(&protogen.HealthResponse{
		IsHealthy: true,
	}, nil)
	clients.ReplaceGlobals(clients.NewClients(
		clients.WithHealthClient(hc),
	))

	// Call the HealthCheckHandler function
	HealthCheckHandler(w, r)
	response := w.Result()
	defer response.Body.Close()

	assert.Equal(t, http.StatusOK, response.StatusCode)
}
