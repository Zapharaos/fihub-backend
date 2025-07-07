package service

import (
	"github.com/Zapharaos/fihub-backend/gen/go/authpb"
	"github.com/Zapharaos/fihub-backend/gen/go/userpb"
	"github.com/Zapharaos/fihub-backend/internal/utils"
	"github.com/spf13/viper"
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
