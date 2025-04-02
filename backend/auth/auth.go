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

func LoginUser(username, password string) (string, error) {
    var user models.User
    err := db.QueryRow("SELECT id, username, password FROM users WHERE username = ?", 
        username).Scan(&user.ID, &user.Username, &user.Password)
    if err != nil {
        return "", errors.New("invalid credentials")
    }

    // Compare password
    err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
    if err != nil {
        return "", errors.New("invalid credentials")
    }

    // Generate JWT token
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user_id":  user.ID,
        "username": user.Username,
        "exp":      time.Now().Add(time.Hour * 24).Unix(),
    })

    tokenString, err := token.SignedString(jwtSecret)
    if err != nil {
        return "", err
    }

    return tokenString, nil
}
