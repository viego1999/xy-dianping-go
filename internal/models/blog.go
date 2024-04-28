package models

import "time"

type Blog struct {
	Id         int64     `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	ShopId     int64     `json:"shopId" gorm:"column:shop_id"`
	UserId     int64     `json:"userId" gorm:"column:user_id"`
	Icon       string    `json:"icon" gorm:"-"`
	Name       string    `json:"name" gorm:"-"`
	IsLike     bool      `json:"isLike" gorm:"-"`
	Title      string    `json:"title" gorm:"column:title"`
	Images     string    `json:"images" gorm:"column:images"`
	Content    string    `json:"content" gorm:"column:content"`
	Liked      int       `json:"liked" gorm:"column:liked"`
	Comments   int       `json:"comments" gorm:"column:comments"`
	CreateTime time.Time `json:"createTime" gorm:"column:create_time;autoCreateTime"` // 注意这里使用了column标签指定字段名
	UpdateTime time.Time `json:"updateTime" gorm:"column:update_time;autoUpdateTime"` // 注意这里使用了column标签指定字段名
}

func (*Blog) TableName() string {
	return "tb_blog"
}
