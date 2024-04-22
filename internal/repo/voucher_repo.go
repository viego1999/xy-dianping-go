package repo

import (
	"gorm.io/gorm"
	"xy-dianping-go/internal/models"
)

type VoucherRepository interface {
	CreateVoucher(voucher *models.Voucher) error
	QueryVoucherByShopId(shopId int64) ([]models.Voucher, error)
}

type VoucherRepositoryImpl struct {
	Db *gorm.DB
}

func NewVoucherRepository(db *gorm.DB) VoucherRepository {
	return &VoucherRepositoryImpl{db}
}

func (r *VoucherRepositoryImpl) CreateVoucher(voucher *models.Voucher) error {
	return r.Db.Create(voucher).Error
}

func (r *VoucherRepositoryImpl) QueryVoucherByShopId(shopId int64) ([]models.Voucher, error) {
	var vouchers []models.Voucher
	err := r.Db.Where("shop_id = ?", shopId).Find(&vouchers).Error
	if err != nil {
		return nil, err
	}

	return vouchers, err
}
