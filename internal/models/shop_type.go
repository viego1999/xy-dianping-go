package models

import "time"

type ShopType struct {
	Id         int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	Name       string    `json:"name" gorm:"column:name"`
	Icon       string    `json:"icon" gorm:"column:icon"`
	Sort       int       `json:"sort" gorm:"column:sort"`
	CreateTime time.Time `json:"createTime" gorm:"column:create_time;autoCreateTime"` // 注意这里使用了column标签指定字段名
	UpdateTime time.Time `json:"updateTime" gorm:"column:update_time;autoUpdateTime"` // 注意这里使用了column标签指定字段名
}

func (ShopType) TableName() string {
	return "tb_shop_type"
}
