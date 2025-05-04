package handlers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/Zapharaos/fihub-backend/cmd/api/app/clients"
	"github.com/Zapharaos/fihub-backend/cmd/api/app/handlers"
	"github.com/Zapharaos/fihub-backend/gen/go/userpb"
	"github.com/Zapharaos/fihub-backend/internal/models"
	"github.com/Zapharaos/fihub-backend/test/mocks"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestCreateUser tests the CreateUser handler
func TestCreateUser(t *testing.T) {
	// Prepare data
	validUser := models.UserInputCreate{
		UserInputPassword: models.UserInputPassword{
			UserWithPassword: models.UserWithPassword{
				User: models.User{
					Email: "email@test.ut",
				},
				Password: "password",
			},
			Confirmation: "password",
		},
		Checkbox: true,
	}
	validUserBody, _ := json.Marshal(validUser)
	validResponse := &userpb.CreateUserResponse{
		User: &userpb.User{
			Id: uuid.New().String(),
		},
	}

	// Test cases
	tests := []struct {
		name           string
		user           []byte
		mockSetup      func(ctrl *gomock.Controller)
		expectedStatus int
	}{
		{
			name: "Fails to decode",
			user: []byte(`invalid json`),
			mockSetup: func(ctrl *gomock.Controller) {
				uc := mocks.NewMockUserServiceClient(ctrl)
				uc.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(0)
				clients.ReplaceGlobals(clients.NewClients(
					clients.WithUserClient(uc),
				))
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Fail at create",
			user: validUserBody,
			mockSetup: func(ctrl *gomock.Controller) {
				uc := mocks.NewMockUserServiceClient(ctrl)
				uc.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(nil, errors.New("error"))
				clients.ReplaceGlobals(clients.NewClients(
					clients.WithUserClient(uc),
				))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "Succeeded",
			user: validUserBody,
			mockSetup: func(ctrl *gomock.Controller) {
				uc := mocks.NewMockUserServiceClient(ctrl)
				uc.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(validResponse, nil)
				clients.ReplaceGlobals(clients.NewClients(
					clients.WithUserClient(uc),
				))
			},
			expectedStatus: http.StatusOK,
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiBasePath := viper.GetString("API_BASE_PATH")
			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", apiBasePath+"/auth/register", bytes.NewBuffer(tt.user))

			ctrl := gomock.NewController(t)
			tt.mockSetup(ctrl)
			defer ctrl.Finish()

			handlers.CreateUser(w, r)
			response := w.Result()
			defer response.Body.Close()

			assert.Equal(t, tt.expectedStatus, response.StatusCode)
		})
	}
}

// TestGetUser tests the GetUser handler
func TestGetUser(t *testing.T) {
	validResponse := &userpb.GetUserResponse{
		User: &userpb.User{
			Id: uuid.New().String(),
		},
	}

	// Test cases
	tests := []struct {
		name           string
		mockSetup      func(ctrl *gomock.Controller)
		expectedStatus int
	}{
		{
			name: "Fails to parse param",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockApiUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), "id").Return(uuid.Nil, false)
				handlers.ReplaceGlobals(m)
				uc := mocks.NewMockUserServiceClient(ctrl)
				uc.EXPECT().GetUser(gomock.Any(), gomock.Any()).Times(0)
				clients.ReplaceGlobals(clients.NewClients(
					clients.WithUserClient(uc),
				))
			},
			expectedStatus: http.StatusOK, // should be http.StatusBadRequest, but not with mock
		},
		{
			name: "Fails to retrieve the user",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockApiUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), "id").Return(uuid.New(), true)
				handlers.ReplaceGlobals(m)
				uc := mocks.NewMockUserServiceClient(ctrl)
				uc.EXPECT().GetUser(gomock.Any(), gomock.Any()).Return(nil, errors.New("error"))
				clients.ReplaceGlobals(clients.NewClients(
					clients.WithUserClient(uc),
				))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "Succeeded",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockApiUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), "id").Return(uuid.New(), true)
				handlers.ReplaceGlobals(m)
				uc := mocks.NewMockUserServiceClient(ctrl)
				uc.EXPECT().GetUser(gomock.Any(), gomock.Any()).Return(validResponse, nil)
				clients.ReplaceGlobals(clients.NewClients(
					clients.WithUserClient(uc),
				))
			},
			expectedStatus: http.StatusOK,
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiBasePath := viper.GetString("API_BASE_PATH")
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", apiBasePath+"/users/{id}", nil)

			// Apply mocks
			ctrl := gomock.NewController(t)
			tt.mockSetup(ctrl)
			defer ctrl.Finish()

			handlers.GetUser(w, r)
			response := w.Result()
			defer response.Body.Close()

			assert.Equal(t, tt.expectedStatus, response.StatusCode)
		})
	}
}

// TestGetUserSelf tests the GetUserSelf handler
func TestGetUserSelf(t *testing.T) {
	validResponse := &userpb.GetUserResponse{
		User: &userpb.User{
			Id: uuid.New().String(),
		},
	}

	// Test cases
	tests := []struct {
		name           string
		mockSetup      func(ctrl *gomock.Controller)
		expectedStatus int
	}{
		{
			name: "Fails to retrieve from context",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockApiUtils(ctrl)
				m.EXPECT().GetUserIDFromContext(gomock.Any()).Return(uuid.Nil.String(), false)
				handlers.ReplaceGlobals(m)
				uc := mocks.NewMockUserServiceClient(ctrl)
				uc.EXPECT().GetUser(gomock.Any(), gomock.Any()).Times(0)
				clients.ReplaceGlobals(clients.NewClients(
					clients.WithUserClient(uc),
				))
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Fails to retrieve user",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockApiUtils(ctrl)
				m.EXPECT().GetUserIDFromContext(gomock.Any()).Return(uuid.New().String(), true)
				handlers.ReplaceGlobals(m)
				uc := mocks.NewMockUserServiceClient(ctrl)
				uc.EXPECT().GetUser(gomock.Any(), gomock.Any()).Return(nil, errors.New("error"))
				clients.ReplaceGlobals(clients.NewClients(
					clients.WithUserClient(uc),
				))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "Succeeded",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockApiUtils(ctrl)
				m.EXPECT().GetUserIDFromContext(gomock.Any()).Return(uuid.New().String(), true)
				handlers.ReplaceGlobals(m)
				uc := mocks.NewMockUserServiceClient(ctrl)
				uc.EXPECT().GetUser(gomock.Any(), gomock.Any()).Return(validResponse, nil)
				clients.ReplaceGlobals(clients.NewClients(
					clients.WithUserClient(uc),
				))
			},
			expectedStatus: http.StatusOK,
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiBasePath := viper.GetString("API_BASE_PATH")
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", apiBasePath+"/user/me", nil)

			// Apply mocks
			ctrl := gomock.NewController(t)
			tt.mockSetup(ctrl)
			defer ctrl.Finish()

			handlers.GetUserSelf(w, r)
			response := w.Result()
			defer response.Body.Close()

			assert.Equal(t, tt.expectedStatus, response.StatusCode)
		})
	}
}

// TestUpdateUserSelf tests the UpdateUserSelf handler
func TestUpdateUserSelf(t *testing.T) {
	// Prepare data
	validUser := models.User{
		Email: "email@test.ut",
	}
	validUserBody, _ := json.Marshal(validUser)
	validResponse := &userpb.UpdateUserResponse{
		User: &userpb.User{
			Id: uuid.New().String(),
		},
	}

	// Test cases
	tests := []struct {
		name           string
		user           []byte
		mockSetup      func(ctrl *gomock.Controller)
		expectedStatus int
	}{
		{
			name: "Fails to retrieve from context",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockApiUtils(ctrl)
				m.EXPECT().GetUserIDFromContext(gomock.Any()).Return(uuid.Nil.String(), false)
				handlers.ReplaceGlobals(m)
				uc := mocks.NewMockUserServiceClient(ctrl)
				uc.EXPECT().UpdateUser(gomock.Any(), gomock.Any()).Times(0)
				clients.ReplaceGlobals(clients.NewClients(
					clients.WithUserClient(uc),
				))
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Fails to decode",
			user: []byte(`invalid json`),
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockApiUtils(ctrl)
				m.EXPECT().GetUserIDFromContext(gomock.Any()).Return(uuid.New().String(), true)
				handlers.ReplaceGlobals(m)
				uc := mocks.NewMockUserServiceClient(ctrl)
				uc.EXPECT().UpdateUser(gomock.Any(), gomock.Any()).Times(0)
				clients.ReplaceGlobals(clients.NewClients(
					clients.WithUserClient(uc),
				))
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Fail at update",
			user: validUserBody,
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockApiUtils(ctrl)
				m.EXPECT().GetUserIDFromContext(gomock.Any()).Return(uuid.New().String(), true)
				handlers.ReplaceGlobals(m)
				uc := mocks.NewMockUserServiceClient(ctrl)
				uc.EXPECT().UpdateUser(gomock.Any(), gomock.Any()).Return(nil, errors.New("error"))
				clients.ReplaceGlobals(clients.NewClients(
					clients.WithUserClient(uc),
				))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "Succeeded",
			user: validUserBody,
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockApiUtils(ctrl)
				m.EXPECT().GetUserIDFromContext(gomock.Any()).Return(uuid.New().String(), true)
				handlers.ReplaceGlobals(m)
				uc := mocks.NewMockUserServiceClient(ctrl)
				uc.EXPECT().UpdateUser(gomock.Any(), gomock.Any()).Return(validResponse, nil)
				clients.ReplaceGlobals(clients.NewClients(
					clients.WithUserClient(uc),
				))
			},
			expectedStatus: http.StatusOK,
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiBasePath := viper.GetString("API_BASE_PATH")
			w := httptest.NewRecorder()
			r := httptest.NewRequest("PUT", apiBasePath+"/user/me", bytes.NewBuffer(tt.user))

			// Apply mocks
			ctrl := gomock.NewController(t)
			tt.mockSetup(ctrl)
			defer ctrl.Finish()

			handlers.UpdateUserSelf(w, r)
			response := w.Result()
			defer response.Body.Close()

			assert.Equal(t, tt.expectedStatus, response.StatusCode)
		})
	}
}

// TestUpdateUserPassword tests the UpdateUserPassword handler
func TestUpdateUserPassword(t *testing.T) {
	// Prepare data
	validPassword := models.UserInputPassword{
		UserWithPassword: models.UserWithPassword{
			User: models.User{
				Email: "email@test.ut",
			},
			Password: "password",
		},
		Confirmation: "password",
	}
	validPasswordBody, _ := json.Marshal(validPassword)

	// Test cases
	tests := []struct {
		name           string
		password       []byte
		mockSetup      func(ctrl *gomock.Controller)
		expectedStatus int
	}{
		{
			name: "Fails to retrieve from context",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockApiUtils(ctrl)
				m.EXPECT().GetUserIDFromContext(gomock.Any()).Return(uuid.Nil.String(), false)
				handlers.ReplaceGlobals(m)
				uc := mocks.NewMockUserServiceClient(ctrl)
				uc.EXPECT().UpdateUserPassword(gomock.Any(), gomock.Any()).Times(0)
				clients.ReplaceGlobals(clients.NewClients(
					clients.WithUserClient(uc),
				))
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:     "Fails to decode",
			password: []byte(`invalid json`),
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockApiUtils(ctrl)
				m.EXPECT().GetUserIDFromContext(gomock.Any()).Return(uuid.New().String(), true)
				handlers.ReplaceGlobals(m)
				uc := mocks.NewMockUserServiceClient(ctrl)
				uc.EXPECT().UpdateUserPassword(gomock.Any(), gomock.Any()).Times(0)
				clients.ReplaceGlobals(clients.NewClients(
					clients.WithUserClient(uc),
				))
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:     "Fail at update password",
			password: validPasswordBody,
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockApiUtils(ctrl)
				m.EXPECT().GetUserIDFromContext(gomock.Any()).Return(uuid.New().String(), true)
				handlers.ReplaceGlobals(m)
				uc := mocks.NewMockUserServiceClient(ctrl)
				uc.EXPECT().UpdateUserPassword(gomock.Any(), gomock.Any()).Return(nil, errors.New("error"))
				clients.ReplaceGlobals(clients.NewClients(
					clients.WithUserClient(uc),
				))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:     "Succeeded",
			password: validPasswordBody,
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockApiUtils(ctrl)
				m.EXPECT().GetUserIDFromContext(gomock.Any()).Return(uuid.New().String(), true)
				handlers.ReplaceGlobals(m)
				uc := mocks.NewMockUserServiceClient(ctrl)
				uc.EXPECT().UpdateUserPassword(gomock.Any(), gomock.Any()).Return(&userpb.UpdateUserPasswordResponse{
					Success: true,
				}, nil)
				clients.ReplaceGlobals(clients.NewClients(
					clients.WithUserClient(uc),
				))
			},
			expectedStatus: http.StatusOK,
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiBasePath := viper.GetString("API_BASE_PATH")
			w := httptest.NewRecorder()
			r := httptest.NewRequest("PUT", apiBasePath+"/user/me/password", bytes.NewBuffer(tt.password))

			// Apply mocks
			ctrl := gomock.NewController(t)
			tt.mockSetup(ctrl)
			defer ctrl.Finish()

			handlers.UpdateUserPassword(w, r)
			response := w.Result()
			defer response.Body.Close()

			assert.Equal(t, tt.expectedStatus, response.StatusCode)
		})
	}
}

// TestDeleteUserSelf tests the DeleteUserSelf handler
func TestDeleteUserSelf(t *testing.T) {
	// Test cases
	tests := []struct {
		name           string
		mockSetup      func(ctrl *gomock.Controller)
		expectedStatus int
	}{
		{
			name: "Fails to retrieve from context",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockApiUtils(ctrl)
				m.EXPECT().GetUserIDFromContext(gomock.Any()).Return(uuid.Nil.String(), false)
				handlers.ReplaceGlobals(m)
				uc := mocks.NewMockUserServiceClient(ctrl)
				uc.EXPECT().DeleteUser(gomock.Any(), gomock.Any()).Times(0)
				clients.ReplaceGlobals(clients.NewClients(
					clients.WithUserClient(uc),
				))
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Fail at delete",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockApiUtils(ctrl)
				m.EXPECT().GetUserIDFromContext(gomock.Any()).Return(uuid.New().String(), true)
				handlers.ReplaceGlobals(m)
				uc := mocks.NewMockUserServiceClient(ctrl)
				uc.EXPECT().DeleteUser(gomock.Any(), gomock.Any()).Return(nil, errors.New("error"))
				clients.ReplaceGlobals(clients.NewClients(
					clients.WithUserClient(uc),
				))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "Succeeded",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockApiUtils(ctrl)
				m.EXPECT().GetUserIDFromContext(gomock.Any()).Return(uuid.New().String(), true)
				handlers.ReplaceGlobals(m)
				uc := mocks.NewMockUserServiceClient(ctrl)
				uc.EXPECT().DeleteUser(gomock.Any(), gomock.Any()).Return(&userpb.DeleteUserResponse{
					Success: true,
				}, nil)
				clients.ReplaceGlobals(clients.NewClients(
					clients.WithUserClient(uc),
				))
			},
			expectedStatus: http.StatusOK,
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiBasePath := viper.GetString("API_BASE_PATH")
			w := httptest.NewRecorder()
			r := httptest.NewRequest("DELETE", apiBasePath+"/user/me", nil)

			// Apply mocks
			ctrl := gomock.NewController(t)
			tt.mockSetup(ctrl)
			defer ctrl.Finish()

			handlers.DeleteUserSelf(w, r)
			response := w.Result()
			defer response.Body.Close()

			assert.Equal(t, tt.expectedStatus, response.StatusCode)
		})
	}
}
