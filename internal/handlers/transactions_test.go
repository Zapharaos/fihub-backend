package handlers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/Zapharaos/fihub-backend/internal/auth/users"
	"github.com/Zapharaos/fihub-backend/internal/brokers"
	"github.com/Zapharaos/fihub-backend/internal/handlers"
	"github.com/Zapharaos/fihub-backend/internal/transactions"
	"github.com/Zapharaos/fihub-backend/test/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// TestCreateTransaction tests the CreateTransaction handler
func TestCreateTransaction(t *testing.T) {
	// Define request bodies
	invalidRequest := transactions.TransactionInput{}
	invalidRequestBody, _ := json.Marshal(invalidRequest)
	validRequest := transactions.TransactionInput{
		ID:       uuid.New(),
		BrokerID: uuid.New(),
		Date:     time.Now().AddDate(-1, 0, 0), // 1 year in the past
		Type:     transactions.BUY,
		Asset:    "asset",
		Quantity: 1,
		Price:    1,
		Fee:      1,
	}
	validRequestBody, _ := json.Marshal(validRequest)

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
				h.EXPECT().GetUserFromContext(gomock.Any()).Return(users.UserWithRoles{}, false)
				handlers.ReplaceGlobals(h)
				bu := mocks.NewBrokerUserRepository(ctrl)
				bu.EXPECT().Exists(gomock.Any()).Times(0)
				brokers.ReplaceGlobals(brokers.NewRepository(nil, bu, nil))
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "fails to decode",
			body: []byte("invalid json"),
			mockSetup: func(ctrl *gomock.Controller) {
				h := mocks.NewMockUtils(ctrl)
				h.EXPECT().GetUserFromContext(gomock.Any()).Return(users.UserWithRoles{}, true)
				handlers.ReplaceGlobals(h)
				bu := mocks.NewBrokerUserRepository(ctrl)
				bu.EXPECT().Exists(gomock.Any()).Times(0)
				brokers.ReplaceGlobals(brokers.NewRepository(nil, bu, nil))
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "fails at bad transaction input",
			body: invalidRequestBody,
			mockSetup: func(ctrl *gomock.Controller) {
				h := mocks.NewMockUtils(ctrl)
				h.EXPECT().GetUserFromContext(gomock.Any()).Return(users.UserWithRoles{}, true)
				handlers.ReplaceGlobals(h)
				bu := mocks.NewBrokerUserRepository(ctrl)
				bu.EXPECT().Exists(gomock.Any()).Times(0)
				brokers.ReplaceGlobals(brokers.NewRepository(nil, bu, nil))
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "fails at user broker existence check",
			body: validRequestBody,
			mockSetup: func(ctrl *gomock.Controller) {
				h := mocks.NewMockUtils(ctrl)
				h.EXPECT().GetUserFromContext(gomock.Any()).Return(users.UserWithRoles{}, true)
				handlers.ReplaceGlobals(h)
				bu := mocks.NewBrokerUserRepository(ctrl)
				bu.EXPECT().Exists(gomock.Any()).Return(false, errors.New("error"))
				brokers.ReplaceGlobals(brokers.NewRepository(nil, bu, nil))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "user broker does not exist",
			body: validRequestBody,
			mockSetup: func(ctrl *gomock.Controller) {
				h := mocks.NewMockUtils(ctrl)
				h.EXPECT().GetUserFromContext(gomock.Any()).Return(users.UserWithRoles{}, true)
				handlers.ReplaceGlobals(h)
				bu := mocks.NewBrokerUserRepository(ctrl)
				bu.EXPECT().Exists(gomock.Any()).Return(false, nil)
				brokers.ReplaceGlobals(brokers.NewRepository(nil, bu, nil))
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "fails at transactions creation",
			body: validRequestBody,
			mockSetup: func(ctrl *gomock.Controller) {
				h := mocks.NewMockUtils(ctrl)
				h.EXPECT().GetUserFromContext(gomock.Any()).Return(users.UserWithRoles{}, true)
				handlers.ReplaceGlobals(h)
				bu := mocks.NewBrokerUserRepository(ctrl)
				bu.EXPECT().Exists(gomock.Any()).Return(true, nil)
				brokers.ReplaceGlobals(brokers.NewRepository(nil, bu, nil))
				tr := mocks.NewTransactionsRepository(ctrl)
				tr.EXPECT().Create(gomock.Any()).Return(uuid.Nil, errors.New("error"))
				tr.EXPECT().Get(gomock.Any()).Times(0)
				transactions.ReplaceGlobals(tr)
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "fails to retrieve the transaction",
			body: validRequestBody,
			mockSetup: func(ctrl *gomock.Controller) {
				h := mocks.NewMockUtils(ctrl)
				h.EXPECT().GetUserFromContext(gomock.Any()).Return(users.UserWithRoles{}, true)
				handlers.ReplaceGlobals(h)
				bu := mocks.NewBrokerUserRepository(ctrl)
				bu.EXPECT().Exists(gomock.Any()).Return(true, nil)
				brokers.ReplaceGlobals(brokers.NewRepository(nil, bu, nil))
				tr := mocks.NewTransactionsRepository(ctrl)
				tr.EXPECT().Create(gomock.Any()).Return(uuid.New(), nil)
				tr.EXPECT().Get(gomock.Any()).Return(transactions.Transaction{}, false, errors.New("error"))
				transactions.ReplaceGlobals(tr)
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "could not find the transaction",
			body: validRequestBody,
			mockSetup: func(ctrl *gomock.Controller) {
				h := mocks.NewMockUtils(ctrl)
				h.EXPECT().GetUserFromContext(gomock.Any()).Return(users.UserWithRoles{}, true)
				handlers.ReplaceGlobals(h)
				bu := mocks.NewBrokerUserRepository(ctrl)
				bu.EXPECT().Exists(gomock.Any()).Return(true, nil)
				brokers.ReplaceGlobals(brokers.NewRepository(nil, bu, nil))
				tr := mocks.NewTransactionsRepository(ctrl)
				tr.EXPECT().Create(gomock.Any()).Return(uuid.New(), nil)
				tr.EXPECT().Get(gomock.Any()).Return(transactions.Transaction{}, false, nil)
				transactions.ReplaceGlobals(tr)
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "succeeded",
			body: validRequestBody,
			mockSetup: func(ctrl *gomock.Controller) {
				h := mocks.NewMockUtils(ctrl)
				h.EXPECT().GetUserFromContext(gomock.Any()).Return(users.UserWithRoles{}, true)
				handlers.ReplaceGlobals(h)
				bu := mocks.NewBrokerUserRepository(ctrl)
				bu.EXPECT().Exists(gomock.Any()).Return(true, nil)
				brokers.ReplaceGlobals(brokers.NewRepository(nil, bu, nil))
				tr := mocks.NewTransactionsRepository(ctrl)
				tr.EXPECT().Create(gomock.Any()).Return(uuid.New(), nil)
				tr.EXPECT().Get(gomock.Any()).Return(transactions.Transaction{}, true, nil)
				transactions.ReplaceGlobals(tr)
			},
			expectedStatus: http.StatusOK,
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", "/api/v1/transactions", bytes.NewBuffer(tt.body))

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
				tr := mocks.NewTransactionsRepository(ctrl)
				tr.EXPECT().Get(gomock.Any()).Times(0)
				transactions.ReplaceGlobals(tr)
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
				tr := mocks.NewTransactionsRepository(ctrl)
				tr.EXPECT().Get(gomock.Any()).Return(transactions.Transaction{}, false, errors.New("error"))
				transactions.ReplaceGlobals(tr)
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
				tr := mocks.NewTransactionsRepository(ctrl)
				tr.EXPECT().Get(gomock.Any()).Return(transactions.Transaction{}, false, nil)
				transactions.ReplaceGlobals(tr)
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name: "fails to retrieve user from context",
			mockSetup: func(ctrl *gomock.Controller) {
				h := mocks.NewMockUtils(ctrl)
				h.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.New(), true)
				h.EXPECT().GetUserFromContext(gomock.Any()).Return(users.UserWithRoles{}, false)
				handlers.ReplaceGlobals(h)
				tr := mocks.NewTransactionsRepository(ctrl)
				tr.EXPECT().Get(gomock.Any()).Return(transactions.Transaction{}, true, nil)
				transactions.ReplaceGlobals(tr)
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "transaction user not corresponds to context user",
			mockSetup: func(ctrl *gomock.Controller) {
				h := mocks.NewMockUtils(ctrl)
				h.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.New(), true)
				h.EXPECT().GetUserFromContext(gomock.Any()).Return(users.UserWithRoles{}, true)
				handlers.ReplaceGlobals(h)
				tr := mocks.NewTransactionsRepository(ctrl)
				tr.EXPECT().Get(gomock.Any()).Return(transactions.Transaction{UserID: uuid.New()}, true, nil)
				transactions.ReplaceGlobals(tr)
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "succeeded",
			mockSetup: func(ctrl *gomock.Controller) {
				h := mocks.NewMockUtils(ctrl)
				h.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.New(), true)
				h.EXPECT().GetUserFromContext(gomock.Any()).Return(users.UserWithRoles{}, true)
				handlers.ReplaceGlobals(h)
				tr := mocks.NewTransactionsRepository(ctrl)
				tr.EXPECT().Get(gomock.Any()).Return(transactions.Transaction{}, true, nil)
				transactions.ReplaceGlobals(tr)
			},
			expectedStatus: http.StatusOK,
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/api/v1/transactions/{id}", nil)

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
	invalidRequest := transactions.TransactionInput{}
	invalidRequestBody, _ := json.Marshal(invalidRequest)
	validRequest := transactions.TransactionInput{
		ID:       uuid.New(),
		BrokerID: uuid.New(),
		Date:     time.Now().AddDate(-1, 0, 0), // 1 year in the past
		Type:     transactions.BUY,
		Asset:    "asset",
		Quantity: 1,
		Price:    1,
		Fee:      1,
	}
	validRequestBody, _ := json.Marshal(validRequest)

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
				h.EXPECT().GetUserFromContext(gomock.Any()).Return(users.UserWithRoles{}, false)
				handlers.ReplaceGlobals(h)
				tr := mocks.NewTransactionsRepository(ctrl)
				tr.EXPECT().Get(gomock.Any()).Times(0)
				transactions.ReplaceGlobals(tr)
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "fails to decode",
			body: []byte("invalid json"),
			mockSetup: func(ctrl *gomock.Controller) {
				h := mocks.NewMockUtils(ctrl)
				h.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.Nil, true)
				h.EXPECT().GetUserFromContext(gomock.Any()).Return(users.UserWithRoles{}, true)
				handlers.ReplaceGlobals(h)
				tr := mocks.NewTransactionsRepository(ctrl)
				tr.EXPECT().Get(gomock.Any()).Times(0)
				transactions.ReplaceGlobals(tr)
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "fails at bad transaction input",
			body: invalidRequestBody,
			mockSetup: func(ctrl *gomock.Controller) {
				h := mocks.NewMockUtils(ctrl)
				h.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.Nil, true)
				h.EXPECT().GetUserFromContext(gomock.Any()).Return(users.UserWithRoles{}, true)
				handlers.ReplaceGlobals(h)
				tr := mocks.NewTransactionsRepository(ctrl)
				tr.EXPECT().Get(gomock.Any()).Times(0)
				transactions.ReplaceGlobals(tr)
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "fails to retrieve the transaction",
			body: validRequestBody,
			mockSetup: func(ctrl *gomock.Controller) {
				h := mocks.NewMockUtils(ctrl)
				h.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.Nil, true)
				h.EXPECT().GetUserFromContext(gomock.Any()).Return(users.UserWithRoles{}, true)
				handlers.ReplaceGlobals(h)
				tr := mocks.NewTransactionsRepository(ctrl)
				tr.EXPECT().Get(gomock.Any()).Return(transactions.Transaction{}, false, errors.New("error"))
				transactions.ReplaceGlobals(tr)
				bu := mocks.NewBrokerUserRepository(ctrl)
				bu.EXPECT().Exists(gomock.Any()).Times(0)
				brokers.ReplaceGlobals(brokers.NewRepository(nil, bu, nil))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "fails to find the transaction",
			body: validRequestBody,
			mockSetup: func(ctrl *gomock.Controller) {
				h := mocks.NewMockUtils(ctrl)
				h.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.Nil, true)
				h.EXPECT().GetUserFromContext(gomock.Any()).Return(users.UserWithRoles{}, true)
				handlers.ReplaceGlobals(h)
				tr := mocks.NewTransactionsRepository(ctrl)
				tr.EXPECT().Get(gomock.Any()).Return(transactions.Transaction{}, false, nil)
				transactions.ReplaceGlobals(tr)
				bu := mocks.NewBrokerUserRepository(ctrl)
				bu.EXPECT().Exists(gomock.Any()).Times(0)
				brokers.ReplaceGlobals(brokers.NewRepository(nil, bu, nil))
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name: "transaction does not belong to user",
			body: validRequestBody,
			mockSetup: func(ctrl *gomock.Controller) {
				h := mocks.NewMockUtils(ctrl)
				h.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.Nil, true)
				h.EXPECT().GetUserFromContext(gomock.Any()).Return(users.UserWithRoles{}, true)
				handlers.ReplaceGlobals(h)
				tr := mocks.NewTransactionsRepository(ctrl)
				tr.EXPECT().Get(gomock.Any()).Return(transactions.Transaction{UserID: uuid.New()}, true, nil)
				transactions.ReplaceGlobals(tr)
				bu := mocks.NewBrokerUserRepository(ctrl)
				bu.EXPECT().Exists(gomock.Any()).Times(0)
				brokers.ReplaceGlobals(brokers.NewRepository(nil, bu, nil))
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "fails to verify user broker existence",
			body: validRequestBody,
			mockSetup: func(ctrl *gomock.Controller) {
				h := mocks.NewMockUtils(ctrl)
				h.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.Nil, true)
				h.EXPECT().GetUserFromContext(gomock.Any()).Return(users.UserWithRoles{}, true)
				handlers.ReplaceGlobals(h)
				tr := mocks.NewTransactionsRepository(ctrl)
				tr.EXPECT().Get(gomock.Any()).Return(transactions.Transaction{}, true, nil)
				tr.EXPECT().Update(gomock.Any()).Times(0)
				transactions.ReplaceGlobals(tr)
				bu := mocks.NewBrokerUserRepository(ctrl)
				bu.EXPECT().Exists(gomock.Any()).Return(false, errors.New("error"))
				brokers.ReplaceGlobals(brokers.NewRepository(nil, bu, nil))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "broker user does not exist",
			body: validRequestBody,
			mockSetup: func(ctrl *gomock.Controller) {
				h := mocks.NewMockUtils(ctrl)
				h.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.Nil, true)
				h.EXPECT().GetUserFromContext(gomock.Any()).Return(users.UserWithRoles{}, true)
				handlers.ReplaceGlobals(h)
				tr := mocks.NewTransactionsRepository(ctrl)
				tr.EXPECT().Get(gomock.Any()).Return(transactions.Transaction{}, true, nil)
				tr.EXPECT().Update(gomock.Any()).Times(0)
				transactions.ReplaceGlobals(tr)
				bu := mocks.NewBrokerUserRepository(ctrl)
				bu.EXPECT().Exists(gomock.Any()).Return(false, nil)
				brokers.ReplaceGlobals(brokers.NewRepository(nil, bu, nil))
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "fails to update the transaction",
			body: validRequestBody,
			mockSetup: func(ctrl *gomock.Controller) {
				h := mocks.NewMockUtils(ctrl)
				h.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.Nil, true)
				h.EXPECT().GetUserFromContext(gomock.Any()).Return(users.UserWithRoles{}, true)
				handlers.ReplaceGlobals(h)
				tr := mocks.NewTransactionsRepository(ctrl)
				tr.EXPECT().Get(gomock.Any()).Return(transactions.Transaction{}, true, nil)
				tr.EXPECT().Update(gomock.Any()).Return(errors.New("error"))
				tr.EXPECT().Get(gomock.Any()).Times(0)
				transactions.ReplaceGlobals(tr)
				bu := mocks.NewBrokerUserRepository(ctrl)
				bu.EXPECT().Exists(gomock.Any()).Return(true, nil)
				brokers.ReplaceGlobals(brokers.NewRepository(nil, bu, nil))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "fails to retrieve the transaction after update",
			body: validRequestBody,
			mockSetup: func(ctrl *gomock.Controller) {
				h := mocks.NewMockUtils(ctrl)
				h.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.Nil, true)
				h.EXPECT().GetUserFromContext(gomock.Any()).Return(users.UserWithRoles{}, true)
				handlers.ReplaceGlobals(h)
				tr := mocks.NewTransactionsRepository(ctrl)
				tr.EXPECT().Get(gomock.Any()).Return(transactions.Transaction{}, true, nil)
				tr.EXPECT().Update(gomock.Any()).Return(nil)
				tr.EXPECT().Get(gomock.Any()).Return(transactions.Transaction{}, false, errors.New("error"))
				transactions.ReplaceGlobals(tr)
				bu := mocks.NewBrokerUserRepository(ctrl)
				bu.EXPECT().Exists(gomock.Any()).Return(true, nil)
				brokers.ReplaceGlobals(brokers.NewRepository(nil, bu, nil))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "fails to find the transaction after update",
			body: validRequestBody,
			mockSetup: func(ctrl *gomock.Controller) {
				h := mocks.NewMockUtils(ctrl)
				h.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.Nil, true)
				h.EXPECT().GetUserFromContext(gomock.Any()).Return(users.UserWithRoles{}, true)
				handlers.ReplaceGlobals(h)
				tr := mocks.NewTransactionsRepository(ctrl)
				tr.EXPECT().Get(gomock.Any()).Return(transactions.Transaction{}, true, nil)
				tr.EXPECT().Update(gomock.Any()).Return(nil)
				tr.EXPECT().Get(gomock.Any()).Return(transactions.Transaction{}, false, nil)
				transactions.ReplaceGlobals(tr)
				bu := mocks.NewBrokerUserRepository(ctrl)
				bu.EXPECT().Exists(gomock.Any()).Return(true, nil)
				brokers.ReplaceGlobals(brokers.NewRepository(nil, bu, nil))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "succeeded",
			body: validRequestBody,
			mockSetup: func(ctrl *gomock.Controller) {
				h := mocks.NewMockUtils(ctrl)
				h.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.Nil, true)
				h.EXPECT().GetUserFromContext(gomock.Any()).Return(users.UserWithRoles{}, true)
				handlers.ReplaceGlobals(h)
				tr := mocks.NewTransactionsRepository(ctrl)
				tr.EXPECT().Get(gomock.Any()).Return(transactions.Transaction{}, true, nil).Times(2)
				tr.EXPECT().Update(gomock.Any()).Return(nil)
				transactions.ReplaceGlobals(tr)
				bu := mocks.NewBrokerUserRepository(ctrl)
				bu.EXPECT().Exists(gomock.Any()).Return(true, nil)
				brokers.ReplaceGlobals(brokers.NewRepository(nil, bu, nil))
			},
			expectedStatus: http.StatusOK,
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest("PUT", "/api/v1/transactions/{id}", bytes.NewBuffer(tt.body))

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
				h.EXPECT().GetUserFromContext(gomock.Any()).Return(users.UserWithRoles{}, false)
				handlers.ReplaceGlobals(h)
				tr := mocks.NewTransactionsRepository(ctrl)
				tr.EXPECT().Get(gomock.Any()).Times(0)
				transactions.ReplaceGlobals(tr)
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "fails to retrieve the transaction",
			mockSetup: func(ctrl *gomock.Controller) {
				h := mocks.NewMockUtils(ctrl)
				h.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.Nil, true)
				h.EXPECT().GetUserFromContext(gomock.Any()).Return(users.UserWithRoles{}, true)
				handlers.ReplaceGlobals(h)
				tr := mocks.NewTransactionsRepository(ctrl)
				tr.EXPECT().Get(gomock.Any()).Return(transactions.Transaction{}, false, errors.New("error"))
				tr.EXPECT().Delete(gomock.Any()).Times(0)
				transactions.ReplaceGlobals(tr)
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "fails to find the transaction",
			mockSetup: func(ctrl *gomock.Controller) {
				h := mocks.NewMockUtils(ctrl)
				h.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.Nil, true)
				h.EXPECT().GetUserFromContext(gomock.Any()).Return(users.UserWithRoles{}, true)
				handlers.ReplaceGlobals(h)
				tr := mocks.NewTransactionsRepository(ctrl)
				tr.EXPECT().Get(gomock.Any()).Return(transactions.Transaction{}, false, nil)
				tr.EXPECT().Delete(gomock.Any()).Times(0)
				transactions.ReplaceGlobals(tr)
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name: "transaction user does not correspond to context user",
			mockSetup: func(ctrl *gomock.Controller) {
				h := mocks.NewMockUtils(ctrl)
				h.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.Nil, true)
				h.EXPECT().GetUserFromContext(gomock.Any()).Return(users.UserWithRoles{}, true)
				handlers.ReplaceGlobals(h)
				tr := mocks.NewTransactionsRepository(ctrl)
				tr.EXPECT().Get(gomock.Any()).Return(transactions.Transaction{UserID: uuid.New()}, true, nil)
				tr.EXPECT().Delete(gomock.Any()).Times(0)
				transactions.ReplaceGlobals(tr)
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "fails to delete",
			mockSetup: func(ctrl *gomock.Controller) {
				h := mocks.NewMockUtils(ctrl)
				h.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.Nil, true)
				h.EXPECT().GetUserFromContext(gomock.Any()).Return(users.UserWithRoles{}, true)
				handlers.ReplaceGlobals(h)
				tr := mocks.NewTransactionsRepository(ctrl)
				tr.EXPECT().Get(gomock.Any()).Return(transactions.Transaction{}, true, nil)
				tr.EXPECT().Delete(gomock.Any()).Return(errors.New("error"))
				transactions.ReplaceGlobals(tr)
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "succeeded",
			mockSetup: func(ctrl *gomock.Controller) {
				h := mocks.NewMockUtils(ctrl)
				h.EXPECT().ParseParamUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(uuid.Nil, true)
				h.EXPECT().GetUserFromContext(gomock.Any()).Return(users.UserWithRoles{}, true)
				handlers.ReplaceGlobals(h)
				tr := mocks.NewTransactionsRepository(ctrl)
				tr.EXPECT().Get(gomock.Any()).Return(transactions.Transaction{}, true, nil)
				tr.EXPECT().Delete(gomock.Any()).Return(nil)
				transactions.ReplaceGlobals(tr)
			},
			expectedStatus: http.StatusOK,
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest("DELETE", "/api/v1/transactions/{id}", nil)

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

// TestGetTransactions tests the GetTransactions handler
func TestGetTransactions(t *testing.T) {
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
				h.EXPECT().GetUserFromContext(gomock.Any()).Return(users.UserWithRoles{}, false)
				handlers.ReplaceGlobals(h)
				tr := mocks.NewTransactionsRepository(ctrl)
				tr.EXPECT().GetAll(gomock.Any()).Times(0)
				transactions.ReplaceGlobals(tr)
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "fails to retrieve all user transactions",
			mockSetup: func(ctrl *gomock.Controller) {
				h := mocks.NewMockUtils(ctrl)
				h.EXPECT().GetUserFromContext(gomock.Any()).Return(users.UserWithRoles{}, true)
				handlers.ReplaceGlobals(h)
				tr := mocks.NewTransactionsRepository(ctrl)
				tr.EXPECT().GetAll(gomock.Any()).Return([]transactions.Transaction{}, errors.New("error"))
				transactions.ReplaceGlobals(tr)
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "succeeded",
			mockSetup: func(ctrl *gomock.Controller) {
				h := mocks.NewMockUtils(ctrl)
				h.EXPECT().GetUserFromContext(gomock.Any()).Return(users.UserWithRoles{}, true)
				handlers.ReplaceGlobals(h)
				tr := mocks.NewTransactionsRepository(ctrl)
				tr.EXPECT().GetAll(gomock.Any()).Return([]transactions.Transaction{}, nil)
				transactions.ReplaceGlobals(tr)
			},
			expectedStatus: http.StatusOK,
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/api/v1/transactions", nil)

			// Apply mocks
			ctrl := gomock.NewController(t)
			tt.mockSetup(ctrl)
			defer ctrl.Finish()

			handlers.GetTransactions(w, r)
			response := w.Result()
			defer response.Body.Close()

			assert.Equal(t, tt.expectedStatus, response.StatusCode)
		})
	}
}
