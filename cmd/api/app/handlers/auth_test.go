package handlers_test

import (
	"bytes"
	"encoding/json"
	"github.com/Zapharaos/fihub-backend/cmd/api/app/clients"
	"github.com/Zapharaos/fihub-backend/cmd/api/app/handlers"
	"github.com/Zapharaos/fihub-backend/gen/go/authpb"
	"github.com/Zapharaos/fihub-backend/internal/models"
	"github.com/Zapharaos/fihub-backend/test/mocks"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetToken(t *testing.T) {
	// Prepare data
	validCreds := models.UserWithPassword{
		Password: "password",
		User: models.User{
			Email: "email",
		},
	}
	validCredsBody, _ := json.Marshal(validCreds)
	validResponse := &authpb.GenerateTokenResponse{
		Token: "valid-token",
	}

	tests := []struct {
		name           string
		body           []byte
		mockSetup      func(ctrl *gomock.Controller)
		expectedStatus int
	}{
		{
			name: "fails to decode",
			body: []byte("invalid"),
			mockSetup: func(ctrl *gomock.Controller) {
				ac := mocks.NewMockAuthServiceClient(ctrl)
				ac.EXPECT().GenerateToken(gomock.Any(), gomock.Any()).Times(0)
				clients.ReplaceGlobals(clients.NewClients(
					clients.WithAuthClient(ac),
				))
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "fails to generate token",
			body: validCredsBody,
			mockSetup: func(ctrl *gomock.Controller) {
				ac := mocks.NewMockAuthServiceClient(ctrl)
				ac.EXPECT().GenerateToken(gomock.Any(), gomock.Any()).Return(nil, status.Error(codes.Unknown, "error"))
				clients.ReplaceGlobals(clients.NewClients(
					clients.WithAuthClient(ac),
				))
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "succeeded",
			body: validCredsBody,
			mockSetup: func(ctrl *gomock.Controller) {
				ac := mocks.NewMockAuthServiceClient(ctrl)
				ac.EXPECT().GenerateToken(gomock.Any(), gomock.Any()).Return(validResponse, nil)
				clients.ReplaceGlobals(clients.NewClients(
					clients.WithAuthClient(ac),
				))
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiBasePath := viper.GetString("API_BASE_PATH")
			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", apiBasePath+"/auth/token", bytes.NewBuffer(tt.body))

			// Apply mocks
			ctrl := gomock.NewController(t)
			tt.mockSetup(ctrl)
			defer ctrl.Finish()

			handlers.GetToken(w, r)
			response := w.Result()
			defer response.Body.Close()

			assert.Equal(t, tt.expectedStatus, response.StatusCode)
		})
	}
}
