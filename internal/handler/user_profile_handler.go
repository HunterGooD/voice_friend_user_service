package handler

import (
	"context"

	"github.com/HunterGooD/voice_friend/user_service/pkg/logger"
	pd "github.com/HunterGooD/voice_friend_contracts/gen/go/user_service"
)

type UserProfileUsecase interface {
	ChangeAvatar(ctx context.Context)
}

type UserProfileHandler struct {
	pd.UnimplementedUserProfileServer
	uu UserProfileUsecase

	log logger.Logger
}

func NewUserProfileHandler(gRPCServer GRPCServer, uu UserProfileUsecase, log logger.Logger) {
	userProfileHandler := &UserProfileHandler{uu: uu, log: log}
	pd.RegisterUserProfileServer(gRPCServer.GetServer(), userProfileHandler)
}

func (up *UserProfileHandler) ChangeAvatar(ctx context.Context, req *pd.AvatarChangeRequest) (*pd.AvatarChangeResponse, error) {
	return nil, nil
}
