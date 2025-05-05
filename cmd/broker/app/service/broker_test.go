package service

import (
	"context"
	"errors"
	"github.com/Zapharaos/fihub-backend/cmd/broker/app/repositories"
	"github.com/Zapharaos/fihub-backend/gen/go/brokerpb"
	"github.com/Zapharaos/fihub-backend/gen/go/securitypb"
	"github.com/Zapharaos/fihub-backend/internal/models"
	"github.com/Zapharaos/fihub-backend/internal/security"
	"github.com/Zapharaos/fihub-backend/test/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"
)

// TestCreateBroker tests the function CreateBroker
func TestCreateBroker(t *testing.T) {
	// Prepare data
	service := &Service{}
	validRequest := &brokerpb.CreateBrokerRequest{
		Name:     "name",
		Disabled: false,
	}

	// Test cases
	tests := []struct {
		name            string
		mockSetup       func(ctrl *gomock.Controller)
		request         *brokerpb.CreateBrokerRequest
		expected        *brokerpb.CreateBrokerResponse
		expectedErrCode codes.Code
	}{
		{
			name: "does not have permission",
			mockSetup: func(ctrl *gomock.Controller) {
				// Mock the public security facade
				publicSecurityClient := mocks.NewMockPublicSecurityServiceClient(ctrl)
				publicSecurityClient.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(&securitypb.CheckPermissionResponse{HasPermission: false}, nil)
				security.ReplaceGlobals(security.NewPublicSecurityFacadeWithGrpcClient(publicSecurityClient))
			},
			request:         validRequest,
			expected:        &brokerpb.CreateBrokerResponse{},
			expectedErrCode: codes.PermissionDenied,
		},
		{
			name: "fails at bad broker input",
			mockSetup: func(ctrl *gomock.Controller) {
				// Mock the public security facade
				publicSecurityClient := mocks.NewMockPublicSecurityServiceClient(ctrl)
				publicSecurityClient.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(&securitypb.CheckPermissionResponse{HasPermission: true}, nil)
				security.ReplaceGlobals(security.NewPublicSecurityFacadeWithGrpcClient(publicSecurityClient))
				// Mock the broker repository
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().ExistsByName(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(bb, nil, nil))
			},
			request: &brokerpb.CreateBrokerRequest{
				Name:     "",
				Disabled: false,
			},
			expected:        &brokerpb.CreateBrokerResponse{},
			expectedErrCode: codes.InvalidArgument,
		},
		{
			name: "fails to verify broker existence",
			mockSetup: func(ctrl *gomock.Controller) {
				// Mock the public security facade
				publicSecurityClient := mocks.NewMockPublicSecurityServiceClient(ctrl)
				publicSecurityClient.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(&securitypb.CheckPermissionResponse{HasPermission: true}, nil)
				security.ReplaceGlobals(security.NewPublicSecurityFacadeWithGrpcClient(publicSecurityClient))
				// Mock the broker repository
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().ExistsByName(gomock.Any()).Return(false, errors.New("error"))
				bb.EXPECT().Create(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(bb, nil, nil))
			},
			request:         validRequest,
			expected:        &brokerpb.CreateBrokerResponse{},
			expectedErrCode: codes.Internal,
		},
		{
			name: "broker already exists",
			mockSetup: func(ctrl *gomock.Controller) {
				// Mock the public security facade
				publicSecurityClient := mocks.NewMockPublicSecurityServiceClient(ctrl)
				publicSecurityClient.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(&securitypb.CheckPermissionResponse{HasPermission: true}, nil)
				security.ReplaceGlobals(security.NewPublicSecurityFacadeWithGrpcClient(publicSecurityClient))
				// Mock the broker repository
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().ExistsByName(gomock.Any()).Return(true, nil)
				bb.EXPECT().Create(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(bb, nil, nil))
			},
			request:         validRequest,
			expected:        &brokerpb.CreateBrokerResponse{},
			expectedErrCode: codes.AlreadyExists,
		},
		{
			name: "fails to create broker",
			mockSetup: func(ctrl *gomock.Controller) {
				// Mock the public security facade
				publicSecurityClient := mocks.NewMockPublicSecurityServiceClient(ctrl)
				publicSecurityClient.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(&securitypb.CheckPermissionResponse{HasPermission: true}, nil)
				security.ReplaceGlobals(security.NewPublicSecurityFacadeWithGrpcClient(publicSecurityClient))
				// Mock the broker repository
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().ExistsByName(gomock.Any()).Return(false, nil)
				bb.EXPECT().Create(gomock.Any()).Return(uuid.Nil, errors.New("error"))
				bb.EXPECT().Get(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(bb, nil, nil))
			},
			request:         validRequest,
			expected:        &brokerpb.CreateBrokerResponse{},
			expectedErrCode: codes.Internal,
		},
		{
			name: "fails to retrieve the broker",
			mockSetup: func(ctrl *gomock.Controller) {
				// Mock the public security facade
				publicSecurityClient := mocks.NewMockPublicSecurityServiceClient(ctrl)
				publicSecurityClient.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(&securitypb.CheckPermissionResponse{HasPermission: true}, nil)
				security.ReplaceGlobals(security.NewPublicSecurityFacadeWithGrpcClient(publicSecurityClient))
				// Mock the broker repository
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().ExistsByName(gomock.Any()).Return(false, nil)
				bb.EXPECT().Create(gomock.Any()).Return(uuid.Nil, nil)
				bb.EXPECT().Get(gomock.Any()).Return(models.Broker{}, false, errors.New("error"))
				repositories.ReplaceGlobals(repositories.NewRepository(bb, nil, nil))
			},
			request:         validRequest,
			expected:        &brokerpb.CreateBrokerResponse{},
			expectedErrCode: codes.Internal,
		},
		{
			name: "fails to find broker",
			mockSetup: func(ctrl *gomock.Controller) {
				// Mock the public security facade
				publicSecurityClient := mocks.NewMockPublicSecurityServiceClient(ctrl)
				publicSecurityClient.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(&securitypb.CheckPermissionResponse{HasPermission: true}, nil)
				security.ReplaceGlobals(security.NewPublicSecurityFacadeWithGrpcClient(publicSecurityClient))
				// Mock the broker repository
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().ExistsByName(gomock.Any()).Return(false, nil)
				bb.EXPECT().Create(gomock.Any()).Return(uuid.Nil, nil)
				bb.EXPECT().Get(gomock.Any()).Return(models.Broker{}, false, nil)
				repositories.ReplaceGlobals(repositories.NewRepository(bb, nil, nil))
			},
			request:         validRequest,
			expected:        &brokerpb.CreateBrokerResponse{},
			expectedErrCode: codes.Internal,
		},
		{
			name: "succeeded",
			mockSetup: func(ctrl *gomock.Controller) {
				// Mock the public security facade
				publicSecurityClient := mocks.NewMockPublicSecurityServiceClient(ctrl)
				publicSecurityClient.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(&securitypb.CheckPermissionResponse{HasPermission: true}, nil)
				security.ReplaceGlobals(security.NewPublicSecurityFacadeWithGrpcClient(publicSecurityClient))
				// Mock the broker repository
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().ExistsByName(gomock.Any()).Return(false, nil)
				bb.EXPECT().Create(gomock.Any()).Return(uuid.Nil, nil)
				bb.EXPECT().Get(gomock.Any()).Return(models.Broker{}, true, nil)
				repositories.ReplaceGlobals(repositories.NewRepository(bb, nil, nil))
			},
			request:         validRequest,
			expected:        &brokerpb.CreateBrokerResponse{},
			expectedErrCode: codes.OK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Apply mocks
			ctrl := gomock.NewController(t)
			tt.mockSetup(ctrl)
			defer ctrl.Finish()

			// Call service
			response, err := service.CreateBroker(context.Background(), tt.request)

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

// TestGetBroker tests the function GetBroker
func TestGetBroker(t *testing.T) {
	// Prepare data
	service := &Service{}
	validRequest := &brokerpb.GetBrokerRequest{
		Id: uuid.New().String(),
	}

	// Test cases
	tests := []struct {
		name            string
		mockSetup       func(ctrl *gomock.Controller)
		request         *brokerpb.GetBrokerRequest
		expected        *brokerpb.GetBrokerResponse
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
			expected:        &brokerpb.GetBrokerResponse{},
			expectedErrCode: codes.InvalidArgument,
		},
		{
			name: "fails to retrieve the broker",
			mockSetup: func(ctrl *gomock.Controller) {
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().Get(gomock.Any()).Return(models.Broker{}, false, errors.New("error"))
				repositories.ReplaceGlobals(repositories.NewRepository(bb, nil, nil))
			},
			request:         validRequest,
			expected:        &brokerpb.GetBrokerResponse{},
			expectedErrCode: codes.Internal,
		},
		{
			name: "fails to find the broker",
			mockSetup: func(ctrl *gomock.Controller) {
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().Get(gomock.Any()).Return(models.Broker{}, false, nil)
				repositories.ReplaceGlobals(repositories.NewRepository(bb, nil, nil))
			},
			request:         validRequest,
			expected:        &brokerpb.GetBrokerResponse{},
			expectedErrCode: codes.NotFound,
		},
		{
			name: "succeeded",
			mockSetup: func(ctrl *gomock.Controller) {
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().Get(gomock.Any()).Return(models.Broker{}, true, nil)
				repositories.ReplaceGlobals(repositories.NewRepository(bb, nil, nil))
			},
			request:         validRequest,
			expected:        &brokerpb.GetBrokerResponse{},
			expectedErrCode: codes.OK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Apply mocks
			ctrl := gomock.NewController(t)
			tt.mockSetup(ctrl)
			defer ctrl.Finish()

			// Call service
			response, err := service.GetBroker(context.Background(), tt.request)

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

// TestUpdateBroker tests the function UpdateBroker
func TestUpdateBroker(t *testing.T) {
	// Prepare data
	service := &Service{}
	validRequest := &brokerpb.UpdateBrokerRequest{
		Id:       uuid.New().String(),
		Name:     "name",
		Disabled: false,
	}
	validBroker := models.Broker{
		ID:       uuid.New(),
		Name:     "name",
		Disabled: false,
	}

	// Test cases
	tests := []struct {
		name            string
		mockSetup       func(ctrl *gomock.Controller)
		request         *brokerpb.UpdateBrokerRequest
		expected        *brokerpb.UpdateBrokerResponse
		expectedErrCode codes.Code
	}{
		{
			name: "does not have permission",
			mockSetup: func(ctrl *gomock.Controller) {
				// Mock the public security facade
				publicSecurityClient := mocks.NewMockPublicSecurityServiceClient(ctrl)
				publicSecurityClient.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(&securitypb.CheckPermissionResponse{HasPermission: false}, nil)
				security.ReplaceGlobals(security.NewPublicSecurityFacadeWithGrpcClient(publicSecurityClient))
			},
			request:         validRequest,
			expected:        &brokerpb.UpdateBrokerResponse{},
			expectedErrCode: codes.PermissionDenied,
		},
		{
			name: "fails to parse ID from request",
			mockSetup: func(ctrl *gomock.Controller) {
				// Mock the public security facade
				publicSecurityClient := mocks.NewMockPublicSecurityServiceClient(ctrl)
				publicSecurityClient.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(&securitypb.CheckPermissionResponse{HasPermission: true}, nil)
				security.ReplaceGlobals(security.NewPublicSecurityFacadeWithGrpcClient(publicSecurityClient))
				// Mock the broker repository
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().Get(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(bb, nil, nil))
			},
			request: &brokerpb.UpdateBrokerRequest{
				Id: "bad-uuid",
			},
			expected:        &brokerpb.UpdateBrokerResponse{},
			expectedErrCode: codes.InvalidArgument,
		},
		{
			name: "fails at bad broker input",
			mockSetup: func(ctrl *gomock.Controller) {
				// Mock the public security facade
				publicSecurityClient := mocks.NewMockPublicSecurityServiceClient(ctrl)
				publicSecurityClient.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(&securitypb.CheckPermissionResponse{HasPermission: true}, nil)
				security.ReplaceGlobals(security.NewPublicSecurityFacadeWithGrpcClient(publicSecurityClient))
				// Mock the broker repository
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().Get(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(bb, nil, nil))
			},
			request: &brokerpb.UpdateBrokerRequest{
				Id: uuid.Nil.String(),
			},
			expected:        &brokerpb.UpdateBrokerResponse{},
			expectedErrCode: codes.InvalidArgument,
		},
		{
			name: "fails to retrieve broker",
			mockSetup: func(ctrl *gomock.Controller) {
				// Mock the public security facade
				publicSecurityClient := mocks.NewMockPublicSecurityServiceClient(ctrl)
				publicSecurityClient.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(&securitypb.CheckPermissionResponse{HasPermission: true}, nil)
				security.ReplaceGlobals(security.NewPublicSecurityFacadeWithGrpcClient(publicSecurityClient))
				// Mock the broker repository
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().Get(gomock.Any()).Return(models.Broker{}, false, errors.New("error"))
				bb.EXPECT().ExistsByName(gomock.Any()).Times(0)
				bb.EXPECT().Update(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(bb, nil, nil))
			},
			request:         validRequest,
			expected:        &brokerpb.UpdateBrokerResponse{},
			expectedErrCode: codes.Internal,
		},
		{
			name: "fails to find broker",
			mockSetup: func(ctrl *gomock.Controller) {
				// Mock the public security facade
				publicSecurityClient := mocks.NewMockPublicSecurityServiceClient(ctrl)
				publicSecurityClient.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(&securitypb.CheckPermissionResponse{HasPermission: true}, nil)
				security.ReplaceGlobals(security.NewPublicSecurityFacadeWithGrpcClient(publicSecurityClient))
				// Mock the broker repository
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().Get(gomock.Any()).Return(models.Broker{}, false, nil)
				bb.EXPECT().ExistsByName(gomock.Any()).Times(0)
				bb.EXPECT().Update(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(bb, nil, nil))
			},
			request:         validRequest,
			expected:        &brokerpb.UpdateBrokerResponse{},
			expectedErrCode: codes.NotFound,
		},
		{
			name: "fails to verify new name usage",
			mockSetup: func(ctrl *gomock.Controller) {
				// Mock the public security facade
				publicSecurityClient := mocks.NewMockPublicSecurityServiceClient(ctrl)
				publicSecurityClient.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(&securitypb.CheckPermissionResponse{HasPermission: true}, nil)
				security.ReplaceGlobals(security.NewPublicSecurityFacadeWithGrpcClient(publicSecurityClient))
				// Mock the broker repository
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().Get(gomock.Any()).Return(models.Broker{}, true, nil)
				bb.EXPECT().ExistsByName(gomock.Any()).Return(true, errors.New("error"))
				bb.EXPECT().Update(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(bb, nil, nil))
			},
			request:         validRequest,
			expected:        &brokerpb.UpdateBrokerResponse{},
			expectedErrCode: codes.Internal,
		},
		{
			name: "new name already in use",
			mockSetup: func(ctrl *gomock.Controller) {
				// Mock the public security facade
				publicSecurityClient := mocks.NewMockPublicSecurityServiceClient(ctrl)
				publicSecurityClient.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(&securitypb.CheckPermissionResponse{HasPermission: true}, nil)
				security.ReplaceGlobals(security.NewPublicSecurityFacadeWithGrpcClient(publicSecurityClient))
				// Mock the broker repository
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().Get(gomock.Any()).Return(models.Broker{}, true, nil)
				bb.EXPECT().ExistsByName(gomock.Any()).Return(true, nil)
				bb.EXPECT().Update(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(bb, nil, nil))
			},
			request:         validRequest,
			expected:        &brokerpb.UpdateBrokerResponse{},
			expectedErrCode: codes.AlreadyExists,
		},
		{
			name: "fails to update broker",
			mockSetup: func(ctrl *gomock.Controller) {
				// Mock the public security facade
				publicSecurityClient := mocks.NewMockPublicSecurityServiceClient(ctrl)
				publicSecurityClient.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(&securitypb.CheckPermissionResponse{HasPermission: true}, nil)
				security.ReplaceGlobals(security.NewPublicSecurityFacadeWithGrpcClient(publicSecurityClient))
				// Mock the broker repository
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().Get(gomock.Any()).Return(models.Broker{}, true, nil)
				bb.EXPECT().ExistsByName(gomock.Any()).Return(false, nil)
				bb.EXPECT().Update(gomock.Any()).Return(errors.New("error"))
				repositories.ReplaceGlobals(repositories.NewRepository(bb, nil, nil))
			},
			request:         validRequest,
			expected:        &brokerpb.UpdateBrokerResponse{},
			expectedErrCode: codes.Internal,
		},
		{
			name: "fails to retrieve the broker",
			mockSetup: func(ctrl *gomock.Controller) {
				// Mock the public security facade
				publicSecurityClient := mocks.NewMockPublicSecurityServiceClient(ctrl)
				publicSecurityClient.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(&securitypb.CheckPermissionResponse{HasPermission: true}, nil)
				security.ReplaceGlobals(security.NewPublicSecurityFacadeWithGrpcClient(publicSecurityClient))
				// Mock the broker repository
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().Get(gomock.Any()).Return(models.Broker{}, true, nil)
				bb.EXPECT().ExistsByName(gomock.Any()).Return(false, nil)
				bb.EXPECT().Update(gomock.Any()).Return(nil)
				bb.EXPECT().Get(gomock.Any()).Return(models.Broker{}, false, errors.New("error"))
				repositories.ReplaceGlobals(repositories.NewRepository(bb, nil, nil))
			},
			request:         validRequest,
			expected:        &brokerpb.UpdateBrokerResponse{},
			expectedErrCode: codes.Internal,
		},
		{
			name: "fails to find broker",
			mockSetup: func(ctrl *gomock.Controller) {
				// Mock the public security facade
				publicSecurityClient := mocks.NewMockPublicSecurityServiceClient(ctrl)
				publicSecurityClient.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(&securitypb.CheckPermissionResponse{HasPermission: true}, nil)
				security.ReplaceGlobals(security.NewPublicSecurityFacadeWithGrpcClient(publicSecurityClient))
				// Mock the broker repository
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().Get(gomock.Any()).Return(models.Broker{}, true, nil)
				bb.EXPECT().ExistsByName(gomock.Any()).Return(false, nil)
				bb.EXPECT().Update(gomock.Any()).Return(nil)
				bb.EXPECT().Get(gomock.Any()).Return(models.Broker{}, false, nil)
				repositories.ReplaceGlobals(repositories.NewRepository(bb, nil, nil))
			},
			request:         validRequest,
			expected:        &brokerpb.UpdateBrokerResponse{},
			expectedErrCode: codes.Internal,
		},
		{
			name: "succeeded",
			mockSetup: func(ctrl *gomock.Controller) {
				// Mock the public security facade
				publicSecurityClient := mocks.NewMockPublicSecurityServiceClient(ctrl)
				publicSecurityClient.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(&securitypb.CheckPermissionResponse{HasPermission: true}, nil)
				security.ReplaceGlobals(security.NewPublicSecurityFacadeWithGrpcClient(publicSecurityClient))
				// Mock the broker repository
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().Get(gomock.Any()).Return(models.Broker{}, true, nil)
				bb.EXPECT().ExistsByName(gomock.Any()).Return(false, nil)
				bb.EXPECT().Update(gomock.Any()).Return(nil)
				bb.EXPECT().Get(gomock.Any()).Return(validBroker, true, nil)
				repositories.ReplaceGlobals(repositories.NewRepository(bb, nil, nil))
			},
			request:         validRequest,
			expected:        &brokerpb.UpdateBrokerResponse{},
			expectedErrCode: codes.OK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Apply mocks
			ctrl := gomock.NewController(t)
			tt.mockSetup(ctrl)
			defer ctrl.Finish()

			// Call service
			response, err := service.UpdateBroker(context.Background(), tt.request)

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

// TestDeleteBroker tests the function DeleteBroker
func TestDeleteBroker(t *testing.T) {
	// Prepare data
	service := &Service{}
	validRequest := &brokerpb.DeleteBrokerRequest{
		Id: uuid.New().String(),
	}

	// Test cases
	tests := []struct {
		name            string
		mockSetup       func(ctrl *gomock.Controller)
		request         *brokerpb.DeleteBrokerRequest
		expected        *brokerpb.DeleteBrokerResponse
		expectedErrCode codes.Code
	}{
		{
			name: "does not have permission",
			mockSetup: func(ctrl *gomock.Controller) {
				// Mock the public security facade
				publicSecurityClient := mocks.NewMockPublicSecurityServiceClient(ctrl)
				publicSecurityClient.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(&securitypb.CheckPermissionResponse{HasPermission: false}, nil)
				security.ReplaceGlobals(security.NewPublicSecurityFacadeWithGrpcClient(publicSecurityClient))
			},
			request:         validRequest,
			expected:        &brokerpb.DeleteBrokerResponse{},
			expectedErrCode: codes.PermissionDenied,
		},
		{
			name: "fails to parse ID from request",
			mockSetup: func(ctrl *gomock.Controller) {
				// Mock the public security facade
				publicSecurityClient := mocks.NewMockPublicSecurityServiceClient(ctrl)
				publicSecurityClient.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(&securitypb.CheckPermissionResponse{HasPermission: true}, nil)
				security.ReplaceGlobals(security.NewPublicSecurityFacadeWithGrpcClient(publicSecurityClient))
				// Mock the broker repository
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().Exists(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(bb, nil, nil))
			},
			request: &brokerpb.DeleteBrokerRequest{
				Id: "bad-uuid",
			},
			expected:        &brokerpb.DeleteBrokerResponse{},
			expectedErrCode: codes.InvalidArgument,
		},
		{
			name: "fails to retrieve the broker",
			mockSetup: func(ctrl *gomock.Controller) {
				// Mock the public security facade
				publicSecurityClient := mocks.NewMockPublicSecurityServiceClient(ctrl)
				publicSecurityClient.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(&securitypb.CheckPermissionResponse{HasPermission: true}, nil)
				security.ReplaceGlobals(security.NewPublicSecurityFacadeWithGrpcClient(publicSecurityClient))
				// Mock the broker repository
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().Exists(gomock.Any()).Return(false, errors.New("error"))
				bb.EXPECT().Delete(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(bb, nil, nil))
			},
			request:         validRequest,
			expected:        &brokerpb.DeleteBrokerResponse{},
			expectedErrCode: codes.Internal,
		},
		{
			name: "fails to find the broker",
			mockSetup: func(ctrl *gomock.Controller) {
				// Mock the public security facade
				publicSecurityClient := mocks.NewMockPublicSecurityServiceClient(ctrl)
				publicSecurityClient.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(&securitypb.CheckPermissionResponse{HasPermission: true}, nil)
				security.ReplaceGlobals(security.NewPublicSecurityFacadeWithGrpcClient(publicSecurityClient))
				// Mock the broker repository
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().Exists(gomock.Any()).Return(false, nil)
				bb.EXPECT().Delete(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(bb, nil, nil))
			},
			request:         validRequest,
			expected:        &brokerpb.DeleteBrokerResponse{},
			expectedErrCode: codes.NotFound,
		},
		{
			name: "fails to delete the broker",
			mockSetup: func(ctrl *gomock.Controller) {
				// Mock the public security facade
				publicSecurityClient := mocks.NewMockPublicSecurityServiceClient(ctrl)
				publicSecurityClient.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(&securitypb.CheckPermissionResponse{HasPermission: true}, nil)
				security.ReplaceGlobals(security.NewPublicSecurityFacadeWithGrpcClient(publicSecurityClient))
				// Mock the broker repository
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().Exists(gomock.Any()).Return(true, nil)
				bb.EXPECT().Delete(gomock.Any()).Return(errors.New("error"))
				repositories.ReplaceGlobals(repositories.NewRepository(bb, nil, nil))
			},
			request:         validRequest,
			expected:        &brokerpb.DeleteBrokerResponse{},
			expectedErrCode: codes.Internal,
		},
		{
			name: "succeeded",
			mockSetup: func(ctrl *gomock.Controller) {
				// Mock the public security facade
				publicSecurityClient := mocks.NewMockPublicSecurityServiceClient(ctrl)
				publicSecurityClient.EXPECT().CheckPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(&securitypb.CheckPermissionResponse{HasPermission: true}, nil)
				security.ReplaceGlobals(security.NewPublicSecurityFacadeWithGrpcClient(publicSecurityClient))
				// Mock the broker repository
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().Exists(gomock.Any()).Return(true, nil)
				bb.EXPECT().Delete(gomock.Any()).Return(nil)
				repositories.ReplaceGlobals(repositories.NewRepository(bb, nil, nil))
			},
			request:         validRequest,
			expected:        &brokerpb.DeleteBrokerResponse{},
			expectedErrCode: codes.OK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Apply mocks
			ctrl := gomock.NewController(t)
			tt.mockSetup(ctrl)
			defer ctrl.Finish()

			// Call service
			response, err := service.DeleteBroker(context.Background(), tt.request)

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

// TestListBrokers tests the function ListBrokers
func TestListBrokers(t *testing.T) {
	// Prepare data
	service := &Service{}

	// Test cases
	tests := []struct {
		name            string
		mockSetup       func(ctrl *gomock.Controller)
		request         *brokerpb.ListBrokersRequest
		expected        *brokerpb.ListBrokersResponse
		expectedErrCode codes.Code
	}{
		{
			name: "fails to retrieve all enabled brokers",
			mockSetup: func(ctrl *gomock.Controller) {
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().GetAllEnabled().Return(nil, errors.New("error"))
				repositories.ReplaceGlobals(repositories.NewRepository(bb, nil, nil))
			},
			request: &brokerpb.ListBrokersRequest{
				EnabledOnly: true,
			},
			expected:        &brokerpb.ListBrokersResponse{},
			expectedErrCode: codes.Internal,
		},
		{
			name: "fails to retrieve all brokers",
			mockSetup: func(ctrl *gomock.Controller) {
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().GetAll().Return(nil, errors.New("error"))
				repositories.ReplaceGlobals(repositories.NewRepository(bb, nil, nil))
			},
			request: &brokerpb.ListBrokersRequest{
				EnabledOnly: false,
			},
			expected:        &brokerpb.ListBrokersResponse{},
			expectedErrCode: codes.Internal,
		},
		{
			name: "succeeded",
			mockSetup: func(ctrl *gomock.Controller) {
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().GetAll().Return([]models.Broker{}, nil)
				repositories.ReplaceGlobals(repositories.NewRepository(bb, nil, nil))
			},
			request: &brokerpb.ListBrokersRequest{
				EnabledOnly: false,
			},
			expected:        &brokerpb.ListBrokersResponse{},
			expectedErrCode: codes.OK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Apply mocks
			ctrl := gomock.NewController(t)
			tt.mockSetup(ctrl)
			defer ctrl.Finish()

			// Call service
			response, err := service.ListBrokers(context.Background(), tt.request)

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
