package auth

import (
	"context"
	"crypto/rsa"
	"os"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

var ErrKeyMissing = errors.New("key missing")

type JWT struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey

	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
	issuer               string
	audience             []string
}

func NewJWTGenerator(privateKey *rsa.PrivateKey, publicKey *rsa.PublicKey, issuer string, accessToken, refreshToken time.Duration, audience []string) *JWT {
	return &JWT{
		privateKey,
		publicKey,
		accessToken,
		refreshToken,
		issuer, audience,
	}
}

func NewJWTGeneratorWithPrivateKey(privateKey *rsa.PrivateKey, issuer string, accessToken, refreshToken time.Duration, audience []string) *JWT {
	return &JWT{
		privateKey,
		nil,
		accessToken,
		refreshToken,
		issuer,
		audience,
	}
}

func NewJWTGeneratorWithPublicKey(publicKey *rsa.PublicKey, issuer string, accessToken, refreshToken time.Duration, audience []string) *JWT {
	return &JWT{
		nil,
		publicKey,
		accessToken,
		refreshToken,
		issuer,
		audience,
	}
}

func (j *JWT) LoadPrivateKeyFromFile(certPath string) (*rsa.PrivateKey, error) {
	keyData, err := os.ReadFile(certPath)
	if err != nil {
		return nil, err
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(keyData)
	if err != nil {
		return nil, err
	}

	j.privateKey = privateKey
	return privateKey, nil
}

func (j *JWT) LoadPublicKeyFromFile(certPath string) (*rsa.PublicKey, error) {
	keyData, err := os.ReadFile(certPath)
	if err != nil {
		return nil, err
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(keyData)
	if err != nil {
		return nil, err
	}

	j.publicKey = publicKey
	return publicKey, nil
}

func (j *JWT) IsValidJWT(ctx context.Context, tokenString string) (bool, error) {
	token, err := j.ParseJWT(ctx, tokenString)
	if err != nil {
		return false, errors.Wrap(err, "error parsing token")
	}

	return token.Valid, nil
}

func (j *JWT) GetClaims(ctx context.Context, tokenString string) (*AuthClaims, error) {
	token, err := j.ParseJWT(ctx, tokenString)
	if err != nil {
		return nil, errors.Wrap(err, "error parsing token")
	}
	claims, ok := token.Claims.(*AuthClaims)
	if !ok || !token.Valid {
		return nil, errors.Wrap(err, "Token not valid")
	}
	return claims, nil
}

func (j *JWT) ParseJWT(ctx context.Context, tokenString string) (*jwt.Token, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		var key any = j.publicKey
		if key == nil {
			key = j.privateKey
		}

		if key == nil {
			return nil, errors.Wrap(ErrKeyMissing, "Error validate jwt keys is nil")
		}

		return key, nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "Error validate jwt token")
	}
	return token, nil
}

// GenerateAllTokensAsync TODO: а надо ли ? может возвращать структуру с access и refresh токеном|
//
//	return array tokens first elem is access token and second if refresh token
func (j *JWT) GenerateAllTokensAsync(ctx context.Context, uid, role, deviceID string) ([]string, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	refreshCh, accessCh := make(chan string, 1), make(chan string, 1)
	errCh := make(chan error, 1)
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		access, err := j.GenerateAccessToken(ctx, uid, role, deviceID)
		if err != nil {
			errCh <- errors.Wrap(err, "Error create access token")
			return
		}
		accessCh <- access
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		refresh, err := j.GenerateRefreshToken(ctx, uid, role, deviceID)
		if err != nil {
			errCh <- errors.Wrap(err, "Error create refresh token")
			return
		}
		refreshCh <- refresh
	}()

	go func() {
		wg.Wait()
		close(errCh)
		close(accessCh)
		close(refreshCh)
	}()

	for err := range errCh {
		if err != nil {
			return nil, err
		}
	}
	return []string{<-accessCh, <-refreshCh}, nil
}

// GenerateAllTokens TODO: а надо ли ? может возвращать структуру с access и refresh токеном|
//
// return array tokens first elem is access token and second if refresh token
func (j *JWT) GenerateAllTokens(ctx context.Context, uid, role, deviceID string) ([]string, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	access, err := j.GenerateAccessToken(ctx, uid, role, deviceID)
	if err != nil {
		return nil, errors.Wrap(err, "Error create access token")
	}

	refresh, err := j.GenerateRefreshToken(ctx, uid, role, deviceID)
	if err != nil {
		return nil, errors.Wrap(err, "Error create refresh token")
	}

	return []string{access, refresh}, nil
}

func (j *JWT) GenerateAccessToken(ctx context.Context, uid, role, deviceID string) (string, error) {
	if ctx.Err() != nil {
		return "", ctx.Err()
	}

	claims := AuthClaims{
		DeviceID: deviceID,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    j.issuer,
			Subject:   uid,
			Audience:  j.audience,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.accessTokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        j.generateJTI(),
		},
	}

	signedToken, err := j.generateJWT(&claims)

	return signedToken, err
}

func (j *JWT) GenerateRefreshToken(ctx context.Context, uid, role, deviceID string) (string, error) {
	if ctx.Err() != nil {
		return "", ctx.Err()
	}

	claims := AuthClaims{
		Role:     role,
		DeviceID: deviceID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    j.issuer,
			Subject:   uid,
			Audience:  j.audience,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.refreshTokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        j.generateJTI(),
		},
	}

	signedToken, err := j.generateJWT(&claims)

	return signedToken, err
}

func (j *JWT) generateJWT(claims *AuthClaims) (string, error) {
	if j.privateKey == nil {
		return "", errors.Wrap(ErrKeyMissing, "Error generate jwt, keys is nil")
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	signedToken, err := token.SignedString(j.privateKey)
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func (j *JWT) generateJTI() string {
	return uuid.New().String()
}
