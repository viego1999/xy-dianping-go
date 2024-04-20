package repo

import (
	"errors"
	"gorm.io/gorm"
	"trpc.group/trpc-go/trpc-go/log"
	"xy-dianping-go/internal/constants"
	"xy-dianping-go/internal/models"
)

type ShopRepository interface {
	QueryById(id int64) (*models.Shop, error)
	QueryByIds(ids []int64, orderClause string) ([]models.Shop, error)
	List() ([]models.Shop, error)
	CreateShop(shop *models.Shop) error
	Update(shop *models.Shop) error
	UpdateColumns(columns []string, shop *models.Shop) error
	UpdateByMap(id int64, shopMap map[string]interface{}) error
	QueryByTypeId(typeId int64, current int) ([]models.Shop, error)
	QueryByName(name string, current int) ([]models.Shop, error)
	ExecuteTransaction(func(repo ShopRepository) error) error
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

func (r *ShopRepositoryImpl) QueryByIds(ids []int64, orderClause string) ([]models.Shop, error) {
	var shops []models.Shop
	err := r.Db.Where("id IN ?", ids).Order(orderClause).Find(&shops).Error
	return shops, err
}

func (r *ShopRepositoryImpl) List() ([]models.Shop, error) {
	var shops []models.Shop
	err := r.Db.Find(&shops).Error
	return shops, err
}

func (r *ShopRepositoryImpl) CreateShop(shop *models.Shop) error {
	return r.Db.Create(shop).Error
}

func (r *ShopRepositoryImpl) Update(shop *models.Shop) error {
	// Updates(model) 不会更新 model 中的零值
	return r.Db.Model(&models.Shop{}).Where("id = ?", shop.Id).Updates(shop).Error
}

func (r *ShopRepositoryImpl) UpdateColumns(columns []string, shop *models.Shop) error {
	// 更新 Select 中指定的 columns
	return r.Db.Model(&models.Shop{}).Where("id = ?", shop.Id).Select(columns).Updates(shop).Error
}

func (r *ShopRepositoryImpl) UpdateByMap(id int64, shopMap map[string]interface{}) error {
	// Updates(map) 更新 map 中所有的字段
	return r.Db.Model(&models.Shop{}).Where("id = ?", id).Updates(shopMap).Error
}

func (r *ShopRepositoryImpl) QueryByTypeId(typeId int64, current int) ([]models.Shop, error) {
	var shops []models.Shop

	result := r.Db.Where("type_id = ?", typeId).
		Offset((current - 1) * constants.DEFAULT_PAGE_SIZE).
		Limit(constants.DEFAULT_PAGE_SIZE).
		Find(&shops)

	// 检查错误
	if result.Error != nil {
		return nil, result.Error
	}

	// 否则返回查询到的记录
	return shops, nil
}

func (r *ShopRepositoryImpl) QueryByName(name string, current int) ([]models.Shop, error) {
	var shops []models.Shop

	result := r.Db.Where("name LIKE ?", "%"+name+"%").
		Offset((current - 1) * constants.DEFAULT_PAGE_SIZE).
		Limit(constants.DEFAULT_PAGE_SIZE).
		Find(&shops)

	// 检查错误
	if result.Error != nil {
		return nil, result.Error
	}

	// 否则返回查询到的记录
	return shops, nil
}

func (r *ShopRepositoryImpl) ExecuteTransaction(fn func(txShopRepo ShopRepository) error) error {
	// 开启事务
	return r.Db.Transaction(func(tx *gorm.DB) error {
		txRepo := &ShopRepositoryImpl{tx} // 创建事务绑定的 repository 实例
		return fn(txRepo)                 // 执行传入的函数，如果返回错误则事务进行回滚
	})
}
