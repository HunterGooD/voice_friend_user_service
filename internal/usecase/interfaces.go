package usecase

import (
	"context"
	"github.com/HunterGooD/voice_friend_user_service/internal/domain/entity"
)

type UserRepository interface {
	AddUser(ctx context.Context, user *entity.User) error
	ExistUser(ctx context.Context, login string) (bool, error)
	GetUserPasswordByLogin(ctx context.Context, login string) (string, error)
}

type TokenRepository interface {
	StoreRefreshToken(ctx context.Context, userID, deviceID, refreshToken string) error
	GetRefreshToken(ctx context.Context, userID, deviceID string) (string, error)
	DeleteRefreshToken(ctx context.Context, userID, deviceID string) error
}

type TokenManager interface {
	GenerateAllTokensAsync(ctx context.Context, uid, role, deviceID string) ([]string, error)
	GenerateAllTokens(ctx context.Context, uid, role, deviceID string) ([]string, error)
	GenerateAccessToken(ctx context.Context, uid, role, deviceID string) (string, error)
	GenerateRefreshToken(ctx context.Context, uid, role, deviceID string) (string, error)
	GetClaims(ctx context.Context, tokenString string) (*entity.AuthClaims, error)
}

type HashManager interface {
	HashPassword(password string) (string, error)
	CheckPassword(password, hashedPassword string) (bool, error)
}
