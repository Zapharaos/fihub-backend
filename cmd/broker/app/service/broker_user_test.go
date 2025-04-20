package service

import (
	"context"
	"errors"
	"github.com/Zapharaos/fihub-backend/cmd/broker/app/repositories"
	"github.com/Zapharaos/fihub-backend/internal/models"
	"github.com/Zapharaos/fihub-backend/protogen"
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
	request := &protogen.CreateBrokerUserRequest{
		UserId:   userID.String(),
		BrokerId: brokerID.String(),
	}

	// Test cases
	tests := []struct {
		name            string
		mockSetup       func(ctrl *gomock.Controller)
		request         *protogen.CreateBrokerUserRequest
		expected        *protogen.CreateBrokerUserResponse
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
			expected:        &protogen.CreateBrokerUserResponse{},
			expectedErrCode: codes.InvalidArgument,
		},
		{
			name: "fails to parse ID from request",
			mockSetup: func(ctrl *gomock.Controller) {
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().Get(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(bb, nil, nil))
			},
			request: &protogen.CreateBrokerUserRequest{
				UserId: "bad-uuid",
			},
			expected:        &protogen.CreateBrokerUserResponse{},
			expectedErrCode: codes.InvalidArgument,
		},
		/*{
			name: "fails at bad broker user input",
			mockSetup: func(ctrl *gomock.Controller) {
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().Get(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(bb, nil, nil))
			},
			request: &protogen.CreateBrokerUserRequest{
				UserId:   uuid.Nil.String(),
				BrokerId: uuid.Nil.String(),
			},
			expected:        &protogen.CreateBrokerUserResponse{},
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
			expected:        &protogen.CreateBrokerUserResponse{},
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
			expected:        &protogen.CreateBrokerUserResponse{},
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
			expected:        &protogen.CreateBrokerUserResponse{},
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
			expected:        &protogen.CreateBrokerUserResponse{},
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
			expected:        &protogen.CreateBrokerUserResponse{},
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
			expected:        &protogen.CreateBrokerUserResponse{},
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
			expected:        &protogen.CreateBrokerUserResponse{},
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
			expected:        &protogen.CreateBrokerUserResponse{},
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
			expected:        &protogen.CreateBrokerUserResponse{},
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

// TestDeleteUserBroker tests the DeleteUserBroker handler
func TestDeleteUserBroker(t *testing.T) {
	// Prepare data
	service := &Service{}

	// Define request data
	userID := uuid.New()
	brokerID := uuid.New()
	request := &protogen.DeleteBrokerUserRequest{
		UserId:   userID.String(),
		BrokerId: brokerID.String(),
	}

	// Test cases
	tests := []struct {
		name            string
		mockSetup       func(ctrl *gomock.Controller)
		request         *protogen.DeleteBrokerUserRequest
		expected        *protogen.DeleteBrokerUserResponse
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
			expected:        &protogen.DeleteBrokerUserResponse{},
			expectedErrCode: codes.InvalidArgument,
		},
		{
			name: "fails to parse ID from request",
			mockSetup: func(ctrl *gomock.Controller) {
				bu := mocks.NewBrokerUserRepository(ctrl)
				bu.EXPECT().Exists(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(nil, bu, nil))
			},
			request: &protogen.DeleteBrokerUserRequest{
				UserId: "bad-uuid",
			},
			expected:        &protogen.DeleteBrokerUserResponse{},
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
			expected:        &protogen.DeleteBrokerUserResponse{},
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
			expected:        &protogen.DeleteBrokerUserResponse{},
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
			expected:        &protogen.DeleteBrokerUserResponse{},
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
			expected:        &protogen.DeleteBrokerUserResponse{},
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

// TestGetUserBrokers tests the GetUserBrokers handler
func TestGetUserBrokers(t *testing.T) {
	// Prepare data
	service := &Service{}

	// Define request data
	userID := uuid.New()
	request := &protogen.ListUserBrokersRequest{
		UserId: userID.String(),
	}

	// Test cases
	tests := []struct {
		name            string
		mockSetup       func(ctrl *gomock.Controller)
		request         *protogen.ListUserBrokersRequest
		expected        *protogen.ListUserBrokersResponse
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
			expected:        &protogen.ListUserBrokersResponse{},
			expectedErrCode: codes.InvalidArgument,
		},
		{
			name: "fails to parse ID from request",
			mockSetup: func(ctrl *gomock.Controller) {
				bu := mocks.NewBrokerUserRepository(ctrl)
				bu.EXPECT().GetAll(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(nil, bu, nil))
			},
			request: &protogen.ListUserBrokersRequest{
				UserId: "bad-uuid",
			},
			expected:        &protogen.ListUserBrokersResponse{},
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
			expected:        &protogen.ListUserBrokersResponse{},
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
			expected:        &protogen.ListUserBrokersResponse{},
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
			response, err := service.GetUserBrokers(context.Background(), tt.request)

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
