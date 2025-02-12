package entity

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Role string

const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
	RoleGuest Role = "guest"
)

func (r Role) IsValid() bool {
	switch r {
	case RoleAdmin, RoleUser, RoleGuest:
		return true
	}
	return false
}

type User struct {
	ID             int64           `db:"id"`
	Login          string          `db:"login"`
	Name           string          `db:"name"`
	UID            uuid.UUID       `db:"uid"`
	Email          *string         `db:"email"`
	Password       string          `db:"password"`
	IsActive       bool            `db:"is_active"`
	LastLogin      *time.Time      `db:"last_login"`
	Role           Role            `db:"role"`
	ProfilePicture *string         `db:"profile_picture"`
	Phone          *string         `db:"phone"`
	Metadata       json.RawMessage `db:"metadata"`
	CreatedAt      time.Time       `db:"created_at"`
	UpdatedAt      time.Time       `db:"updated_at"`
	DeletedAt      *time.Time      `db:"deleted_at"`
}

type AuthUserResponse struct {
	AccessToken  string
	RefreshToken string
}
