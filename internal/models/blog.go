package models

import "time"

type Blog struct {
	Id         int64     `json:"id" grom:"column:id"`
	ShopId     int64     `json:"shopId" grom:"column:shop_id"`
	UserId     int64     `json:"userId" grom:"column:user_id"`
	Icon       string    `json:"icon" grom:"-"`
	Name       string    `json:"name" grom:"-"`
	IsLike     bool      `json:"isLike" grom:"-"`
	Title      string    `json:"title" grom:"column:title"`
	Images     string    `json:"images" grom:"column:images"`
	Content    string    `json:"content" grom:"column:content"`
	Liked      int       `json:"liked" grom:"column:liked"`
	Comments   int       `json:"comments" grom:"column:comments"`
	CreateTime time.Time `json:"createTime" grom:"column:create_time"`
	UpdateTime time.Time `json:"updateTime" grom:"column:update_time"`
}

func (Blog) TableName() string {
	return "tb_blog"
}
