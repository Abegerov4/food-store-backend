package middleware

import "net/http"

func AdminOnly(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("role") != "admin" {
			http.Error(w, "Admin only", http.StatusForbidden)
			return
		}
		next(w, r)
	}
}