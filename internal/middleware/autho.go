package middleware

import (
	"context"
	"net/http"
	"xy-dianping-go/internal/dto"
)

type userCtxKey struct{}

func NewUserContext(ctx context.Context, user dto.UserDTO) context.Context {
	return context.WithValue(ctx, userCtxKey{}, user)
}

func FromUserContext(ctx context.Context) (dto.UserDTO, bool) {
	user, ok := ctx.Value(userCtxKey{}).(dto.UserDTO)
	return user, ok
}

func AuthenticateMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("authorization")
		if token == "" {
			next.ServeHTTP(w, r)
			return
		}

		user := dto.UserDTO{
			Id:       1,
			NickName: "Jack",
		}

		ctx := NewUserContext(r.Context(), user)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
