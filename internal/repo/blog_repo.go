package repo

import (
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"trpc.group/trpc-go/trpc-go/log"
	"xy-dianping-go/internal/constants"
	"xy-dianping-go/internal/models"
)

// BlogRepository 定义了 models.Blog 模型的存储库接口
type BlogRepository interface {
	QueryById(id int64) (*models.Blog, error)
	CreateBlog(blog *models.Blog) error
	UpdateById(id int64, column string, expr clause.Expr) (int64, error)
	QueryByUserId(userId int64, current int) ([]models.Blog, error)
	PageQuery(current int) ([]models.Blog, error)
	ListByIds(sql string, ids []int64, idsStr string) ([]models.Blog, error)
}

type BlogRepositoryImpl struct {
	Db *gorm.DB
}

func NewBlogRepository(db *gorm.DB) BlogRepository {
	return &BlogRepositoryImpl{Db: db}
}

func (r *BlogRepositoryImpl) CreateBlog(blog *models.Blog) error {
	return r.Db.Create(blog).Error
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

func (r *BlogRepositoryImpl) UpdateById(id int64, column string, expr clause.Expr) (int64, error) {
	result := r.Db.Model(&models.Blog{}).Where("id = ?", id).Update(column, expr)
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}

func (r *BlogRepositoryImpl) QueryByUserId(userId int64, current int) ([]models.Blog, error) {
	var blogs []models.Blog
	err := r.Db.Where("user_id = ?", userId).
		Offset((current - 1) * constants.MAX_PAGE_SIZE).
		Limit(constants.MAX_PAGE_SIZE).Find(&blogs).Error
	return blogs, err
}

func (r *BlogRepositoryImpl) PageQuery(current int) ([]models.Blog, error) {
	var blogs []models.Blog
	err := r.Db.Order("liked DESC").
		Offset((current - 1) * constants.MAX_PAGE_SIZE).
		Limit(constants.MAX_PAGE_SIZE).Find(&blogs).Error
	return blogs, err
}

// ListByIds 根据博客的 id 列表批量查询博客记录
func (r *BlogRepositoryImpl) ListByIds(sql string, ids []int64, idsStr string) ([]models.Blog, error) {
	var blogs []models.Blog
	err := r.Db.Raw(sql, ids, idsStr).Scan(&blogs).Error
	return blogs, err
}
