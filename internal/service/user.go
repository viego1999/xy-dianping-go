package service

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/jinzhu/copier"
	"github.com/redis/go-redis/v9"
	"strconv"
	"time"
	"trpc.group/trpc-go/trpc-go/log"
	"xy-dianping-go/internal/common"
	"xy-dianping-go/internal/constants"
	"xy-dianping-go/internal/dto"
	"xy-dianping-go/internal/models"
	"xy-dianping-go/internal/repo"
	"xy-dianping-go/pkg/utils"
)

type UserService interface {
	GetUserById(id int64) *dto.Result
	SendCode(ctx context.Context, phone string) *dto.Result
	Login(ctx context.Context, loginForm *dto.LoginFormDTO) *dto.Result
	Me(ctx context.Context) *dto.Result
	Sign(ctx context.Context) *dto.Result
	SignCount(ctx context.Context) *dto.Result
}

type UserServiceImpl struct {
	redisClient redis.UniversalClient
	userRepo    repo.UserRepository
}

func NewUserService(redisClient redis.UniversalClient, userRepo repo.UserRepository) UserService {
	return &UserServiceImpl{redisClient: redisClient, userRepo: userRepo}
}

func (s *UserServiceImpl) GetUserById(id int64) *dto.Result {
	// 查询详情
	user, _ := s.userRepo.QueryById(id)
	if user == nil {
		return common.Ok()
	}

	userDTO := dto.UserDTO{}
	_ = copier.Copy(&userDTO, user) // 值填充

	return common.OkWithData(userDTO)
}

func (s *UserServiceImpl) SendCode(ctx context.Context, phone string) *dto.Result {
	// 1.校验手机号
	if utils.IsPhoneInvalid(phone) {
		// 2.手机号格式错误
		return common.Fail("手机号格式错误")
	}

	// 3.符合，生成验证码
	code := utils.RandomNumbers(6)
	// 4.保存验证码到 redis
	s.redisClient.Set(ctx, constants.LOGIN_CODE_KEY+phone, code, time.Minute*constants.LOGIN_CODE_TTL)
	// 5.发送短信验证码
	log.Debugf("发送短信验证码成功，验证码：{%s}", code)
	// 返回 Ok
	return common.OkWithData(code)
}

func (s *UserServiceImpl) Login(ctx context.Context, loginForm *dto.LoginFormDTO) *dto.Result {
	// 1.校验手机号
	phone := loginForm.Phone
	if utils.IsPhoneInvalid(phone) {
		// 2.手机号格式错误
		return common.Fail("手机号格式错误")
	}

	// 2.从 redis 中取出验证码并校验
	if cacheCode, err := s.redisClient.Get(ctx, constants.LOGIN_CODE_KEY+phone).Result(); err != nil || cacheCode != loginForm.Code {
		// 3.不一致，报错
		log.Debugf("验证码不一致，cacheCode:{%s}, code:{%s}。", cacheCode, loginForm.Code)
		return common.Fail("验证码不一致")
	}
	// 4.一致，根据手机号查询用户
	user, _ := s.userRepo.QueryByPhone(phone)
	// 5.判断用户是否存在
	if user == nil {
		// 6.不存在，创建用户
		user = &models.User{
			Phone:      phone,
			NickName:   "user_" + utils.RandomString(10),
			CreateTime: time.Now(),
			UpdateTime: time.Now(),
		}
		_ = s.userRepo.CreateUser(user)
	}
	// 7.保存信息到 redis 中
	userMap := utils.StructToMap(&dto.UserDTO{
		Id:       user.Id,
		NickName: user.NickName,
		Icon:     user.Icon,
	})

	// 生成 token
	uid, _ := uuid.NewRandom()
	token := uid.String()
	tokenKey := constants.LOGIN_USER_KEY + token
	// 存储到 redis，同时设置多个字段
	if err := s.redisClient.HSet(ctx, tokenKey, userMap).Err(); err != nil {
		return common.Fail("Redis存储用户信息失败：" + err.Error())
	}
	// 设置 token 过期时间
	s.redisClient.Expire(ctx, tokenKey, constants.LOGIN_USER_TTL*time.Second)
	return common.OkWithData(token)
}

// SendCodeWithSession 将验证码存储到 http.session 中
//
// Deprecated: 由于 http.session 中容易出现共享问题，因此使用基于 redis 的 SendCode 方法
func (s *UserServiceImpl) SendCodeWithSession(phone string, session *sessions.Session) *dto.Result {
	if utils.IsPhoneInvalid(phone) {
		// 手机号格式错误
		return common.Fail("手机号格式错误")
	}

	// 符合，生成验证码
	code := utils.RandomNumbers(6)
	// 保存到验证码 session
	session.Values[phone+":code"] = code
	// 发送验证码
	log.Debugf("发送验证码成功，验证码：{%s}", code)
	// 返回结果
	return common.OkWithData(code)
}

// LoginWithSession 基于 http.session 的登录功能
//
// Deprecated: 使用基于 redis 的 Login 方法代替
func (s *UserServiceImpl) LoginWithSession(loginForm dto.LoginFormDTO, session *sessions.Session) *dto.Result {
	phone := loginForm.Phone
	if utils.IsPhoneInvalid(phone) {
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
			NickName:   constants.USER_NICK_NAME_PREFIX + utils.RandomString(10),
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

func (s *UserServiceImpl) Me(ctx context.Context) *dto.Result {
	// 获取当前登录用户
	userDTO, ok := common.GetUserFromContext(ctx)
	if !ok {
		return common.Fail("获取当前用户失败")
	}

	return common.OkWithData(userDTO)
}

func (s *UserServiceImpl) Sign(ctx context.Context) *dto.Result {
	// 1.获取当前登录用户
	userDTO, ok := common.GetUserFromContext(ctx)
	if !ok {
		return common.Fail("获取当前用户失败")
	}
	// 2.获取日期
	now := time.Now()
	// 3.拼接key
	keySuffix := now.Format(":200601")
	key := constants.USER_SIGN_KEY + strconv.Itoa(int(userDTO.Id)) + keySuffix
	// 4.获取今天时本月的第几天
	dayOfMoth := now.Day()
	// 5.写入 redis `SETBIT key offset 1`
	s.redisClient.SetBit(ctx, key, int64(dayOfMoth-1), 1)
	// 回复结果
	return common.Ok()
}

func (s *UserServiceImpl) SignCount(ctx context.Context) *dto.Result {
	// 1.获取当前登录用户
	userDTO, _ := common.GetUserFromContext(ctx)
	// 2.获取日期
	now := time.Now()
	// 3.拼接 key
	keySuffix := now.Format(":200601")
	key := constants.USER_SIGN_KEY + strconv.Itoa(int(userDTO.Id)) + keySuffix
	// 4.获取今天是本月第几天
	dayOfMonth := now.Day()
	// 5.获取本月截止到今天为止的所有签到记录
	// 构造 BITFIELD 命令
	cmd := []interface{}{"GET", fmt.Sprintf("u%d", dayOfMonth), 0} // ud 表示无符号d位整数，0 为偏移量
	result, err := s.redisClient.BitField(ctx, key, cmd...).Result()

	if err != nil {
		return common.Fail("Redis bitField error: " + err.Error())
	}

	if len(result) == 0 {
		return common.OkWithData(0)
	}

	num := result[0]
	if num == 0 {
		return common.OkWithData(0)
	}
	count := 0
	// 6.循环遍历
	for num != 0 {
		// 7.位运算遍历
		if (num & 1) == 0 { // 未签到
			break
		}
		count++
		num >>= 1
	}
	return common.OkWithData(count)
}
