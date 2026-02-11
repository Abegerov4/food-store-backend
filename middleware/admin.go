package middleware

import (
	"net/http"

	"food-store-backend/utils"
)


func AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims := r.Context().Value(UserContextKey)
		if claims == nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		userClaims := claims.(*utils.Claims)
		if userClaims.Role != "admin" {
			http.Error(w, "Forbidden (admin only)", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}