package middleware

import (
	"github.com/gorilla/sessions"
	"net/http"
	"regexp"
	"xy-dianping-go/internal/common"
	"xy-dianping-go/internal/dto"
)

var (
	excludePathPatterns = []*regexp.Regexp{
		regexp.MustCompile("^/shop/.*"),
		regexp.MustCompile("^/voucher/.*"),
		regexp.MustCompile("^/shop-type/.*"),
		regexp.MustCompile("^/upload/.*"),
		regexp.MustCompile("^/blog/hot/?$"),
		regexp.MustCompile("^/user/code/?$"),
		regexp.MustCompile("^/user/login/?$"),
	}
)

func AuthenticateMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 遍历不需要登录的请求列表
		for _, p := range excludePathPatterns {
			isMatch := p.MatchString(r.URL.Path)
			if isMatch { // 匹配，直接放行
				next.ServeHTTP(w, r)
				return
			}
		}

		// 获取 session
		session := r.Context().Value("session").(*sessions.Session)
		// 获取 session 中的用户
		user, ok := session.Values["user"]

		// 判断是否需要进行拦截（是否存在用户）
		if !ok {
			// 没有用户，需要拦截，设置状态码
			common.SendResponseWithCode(w, common.Fail("未登录"), 401)
			// 拦截
			return
		}
		// 存在，保存用户信息到 ctx
		ctx := common.SetUserToContext(r.Context(), user.(*dto.UserDTO))
		// 则进行放行
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
