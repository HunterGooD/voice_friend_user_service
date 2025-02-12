package usecase

import (
	"context"

	"github.com/HunterGooD/voice_friend/user_service/pkg/logger"
)

type UserProfileUsecase struct {
	ur UserRepository
	tm TokenManager

	log logger.Logger
}

func NewUserProfileUsecase(ur UserRepository, tm TokenManager, log logger.Logger) *UserProfileUsecase {
	return &UserProfileUsecase{ur, tm, log}
}

func (uu *UserProfileUsecase) ChangeAvatar(ctx context.Context) {

}
