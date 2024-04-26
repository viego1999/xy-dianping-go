//Package middleware 提供一系列 HTTP 服务的中间件组件。
/*
主要包括了 LoginMiddleware， RecoverMiddleware， RefreshTokenMiddleware， SessionMiddleware

Example：

// 导入 mux 包

import "github.com/gorilla/mux"

// 路由注册

router := mux.NewRouter()

// Panic 处理中间件

router.Use(middleware.RecoverMiddleware)
*/
package middleware
