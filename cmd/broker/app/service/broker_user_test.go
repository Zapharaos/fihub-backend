package service

import (
	"context"
	"errors"
	"github.com/Zapharaos/fihub-backend/cmd/broker/app/repositories"
	"github.com/Zapharaos/fihub-backend/gen/go/brokerpb"
	"github.com/Zapharaos/fihub-backend/internal/models"
	"github.com/Zapharaos/fihub-backend/test/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"
)

// TestCreateUserBroker tests the CreateUserBroker handler
func TestCreateUserBroker(t *testing.T) {
	// Prepare data
	service := &Service{}

	// Define request data
	userID := uuid.New()
	brokerID := uuid.New()
	request := &brokerpb.CreateBrokerUserRequest{
		UserId:   userID.String(),
		BrokerId: brokerID.String(),
	}

	// Test cases
	tests := []struct {
		name            string
		mockSetup       func(ctrl *gomock.Controller)
		request         *brokerpb.CreateBrokerUserRequest
		expected        *brokerpb.CreateBrokerUserResponse
		expectedErrCode codes.Code
	}{
		{
			name: "missing request body",
			mockSetup: func(ctrl *gomock.Controller) {
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().Get(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(bb, nil, nil))
			},
			request:         nil,
			expected:        &brokerpb.CreateBrokerUserResponse{},
			expectedErrCode: codes.InvalidArgument,
		},
		{
			name: "fails to parse ID from request",
			mockSetup: func(ctrl *gomock.Controller) {
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().Get(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(bb, nil, nil))
			},
			request: &brokerpb.CreateBrokerUserRequest{
				UserId: "bad-uuid",
			},
			expected:        &brokerpb.CreateBrokerUserResponse{},
			expectedErrCode: codes.InvalidArgument,
		},
		/*{
			name: "fails at bad broker user input",
			mockSetup: func(ctrl *gomock.Controller) {
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().Get(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(bb, nil, nil))
			},
			request: &gen.CreateBrokerUserRequest{
				UserId:   uuid.Nil.String(),
				BrokerId: uuid.Nil.String(),
			},
			expected:        &gen.CreateBrokerUserResponse{},
			expectedErrCode: codes.InvalidArgument,
		},*/
		{
			name: "Fails to verify the broker existence",
			mockSetup: func(ctrl *gomock.Controller) {
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().Get(gomock.Any()).Return(models.Broker{}, false, errors.New("error"))
				bu := mocks.NewBrokerUserRepository(ctrl)
				bu.EXPECT().Exists(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(bb, bu, nil))
			},
			request:         request,
			expected:        &brokerpb.CreateBrokerUserResponse{},
			expectedErrCode: codes.Internal,
		},
		{
			name: "Fails to find the broker",
			mockSetup: func(ctrl *gomock.Controller) {
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().Get(gomock.Any()).Return(models.Broker{}, false, nil)
				bu := mocks.NewBrokerUserRepository(ctrl)
				bu.EXPECT().Exists(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(bb, bu, nil))
			},
			request:         request,
			expected:        &brokerpb.CreateBrokerUserResponse{},
			expectedErrCode: codes.NotFound,
		},
		{
			name: "Broker is not enabled",
			mockSetup: func(ctrl *gomock.Controller) {
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().Get(gomock.Any()).Return(models.Broker{Disabled: true}, true, nil)
				bu := mocks.NewBrokerUserRepository(ctrl)
				bu.EXPECT().Exists(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(bb, bu, nil))
			},
			request:         request,
			expected:        &brokerpb.CreateBrokerUserResponse{},
			expectedErrCode: codes.FailedPrecondition,
		},
		{
			name: "Fails to verify broker user existence",
			mockSetup: func(ctrl *gomock.Controller) {
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().Get(gomock.Any()).Return(models.Broker{}, true, nil)
				bu := mocks.NewBrokerUserRepository(ctrl)
				bu.EXPECT().Exists(gomock.Any()).Return(true, errors.New("error"))
				bu.EXPECT().Create(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(bb, bu, nil))
			},
			request:         request,
			expected:        &brokerpb.CreateBrokerUserResponse{},
			expectedErrCode: codes.Internal,
		},
		{
			name: "BrokerUser already exists",
			mockSetup: func(ctrl *gomock.Controller) {
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().Get(gomock.Any()).Return(models.Broker{}, true, nil)
				bu := mocks.NewBrokerUserRepository(ctrl)
				bu.EXPECT().Exists(gomock.Any()).Return(true, nil)
				bu.EXPECT().Create(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(bb, bu, nil))
			},
			request:         request,
			expected:        &brokerpb.CreateBrokerUserResponse{},
			expectedErrCode: codes.AlreadyExists,
		},
		{
			name: "Fails to create user broker",
			mockSetup: func(ctrl *gomock.Controller) {
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().Get(gomock.Any()).Return(models.Broker{}, true, nil)
				bu := mocks.NewBrokerUserRepository(ctrl)
				bu.EXPECT().Exists(gomock.Any()).Return(false, nil)
				bu.EXPECT().Create(gomock.Any()).Return(errors.New("error"))
				bu.EXPECT().GetAll(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(bb, bu, nil))
			},
			request:         request,
			expected:        &brokerpb.CreateBrokerUserResponse{},
			expectedErrCode: codes.Internal,
		},
		{
			name: "Fails to retrieve all user brokers",
			mockSetup: func(ctrl *gomock.Controller) {
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().Get(gomock.Any()).Return(models.Broker{}, true, nil)
				bu := mocks.NewBrokerUserRepository(ctrl)
				bu.EXPECT().Exists(gomock.Any()).Return(false, nil)
				bu.EXPECT().Create(gomock.Any()).Return(nil)
				bu.EXPECT().GetAll(gomock.Any()).Return([]models.BrokerUser{}, errors.New("error"))
				repositories.ReplaceGlobals(repositories.NewRepository(bb, bu, nil))
			},
			request:         request,
			expected:        &brokerpb.CreateBrokerUserResponse{},
			expectedErrCode: codes.Internal,
		},
		{
			name: "Missing new broker user",
			mockSetup: func(ctrl *gomock.Controller) {
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().Get(gomock.Any()).Return(models.Broker{}, true, nil)
				bu := mocks.NewBrokerUserRepository(ctrl)
				bu.EXPECT().Exists(gomock.Any()).Return(false, nil)
				bu.EXPECT().Create(gomock.Any()).Return(nil)
				bu.EXPECT().GetAll(gomock.Any()).Return([]models.BrokerUser{}, nil)
				repositories.ReplaceGlobals(repositories.NewRepository(bb, bu, nil))
			},
			request:         request,
			expected:        &brokerpb.CreateBrokerUserResponse{},
			expectedErrCode: codes.Internal,
		},
		{
			name: "Succeeded",
			mockSetup: func(ctrl *gomock.Controller) {
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().Get(gomock.Any()).Return(models.Broker{}, true, nil)
				bu := mocks.NewBrokerUserRepository(ctrl)
				bu.EXPECT().Exists(gomock.Any()).Return(false, nil)
				bu.EXPECT().Create(gomock.Any()).Return(nil)
				bu.EXPECT().GetAll(gomock.Any()).Return([]models.BrokerUser{{
					UserID: userID,
					Broker: models.Broker{
						ID: brokerID,
					},
				}}, nil)
				repositories.ReplaceGlobals(repositories.NewRepository(bb, bu, nil))
			},
			request:         request,
			expected:        &brokerpb.CreateBrokerUserResponse{},
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
			response, err := service.CreateBrokerUser(context.Background(), tt.request)

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

// TestGetUserBroker tests the GetUserBroker handler
func TestGetUserBroker(t *testing.T) {
	// Prepare data
	service := &Service{}
	validRequest := &brokerpb.GetBrokerUserRequest{
		UserId:   uuid.New().String(),
		BrokerId: uuid.New().String(),
	}
	validResponse := models.BrokerUser{
		UserID: uuid.New(),
		Broker: models.Broker{
			ID:       uuid.New(),
			Name:     "name",
			Disabled: false,
		},
	}

	// Test cases
	tests := []struct {
		name            string
		mockSetup       func(ctrl *gomock.Controller)
		request         *brokerpb.GetBrokerUserRequest
		expected        *brokerpb.GetBrokerUserResponse
		expectedErrCode codes.Code
	}{
		{
			name: "missing request body",
			mockSetup: func(ctrl *gomock.Controller) {
				bu := mocks.NewBrokerUserRepository(ctrl)
				bu.EXPECT().Get(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(nil, bu, nil))
			},
			request:         nil,
			expected:        &brokerpb.GetBrokerUserResponse{},
			expectedErrCode: codes.InvalidArgument,
		},
		{
			name: "fails to parse ID from request",
			mockSetup: func(ctrl *gomock.Controller) {
				bu := mocks.NewBrokerUserRepository(ctrl)
				bu.EXPECT().Get(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(nil, bu, nil))
			},
			request: &brokerpb.GetBrokerUserRequest{
				UserId: "bad-uuid",
			},
			expected:        &brokerpb.GetBrokerUserResponse{},
			expectedErrCode: codes.InvalidArgument,
		},
		{
			name: "Fails to retrieve the user broker",
			mockSetup: func(ctrl *gomock.Controller) {
				bu := mocks.NewBrokerUserRepository(ctrl)
				bu.EXPECT().Get(gomock.Any()).Return(models.BrokerUser{}, false, errors.New("error"))
				repositories.ReplaceGlobals(repositories.NewRepository(nil, bu, nil))
			},
			request:         validRequest,
			expected:        &brokerpb.GetBrokerUserResponse{},
			expectedErrCode: codes.Internal,
		},
		{
			name: "Fails to find the user broker",
			mockSetup: func(ctrl *gomock.Controller) {
				bu := mocks.NewBrokerUserRepository(ctrl)
				bu.EXPECT().Get(gomock.Any()).Return(models.BrokerUser{}, false, nil)
				repositories.ReplaceGlobals(repositories.NewRepository(nil, bu, nil))
			},
			request:         validRequest,
			expected:        &brokerpb.GetBrokerUserResponse{},
			expectedErrCode: codes.NotFound,
		},
		{
			name: "Succeeded",
			mockSetup: func(ctrl *gomock.Controller) {
				bu := mocks.NewBrokerUserRepository(ctrl)
				bu.EXPECT().Get(gomock.Any()).Return(validResponse, true, nil)
				repositories.ReplaceGlobals(repositories.NewRepository(nil, bu, nil))
			},
			request:         validRequest,
			expected:        &brokerpb.GetBrokerUserResponse{},
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
			response, err := service.GetBrokerUser(context.Background(), tt.request)

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

// TestDeleteUserBroker tests the DeleteUserBroker handler
func TestDeleteUserBroker(t *testing.T) {
	// Prepare data
	service := &Service{}

	// Define request data
	userID := uuid.New()
	brokerID := uuid.New()
	request := &brokerpb.DeleteBrokerUserRequest{
		UserId:   userID.String(),
		BrokerId: brokerID.String(),
	}

	// Test cases
	tests := []struct {
		name            string
		mockSetup       func(ctrl *gomock.Controller)
		request         *brokerpb.DeleteBrokerUserRequest
		expected        *brokerpb.DeleteBrokerUserResponse
		expectedErrCode codes.Code
	}{
		{
			name: "missing request body",
			mockSetup: func(ctrl *gomock.Controller) {
				bu := mocks.NewBrokerUserRepository(ctrl)
				bu.EXPECT().Exists(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(nil, bu, nil))
			},
			request:         nil,
			expected:        &brokerpb.DeleteBrokerUserResponse{},
			expectedErrCode: codes.InvalidArgument,
		},
		{
			name: "fails to parse ID from request",
			mockSetup: func(ctrl *gomock.Controller) {
				bu := mocks.NewBrokerUserRepository(ctrl)
				bu.EXPECT().Exists(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(nil, bu, nil))
			},
			request: &brokerpb.DeleteBrokerUserRequest{
				UserId: "bad-uuid",
			},
			expected:        &brokerpb.DeleteBrokerUserResponse{},
			expectedErrCode: codes.InvalidArgument,
		},
		{
			name: "Fails to verify the user broker existence",
			mockSetup: func(ctrl *gomock.Controller) {
				bu := mocks.NewBrokerUserRepository(ctrl)
				bu.EXPECT().Exists(gomock.Any()).Return(false, errors.New("error"))
				bu.EXPECT().Delete(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(nil, bu, nil))
			},
			request:         request,
			expected:        &brokerpb.DeleteBrokerUserResponse{},
			expectedErrCode: codes.Internal,
		},
		{
			name: "Fails to find the user broker",
			mockSetup: func(ctrl *gomock.Controller) {
				bu := mocks.NewBrokerUserRepository(ctrl)
				bu.EXPECT().Exists(gomock.Any()).Return(false, nil)
				bu.EXPECT().Delete(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(nil, bu, nil))
			},
			request:         request,
			expected:        &brokerpb.DeleteBrokerUserResponse{},
			expectedErrCode: codes.NotFound,
		},
		{
			name: "Fails to delete the user broker",
			mockSetup: func(ctrl *gomock.Controller) {
				bu := mocks.NewBrokerUserRepository(ctrl)
				bu.EXPECT().Exists(gomock.Any()).Return(true, nil)
				bu.EXPECT().Delete(gomock.Any()).Return(errors.New("error"))
				repositories.ReplaceGlobals(repositories.NewRepository(nil, bu, nil))
			},
			request:         request,
			expected:        &brokerpb.DeleteBrokerUserResponse{},
			expectedErrCode: codes.Internal,
		},
		{
			name: "Succeeded",
			mockSetup: func(ctrl *gomock.Controller) {
				bu := mocks.NewBrokerUserRepository(ctrl)
				bu.EXPECT().Exists(gomock.Any()).Return(true, nil)
				bu.EXPECT().Delete(gomock.Any()).Return(nil)
				repositories.ReplaceGlobals(repositories.NewRepository(nil, bu, nil))
			},
			request:         request,
			expected:        &brokerpb.DeleteBrokerUserResponse{},
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
			response, err := service.DeleteBrokerUser(context.Background(), tt.request)

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

// TestListUserBrokers tests the ListUserBrokers handler
func TestListUserBrokers(t *testing.T) {
	// Prepare data
	service := &Service{}

	// Define request data
	userID := uuid.New()
	request := &brokerpb.ListUserBrokersRequest{
		UserId: userID.String(),
	}

	// Test cases
	tests := []struct {
		name            string
		mockSetup       func(ctrl *gomock.Controller)
		request         *brokerpb.ListUserBrokersRequest
		expected        *brokerpb.ListUserBrokersResponse
		expectedErrCode codes.Code
	}{
		{
			name: "missing request body",
			mockSetup: func(ctrl *gomock.Controller) {
				bu := mocks.NewBrokerUserRepository(ctrl)
				bu.EXPECT().GetAll(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(nil, bu, nil))
			},
			request:         nil,
			expected:        &brokerpb.ListUserBrokersResponse{},
			expectedErrCode: codes.InvalidArgument,
		},
		{
			name: "fails to parse ID from request",
			mockSetup: func(ctrl *gomock.Controller) {
				bu := mocks.NewBrokerUserRepository(ctrl)
				bu.EXPECT().GetAll(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(nil, bu, nil))
			},
			request: &brokerpb.ListUserBrokersRequest{
				UserId: "bad-uuid",
			},
			expected:        &brokerpb.ListUserBrokersResponse{},
			expectedErrCode: codes.InvalidArgument,
		},
		{
			name: "Fails to retrieve all user brokers",
			mockSetup: func(ctrl *gomock.Controller) {
				bu := mocks.NewBrokerUserRepository(ctrl)
				bu.EXPECT().GetAll(gomock.Any()).Return([]models.BrokerUser{}, errors.New("error"))
				repositories.ReplaceGlobals(repositories.NewRepository(nil, bu, nil))
			},
			request:         request,
			expected:        &brokerpb.ListUserBrokersResponse{},
			expectedErrCode: codes.Internal,
		},
		{
			name: "Succeeded",
			mockSetup: func(ctrl *gomock.Controller) {
				bu := mocks.NewBrokerUserRepository(ctrl)
				bu.EXPECT().GetAll(gomock.Any()).Return([]models.BrokerUser{}, nil)
				repositories.ReplaceGlobals(repositories.NewRepository(nil, bu, nil))
			},
			request:         request,
			expected:        &brokerpb.ListUserBrokersResponse{},
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
			response, err := service.ListUserBrokers(context.Background(), tt.request)

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
