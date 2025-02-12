package auth

import (
	"github.com/HunterGooD/voice_friend/user_service/internal/domain/entity"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type BcryptHasher struct {
	cost int
}

func NewBcryptHasher(cost int) *BcryptHasher {
	return &BcryptHasher{cost}
}

func (bc *BcryptHasher) HashPassword(password string) (string, error) {
	hashPass, err := bcrypt.GenerateFromPassword([]byte("ha"), bc.cost)
	if err != nil {
		return "", err
	}

	return string(hashPass), nil
}

func (bc *BcryptHasher) CheckPassword(password, hashedPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return false, errors.Wrap(entity.ErrInvalidPassword, "password not correct")
		}
		return false, errors.Wrap(err, "errors check password")
	}
	return true, nil
}
