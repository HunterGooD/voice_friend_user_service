package usecase

import (
	"context"
	"time"

	"github.com/HunterGooD/voice_friend/user_service/internal/domain/entity"
	"github.com/pkg/errors"
)

type AuthUsecase struct {
	userRepo  UserRepository
	tokenRepo TokenRepository
	tokenMng  TokenManager
	hashMng   HashManager
}

func NewAuthUsecase(ur UserRepository, tr TokenRepository, tm TokenManager, hs HashManager) *AuthUsecase {
	return &AuthUsecase{ur, tr, tm, hs}
}

func (u *AuthUsecase) RegisterUserUsecase(ctx context.Context, user *entity.User, deviceID string) (*entity.AuthUserResponse, error) {
	ok, err := u.userRepo.ExistUser(ctx, user.Login)
	if err != nil {
		return nil, errors.Wrap(err, "Error check user existence")
	}
	if ok {
		return nil, entity.ErrUserAlreadyExists
	}

	hashPassword, err := u.hashMng.HashPassword(user.Password)
	if err != nil {
		return nil, errors.Wrap(err, "Error create hash")
	}
	user.Password = hashPassword

	if err := u.userRepo.AddUser(ctx, user); err != nil {
		return nil, errors.Wrap(err, "Error create user")
	}

	return u.generateAuthResponse(ctx, user.UID.String(), string(user.Role), deviceID)
}

func (u *AuthUsecase) LoginUserUsecase(ctx context.Context, user *entity.User, deviceID string) (*entity.AuthUserResponse, error) {
	password, err := u.userRepo.GetUserPasswordByLogin(ctx, user.Login)
	if err != nil {
		return nil, errors.Wrap(err, "Error get user password")
	}

	isCorrect, err := u.hashMng.CheckPassword(user.Password, password)
	if err != nil {
		return nil, errors.Wrap(err, "Error check password")
	}

	if !isCorrect {
		return nil, errors.Wrap(entity.ErrInvalidPassword, "Error password not correct")
	}

	return u.generateAuthResponse(ctx, user.UID.String(), string(user.Role), deviceID)
}

func (u *AuthUsecase) UpdateAccessTokenUsecase(ctx context.Context, refreshToken string) (*entity.AuthUserResponse, error) {
	claims, err := u.tokenMng.GetClaims(ctx, refreshToken)
	if err != nil {
		return nil, errors.Wrap(err, "Error verify token")
	}
	expTime := claims.ExpireTime
	timeUntilExp := expTime.Sub(time.Now())
	if timeUntilExp <= 3*24*time.Hour {
		return u.generateAuthResponse(ctx, claims.Subject, claims.Role, claims.DeviceID)
	}

	accessToken, err := u.tokenMng.GenerateAccessToken(ctx, claims.Subject, claims.Role, claims.DeviceID)
	if err != nil {
		return nil, errors.Wrap(err, "Error generate access token")
	}

	return &entity.AuthUserResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (u *AuthUsecase) UpdateRefreshTokenUsecase(ctx context.Context, refreshToken string) (*entity.AuthUserResponse, error) {
	claims, err := u.tokenMng.GetClaims(ctx, refreshToken)
	if err != nil {
		return nil, errors.Wrap(err, "Error verify token")
	}

	return u.generateAuthResponse(ctx, claims.Subject, claims.Role, claims.DeviceID)
}

func (u *AuthUsecase) LogoutUserUsecase(ctx context.Context, refreshToken string) error {
	claims, err := u.tokenMng.GetClaims(ctx, refreshToken)
	if err != nil {
		return errors.Wrap(err, "Error verify token")
	}

	return u.tokenRepo.DeleteRefreshToken(ctx, claims.Subject, claims.DeviceID)
}

func (u *AuthUsecase) generateAuthResponse(ctx context.Context, uid, role, deviceID string) (*entity.AuthUserResponse, error) {
	tokens, err := u.tokenMng.GenerateAllTokensAsync(ctx, uid, role, deviceID)
	if err != nil {
		return nil, errors.Wrap(err, "Error create jwt")
	}

	err = u.tokenRepo.StoreRefreshToken(ctx, uid, deviceID, tokens[0])
	if err != nil {
		return nil, errors.Wrap(err, "Error store refresh token")
	}

	return &entity.AuthUserResponse{
		AccessToken:  tokens[0],
		RefreshToken: tokens[1],
	}, nil
}
