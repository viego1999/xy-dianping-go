package dto

import "encoding/gob"

type LoginFormDTO struct {
	Phone    string `json:"phone"`
	Code     string `json:"code"`
	Password string `json:"password"`
}

type Result struct {
	Success  bool   `json:"success"`
	ErrorMsg string `json:"errorMsg"`
	Data     any    `json:"data"`
	Total    int64  `json:"total"`
	Code     int64  `json:"code"`
}

type ScrollResult struct {
	List    []any `json:"list"`
	MinTime int64 `json:"minTime"`
	Offset  int   `json:"offset"`
}

type SeckillOrderDTO struct {
	UserId    int   `json:"userId"`
	VoucherId int   `json:"voucherId"`
	OrderId   int64 `json:"orderId"`
}

type UserDTO struct {
	Id       int64  `json:"id"`
	NickName string `json:"nickName"`
	Icon     string `json:"icon"`
}

func init() {
	// 注册类型以便在使用 securecookie 时可以正确编码和解码
	gob.Register(&UserDTO{})
}
