package middleware

import (
	"net/http"
)

func CorsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 檢查請求來源
		origin := r.Header.Get("Origin")
		if origin == "http://localhost:5173" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}

		// 允許認證資訊
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// 允許的方法和標頭
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
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
