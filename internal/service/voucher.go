package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"strconv"
	"time"
	"xy-dianping-go/internal/common"
	"xy-dianping-go/internal/constants"
	"xy-dianping-go/internal/dto"
	"xy-dianping-go/internal/models"
	"xy-dianping-go/internal/repo"
)

type VoucherService interface {
	SaveVoucher(voucher *models.Voucher) *dto.Result
	SaveSeckillVoucher(ctx context.Context, voucher *models.Voucher) *dto.Result
	QueryVoucherOfShop(shopId int64) *dto.Result
}

type VoucherServiceImpl struct {
	redisClient        redis.UniversalClient
	voucherRepo        repo.VoucherRepository
	seckillVoucherRepo repo.SeckillVoucherRepository
}

func NewVoucherService(redisClient redis.UniversalClient, voucherRepo repo.VoucherRepository, seckillVoucherRepo repo.SeckillVoucherRepository) VoucherService {
	return &VoucherServiceImpl{redisClient, voucherRepo, seckillVoucherRepo}
}

func (s *VoucherServiceImpl) SaveVoucher(voucher *models.Voucher) *dto.Result {
	err := s.voucherRepo.CreateVoucher(voucher)
	if err != nil {
		panic(err)
	}
	return common.OkWithData(voucher.Id)
}

func (s *VoucherServiceImpl) SaveSeckillVoucher(ctx context.Context, voucher *models.Voucher) *dto.Result {
	// 开启事务
	err := s.voucherRepo.ExecuteTransaction(func(txRepo repo.VoucherRepository) error {
		// 保存优惠券
		if err := s.voucherRepo.CreateVoucher(voucher); err != nil {
			return err
		}
		// 保存秒杀信息
		seckillVoucher := models.SeckillVoucher{
			VoucherId: voucher.Id,
			Stock:     voucher.Stock,
			BeginTime: voucher.BeginTime,
			EndTime:   voucher.EndTime,
		}
		// 未设定活动起止时间则自动设置为当前时间
		now := time.Now()
		if seckillVoucher.BeginTime.IsZero() {
			seckillVoucher.BeginTime = now
		}
		if seckillVoucher.EndTime.IsZero() {
			seckillVoucher.EndTime = now
		}
		// 保存秒杀优惠券
		if err := s.seckillVoucherRepo.CreateSeckillVoucher(&seckillVoucher); err != nil {
			return err
		}
		// 保存到 redis 中
		m := map[string]interface{}{
			"stock": strconv.Itoa(seckillVoucher.Stock),
			"begin": strconv.FormatInt(seckillVoucher.BeginTime.UnixNano()/int64(time.Millisecond), 10),
			"end":   strconv.FormatInt(seckillVoucher.EndTime.UnixNano()/int64(time.Millisecond), 10),
		}
		key := constants.SECKILL + fmt.Sprintf("%d", seckillVoucher.VoucherId)
		if _, err := s.redisClient.HSet(ctx, key, m).Result(); err != nil {
			return errors.New(fmt.Sprintf("Error saving to Redis:%+v.", err))
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	return common.OkWithData(voucher.Id)
}

func (s *VoucherServiceImpl) QueryVoucherOfShop(shopId int64) *dto.Result {
	// 查询优惠券信息
	vouchers, err := s.voucherRepo.QueryVoucherByShopId(shopId)
	if err != nil {
		panic(err)
	}
	// 返回结果
	return common.OkWithData(vouchers)
}
