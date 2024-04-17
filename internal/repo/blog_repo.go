package repo

import (
	"errors"

	"gorm.io/gorm"
	"trpc.group/trpc-go/trpc-go/log"
	"xy-dianping-go/internal/models"
)

// BlogRepository 定义了 models.Blog 模型的存储库接口
type BlogRepository interface {
	QueryById(id int64) (*models.Blog, error)
	UpdateBySql(sql string) (int, error)
}

type BlogRepositoryImpl struct {
	Db *gorm.DB
}

func NewBlogRepository(db *gorm.DB) BlogRepository {
	return &BlogRepositoryImpl{Db: db}
}

func (r *BlogRepositoryImpl) QueryById(id int64) (*models.Blog, error) {
	var blog models.Blog
	err := r.Db.Where("id = ?", id).First(&blog).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Info("record not found.")
		} else {
			log.Error("QueryBlogById error:", err)
		}
		return nil, err
	}
	return &blog, err
}

func (r *BlogRepositoryImpl) UpdateBySql(sql string) (int, error) {
	return 0, nil
}
