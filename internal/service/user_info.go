package service

import (
	"xy-dianping-go/internal/models"
	"xy-dianping-go/internal/repo"
)

type UserInfoService interface {
	GetUserInfoByUserId(userId int64) (*models.UserInfo, error)
}

type UserInfoServiceImpl struct {
	userInfoRepo repo.UserInfoRepository
}

func NewUserInfoService(userInfoRepo repo.UserInfoRepository) UserInfoService {
	return &UserInfoServiceImpl{userInfoRepo: userInfoRepo}
}

func (s *UserInfoServiceImpl) GetUserInfoByUserId(userId int64) (*models.UserInfo, error) {

	return s.userInfoRepo.QueryByUserId(userId)
}
