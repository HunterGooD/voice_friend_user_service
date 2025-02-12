package handler

import (
	"context"
	"google.golang.org/grpc"

	"github.com/HunterGooD/voice_friend/user_service/internal/domain/entity"
	// TODO: domain interface
	"github.com/HunterGooD/voice_friend/user_service/pkg/logger"
	"github.com/HunterGooD/voice_friend/user_service/pkg/utils"
	pd "github.com/HunterGooD/voice_friend_contracts/gen/go/user_service"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthUsecase interface {
	RegisterUserUsecase(ctx context.Context, user *entity.User, deviceID string) (*entity.AuthUserResponse, error)
	LoginUserUsecase(ctx context.Context, user *entity.User, deviceID string) (*entity.AuthUserResponse, error)
	LogoutUserUsecase(ctx context.Context, refreshToken string) error
	UpdateRefreshTokenUsecase(ctx context.Context, refreshToken string) (*entity.AuthUserResponse, error)
	UpdateAccessTokenUsecase(ctx context.Context, refreshToken string) (*entity.AuthUserResponse, error)
}

type GRPCServer interface {
	GetServer() *grpc.Server
}

type AuthHandler struct {
	pd.UnimplementedAuthServer
	authUsecase AuthUsecase

	log logger.Logger
}

func NewAuthHandler(gRPCServer GRPCServer, uu AuthUsecase, log logger.Logger) {
	authHandler := &AuthHandler{authUsecase: uu, log: log}
	pd.RegisterAuthServer(gRPCServer.GetServer(), authHandler)
}

func (h *AuthHandler) Register(ctx context.Context, req *pd.RegisterRequest) (*pd.AuthResponse, error) {
	if req.Login == "" || req.Name == "" || req.Password == "" || req.DeviceID == "" {
		h.log.Warn("Request without params")
		return nil, status.Errorf(codes.InvalidArgument, "Request missing required field: Name or Login or Password or DeviceID")
	}

	err := h.validator(req.Email, req.Phone)
	if err != nil {
		return nil, err
	}

	user := entity.User{
		Login:          req.Login,
		Name:           req.Name,
		Email:          req.Email,
		Password:       req.Password,
		ProfilePicture: req.ProfilePicture,
		Phone:          req.Phone,
	}
	res, err := h.authUsecase.RegisterUserUsecase(ctx, &user, req.DeviceID)
	if err != nil {
		if errors.Is(err, entity.ErrUserAlreadyExists) {
			h.log.Warn("User already exists", map[string]interface{}{
				"user":  user.UID.String(),
				"error": err,
			})
			return nil, status.Errorf(codes.AlreadyExists, "Request user exists")
		}
		h.log.Error("Error unknown ", err)
		return nil, status.Errorf(codes.Internal, "Unknown error %+v", err)
	}

	return &pd.AuthResponse{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
	}, nil
}

func (h *AuthHandler) Login(ctx context.Context, req *pd.LoginRequest) (*pd.AuthResponse, error) {
	if req.Login == "" || req.Password == "" || req.DeviceID == "" {
		h.log.Warn("Request without params")
		return nil, status.Errorf(codes.InvalidArgument, "Request missing required field: Login or Password or DeviceID")
	}

	err := h.validator(req.Email, req.Phone)
	if err != nil {
		return nil, err
	}

	user := entity.User{
		Login:    req.Login,
		Email:    req.Email,
		Password: req.Password,
		Phone:    req.Phone,
	}

	res, err := h.authUsecase.LoginUserUsecase(ctx, &user, req.DeviceID)
	if err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			h.log.Warn("User not found", map[string]interface{}{
				"user":  user.UID.String(),
				"error": err,
			})
			return nil, status.Errorf(codes.NotFound, "Request user not found")
		}
		h.log.Error("Error unknown ", err)
		return nil, status.Errorf(codes.Internal, "Unknown error %+v", err)
	}

	return &pd.AuthResponse{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
	}, nil
}

func (h *AuthHandler) validator(email, phone *string) error {
	if email != nil {
		if utils.ValidateEmail(*email) != true {
			h.log.Warn("Email dont valid", map[string]interface{}{
				"email": *email,
			})
			return status.Errorf(codes.InvalidArgument, "Request email invalid validation")
		}
	}
	if phone != nil {
		if utils.ValidatePhone(*phone) != true {
			h.log.Warn("Phone number dont valid", map[string]interface{}{
				"phone": *phone,
			})
			return status.Errorf(codes.InvalidArgument, "Request phone invalid validation")
		}
	}
	return nil
}

func (h *AuthHandler) UpdateAccessToken(ctx context.Context, req *pd.RefreshToken) (*pd.AuthResponse, error) {
	res, err := h.authUsecase.UpdateAccessTokenUsecase(ctx, req.Token)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Internal error %v", err)
	}
	return &pd.AuthResponse{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
	}, nil
}

func (h *AuthHandler) UpdateRefreshToken(ctx context.Context, req *pd.RefreshToken) (*pd.AuthResponse, error) {
	res, err := h.authUsecase.UpdateRefreshTokenUsecase(ctx, req.Token)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Internal error %v", err)
	}
	return &pd.AuthResponse{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
	}, nil
}

func (h *AuthHandler) LogOut(ctx context.Context, req *pd.LogoutRequest) (*pd.LogoutResponse, error) {
	// TODO: user uid, deviceID
	err := h.authUsecase.LogoutUserUsecase(ctx, req.RefreshToken)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Unknown error %+v", err)
	}
	return &pd.LogoutResponse{Success: true}, nil
}
