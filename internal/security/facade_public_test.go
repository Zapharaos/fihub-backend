package security

import (
	"context"
	"errors"
	"github.com/Zapharaos/fihub-backend/gen/go/securitypb"
	"github.com/Zapharaos/fihub-backend/test/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/metadata"
	"testing"
)

// TestPublicSecurityFacade_ReplaceGlobals tests the ReplaceGlobals function
func TestPublicSecurityFacade_ReplaceGlobals(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFacade := NewPublicSecurityFacade(mocks.NewMockPermissionChecker(ctrl))
	ReplaceGlobals(mockFacade)

	assert.Equal(t, mockFacade, Facade(), "The global facade instance should be replaced correctly")
}

// TestPublicSecurityFacade_CheckPermission tests the PublicSecurityFacade.CheckPermission method
func TestPublicSecurityFacade_CheckPermission(t *testing.T) {
	facade := NewPublicSecurityFacade(nil)
	// Using incoming context here as we are testing the facade
	ctx := metadata.NewIncomingContext(context.Background(), metadata.MD{})

	// Prepare tests
	tests := []struct {
		name        string
		mockSetup   func(ctrl *gomock.Controller)
		userIDs     []uuid.UUID
		expectError bool
	}{
		{
			name: "Failure - CheckPermission Error",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockPermissionChecker(ctrl)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any()).Return((*securitypb.CheckPermissionResponse)(nil), errors.New("internal error"))
				facade = NewPublicSecurityFacade(m)
			},
			expectError: true,
		},
		{
			name: "Failure - No Permission",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockPermissionChecker(ctrl)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any()).Return(&securitypb.CheckPermissionResponse{HasPermission: false}, nil)
				facade = NewPublicSecurityFacade(m)
			},
			expectError: true,
		},
		{
			name: "Success - Self Has Permission",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockPermissionChecker(ctrl)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any()).Return(&securitypb.CheckPermissionResponse{HasPermission: true}, nil)
				facade = NewPublicSecurityFacade(m)
			},
			userIDs:     []uuid.UUID{uuid.New()},
			expectError: false,
		},
		{
			name: "Success - Has Permission",
			mockSetup: func(ctrl *gomock.Controller) {
				m := mocks.NewMockPermissionChecker(ctrl)
				m.EXPECT().CheckPermission(gomock.Any(), gomock.Any()).Return(&securitypb.CheckPermissionResponse{HasPermission: true}, nil)
				facade = NewPublicSecurityFacade(m)
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Apply mocks
			ctrl := gomock.NewController(t)
			tt.mockSetup(ctrl)
			defer ctrl.Finish()

			err := facade.CheckPermission(ctx, "permission", tt.userIDs...)
			assert.Equal(t, tt.expectError, err != nil)
		})
	}
}

// TestPublicSecurityFacade_GrpcClientAdapter tests the GrpcClientAdapter methods
func TestPublicSecurityFacade_GrpcClientAdapter(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockPublicSecurityServiceClient(ctrl)
	facade := NewPublicSecurityFacadeWithGrpcClient(mockClient)

	ctx := context.Background()
	req := &securitypb.CheckPermissionRequest{}

	mockClient.EXPECT().CheckPermission(ctx, req).Return(&securitypb.CheckPermissionResponse{HasPermission: true}, nil)

	resp, err := facade.service.CheckPermission(ctx, req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.GetHasPermission())
}
