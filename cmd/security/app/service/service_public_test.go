package service

import (
	"context"
	"errors"
	"github.com/Zapharaos/fihub-backend/cmd/security/app/repositories"
	"github.com/Zapharaos/fihub-backend/gen/go/securitypb"
	"github.com/Zapharaos/fihub-backend/internal/models"
	"github.com/Zapharaos/fihub-backend/test/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"testing"
)

func TestPublicService_CheckPermission(t *testing.T) {
	// Prepare data
	service := &PublicService{}
	userID := uuid.New()
	validContext := metadata.NewIncomingContext(context.Background(), metadata.MD{
		"x-user-id": {uuid.New().String()},
	})
	validRequest := &securitypb.CheckPermissionRequest{
		UserId:     userID.String(),
		Permission: "example.permission",
	}
	validResponse := []models.RoleWithPermissions{
		{
			Permissions: []models.Permission{
				{
					Value: "example.permission",
				},
			},
		},
	}

	// Define the test cases
	tests := []struct {
		name            string
		mockSetup       func(ctrl *gomock.Controller) context.Context
		request         *securitypb.CheckPermissionRequest
		expected        bool
		expectedErrCode codes.Code
	}{
		{
			name: "missing metadata in context",
			mockSetup: func(ctrl *gomock.Controller) context.Context {
				r := mocks.NewSecurityRoleRepository(ctrl)
				r.EXPECT().ListWithPermissionsByUserId(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(r, nil))
				// Create a new context without metadata
				return context.Background()
			},
			request:         validRequest,
			expected:        false,
			expectedErrCode: codes.Unauthenticated,
		},
		{
			name: "missing user ID in metadata",
			mockSetup: func(ctrl *gomock.Controller) context.Context {
				r := mocks.NewSecurityRoleRepository(ctrl)
				r.EXPECT().ListWithPermissionsByUserId(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(r, nil))
				// Create a new context without userID in metadata
				return metadata.NewIncomingContext(context.Background(), metadata.MD{})
			},
			request:         validRequest,
			expected:        false,
			expectedErrCode: codes.Unauthenticated,
		},
		{
			name: "fails to list user permissions",
			mockSetup: func(ctrl *gomock.Controller) context.Context {
				r := mocks.NewSecurityRoleRepository(ctrl)
				r.EXPECT().ListWithPermissionsByUserId(gomock.Any()).Return(nil, errors.New("error"))
				repositories.ReplaceGlobals(repositories.NewRepository(r, nil))
				// Create a new context without userID in metadata
				return validContext
			},
			request:         validRequest,
			expected:        false,
			expectedErrCode: codes.Internal,
		},
		{
			name: "user has no permission",
			mockSetup: func(ctrl *gomock.Controller) context.Context {
				r := mocks.NewSecurityRoleRepository(ctrl)
				r.EXPECT().ListWithPermissionsByUserId(gomock.Any()).Return(models.RolesWithPermissions{}, nil)
				repositories.ReplaceGlobals(repositories.NewRepository(r, nil))
				// Create a new context without userID in metadata
				return validContext
			},
			request:         validRequest,
			expected:        false,
			expectedErrCode: codes.PermissionDenied,
		},
		{
			name: "successful self permission check",
			mockSetup: func(ctrl *gomock.Controller) context.Context {
				r := mocks.NewSecurityRoleRepository(ctrl)
				r.EXPECT().ListWithPermissionsByUserId(gomock.Any()).Times(0)
				repositories.ReplaceGlobals(repositories.NewRepository(r, nil))
				// Create a new context without userID in metadata
				return metadata.NewIncomingContext(context.Background(), metadata.MD{
					"x-user-id": {userID.String()},
				})
			},
			request:         validRequest,
			expected:        true,
			expectedErrCode: codes.OK,
		},
		{
			name: "successful permission check",
			mockSetup: func(ctrl *gomock.Controller) context.Context {
				r := mocks.NewSecurityRoleRepository(ctrl)
				r.EXPECT().ListWithPermissionsByUserId(gomock.Any()).Return(validResponse, nil)
				repositories.ReplaceGlobals(repositories.NewRepository(r, nil))
				// Create a new context without userID in metadata
				return validContext
			},
			request:         validRequest,
			expected:        true,
			expectedErrCode: codes.OK,
		},
	}

	// Run the test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Apply mocks
			ctrl := gomock.NewController(t)
			ctx := tt.mockSetup(ctrl)
			defer ctrl.Finish()

			// Call service
			response, err := service.CheckPermission(ctx, tt.request)

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
			assert.Equal(t, tt.expected, response.GetHasPermission())
		})
	}
}
