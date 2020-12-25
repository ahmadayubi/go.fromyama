package middleware

import (
	"context"
	"net/http"
	"strings"

	"../utils"
	"../utils/jwtUtil"
)

func ProtectedRoute(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqToken := r.Header.Get("Authorization")
		token := strings.Split(reqToken, "Bearer ")[1]
		claims, err := jwtUtil.CheckAndParseToken(token)
		if err != nil {
			utils.ForbiddenResponse(w)
			return
		}
		ctx := context.WithValue(r.Context(), "claims",claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}