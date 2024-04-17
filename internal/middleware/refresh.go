package middleware

import (
	"net/http"
)

func RefreshTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("authorization")
		if token == "" {
			next.ServeHTTP(w, r)
			return
		}
	})
}
