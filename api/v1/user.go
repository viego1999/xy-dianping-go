package v1

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"xy-dianping-go/internal/common"
	"xy-dianping-go/internal/dto"
	"xy-dianping-go/internal/service"
)

type UserController struct {
	userService     service.UserService
	userInfoService service.UserInfoService
}

func NewUserController(userService service.UserService, userInfoService service.UserInfoService) *UserController {
	return &UserController{userService, userInfoService}
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

	// ======= 将验证码存储到 http.session 中 ======
	//// 获取 session
	//session := r.Context().Value("session").(*sessions.Session)
	//result := c.userService.SendCodeWithSession(phone, session)
	//// 保存 session 设置的值
	//common.SessionSave(session, r, w)
	// ===========    END Send Code   ============

	result := c.userService.SendCode(r.Context(), phone)

	// 发送响应
	common.SendResponse(w, result)
}

func (c *UserController) Login(w http.ResponseWriter, r *http.Request) {
	var loginForm dto.LoginFormDTO
	// 获取前端登录请求信息
	err := json.NewDecoder(r.Body).Decode(&loginForm)
	if err != nil {
		common.SendResponseWithCode(w, common.Fail("Bad request"), http.StatusBadRequest)
		return
	}

	// =========== 基于 http.session 进行登录验证 ==========
	//// 获取 session
	//session := r.Context().Value("session").(*sessions.Session)
	//
	//// 一致，根据手机号码查询用户
	//result := c.userService.LoginWithSession(loginForm, session)
	//// 保存 session 的值
	//common.SessionSave(session, r, w)
	// ==========  END LOGIN WiTH SESSION ==============

	result := c.userService.Login(r.Context(), &loginForm)

	// 发送响应
	common.SendResponse(w, result)
}

func (c *UserController) Me(w http.ResponseWriter, r *http.Request) {

	common.SendResponse(w, c.userService.Me(r.Context()))
}

func (c *UserController) Info(w http.ResponseWriter, r *http.Request) {
	// 获取用户id
	vars := mux.Vars(r)
	userIdStr := vars["id"]
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		common.SendResponse(w, common.FailWithCode("Invalid user id", http.StatusBadRequest))
		return
	}
	userInfo, err := c.userInfoService.GetUserInfoByUserId(int64(userId))
	if userInfo == nil {
		// 没有详情，应该是第一次查看详情
		common.SendResponse(w, common.Ok())
		return
	}
	// 返回
	common.SendResponse(w, common.OkWithData(userInfo))
}

func (c *UserController) QueryUserById(w http.ResponseWriter, r *http.Request) {
	// 获取用户id
	vars := mux.Vars(r)
	userIdStr := vars["id"]
	id, err := strconv.Atoi(userIdStr)
	if err != nil {
		common.SendResponse(w, common.FailWithCode("Invalid user id", http.StatusBadRequest))
		return
	}

	common.SendResponse(w, c.userService.GetUserById(int64(id)))
}

func (c *UserController) Sign(w http.ResponseWriter, r *http.Request) {

	common.SendResponse(w, c.userService.Sign(r.Context()))
}

func (c *UserController) SignCount(w http.ResponseWriter, r *http.Request) {

	common.SendResponse(w, c.userService.SignCount(r.Context()))
}
