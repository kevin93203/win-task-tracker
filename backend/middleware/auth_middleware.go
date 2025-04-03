package middleware

import (
	"net/http"

	"github.com/golang-jwt/jwt/v4"
)

var jwtSecret = []byte("your-secret-key") // Should match the secret in auth package

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get token from cookie
		cookie, err := r.Cookie("jwt")
		if err != nil {
			http.Error(w, "Unauthorized: No JWT cookie", http.StatusUnauthorized)
			return
		}

		tokenString := cookie.Value
		if tokenString == "" {
			http.Error(w, "Unauthorized: Empty JWT", http.StatusUnauthorized)
			return
		}

		// Parse and validate the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.NewValidationError("invalid signing method", jwt.ValidationErrorSignatureInvalid)
			}
			return jwtSecret, nil
		})

		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		if !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Token is valid, call the next handler
		next(w, r)
	}
}
