package models

import (
	"encoding/json"
	"time"
)

type Voucher struct {
	Id          int64     `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	ShopId      int64     `json:"shopId" gorm:"column:shop_id"`
	Title       string    `json:"title" gorm:"column:title"`
	SubTitle    string    `json:"subTitle" gorm:"column:sub_title"`
	Rules       string    `json:"rules" gorm:"column:rules"`
	PayValue    int64     `json:"payValue" gorm:"column:pay_value"`
	ActualValue int64     `json:"actualValue" gorm:"column:actual_value"`
	Type        int       `json:"type" gorm:"column:type"`
	Status      int       `json:"status" gorm:"column:status"`
	Stock       int       `json:"stock" gorm:"-"`
	BeginTime   time.Time `json:"beginTime" gorm:"-"`
	EndTime     time.Time `json:"endTime" gorm:"-"`
	CreateTime  time.Time `json:"createTime" gorm:"column:create_time;autoCreateTime"` // 注意这里使用了column标签指定字段名
	UpdateTime  time.Time `json:"updateTime" gorm:"column:update_time;autoUpdateTime"` // 注意这里使用了column标签指定字段名
}

func (*Voucher) TableName() string {
	return "tb_voucher"
}

func (v *Voucher) UnmarshalJSON(data []byte) error {
	type Alias Voucher
	aux := &struct {
		BeginTime string `json:"beginTime"`
		EndTime   string `json:"endTime"`
		*Alias
	}{Alias: (*Alias)(v)}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// 然后使用预定义的格式解析时间字符串
	// time.Parse() 默认解析为 UTC 时区
	location, err := time.LoadLocation("Asia/Shanghai")
	beginTime, endTime := v.BeginTime, v.EndTime
	if v.BeginTime, err = time.ParseInLocation(time.DateTime, aux.BeginTime, location); err != nil {
		v.BeginTime = beginTime
	}

	if v.EndTime, err = time.ParseInLocation(time.DateTime, aux.EndTime, location); err != nil {
		v.EndTime = endTime
	}
	return nil
}
