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
	layout := "2006-01-02 15:04:05"
	if v.BeginTime, err = time.ParseInLocation(layout, aux.BeginTime, location); err != nil {
		return err
	}

	if v.EndTime, err = time.ParseInLocation(layout, aux.EndTime, location); err != nil {
		return err
	}
	return nil
}
