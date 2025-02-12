package repository

import (
	"context"
	"fmt"
	"github.com/HunterGooD/voice_friend/user_service/internal/domain/entity"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"time"
)

/**
POSTGRESQL scheme

CREATE TABLE user_tokens (
  id SERIAL PRIMARY KEY,
  user_id INT NOT NULL,
  token_id UUID NOT NULL,
  token_type VARCHAR(10) NOT NULL, -- access или refresh
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  expires_at TIMESTAMP NOT NULL
);

CREATE TABLE blacklisted_tokens (
  id SERIAL PRIMARY KEY,
  -- FOREIGN KEY uuid ?
  token_id UUID NOT NULL,
  added_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
*/

const MaxTokensPerUser = 5
const Day = 24 * time.Hour

type TokenRepository struct {
	conn *redis.Client
}

func NewTokenRepository(conn *redis.Client) *TokenRepository {
	return &TokenRepository{conn: conn}
}

func (r *TokenRepository) StoreRefreshToken(ctx context.Context, userID, deviceID, refreshToken string) error {
	key := fmt.Sprintf("refresh_token:%s:%s", userID, deviceID)
	// TODO: MaxTokensPerUser check if new deviceID
	return r.conn.Set(ctx, key, refreshToken, 30*Day).Err()
}

func (r *TokenRepository) GetRefreshToken(ctx context.Context, userID, deviceID string) (string, error) {
	key := fmt.Sprintf("refresh_token:%s:%s", userID, deviceID)
	val, err := r.conn.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", errors.Wrap(entity.ErrNotFound, fmt.Sprintf("Error redis get user %s", userID))
		}
		return "", errors.Wrap(err, fmt.Sprintf("Unknown error redis get user %s", userID))
	}
	return val, nil
}

func (r *TokenRepository) DeleteRefreshToken(ctx context.Context, userID, deviceID string) error {
	key := fmt.Sprintf("refresh_token:%s:%s", userID, deviceID)
	if err := r.conn.Del(ctx, key).Err(); err != nil {
		if errors.Is(err, redis.Nil) {
			return errors.Wrap(entity.ErrNotFound, fmt.Sprintf("Error redis delete refresh token %s", userID))
		}
		return errors.Wrap(err, fmt.Sprintf("Unknown error redis delete refresh token %s", userID))
	}
	return nil
}
