package auth

import (
    "database/sql"
    "time"
    "errors"
    "golang.org/x/crypto/bcrypt"
    "github.com/golang-jwt/jwt/v4"
    _ "github.com/mattn/go-sqlite3"
    "github.com/kevin93203/win-task-tracker/models"
)

var db *sql.DB
var jwtSecret = []byte("your-secret-key") // In production, use environment variable

type Claims struct {
    UserID   int    `json:"user_id"`
    Username string `json:"username"`
    jwt.RegisteredClaims
}

func InitDB() error {
    var err error
    db, err = sql.Open("sqlite3", "./users.db")
    if err != nil {
        return err
    }

    // Create users table if not exists
    createTable := `
    CREATE TABLE IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        username TEXT UNIQUE NOT NULL,
        password TEXT NOT NULL
    );`

    _, err = db.Exec(createTable)
    return err
}

func RegisterUser(username, password string) error {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return err
    }

    _, err = db.Exec("INSERT INTO users (username, password) VALUES (?, ?)", 
        username, string(hashedPassword))
    return err
}

func LoginUser(username, password string) (string, time.Time, error) {
    var user models.User
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
