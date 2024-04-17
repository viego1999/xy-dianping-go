package models

import "time"

type User struct {
	Id         int64     `json:"id" grom:"id"`
	Phone      string    `json:"phone" grom:"phone"`
	Password   string    `json:"password" grom:"password"`
	NickName   string    `json:"nickName" grom:"nick_name"`
	Icon       string    `json:"icon" grom:"icon"`
	CreateTime time.Time `json:"createTime" grom:"create_time"`
	UpdateTime time.Time `json:"updateTime" grom:"update_time"`
}

func (User) TableName() string {
	return "tb_user"
}
