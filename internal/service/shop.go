package service

import (
	"context"
	"encoding/json"
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
