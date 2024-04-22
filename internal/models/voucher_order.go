package models

import "time"

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

func (VoucherOrder) TableName() string {
	return "tb_voucher_order"
}
