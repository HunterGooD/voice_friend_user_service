package repository

import (
	"context"
	"database/sql"
	"github.com/HunterGooD/voice_friend/user_service/internal/domain/entity"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type UserRepository struct {
	db *sqlx.DB
	tx *sqlx.Tx
}

func NewUserRepository(conn *sqlx.DB) *UserRepository {
	return &UserRepository{
		db: conn,
	}
}

func (r *UserRepository) WithTransaction(tx *sqlx.Tx) *UserRepository {
	return &UserRepository{r.db, tx}
}

func (r *UserRepository) AddUser(ctx context.Context, user *entity.User) error {
	query := `INSERT INTO users (
			login, name, email, password, profile_picture, phone
		) VALUES (
			:login, :name, :email, :password, :profile_picture, :phone
		)
		RETURNING id, uid, created_at, updated_at;`
	rows, err := r.db.NamedQueryContext(ctx, query, user)
	if err != nil {
		return errors.Wrap(err, "Error execute insert query for new user")
	}
	defer rows.Close()

	if rows.Next() {
		err = rows.StructScan(&user)
		if err != nil {
			return errors.Wrap(err, "Failed to scan result")
		}
	}
	return nil
}

func (r *UserRepository) ExistUser(ctx context.Context, login string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE login = $1)`
	err := r.db.GetContext(ctx, &exists, query, login)
	if err != nil {
		return false, errors.Wrap(err, "Error check exists user")
	}
	return exists, nil
}

func (r *UserRepository) GetUserPasswordByLogin(ctx context.Context, login string) (string, error) {
	var password string
	query := `SELECT password FROM users WHERE login = $1`
	err := r.db.GetContext(ctx, &password, query, login)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", errors.Wrap(entity.ErrNotFound, "Error get user password")
		}
		return "", errors.Wrap(err, "Unknown error get user password")
	}
	return password, nil
}
