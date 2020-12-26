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
		if reqToken == ""{
			utils.ForbiddenResponse(w)
			return
		}
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

func ProtectedApprovedUserRoute(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqToken := r.Header.Get("Authorization")
		if reqToken == ""{
			utils.ForbiddenResponse(w)
			return
		}
		token := strings.Split(reqToken, "Bearer ")[1]
		claims, err := jwtUtil.CheckAndParseToken(token)
		if err != nil {
			utils.ForbiddenResponse(w)
			return
		}
		if !claims.Approved {
			utils.JSONResponse(w, http.StatusUnauthorized, "Not Approved")
			return
		}
		ctx := context.WithValue(r.Context(), "claims",claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}