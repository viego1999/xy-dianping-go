package repo

import (
	"errors"

	"gorm.io/gorm"
	"trpc.group/trpc-go/trpc-go/log"
	"xy-dianping-go/internal/models"
)

type UserRepository interface {
	QueryById(id int64) (*models.User, error)
	QueryByPhone(phone string) (*models.User, error)
	CreateUser(user *models.User) error
}

type UserRepositoryImpl struct {
	Db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &UserRepositoryImpl{Db: db}
}

func (r *UserRepositoryImpl) QueryById(id int64) (*models.User, error) {
	var user models.User
	err := r.Db.Where("id = ?", id).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Info("QueryById: record not found.")
		} else {
			log.Error("QueryById error:", err)
		}
		return nil, err
	}
	return &user, err
}

func (r *UserRepositoryImpl) QueryByPhone(phone string) (*models.User, error) {
	var user models.User
	err := r.Db.Where("phone = ?", phone).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Info("QueryByPhone: record not found.")
		} else {
			log.Error("QueryByPhone error:", err)
		}
		return nil, err
	}
	return &user, err
}

func (r *UserRepositoryImpl) CreateUser(user *models.User) error {
	return r.Db.Create(user).Error
}
