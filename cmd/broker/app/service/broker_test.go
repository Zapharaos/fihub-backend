package service

import (
	"context"
	"errors"
	"github.com/Zapharaos/fihub-backend/cmd/broker/app/repositories"
	"github.com/Zapharaos/fihub-backend/gen"
	"github.com/Zapharaos/fihub-backend/internal/models"
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
	validRequest := &protogen.CreateBrokerRequest{
		Name:     "name",
		Disabled: false,
	}

	// Test cases
	tests := []struct {
		name            string
		mockSetup       func(ctrl *gomock.Controller)
		request         *protogen.CreateBrokerRequest
		expected        *protogen.CreateBrokerResponse
		expectedErrCode codes.Code
	}{
		{
			name: "missing request body",
			mockSetup: func(ctrl *gomock.Controller) {
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().ExistsByName(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(bb, nil, nil))
			},
			request:         nil,
			expected:        &protogen.CreateBrokerResponse{},
			expectedErrCode: codes.InvalidArgument,
		},
		{
			name: "fails at bad broker input",
			mockSetup: func(ctrl *gomock.Controller) {
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().ExistsByName(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(bb, nil, nil))
			},
			request: &protogen.CreateBrokerRequest{
				Name:     "",
				Disabled: false,
			},
			expected:        &protogen.CreateBrokerResponse{},
			expectedErrCode: codes.InvalidArgument,
		},
		{
			name: "fails to verify broker existence",
			mockSetup: func(ctrl *gomock.Controller) {
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().ExistsByName(gomock.Any()).Return(false, errors.New("error"))
				bb.EXPECT().Create(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(bb, nil, nil))
			},
			request:         validRequest,
			expected:        &protogen.CreateBrokerResponse{},
			expectedErrCode: codes.Internal,
		},
		{
			name: "broker already exists",
			mockSetup: func(ctrl *gomock.Controller) {
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().ExistsByName(gomock.Any()).Return(true, nil)
				bb.EXPECT().Create(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(bb, nil, nil))
			},
			request:         validRequest,
			expected:        &protogen.CreateBrokerResponse{},
			expectedErrCode: codes.AlreadyExists,
		},
		{
			name: "fails to create broker",
			mockSetup: func(ctrl *gomock.Controller) {
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().ExistsByName(gomock.Any()).Return(false, nil)
				bb.EXPECT().Create(gomock.Any()).Return(uuid.Nil, errors.New("error"))
				bb.EXPECT().Get(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(bb, nil, nil))
			},
			request:         validRequest,
			expected:        &protogen.CreateBrokerResponse{},
			expectedErrCode: codes.Internal,
		},
		{
			name: "fails to retrieve the broker",
			mockSetup: func(ctrl *gomock.Controller) {
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().ExistsByName(gomock.Any()).Return(false, nil)
				bb.EXPECT().Create(gomock.Any()).Return(uuid.Nil, nil)
				bb.EXPECT().Get(gomock.Any()).Return(models.Broker{}, false, errors.New("error"))
				repositories.ReplaceGlobals(repositories.NewRepository(bb, nil, nil))
			},
			request:         validRequest,
			expected:        &protogen.CreateBrokerResponse{},
			expectedErrCode: codes.Internal,
		},
		{
			name: "fails to find broker",
			mockSetup: func(ctrl *gomock.Controller) {
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().ExistsByName(gomock.Any()).Return(false, nil)
				bb.EXPECT().Create(gomock.Any()).Return(uuid.Nil, nil)
				bb.EXPECT().Get(gomock.Any()).Return(models.Broker{}, false, nil)
				repositories.ReplaceGlobals(repositories.NewRepository(bb, nil, nil))
			},
			request:         validRequest,
			expected:        &protogen.CreateBrokerResponse{},
			expectedErrCode: codes.Internal,
		},
		{
			name: "succeeded",
			mockSetup: func(ctrl *gomock.Controller) {
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().ExistsByName(gomock.Any()).Return(false, nil)
				bb.EXPECT().Create(gomock.Any()).Return(uuid.Nil, nil)
				bb.EXPECT().Get(gomock.Any()).Return(models.Broker{}, true, nil)
				repositories.ReplaceGlobals(repositories.NewRepository(bb, nil, nil))
			},
			request:         validRequest,
			expected:        &protogen.CreateBrokerResponse{},
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
	validRequest := &protogen.GetBrokerRequest{
		Id: uuid.New().String(),
	}

	// Test cases
	tests := []struct {
		name            string
		mockSetup       func(ctrl *gomock.Controller)
		request         *protogen.GetBrokerRequest
		expected        *protogen.GetBrokerResponse
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
			expected:        &protogen.GetBrokerResponse{},
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
			expected:        &protogen.GetBrokerResponse{},
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
			expected:        &protogen.GetBrokerResponse{},
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
			expected:        &protogen.GetBrokerResponse{},
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
	validRequest := &protogen.UpdateBrokerRequest{
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
		request         *protogen.UpdateBrokerRequest
		expected        *protogen.UpdateBrokerResponse
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
			expected:        &protogen.UpdateBrokerResponse{},
			expectedErrCode: codes.InvalidArgument,
		},
		{
			name: "fails to parse ID from request",
			mockSetup: func(ctrl *gomock.Controller) {
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().Get(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(bb, nil, nil))
			},
			request: &protogen.UpdateBrokerRequest{
				Id: "bad-uuid",
			},
			expected:        &protogen.UpdateBrokerResponse{},
			expectedErrCode: codes.InvalidArgument,
		},
		{
			name: "fails at bad broker input",
			mockSetup: func(ctrl *gomock.Controller) {
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().Get(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(bb, nil, nil))
			},
			request: &protogen.UpdateBrokerRequest{
				Id: uuid.Nil.String(),
			},
			expected:        &protogen.UpdateBrokerResponse{},
			expectedErrCode: codes.InvalidArgument,
		},
		{
			name: "fails to retrieve broker",
			mockSetup: func(ctrl *gomock.Controller) {
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().Get(gomock.Any()).Return(models.Broker{}, false, errors.New("error"))
				bb.EXPECT().ExistsByName(gomock.Any()).Times(0)
				bb.EXPECT().Update(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(bb, nil, nil))
			},
			request:         validRequest,
			expected:        &protogen.UpdateBrokerResponse{},
			expectedErrCode: codes.Internal,
		},
		{
			name: "fails to find broker",
			mockSetup: func(ctrl *gomock.Controller) {
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().Get(gomock.Any()).Return(models.Broker{}, false, nil)
				bb.EXPECT().ExistsByName(gomock.Any()).Times(0)
				bb.EXPECT().Update(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(bb, nil, nil))
			},
			request:         validRequest,
			expected:        &protogen.UpdateBrokerResponse{},
			expectedErrCode: codes.NotFound,
		},
		{
			name: "fails to verify new name usage",
			mockSetup: func(ctrl *gomock.Controller) {
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().Get(gomock.Any()).Return(models.Broker{}, true, nil)
				bb.EXPECT().ExistsByName(gomock.Any()).Return(true, errors.New("error"))
				bb.EXPECT().Update(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(bb, nil, nil))
			},
			request:         validRequest,
			expected:        &protogen.UpdateBrokerResponse{},
			expectedErrCode: codes.Internal,
		},
		{
			name: "new name already in use",
			mockSetup: func(ctrl *gomock.Controller) {
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().Get(gomock.Any()).Return(models.Broker{}, true, nil)
				bb.EXPECT().ExistsByName(gomock.Any()).Return(true, nil)
				bb.EXPECT().Update(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(bb, nil, nil))
			},
			request:         validRequest,
			expected:        &protogen.UpdateBrokerResponse{},
			expectedErrCode: codes.AlreadyExists,
		},
		{
			name: "fails to update broker",
			mockSetup: func(ctrl *gomock.Controller) {
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().Get(gomock.Any()).Return(models.Broker{}, true, nil)
				bb.EXPECT().ExistsByName(gomock.Any()).Return(false, nil)
				bb.EXPECT().Update(gomock.Any()).Return(errors.New("error"))
				repositories.ReplaceGlobals(repositories.NewRepository(bb, nil, nil))
			},
			request:         validRequest,
			expected:        &protogen.UpdateBrokerResponse{},
			expectedErrCode: codes.Internal,
		},
		{
			name: "fails to retrieve the broker",
			mockSetup: func(ctrl *gomock.Controller) {
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().Get(gomock.Any()).Return(models.Broker{}, true, nil)
				bb.EXPECT().ExistsByName(gomock.Any()).Return(false, nil)
				bb.EXPECT().Update(gomock.Any()).Return(nil)
				bb.EXPECT().Get(gomock.Any()).Return(models.Broker{}, false, errors.New("error"))
				repositories.ReplaceGlobals(repositories.NewRepository(bb, nil, nil))
			},
			request:         validRequest,
			expected:        &protogen.UpdateBrokerResponse{},
			expectedErrCode: codes.Internal,
		},
		{
			name: "fails to find broker",
			mockSetup: func(ctrl *gomock.Controller) {
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().Get(gomock.Any()).Return(models.Broker{}, true, nil)
				bb.EXPECT().ExistsByName(gomock.Any()).Return(false, nil)
				bb.EXPECT().Update(gomock.Any()).Return(nil)
				bb.EXPECT().Get(gomock.Any()).Return(models.Broker{}, false, nil)
				repositories.ReplaceGlobals(repositories.NewRepository(bb, nil, nil))
			},
			request:         validRequest,
			expected:        &protogen.UpdateBrokerResponse{},
			expectedErrCode: codes.Internal,
		},
		{
			name: "succeeded",
			mockSetup: func(ctrl *gomock.Controller) {
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().Get(gomock.Any()).Return(models.Broker{}, true, nil)
				bb.EXPECT().ExistsByName(gomock.Any()).Return(false, nil)
				bb.EXPECT().Update(gomock.Any()).Return(nil)
				bb.EXPECT().Get(gomock.Any()).Return(validBroker, true, nil)
				repositories.ReplaceGlobals(repositories.NewRepository(bb, nil, nil))
			},
			request:         validRequest,
			expected:        &protogen.UpdateBrokerResponse{},
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
	validRequest := &protogen.DeleteBrokerRequest{
		Id: uuid.New().String(),
	}

	// Test cases
	tests := []struct {
		name            string
		mockSetup       func(ctrl *gomock.Controller)
		request         *protogen.DeleteBrokerRequest
		expected        *protogen.DeleteBrokerResponse
		expectedErrCode codes.Code
	}{
		{
			name: "missing request body",
			mockSetup: func(ctrl *gomock.Controller) {
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().Exists(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(bb, nil, nil))
			},
			request:         nil,
			expected:        &protogen.DeleteBrokerResponse{},
			expectedErrCode: codes.InvalidArgument,
		},
		{
			name: "fails to parse ID from request",
			mockSetup: func(ctrl *gomock.Controller) {
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().Exists(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(bb, nil, nil))
			},
			request: &protogen.DeleteBrokerRequest{
				Id: "bad-uuid",
			},
			expected:        &protogen.DeleteBrokerResponse{},
			expectedErrCode: codes.InvalidArgument,
		},
		{
			name: "fails to retrieve the broker",
			mockSetup: func(ctrl *gomock.Controller) {
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().Exists(gomock.Any()).Return(false, errors.New("error"))
				bb.EXPECT().Delete(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(bb, nil, nil))
			},
			request:         validRequest,
			expected:        &protogen.DeleteBrokerResponse{},
			expectedErrCode: codes.Internal,
		},
		{
			name: "fails to find the broker",
			mockSetup: func(ctrl *gomock.Controller) {
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().Exists(gomock.Any()).Return(false, nil)
				bb.EXPECT().Delete(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(bb, nil, nil))
			},
			request:         validRequest,
			expected:        &protogen.DeleteBrokerResponse{},
			expectedErrCode: codes.NotFound,
		},
		{
			name: "fails to delete the broker",
			mockSetup: func(ctrl *gomock.Controller) {
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().Exists(gomock.Any()).Return(true, nil)
				bb.EXPECT().Delete(gomock.Any()).Return(errors.New("error"))
				repositories.ReplaceGlobals(repositories.NewRepository(bb, nil, nil))
			},
			request:         validRequest,
			expected:        &protogen.DeleteBrokerResponse{},
			expectedErrCode: codes.Internal,
		},
		{
			name: "succeeded",
			mockSetup: func(ctrl *gomock.Controller) {
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().Exists(gomock.Any()).Return(true, nil)
				bb.EXPECT().Delete(gomock.Any()).Return(nil)
				repositories.ReplaceGlobals(repositories.NewRepository(bb, nil, nil))
			},
			request:         validRequest,
			expected:        &protogen.DeleteBrokerResponse{},
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
		request         *protogen.ListBrokersRequest
		expected        *protogen.ListBrokersResponse
		expectedErrCode codes.Code
	}{
		{
			name: "fails to retrieve all enabled brokers",
			mockSetup: func(ctrl *gomock.Controller) {
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().GetAllEnabled().Return(nil, errors.New("error"))
				repositories.ReplaceGlobals(repositories.NewRepository(bb, nil, nil))
			},
			request: &protogen.ListBrokersRequest{
				EnabledOnly: true,
			},
			expected:        &protogen.ListBrokersResponse{},
			expectedErrCode: codes.Internal,
		},
		{
			name: "fails to retrieve all brokers",
			mockSetup: func(ctrl *gomock.Controller) {
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().GetAll().Return(nil, errors.New("error"))
				repositories.ReplaceGlobals(repositories.NewRepository(bb, nil, nil))
			},
			request: &protogen.ListBrokersRequest{
				EnabledOnly: false,
			},
			expected:        &protogen.ListBrokersResponse{},
			expectedErrCode: codes.Internal,
		},
		{
			name: "succeeded",
			mockSetup: func(ctrl *gomock.Controller) {
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().GetAll().Return([]models.Broker{}, nil)
				repositories.ReplaceGlobals(repositories.NewRepository(bb, nil, nil))
			},
			request: &protogen.ListBrokersRequest{
				EnabledOnly: false,
			},
			expected:        &protogen.ListBrokersResponse{},
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
