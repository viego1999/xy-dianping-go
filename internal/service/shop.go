package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"strconv"
	"time"
	"trpc.group/trpc-go/trpc-database/goredis/redlock"
	"xy-dianping-go/internal/common"
	"xy-dianping-go/internal/constants"
	"xy-dianping-go/internal/dto"
	"xy-dianping-go/internal/models"
	"xy-dianping-go/internal/repo"
)

type ShopService interface {
	QueryShopById(ctx context.Context, id int64) *dto.Result
	SaveShop(shop *models.Shop) *dto.Result
	UpdateShop(ctx context.Context, shop *models.Shop) *dto.Result
}

type ShopServiceImpl struct {
	redisClient redis.UniversalClient
	shopRepo    repo.ShopRepository
}

func NewShopService(redisClient redis.UniversalClient, shopRepo repo.ShopRepository) ShopService {
	return &ShopServiceImpl{redisClient, shopRepo}
}

// QueryShopById 获取店铺信息，基于 redis 的分布式锁解决缓存穿透
func (s *ShopServiceImpl) QueryShopById(ctx context.Context, id int64) *dto.Result {
	var shop = &models.Shop{}
	idStr := strconv.Itoa(int(id))
	// 1.从 redis 缓存中查询商铺
	key := constants.CACHE_SHOP_KEY + idStr
	shopJson, err := s.redisClient.Get(ctx, key).Result()
	// 2.判断是否存在
	if err != redis.Nil { // redis 中存在记录
		if err != nil { // redis 操作出现其他错误
			panic(fmt.Sprintf("QueryShopById - Redis Get error: %v", err))
		}
		// 3.存在，直接返回
		// 判断是否为 null，缓存的空值
		if shopJson == "null" {
			return common.Fail("店铺不存在！")
		}
		// 解码 json 为 shop
		if err = json.Unmarshal([]byte(shopJson), shop); err != nil {
			panic(fmt.Sprintf("QueryShopById - JSON Unmarshal error: %v", err))
		} else {
			return common.OkWithData(shop)
		}
	}
	// 4.实现缓存重建
	// 4.1 获取互斥锁
	lockKey := constants.LOCK_SHOP_KEY + idStr
	// 创建分布式锁
	redLock, err := redlock.New(s.redisClient)
	mu, err := redLock.TryLock(ctx, lockKey)
	// 4.2 判断锁是否获取成功
	if err != nil {
		// 4.3 失败，则休眠并重试
		time.Sleep(50 * time.Millisecond)
		return s.QueryShopById(ctx, id)
	}

	// 加锁成功后，最后需要释放锁
	defer func(mu redlock.Mutex, ctx context.Context) {
		if err := mu.Unlock(ctx); err != nil {
			panic(fmt.Sprintf("QeuryShopById - lock fail %v", err))
		}
	}(mu, ctx)

	// 4.4 加锁成功，根据 id 查询数据库
	shop, _ = s.shopRepo.QueryById(id)
	// 5.不存在，返回错误
	if shop == nil {
		// 【缓存穿透】将空值写入 redis
		s.redisClient.Set(ctx, key, "null", constants.CACHE_NULL_TTL*time.Minute)
		// 返回错误信息
		return common.Fail("店铺不存在！")
	}
	// 6.存在，写入 redis
	shopBytes, _ := json.Marshal(shop)
	s.redisClient.Set(ctx, key, string(shopBytes), constants.CACHE_SHOP_TTL*time.Minute)
	// 返回结果
	return common.OkWithData(shop)
}

func (s *ShopServiceImpl) SaveShop(shop *models.Shop) *dto.Result {
	err := s.shopRepo.CreateShop(shop)
	if err != nil {
		panic(fmt.Sprintf("SaveShop - save shop error: %v", err))
	}
	return common.OkWithData(shop.Id)
}

func (s *ShopServiceImpl) UpdateShop(ctx context.Context, shop *models.Shop) *dto.Result {
	// 获取 id
	id := shop.Id
	if id == 0 {
		return common.Fail("店铺id不能为空！")
	}
	// 使用事务进行数据库操作，当缓存删除发生异常时方便回滚
	// 开始事务，返回错误进行回滚，否则提交事务
	err := s.shopRepo.ExecuteTransaction(func(repo repo.ShopRepository) error {
		// 1.更新数据库
		err := repo.Update(shop)
		if err != nil {
			return errors.New(fmt.Sprintf("UpdateShop - update shop by map error: %+v.", err))
		}
		// 2.删除缓存（当删除失败可能导致缓存中为旧数据）
		key := constants.CACHE_SHOP_KEY + strconv.Itoa(int(id))
		_, err = s.redisClient.Del(ctx, key).Result()
		if err != nil {
			return errors.New(fmt.Sprintf("UpdateShop - Redis delete error: %+v.", err))
		}
		// 如果没有错误，返回 nil 以提交事务
		return nil
	})
	if err != nil {
		panic(err)
	}
	return common.Ok()
}
