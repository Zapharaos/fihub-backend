package service

import (
	"context"
	"github.com/Zapharaos/fihub-backend/gen/go/authpb"
	"github.com/Zapharaos/fihub-backend/gen/go/userpb"
	"github.com/Zapharaos/fihub-backend/internal/mappers"
	"github.com/Zapharaos/fihub-backend/internal/models"
	"github.com/Zapharaos/fihub-backend/internal/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

type AuthService struct {
	authpb.UnimplementedAuthServiceServer
	signingKey []byte
	userClient userpb.UserServiceClient
}

const (
	JwtUserIDKey = "id"
)

// NewAuthService creates a new AuthService instance
func NewAuthService(userClient userpb.UserServiceClient) *AuthService {
	var signingKey []byte
	if viper.GetString("APP_ENV") != "production" {
		signingKey = []byte("dev-signing-key")
	} else {
		signingKey = []byte(utils.RandString(128))
	}

	return &AuthService{
		signingKey: signingKey,
		userClient: userClient,
	}
}

// GenerateToken authenticates a user and generates a JWT token for them
func (s *AuthService) GenerateToken(ctx context.Context, req *authpb.GenerateTokenRequest) (*authpb.GenerateTokenResponse, error) {
	// Try to authenticate the user
	response, err := s.userClient.AuthenticateUser(ctx, &userpb.AuthenticateUserRequest{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		zap.L().Error("failed to authenticate user", zap.Error(err))
		return nil, err
	}

	// Generate a token for the authenticated user
	user := mappers.UserFromProto(response.GetUser())
	token, err := s.createToken(user)
	if err != nil {
		zap.L().Error("failed to create token", zap.Error(err))
		return nil, err
	}

	return &authpb.GenerateTokenResponse{Token: token}, nil
}

// ValidateToken validates the JWT token and extracts the user ID
func (s *AuthService) ValidateToken(ctx context.Context, req *authpb.ValidateTokenRequest) (*authpb.ValidateTokenResponse, error) {
	claims, err := s.parseToken(req.Token)
	if err != nil {
		return nil, err
	}

	userID, ok := claims[JwtUserIDKey].(string)
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "invalid token claims")
	}

	return &authpb.ValidateTokenResponse{UserId: userID}, nil
}

// ExtractUserID extracts the user ID from the JWT token without verifying the signature
func (s *AuthService) ExtractUserID(ctx context.Context, req *authpb.ExtractUserIDRequest) (*authpb.ExtractUserIDResponse, error) {
	// Decode token without verifying signature
	token, _, err := jwt.NewParser().ParseUnverified(req.GetToken(), jwt.MapClaims{})
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid parse token unverified")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "invalid token")
	}

	userID, ok := claims[JwtUserIDKey].(string)
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "invalid token claims")
	}

	return &authpb.ExtractUserIDResponse{UserId: userID}, nil
}

func (s *AuthService) createToken(user models.User) (string, error) {
	if s.signingKey == nil {
		return "", status.Error(codes.FailedPrecondition, "signing key is nil")
	}

	claims := jwt.MapClaims{
		"exp":        jwt.NewNumericDate(time.Now().Add(time.Hour * 12)),
		"iat":        jwt.NewNumericDate(time.Now()),
		"nbf":        jwt.NewNumericDate(time.Now()),
		JwtUserIDKey: user.ID.String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.signingKey)
}

func (s *AuthService) parseToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return s.signingKey, nil
	})
	if err != nil || !token.Valid {
		return nil, status.Error(codes.InvalidArgument, "invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "invalid token claims")
	}

	return claims, nil
}
