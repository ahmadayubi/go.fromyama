package middleware

import (
	"context"
	"net/http"
	"strings"

	"go.fromyama/utils/jwtUtil"
	"go.fromyama/utils/response"
)

func ProtectedRoute(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqToken := r.Header.Get("Authorization")
		splitHeader := strings.Split(reqToken, "Bearer ")
		if reqToken == "" || len(splitHeader) != 2{
			response.Forbidden(w)
			return
		}

		token := splitHeader[1]
		claims, err := jwtUtil.CheckAndParseToken(token)
		if err != nil {
			response.Forbidden(w)
		} else {
			ctx := context.WithValue(r.Context(), "claims", *claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	})
}

func ProtectedApprovedUserRoute(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqToken := r.Header.Get("Authorization")
		splitHeader := strings.Split(reqToken, "Bearer ")
		if reqToken == "" || len(splitHeader) != 2{
			response.Forbidden(w)
			return
		}

		token := splitHeader[1]
		claims, err := jwtUtil.CheckAndParseToken(token)
		if err != nil {
			response.Forbidden(w)
		} else if !claims.Approved {
			response.JSON(w, http.StatusUnauthorized, response.BasicMessage{
				Message: "Not Approved",
			})
		} else {
			ctx := context.WithValue(r.Context(), "claims", *claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	})
}