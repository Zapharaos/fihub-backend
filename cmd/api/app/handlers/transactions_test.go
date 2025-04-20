package handlers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/Zapharaos/fihub-backend/cmd/api/app/clients"
	"github.com/Zapharaos/fihub-backend/cmd/api/app/handlers"
	"github.com/Zapharaos/fihub-backend/cmd/broker/app/repositories"
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
	"time"
)

// TestCreateTransaction tests the CreateTransaction handler
func TestCreateTransaction(t *testing.T) {
	// Define request bodies
	validRequest := models.TransactionInput{
		ID:       uuid.New(),
		BrokerID: uuid.New(),
		Date:     time.Now().AddDate(-1, 0, 0), // 1 year in the past
		Type:     models.BUY,
		Asset:    "asset",
		Quantity: 1,
		Price:    1,
		Fee:      1,
	}
	validRequestBody, _ := json.Marshal(validRequest)
	validResponse := &protogen.CreateTransactionResponse{
		Transaction: &protogen.Transaction{
			Id:       uuid.New().String(),
			UserId:   uuid.New().String(),
			BrokerId: uuid.New().String(),
		},
	}

	// Define tests
	tests := []struct {
		name           string
		body           []byte
		mockSetup      func(ctrl *gomock.Controller)
		expectedStatus int
	}{
		{
			name: "fails to retrieve user from context",
			mockSetup: func(ctrl *gomock.Controller) {
				h := mocks.NewMockUtils(ctrl)
				h.EXPECT().GetUserFromContext(gomock.Any()).Return(models.UserWithRoles{}, false)
				handlers.ReplaceGlobals(h)
				bu := mocks.NewBrokerUserRepository(ctrl)
				bu.EXPECT().Exists(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(nil, bu, nil))
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "fails to decode",
			body: []byte("invalid json"),
			mockSetup: func(ctrl *gomock.Controller) {
				h := mocks.NewMockUtils(ctrl)
				h.EXPECT().GetUserFromContext(gomock.Any()).Return(models.UserWithRoles{}, true)
				handlers.ReplaceGlobals(h)
				bu := mocks.NewBrokerUserRepository(ctrl)
				bu.EXPECT().Exists(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(nil, bu, nil))
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "fails at user broker existence check",
			body: validRequestBody,
			mockSetup: func(ctrl *gomock.Controller) {
				h := mocks.NewMockUtils(ctrl)
				h.EXPECT().GetUserFromContext(gomock.Any()).Return(models.UserWithRoles{}, true)
				handlers.ReplaceGlobals(h)
				bu := mocks.NewBrokerUserRepository(ctrl)
				bu.EXPECT().Exists(gomock.Any()).Return(false, errors.New("error"))
				repositories.ReplaceGlobals(repositories.NewRepository(nil, bu, nil))
				tc := mocks.NewTransactionServiceClient(ctrl)
				tc.EXPECT().CreateTransaction(gomock.Any(), gomock.Any()).Times(0)
				clients.ReplaceGlobals(clients.NewClients(nil, nil, tc))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "user broker does not exist",
			body: validRequestBody,
			mockSetup: func(ctrl *gomock.Controller) {
				h := mocks.NewMockUtils(ctrl)
				h.EXPECT().GetUserFromContext(gomock.Any()).Return(models.UserWithRoles{}, true)
				handlers.ReplaceGlobals(h)
				bu := mocks.NewBrokerUserRepository(ctrl)
				bu.EXPECT().Exists(gomock.Any()).Return(false, nil)
				repositories.ReplaceGlobals(repositories.NewRepository(nil, bu, nil))
				tc := mocks.NewTransactionServiceClient(ctrl)
				tc.EXPECT().CreateTransaction(gomock.Any(), gomock.Any()).Times(0)
				clients.ReplaceGlobals(clients.NewClients(nil, nil, tc))
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "fails at transactions creation",
			body: validRequestBody,
			mockSetup: func(ctrl *gomock.Controller) {
				h := mocks.NewMockUtils(ctrl)
				h.EXPECT().GetUserFromContext(gomock.Any()).Return(models.UserWithRoles{}, true)
				handlers.ReplaceGlobals(h)
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().Get(gomock.Any()).Times(0)
				bu := mocks.NewBrokerUserRepository(ctrl)
				bu.EXPECT().Exists(gomock.Any()).Return(true, nil)
				repositories.ReplaceGlobals(repositories.NewRepository(bb, bu, nil))
				tc := mocks.NewTransactionServiceClient(ctrl)
				tc.EXPECT().CreateTransaction(gomock.Any(), gomock.Any()).Return(nil, status.Error(codes.Unknown, "error"))
				clients.ReplaceGlobals(clients.NewClients(nil, nil, tc))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "fails to retrieve broker",
			body: validRequestBody,
			mockSetup: func(ctrl *gomock.Controller) {
				h := mocks.NewMockUtils(ctrl)
				h.EXPECT().GetUserFromContext(gomock.Any()).Return(models.UserWithRoles{}, true)
				handlers.ReplaceGlobals(h)
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().Get(gomock.Any()).Return(models.Broker{}, false, errors.New("error"))
				bu := mocks.NewBrokerUserRepository(ctrl)
				bu.EXPECT().Exists(gomock.Any()).Return(true, nil)
				repositories.ReplaceGlobals(repositories.NewRepository(bb, bu, nil))
				tc := mocks.NewTransactionServiceClient(ctrl)
				tc.EXPECT().CreateTransaction(gomock.Any(), gomock.Any()).Return(validResponse, nil)
				clients.ReplaceGlobals(clients.NewClients(nil, nil, tc))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "fails to find broker",
			body: validRequestBody,
			mockSetup: func(ctrl *gomock.Controller) {
				h := mocks.NewMockUtils(ctrl)
				h.EXPECT().GetUserFromContext(gomock.Any()).Return(models.UserWithRoles{}, true)
				handlers.ReplaceGlobals(h)
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().Get(gomock.Any()).Return(models.Broker{}, false, nil)
				bu := mocks.NewBrokerUserRepository(ctrl)
				bu.EXPECT().Exists(gomock.Any()).Return(true, nil)
				repositories.ReplaceGlobals(repositories.NewRepository(bb, bu, nil))
				tc := mocks.NewTransactionServiceClient(ctrl)
				tc.EXPECT().CreateTransaction(gomock.Any(), gomock.Any()).Return(validResponse, nil)
				clients.ReplaceGlobals(clients.NewClients(nil, nil, tc))
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name: "succeeded",
			body: validRequestBody,
			mockSetup: func(ctrl *gomock.Controller) {
				h := mocks.NewMockUtils(ctrl)
				h.EXPECT().GetUserFromContext(gomock.Any()).Return(models.UserWithRoles{}, true)
				handlers.ReplaceGlobals(h)
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().Get(gomock.Any()).Return(models.Broker{}, true, nil)
				bu := mocks.NewBrokerUserRepository(ctrl)
				bu.EXPECT().Exists(gomock.Any()).Return(true, nil)
				repositories.ReplaceGlobals(repositories.NewRepository(bb, bu, nil))
				tc := mocks.NewTransactionServiceClient(ctrl)
				tc.EXPECT().CreateTransaction(gomock.Any(), gomock.Any()).Return(validResponse, nil)
				clients.ReplaceGlobals(clients.NewClients(nil, nil, tc))
			},
			expectedStatus: http.StatusOK,
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiBasePath := viper.GetString("API_BASE_PATH")
			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", apiBasePath+"/transactions", bytes.NewBuffer(tt.body))

			// Apply mocks
			ctrl := gomock.NewController(t)
			tt.mockSetup(ctrl)
			defer ctrl.Finish()

			handlers.CreateTransaction(w, r)
			response := w.Result()
			defer response.Body.Close()

			assert.Equal(t, tt.expectedStatus, response.StatusCode)
		})
	}
}

// TestGetTransaction tests the GetTransaction handler
func TestGetTransaction(t *testing.T) {
	userId := uuid.New()
	validContextUser := models.UserWithRoles{
		User: models.User{
			ID: userId,
		},
	}
	validResponse := &protogen.GetTransactionResponse{
		Transaction: &protogen.Transaction{
			Id:       uuid.New().String(),
			UserId:   userId.String(),
			BrokerId: uuid.New().String(),
		},
	}

	// Define tests
	tests := []struct {
		name           string
		mockSetup      func(ctrl *gomock.Controller)
		expectedStatus int
	}{
		{
			name: "fails to parse param",
			mockSetup: func(ctrl *gomock.Controller) {
				h := mocks.NewMockUtils(ctrl)
				h.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.Nil, false)
				handlers.ReplaceGlobals(h)
				tc := mocks.NewTransactionServiceClient(ctrl)
				tc.EXPECT().GetTransaction(gomock.Any(), gomock.Any()).Times(0)
				clients.ReplaceGlobals(clients.NewClients(nil, nil, tc))
			},
			expectedStatus: http.StatusOK, // should be http.StatusBadRequest, but not with mock
		},
		{
			name: "fails to retrieve the transaction",
			mockSetup: func(ctrl *gomock.Controller) {
				h := mocks.NewMockUtils(ctrl)
				h.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.New(), true)
				h.EXPECT().GetUserFromContext(gomock.Any()).Times(0)
				handlers.ReplaceGlobals(h)
				tc := mocks.NewTransactionServiceClient(ctrl)
				tc.EXPECT().GetTransaction(gomock.Any(), gomock.Any()).Return(nil, status.Error(codes.Unknown, "error"))
				clients.ReplaceGlobals(clients.NewClients(nil, nil, tc))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "fails to find the transaction",
			mockSetup: func(ctrl *gomock.Controller) {
				h := mocks.NewMockUtils(ctrl)
				h.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.New(), true)
				h.EXPECT().GetUserFromContext(gomock.Any()).Times(0)
				handlers.ReplaceGlobals(h)
				tc := mocks.NewTransactionServiceClient(ctrl)
				tc.EXPECT().GetTransaction(gomock.Any(), gomock.Any()).Return(nil, status.Error(codes.NotFound, "error"))
				clients.ReplaceGlobals(clients.NewClients(nil, nil, tc))
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name: "fails to retrieve user from context",
			mockSetup: func(ctrl *gomock.Controller) {
				h := mocks.NewMockUtils(ctrl)
				h.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.New(), true)
				h.EXPECT().GetUserFromContext(gomock.Any()).Return(models.UserWithRoles{}, false)
				handlers.ReplaceGlobals(h)
				tc := mocks.NewTransactionServiceClient(ctrl)
				tc.EXPECT().GetTransaction(gomock.Any(), gomock.Any()).Return(nil, nil)
				clients.ReplaceGlobals(clients.NewClients(nil, nil, tc))
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "transaction user not corresponds to context user",
			mockSetup: func(ctrl *gomock.Controller) {
				h := mocks.NewMockUtils(ctrl)
				h.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.New(), true)
				h.EXPECT().GetUserFromContext(gomock.Any()).Return(models.UserWithRoles{}, true)
				handlers.ReplaceGlobals(h)
				tc := mocks.NewTransactionServiceClient(ctrl)
				tc.EXPECT().GetTransaction(gomock.Any(), gomock.Any()).Return(validResponse, nil)
				clients.ReplaceGlobals(clients.NewClients(nil, nil, tc))
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().Get(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(bb, nil, nil))
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "fails to retrieve broker",
			mockSetup: func(ctrl *gomock.Controller) {
				h := mocks.NewMockUtils(ctrl)
				h.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.New(), true)
				h.EXPECT().GetUserFromContext(gomock.Any()).Return(validContextUser, true)
				handlers.ReplaceGlobals(h)
				tc := mocks.NewTransactionServiceClient(ctrl)
				tc.EXPECT().GetTransaction(gomock.Any(), gomock.Any()).Return(validResponse, nil)
				clients.ReplaceGlobals(clients.NewClients(nil, nil, tc))
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().Get(gomock.Any()).Return(models.Broker{}, false, errors.New("error"))
				repositories.ReplaceGlobals(repositories.NewRepository(bb, nil, nil))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "fails to find broker",
			mockSetup: func(ctrl *gomock.Controller) {
				h := mocks.NewMockUtils(ctrl)
				h.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.New(), true)
				h.EXPECT().GetUserFromContext(gomock.Any()).Return(validContextUser, true)
				handlers.ReplaceGlobals(h)
				tc := mocks.NewTransactionServiceClient(ctrl)
				tc.EXPECT().GetTransaction(gomock.Any(), gomock.Any()).Return(validResponse, nil)
				clients.ReplaceGlobals(clients.NewClients(nil, nil, tc))
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().Get(gomock.Any()).Return(models.Broker{}, false, nil)
				repositories.ReplaceGlobals(repositories.NewRepository(bb, nil, nil))
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name: "succeeded",
			mockSetup: func(ctrl *gomock.Controller) {
				h := mocks.NewMockUtils(ctrl)
				h.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.New(), true)
				h.EXPECT().GetUserFromContext(gomock.Any()).Return(validContextUser, true)
				handlers.ReplaceGlobals(h)
				tc := mocks.NewTransactionServiceClient(ctrl)
				tc.EXPECT().GetTransaction(gomock.Any(), gomock.Any()).Return(validResponse, nil)
				clients.ReplaceGlobals(clients.NewClients(nil, nil, tc))
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().Get(gomock.Any()).Return(models.Broker{}, true, nil)
				repositories.ReplaceGlobals(repositories.NewRepository(bb, nil, nil))
			},
			expectedStatus: http.StatusOK,
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiBasePath := viper.GetString("API_BASE_PATH")
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", apiBasePath+"/transactions/{id}", nil)

			// Apply mocks
			ctrl := gomock.NewController(t)
			tt.mockSetup(ctrl)
			defer ctrl.Finish()

			handlers.GetTransaction(w, r)
			response := w.Result()
			defer response.Body.Close()

			assert.Equal(t, tt.expectedStatus, response.StatusCode)
		})
	}
}

// TestUpdateTransaction tests the UpdateTransaction handler
func TestUpdateTransaction(t *testing.T) {
	// Define request bodies
	validRequest := models.TransactionInput{
		ID:       uuid.New(),
		BrokerID: uuid.New(),
		Date:     time.Now().AddDate(-1, 0, 0), // 1 year in the past
		Type:     models.BUY,
		Asset:    "asset",
		Quantity: 1,
		Price:    1,
		Fee:      1,
	}
	validRequestBody, _ := json.Marshal(validRequest)
	validResponse := &protogen.UpdateTransactionResponse{
		Transaction: &protogen.Transaction{
			Id:       uuid.New().String(),
			UserId:   uuid.New().String(),
			BrokerId: uuid.New().String(),
		},
	}

	// Define tests
	tests := []struct {
		name           string
		body           []byte
		mockSetup      func(ctrl *gomock.Controller)
		expectedStatus int
	}{
		{
			name: "fails to parse param",
			mockSetup: func(ctrl *gomock.Controller) {
				h := mocks.NewMockUtils(ctrl)
				h.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.Nil, false)
				h.EXPECT().GetUserFromContext(gomock.Any()).Times(0)
				handlers.ReplaceGlobals(h)
			},
			expectedStatus: http.StatusOK, // should be http.StatusBadRequest, but not with mock
		},
		{
			name: "fails to retrieve user from context",
			mockSetup: func(ctrl *gomock.Controller) {
				h := mocks.NewMockUtils(ctrl)
				h.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.Nil, true)
				h.EXPECT().GetUserFromContext(gomock.Any()).Return(models.UserWithRoles{}, false)
				handlers.ReplaceGlobals(h)
				bu := mocks.NewBrokerUserRepository(ctrl)
				bu.EXPECT().Exists(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(nil, bu, nil))
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "fails to decode",
			body: []byte("invalid json"),
			mockSetup: func(ctrl *gomock.Controller) {
				h := mocks.NewMockUtils(ctrl)
				h.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.Nil, true)
				h.EXPECT().GetUserFromContext(gomock.Any()).Return(models.UserWithRoles{}, true)
				handlers.ReplaceGlobals(h)
				bu := mocks.NewBrokerUserRepository(ctrl)
				bu.EXPECT().Exists(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(nil, bu, nil))
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "fails to verify user broker existence",
			body: validRequestBody,
			mockSetup: func(ctrl *gomock.Controller) {
				h := mocks.NewMockUtils(ctrl)
				h.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.Nil, true)
				h.EXPECT().GetUserFromContext(gomock.Any()).Return(models.UserWithRoles{}, true)
				handlers.ReplaceGlobals(h)
				bu := mocks.NewBrokerUserRepository(ctrl)
				bu.EXPECT().Exists(gomock.Any()).Return(false, errors.New("error"))
				repositories.ReplaceGlobals(repositories.NewRepository(nil, bu, nil))
				tc := mocks.NewTransactionServiceClient(ctrl)
				tc.EXPECT().UpdateTransaction(gomock.Any(), gomock.Any()).Times(0)
				clients.ReplaceGlobals(clients.NewClients(nil, nil, tc))

			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "broker user does not exist",
			body: validRequestBody,
			mockSetup: func(ctrl *gomock.Controller) {
				h := mocks.NewMockUtils(ctrl)
				h.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.Nil, true)
				h.EXPECT().GetUserFromContext(gomock.Any()).Return(models.UserWithRoles{}, true)
				handlers.ReplaceGlobals(h)
				bu := mocks.NewBrokerUserRepository(ctrl)
				bu.EXPECT().Exists(gomock.Any()).Return(false, nil)
				repositories.ReplaceGlobals(repositories.NewRepository(nil, bu, nil))
				tc := mocks.NewTransactionServiceClient(ctrl)
				tc.EXPECT().UpdateTransaction(gomock.Any(), gomock.Any()).Times(0)
				clients.ReplaceGlobals(clients.NewClients(nil, nil, tc))
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "fails to update the transaction",
			body: validRequestBody,
			mockSetup: func(ctrl *gomock.Controller) {
				h := mocks.NewMockUtils(ctrl)
				h.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.Nil, true)
				h.EXPECT().GetUserFromContext(gomock.Any()).Return(models.UserWithRoles{}, true)
				handlers.ReplaceGlobals(h)
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().Get(gomock.Any()).Times(0)
				bu := mocks.NewBrokerUserRepository(ctrl)
				bu.EXPECT().Exists(gomock.Any()).Return(true, nil)
				repositories.ReplaceGlobals(repositories.NewRepository(bb, bu, nil))
				tc := mocks.NewTransactionServiceClient(ctrl)
				tc.EXPECT().UpdateTransaction(gomock.Any(), gomock.Any()).Return(nil, status.Error(codes.Unknown, "error"))
				clients.ReplaceGlobals(clients.NewClients(nil, nil, tc))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "fails to retrieve broker",
			body: validRequestBody,
			mockSetup: func(ctrl *gomock.Controller) {
				h := mocks.NewMockUtils(ctrl)
				h.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.Nil, true)
				h.EXPECT().GetUserFromContext(gomock.Any()).Return(models.UserWithRoles{}, true)
				handlers.ReplaceGlobals(h)
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().Get(gomock.Any()).Return(models.Broker{}, false, errors.New("error"))
				bu := mocks.NewBrokerUserRepository(ctrl)
				bu.EXPECT().Exists(gomock.Any()).Return(true, nil)
				repositories.ReplaceGlobals(repositories.NewRepository(bb, bu, nil))
				tc := mocks.NewTransactionServiceClient(ctrl)
				tc.EXPECT().UpdateTransaction(gomock.Any(), gomock.Any()).Return(validResponse, nil)
				clients.ReplaceGlobals(clients.NewClients(nil, nil, tc))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "fails to find broker",
			body: validRequestBody,
			mockSetup: func(ctrl *gomock.Controller) {
				h := mocks.NewMockUtils(ctrl)
				h.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.Nil, true)
				h.EXPECT().GetUserFromContext(gomock.Any()).Return(models.UserWithRoles{}, true)
				handlers.ReplaceGlobals(h)
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().Get(gomock.Any()).Return(models.Broker{}, false, nil)
				bu := mocks.NewBrokerUserRepository(ctrl)
				bu.EXPECT().Exists(gomock.Any()).Return(true, nil)
				repositories.ReplaceGlobals(repositories.NewRepository(bb, bu, nil))
				tc := mocks.NewTransactionServiceClient(ctrl)
				tc.EXPECT().UpdateTransaction(gomock.Any(), gomock.Any()).Return(validResponse, nil)
				clients.ReplaceGlobals(clients.NewClients(nil, nil, tc))
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name: "succeeded",
			body: validRequestBody,
			mockSetup: func(ctrl *gomock.Controller) {
				h := mocks.NewMockUtils(ctrl)
				h.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.Nil, true)
				h.EXPECT().GetUserFromContext(gomock.Any()).Return(models.UserWithRoles{}, true)
				handlers.ReplaceGlobals(h)
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().Get(gomock.Any()).Return(models.Broker{}, true, nil)
				bu := mocks.NewBrokerUserRepository(ctrl)
				bu.EXPECT().Exists(gomock.Any()).Return(true, nil)
				repositories.ReplaceGlobals(repositories.NewRepository(bb, bu, nil))
				tc := mocks.NewTransactionServiceClient(ctrl)
				tc.EXPECT().UpdateTransaction(gomock.Any(), gomock.Any()).Return(validResponse, nil)
				clients.ReplaceGlobals(clients.NewClients(nil, nil, tc))
			},
			expectedStatus: http.StatusOK,
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiBasePath := viper.GetString("API_BASE_PATH")
			w := httptest.NewRecorder()
			r := httptest.NewRequest("PUT", apiBasePath+"/transactions/{id}", bytes.NewBuffer(tt.body))

			// Apply mocks
			ctrl := gomock.NewController(t)
			tt.mockSetup(ctrl)
			defer ctrl.Finish()

			handlers.UpdateTransaction(w, r)
			response := w.Result()
			defer response.Body.Close()

			assert.Equal(t, tt.expectedStatus, response.StatusCode)
		})
	}
}

// TestDeleteTransaction tests the DeleteTransaction handler
func TestDeleteTransaction(t *testing.T) {
	// Define tests
	tests := []struct {
		name           string
		mockSetup      func(ctrl *gomock.Controller)
		expectedStatus int
	}{
		{
			name: "fails to parse param",
			mockSetup: func(ctrl *gomock.Controller) {
				h := mocks.NewMockUtils(ctrl)
				h.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.Nil, false)
				h.EXPECT().GetUserFromContext(gomock.Any()).Times(0)
				handlers.ReplaceGlobals(h)
			},
			expectedStatus: http.StatusOK, // should be http.StatusBadRequest, but not with mock
		},
		{
			name: "fails to retrieve user from context",
			mockSetup: func(ctrl *gomock.Controller) {
				h := mocks.NewMockUtils(ctrl)
				h.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.Nil, true)
				h.EXPECT().GetUserFromContext(gomock.Any()).Return(models.UserWithRoles{}, false)
				handlers.ReplaceGlobals(h)
				tc := mocks.NewTransactionServiceClient(ctrl)
				tc.EXPECT().DeleteTransaction(gomock.Any(), gomock.Any()).Times(0)
				clients.ReplaceGlobals(clients.NewClients(nil, nil, tc))
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "fails to delete",
			mockSetup: func(ctrl *gomock.Controller) {
				h := mocks.NewMockUtils(ctrl)
				h.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.Nil, true)
				h.EXPECT().GetUserFromContext(gomock.Any()).Return(models.UserWithRoles{}, true)
				handlers.ReplaceGlobals(h)
				tc := mocks.NewTransactionServiceClient(ctrl)
				tc.EXPECT().DeleteTransaction(gomock.Any(), gomock.Any()).Return(
					&protogen.DeleteTransactionResponse{},
					status.Error(codes.Unknown, "error"))
				clients.ReplaceGlobals(clients.NewClients(nil, nil, tc))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "succeeded",
			mockSetup: func(ctrl *gomock.Controller) {
				h := mocks.NewMockUtils(ctrl)
				h.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.Nil, true)
				h.EXPECT().GetUserFromContext(gomock.Any()).Return(models.UserWithRoles{}, true)
				handlers.ReplaceGlobals(h)
				tc := mocks.NewTransactionServiceClient(ctrl)
				tc.EXPECT().DeleteTransaction(gomock.Any(), gomock.Any()).Return(
					&protogen.DeleteTransactionResponse{}, nil)
				clients.ReplaceGlobals(clients.NewClients(nil, nil, tc))
			},
			expectedStatus: http.StatusOK,
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiBasePath := viper.GetString("API_BASE_PATH")
			w := httptest.NewRecorder()
			r := httptest.NewRequest("DELETE", apiBasePath+"/transactions/{id}", nil)

			// Apply mocks
			ctrl := gomock.NewController(t)
			tt.mockSetup(ctrl)
			defer ctrl.Finish()

			handlers.DeleteTransaction(w, r)
			response := w.Result()
			defer response.Body.Close()

			assert.Equal(t, tt.expectedStatus, response.StatusCode)
		})
	}
}

// TestListTransactions tests the ListTransactions handler
func TestListTransactions(t *testing.T) {
	// Define tests
	tests := []struct {
		name           string
		mockSetup      func(ctrl *gomock.Controller)
		expectedStatus int
	}{
		{
			name: "fails to retrieve user from context",
			mockSetup: func(ctrl *gomock.Controller) {
				h := mocks.NewMockUtils(ctrl)
				h.EXPECT().GetUserFromContext(gomock.Any()).Return(models.UserWithRoles{}, false)
				handlers.ReplaceGlobals(h)
				tc := mocks.NewTransactionServiceClient(ctrl)
				tc.EXPECT().ListTransactions(gomock.Any(), gomock.Any()).Times(0)
				clients.ReplaceGlobals(clients.NewClients(nil, nil, tc))
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "fails to retrieve all user transactions",
			mockSetup: func(ctrl *gomock.Controller) {
				h := mocks.NewMockUtils(ctrl)
				h.EXPECT().GetUserFromContext(gomock.Any()).Return(models.UserWithRoles{}, true)
				handlers.ReplaceGlobals(h)
				tc := mocks.NewTransactionServiceClient(ctrl)
				tc.EXPECT().ListTransactions(gomock.Any(), gomock.Any()).Return(nil, status.Error(codes.Unknown, "error"))
				clients.ReplaceGlobals(clients.NewClients(nil, nil, tc))
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().GetAll().Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(bb, nil, nil))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "fails to retrieve all brokers",
			mockSetup: func(ctrl *gomock.Controller) {
				h := mocks.NewMockUtils(ctrl)
				h.EXPECT().GetUserFromContext(gomock.Any()).Return(models.UserWithRoles{}, true)
				handlers.ReplaceGlobals(h)
				tc := mocks.NewTransactionServiceClient(ctrl)
				tc.EXPECT().ListTransactions(gomock.Any(), gomock.Any()).Return(&protogen.ListTransactionsResponse{}, nil)
				clients.ReplaceGlobals(clients.NewClients(nil, nil, tc))
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().GetAll().Return(nil, errors.New("error"))
				repositories.ReplaceGlobals(repositories.NewRepository(bb, nil, nil))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "succeeded",
			mockSetup: func(ctrl *gomock.Controller) {
				h := mocks.NewMockUtils(ctrl)
				h.EXPECT().GetUserFromContext(gomock.Any()).Return(models.UserWithRoles{}, true)
				handlers.ReplaceGlobals(h)
				tc := mocks.NewTransactionServiceClient(ctrl)
				tc.EXPECT().ListTransactions(gomock.Any(), gomock.Any()).Return(&protogen.ListTransactionsResponse{}, nil)
				clients.ReplaceGlobals(clients.NewClients(nil, nil, tc))
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().GetAll().Return([]models.Broker{}, nil)
				repositories.ReplaceGlobals(repositories.NewRepository(bb, nil, nil))
			},
			expectedStatus: http.StatusOK,
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiBasePath := viper.GetString("API_BASE_PATH")
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", apiBasePath+"/transactions", nil)

			// Apply mocks
			ctrl := gomock.NewController(t)
			tt.mockSetup(ctrl)
			defer ctrl.Finish()

			handlers.ListTransactions(w, r)
			response := w.Result()
			defer response.Body.Close()

			assert.Equal(t, tt.expectedStatus, response.StatusCode)
		})
	}
}
