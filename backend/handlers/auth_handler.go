package handlers

import (
    "encoding/json"
    "net/http"
    "github.com/kevin93203/win-task-tracker/auth"
)

type LoginRequest struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

type LoginResponse struct {
    Token string `json:"token"`
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var req LoginRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    if err := auth.RegisterUser(req.Username, req.Password); err != nil {
        http.Error(w, "Registration failed", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var req LoginRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    token, err := auth.LoginUser(req.Username, req.Password)
    if err != nil {
        http.Error(w, "Login failed", http.StatusUnauthorized)
        return
    }

    response := LoginResponse{Token: token}
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}
