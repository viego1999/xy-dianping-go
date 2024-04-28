package repo

import (
	"gorm.io/gorm"
	"xy-dianping-go/internal/models"
)

type FollowRepository interface {
	QueryFollows(query string, ids ...interface{}) ([]models.Follow, error)
	CreateFollow(follow *models.Follow) error
	DeleteFollow(query string, ids ...interface{}) (int64, error)
}

type FollowRepositoryImpl struct {
	Db *gorm.DB
}

func NewFollowRepository(db *gorm.DB) FollowRepository {
	return &FollowRepositoryImpl{db}
}

// QueryFollows 查询关注列表通过 user_id 和 follow_user_id，
//
// 例如： follows, err := QueryFollows("user_id = ? AND follow_user_id = ?", userId, follow_user_id)
func (r *FollowRepositoryImpl) QueryFollows(query string, ids ...interface{}) ([]models.Follow, error) {
	var follows []models.Follow
	err := r.Db.Where(query, ids).Find(&follows).Error
	return follows, err
}

func (r *FollowRepositoryImpl) CreateFollow(follow *models.Follow) error {
	return r.Db.Create(follow).Error
}

// DeleteFollow 删除关注列表通过 user_id 和 follow_user_id，
//
// 例如： rows, err := DeleteFollow("user_id = ? AND follow_user_id = ?", userId, follow_user_id)
func (r *FollowRepositoryImpl) DeleteFollow(query string, ids ...interface{}) (int64, error) {
	result := r.Db.Where(query, ids).Delete(&models.Follow{})
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}
