package models

import "time"

type UserInfo struct {
	UserId     int64     `json:"userId" gorm:"column:user_id;primaryKey;autoIncrement"`
	City       string    `json:"city" gorm:"column:city"`
	Introduce  string    `json:"introduce" gorm:"column:introduce"`
	Fans       int       `json:"fans" gorm:"column:fans"`
	Followee   int       `json:"followee" gorm:"column:followee"`
	Gender     bool      `json:"gender" gorm:"column:gender"`
	Birthday   time.Time `json:"birthday" gorm:"column:birthday"`
	Credits    int       `json:"credits" gorm:"column:credits"`
	Level      bool      `json:"level" gorm:"column:level"`
	CreateTime time.Time `json:"createTime" gorm:"column:create_time"`
	UpdateTime time.Time `json:"updateTime" gorm:"column:update_time"`
}

func (UserInfo) TableName() string {
	return "tb_user_info"
}
