package handlers_test

import (
	"bytes"
	"encoding/json"
	"errors"
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
	validResponse := &transactionpb.CreateTransactionResponse{
		Transaction: &transactionpb.Transaction{
			Id:       uuid.New().String(),
			UserId:   uuid.New().String(),
			BrokerId: uuid.New().String(),
		},
	}
	validResponseBroker := &brokerpb.GetBrokerResponse{
		Broker: &brokerpb.Broker{
			Id:       uuid.New().String(),
			Name:     "broker",
			Disabled: false,
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
				m := mocks.NewMockApiUtils(ctrl)
				m.EXPECT().GetUserIDFromContext(gomock.Any()).Return(uuid.Nil.String(), false)
				handlers.ReplaceGlobals(m)
				bc := mocks.NewMockBrokerServiceClient(ctrl)
				bc.EXPECT().GetBrokerUser(gomock.Any(), gomock.Any()).Times(0)
				clients.ReplaceGlobals(clients.NewClients(
					clients.WithBrokerClient(bc),
				))
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "fails to decode",
			body: []byte("invalid json"),
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockApiUtils(ctrl)
				m.EXPECT().GetUserIDFromContext(gomock.Any()).Return(uuid.New().String(), true)
				handlers.ReplaceGlobals(m)
				bc := mocks.NewMockBrokerServiceClient(ctrl)
				bc.EXPECT().GetBrokerUser(gomock.Any(), gomock.Any()).Times(0)
				clients.ReplaceGlobals(clients.NewClients(
					clients.WithBrokerClient(bc),
				))
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "user broker existence check",
			body: validRequestBody,
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockApiUtils(ctrl)
				m.EXPECT().GetUserIDFromContext(gomock.Any()).Return(uuid.New().String(), true)
				handlers.ReplaceGlobals(m)
				bc := mocks.NewMockBrokerServiceClient(ctrl)
				bc.EXPECT().GetBrokerUser(gomock.Any(), gomock.Any()).Return(nil, status.Error(codes.Unknown, "error"))
				tc := mocks.NewMockTransactionServiceClient(ctrl)
				tc.EXPECT().CreateTransaction(gomock.Any(), gomock.Any()).Times(0)
				clients.ReplaceGlobals(clients.NewClients(
					clients.WithBrokerClient(bc),
					clients.WithTransactionClient(tc),
				))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "fails at transactions creation",
			body: validRequestBody,
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockApiUtils(ctrl)
				m.EXPECT().GetUserIDFromContext(gomock.Any()).Return(uuid.New().String(), true)
				handlers.ReplaceGlobals(m)
				bc := mocks.NewMockBrokerServiceClient(ctrl)
				bc.EXPECT().GetBrokerUser(gomock.Any(), gomock.Any()).Return(nil, nil)
				bc.EXPECT().GetBroker(gomock.Any(), gomock.Any()).Times(0)
				tc := mocks.NewMockTransactionServiceClient(ctrl)
				tc.EXPECT().CreateTransaction(gomock.Any(), gomock.Any()).Return(nil, status.Error(codes.Unknown, "error"))
				clients.ReplaceGlobals(clients.NewClients(
					clients.WithBrokerClient(bc),
					clients.WithTransactionClient(tc),
				))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "fails to retrieve broker",
			body: validRequestBody,
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockApiUtils(ctrl)
				m.EXPECT().GetUserIDFromContext(gomock.Any()).Return(uuid.New().String(), true)
				handlers.ReplaceGlobals(m)
				bc := mocks.NewMockBrokerServiceClient(ctrl)
				bc.EXPECT().GetBrokerUser(gomock.Any(), gomock.Any()).Return(nil, nil)
				bc.EXPECT().GetBroker(gomock.Any(), gomock.Any()).Return(&brokerpb.GetBrokerResponse{}, status.Error(codes.Unknown, "error"))
				tc := mocks.NewMockTransactionServiceClient(ctrl)
				tc.EXPECT().CreateTransaction(gomock.Any(), gomock.Any()).Return(validResponse, nil)
				clients.ReplaceGlobals(clients.NewClients(
					clients.WithBrokerClient(bc),
					clients.WithTransactionClient(tc),
				))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "succeeded",
			body: validRequestBody,
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockApiUtils(ctrl)
				m.EXPECT().GetUserIDFromContext(gomock.Any()).Return(uuid.New().String(), true)
				handlers.ReplaceGlobals(m)
				bc := mocks.NewMockBrokerServiceClient(ctrl)
				bc.EXPECT().GetBrokerUser(gomock.Any(), gomock.Any()).Return(nil, nil)
				bc.EXPECT().GetBroker(gomock.Any(), gomock.Any()).Return(validResponseBroker, nil)
				tc := mocks.NewMockTransactionServiceClient(ctrl)
				tc.EXPECT().CreateTransaction(gomock.Any(), gomock.Any()).Return(validResponse, nil)
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
			r := httptest.NewRequest("POST", apiBasePath+"/transaction", bytes.NewBuffer(tt.body))

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
	validResponse := &transactionpb.GetTransactionResponse{
		Transaction: &transactionpb.Transaction{
			Id:       uuid.New().String(),
			UserId:   userId.String(),
			BrokerId: uuid.New().String(),
		},
	}
	validResponseBroker := &brokerpb.GetBrokerResponse{
		Broker: &brokerpb.Broker{
			Id:       uuid.New().String(),
			Name:     "broker",
			Disabled: false,
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
				m := mocks.NewMockApiUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.Nil, false)
				handlers.ReplaceGlobals(m)
				tc := mocks.NewMockTransactionServiceClient(ctrl)
				tc.EXPECT().GetTransaction(gomock.Any(), gomock.Any()).Times(0)
				clients.ReplaceGlobals(clients.NewClients(
					clients.WithTransactionClient(tc),
				))
			},
			expectedStatus: http.StatusOK, // should be http.StatusBadRequest, but not with mock
		},
		{
			name: "fails to retrieve the transaction",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockApiUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.New(), true)
				m.EXPECT().GetUserIDFromContext(gomock.Any()).Times(0)
				handlers.ReplaceGlobals(m)
				tc := mocks.NewMockTransactionServiceClient(ctrl)
				tc.EXPECT().GetTransaction(gomock.Any(), gomock.Any()).Return(nil, status.Error(codes.Unknown, "error"))
				clients.ReplaceGlobals(clients.NewClients(
					clients.WithTransactionClient(tc),
				))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "fails to retrieve user from context",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockApiUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.New(), true)
				m.EXPECT().GetUserIDFromContext(gomock.Any()).Return(uuid.Nil.String(), false)
				handlers.ReplaceGlobals(m)
				tc := mocks.NewMockTransactionServiceClient(ctrl)
				tc.EXPECT().GetTransaction(gomock.Any(), gomock.Any()).Return(nil, nil)
				clients.ReplaceGlobals(clients.NewClients(
					clients.WithTransactionClient(tc),
				))
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "transaction user not corresponds to context user",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockApiUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.New(), true)
				m.EXPECT().GetUserIDFromContext(gomock.Any()).Return(uuid.New().String(), true)
				handlers.ReplaceGlobals(m)
				tc := mocks.NewMockTransactionServiceClient(ctrl)
				tc.EXPECT().GetTransaction(gomock.Any(), gomock.Any()).Return(validResponse, nil)
				bc := mocks.NewMockBrokerServiceClient(ctrl)
				bc.EXPECT().GetBroker(gomock.Any(), gomock.Any()).Times(0)
				clients.ReplaceGlobals(clients.NewClients(
					clients.WithBrokerClient(bc),
					clients.WithTransactionClient(tc),
				))
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "fails to retrieve broker",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockApiUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.New(), true)
				m.EXPECT().GetUserIDFromContext(gomock.Any()).Return(userId.String(), true)
				handlers.ReplaceGlobals(m)
				tc := mocks.NewMockTransactionServiceClient(ctrl)
				tc.EXPECT().GetTransaction(gomock.Any(), gomock.Any()).Return(validResponse, nil)
				bc := mocks.NewMockBrokerServiceClient(ctrl)
				bc.EXPECT().GetBroker(gomock.Any(), gomock.Any()).Return(nil, errors.New("error"))
				clients.ReplaceGlobals(clients.NewClients(
					clients.WithBrokerClient(bc),
					clients.WithTransactionClient(tc),
				))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "succeeded",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockApiUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.New(), true)
				m.EXPECT().GetUserIDFromContext(gomock.Any()).Return(userId.String(), true)
				handlers.ReplaceGlobals(m)
				tc := mocks.NewMockTransactionServiceClient(ctrl)
				tc.EXPECT().GetTransaction(gomock.Any(), gomock.Any()).Return(validResponse, nil)
				bc := mocks.NewMockBrokerServiceClient(ctrl)
				bc.EXPECT().GetBroker(gomock.Any(), gomock.Any()).Return(validResponseBroker, nil)
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
			r := httptest.NewRequest("GET", apiBasePath+"/transaction/"+uuid.New().String(), nil)

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
	validResponse := &transactionpb.UpdateTransactionResponse{
		Transaction: &transactionpb.Transaction{
			Id:       uuid.New().String(),
			UserId:   uuid.New().String(),
			BrokerId: uuid.New().String(),
		},
	}
	validResponseBroker := &brokerpb.GetBrokerResponse{
		Broker: &brokerpb.Broker{
			Id:       uuid.New().String(),
			Name:     "broker",
			Disabled: false,
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
				m := mocks.NewMockApiUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.Nil, false)
				m.EXPECT().GetUserIDFromContext(gomock.Any()).Times(0)
				handlers.ReplaceGlobals(m)
			},
			expectedStatus: http.StatusOK, // should be http.StatusBadRequest, but not with mock
		},
		{
			name: "fails to retrieve user from context",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockApiUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.New(), true)
				m.EXPECT().GetUserIDFromContext(gomock.Any()).Return(uuid.Nil.String(), false)
				handlers.ReplaceGlobals(m)
				bc := mocks.NewMockBrokerServiceClient(ctrl)
				bc.EXPECT().GetBrokerUser(gomock.Any(), gomock.Any()).Times(0)
				clients.ReplaceGlobals(clients.NewClients(
					clients.WithBrokerClient(bc),
				))
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "fails to decode",
			body: []byte("invalid json"),
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockApiUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.New(), true)
				m.EXPECT().GetUserIDFromContext(gomock.Any()).Return(uuid.New().String(), true)
				handlers.ReplaceGlobals(m)
				bc := mocks.NewMockBrokerServiceClient(ctrl)
				bc.EXPECT().GetBrokerUser(gomock.Any(), gomock.Any()).Times(0)
				clients.ReplaceGlobals(clients.NewClients(
					clients.WithBrokerClient(bc),
				))
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "fails to verify user broker existence",
			body: validRequestBody,
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockApiUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.New(), true)
				m.EXPECT().GetUserIDFromContext(gomock.Any()).Return(uuid.New().String(), true)
				handlers.ReplaceGlobals(m)
				bc := mocks.NewMockBrokerServiceClient(ctrl)
				bc.EXPECT().GetBrokerUser(gomock.Any(), gomock.Any()).Return(nil, status.Error(codes.Unknown, "error"))
				tc := mocks.NewMockTransactionServiceClient(ctrl)
				tc.EXPECT().UpdateTransaction(gomock.Any(), gomock.Any()).Times(0)
				clients.ReplaceGlobals(clients.NewClients(
					clients.WithBrokerClient(bc),
					clients.WithTransactionClient(tc),
				))

			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "fails to update the transaction",
			body: validRequestBody,
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockApiUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.New(), true)
				m.EXPECT().GetUserIDFromContext(gomock.Any()).Return(uuid.New().String(), true)
				handlers.ReplaceGlobals(m)
				bc := mocks.NewMockBrokerServiceClient(ctrl)
				bc.EXPECT().GetBrokerUser(gomock.Any(), gomock.Any()).Return(nil, nil)
				bc.EXPECT().GetBroker(gomock.Any(), gomock.Any()).Times(0)
				tc := mocks.NewMockTransactionServiceClient(ctrl)
				tc.EXPECT().UpdateTransaction(gomock.Any(), gomock.Any()).Return(nil, status.Error(codes.Unknown, "error"))
				clients.ReplaceGlobals(clients.NewClients(
					clients.WithBrokerClient(bc),
					clients.WithTransactionClient(tc),
				))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "fails to retrieve broker",
			body: validRequestBody,
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockApiUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.New(), true)
				m.EXPECT().GetUserIDFromContext(gomock.Any()).Return(uuid.New().String(), true)
				handlers.ReplaceGlobals(m)
				bc := mocks.NewMockBrokerServiceClient(ctrl)
				bc.EXPECT().GetBrokerUser(gomock.Any(), gomock.Any()).Return(nil, nil)
				bc.EXPECT().GetBroker(gomock.Any(), gomock.Any()).Return(nil, status.Error(codes.Unknown, "error"))
				tc := mocks.NewMockTransactionServiceClient(ctrl)
				tc.EXPECT().UpdateTransaction(gomock.Any(), gomock.Any()).Return(validResponse, nil)
				clients.ReplaceGlobals(clients.NewClients(
					clients.WithBrokerClient(bc),
					clients.WithTransactionClient(tc),
				))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "succeeded",
			body: validRequestBody,
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockApiUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.New(), true)
				m.EXPECT().GetUserIDFromContext(gomock.Any()).Return(uuid.New().String(), true)
				handlers.ReplaceGlobals(m)
				bc := mocks.NewMockBrokerServiceClient(ctrl)
				bc.EXPECT().GetBrokerUser(gomock.Any(), gomock.Any()).Return(nil, nil)
				bc.EXPECT().GetBroker(gomock.Any(), gomock.Any()).Return(validResponseBroker, nil)
				tc := mocks.NewMockTransactionServiceClient(ctrl)
				tc.EXPECT().UpdateTransaction(gomock.Any(), gomock.Any()).Return(validResponse, nil)
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
			r := httptest.NewRequest("PUT", apiBasePath+"/transaction/"+uuid.New().String(), bytes.NewBuffer(tt.body))

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
				m := mocks.NewMockApiUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.Nil, false)
				m.EXPECT().GetUserIDFromContext(gomock.Any()).Times(0)
				handlers.ReplaceGlobals(m)
			},
			expectedStatus: http.StatusOK, // should be http.StatusBadRequest, but not with mock
		},
		{
			name: "fails to retrieve user from context",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockApiUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.New(), true)
				m.EXPECT().GetUserIDFromContext(gomock.Any()).Return(uuid.Nil.String(), false)
				handlers.ReplaceGlobals(m)
				tc := mocks.NewMockTransactionServiceClient(ctrl)
				tc.EXPECT().DeleteTransaction(gomock.Any(), gomock.Any()).Times(0)
				clients.ReplaceGlobals(clients.NewClients(
					clients.WithTransactionClient(tc),
				))
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "fails to delete",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockApiUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.New(), true)
				m.EXPECT().GetUserIDFromContext(gomock.Any()).Return(uuid.New().String(), true)
				handlers.ReplaceGlobals(m)
				tc := mocks.NewMockTransactionServiceClient(ctrl)
				tc.EXPECT().DeleteTransaction(gomock.Any(), gomock.Any()).Return(
					&transactionpb.DeleteTransactionResponse{},
					status.Error(codes.Unknown, "error"))
				clients.ReplaceGlobals(clients.NewClients(
					clients.WithTransactionClient(tc),
				))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "succeeded",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockApiUtils(ctrl)
				m.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.New(), true)
				m.EXPECT().GetUserIDFromContext(gomock.Any()).Return(uuid.New().String(), true)
				handlers.ReplaceGlobals(m)
				tc := mocks.NewMockTransactionServiceClient(ctrl)
				tc.EXPECT().DeleteTransaction(gomock.Any(), gomock.Any()).Return(
					&transactionpb.DeleteTransactionResponse{}, nil)
				clients.ReplaceGlobals(clients.NewClients(
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
			r := httptest.NewRequest("DELETE", apiBasePath+"/transaction/"+uuid.New().String(), nil)

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
				m := mocks.NewMockApiUtils(ctrl)
				m.EXPECT().GetUserIDFromContext(gomock.Any()).Return(uuid.Nil.String(), false)
				handlers.ReplaceGlobals(m)
				tc := mocks.NewMockTransactionServiceClient(ctrl)
				tc.EXPECT().ListTransactions(gomock.Any(), gomock.Any()).Times(0)
				clients.ReplaceGlobals(clients.NewClients(
					clients.WithTransactionClient(tc),
				))
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "fails to retrieve all user transactions",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockApiUtils(ctrl)
				m.EXPECT().GetUserIDFromContext(gomock.Any()).Return(uuid.New().String(), true)
				handlers.ReplaceGlobals(m)
				tc := mocks.NewMockTransactionServiceClient(ctrl)
				tc.EXPECT().ListTransactions(gomock.Any(), gomock.Any()).Return(nil, status.Error(codes.Unknown, "error"))
				bc := mocks.NewMockBrokerServiceClient(ctrl)
				bc.EXPECT().ListBrokers(gomock.Any(), gomock.Any()).Times(0)
				clients.ReplaceGlobals(clients.NewClients(
					clients.WithBrokerClient(bc),
					clients.WithTransactionClient(tc),
				))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "fails to retrieve all brokers",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockApiUtils(ctrl)
				m.EXPECT().GetUserIDFromContext(gomock.Any()).Return(uuid.New().String(), true)
				handlers.ReplaceGlobals(m)
				tc := mocks.NewMockTransactionServiceClient(ctrl)
				tc.EXPECT().ListTransactions(gomock.Any(), gomock.Any()).Return(&transactionpb.ListTransactionsResponse{}, nil)
				bc := mocks.NewMockBrokerServiceClient(ctrl)
				bc.EXPECT().ListBrokers(gomock.Any(), gomock.Any()).Return(nil, status.Error(codes.Unknown, "error"))
				clients.ReplaceGlobals(clients.NewClients(
					clients.WithBrokerClient(bc),
					clients.WithTransactionClient(tc),
				))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "succeeded",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockApiUtils(ctrl)
				m.EXPECT().GetUserIDFromContext(gomock.Any()).Return(uuid.New().String(), true)
				handlers.ReplaceGlobals(m)
				tc := mocks.NewMockTransactionServiceClient(ctrl)
				tc.EXPECT().ListTransactions(gomock.Any(), gomock.Any()).Return(&transactionpb.ListTransactionsResponse{}, nil)
				bc := mocks.NewMockBrokerServiceClient(ctrl)
				bc.EXPECT().ListBrokers(gomock.Any(), gomock.Any()).Return(&brokerpb.ListBrokersResponse{
					Brokers: []*brokerpb.Broker{},
				}, nil)
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
			r := httptest.NewRequest("GET", apiBasePath+"/transaction", nil)

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
