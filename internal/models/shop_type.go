package models

import "time"

type ShopType struct {
	Id         int64     `json:"id" grom:"primaryKey;autoIncrement"`
	Name       string    `json:"name" grom:"column:name"`
	Icon       string    `json:"icon" grom:"column:icon"`
	Sort       int       `json:"sort" grom:"column:sort"`
	CreateTime time.Time `json:"createTime" grom:"column:create_time"`
	UpdateTime time.Time `json:"updateTime" grom:"column:update_time"`
}

func (ShopType) TableName() string {
	return "tb_shop_type"
}
