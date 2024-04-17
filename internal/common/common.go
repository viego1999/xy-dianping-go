package common

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/sessions"
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
	// 设置 Content-Type 头信息，表明返回的内容类型为 JSON
	w.Header().Set("Content-Type", "application/json")
	// 将结构体编码为 JSON，并写入响应中
	if err := json.NewEncoder(w).Encode(result); err != nil {
		log.Errorf("Error encoding JSON: %v", err)
		panic(fmt.Sprintf("Error encoding JSON: %v", result))
	}
}

func SendResponseWithCode(w http.ResponseWriter, result *dto.Result, statusCode int) {
	// 设置 Content-Type 头信息，表明返回的内容类型为 JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode) // 设置状态码

	// 将结构体编码为 JSON，并写入响应中
	if err := json.NewEncoder(w).Encode(result); err != nil {
		log.Errorf("Error encoding JSON: %v", err)
		panic(fmt.Sprintf("Error encoding JSON: %v", result))
	}
}

func SessionSave(session *sessions.Session, r *http.Request, w http.ResponseWriter) {
	// 保存 session 设置的值
	if err := session.Save(r, w); err != nil {
		log.Errorf("Session save data failed, error: %v.", err)
		panic("Session数据设置失败：" + err.Error())
	}
}

// 实现 UserHolder 功能，ThreadLocal
// user 上下文键，用来保存当前登录用户的键
type userCtxKey struct{}

// SetUserToContext 保存当前上下文用户信息
func SetUserToContext(ctx context.Context, user *dto.UserDTO) context.Context {
	return context.WithValue(ctx, userCtxKey{}, user)
}

// GetUserFromContext 获取当前上下文用户信息
func GetUserFromContext(ctx context.Context) (*dto.UserDTO, bool) {
	user, ok := ctx.Value(userCtxKey{}).(*dto.UserDTO)
	return user, ok
}
