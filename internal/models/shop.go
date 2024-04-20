package models

import "time"

type Shop struct {
	Id         int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	Name       string    `json:"name" gorm:"column:name"`
	TypeId     int64     `json:"typeId" gorm:"column:type_id"`
	Images     string    `json:"images" gorm:"column:images"`
	Area       string    `json:"area" gorm:"column:area"`
	Address    string    `json:"address" gorm:"column:address"`
	X          float64   `json:"x" gorm:"column:x"`
	Y          float64   `json:"y" gorm:"column:y"`
	AvgPrice   int64     `json:"avgPrice" gorm:"column:avg_price"`
	Sold       int       `json:"sold" gorm:"column:sold"`
	Comments   int       `json:"comments" gorm:"column:comments"`
	Score      int       `json:"score" gorm:"column:score"`
	OpenHours  string    `json:"openHours" gorm:"column:open_hours"`
	CreateTime time.Time `json:"createTime" gorm:"column:create_time;autoCreateTime"` // 注意这里使用了column标签指定字段名
	UpdateTime time.Time `json:"updateTime" gorm:"column:update_time;autoUpdateTime"` // 注意这里使用了column标签指定字段名
	Distance   float64   `json:"distance" gorm:"-"`
}

func (Shop) TableName() string {
	return "tb_shop"
}
