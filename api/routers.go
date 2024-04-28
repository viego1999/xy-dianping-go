package api

import (
	"github.com/gorilla/mux"
	"go.uber.org/fx"
	"net/http"
	"xy-dianping-go/api/v1"
	"xy-dianping-go/internal/middleware"
)

var (
	// Module Router 注入模块
	Module = fx.Options(fx.Provide(Router))
)

// Router 路由统一管理，返回 *mux.Router
func Router(userController *v1.UserController, shopController *v1.ShopController,
	voucherController *v1.VoucherController, voucherOrderController *v1.VoucherOrderController,
	shopTypeController *v1.ShopTypeController, blogController *v1.BlogController) *mux.Router {
	// 路由注册
	router := mux.NewRouter()

	// 全局中间件
	router.Use(middleware.RecoverMiddleware)
	// Session 中间件
	router.Use(middleware.RefreshTokenMiddleware)
	// 登录中间件
	router.Use(middleware.LoginMiddleware)

	// 注册路由
	router.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		_, err := writer.Write([]byte("hello world"))
		if err != nil {
			return
		}
	}).Methods("GET")

	// 注册 user 子路由器
	userRouter := router.PathPrefix("/user").Subrouter()
	userRouter.HandleFunc("/code", userController.SendCode).Methods("POST")
	userRouter.HandleFunc("/login", userController.Login).Methods("POST")
	userRouter.HandleFunc("/me", userController.Me).Methods("GET")
	userRouter.HandleFunc("/info/{id}", userController.Info).Methods("GET")
	userRouter.HandleFunc("/{id}", userController.QueryUserById).Methods("GET")
	userRouter.HandleFunc("/sign", userController.Sign).Methods("POST")
	userRouter.HandleFunc("/sign/count", userController.SignCount).Methods("GET")

	// 注册 shop 子路由器
	shopRouter := router.PathPrefix("/shop").Subrouter()
	shopRouter.HandleFunc("/{id}", shopController.QueryShopById).Methods("GET")
	shopRouter.HandleFunc("", shopController.SaveShop).Methods("POST")
	shopRouter.HandleFunc("", shopController.UpdateShop).Methods("PUT")
	shopRouter.HandleFunc("/of/type", shopController.QueryShopByType).Methods("GET")
	shopRouter.HandleFunc("/of/name", shopController.QueryShopByName).Methods("GET")

	//注册 shopType 子路由器
	shopTypeRouter := router.PathPrefix("/shop-type").Subrouter()
	shopTypeRouter.HandleFunc("/list", shopTypeController.QueryTypeList).Methods("GET")

	// 注册 voucher 子路由器
	voucherRouter := router.PathPrefix("/voucher").Subrouter()
	voucherRouter.HandleFunc("", voucherController.AddVoucher).Methods("POST")
	voucherRouter.HandleFunc("/seckill", voucherController.AddSeckillVoucher).Methods("POST")
	voucherRouter.HandleFunc("/list/{shopId}", voucherController.QueryVoucherOfShop).Methods("GET")

	// 注册 voucherOrder 子路由器
	voucherOrderRouter := router.PathPrefix("/voucher-order").Subrouter()
	voucherOrderRouter.HandleFunc("/seckill/{id}", voucherOrderController.SeckillVoucher).Methods("POST")

	// 注册 blog 子路由器
	blogRouter := router.PathPrefix("/blog").Subrouter()
	blogRouter.HandleFunc("", blogController.SaveBlog).Methods("POST")
	blogRouter.HandleFunc("/like/{id}", blogController.LikeBlog).Methods("PUT")
	blogRouter.HandleFunc("/of/me", blogController.QueryMyBlog).Methods("GET")
	blogRouter.HandleFunc("/hot", blogController.QueryHotBlog).Methods("GET")
	blogRouter.HandleFunc("/{id}", blogController.QueryBlogById).Methods("GET")
	blogRouter.HandleFunc("/likes/{id}", blogController.QueryBlogLikes).Methods("GET")
	blogRouter.HandleFunc("/of/user", blogController.QueryBlogByUserId).Methods("GET")
	blogRouter.HandleFunc("/of/follow", blogController.QueryBlogOfFollow).Methods("GET")

	// 注册

	return router
}
