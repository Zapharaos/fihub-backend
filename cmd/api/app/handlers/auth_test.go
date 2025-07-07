package handlers_test

import (
	"bytes"
	"encoding/json"
	"github.com/Zapharaos/fihub-backend/cmd/api/app/clients"
	"github.com/Zapharaos/fihub-backend/cmd/api/app/handlers"
	"github.com/Zapharaos/fihub-backend/gen/go/authpb"
	"github.com/Zapharaos/fihub-backend/internal/models"
	"github.com/Zapharaos/fihub-backend/test/mocks"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golang.org/x/text/language"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
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

func TestGenerateOTP(t *testing.T) {
	// Prepare data
	validLanguageTag := language.English
	emptyRequest := models.RequestUserOtp{
		Email: "",
	}
	validRequest := models.RequestUserOtp{
		Email: "email",
	}
	emptyRequestBody, _ := json.Marshal(emptyRequest)
	validRequestBody, _ := json.Marshal(validRequest)

	validResponse := &authpb.GenerateOTPResponse{
		Identifier: "identifier",
		ExpiresAt:  timestamppb.New(time.Unix(1717171717, 0)),
		Error:      nil,
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
				ac.EXPECT().GenerateOTP(gomock.Any(), gomock.Any()).Times(0)
				clients.ReplaceGlobals(clients.NewClients(
					clients.WithAuthClient(ac),
				))
				u := mocks.NewMockApiUtils(ctrl)
				u.EXPECT().ParseParamLanguage(gomock.Any(), gomock.Any()).Times(0)
				handlers.ReplaceGlobals(u)
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "empty request - fails to retrieve userID from context",
			body: emptyRequestBody,
			mockSetup: func(ctrl *gomock.Controller) {
				ac := mocks.NewMockAuthServiceClient(ctrl)
				ac.EXPECT().GenerateOTP(gomock.Any(), gomock.Any()).Times(0)
				clients.ReplaceGlobals(clients.NewClients(
					clients.WithAuthClient(ac),
				))
				u := mocks.NewMockApiUtils(ctrl)
				u.EXPECT().ParseParamLanguage(gomock.Any(), gomock.Any()).Return(validLanguageTag)
				u.EXPECT().GetUserIDFromContext(gomock.Any()).Return("", false)
				handlers.ReplaceGlobals(u)
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "fails to generate token",
			body: validRequestBody,
			mockSetup: func(ctrl *gomock.Controller) {
				ac := mocks.NewMockAuthServiceClient(ctrl)
				ac.EXPECT().GenerateOTP(gomock.Any(), gomock.Any()).Return(nil, status.Error(codes.Unknown, "error"))
				clients.ReplaceGlobals(clients.NewClients(
					clients.WithAuthClient(ac),
				))
				u := mocks.NewMockApiUtils(ctrl)
				u.EXPECT().ParseParamLanguage(gomock.Any(), gomock.Any()).Return(validLanguageTag)
				handlers.ReplaceGlobals(u)
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "succeeded - empty request must find userID in context",
			body: emptyRequestBody,
			mockSetup: func(ctrl *gomock.Controller) {
				ac := mocks.NewMockAuthServiceClient(ctrl)
				ac.EXPECT().GenerateOTP(gomock.Any(), gomock.Any()).Return(validResponse, nil)
				clients.ReplaceGlobals(clients.NewClients(
					clients.WithAuthClient(ac),
				))
				u := mocks.NewMockApiUtils(ctrl)
				u.EXPECT().ParseParamLanguage(gomock.Any(), gomock.Any()).Return(validLanguageTag)
				u.EXPECT().GetUserIDFromContext(gomock.Any()).Return("user-id", true)
				handlers.ReplaceGlobals(u)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "succeeded - default request",
			body: validRequestBody,
			mockSetup: func(ctrl *gomock.Controller) {
				ac := mocks.NewMockAuthServiceClient(ctrl)
				ac.EXPECT().GenerateOTP(gomock.Any(), gomock.Any()).Return(validResponse, nil)
				clients.ReplaceGlobals(clients.NewClients(
					clients.WithAuthClient(ac),
				))
				u := mocks.NewMockApiUtils(ctrl)
				u.EXPECT().ParseParamLanguage(gomock.Any(), gomock.Any()).Return(validLanguageTag)
				handlers.ReplaceGlobals(u)
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiBasePath := viper.GetString("API_BASE_PATH")
			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", apiBasePath+"/auth/register/otp", bytes.NewBuffer(tt.body))

			// Apply mocks
			ctrl := gomock.NewController(t)
			tt.mockSetup(ctrl)
			defer ctrl.Finish()

			handlers.GenerateSignupOTP(w, r)
			response := w.Result()
			defer response.Body.Close()

			assert.Equal(t, tt.expectedStatus, response.StatusCode)
		})
	}
}

func TestValidateOTP(t *testing.T) {
	// Prepare data
	validRequest := models.ValidateUserOtp{
		Identifier: "identifier",
		Otp:        "otp-code",
	}
	validRequestBody, _ := json.Marshal(validRequest)

	validResponse := &authpb.ValidateOTPResponse{
		RequestId: "request-id",
		ExpiresAt: timestamppb.New(time.Unix(1717171717, 0)),
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
				ac.EXPECT().ValidateOTP(gomock.Any(), gomock.Any()).Times(0)
				clients.ReplaceGlobals(clients.NewClients(
					clients.WithAuthClient(ac),
				))
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "fails to generate token",
			body: validRequestBody,
			mockSetup: func(ctrl *gomock.Controller) {
				ac := mocks.NewMockAuthServiceClient(ctrl)
				ac.EXPECT().ValidateOTP(gomock.Any(), gomock.Any()).Return(nil, status.Error(codes.Unknown, "error"))
				clients.ReplaceGlobals(clients.NewClients(
					clients.WithAuthClient(ac),
				))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "succeeded",
			body: validRequestBody,
			mockSetup: func(ctrl *gomock.Controller) {
				ac := mocks.NewMockAuthServiceClient(ctrl)
				ac.EXPECT().ValidateOTP(gomock.Any(), gomock.Any()).Return(validResponse, nil)
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
			r := httptest.NewRequest("POST", apiBasePath+"/auth/register/otp/validate", bytes.NewBuffer(tt.body))

			// Apply mocks
			ctrl := gomock.NewController(t)
			tt.mockSetup(ctrl)
			defer ctrl.Finish()

			handlers.ValidateSignupOTP(w, r)
			response := w.Result()
			defer response.Body.Close()

			assert.Equal(t, tt.expectedStatus, response.StatusCode)
		})
	}
}

func TestRegisterUser(t *testing.T) {
	// Prepare data
	validRequest := models.UserInputCreate{
		OtpRequestID: uuid.New(),
	}
	validRequestBody, _ := json.Marshal(validRequest)

	validResponse := &authpb.CreateUserResponse{
		Success: true,
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
				ac.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(0)
				clients.ReplaceGlobals(clients.NewClients(
					clients.WithAuthClient(ac),
				))
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "fails to generate token",
			body: validRequestBody,
			mockSetup: func(ctrl *gomock.Controller) {
				ac := mocks.NewMockAuthServiceClient(ctrl)
				ac.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(nil, status.Error(codes.Unknown, "error"))
				clients.ReplaceGlobals(clients.NewClients(
					clients.WithAuthClient(ac),
				))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "succeeded",
			body: validRequestBody,
			mockSetup: func(ctrl *gomock.Controller) {
				ac := mocks.NewMockAuthServiceClient(ctrl)
				ac.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(validResponse, nil)
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
			r := httptest.NewRequest("POST", apiBasePath+"/auth/register", bytes.NewBuffer(tt.body))

			// Apply mocks
			ctrl := gomock.NewController(t)
			tt.mockSetup(ctrl)
			defer ctrl.Finish()

			handlers.RegisterUser(w, r)
			response := w.Result()
			defer response.Body.Close()

			assert.Equal(t, tt.expectedStatus, response.StatusCode)
		})
	}
}

func TestSubmitChangePassword(t *testing.T) {
	// Prepare data
	validRequest := models.UserInputChangePassword{
		OtpRequestID: uuid.New(),
	}
	validRequestBody, _ := json.Marshal(validRequest)

	validResponse := &authpb.UpdatePasswordResponse{
		Success: true,
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
				ac.EXPECT().UpdatePassword(gomock.Any(), gomock.Any()).Times(0)
				clients.ReplaceGlobals(clients.NewClients(
					clients.WithAuthClient(ac),
				))
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "fails to generate token",
			body: validRequestBody,
			mockSetup: func(ctrl *gomock.Controller) {
				ac := mocks.NewMockAuthServiceClient(ctrl)
				ac.EXPECT().UpdatePassword(gomock.Any(), gomock.Any()).Return(nil, status.Error(codes.Unknown, "error"))
				clients.ReplaceGlobals(clients.NewClients(
					clients.WithAuthClient(ac),
				))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "succeeded",
			body: validRequestBody,
			mockSetup: func(ctrl *gomock.Controller) {
				ac := mocks.NewMockAuthServiceClient(ctrl)
				ac.EXPECT().UpdatePassword(gomock.Any(), gomock.Any()).Return(validResponse, nil)
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
			r := httptest.NewRequest("POST", apiBasePath+"/user/me/password", bytes.NewBuffer(tt.body))

			// Apply mocks
			ctrl := gomock.NewController(t)
			tt.mockSetup(ctrl)
			defer ctrl.Finish()

			handlers.SubmitChangePassword(w, r)
			response := w.Result()
			defer response.Body.Close()

			assert.Equal(t, tt.expectedStatus, response.StatusCode)
		})
	}
}

func TestResetForgottenPassword(t *testing.T) {
	// Prepare data
	validRequest := models.UserInputChangePassword{
		OtpRequestID: uuid.New(),
	}
	validRequestBody, _ := json.Marshal(validRequest)

	validResponse := &authpb.ResetForgottenPasswordResponse{
		Success: true,
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
				ac.EXPECT().ResetForgottenPassword(gomock.Any(), gomock.Any()).Times(0)
				clients.ReplaceGlobals(clients.NewClients(
					clients.WithAuthClient(ac),
				))
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "fails to generate token",
			body: validRequestBody,
			mockSetup: func(ctrl *gomock.Controller) {
				ac := mocks.NewMockAuthServiceClient(ctrl)
				ac.EXPECT().ResetForgottenPassword(gomock.Any(), gomock.Any()).Return(nil, status.Error(codes.Unknown, "error"))
				clients.ReplaceGlobals(clients.NewClients(
					clients.WithAuthClient(ac),
				))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "succeeded",
			body: validRequestBody,
			mockSetup: func(ctrl *gomock.Controller) {
				ac := mocks.NewMockAuthServiceClient(ctrl)
				ac.EXPECT().ResetForgottenPassword(gomock.Any(), gomock.Any()).Return(validResponse, nil)
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
			r := httptest.NewRequest("POST", apiBasePath+"/auth/password/reset", bytes.NewBuffer(tt.body))

			// Apply mocks
			ctrl := gomock.NewController(t)
			tt.mockSetup(ctrl)
			defer ctrl.Finish()

			handlers.ResetForgottenPassword(w, r)
			response := w.Result()
			defer response.Body.Close()

			assert.Equal(t, tt.expectedStatus, response.StatusCode)
		})
	}
}
