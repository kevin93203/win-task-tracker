package main

import (
	"fmt"
	"net/http"

	"github.com/kevin93203/win-task-tracker/handlers"
	"github.com/kevin93203/win-task-tracker/middleware"
	"github.com/kevin93203/win-task-tracker/models"
)

func main() {
	// Initialize all model tables
	if err := models.InitDB(); err != nil {
		fmt.Printf("Failed to initialize model tables: %v\n", err)
		return
	}

	// Create handlers
	remoteComputerHandler := handlers.NewRemoteComputerHandler(models.GetDB())

	// Authentication endpoints
	http.HandleFunc("/api/register", middleware.CorsMiddleware(handlers.RegisterHandler))
	http.HandleFunc("/api/login", middleware.CorsMiddleware(handlers.LoginHandler))
	http.HandleFunc("/api/verify", middleware.CorsMiddleware(middleware.AuthMiddleware(handlers.VerifyHandler)))
	http.HandleFunc("/api/logout", middleware.CorsMiddleware(handlers.LogoutHandler))

	// Tasks endpoints
	http.HandleFunc("/api/tasks", middleware.CorsMiddleware(middleware.AuthMiddleware(handlers.GetTasksHandler)))

	// Remote computer management endpoints
	// Get method endpoints
	http.HandleFunc("/api/computers/list", middleware.CorsMiddleware(middleware.AuthMiddleware(remoteComputerHandler.HandleGetUserComputers)))
	http.HandleFunc("/api/computers/credentials", middleware.CorsMiddleware(middleware.AuthMiddleware(remoteComputerHandler.HandleGetComputerCredentials)))
	http.HandleFunc("/api/credentials/list", middleware.CorsMiddleware(middleware.AuthMiddleware(remoteComputerHandler.HandleListUserCredentials)))

	// Post method endpoints
	http.HandleFunc("/api/computers", middleware.CorsMiddleware(middleware.AuthMiddleware(remoteComputerHandler.HandleCreateRemoteComputer)))
	http.HandleFunc("/api/credentials", middleware.CorsMiddleware(middleware.AuthMiddleware(remoteComputerHandler.HandleCreateCredential)))
	http.HandleFunc("/api/computers/map-credential", middleware.CorsMiddleware(middleware.AuthMiddleware(remoteComputerHandler.HandleMapComputerCredential)))

	// Patch method endpoints
	http.HandleFunc("/api/credentials/update", middleware.CorsMiddleware(middleware.AuthMiddleware(remoteComputerHandler.HandleUpdateCredential)))
	http.HandleFunc("/api/computers/map-credential/update", middleware.CorsMiddleware(middleware.AuthMiddleware(remoteComputerHandler.HandleUpdateComputerCredentialMapping)))

	// Delete method endpoints
	http.HandleFunc("/api/computers/delete", middleware.CorsMiddleware(middleware.AuthMiddleware(remoteComputerHandler.HandleDeleteComputer)))
	http.HandleFunc("/api/credentials/delete", middleware.CorsMiddleware(middleware.AuthMiddleware(remoteComputerHandler.HandleDeleteCredential)))

	fmt.Println("Server is running on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Server failed to start: %v\n", err)
	}
}
