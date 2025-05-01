package service

import (
	"context"
	"errors"
	"github.com/Zapharaos/fihub-backend/cmd/user/app/repositories"
	"github.com/Zapharaos/fihub-backend/internal/models"
	"github.com/Zapharaos/fihub-backend/internal/utils"
	"github.com/Zapharaos/fihub-backend/protogen"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"time"
)

type AuthService struct {
	protogen.UnimplementedAuthServiceServer
	signingKey []byte
}

const (
	JwtUserIDKey = "id"
)

func NewAuthService() *AuthService {
	var signingKey []byte
	if viper.GetString("APP_ENV") != "production" {
		signingKey = []byte("dev-signing-key")
	} else {
		signingKey = []byte(utils.RandString(128))
	}

	return &AuthService{
		signingKey: signingKey,
	}
}

func (s *AuthService) GenerateToken(ctx context.Context, req *protogen.GenerateTokenRequest) (*protogen.GenerateTokenResponse, error) {
	// TODO : replace with call to user service
	user, found, err := repositories.R().Authenticate(req.Email, req.Password)
	if err != nil || !found {
		zap.L().Warn("AuthService.GenerateToken.Authenticate", zap.Error(err))
		return nil, errors.New("invalid credentials")
	}

	token, err := s.createToken(user)
	if err != nil {
		zap.L().Error("AuthService.GenerateToken", zap.Error(err))
		return nil, err
	}

	return &protogen.GenerateTokenResponse{Token: token}, nil
}

func (s *AuthService) ValidateToken(ctx context.Context, req *protogen.ValidateTokenRequest) (*protogen.ValidateTokenResponse, error) {
	claims, err := s.parseToken(req.Token)
	if err != nil {
		return nil, err
	}

	userID, ok := claims[JwtUserIDKey].(string)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	return &protogen.ValidateTokenResponse{UserId: userID}, nil
}

func (s *AuthService) ExtractUserID(ctx context.Context, req *protogen.ExtractUserIDRequest) (*protogen.ExtractUserIDResponse, error) {
	// Decode token without verifying signature
	token, _, err := jwt.NewParser().ParseUnverified(req.GetToken(), jwt.MapClaims{})
	if err != nil {
		return nil, errors.New("invalid parse token unverified")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid claims")
	}

	userID, ok := claims[JwtUserIDKey].(string)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	return &protogen.ExtractUserIDResponse{UserId: userID}, nil
}

func (s *AuthService) createToken(user models.User) (string, error) {
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
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid claims")
	}

	return claims, nil
}
