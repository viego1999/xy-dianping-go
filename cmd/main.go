package main

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/redis/go-redis/v9"
	"github.com/streadway/amqp"
	"go.uber.org/fx"
	"gorm.io/gorm"
	"trpc.group/trpc-go/trpc-go"
	thttp "trpc.group/trpc-go/trpc-go/http"
	"trpc.group/trpc-go/trpc-go/log"
	"trpc.group/trpc-go/trpc-go/server"
	"xy-dianping-go/api"
	v1 "xy-dianping-go/api/v1"
	"xy-dianping-go/internal/config"
	"xy-dianping-go/internal/db"
	"xy-dianping-go/internal/mq"
	"xy-dianping-go/internal/mq/receiver"
	"xy-dianping-go/internal/mq/sender"
	"xy-dianping-go/internal/repo"
	"xy-dianping-go/internal/service"
)

func main() {
	// 定义 fx 应用
	// fx.New 返回的应用程序会在调用.Run()时开始运行，同时等待所有注册的OnStart hooks执行完成。
	// 当所有的OnStart hooks执行完成后，应用程序会继续运行，直到调用了应用程序的Stop方法（通常是通过调用OnStop hooks来触发）。
	// 所以，当fx应用程序运行时，main goroutine会阻塞，直到应用程序停止。
	fx.New(
		// 创建并导入 trpc 服务
		fx.Provide(trpc.NewServer),
		// 提供 *gorm.DB 实例
		// 初始化数据库连接
		fx.Provide(db.InitDatabase, db.InitRedisClient, mq.InitAMQPConnection),
		// 提供 Repository 的实例，依赖于 *gorm.DB
		fx.Provide(repo.NewUserRepository, repo.NewUserInfoRepository, repo.NewShopRepository, repo.NewVoucherRepository,
			repo.NewSeckillVoucherRepository, repo.NewVoucherOrderRepository, repo.NewBlogRepository, sender.NewMqSender),
		// 提供 Service 的实例，依赖于 Repository
		fx.Provide(service.NewUserService, service.NewUserInfoService, service.NewShopService, service.NewVoucherService,
			service.NewVoucherOrderService, receiver.NewMqReceiver, service.NewBlogService),
		// 提供 Controller 的实例，依赖于 Service
		fx.Provide(v1.NewUserController, v1.NewShopController, v1.NewVoucherController, v1.NewVoucherOrderController,
			v1.NewBlogController),
		// 导入路由模块
		api.Module,

		// 定义一个 Invoke 函数，用于在应用启动时执行一些操作
		fx.Invoke(func(lc fx.Lifecycle, s *server.Server, router *mux.Router, receiver receiver.MqReceiver, Db *gorm.DB,
			redisClient redis.UniversalClient, amqpConn *amqp.Connection) {
			// 示例：添加生命周期钩子，在 app 关闭时执行清理操作
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					// 向 trpc 中注册 HTTP 服务
					thttp.RegisterNoProtocolServiceMux(s.Service(config.HttpServiceName), router)

					// 启动 rabbitmq 消费者监听秒杀队列进行消费
					go receiver.ReceiveSeckillOrder(context.Background())

					// 启动 trpc 服务
					go func() {
						if err := s.Serve(); err != nil {
							log.Fatal("trpc server failed to launch:", err)
							return
						}
					}()
					log.Info("trpc server started successfully.")
					return nil
				},
				OnStop: func(ctx context.Context) error {
					log.Info("Cleaning up resources...")
					// 注入关闭资源
					log.Info("Close database connection...")
					sqlDB, err := Db.DB()
					if err != nil {
						log.Errorf("Get sqlDB error: %+v", err)
						return err
					}
					if err = sqlDB.Close(); err != nil {
						log.Errorf("Failed to close database connection: %+v", err)
						return err
					}
					log.Info("Close redis client...")
					if err = redisClient.Close(); err != nil {
						log.Errorf("Failed to close redis client: %+v", err)
						return err
					}
					log.Info("Close amqp connection...")
					if err = amqpConn.Close(); err != nil {
						log.Errorf("Failed to close amqp connection: %+v", err)
						return err
					}
					return nil
				},
			})
		}),
	).Run()
}
