package middleware

import (
	"net/http"
)

// 定义上下文键,创建一个新的上下文键类型，防止上下文键冲突
var userContextKey struct{}

func RefreshTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("authorization")
		if token == "" {
			next.ServeHTTP(w, r)
			return
		}
	})
}
