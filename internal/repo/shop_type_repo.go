package repo

import (
	"gorm.io/gorm"
	"xy-dianping-go/internal/models"
)

type ShopTypeRepository interface {
	List() ([]models.ShopType, error)
}

type ShopTypeRepositoryImpl struct {
	Db *gorm.DB
}

func NewShopTypeRepository(db *gorm.DB) ShopTypeRepository {
	return &ShopTypeRepositoryImpl{db}
}

func (r *ShopTypeRepositoryImpl) List() ([]models.ShopType, error) {
	var shopTypes []models.ShopType
	err := r.Db.Order("sort ASC").Find(&shopTypes).Error
	return shopTypes, err
}
