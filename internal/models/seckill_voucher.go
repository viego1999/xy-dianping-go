package models

import "time"

type SeckillVoucher struct {
	VoucherId  int64     `json:"voucherId" gorm:"column:voucher_id;primaryKey"`
	Stock      int       `json:"stock" gorm:"column:stock"`
	BeginTime  time.Time `json:"beginTime" gorm:"column:begin_time;autoBeginTime"`
	EndTime    time.Time `json:"endTime" gorm:"column:end_time;autoEndTime"`
	CreateTime time.Time `json:"createTime" gorm:"column:create_time;autoCreateTime"` // 注意这里使用了column标签指定字段名
	UpdateTime time.Time `json:"updateTime" gorm:"column:update_time;autoUpdateTime"` // 注意这里使用了column标签指定字段名
}

func (SeckillVoucher) TableName() string {
	return "tb_seckill_voucher"
}
