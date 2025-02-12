package adapter

import (
	"context"
	"github.com/HunterGooD/voice_friend/user_service/internal/domain/entity"
	"github.com/HunterGooD/voice_friend/user_service/pkg/auth"
)

type TokenManagerAdapter struct {
	jwt *auth.JWT
}

func NewTokenManagerAdapter(jwt *auth.JWT) *TokenManagerAdapter {
	return &TokenManagerAdapter{jwt: jwt}
}

func (t TokenManagerAdapter) GenerateAllTokensAsync(ctx context.Context, uid, role, deviceID string) ([]string, error) {
	return t.jwt.GenerateAllTokensAsync(ctx, uid, role, deviceID)
}

func (t TokenManagerAdapter) GenerateAllTokens(ctx context.Context, uid, role, deviceID string) ([]string, error) {
	return t.jwt.GenerateAllTokens(ctx, uid, role, deviceID)
}

func (t TokenManagerAdapter) GenerateAccessToken(ctx context.Context, uid, role, deviceID string) (string, error) {
	return t.jwt.GenerateAccessToken(ctx, uid, role, deviceID)
}

func (t TokenManagerAdapter) GenerateRefreshToken(ctx context.Context, uid, role, deviceID string) (string, error) {
	return t.jwt.GenerateRefreshToken(ctx, uid, role, deviceID)
}

func (t TokenManagerAdapter) GetClaims(ctx context.Context, tokenString string) (*entity.AuthClaims, error) {
	claims, err := t.jwt.GetClaims(ctx, tokenString)
	if err != nil {
		return nil, err
	}

	return &entity.AuthClaims{
		Role:       claims.Role,
		DeviceID:   claims.DeviceID,
		Subject:    claims.GetUID(),
		JTIDY:      claims.ID,
		ExpireTime: claims.GetExpireTime(),
	}, nil
}
