package models

import "time"

type Follow struct {
	Id           int64     `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	UserId       int64     `json:"userId" gorm:"column:user_id"`
	FollowUserId int64     `json:"followUserId" gorm:"column:follow_user_id"`
	CreateTime   time.Time `json:"createTime" gorm:"column:create_time;autoCreateTime"`
}

func (*Follow) TableName() string {
	return "tb_follow"
}
