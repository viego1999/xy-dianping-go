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
func Router(userController *v1.UserController, shopController *v1.ShopController, blogController *v1.BlogController) *mux.Router {
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

	// 注册 blog 子路由器
	blogRouter := router.PathPrefix("/blog").Subrouter()
	blogRouter.HandleFunc("/{id}", blogController.QueryBlogById).Methods("GET")

	//

	return router
}
