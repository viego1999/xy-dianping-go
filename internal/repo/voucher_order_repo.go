package repo

import (
	"gorm.io/gorm"
	"xy-dianping-go/internal/models"
)

type VoucherOrderRepository interface {
	QueryOrderByQuery(query interface{}, args ...interface{}) (*models.VoucherOrder, error)
	CreateVoucherOrder(order *models.VoucherOrder) error
	ExecuteTransaction(fn func(txRepo VoucherOrderRepository) error) error
}

type VoucherOrderRepositoryImpl struct {
	Db *gorm.DB
}

func NewVoucherOrderRepository(db *gorm.DB) VoucherOrderRepository {
	return &VoucherOrderRepositoryImpl{db}
}

func (r *VoucherOrderRepositoryImpl) QueryOrderByQuery(query interface{}, args ...interface{}) (*models.VoucherOrder, error) {
	var order models.VoucherOrder
	err := r.Db.Where(query, args...).First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, err
}

func (r *VoucherOrderRepositoryImpl) CreateVoucherOrder(order *models.VoucherOrder) error {
	return r.Db.Create(&order).Error
}

func (r *VoucherOrderRepositoryImpl) ExecuteTransaction(fn func(txRepo VoucherOrderRepository) error) error {
	return r.Db.Transaction(func(tx *gorm.DB) error {
		txVoucherOrderRepo := &VoucherOrderRepositoryImpl{tx}
		return fn(txVoucherOrderRepo)
	})
}
