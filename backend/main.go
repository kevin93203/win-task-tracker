package main

import (
	"fmt"
	"net/http"

	"github.com/kevin93203/win-task-tracker/auth"
	"github.com/kevin93203/win-task-tracker/handlers"
	"github.com/kevin93203/win-task-tracker/middleware"
)

func main() {
	// Initialize the database
	if err := auth.InitDB(); err != nil {
		fmt.Printf("Failed to initialize database: %v\n", err)
		return
	}

	// Authentication endpoints
	http.HandleFunc("/api/register", handlers.RegisterHandler)
	http.HandleFunc("/api/login", handlers.LoginHandler)

	// Protected endpoints with JWT authentication
	http.HandleFunc("/api/tasks", middleware.AuthMiddleware(handlers.GetTasksHandler))

	fmt.Println("Server is running on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Server failed to start: %v\n", err)
	}
}
