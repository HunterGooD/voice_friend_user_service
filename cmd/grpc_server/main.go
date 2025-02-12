package main

import (
	"github.com/HunterGooD/voice_friend_user_service/internal/adapter"
	auth2 "github.com/HunterGooD/voice_friend_user_service/pkg/auth"
	"os"
	"time"

	"github.com/HunterGooD/voice_friend_user_service/config"
	"github.com/HunterGooD/voice_friend_user_service/internal/handler"
	"github.com/HunterGooD/voice_friend_user_service/internal/repository"
	"github.com/HunterGooD/voice_friend_user_service/internal/usecase"
	"github.com/HunterGooD/voice_friend_user_service/pkg/database"
	"github.com/HunterGooD/voice_friend_user_service/pkg/logger"
	"github.com/HunterGooD/voice_friend_user_service/pkg/server"
)

func main() {
	// load config
	log := logger.NewJsonLogrusLogger(os.Stdout, os.Getenv("LOG_LEVEL"))

	configPath := os.Getenv("CONFIG_PATH")
	cfg, err := config.NewConfig(configPath)
	if err != nil {
		log.Error("Error init config", err)
		panic(err)
	}

	db, err := database.NewPostgresConnection(
		cfg.BuildDSN(),
		cfg.Database.PoolConnection.MaxOpenConns,
		cfg.Database.PoolConnection.MaxIdleConns,
		cfg.Database.PoolConnection.MaxLifeTime,
	)
	if err != nil {
		log.Error("Error init postgresql database", err)
		panic(err)
	}

	conn, err := database.NewRedisConn(cfg.GetRedisAddr(), cfg.Redis.User, cfg.Redis.Password, cfg.Redis.DBIdx)
	if err != nil {
		log.Error("Error init redis", err)
		panic(err)
	}

	userRepository := repository.NewUserRepository(db)
	tokenRepository := repository.NewTokenRepository(conn)

	tokenManager := auth2.NewJWTGenerator(
		nil,
		nil,
		cfg.JWT.Issuer,
		cfg.JWT.AccessTokenDuration,
		cfg.GetRefreshTokenTime(),
		[]string{""},
	)

	_, err = tokenManager.LoadPrivateKeyFromFile(cfg.App.CertFilePath)
	if err != nil {
		log.Error("Error load private key for token manager", err)
		panic(err)
	}

	tokenAdapter := adapter.NewTokenManagerAdapter(tokenManager)

	hasher := auth2.NewArgon2Hasher(
		cfg.Argon2.Times,
		cfg.Argon2.Memory*1024,
		cfg.Argon2.KeyLen,
		cfg.Argon2.SaltLen,
		cfg.Argon2.Threads,
	)

	authUsecase := usecase.NewAuthUsecase(userRepository, tokenRepository, tokenAdapter, hasher)
	userProfileUsecase := usecase.NewUserProfileUsecase(userRepository, tokenAdapter, log)

	// init gRPC server
	gRPCServer := server.NewGRPCServer(log, 5, time.Duration(30)*time.Second)

	// register handlers
	handler.NewAuthHandler(gRPCServer, authUsecase, log)
	handler.NewUserProfileHandler(gRPCServer, userProfileUsecase, log)

	if err := gRPCServer.Start(cfg.GetAddress()); err != nil {
		panic(err)
	}
}
