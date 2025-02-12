package entity

import (
	"github.com/pkg/errors"
)

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrNotFound          = errors.New("not found")
	ErrUnauthorized      = errors.New("unauthorization user")
	ErrDataNotValid      = errors.New("not valid user data")
	ErrInvalidPassword   = errors.New("invalid password")
)
