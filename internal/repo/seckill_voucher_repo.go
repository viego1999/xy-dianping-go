package repo

import (
	"gorm.io/gorm"
	"xy-dianping-go/internal/models"
)

type SeckillVoucherRepository interface {
	CreateSeckillVoucher(seckillVoucher *models.SeckillVoucher) error
	UpdateSeckillVoucher(query interface{}, arg interface{}, column string, value interface{}) (int64, error)
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

func (r *SeckillVoucherRepositoryImpl) UpdateSeckillVoucher(query interface{}, arg interface{}, column string, value interface{}) (int64, error) {
	result := r.Db.Model(&models.SeckillVoucher{}).Where(query, arg).Update(column, value)
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}
