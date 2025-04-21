package service

import (
	"context"
	"errors"
	"github.com/Zapharaos/fihub-backend/cmd/transaction/app/repositories"
	"github.com/Zapharaos/fihub-backend/internal/models"
	"github.com/Zapharaos/fihub-backend/protogen"
	"github.com/Zapharaos/fihub-backend/test/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"testing"
	"time"
)

// TestCreateTransaction tests the CreateTransaction service
func TestCreateTransaction(t *testing.T) {
	service := &Service{}

	// Define request data
	userID := uuid.New()
	brokerID := uuid.New()
	date := timestamppb.New(time.Now().AddDate(-1, 0, 0)) // 1 year in the past
	request := &protogen.CreateTransactionRequest{
		UserId:          userID.String(),
		BrokerId:        brokerID.String(),
		Date:            date,
		TransactionType: protogen.TransactionType_BUY,
		Asset:           "asset",
		Quantity:        1,
		Price:           1,
		Fee:             1,
	}

	// Define tests
	tests := []struct {
		name            string
		mockSetup       func(ctrl *gomock.Controller)
		request         *protogen.CreateTransactionRequest
		expected        *protogen.CreateTransactionResponse
		expectedErrCode codes.Code
	}{
		{
			name: "missing request body",
			mockSetup: func(ctrl *gomock.Controller) {
				tr := mocks.NewTransactionsRepository(ctrl)
				tr.EXPECT().Get(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(tr)
			},
			request: nil,
			expected: &protogen.CreateTransactionResponse{
				Transaction: nil,
			},
			expectedErrCode: codes.InvalidArgument,
		},
		{
			name: "fails to parse ID from request",
			mockSetup: func(ctrl *gomock.Controller) {
				tr := mocks.NewTransactionsRepository(ctrl)
				tr.EXPECT().Get(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(tr)
			},
			request: &protogen.CreateTransactionRequest{
				UserId: "bad-uuid",
			},
			expected: &protogen.CreateTransactionResponse{
				Transaction: nil,
			},
			expectedErrCode: codes.InvalidArgument,
		},
		{
			name: "fails at bad transaction input",
			mockSetup: func(ctrl *gomock.Controller) {
				tr := mocks.NewTransactionsRepository(ctrl)
				tr.EXPECT().Get(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(tr)
			},
			request: &protogen.CreateTransactionRequest{
				UserId:          userID.String(),
				BrokerId:        brokerID.String(),
				TransactionType: protogen.TransactionType_TRANSACTION_TYPE_UNSPECIFIED,
			},
			expected: &protogen.CreateTransactionResponse{
				Transaction: nil,
			},
			expectedErrCode: codes.InvalidArgument,
		},
		{
			name: "fails at transactions creation",
			mockSetup: func(ctrl *gomock.Controller) {
				tr := mocks.NewTransactionsRepository(ctrl)
				tr.EXPECT().Create(gomock.Any()).Return(uuid.Nil, errors.New("error"))
				tr.EXPECT().Get(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(tr)
			},
			request: request,
			expected: &protogen.CreateTransactionResponse{
				Transaction: nil,
			},
			expectedErrCode: codes.Internal,
		},
		{
			name: "fails to retrieve the transaction",
			mockSetup: func(ctrl *gomock.Controller) {
				tr := mocks.NewTransactionsRepository(ctrl)
				tr.EXPECT().Create(gomock.Any()).Return(uuid.New(), nil)
				tr.EXPECT().Get(gomock.Any()).Return(models.Transaction{}, false, errors.New("error"))
				repositories.ReplaceGlobals(tr)
			},
			request: request,
			expected: &protogen.CreateTransactionResponse{
				Transaction: nil,
			},
			expectedErrCode: codes.Internal,
		},
		{
			name: "could not find the transaction",
			mockSetup: func(ctrl *gomock.Controller) {
				tr := mocks.NewTransactionsRepository(ctrl)
				tr.EXPECT().Create(gomock.Any()).Return(uuid.New(), nil)
				tr.EXPECT().Get(gomock.Any()).Return(models.Transaction{}, false, nil)
				repositories.ReplaceGlobals(tr)
			},
			request: request,
			expected: &protogen.CreateTransactionResponse{
				Transaction: nil,
			},
			expectedErrCode: codes.NotFound,
		},
		{
			name: "succeeded",
			mockSetup: func(ctrl *gomock.Controller) {
				tr := mocks.NewTransactionsRepository(ctrl)
				tr.EXPECT().Create(gomock.Any()).Return(uuid.New(), nil)
				tr.EXPECT().Get(gomock.Any()).Return(models.Transaction{}, true, nil)
				repositories.ReplaceGlobals(tr)
			},
			request: request,
			expected: &protogen.CreateTransactionResponse{
				Transaction: &protogen.Transaction{},
			},
			expectedErrCode: codes.OK,
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Apply mocks
			ctrl := gomock.NewController(t)
			tt.mockSetup(ctrl)
			defer ctrl.Finish()

			// Call service
			response, err := service.CreateTransaction(context.Background(), tt.request)

			// Handle errors
			if err != nil && tt.expectedErrCode == codes.OK {
				assert.Fail(t, "unexpected error", err)
			} else if err != nil {
				if s, ok := status.FromError(err); ok {
					assert.Equal(t, tt.expectedErrCode, s.Code())
				} else {
					assert.Fail(t, "failed to get status from error")
				}
			}

			// Handle response
			if tt.expectedErrCode == codes.OK {
				assert.NotNil(t, response)
			} else {
				assert.Equal(t, tt.expected, response)
			}
		})
	}
}

// TestGetTransaction tests the GetTransaction service
func TestGetTransaction(t *testing.T) {
	service := &Service{}

	// Define request data
	transactionID := uuid.New()
	request := &protogen.GetTransactionRequest{
		TransactionId: transactionID.String(),
	}

	// Define tests
	tests := []struct {
		name            string
		mockSetup       func(ctrl *gomock.Controller)
		request         *protogen.GetTransactionRequest
		expected        *protogen.GetTransactionResponse
		expectedErrCode codes.Code
	}{
		{
			name: "missing request body",
			mockSetup: func(ctrl *gomock.Controller) {
				tr := mocks.NewTransactionsRepository(ctrl)
				tr.EXPECT().Get(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(tr)
			},
			request: nil,
			expected: &protogen.GetTransactionResponse{
				Transaction: nil,
			},
			expectedErrCode: codes.InvalidArgument,
		},
		{
			name: "fails to parse ID from request",
			mockSetup: func(ctrl *gomock.Controller) {
				tr := mocks.NewTransactionsRepository(ctrl)
				tr.EXPECT().Get(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(tr)
			},
			request: &protogen.GetTransactionRequest{
				TransactionId: "bad-uuid",
			},
			expected: &protogen.GetTransactionResponse{
				Transaction: nil,
			},
			expectedErrCode: codes.InvalidArgument,
		},
		{
			name: "fails to retrieve the transaction",
			mockSetup: func(ctrl *gomock.Controller) {
				tr := mocks.NewTransactionsRepository(ctrl)
				tr.EXPECT().Get(gomock.Any()).Return(models.Transaction{}, false, errors.New("error"))
				repositories.ReplaceGlobals(tr)
			},
			request: request,
			expected: &protogen.GetTransactionResponse{
				Transaction: nil,
			},
			expectedErrCode: codes.Internal,
		},
		{
			name: "could not find the transaction",
			mockSetup: func(ctrl *gomock.Controller) {
				tr := mocks.NewTransactionsRepository(ctrl)
				tr.EXPECT().Get(gomock.Any()).Return(models.Transaction{}, false, nil)
				repositories.ReplaceGlobals(tr)
			},
			request: request,
			expected: &protogen.GetTransactionResponse{
				Transaction: nil,
			},
			expectedErrCode: codes.NotFound,
		},
		{
			name: "succeeded",
			mockSetup: func(ctrl *gomock.Controller) {
				tr := mocks.NewTransactionsRepository(ctrl)
				tr.EXPECT().Get(gomock.Any()).Return(models.Transaction{
					ID: transactionID,
				}, true, nil)
				repositories.ReplaceGlobals(tr)
			},
			request: request,
			expected: &protogen.GetTransactionResponse{
				Transaction: &protogen.Transaction{
					Id: transactionID.String(),
				},
			},
			expectedErrCode: codes.OK,
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Apply mocks
			ctrl := gomock.NewController(t)
			tt.mockSetup(ctrl)
			defer ctrl.Finish()

			// Call service
			response, err := service.GetTransaction(context.Background(), tt.request)

			// Handle errors
			if err != nil && tt.expectedErrCode == codes.OK {
				assert.Fail(t, "unexpected error", err)
			} else if err != nil {
				if s, ok := status.FromError(err); ok {
					assert.Equal(t, tt.expectedErrCode, s.Code())
				} else {
					assert.Fail(t, "failed to get status from error")
				}
			}

			// Handle response
			if tt.expectedErrCode == codes.OK {
				assert.NotNil(t, response)
				assert.Equal(t, tt.expected.Transaction.Id, response.Transaction.Id)
			} else {
				assert.Equal(t, tt.expected, response)
			}
		})
	}
}

// TestListTransactions tests the ListTransactions service
func TestListTransactions(t *testing.T) {
	service := &Service{}

	// Define request data
	userID := uuid.New()
	request := &protogen.ListTransactionsRequest{
		UserId: userID.String(),
	}

	// Define tests
	tests := []struct {
		name            string
		mockSetup       func(ctrl *gomock.Controller)
		request         *protogen.ListTransactionsRequest
		expected        *protogen.ListTransactionsResponse
		expectedErrCode codes.Code
	}{
		{
			name: "missing request body",
			mockSetup: func(ctrl *gomock.Controller) {
				tr := mocks.NewTransactionsRepository(ctrl)
				tr.EXPECT().GetAll(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(tr)
			},
			request: nil,
			expected: &protogen.ListTransactionsResponse{
				Transactions: nil,
			},
			expectedErrCode: codes.InvalidArgument,
		},
		{
			name: "fails to parse ID from request",
			mockSetup: func(ctrl *gomock.Controller) {
				tr := mocks.NewTransactionsRepository(ctrl)
				tr.EXPECT().GetAll(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(tr)
			},
			request: &protogen.ListTransactionsRequest{
				UserId: "bad-uuid",
			},
			expected: &protogen.ListTransactionsResponse{
				Transactions: nil,
			},
			expectedErrCode: codes.InvalidArgument,
		},
		{
			name: "fails to list the transactions",
			mockSetup: func(ctrl *gomock.Controller) {
				tr := mocks.NewTransactionsRepository(ctrl)
				tr.EXPECT().GetAll(gomock.Any()).Return([]models.Transaction{}, errors.New("error"))
				repositories.ReplaceGlobals(tr)
			},
			request: request,
			expected: &protogen.ListTransactionsResponse{
				Transactions: nil,
			},
			expectedErrCode: codes.Internal,
		},
		{
			name: "succeeded",
			mockSetup: func(ctrl *gomock.Controller) {
				tr := mocks.NewTransactionsRepository(ctrl)
				tr.EXPECT().GetAll(gomock.Any()).Return([]models.Transaction{
					{UserID: userID},
					{UserID: userID},
					{UserID: userID},
				}, nil)
				repositories.ReplaceGlobals(tr)
			},
			request: request,
			expected: &protogen.ListTransactionsResponse{
				Transactions: []*protogen.Transaction{
					{UserId: userID.String()},
					{UserId: userID.String()},
					{UserId: userID.String()},
				},
			},
			expectedErrCode: codes.OK,
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Apply mocks
			ctrl := gomock.NewController(t)
			tt.mockSetup(ctrl)
			defer ctrl.Finish()

			// Call service
			response, err := service.ListTransactions(context.Background(), tt.request)

			// Handle errors
			if err != nil && tt.expectedErrCode == codes.OK {
				assert.Fail(t, "unexpected error", err)
			} else if err != nil {
				if s, ok := status.FromError(err); ok {
					assert.Equal(t, tt.expectedErrCode, s.Code())
				} else {
					assert.Fail(t, "failed to get status from error")
				}
			}

			// Handle response
			if tt.expectedErrCode == codes.OK {
				assert.NotNil(t, response)
				assert.Equal(t, len(tt.expected.Transactions), len(response.Transactions))
				assert.Equal(t, tt.expected.Transactions[0].UserId, response.Transactions[0].UserId)
			} else {
				assert.Equal(t, tt.expected, response)
			}
		})
	}
}

// TestUpdateTransaction tests the UpdateTransaction handler
func TestUpdateTransaction(t *testing.T) {
	service := &Service{}

	// Define request data
	transactionID := uuid.New()
	userID := uuid.New()
	brokerID := uuid.New()
	date := timestamppb.New(time.Now().AddDate(-1, 0, 0)) // 1 year in the past
	request := &protogen.UpdateTransactionRequest{
		TransactionId:   transactionID.String(),
		UserId:          userID.String(),
		BrokerId:        brokerID.String(),
		Date:            date,
		TransactionType: protogen.TransactionType_BUY,
		Asset:           "asset",
		Quantity:        1,
		Price:           1,
		Fee:             1,
	}

	// Define tests
	tests := []struct {
		name            string
		mockSetup       func(ctrl *gomock.Controller)
		request         *protogen.UpdateTransactionRequest
		expected        *protogen.UpdateTransactionResponse
		expectedErrCode codes.Code
	}{
		{
			name: "missing request body",
			mockSetup: func(ctrl *gomock.Controller) {
				tr := mocks.NewTransactionsRepository(ctrl)
				tr.EXPECT().Get(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(tr)
			},
			request:         nil,
			expectedErrCode: codes.InvalidArgument,
		},
		{
			name: "fails to parse ID from request",
			mockSetup: func(ctrl *gomock.Controller) {
				tr := mocks.NewTransactionsRepository(ctrl)
				tr.EXPECT().Get(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(tr)
			},
			request: &protogen.UpdateTransactionRequest{
				UserId: "bad-uuid",
			},
			expectedErrCode: codes.InvalidArgument,
		},
		{
			name: "fails at bad transaction input",
			mockSetup: func(ctrl *gomock.Controller) {
				tr := mocks.NewTransactionsRepository(ctrl)
				tr.EXPECT().Get(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(tr)
			},
			request: &protogen.UpdateTransactionRequest{
				TransactionId:   transactionID.String(),
				UserId:          userID.String(),
				BrokerId:        brokerID.String(),
				TransactionType: protogen.TransactionType_TRANSACTION_TYPE_UNSPECIFIED,
			},
			expected: &protogen.UpdateTransactionResponse{
				Transaction: nil,
			},
			expectedErrCode: codes.InvalidArgument,
		},
		{
			name: "fails to retrieve the transaction",
			mockSetup: func(ctrl *gomock.Controller) {
				tr := mocks.NewTransactionsRepository(ctrl)
				tr.EXPECT().Get(gomock.Any()).Return(models.Transaction{}, false, errors.New("error"))
				tr.EXPECT().Update(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(tr)
			},
			request:         request,
			expectedErrCode: codes.Internal,
		},
		{
			name: "could not find the transaction",
			mockSetup: func(ctrl *gomock.Controller) {
				tr := mocks.NewTransactionsRepository(ctrl)
				tr.EXPECT().Get(gomock.Any()).Return(models.Transaction{}, false, nil)
				tr.EXPECT().Update(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(tr)
			},
			request:         request,
			expectedErrCode: codes.NotFound,
		},
		{
			name: "transaction user does not correspond to context user",
			mockSetup: func(ctrl *gomock.Controller) {
				tr := mocks.NewTransactionsRepository(ctrl)
				tr.EXPECT().Get(gomock.Any()).Return(models.Transaction{UserID: uuid.New()}, true, nil)
				tr.EXPECT().Update(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(tr)
			},
			request:         request,
			expectedErrCode: codes.PermissionDenied,
		},
		{
			name: "fails to update the transaction",
			mockSetup: func(ctrl *gomock.Controller) {
				tr := mocks.NewTransactionsRepository(ctrl)
				tr.EXPECT().Get(gomock.Any()).Return(models.Transaction{UserID: userID}, true, nil)
				tr.EXPECT().Update(gomock.Any()).Return(errors.New("error"))
				tr.EXPECT().Get(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(tr)
			},
			request:         request,
			expectedErrCode: codes.Internal,
		},
		{
			name: "fails to retrieve the transaction after update",
			mockSetup: func(ctrl *gomock.Controller) {
				tr := mocks.NewTransactionsRepository(ctrl)
				tr.EXPECT().Get(gomock.Any()).Return(models.Transaction{UserID: userID}, true, nil)
				tr.EXPECT().Update(gomock.Any()).Return(nil)
				tr.EXPECT().Get(gomock.Any()).Return(models.Transaction{}, false, errors.New("error"))
				repositories.ReplaceGlobals(tr)
			},
			request:         request,
			expectedErrCode: codes.Internal,
		},
		{
			name: "fails to find the transaction after update",
			mockSetup: func(ctrl *gomock.Controller) {
				tr := mocks.NewTransactionsRepository(ctrl)
				tr.EXPECT().Get(gomock.Any()).Return(models.Transaction{UserID: userID}, true, nil)
				tr.EXPECT().Update(gomock.Any()).Return(nil)
				tr.EXPECT().Get(gomock.Any()).Return(models.Transaction{}, false, nil)
				repositories.ReplaceGlobals(tr)
			},
			request:         request,
			expectedErrCode: codes.NotFound,
		},
		{
			name: "succeeded",
			mockSetup: func(ctrl *gomock.Controller) {
				tr := mocks.NewTransactionsRepository(ctrl)
				tr.EXPECT().Get(gomock.Any()).Return(models.Transaction{
					UserID: userID,
					Price:  request.Price + 1,
				}, true, nil).Times(2)
				tr.EXPECT().Update(gomock.Any()).Return(nil)
				repositories.ReplaceGlobals(tr)
			},
			request: request,
			expected: &protogen.UpdateTransactionResponse{
				Transaction: &protogen.Transaction{
					Price: request.Price + 1,
				},
			},
			expectedErrCode: codes.OK,
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Apply mocks
			ctrl := gomock.NewController(t)
			tt.mockSetup(ctrl)
			defer ctrl.Finish()

			// Call service
			response, err := service.UpdateTransaction(context.Background(), tt.request)

			// Handle errors
			if err != nil && tt.expectedErrCode == codes.OK {
				assert.Fail(t, "unexpected error", err)
			} else if err != nil {
				if s, ok := status.FromError(err); ok {
					assert.Equal(t, tt.expectedErrCode, s.Code())
				} else {
					assert.Fail(t, "failed to get status from error")
				}
			}

			// Handle response
			if tt.expectedErrCode == codes.OK {
				assert.NotNil(t, response)
				assert.Equal(t, tt.expected.Transaction.Price, response.Transaction.Price)
				assert.NotEqual(t, tt.expected.Transaction.Price, request.Price)
			}
		})
	}
}

// TestDeleteTransaction tests the DeleteTransaction handler
func TestDeleteTransaction(t *testing.T) {
	service := &Service{}

	// Define request data
	userID := uuid.New()
	transactionID := uuid.New()
	request := &protogen.DeleteTransactionRequest{
		UserId:        userID.String(),
		TransactionId: transactionID.String(),
	}

	// Define tests
	tests := []struct {
		name            string
		mockSetup       func(ctrl *gomock.Controller)
		request         *protogen.DeleteTransactionRequest
		expectedErrCode codes.Code
	}{
		{
			name: "missing request body",
			mockSetup: func(ctrl *gomock.Controller) {
				tr := mocks.NewTransactionsRepository(ctrl)
				tr.EXPECT().Get(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(tr)
			},
			request:         nil,
			expectedErrCode: codes.InvalidArgument,
		},
		{
			name: "fails to parse ID from request",
			mockSetup: func(ctrl *gomock.Controller) {
				tr := mocks.NewTransactionsRepository(ctrl)
				tr.EXPECT().Get(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(tr)
			},
			request: &protogen.DeleteTransactionRequest{
				UserId: "bad-uuid",
			},
			expectedErrCode: codes.InvalidArgument,
		},
		{
			name: "fails to retrieve the transaction",
			mockSetup: func(ctrl *gomock.Controller) {
				tr := mocks.NewTransactionsRepository(ctrl)
				tr.EXPECT().Get(gomock.Any()).Return(models.Transaction{}, false, errors.New("error"))
				tr.EXPECT().Delete(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(tr)
			},
			request:         request,
			expectedErrCode: codes.Internal,
		},
		{
			name: "could not find the transaction",
			mockSetup: func(ctrl *gomock.Controller) {
				tr := mocks.NewTransactionsRepository(ctrl)
				tr.EXPECT().Get(gomock.Any()).Return(models.Transaction{}, false, nil)
				tr.EXPECT().Delete(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(tr)
			},
			request:         request,
			expectedErrCode: codes.NotFound,
		},
		{
			name: "transaction user does not correspond to context user",
			mockSetup: func(ctrl *gomock.Controller) {
				tr := mocks.NewTransactionsRepository(ctrl)
				tr.EXPECT().Get(gomock.Any()).Return(models.Transaction{UserID: uuid.New()}, true, nil)
				tr.EXPECT().Delete(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(tr)
			},
			request:         request,
			expectedErrCode: codes.PermissionDenied,
		},
		{
			name: "fails to delete the transaction",
			mockSetup: func(ctrl *gomock.Controller) {
				tr := mocks.NewTransactionsRepository(ctrl)
				tr.EXPECT().Get(gomock.Any()).Return(models.Transaction{UserID: userID}, true, nil)
				tr.EXPECT().Delete(gomock.Any()).Return(errors.New("error"))
				repositories.ReplaceGlobals(tr)
			},
			request:         request,
			expectedErrCode: codes.Internal,
		},
		{
			name: "succeeded",
			mockSetup: func(ctrl *gomock.Controller) {
				tr := mocks.NewTransactionsRepository(ctrl)
				tr.EXPECT().Get(gomock.Any()).Return(models.Transaction{UserID: userID}, true, nil)
				tr.EXPECT().Delete(gomock.Any()).Return(nil)
				repositories.ReplaceGlobals(tr)
			},
			request:         request,
			expectedErrCode: codes.OK,
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Apply mocks
			ctrl := gomock.NewController(t)
			tt.mockSetup(ctrl)
			defer ctrl.Finish()

			// Call service
			response, err := service.DeleteTransaction(context.Background(), tt.request)

			// Handle errors
			if err != nil && tt.expectedErrCode == codes.OK {
				assert.Fail(t, "unexpected error", err)
			} else if err != nil {
				if s, ok := status.FromError(err); ok {
					assert.Equal(t, tt.expectedErrCode, s.Code())
				} else {
					assert.Fail(t, "failed to get status from error")
				}
			}

			// Handle response
			if tt.expectedErrCode == codes.OK {
				assert.NotNil(t, response)
			}
		})
	}
}

// TestDeleteTransactionByBroker tests the DeleteTransactionByBroker handler
func TestDeleteTransactionByBroker(t *testing.T) {
	service := &Service{}

	// Define request data
	userID := uuid.New()
	brokerID := uuid.New()
	request := &protogen.DeleteTransactionByBrokerRequest{
		UserId:   userID.String(),
		BrokerId: brokerID.String(),
	}

	// Define tests
	tests := []struct {
		name            string
		mockSetup       func(ctrl *gomock.Controller)
		request         *protogen.DeleteTransactionByBrokerRequest
		expectedErrCode codes.Code
	}{
		{
			name: "missing request body",
			mockSetup: func(ctrl *gomock.Controller) {
				tr := mocks.NewTransactionsRepository(ctrl)
				tr.EXPECT().DeleteByBroker(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(tr)
			},
			request:         nil,
			expectedErrCode: codes.InvalidArgument,
		},
		{
			name: "fails to parse ID from request",
			mockSetup: func(ctrl *gomock.Controller) {
				tr := mocks.NewTransactionsRepository(ctrl)
				tr.EXPECT().DeleteByBroker(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(tr)
			},
			request: &protogen.DeleteTransactionByBrokerRequest{
				UserId: "bad-uuid",
			},
			expectedErrCode: codes.InvalidArgument,
		},
		{
			name: "fails to delete the transaction",
			mockSetup: func(ctrl *gomock.Controller) {
				tr := mocks.NewTransactionsRepository(ctrl)
				tr.EXPECT().DeleteByBroker(gomock.Any()).Return(errors.New("error"))
				repositories.ReplaceGlobals(tr)
			},
			request:         request,
			expectedErrCode: codes.Internal,
		},
		{
			name: "succeeded",
			mockSetup: func(ctrl *gomock.Controller) {
				tr := mocks.NewTransactionsRepository(ctrl)
				tr.EXPECT().DeleteByBroker(gomock.Any()).Return(nil)
				repositories.ReplaceGlobals(tr)
			},
			request:         request,
			expectedErrCode: codes.OK,
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Apply mocks
			ctrl := gomock.NewController(t)
			tt.mockSetup(ctrl)
			defer ctrl.Finish()

			// Call service
			response, err := service.DeleteTransactionByBroker(context.Background(), tt.request)

			// Handle errors
			if err != nil && tt.expectedErrCode == codes.OK {
				assert.Fail(t, "unexpected error", err)
			} else if err != nil {
				if s, ok := status.FromError(err); ok {
					assert.Equal(t, tt.expectedErrCode, s.Code())
				} else {
					assert.Fail(t, "failed to get status from error")
				}
			}

			// Handle response
			if tt.expectedErrCode == codes.OK {
				assert.NotNil(t, response)
			}
		})
	}
}
