package handlers_test

import (
	"bytes"
	"encoding/json"
	"github.com/Zapharaos/fihub-backend/cmd/api/app/clients"
	"github.com/Zapharaos/fihub-backend/cmd/api/app/handlers"
	"github.com/Zapharaos/fihub-backend/gen/go/brokerpb"
	"github.com/Zapharaos/fihub-backend/gen/go/transactionpb"
	"github.com/Zapharaos/fihub-backend/internal/models"
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
	validResponse := &brokerpb.CreateBrokerUserResponse{
		UserBrokers: []*brokerpb.BrokerUser{},
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
				m := mocks.NewMockApiUtils(ctrl)
				m.EXPECT().GetUserIDFromContext(gomock.Any()).Return(uuid.Nil.String(), false)
				handlers.ReplaceGlobals(m)
				bc := mocks.NewMockBrokerServiceClient(ctrl)
				bc.EXPECT().CreateBrokerUser(gomock.Any(), gomock.Any()).Times(0)
				clients.ReplaceGlobals(clients.NewClients(
					clients.WithBrokerClient(bc),
				))
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:   "Fails to decode",
			broker: []byte(`invalid json`),
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockApiUtils(ctrl)
				m.EXPECT().GetUserIDFromContext(gomock.Any()).Return(uuid.New().String(), true)
				handlers.ReplaceGlobals(m)
				bc := mocks.NewMockBrokerServiceClient(ctrl)
				bc.EXPECT().CreateBrokerUser(gomock.Any(), gomock.Any()).Times(0)
				clients.ReplaceGlobals(clients.NewClients(
					clients.WithBrokerClient(bc),
				))
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "Fails to create user broker",
			broker: validBrokerBody,
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockApiUtils(ctrl)
				m.EXPECT().GetUserIDFromContext(gomock.Any()).Return(uuid.New().String(), true)
				handlers.ReplaceGlobals(m)
				bc := mocks.NewMockBrokerServiceClient(ctrl)
				bc.EXPECT().CreateBrokerUser(gomock.Any(), gomock.Any()).Return(nil, status.Error(codes.Unknown, "error"))
				clients.ReplaceGlobals(clients.NewClients(
					clients.WithBrokerClient(bc),
				))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:   "Succeeded",
			broker: validBrokerBody,
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockApiUtils(ctrl)
				m.EXPECT().GetUserIDFromContext(gomock.Any()).Return(uuid.New().String(), true)
				handlers.ReplaceGlobals(m)
				bc := mocks.NewMockBrokerServiceClient(ctrl)
				bc.EXPECT().CreateBrokerUser(gomock.Any(), gomock.Any()).Return(validResponse, nil)
				clients.ReplaceGlobals(clients.NewClients(
					clients.WithBrokerClient(bc),
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
			r := httptest.NewRequest("POST", apiBasePath+"/broker/user", bytes.NewBuffer(tt.broker))

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
				m := mocks.NewMockApiUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.Nil, false)
				m.EXPECT().GetUserIDFromContext(gomock.Any()).Times(0)
				handlers.ReplaceGlobals(m)
			},
			expectedStatus: http.StatusOK, // should be http.StatusBadRequest, but not with mock
		},
		{
			name: "Fails to retrieve user from context",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockApiUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.New(), true)
				m.EXPECT().GetUserIDFromContext(gomock.Any()).Return(uuid.Nil.String(), false)
				handlers.ReplaceGlobals(m)
				bc := mocks.NewMockBrokerServiceClient(ctrl)
				bc.EXPECT().DeleteBrokerUser(gomock.Any(), gomock.Any()).Times(0)
				clients.ReplaceGlobals(clients.NewClients(
					clients.WithBrokerClient(bc),
				))
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Fails to delete the user broker",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockApiUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.New(), true)
				m.EXPECT().GetUserIDFromContext(gomock.Any()).Return(uuid.New().String(), true)
				handlers.ReplaceGlobals(m)
				bc := mocks.NewMockBrokerServiceClient(ctrl)
				bc.EXPECT().DeleteBrokerUser(gomock.Any(), gomock.Any()).Return(nil, status.Error(codes.Unknown, "error"))
				tc := mocks.NewMockTransactionServiceClient(ctrl)
				tc.EXPECT().DeleteTransactionByBroker(gomock.Any(), gomock.Any()).Times(0)
				clients.ReplaceGlobals(clients.NewClients(
					clients.WithBrokerClient(bc),
					clients.WithTransactionClient(tc),
				))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "Fails to delete transactions related to the user broker",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockApiUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.New(), true)
				m.EXPECT().GetUserIDFromContext(gomock.Any()).Return(uuid.New().String(), true)
				handlers.ReplaceGlobals(m)
				bc := mocks.NewMockBrokerServiceClient(ctrl)
				bc.EXPECT().DeleteBrokerUser(gomock.Any(), gomock.Any()).Return(&brokerpb.DeleteBrokerUserResponse{}, nil)
				tc := mocks.NewMockTransactionServiceClient(ctrl)
				tc.EXPECT().DeleteTransactionByBroker(gomock.Any(), gomock.Any()).Return(nil, status.Error(codes.Unknown, "error"))
				clients.ReplaceGlobals(clients.NewClients(
					clients.WithBrokerClient(bc),
					clients.WithTransactionClient(tc),
				))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "Succeeded",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockApiUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.New(), true)
				m.EXPECT().GetUserIDFromContext(gomock.Any()).Return(uuid.New().String(), true)
				handlers.ReplaceGlobals(m)
				bc := mocks.NewMockBrokerServiceClient(ctrl)
				bc.EXPECT().DeleteBrokerUser(gomock.Any(), gomock.Any()).Return(&brokerpb.DeleteBrokerUserResponse{}, nil)
				tc := mocks.NewMockTransactionServiceClient(ctrl)
				tc.EXPECT().DeleteTransactionByBroker(gomock.Any(), gomock.Any()).Return(&transactionpb.DeleteTransactionByBrokerResponse{}, nil)
				clients.ReplaceGlobals(clients.NewClients(
					clients.WithBrokerClient(bc),
					clients.WithTransactionClient(tc),
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
			r := httptest.NewRequest("DELETE", apiBasePath+"/broker/"+uuid.New().String()+"/user", nil)

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
				m := mocks.NewMockApiUtils(ctrl)
				m.EXPECT().GetUserIDFromContext(gomock.Any()).Return(uuid.Nil.String(), false)
				handlers.ReplaceGlobals(m)
				bc := mocks.NewMockBrokerServiceClient(ctrl)
				bc.EXPECT().ListUserBrokers(gomock.Any(), gomock.Any()).Times(0)
				clients.ReplaceGlobals(clients.NewClients(
					clients.WithBrokerClient(bc),
				))
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Fails to retrieve all user brokers",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockApiUtils(ctrl)
				m.EXPECT().GetUserIDFromContext(gomock.Any()).Return(uuid.New().String(), true)
				handlers.ReplaceGlobals(m)
				bc := mocks.NewMockBrokerServiceClient(ctrl)
				bc.EXPECT().ListUserBrokers(gomock.Any(), gomock.Any()).Return(nil, status.Error(codes.Unknown, "error"))
				clients.ReplaceGlobals(clients.NewClients(
					clients.WithBrokerClient(bc),
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
				bc := mocks.NewMockBrokerServiceClient(ctrl)
				bc.EXPECT().ListUserBrokers(gomock.Any(), gomock.Any()).Return(&brokerpb.ListUserBrokersResponse{
					UserBrokers: []*brokerpb.BrokerUser{},
				}, nil)
				clients.ReplaceGlobals(clients.NewClients(
					clients.WithBrokerClient(bc),
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
			r := httptest.NewRequest("GET", apiBasePath+"/broker/user", nil)

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
