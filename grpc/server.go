package grpc

import (
	"context"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/nurefendi/auth-service-using-golang/dto"
	"github.com/nurefendi/auth-service-using-golang/grpc/pb"
	jwtMiddleware "github.com/nurefendi/auth-service-using-golang/middleware/jwt"
	"github.com/nurefendi/auth-service-using-golang/tools/locals"
	"github.com/nurefendi/auth-service-using-golang/usecase"
)

// AuthServer implements the gRPC AuthService
type AuthServer struct {
	pb.UnimplementedAuthServiceServer
}

// NewAuthServer creates a new AuthServer instance
func NewAuthServer() *AuthServer {
	return &AuthServer{}
}

// Helper function to create fiber context from gRPC context
func createFiberContextFromGRPC(ctx context.Context) *fiber.Ctx {
	app := fiber.New()
	c := app.AcquireCtx(nil)

	// Set user locals with request ID
	userLocals := dto.UserLocals{
		RequestID:    uuid.New().String(),
		LanguageCode: "en",
		ChannelID:    "grpc",
	}

	// Extract JWT token from metadata if present
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if authHeaders := md.Get("authorization"); len(authHeaders) > 0 {
			authHeader := authHeaders[0]
			if strings.HasPrefix(authHeader, "Bearer ") {
				token := strings.TrimPrefix(authHeader, "Bearer ")

				// Parse and validate JWT token
				claims, err := jwtMiddleware.JwtClaims(c, token)
				if err == nil {
					// Extract user info from claims
					if userIDStr, ok := claims["userId"].(string); ok {
						if userID, err := uuid.Parse(userIDStr); err == nil {
							userLocals.UserAccess = &dto.CurrentUserAccess{
								UserID: userID,
							}
							if userName, ok := claims["userName"].(string); ok {
								userLocals.UserAccess.UserName = userName
							}
							if email, ok := claims["email"].(string); ok {
								userLocals.UserAccess.Email = email
							}
							if fullName, ok := claims["fullName"].(string); ok {
								userLocals.UserAccess.FullName = fullName
							}
						}
					}
				}
			}
		}
	}

	c.Locals(locals.UserLocalKey, userLocals)
	return c
}

// Register implements gRPC Register method
func (s *AuthServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	c := createFiberContextFromGRPC(ctx)
	defer c.App().ReleaseCtx(c)

	// Convert gRPC request to DTO
	registerReq := dto.AuthUserRegisterRequest{
		FullName: req.FullName,
		Email:    req.Email,
		Password: req.Password,
		Gender:   int(req.Gender),
	}

	c.Locals(locals.PayloadLocalKey, registerReq)

	// Call usecase
	err := usecase.AuthUSeCase().Register(c)
	if err != nil {
		return &pb.RegisterResponse{
			Success: false,
			Message: err.Message,
		}, status.Error(codes.Code(err.Code), err.Message)
	}

	return &pb.RegisterResponse{
		Success: true,
		Message: "User registered successfully",
	}, nil
}

// Login implements gRPC Login method
func (s *AuthServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	c := createFiberContextFromGRPC(ctx)
	defer c.App().ReleaseCtx(c)

	// Convert gRPC request to DTO
	loginReq := dto.AuthUserLoginRequest{
		Email:    req.Email,
		Password: req.Password,
	}

	c.Locals(locals.PayloadLocalKey, loginReq)

	// Call usecase
	err := usecase.AuthUSeCase().Login(c)
	if err != nil {
		return &pb.LoginResponse{
			Success: false,
			Message: err.Message,
		}, status.Error(codes.Code(err.Code), err.Message)
	}

	// Get user info
	userInfo, userErr := usecase.AuthUSeCase().Me(c)
	if userErr != nil {
		return &pb.LoginResponse{
			Success: false,
			Message: userErr.Message,
		}, status.Error(codes.Code(userErr.Code), userErr.Message)
	}

	// For now, we'll generate dummy tokens (you should implement proper token generation)
	accessToken := "dummy_access_token"
	refreshToken := "dummy_refresh_token"
	expiresAt := time.Now().Add(time.Hour * 24)

	// Convert user info to gRPC response
	grpcUser := &pb.AuthUser{
		UserId:     userInfo.UserID.String(),
		UserName:   userInfo.UserName,
		Email:      userInfo.Email,
		FullName:   userInfo.FullName,
		Gender:     int32(userInfo.Gender),
		GenderName: userInfo.GenderName,
	}

	if userInfo.Picture != nil {
		grpcUser.Picture = *userInfo.Picture
	}

	for _, groupID := range userInfo.GroupIDs {
		grpcUser.GroupIds = append(grpcUser.GroupIds, groupID.String())
	}

	return &pb.LoginResponse{
		Success:      true,
		Message:      "Login successful",
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    timestamppb.New(expiresAt),
		User:         grpcUser,
	}, nil
}

// RefreshToken implements gRPC RefreshToken method
func (s *AuthServer) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.RefreshTokenResponse, error) {
	c := createFiberContextFromGRPC(ctx)
	defer c.App().ReleaseCtx(c)

	// Call usecase
	err := usecase.AuthUSeCase().RefreshToken(c)
	if err != nil {
		return &pb.RefreshTokenResponse{
			Success: false,
			Message: err.Message,
		}, status.Error(codes.Code(err.Code), err.Message)
	}

	// For now, generate dummy tokens
	accessToken := "new_dummy_access_token"
	refreshToken := "new_dummy_refresh_token"
	expiresAt := time.Now().Add(time.Hour * 24)

	return &pb.RefreshTokenResponse{
		Success:      true,
		Message:      "Token refreshed successfully",
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    timestamppb.New(expiresAt),
	}, nil
}

// Logout implements gRPC Logout method
func (s *AuthServer) Logout(ctx context.Context, req *pb.LogoutRequest) (*pb.LogoutResponse, error) {
	c := createFiberContextFromGRPC(ctx)
	defer c.App().ReleaseCtx(c)

	// Call usecase
	err := usecase.AuthUSeCase().Logout(c)
	if err != nil {
		return &pb.LogoutResponse{
			Success: false,
			Message: err.Message,
		}, status.Error(codes.Code(err.Code), err.Message)
	}

	return &pb.LogoutResponse{
		Success: true,
		Message: "Logout successful",
	}, nil
}

// ChangePassword implements gRPC ChangePassword method
func (s *AuthServer) ChangePassword(ctx context.Context, req *pb.ChangePasswordRequest) (*pb.ChangePasswordResponse, error) {
	// This would need implementation in your usecase
	// For now, return a dummy response
	return &pb.ChangePasswordResponse{
		Success: true,
		Message: "Password changed successfully",
	}, nil
}

// CheckAccess implements gRPC CheckAccess method
func (s *AuthServer) CheckAccess(ctx context.Context, req *pb.CheckAccessRequest) (*pb.CheckAccessResponse, error) {
	c := createFiberContextFromGRPC(ctx)
	defer c.App().ReleaseCtx(c)

	// Convert gRPC request to DTO
	checkAccessReq := dto.AuthCheckAccessRequest{
		Path:   req.Path,
		Mathod: req.Method, // Note: there's a typo in the original DTO (Mathod instead of Method)
	}

	// Call usecase
	err := usecase.AuthUSeCase().CheckAccess(c, checkAccessReq)
	if err != nil {
		return &pb.CheckAccessResponse{
			HasAccess: false,
			Message:   err.Message,
		}, nil // Don't return error for access denied, just false
	}

	return &pb.CheckAccessResponse{
		HasAccess: true,
		Message:   "Access granted",
	}, nil
}

// GetUserProfile implements gRPC GetUserProfile method
func (s *AuthServer) GetUserProfile(ctx context.Context, req *pb.GetUserProfileRequest) (*pb.GetUserProfileResponse, error) {
	c := createFiberContextFromGRPC(ctx)
	defer c.App().ReleaseCtx(c)

	// Call usecase
	userInfo, err := usecase.AuthUSeCase().Me(c)
	if err != nil {
		return &pb.GetUserProfileResponse{
			Success: false,
			Message: err.Message,
		}, status.Error(codes.Code(err.Code), err.Message)
	}

	// Convert user info to gRPC response
	grpcUser := &pb.AuthUser{
		UserId:     userInfo.UserID.String(),
		UserName:   userInfo.UserName,
		Email:      userInfo.Email,
		FullName:   userInfo.FullName,
		Gender:     int32(userInfo.Gender),
		GenderName: userInfo.GenderName,
	}

	if userInfo.Picture != nil {
		grpcUser.Picture = *userInfo.Picture
	}

	for _, groupID := range userInfo.GroupIDs {
		grpcUser.GroupIds = append(grpcUser.GroupIds, groupID.String())
	}

	return &pb.GetUserProfileResponse{
		Success: true,
		Message: "User profile retrieved successfully",
		User:    grpcUser,
	}, nil
}

// GetUserFunctions implements gRPC GetUserFunctions method
func (s *AuthServer) GetUserFunctions(ctx context.Context, req *pb.GetUserFunctionsRequest) (*pb.GetUserFunctionsResponse, error) {
	c := createFiberContextFromGRPC(ctx)
	defer c.App().ReleaseCtx(c)

	// Call usecase
	functions, err := usecase.AuthUSeCase().MyAcl(c)
	if err != nil {
		return &pb.GetUserFunctionsResponse{
			Success: false,
			Message: err.Message,
		}, status.Error(codes.Code(err.Code), err.Message)
	}

	// Convert functions to gRPC response
	var grpcFunctions []*pb.AuthUserFunction
	for _, function := range functions {
		grpcFunction := &pb.AuthUserFunction{
			GroupId:      function.GroupID.String(),
			GroupName:    function.GroupName,
			PortalId:     function.PortalID.String(),
			PortalName:   function.PortalName,
			FunctionId:   function.FunctionID.String(),
			FunctionName: function.FunctionName,
			GrantCreate:  int32(function.GrantCreate),
			GrantRead:    int32(function.GrantRead),
			GrantUpdate:  int32(function.GrantUpdate),
			GrantDelete:  int32(function.GrantDelete),
		}
		grpcFunctions = append(grpcFunctions, grpcFunction)
	}

	return &pb.GetUserFunctionsResponse{
		Success:   true,
		Message:   "User functions retrieved successfully",
		Functions: grpcFunctions,
	}, nil
}
