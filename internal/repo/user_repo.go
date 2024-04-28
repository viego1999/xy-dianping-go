package repo

import (
	"errors"
	"gorm.io/gorm"
	"trpc.group/trpc-go/trpc-go/log"
	"xy-dianping-go/internal/models"
)

type UserRepository interface {
	QueryById(id int64) (*models.User, error)
	QueryByIds(ids []int64) ([]models.User, error)
	QueryByPhone(phone string) (*models.User, error)
	CreateUser(user *models.User) error
	ListByIds(sql string, ids []int64, idsStr string) ([]models.User, error)
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

func (r *UserRepositoryImpl) QueryByIds(ids []int64) ([]models.User, error) {
	var users []models.User
	err := r.Db.Where("id IN ?", ids).Find(&users).Error
	return users, err
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

// ListByIds 根据用户的 id 列表批量查询用户记录
func (r *UserRepositoryImpl) ListByIds(sql string, ids []int64, idsStr string) ([]models.User, error) {
	var users []models.User
	err := r.Db.Raw(sql, ids, idsStr).Scan(&users).Error
	return users, err
}
