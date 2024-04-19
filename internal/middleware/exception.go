package middleware

import (
	"fmt"
	"net/http"
	"runtime"

	"trpc.group/trpc-go/trpc-go/log"
	v1 "xy-dianping-go/internal/common"
)

// RecoverMiddleware 捕获处理器中的panic，并返回500内部服务器错误。
func RecoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				buf := make([]byte, 1024)
				n := runtime.Stack(buf, false)
				stackInfo := string(buf[:n])
				log.Errorf("Recovered from panic: %+v", err)
				log.Error("Stack trace:")
				log.Error(stackInfo)

				// 返回内部错误响应
				v1.SendResponseWithCode(w, v1.Fail(fmt.Sprintf("请求失败：%v", err)), http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
