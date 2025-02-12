package database

import (
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Подключение PostgreSQL драйвера
)

func NewPostgresConnection(dsn string, maxOpenConns, maxIdleConns int, maxLifeTime time.Duration) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, err
	}

	// settings pool connections
	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxLifetime(maxLifeTime)

	return db, nil
}
