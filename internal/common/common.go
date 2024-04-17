package common

import (
	"encoding/json"
	"fmt"
	"net/http"
	"trpc.group/trpc-go/trpc-go/log"
	"xy-dianping-go/internal/dto"
)

func Ok() *dto.Result {
	return &dto.Result{
		Success: true,
		Code:    200,
	}
}

func OkWithData(data any) *dto.Result {
	return &dto.Result{
		Success: true,
		Data:    data,
		Code:    200,
	}
}

func OkWithDataCode(data any, code int64) *dto.Result {
	return &dto.Result{
		Success: true,
		Data:    data,
		Code:    code,
	}
}

func OkWithDataTotal(data any, total int64) *dto.Result {
	return &dto.Result{
		Success: true,
		Data:    data,
		Total:   total,
		Code:    200,
	}
}

func Fail(errorMsg string) *dto.Result {
	return &dto.Result{
		ErrorMsg: errorMsg,
		Code:     500,
	}
}

func FailWithCode(errorMsg string, code int64) *dto.Result {
	return &dto.Result{
		ErrorMsg: errorMsg,
		Code:     code,
	}
}

func SendResponse(w http.ResponseWriter, result *dto.Result) {
	w.Header().Set("Content-Type", "application/json")

	// 将结构体编码为 JSON，并写入响应中
	err := json.NewEncoder(w).Encode(result)
	if err != nil {
		log.Errorf("Error encoding JSON: %v", err)
		panic(fmt.Sprintf("Error encoding JSON: %v", result))
	}
}

func SendResponseWithCode(w http.ResponseWriter, result *dto.Result, statusCode int) {
	// 设置 Content-Type 头信息，表明返回的内容类型为 JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode) // 设置状态码

	// 将结构体编码为 JSON，并写入响应中
	err := json.NewEncoder(w).Encode(result)
	if err != nil {
		log.Errorf("Error encoding JSON: %v", err)
		panic(fmt.Sprintf("Error encoding JSON: %v", result))
	}
}
