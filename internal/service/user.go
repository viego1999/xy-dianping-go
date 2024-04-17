package service

import (
	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"time"
	"trpc.group/trpc-go/trpc-go/log"
	"xy-dianping-go/internal/common"
	"xy-dianping-go/internal/dto"
	"xy-dianping-go/internal/models"
	"xy-dianping-go/internal/repo"
	"xy-dianping-go/pkg/util"
)

type UserService interface {
	GetUserById(id int64) (*models.User, error)
	SendCode(phone string, session *sessions.Session) *dto.Result
	Login(dto dto.LoginFormDTO, session *sessions.Session) *dto.Result
}

type UserServiceImpl struct {
	userRepo repo.UserRepository
}

func NewUserService(userRepo repo.UserRepository) UserService {
	return &UserServiceImpl{userRepo: userRepo}
}

func (s *UserServiceImpl) GetUserById(id int64) (*models.User, error) {
	return s.userRepo.QueryById(id)
}

func (s *UserServiceImpl) SendCode(phone string, session *sessions.Session) *dto.Result {
	if util.IsPhoneInvalid(phone) {
		// 手机号格式错误
		return common.Fail("手机号格式错误")
	}

	// 符合，生成验证码
	code := util.RandomNumbers(6)
	// 保存到验证码 session
	session.Values[phone+":code"] = code
	// 发送验证码
	log.Debugf("发送验证码成功，验证码：{%s}", code)
	// 返回结果
	return common.OkWithData(code)
}

func (s *UserServiceImpl) Login(loginForm dto.LoginFormDTO, session *sessions.Session) *dto.Result {
	phone := loginForm.Phone
	if util.IsPhoneInvalid(phone) {
		// 手机号格式错误
		return common.Fail("手机号格式错误")
	}

	// 校验验证码
	if cacheCode, ok := session.Values[phone+":code"]; !ok || cacheCode != loginForm.Code {
		// 不一致，报错
		log.Debugf("验证码不一致，cacheCode:{%s}, code:{%s}。", cacheCode, loginForm.Code)
		return common.Fail("验证码不一致")
	}
	// 根据手机号查询用户
	user, _ := s.userRepo.QueryByPhone(phone)
	// 判断用户是否存在
	if user == nil {
		// 不存在，创建用户
		user = &models.User{
			Phone:      phone,
			NickName:   "user_" + util.RandomString(10),
			CreateTime: time.Now(),
			UpdateTime: time.Now(),
		}
		_ = s.userRepo.CreateUser(user)
	}
	// 保存信息到 session 中
	session.Values["user"] = &dto.UserDTO{
		Id:       user.Id,
		NickName: user.NickName,
		Icon:     user.Icon,
	}

	// 生成 token
	uid, _ := uuid.NewRandom()
	token := uid.String()
	return common.OkWithData(token)
}
