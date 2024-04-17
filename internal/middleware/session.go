package middleware

import (
	"context"
	"github.com/gorilla/sessions"
	"net/http"
	"xy-dianping-go/internal/common"
)

var (
	// Store 初始化 session 存储
	Store = sessions.NewCookieStore([]byte("xy-dianping-go"))
	// SessionName session name
	SessionName = "wxy"
)

// SessionMiddleware 是一个中间件，用于将 session 存储添加到请求上下文中
func SessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := Store.Get(r, SessionName)
		if err != nil {
			common.SendResponseWithCode(w, common.Fail(err.Error()), http.StatusInternalServerError)
			return
		}
		// 将 session 存储到上下文请求
		ctx := context.WithValue(r.Context(), "session", session)
		// 调用下一个处理函数
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
