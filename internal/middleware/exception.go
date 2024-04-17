package middleware

import (
	"fmt"
	"net/http"
	v1 "xy-dianping-go/internal/common"

	"trpc.group/trpc-go/trpc-go/log"
)

// RecoverMiddleware 捕获处理器中的panic，并返回500内部服务器错误。
func RecoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Warnf("Recovered from panic: %+v", err)

				v1.SendResponseWithCode(w, v1.Fail(fmt.Sprintf("%v", err)), http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
