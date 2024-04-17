package v1

import (
	"encoding/json"
	"github.com/gorilla/sessions"
	"net/http"
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
	if err := session.Save(r, w); err != nil {
		common.SendResponseWithCode(w, common.Fail("Session数据设置失败:"+err.Error()), http.StatusInternalServerError)
		return
	}
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
	common.SendResponse(w, c.userService.Login(loginForm, session))
}
