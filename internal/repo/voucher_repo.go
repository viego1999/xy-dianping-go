package repo

import (
	"gorm.io/gorm"
	"xy-dianping-go/internal/models"
)

type VoucherRepository interface {
	CreateVoucher(voucher *models.Voucher) error
	QueryVoucherByShopId(shopId int64) ([]models.Voucher, error)
	// ExecuteTransaction 开启事务操作，返回结果为 nil 时，提交并执行事务，否则进行回滚
	ExecuteTransaction(fn func(txRepo VoucherRepository) error) error
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

func (r *VoucherRepositoryImpl) ExecuteTransaction(fn func(txRepo VoucherRepository) error) error {
	return r.Db.Transaction(func(tx *gorm.DB) error {
		txVoucherRepo := &VoucherRepositoryImpl{tx}
		return fn(txVoucherRepo)
	})
}
