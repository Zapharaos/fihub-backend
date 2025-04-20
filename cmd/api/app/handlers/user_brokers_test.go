package handlers_test

import (
	"bytes"
	"encoding/json"
	"github.com/Zapharaos/fihub-backend/cmd/api/app/clients"
	"github.com/Zapharaos/fihub-backend/cmd/api/app/handlers"
	"github.com/Zapharaos/fihub-backend/internal/models"
	"github.com/Zapharaos/fihub-backend/protogen"
	"github.com/Zapharaos/fihub-backend/test/mocks"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestCreateUserBroker tests the CreateUserBroker handler
func TestCreateUserBroker(t *testing.T) {
	// Prepare data
	validBroker := models.BrokerUserInput{
		BrokerID: uuid.New().String(),
	}
	validBrokerBody, _ := json.Marshal(validBroker)
	validResponse := &protogen.CreateBrokerUserResponse{
		UserBrokers: []*protogen.BrokerUser{},
	}

	// Test cases
	tests := []struct {
		name           string
		broker         []byte
		mockSetup      func(ctrl *gomock.Controller)
		expectedStatus int
	}{
		{
			name: "Fails to retrieve user from context",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().GetUserFromContext(gomock.Any()).Return(models.UserWithRoles{}, false)
				handlers.ReplaceGlobals(m)
				bc := mocks.NewBrokerServiceClient(ctrl)
				bc.EXPECT().CreateBrokerUser(gomock.Any(), gomock.Any()).Times(0)
				clients.ReplaceGlobals(clients.NewClients(nil, bc, nil))
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:   "Fails to decode",
			broker: []byte(`invalid json`),
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().GetUserFromContext(gomock.Any()).Return(models.UserWithRoles{}, true)
				handlers.ReplaceGlobals(m)
				bc := mocks.NewBrokerServiceClient(ctrl)
				bc.EXPECT().CreateBrokerUser(gomock.Any(), gomock.Any()).Times(0)
				clients.ReplaceGlobals(clients.NewClients(nil, bc, nil))
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "Fails to create user broker",
			broker: validBrokerBody,
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().GetUserFromContext(gomock.Any()).Return(models.UserWithRoles{}, true)
				handlers.ReplaceGlobals(m)
				bc := mocks.NewBrokerServiceClient(ctrl)
				bc.EXPECT().CreateBrokerUser(gomock.Any(), gomock.Any()).Return(nil, status.Error(codes.Unknown, "error"))
				clients.ReplaceGlobals(clients.NewClients(nil, bc, nil))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:   "Succeeded",
			broker: validBrokerBody,
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().GetUserFromContext(gomock.Any()).Return(models.UserWithRoles{}, true)
				handlers.ReplaceGlobals(m)
				bc := mocks.NewBrokerServiceClient(ctrl)
				bc.EXPECT().CreateBrokerUser(gomock.Any(), gomock.Any()).Return(validResponse, nil)
				clients.ReplaceGlobals(clients.NewClients(nil, bc, nil))
			},
			expectedStatus: http.StatusOK,
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiBasePath := viper.GetString("API_BASE_PATH")
			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", apiBasePath+"/users/brokers", bytes.NewBuffer(tt.broker))

			ctrl := gomock.NewController(t)
			tt.mockSetup(ctrl)
			defer ctrl.Finish()

			handlers.CreateUserBroker(w, r)
			response := w.Result()
			defer response.Body.Close()

			assert.Equal(t, tt.expectedStatus, response.StatusCode)
		})
	}
}

// TestDeleteUserBroker tests the DeleteUserBroker handler
func TestDeleteUserBroker(t *testing.T) {
	// Test cases
	tests := []struct {
		name           string
		mockSetup      func(ctrl *gomock.Controller)
		expectedStatus int
	}{
		{
			name: "Fails to parse param",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.Nil, false)
				m.EXPECT().GetUserFromContext(gomock.Any()).Times(0)
				handlers.ReplaceGlobals(m)
			},
			expectedStatus: http.StatusOK, // should be http.StatusBadRequest, but not with mock
		},
		{
			name: "Fails to retrieve user from context",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.Nil, true)
				m.EXPECT().GetUserFromContext(gomock.Any()).Return(models.UserWithRoles{}, false)
				handlers.ReplaceGlobals(m)
				bc := mocks.NewBrokerServiceClient(ctrl)
				bc.EXPECT().DeleteBrokerUser(gomock.Any(), gomock.Any()).Times(0)
				clients.ReplaceGlobals(clients.NewClients(nil, bc, nil))
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Fails to delete the user broker",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.Nil, true)
				m.EXPECT().GetUserFromContext(gomock.Any()).Return(models.UserWithRoles{}, true)
				handlers.ReplaceGlobals(m)
				bc := mocks.NewBrokerServiceClient(ctrl)
				bc.EXPECT().DeleteBrokerUser(gomock.Any(), gomock.Any()).Return(nil, status.Error(codes.Unknown, "error"))
				clients.ReplaceGlobals(clients.NewClients(nil, bc, nil))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "Succeeded",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.Nil, true)
				m.EXPECT().GetUserFromContext(gomock.Any()).Return(models.UserWithRoles{}, true)
				handlers.ReplaceGlobals(m)
				bc := mocks.NewBrokerServiceClient(ctrl)
				bc.EXPECT().DeleteBrokerUser(gomock.Any(), gomock.Any()).Return(&protogen.DeleteBrokerUserResponse{}, nil)
				clients.ReplaceGlobals(clients.NewClients(nil, bc, nil))
			},
			expectedStatus: http.StatusOK,
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiBasePath := viper.GetString("API_BASE_PATH")
			w := httptest.NewRecorder()
			r := httptest.NewRequest("DELETE", apiBasePath+"/users/brokers/"+uuid.New().String(), nil)

			// Apply mocks
			ctrl := gomock.NewController(t)
			tt.mockSetup(ctrl)
			defer ctrl.Finish()

			handlers.DeleteUserBroker(w, r)
			response := w.Result()
			defer response.Body.Close()

			assert.Equal(t, tt.expectedStatus, response.StatusCode)
		})
	}
}

// TestListUserBrokers tests the ListUserBrokers handler
func TestListUserBrokers(t *testing.T) {
	// Test cases
	tests := []struct {
		name           string
		mockSetup      func(ctrl *gomock.Controller)
		expectedStatus int
	}{
		{
			name: "Fails to retrieve user from context",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().GetUserFromContext(gomock.Any()).Return(models.UserWithRoles{}, false)
				handlers.ReplaceGlobals(m)
				bc := mocks.NewBrokerServiceClient(ctrl)
				bc.EXPECT().ListUserBrokers(gomock.Any(), gomock.Any()).Times(0)
				clients.ReplaceGlobals(clients.NewClients(nil, bc, nil))
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Fails to retrieve all user brokers",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().GetUserFromContext(gomock.Any()).Return(models.UserWithRoles{}, true)
				handlers.ReplaceGlobals(m)
				bc := mocks.NewBrokerServiceClient(ctrl)
				bc.EXPECT().ListUserBrokers(gomock.Any(), gomock.Any()).Return(nil, status.Error(codes.Unknown, "error"))
				clients.ReplaceGlobals(clients.NewClients(nil, bc, nil))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "Succeeded",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockUtils(ctrl)
				m.EXPECT().GetUserFromContext(gomock.Any()).Return(models.UserWithRoles{}, true)
				handlers.ReplaceGlobals(m)
				bc := mocks.NewBrokerServiceClient(ctrl)
				bc.EXPECT().ListUserBrokers(gomock.Any(), gomock.Any()).Return(&protogen.ListUserBrokersResponse{
					UserBrokers: []*protogen.BrokerUser{},
				}, nil)
				clients.ReplaceGlobals(clients.NewClients(nil, bc, nil))
			},
			expectedStatus: http.StatusOK,
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiBasePath := viper.GetString("API_BASE_PATH")
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", apiBasePath+"/users/brokers", nil)

			// Apply mocks
			ctrl := gomock.NewController(t)
			tt.mockSetup(ctrl)
			defer ctrl.Finish()

			handlers.ListUserBrokers(w, r)
			response := w.Result()
			defer response.Body.Close()

			assert.Equal(t, tt.expectedStatus, response.StatusCode)
		})
	}
}
