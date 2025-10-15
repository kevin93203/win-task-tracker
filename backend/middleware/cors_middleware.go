package middleware

import (
	"net/http"
	"os"
	"strings"
)

func CorsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	allowedOriginsStr := os.Getenv("ALLOWED_ORIGINS")
	var allowedOrigins []string
	if allowedOriginsStr != "" {
		allowedOrigins = strings.Split(allowedOriginsStr, ",")
	}

	return func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		isAllowed := false
		for _, allowed := range allowedOrigins {
			if origin == allowed {
				isAllowed = true
				break
			}
		}

		if isAllowed {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}

		// 允許認證資訊
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// 允許的方法和標頭
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, PATCH,DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		// 暴露 Authorization 標頭
		w.Header().Set("Access-Control-Expose-Headers", "Authorization")

		// 處理 OPTIONS 請求
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next(w, r)
	}
}
