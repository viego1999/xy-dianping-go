package v1

import (
	"encoding/json"
	"github.com/gorilla/sessions"
	"net/http"
	"trpc.group/trpc-go/trpc-go/log"
	"xy-dianping-go/internal/common"
	"xy-dianping-go/internal/dto"
	"xy-dianping-go/internal/service"
)

type UserController struct {
	userService service.UserService
}

func NewUserController(userService service.UserService) *UserController {
	return &UserController{userService}
}

func (c *UserController) SendCode(w http.ResponseWriter, r *http.Request) {
	// 解析表单数据
	err := r.ParseForm()
	if err != nil {
		common.SendResponseWithCode(w, common.Fail("表单数据解析失败："+err.Error()), http.StatusBadRequest)
		return
	}

	// 获取表单值
	phone := r.FormValue("phone")

	// 获取 session
	session := r.Context().Value("session").(*sessions.Session)

	result := c.userService.SendCode(phone, session)
	// 保存 session 设置的值
	common.SessionSave(session, r, w)
	// 发送响应
	common.SendResponse(w, result)
}

func (c *UserController) Login(w http.ResponseWriter, r *http.Request) {
	var loginForm dto.LoginFormDTO
	// 检验手机号
	err := json.NewDecoder(r.Body).Decode(&loginForm)
	if err != nil {
		common.SendResponseWithCode(w, common.Fail("Bad request"), http.StatusBadRequest)
		return
	}

	// 获取 session
	session := r.Context().Value("session").(*sessions.Session)

	// 一致，根据手机号码查询用户
	result := c.userService.Login(loginForm, session)
	// 保存 session 的值
	common.SessionSave(session, r, w)
	// 发送响应
	common.SendResponse(w, result)
}

func (c *UserController) Sign(w http.ResponseWriter, r *http.Request) {
	// 获取当前登录用户
	userDTO, ok := common.GetUserFromContext(r.Context())
	if !ok {
		common.SendResponse(w, common.Fail("获取当前用户失败"))
		return
	}
	// 获取日期
	log.Infof("userDTO: %v, ok: %t.", userDTO, ok)
	// 进行签到

	// 回复结果
	common.SendResponse(w, common.Ok())
}
