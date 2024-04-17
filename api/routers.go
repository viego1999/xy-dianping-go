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
func Router(userController *v1.UserController, blogController *v1.BlogController) *mux.Router {
	// 路由注册
	router := mux.NewRouter()

	// 全局中间件
	router.Use(middleware.RecoverMiddleware)
	// Session 中间件
	router.Use(middleware.SessionMiddleware)

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

	// 注册 blog 子路由器
	blogRouter := router.PathPrefix("/blog").Subrouter()
	blogRouter.HandleFunc("/{id}", blogController.QueryBlogById).Methods("GET")

	//

	return router
}
