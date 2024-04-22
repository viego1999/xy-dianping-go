package service

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"strconv"
	"time"
	"trpc.group/trpc-go/trpc-go/log"
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
	// 保存优惠券
	err := s.voucherRepo.CreateVoucher(voucher)
	if err != nil {
		panic(err)
	}
	// 保存秒杀信息
	seckillVoucher := models.SeckillVoucher{
		VoucherId: voucher.Id,
		Stock:     voucher.Stock,
		BeginTime: voucher.BeginTime,
		EndTime:   voucher.EndTime,
	}
	log.Debugf("seckillVoucher: %+v", seckillVoucher)
	now := time.Now()
	if seckillVoucher.BeginTime.IsZero() {
		seckillVoucher.BeginTime = now
	}
	if seckillVoucher.EndTime.IsZero() {
		seckillVoucher.EndTime = now
	}
	err = s.seckillVoucherRepo.CreateSeckillVoucher(&seckillVoucher)
	log.Debugf("seckillVoucher: %+v", seckillVoucher)
	if err != nil {
		panic(err)
	}
	// 保存到 redis 中
	m := map[string]interface{}{
		"stock": strconv.Itoa(seckillVoucher.Stock),
		"begin": strconv.FormatInt(seckillVoucher.BeginTime.UnixNano()/int64(time.Millisecond), 10),
		"end":   strconv.FormatInt(seckillVoucher.EndTime.UnixNano()/int64(time.Millisecond), 10),
	}
	key := constants.SECKILL + fmt.Sprintf("%d", seckillVoucher.VoucherId)
	if _, err := s.redisClient.HSet(ctx, key, m).Result(); err != nil {
		panic(fmt.Sprintf("Error saving to Redis:%+v.", err))
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
