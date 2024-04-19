package models

import "time"

type Shop struct {
	Id         int64     `json:"id" grom:"primaryKey;autoIncrement"`
	Name       string    `json:"name" grom:"column:name"`
	TypeId     int64     `json:"typeId" grom:"column:type_id"`
	Images     string    `json:"images" grom:"column:images"`
	Area       string    `json:"area" grom:"column:area"`
	Address    string    `json:"address" grom:"column:address"`
	X          float64   `json:"x" grom:"column:x"`
	Y          float64   `json:"y" grom:"column:y"`
	AvgPrice   int64     `json:"avgPrice" grom:"column:avg_price"`
	Sold       int       `json:"sold" grom:"column:sold"`
	Comments   int       `json:"comments" grom:"column:comments"`
	Score      int       `json:"score" grom:"column:score"`
	OpenHours  string    `json:"openHours" grom:"column:open_hours"`
	CreateTime time.Time `json:"createTime" grom:"column:create_time"`
	UpdateTime time.Time `json:"updateTime" grom:"column:update_time"`
}

func (Shop) TableName() string {
	return "tb_shop"
}
