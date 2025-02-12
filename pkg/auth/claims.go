package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type AuthClaims struct {
	Role     string `json:"role"`
	DeviceID string `json:"device_id"`
	jwt.RegisteredClaims
}

func (a *AuthClaims) GetUID() string {
	return a.Subject
}

func (a *AuthClaims) GetRole() string {
	return a.Role
}

func (a *AuthClaims) GetDeviceId() string {
	return a.DeviceID
}

func (a *AuthClaims) GetExpireTime() time.Time {
	return a.ExpiresAt.Time
}
