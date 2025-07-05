package service

import (
	"context"
	"github.com/Zapharaos/fihub-backend/cmd/auth/app/otp"
	"github.com/Zapharaos/fihub-backend/gen/go/authpb"
	"github.com/Zapharaos/fihub-backend/gen/go/userpb"
	"github.com/Zapharaos/fihub-backend/internal/database"
	"github.com/Zapharaos/fihub-backend/pkg/email"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"golang.org/x/text/language"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

func (s *AuthService) findUserIdentifier(ctx context.Context, req *authpb.GenerateOTPRequest) (string, error) {

	switch req.Purpose {
	case authpb.OtpPurpose_PASSWORD_RESET:
		// Must include email in request and find userID from db with email

		// Verify if the user exists
		response, err := s.userClient.GetByEmail(ctx, &userpb.GetByEmailRequest{
			Email: req.GetEmail(),
		})
		if err != nil {
			zap.L().Error("userClient.GetByEmail", zap.Error(err))
			return "", err
		}
		if response.GetUser() == nil || response.GetUser().GetId() == "" {
			return "", status.Error(codes.NotFound, otp.ErrSilentPrivate.Error())
		}

		// Return the user ID as the identifier
		return response.GetUser().GetId(), nil

	case authpb.OtpPurpose_PASSWORD_CHANGE:
		// Must retrieve userID from context

		userID, ok := ctx.Value("userID").(string)
		if !ok || userID == "" {
			return "", status.Error(codes.Unauthenticated, "user ID not found in context")
		}

		// Retrieve the metadata
		/*md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			"", status.Error(codes.Unauthenticated, "Missing metadata")
		}

		// Check if the user ID is provided in the metadata
		userIDs := md.Get("x-user-id")
		if len(userIDs) == 0 {
			"", status.Error(codes.Unauthenticated, "Missing user ID in metadata")
		}*/

		// Return the user ID from context as the identifier
		return userID, nil

	case authpb.OtpPurpose_USER_SIGNUP:
		// Must assure that the provided email is not already associated with a user

		// Verify if the user exists
		response, err := s.userClient.GetByEmail(ctx, &userpb.GetByEmailRequest{
			Email: req.GetEmail(),
		})
		if err != nil {
			st, ok := status.FromError(err)
			if ok && (st.Code() == codes.NotFound) {
				// If the user does not exist, return the email as the identifier
				return req.GetEmail(), nil
			}

			zap.L().Error("userClient.GetByEmail", zap.Error(err))
			return "", err
		}

		// If the user exists, return an error
		if response.GetUser() != nil {
			return "", status.Error(codes.AlreadyExists, otp.ErrSilentPrivate.Error())
		}

		// Return the input email as the identifier
		return req.GetEmail(), nil

	default:
		return "", status.Error(codes.InvalidArgument, otp.ErrArgumentInvalid.Error())
	}
}

func (s *AuthService) setupForFinalRequest(ctx context.Context, purpose authpb.OtpPurpose, identifier string) (*authpb.ValidateOTPResponse, error) {
	// Prepare next step data
	requestID := uuid.New().String()
	requestTimeLimit := otp.GetFinalRequestTimeLimit()
	requestKey := otp.BuildOtpRequestKey(identifier, purpose)

	// Prepare pipeline to store request ID and delete OTP
	pipe := database.DB().Redis().Client.TxPipeline()
	pipe.Del(ctx, otp.BuildOtpKey(identifier, purpose))
	pipe.SetEx(ctx, requestKey, requestID, requestTimeLimit)

	// Execute pipeline
	if _, err := pipe.Exec(ctx); err != nil {
		zap.L().Error("failed to store request ID", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to store request ID")
	}

	return &authpb.ValidateOTPResponse{
		RequestId: requestID,
		ExpiresAt: timestamppb.New(time.Now().Add(requestTimeLimit)),
	}, nil
}

// GenerateOTP generates a one-time password (OTP) for the user
func (s *AuthService) GenerateOTP(ctx context.Context, req *authpb.GenerateOTPRequest) (*authpb.GenerateOTPResponse, error) {
	// TODO : move handlers middleware rate limiter to here? attempts count?

	// Retrieve the user identifier based on the request purpose
	identifier, err := s.findUserIdentifier(ctx, req)
	if err != nil {
		st, ok := status.FromError(err)
		if ok && (st.Code() == codes.AlreadyExists || st.Code() == codes.NotFound) {
			// Return a generic response to avoid leaking info regarding the user existence
			return &authpb.GenerateOTPResponse{}, nil
		}
		return nil, err
	}

	// Check for existing OTP with identifier and purpose
	otpKey := otp.BuildOtpKey(identifier, req.GetPurpose())
	ttl, err := otp.GetTtlForRedisKey(ctx, otpKey)
	if err != nil {
		zap.L().Error("failed to get OTP expiration", zap.Error(err))
		return nil, err
	}
	if ttl > 0 {
		// If an OTP already exists, return the existing OTP expiration time
		return &authpb.GenerateOTPResponse{
			Identifier: identifier,
			ExpiresAt:  timestamppb.New(time.Now().Add(ttl)),
			Error:      proto.String(otp.ErrRequestActive),
		}, nil
	}

	// Prepare OTP data
	otpTimeLimit := otp.GetOtpTimeLimit()
	expiresAt := timestamppb.New(time.Now().Add(otpTimeLimit))
	otpValue, otpHash := otp.Generate()

	// Store OTP in Redis with identifier and purpose as key
	err = database.DB().Redis().Client.Set(ctx, otpKey, otpHash, otpTimeLimit).Err()
	if err != nil {
		zap.L().Error("failed to store OTP", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to store OTP")
	}

	// Prepare otp email content
	userLanguage := language.MustParse(req.GetLanguage())
	subject, plainTextContent, htmlContent, err := otp.BuildOtpEmailContents(userLanguage, otpValue, otpTimeLimit)
	if err != nil {
		// Delete the request since the email could not be sent
		otp.CleanupRedisKey(ctx, otpKey)

		zap.L().Error("failed to build OTP email content", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to build OTP email content")
	}

	// If debug mode is on, log the OTP value and skip sending email
	if viper.GetString("APP_ENV") != "production" {
		zap.L().Info("OTP generated (debug mode)", zap.String("identifier", identifier), zap.String("otp", otpValue), zap.String("limit", otpTimeLimit.String()))
		return &authpb.GenerateOTPResponse{
			Identifier: identifier,
			ExpiresAt:  expiresAt,
		}, nil
	}

	// Send email
	err = email.S().Send(req.GetEmail(), subject, plainTextContent, htmlContent)
	if err != nil {
		// Delete the request since the email could not be sent
		otp.CleanupRedisKey(ctx, otpKey)

		zap.L().Error("Failed to send OTP email", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to send OTP email")
	}

	return &authpb.GenerateOTPResponse{
		Identifier: identifier,
		ExpiresAt:  expiresAt,
	}, nil
}

// ValidateOTP validates the one-time password (OTP) for the user
func (s *AuthService) ValidateOTP(ctx context.Context, req *authpb.ValidateOTPRequest) (*authpb.ValidateOTPResponse, error) {
	// TODO : same as for GenerateOTP "todo" for rate limiting

	var identifier string
	switch req.Purpose {
	case authpb.OtpPurpose_PASSWORD_CHANGE:
		// User authenticated, retrieve his identifier from context
		userID, ok := ctx.Value("userID").(string)
		if !ok || userID == "" {
			return nil, status.Error(codes.Unauthenticated, "user ID not found in context")
		}

		// Retrieve the metadata
		/*md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			nil, status.Error(codes.Unauthenticated, "Missing metadata")
		}

		// Check if the user ID is provided in the metadata
		userIDs := md.Get("x-user-id")
		if len(userIDs) == 0 {
			nil, status.Error(codes.Unauthenticated, "Missing user ID in metadata")
		}*/
	case authpb.OtpPurpose_PASSWORD_RESET, authpb.OtpPurpose_USER_SIGNUP:
		// User not authenticated, must include his identifier in request
		identifier = req.GetIdentifier()
	default:
		return nil, status.Error(codes.InvalidArgument, otp.ErrArgumentInvalid.Error())
	}

	// Validate the OTP
	err := otp.IsOtpValid(ctx, req.GetPurpose(), identifier, req.GetOtp())
	if err != nil {
		return nil, err
	}

	// Setup for the final request
	return s.setupForFinalRequest(ctx, req.GetPurpose(), identifier)
}
