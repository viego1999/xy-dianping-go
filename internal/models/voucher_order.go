package models

import (
	"encoding/json"
	"time"
)

type VoucherOrder struct {
	Id         int64     `json:"id" gorm:"column:id"`
	UserId     int64     `json:"userId" gorm:"column:user_id"`
	VoucherId  int64     `json:"voucherId" gorm:"column:voucher_id"`
	PayType    int       `json:"payType" gorm:"column:pay_type"`
	Status     int       `json:"status" gorm:"column:status"`
	PayTime    time.Time `json:"payTime" gorm:"column:pay_time"`
	UseTime    time.Time `json:"useTime" gorm:"column:use_time"`
	RefundTime time.Time `json:"refundTime" gorm:"column:refund_time"`
	CreateTime time.Time `json:"createTime" gorm:"column:create_time;autoCreateTime"` // 注意这里使用了column标签指定字段名
	UpdateTime time.Time `json:"updateTime" gorm:"column:update_time;autoUpdateTime"` // 注意这里使用了column标签指定字段名
}

func (*VoucherOrder) TableName() string {
	return "tb_voucher_order"
}

func (v *VoucherOrder) UnmarshalJSON(data []byte) error {
	type Alias VoucherOrder
	aux := &struct {
		PayTime    string `json:"beginTime"`
		UseTime    string `json:"useTime"`
		RefundTime string `json:"refundTime"`
		*Alias
	}{Alias: (*Alias)(v)}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// 然后使用预定义的格式解析时间字符串
	// time.Parse() 默认解析为 UTC 时区
	location, err := time.LoadLocation("Asia/Shanghai")
	payTime, useTime, refundTime := v.PayTime, v.UseTime, v.RefundTime
	if v.PayTime, err = time.ParseInLocation(time.DateTime, aux.PayTime, location); err != nil {
		v.PayTime = payTime // 解析失败返回原始值
	}
	if v.UseTime, err = time.ParseInLocation(time.DateTime, aux.UseTime, location); err != nil {
		v.UseTime = useTime
	}
	if v.RefundTime, err = time.ParseInLocation(time.DateTime, aux.RefundTime, location); err != nil {
		v.UseTime = refundTime
	}
	return nil
}
