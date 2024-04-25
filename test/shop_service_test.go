package test

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"os"
	"testing"
	"trpc.group/trpc-go/trpc-go"
	"trpc.group/trpc-go/trpc-go/log"
	"xy-dianping-go/internal/constants"
	"xy-dianping-go/internal/db"
	"xy-dianping-go/internal/models"
	"xy-dianping-go/internal/repo"
	"xy-dianping-go/internal/service"
	"xy-dianping-go/pkg/utils"
)

var (
	ctx         = context.Background()
	Db          *gorm.DB
	redisClient redis.UniversalClient
	shopRepo    repo.ShopRepository
	shopService service.ShopService
)

func TestMain(m *testing.M) {
	os.Args = append(os.Args, "-conf", "../conf/trpc_go.yaml")

	trpc.NewServer() // 主要作用是读取数据库的配置信息
	// 执行启动逻辑
	log.Info("执行测试启动逻辑，初始化...")
	Db = db.InitDatabase()
	redisClient = db.InitRedisClient()
	shopRepo = repo.NewShopRepository(Db)
	shopService = service.NewShopService(redisClient, shopRepo)

	// 运行测试
	code := m.Run()

	// 测试完后的逻辑
	log.Info("测试结束。")

	// 退出测试
	os.Exit(code)
}

// 提前将所有店铺信息按照店铺类型以 GEO 格式存储到 redis 中
func TestLoadShopData(t *testing.T) {
	log.Info("测试 LoadShopData().")
	// 1.查询店铺信息
	list, err := shopRepo.List()
	if err != nil {
		log.Error("查询店铺信息错误：", err)
		return
	}
	// 2.把店铺分组，按照 typeId 分组，id 一致的放到一个集合
	mapOfShops := utils.GroupBy(list, func(s models.Shop) int64 { return s.TypeId })
	// 3.分批完成写入 Redis
	for typeId, shops := range mapOfShops {
		// 拼接 key
		key := constants.SHOP_GEO_KEY + fmt.Sprintf("%d", typeId)
		// 构造经纬度集合
		geoLocations := make([]*redis.GeoLocation, 0, len(shops))
		for _, shop := range shops {
			geoLocations = append(geoLocations, &redis.GeoLocation{
				Name:      fmt.Sprintf("%d", shop.Id),
				Longitude: shop.X,
				Latitude:  shop.Y,
			})
		}

		// 4.批量写入 redis
		_, err = redisClient.GeoAdd(ctx, key, geoLocations...).Result()
		if err != nil {
			log.Error("GeoAdd error:", err)
		}
	}
}

func TestQueryShopById(t *testing.T) {
	log.Info("test QueryShopById()")
	result := shopService.QueryShopById(ctx, 1)
	t.Logf("result: %+v", result)
}
