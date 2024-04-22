package repo

import (
	"gorm.io/gorm"
	"xy-dianping-go/internal/models"
)

type SeckillVoucherRepository interface {
	CreateSeckillVoucher(seckillVoucher *models.SeckillVoucher) error
}

type SeckillVoucherRepositoryImpl struct {
	Db *gorm.DB
}

func NewSeckillVoucherRepository(Db *gorm.DB) SeckillVoucherRepository {
	return &SeckillVoucherRepositoryImpl{Db}
}

func (r *SeckillVoucherRepositoryImpl) CreateSeckillVoucher(seckillVoucher *models.SeckillVoucher) error {
	return r.Db.Create(seckillVoucher).Error
}
