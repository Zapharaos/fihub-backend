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
	"strings"
	"testing"
)

// TestCreateBrokerImage tests the CreateBrokerImage function
func TestCreateBrokerImage(t *testing.T) {
	// Prepare data
	service := &Service{}
	fileData := []byte{0x00, 0x01, 0x02, 0x03}
	fileName := strings.Repeat("a", models.ImageNameMinLength)
	validRequest := &protogen.CreateBrokerImageRequest{
		BrokerId: uuid.New().String(),
		Name:     fileName,
		Data:     fileData,
	}
	validImage := models.BrokerImage{
		ID:       uuid.New(),
		BrokerID: uuid.New(),
		Name:     fileName,
		Data:     fileData,
	}

	// Prepare tests
	tests := []struct {
		name            string
		mockSetup       func(ctrl *gomock.Controller)
		request         *protogen.CreateBrokerImageRequest
		expected        *protogen.CreateBrokerImageResponse
		expectedErrCode codes.Code
	}{
		{
			name: "missing request body",
			mockSetup: func(ctrl *gomock.Controller) {
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().HasImage(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(bb, nil, nil))
			},
			request:         nil,
			expected:        &protogen.CreateBrokerImageResponse{},
			expectedErrCode: codes.InvalidArgument,
		},
		{
			name: "fails to parse ID from request",
			mockSetup: func(ctrl *gomock.Controller) {
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().HasImage(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(bb, nil, nil))
			},
			request: &protogen.CreateBrokerImageRequest{
				BrokerId: "bad-uuid",
			},
			expected:        &protogen.CreateBrokerImageResponse{},
			expectedErrCode: codes.InvalidArgument,
		},
		{
			name: "fails at bad image input",
			mockSetup: func(ctrl *gomock.Controller) {
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().HasImage(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(bb, nil, nil))
			},
			request: &protogen.CreateBrokerImageRequest{
				BrokerId: uuid.Nil.String(),
			},
			expected:        &protogen.CreateBrokerImageResponse{},
			expectedErrCode: codes.InvalidArgument,
		},
		{
			name: "fails to verify broker image existence",
			mockSetup: func(ctrl *gomock.Controller) {
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().HasImage(gomock.Any()).Return(true, errors.New("error"))
				bi := mocks.NewBrokerImageRepository(ctrl)
				bi.EXPECT().Create(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(bb, nil, nil))
			},
			request:         validRequest,
			expected:        &protogen.CreateBrokerImageResponse{},
			expectedErrCode: codes.Internal,
		},
		{
			name: "broker already has an image",
			mockSetup: func(ctrl *gomock.Controller) {
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().HasImage(gomock.Any()).Return(true, nil)
				bi := mocks.NewBrokerImageRepository(ctrl)
				bi.EXPECT().Create(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(bb, nil, nil))
			},
			request:         validRequest,
			expected:        &protogen.CreateBrokerImageResponse{},
			expectedErrCode: codes.InvalidArgument,
		},
		{
			name: "fails to create an image",
			mockSetup: func(ctrl *gomock.Controller) {
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().HasImage(gomock.Any()).Return(false, nil)
				bi := mocks.NewBrokerImageRepository(ctrl)
				bi.EXPECT().Create(gomock.Any()).Return(errors.New("error"))
				bi.EXPECT().Get(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(bb, nil, bi))
			},
			request:         validRequest,
			expected:        &protogen.CreateBrokerImageResponse{},
			expectedErrCode: codes.Internal,
		},
		{
			name: "fails to retrieve the broker image",
			mockSetup: func(ctrl *gomock.Controller) {
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().HasImage(gomock.Any()).Return(false, nil)
				bb.EXPECT().SetImage(gomock.Any(), gomock.Any()).Times(0)
				bi := mocks.NewBrokerImageRepository(ctrl)
				bi.EXPECT().Create(gomock.Any()).Return(nil)
				bi.EXPECT().Get(gomock.Any()).Return(models.BrokerImage{}, false, errors.New("error"))
				repositories.ReplaceGlobals(repositories.NewRepository(bb, nil, bi))
			},
			request:         validRequest,
			expected:        &protogen.CreateBrokerImageResponse{},
			expectedErrCode: codes.Internal,
		},
		{
			name: "fails to find the broker image",
			mockSetup: func(ctrl *gomock.Controller) {
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().HasImage(gomock.Any()).Return(false, nil)
				bb.EXPECT().SetImage(gomock.Any(), gomock.Any()).Times(0)
				bi := mocks.NewBrokerImageRepository(ctrl)
				bi.EXPECT().Create(gomock.Any()).Return(nil)
				bi.EXPECT().Get(gomock.Any()).Return(models.BrokerImage{}, false, nil)
				repositories.ReplaceGlobals(repositories.NewRepository(bb, nil, bi))
			},
			request:         validRequest,
			expected:        &protogen.CreateBrokerImageResponse{},
			expectedErrCode: codes.Internal,
		},
		{
			name: "fails to set the image to broker",
			mockSetup: func(ctrl *gomock.Controller) {
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().HasImage(gomock.Any()).Return(false, nil)
				bb.EXPECT().SetImage(gomock.Any(), gomock.Any()).Return(errors.New("error"))
				bi := mocks.NewBrokerImageRepository(ctrl)
				bi.EXPECT().Create(gomock.Any()).Return(nil)
				bi.EXPECT().Get(gomock.Any()).Return(models.BrokerImage{}, true, nil)
				repositories.ReplaceGlobals(repositories.NewRepository(bb, nil, bi))
			},
			request:         validRequest,
			expected:        &protogen.CreateBrokerImageResponse{},
			expectedErrCode: codes.Internal,
		},
		{
			name: "succeeded",
			mockSetup: func(ctrl *gomock.Controller) {
				bb := mocks.NewBrokerRepository(ctrl)
				bb.EXPECT().HasImage(gomock.Any()).Return(false, nil)
				bb.EXPECT().SetImage(gomock.Any(), gomock.Any()).Return(nil)
				bi := mocks.NewBrokerImageRepository(ctrl)
				bi.EXPECT().Create(gomock.Any()).Return(nil)
				bi.EXPECT().Get(gomock.Any()).Return(validImage, true, nil)
				repositories.ReplaceGlobals(repositories.NewRepository(bb, nil, bi))
			},
			request:         validRequest,
			expected:        &protogen.CreateBrokerImageResponse{},
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
			response, err := service.CreateBrokerImage(context.Background(), tt.request)

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

// TestGetBrokerImage tests the GetBrokerImage function
func TestGetBrokerImage(t *testing.T) {
	// Prepare data
	service := &Service{}
	validRequest := &protogen.GetBrokerImageRequest{
		ImageId: uuid.New().String(),
	}

	// Prepare tests
	tests := []struct {
		name            string
		mockSetup       func(ctrl *gomock.Controller)
		request         *protogen.GetBrokerImageRequest
		expected        *protogen.GetBrokerImageResponse
		expectedErrCode codes.Code
	}{
		{
			name: "missing request body",
			mockSetup: func(ctrl *gomock.Controller) {
				bi := mocks.NewBrokerImageRepository(ctrl)
				bi.EXPECT().Get(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(nil, nil, bi))
			},
			request:         nil,
			expected:        &protogen.GetBrokerImageResponse{},
			expectedErrCode: codes.InvalidArgument,
		},
		{
			name: "fails to parse ID from request",
			mockSetup: func(ctrl *gomock.Controller) {
				bi := mocks.NewBrokerImageRepository(ctrl)
				bi.EXPECT().Get(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(nil, nil, bi))
			},
			request: &protogen.GetBrokerImageRequest{
				ImageId: "bad-uuid",
			},
			expected:        &protogen.GetBrokerImageResponse{},
			expectedErrCode: codes.InvalidArgument,
		},
		{
			name: "fails to retrieve broker image",
			mockSetup: func(ctrl *gomock.Controller) {
				bi := mocks.NewBrokerImageRepository(ctrl)
				bi.EXPECT().Get(gomock.Any()).Return(models.BrokerImage{}, false, errors.New("error"))
				repositories.ReplaceGlobals(repositories.NewRepository(nil, nil, bi))
			},
			request:         validRequest,
			expected:        &protogen.GetBrokerImageResponse{},
			expectedErrCode: codes.Internal,
		},
		{
			name: "fails to find broker image",
			mockSetup: func(ctrl *gomock.Controller) {
				bi := mocks.NewBrokerImageRepository(ctrl)
				bi.EXPECT().Get(gomock.Any()).Return(models.BrokerImage{}, false, nil)
				repositories.ReplaceGlobals(repositories.NewRepository(nil, nil, bi))
			},
			request:         validRequest,
			expected:        &protogen.GetBrokerImageResponse{},
			expectedErrCode: codes.NotFound,
		},
		{
			name: "succeeded",
			mockSetup: func(ctrl *gomock.Controller) {
				bi := mocks.NewBrokerImageRepository(ctrl)
				bi.EXPECT().Get(gomock.Any()).Return(models.BrokerImage{Data: []byte{0x00, 0x01}}, true, nil)
				repositories.ReplaceGlobals(repositories.NewRepository(nil, nil, bi))
			},
			request:         validRequest,
			expected:        &protogen.GetBrokerImageResponse{},
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
			response, err := service.GetBrokerImage(context.Background(), tt.request)

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

// TestUpdateBrokerImage tests the UpdateBrokerImage function
func TestUpdateBrokerImage(t *testing.T) {
	// Prepare data
	service := &Service{}
	fileData := []byte{0x00, 0x01, 0x02, 0x03}
	fileName := strings.Repeat("a", models.ImageNameMinLength)
	validRequest := &protogen.UpdateBrokerImageRequest{
		ImageId:  uuid.New().String(),
		BrokerId: uuid.New().String(),
		Name:     fileName,
		Data:     fileData,
	}
	validImage := models.BrokerImage{
		ID:       uuid.New(),
		BrokerID: uuid.New(),
		Name:     fileName,
		Data:     fileData,
	}

	// Prepare tests
	tests := []struct {
		name            string
		mockSetup       func(ctrl *gomock.Controller)
		request         *protogen.UpdateBrokerImageRequest
		expected        *protogen.UpdateBrokerImageResponse
		expectedErrCode codes.Code
	}{
		{
			name: "missing request body",
			mockSetup: func(ctrl *gomock.Controller) {
				bi := mocks.NewBrokerImageRepository(ctrl)
				bi.EXPECT().Exists(gomock.Any(), gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(nil, nil, bi))
			},
			request:         nil,
			expected:        &protogen.UpdateBrokerImageResponse{},
			expectedErrCode: codes.InvalidArgument,
		},
		{
			name: "fails to parse ID from request",
			mockSetup: func(ctrl *gomock.Controller) {
				bi := mocks.NewBrokerImageRepository(ctrl)
				bi.EXPECT().Exists(gomock.Any(), gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(nil, nil, bi))
			},
			request: &protogen.UpdateBrokerImageRequest{
				ImageId: "bad-uuid",
			},
			expected:        &protogen.UpdateBrokerImageResponse{},
			expectedErrCode: codes.InvalidArgument,
		},
		{
			name: "image input is invalid",
			mockSetup: func(ctrl *gomock.Controller) {
				bi := mocks.NewBrokerImageRepository(ctrl)
				bi.EXPECT().Exists(gomock.Any(), gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(nil, nil, bi))
			},
			request: &protogen.UpdateBrokerImageRequest{
				ImageId: uuid.Nil.String(),
			},
			expected:        &protogen.UpdateBrokerImageResponse{},
			expectedErrCode: codes.InvalidArgument,
		},
		{
			name: "fails to verify image broker existence",
			mockSetup: func(ctrl *gomock.Controller) {
				bi := mocks.NewBrokerImageRepository(ctrl)
				bi.EXPECT().Exists(gomock.Any(), gomock.Any()).Return(false, errors.New("error"))
				bi.EXPECT().Update(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(nil, nil, bi))
			},
			request:         validRequest,
			expected:        &protogen.UpdateBrokerImageResponse{},
			expectedErrCode: codes.Internal,
		},
		{
			name: "fails to find image broker",
			mockSetup: func(ctrl *gomock.Controller) {
				bi := mocks.NewBrokerImageRepository(ctrl)
				bi.EXPECT().Exists(gomock.Any(), gomock.Any()).Return(false, nil)
				bi.EXPECT().Update(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(nil, nil, bi))
			},
			request:         validRequest,
			expected:        &protogen.UpdateBrokerImageResponse{},
			expectedErrCode: codes.NotFound,
		},
		{
			name: "fails to update image",
			mockSetup: func(ctrl *gomock.Controller) {
				bi := mocks.NewBrokerImageRepository(ctrl)
				bi.EXPECT().Exists(gomock.Any(), gomock.Any()).Return(true, nil)
				bi.EXPECT().Update(gomock.Any()).Return(errors.New("error"))
				repositories.ReplaceGlobals(repositories.NewRepository(nil, nil, bi))
			},
			request:         validRequest,
			expected:        &protogen.UpdateBrokerImageResponse{},
			expectedErrCode: codes.Internal,
		},
		{
			name: "fails to retrieve broker image",
			mockSetup: func(ctrl *gomock.Controller) {
				bi := mocks.NewBrokerImageRepository(ctrl)
				bi.EXPECT().Exists(gomock.Any(), gomock.Any()).Return(true, nil)
				bi.EXPECT().Update(gomock.Any()).Return(nil)
				bi.EXPECT().Get(gomock.Any()).Return(models.BrokerImage{}, false, errors.New("error"))
				repositories.ReplaceGlobals(repositories.NewRepository(nil, nil, bi))
			},
			request:         validRequest,
			expected:        &protogen.UpdateBrokerImageResponse{},
			expectedErrCode: codes.Internal,
		},
		{
			name: "fails to find broker image",
			mockSetup: func(ctrl *gomock.Controller) {
				bi := mocks.NewBrokerImageRepository(ctrl)
				bi.EXPECT().Exists(gomock.Any(), gomock.Any()).Return(true, nil)
				bi.EXPECT().Update(gomock.Any()).Return(nil)
				bi.EXPECT().Get(gomock.Any()).Return(models.BrokerImage{}, false, nil)
				repositories.ReplaceGlobals(repositories.NewRepository(nil, nil, bi))
			},
			request:         validRequest,
			expected:        &protogen.UpdateBrokerImageResponse{},
			expectedErrCode: codes.Internal,
		},
		{
			name: "succeeded",
			mockSetup: func(ctrl *gomock.Controller) {
				bi := mocks.NewBrokerImageRepository(ctrl)
				bi.EXPECT().Exists(gomock.Any(), gomock.Any()).Return(true, nil)
				bi.EXPECT().Update(gomock.Any()).Return(nil)
				bi.EXPECT().Get(gomock.Any()).Return(validImage, true, nil)
				repositories.ReplaceGlobals(repositories.NewRepository(nil, nil, bi))
			},
			request:         validRequest,
			expected:        &protogen.UpdateBrokerImageResponse{},
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
			response, err := service.UpdateBrokerImage(context.Background(), tt.request)

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

// TestDeleteBrokerImage tests the DeleteBrokerImage function
func TestDeleteBrokerImage(t *testing.T) {
	// Prepare data
	service := &Service{}
	validRequest := &protogen.DeleteBrokerImageRequest{
		ImageId:  uuid.New().String(),
		BrokerId: uuid.New().String(),
	}

	// Prepare tests
	tests := []struct {
		name            string
		mockSetup       func(ctrl *gomock.Controller)
		request         *protogen.DeleteBrokerImageRequest
		expected        *protogen.DeleteBrokerImageResponse
		expectedErrCode codes.Code
	}{
		{
			name: "missing request body",
			mockSetup: func(ctrl *gomock.Controller) {
				bi := mocks.NewBrokerImageRepository(ctrl)
				bi.EXPECT().Exists(gomock.Any(), gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(nil, nil, bi))
			},
			request:         nil,
			expected:        &protogen.DeleteBrokerImageResponse{},
			expectedErrCode: codes.InvalidArgument,
		},
		{
			name: "fails to parse ID from request",
			mockSetup: func(ctrl *gomock.Controller) {
				bi := mocks.NewBrokerImageRepository(ctrl)
				bi.EXPECT().Exists(gomock.Any(), gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(nil, nil, bi))
			},
			request: &protogen.DeleteBrokerImageRequest{
				BrokerId: "bad-uuid",
			},
			expected:        &protogen.DeleteBrokerImageResponse{},
			expectedErrCode: codes.InvalidArgument,
		},
		{
			name: "fails to verify the image broker existence",
			mockSetup: func(ctrl *gomock.Controller) {
				bi := mocks.NewBrokerImageRepository(ctrl)
				bi.EXPECT().Exists(gomock.Any(), gomock.Any()).Return(false, errors.New("error"))
				bi.EXPECT().Delete(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(nil, nil, bi))
			},
			request:         validRequest,
			expected:        &protogen.DeleteBrokerImageResponse{},
			expectedErrCode: codes.Internal,
		},
		{
			name: "fails to find the image broker",
			mockSetup: func(ctrl *gomock.Controller) {
				bi := mocks.NewBrokerImageRepository(ctrl)
				bi.EXPECT().Exists(gomock.Any(), gomock.Any()).Return(false, nil)
				bi.EXPECT().Delete(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(nil, nil, bi))
			},
			request:         validRequest,
			expected:        &protogen.DeleteBrokerImageResponse{},
			expectedErrCode: codes.NotFound,
		},
		{
			name: "fails to delete the broker image",
			mockSetup: func(ctrl *gomock.Controller) {
				bi := mocks.NewBrokerImageRepository(ctrl)
				bi.EXPECT().Exists(gomock.Any(), gomock.Any()).Return(true, nil)
				bi.EXPECT().Delete(gomock.Any()).Return(errors.New("error"))
				repositories.ReplaceGlobals(repositories.NewRepository(nil, nil, bi))
			},
			request:         validRequest,
			expected:        &protogen.DeleteBrokerImageResponse{},
			expectedErrCode: codes.Internal,
		},
		{
			name: "succeeded",
			mockSetup: func(ctrl *gomock.Controller) {
				bi := mocks.NewBrokerImageRepository(ctrl)
				bi.EXPECT().Exists(gomock.Any(), gomock.Any()).Return(true, nil)
				bi.EXPECT().Delete(gomock.Any()).Return(nil)
				repositories.ReplaceGlobals(repositories.NewRepository(nil, nil, bi))
			},
			request:         validRequest,
			expected:        &protogen.DeleteBrokerImageResponse{},
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
			response, err := service.DeleteBrokerImage(context.Background(), tt.request)

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
