package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"win-task-tracker/backend/models"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte("your-secret-key") // In production, use environment variable

type Claims struct {
	UserID   int    `json:"sub"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func RegisterUser(username, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	db := models.GetDB()

	_, err = db.Exec("INSERT INTO users (username, password) VALUES (?, ?)",
		username, string(hashedPassword))
	return err
}

func LoginUser(username, password string) (string, time.Time, error) {
	var user models.User

	db := models.GetDB()
	err := db.QueryRow("SELECT id, username, password FROM users WHERE username = ?",
		username).Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		return "", time.Time{}, errors.New("invalid credentials")
	}

	// Compare password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", time.Time{}, errors.New("invalid credentials")
	}

	// Set expiration time
	expTime := time.Now().Add(time.Hour * 24)

	// Generate JWT token
	claims := &Claims{
		UserID:   user.ID,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expTime, nil
}

func VerifyToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
