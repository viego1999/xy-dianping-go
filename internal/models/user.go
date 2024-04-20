package models

import "time"

type User struct {
	Id         int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	Phone      string    `json:"phone" gorm:"phone"`
	Password   string    `json:"password" gorm:"password"`
	NickName   string    `json:"nickName" gorm:"nick_name"`
	Icon       string    `json:"icon" gorm:"icon"`
	CreateTime time.Time `json:"createTime" gorm:"column:create_time;autoCreateTime"` // 注意这里使用了column标签指定字段名
	UpdateTime time.Time `json:"updateTime" gorm:"column:update_time;autoUpdateTime"` // 注意这里使用了column标签指定字段名
}

func (User) TableName() string {
	return "tb_user"
}
