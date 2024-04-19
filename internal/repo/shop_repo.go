package repo

import (
	"errors"
	"gorm.io/gorm"
	"trpc.group/trpc-go/trpc-go/log"
	"xy-dianping-go/internal/models"
)

type ShopRepository interface {
	QueryById(id int64) (*models.Shop, error)
}

type ShopRepositoryImpl struct {
	Db *gorm.DB
}

func NewShopRepository(db *gorm.DB) ShopRepository {
	return &ShopRepositoryImpl{Db: db}
}

func (r *ShopRepositoryImpl) QueryById(id int64) (*models.Shop, error) {
	var shop models.Shop
	err := r.Db.Where("id = ?", id).First(&shop).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Info("QueryById: record not found.")
		} else {
			log.Error("QueryById error:", err)
		}
		return nil, err
	}
	return &shop, err
}
