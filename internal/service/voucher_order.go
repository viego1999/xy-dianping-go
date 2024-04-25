package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"os"
	"strconv"
	"time"
	"trpc.group/trpc-go/trpc-database/goredis/redlock"
	"trpc.group/trpc-go/trpc-go/log"
	"xy-dianping-go/internal/common"
	"xy-dianping-go/internal/dto"
	"xy-dianping-go/internal/models"
	"xy-dianping-go/internal/mq/sender"
	"xy-dianping-go/internal/repo"
	"xy-dianping-go/pkg/utils"
	"xy-dianping-go/script/lua"
)

type VoucherOrderService interface {
	SeckillVoucher(ctx context.Context, voucherId int64) *dto.Result
	CreateVoucherOrder(ctx context.Context, order *models.VoucherOrder) error
}

type VoucherOrderServiceImpl struct {
	redisClient        redis.UniversalClient
	voucherOrderRepo   repo.VoucherOrderRepository
	seckillVoucherRepo repo.SeckillVoucherRepository
	sender             sender.MqSender
}

func NewVoucherOrderService(redisClient redis.UniversalClient, voucherOrderRepo repo.VoucherOrderRepository, seckillVoucherRepo repo.SeckillVoucherRepository, sender sender.MqSender) VoucherOrderService {
	return &VoucherOrderServiceImpl{redisClient, voucherOrderRepo, seckillVoucherRepo, sender}
}

func (s *VoucherOrderServiceImpl) SeckillVoucher(ctx context.Context, voucherId int64) *dto.Result {
	// 获取当前登录用户
	userDTO, _ := common.GetUserFromContext(ctx)
	userId := userDTO.Id
	// 1.执行 lua 脚本
	script, err := os.ReadFile(lua.SeckillLuaPath)
	if err != nil {
		panic(err)
	}
	result, err := s.redisClient.Eval(ctx, string(script), []string{}, voucherId, userId,
		time.Now().UnixNano()/int64(time.Millisecond)).Result()
	if err != nil {
		panic(err)
	}
	if result == nil {
		result = 0
	}
	result = int(result.(int64)) // int64(0) == int(0) 为 false

	// 2.判断结果是否为 0
	if result != 0 {
		// 2.1 不为 0，代表没有购买资格
		errorMsg := ""
		switch result {
		case 1:
			errorMsg = "秒杀尚未开始！"
		case 2:
			errorMsg = "秒杀已经结束！"
		case 3:
			errorMsg = "库存不足！"
		default:
			errorMsg = "不能重复下单！"
		}
		log.Warn(errorMsg)
		return common.Fail(errorMsg)
	}
	// 2.2 为 0，有购买资格，把下单信息保存到阻塞队列中
	// 2.3 订单 id
	orderId := utils.NextId(ctx, "order")
	voucherOrder := models.VoucherOrder{
		Id:        orderId,
		UserId:    userId,
		VoucherId: voucherId,
	}
	log.Debugf("voucherOrder: %+v", voucherOrder)
	// 发送到消息队列
	s.sender.SendSeckillMessage(&voucherOrder, false)
	// 3.返回订单 id
	return common.OkWithData(orderId)
}

func (s *VoucherOrderServiceImpl) CreateVoucherOrder(ctx context.Context, order *models.VoucherOrder) error {
	userId, voucherId := order.UserId, order.VoucherId
	// 创建锁对象
	redLock, err := redlock.New(s.redisClient)
	if err != nil {
		panic(fmt.Sprintf("redlock.New error: %+v", err))
	}
	lockKey := "order" + strconv.FormatInt(userId, 10)
	mu, err := redLock.TryLock(ctx, lockKey, []redlock.Option{redlock.WithKeyExpiration(10 * time.Second)}...)
	if err != nil {
		// 获取锁失败
		log.Errorf("不允许重复下单：%+v", err)
		return err
	}
	// 函数结束时释放锁
	defer func(mu redlock.Mutex, ctx context.Context) {
		if err := mu.Unlock(ctx); err != nil {
			panic(fmt.Sprintf("RedLock unlock error: %+v", err))
		}
	}(mu, ctx)

	// 进行事务操作
	err = s.voucherOrderRepo.ExecuteTransaction(func(txRepo repo.VoucherOrderRepository) error {
		// 5.1 查询订单
		o, _ := txRepo.QueryOrderByQuery("user_id = ? AND voucher_id = ?", userId, voucherId)
		// 5.2 判断是否存在
		if o != nil {
			// 已经购买过一次
			log.Errorf("该用户已经购买过一次！")
			return errors.New("该用户已经购买过一次！")
		}
		// 6.扣减库存
		rows, e := s.seckillVoucherRepo.UpdateSeckillVoucher(
			"voucher_id = ? AND stock > 0",
			voucherId,
			"stock",
			gorm.Expr("stock - 1"))
		if e != nil {
			return errors.New(fmt.Sprintf("UpdateSeckillVoucher error: %+v", e))
		}
		if rows == 0 {
			// 扣减失败
			return errors.New("库存不足")
		}
		// 保存订单
		e = s.voucherOrderRepo.CreateVoucherOrder(order)
		log.Debugf("CreateVoucherOrder success: %t.", e == nil)
		return e
	})
	return err
}
