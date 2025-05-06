package service

import (
	"context"
	"github.com/Zapharaos/fihub-backend/gen/go/authpb"
	"github.com/Zapharaos/fihub-backend/gen/go/userpb"
	"github.com/Zapharaos/fihub-backend/internal/models"
	"github.com/Zapharaos/fihub-backend/test/mocks"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"
	"time"
)

// TestGenerateToken tests the AuthService.GenerateToken service
func TestGenerateToken(t *testing.T) {
	validRequest := &authpb.GenerateTokenRequest{
		Email:    "email",
		Password: "password",
	}

	// Define tests
	tests := []struct {
		name            string
		serviceSetup    func(ctrl *gomock.Controller) *AuthService
		request         *authpb.GenerateTokenRequest
		expectToken     bool
		expectedErrCode codes.Code
	}{
		{
			name: "fails to authenticate user",
			serviceSetup: func(ctrl *gomock.Controller) *AuthService {
				userClient := mocks.NewMockUserServiceClient(gomock.NewController(t))
				userClient.EXPECT().AuthenticateUser(gomock.Any(), gomock.Any()).Return(nil, status.Error(codes.Internal, "internal error"))
				return NewAuthService(userClient)
			},
			request:         validRequest,
			expectToken:     false,
			expectedErrCode: codes.Internal,
		},
		{
			name: "fails to create token",
			serviceSetup: func(ctrl *gomock.Controller) *AuthService {
				userClient := mocks.NewMockUserServiceClient(gomock.NewController(t))
				userClient.EXPECT().AuthenticateUser(gomock.Any(), gomock.Any()).Return(&userpb.AuthenticateUserResponse{
					User: &userpb.User{Id: uuid.New().String()},
				}, nil)
				service := NewAuthService(userClient)
				service.signingKey = nil
				return service
			},
			request:         validRequest,
			expectToken:     false,
			expectedErrCode: codes.FailedPrecondition,
		},
		{
			name: "succeeds",
			serviceSetup: func(ctrl *gomock.Controller) *AuthService {
				userClient := mocks.NewMockUserServiceClient(gomock.NewController(t))
				userClient.EXPECT().AuthenticateUser(gomock.Any(), gomock.Any()).Return(&userpb.AuthenticateUserResponse{
					User: &userpb.User{Id: uuid.New().String()},
				}, nil)
				return NewAuthService(userClient)
			},
			request:         validRequest,
			expectToken:     true,
			expectedErrCode: codes.Internal,
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Prepare mocks
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// Call service
			service := tt.serviceSetup(ctrl)
			response, err := service.GenerateToken(context.Background(), tt.request)

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
			if tt.expectToken {
				assert.NotEmpty(t, response.Token)
			}
		})
	}
}

// TestValidateToken tests the AuthService.ValidateToken service
func TestValidateToken(t *testing.T) {
	// Data
	userID := uuid.New()

	tests := []struct {
		name           string
		signingKey     []byte
		token          string
		expectedUserID string
		expectError    bool
	}{
		{
			name:        "fails with invalid token",
			signingKey:  []byte("test-signing-key"),
			token:       "invalid-token",
			expectError: true,
		},
		{
			name:       "fails with missing user ID claim",
			signingKey: []byte("test-signing-key"),
			token: func() string {
				service := &AuthService{signingKey: []byte("test-signing-key")}
				claims := jwt.MapClaims{
					"exp": jwt.NewNumericDate(time.Now().Add(time.Hour)),
				}
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
				signedToken, _ := token.SignedString(service.signingKey)
				return signedToken
			}(),
			expectError: true,
		},
		{
			name:       "successfully validates token",
			signingKey: []byte("test-signing-key"),
			token: func() string {
				service := &AuthService{signingKey: []byte("test-signing-key")}
				user := models.User{ID: userID}
				token, _ := service.createToken(user)
				return token
			}(),
			expectedUserID: func() string {
				return userID.String()
			}(),
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &AuthService{
				signingKey: tt.signingKey,
			}

			response, err := service.ValidateToken(context.Background(), &authpb.ValidateTokenRequest{
				Token: tt.token,
			})

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, response)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, response)
				assert.Equal(t, tt.expectedUserID, response.UserId)
			}
		})
	}
}

// TestExtractUserID tests the AuthService.ExtractUserID service
func TestExtractUserID(t *testing.T) {
	// Data
	userID := uuid.New()

	tests := []struct {
		name           string
		token          string
		expectedUserID string
		expectError    bool
	}{
		{
			name:        "fails with invalid token",
			token:       "invalid-token",
			expectError: true,
		},
		{
			name: "fails with missing user ID claim",
			token: func() string {
				claims := jwt.MapClaims{
					"exp": jwt.NewNumericDate(time.Now().Add(time.Hour)),
				}
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
				signedToken, _ := token.SignedString([]byte("test-signing-key"))
				return signedToken
			}(),
			expectError: true,
		},
		{
			name: "successfully extracts user ID",
			token: func() string {
				claims := jwt.MapClaims{
					"exp":        jwt.NewNumericDate(time.Now().Add(time.Hour)),
					JwtUserIDKey: userID.String(),
				}
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
				signedToken, _ := token.SignedString([]byte("test-signing-key"))
				return signedToken
			}(),
			expectedUserID: func() string {
				return userID.String()
			}(),
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &AuthService{}

			response, err := service.ExtractUserIDFromToken(context.Background(), &authpb.ExtractUserIDFromTokenRequest{
				Token: tt.token,
			})

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, response)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, response)
				assert.Equal(t, tt.expectedUserID, response.UserId)
			}
		})
	}
}

func TestCreateToken(t *testing.T) {
	tests := []struct {
		name        string
		signingKey  []byte
		user        models.User
		expectError bool
	}{
		{
			name:        "fails with nil signing key",
			signingKey:  nil,
			user:        models.User{ID: uuid.New()},
			expectError: true,
		},
		{
			name:       "successfully creates token",
			signingKey: []byte("test-signing-key"),
			user: models.User{
				ID: uuid.New(),
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &AuthService{
				signingKey: tt.signingKey,
			}

			token, err := service.createToken(tt.user)

			if tt.expectError {
				assert.Error(t, err)
				assert.Empty(t, token)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, token)
			}
		})
	}
}

func TestParseToken(t *testing.T) {
	tests := []struct {
		name        string
		signingKey  []byte
		token       string
		expectError bool
	}{
		{
			name:        "fails with invalid token",
			signingKey:  []byte("test-signing-key"),
			token:       "invalid-token",
			expectError: true,
		},
		{
			name:        "fails with nil signing key",
			signingKey:  nil,
			token:       "",
			expectError: true,
		},
		{
			name:       "successfully parses valid token",
			signingKey: []byte("test-signing-key"),
			token: func() string {
				service := &AuthService{signingKey: []byte("test-signing-key")}
				user := models.User{ID: uuid.New()}
				token, _ := service.createToken(user)
				return token
			}(),
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &AuthService{
				signingKey: tt.signingKey,
			}

			claims, err := service.parseToken(tt.token)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, claims)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, claims)
			}
		})
	}
}
