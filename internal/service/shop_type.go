package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"xy-dianping-go/internal/common"
	"xy-dianping-go/internal/constants"
	"xy-dianping-go/internal/dto"
	"xy-dianping-go/internal/models"
	"xy-dianping-go/internal/repo"
)

type ShopTypeService interface {
	QueryTypeList(ctx context.Context) *dto.Result
}

type ShopTypeServiceImpl struct {
	redisClient  redis.UniversalClient
	shopTypeRepo repo.ShopTypeRepository
}

func NewShopTypeService(redisClient redis.UniversalClient, repo repo.ShopTypeRepository) ShopTypeService {
	return &ShopTypeServiceImpl{redisClient, repo}
}

func (s *ShopTypeServiceImpl) QueryTypeList(ctx context.Context) *dto.Result {
	typeString, err := s.redisClient.LRange(ctx, constants.CACHE_SHOP_TYPE_KEY, 0, -1).Result()
	if err != nil {
		panic(fmt.Sprintf("QueryTypeList - redis LRange error: %+v", err))
	}
	types := make([]models.ShopType, 0, len(typeString))
	// 判断 redis 是否存在数据
	if len(typeString) > 0 {
		for _, str := range typeString { // 存在
			var shopType models.ShopType
			if err := json.Unmarshal([]byte(str), &shopType); err != nil {
				panic(fmt.Sprintf("QueryTypeList - json unmarshal error: %+v", err))
			}
			types = append(types, shopType)
		}
		// 返回结果
		return common.OkWithData(types)
	}
	// 不存在，查询数据库
	types, err = s.shopTypeRepo.List()
	if err != nil {
		panic(fmt.Sprintf("QueryTypeList repo List error: %+v", err))
	}
	// 将查询结果存储到 redis
	for _, shopType := range types {
		jsonByts, err := json.Marshal(shopType)
		if err != nil {
			panic(fmt.Sprintf("QueryTypeList json marshal error: %+v", err))
		}
		if err = s.redisClient.RPush(ctx, constants.CACHE_SHOP_TYPE_KEY, jsonByts).Err(); err != nil {
			panic(fmt.Sprintf("QueryTypeList redis RPush error: %+v", err))
		}
	}

	return common.OkWithData(types)
}
