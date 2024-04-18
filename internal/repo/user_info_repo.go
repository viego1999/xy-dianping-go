package repo

import (
	"errors"
	"gorm.io/gorm"
	"trpc.group/trpc-go/trpc-go/log"
	"xy-dianping-go/internal/models"
)

type UserInfoRepository interface {
	QueryByUserId(userId int64) (*models.UserInfo, error)
}

type UserInfoRepositoryImpl struct {
	Db *gorm.DB
}

func NewUserInfoRepository(db *gorm.DB) UserInfoRepository {
	return &UserRepositoryImpl{Db: db}
}

func (r *UserRepositoryImpl) QueryByUserId(userId int64) (*models.UserInfo, error) {
	var userInfo models.UserInfo
	err := r.Db.Where("user_id = ?", userId).First(&userInfo).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Info("record not found.")
		} else {
			log.Error("QueryByUserId error:", err)
		}
		return nil, err
	}
	return &userInfo, err
}
