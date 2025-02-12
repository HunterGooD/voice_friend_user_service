package database

import (
	"context"
	"github.com/redis/go-redis/v9"
)

func NewRedisConn(addr, username, password string, dbIdx int) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Username: username,
		Password: password,
		DB:       dbIdx,
	})
	_, err := client.Ping(context.Background()).Result()
	return client, err
}
