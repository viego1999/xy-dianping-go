package middleware

import (
	"net/http"
	"time"
	"xy-dianping-go/internal/common"
	"xy-dianping-go/internal/constants"
	"xy-dianping-go/internal/db"
	"xy-dianping-go/internal/dto"
	"xy-dianping-go/pkg/utils"
)

func RefreshTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1.获取请求头中的 token
		token := r.Header.Get("authorization")
		if token == "" {
			next.ServeHTTP(w, r)
			return
		}

		// 2.基于 token 获取 redis 中的用户
		key := constants.LOGIN_USER_KEY + token
		// 获取哈希表中所有字段
		userMap, err := db.RedisClient.HGetAll(r.Context(), key).Result()
		if err != nil {
			panic("Redis HGetAll error:" + err.Error())
		}
		// 3.判断用户是否存在
		if userMap == nil || len(userMap) == 0 {
			// 4.不存在，放行
			next.ServeHTTP(w, r)
			return
		}
		userDTO := dto.UserDTO{}
		// 5.将查询到的数据转化为 userDTO
		err = utils.MapToStruct(utils.MapValueToAny(userMap), &userDTO)
		if err != nil {
			panic("MapToStruct error:" + err.Error())
		}

		// 6.保存信息到 ThreadLocal 中
		ctx := common.SetUserToContext(r.Context(), &userDTO)
		// 7.刷新 token 有效期
		db.RedisClient.Expire(r.Context(), key, constants.LOGIN_USER_TTL*time.Second)
		// 8.放行
		next.ServeHTTP(w, r.WithContext(ctx))
		return
	})
}
