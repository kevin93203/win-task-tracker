package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"win-task-tracker/backend/handlers"
	"win-task-tracker/backend/middleware"
	"win-task-tracker/backend/models"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize all model tables
	if err := models.InitDB(); err != nil {
		fmt.Printf("Failed to initialize model tables: %v\n", err)
		return
	}

	// Create handlers
	remoteComputerHandler := handlers.NewRemoteComputerHandler(models.GetDB())

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port if not specified
	}

	// Authentication endpoints
	http.HandleFunc("/api/register", middleware.CorsMiddleware(handlers.RegisterHandler))
	http.HandleFunc("/api/login", middleware.CorsMiddleware(handlers.LoginHandler))
	http.HandleFunc("/api/verify", middleware.CorsMiddleware(middleware.AuthMiddleware(handlers.VerifyHandler)))
	http.HandleFunc("/api/logout", middleware.CorsMiddleware(handlers.LogoutHandler))

	// Tasks endpoints
	http.HandleFunc("/api/tasks", middleware.CorsMiddleware(middleware.AuthMiddleware(handlers.GetTasksHandler)))
	http.HandleFunc("/api/tasks/", middleware.CorsMiddleware(middleware.AuthMiddleware(handlers.GetTaskHandler)))
	http.HandleFunc("/api/tasks/disable", middleware.CorsMiddleware(middleware.AuthMiddleware(handlers.DisableTaskHandler)))
	http.HandleFunc("/api/tasks/enable", middleware.CorsMiddleware(middleware.AuthMiddleware(handlers.EnableTaskHandler)))
	http.HandleFunc("/api/tasks/start", middleware.CorsMiddleware(middleware.AuthMiddleware(handlers.StartTaskHandler)))
	http.HandleFunc("/api/tasks/stop", middleware.CorsMiddleware(middleware.AuthMiddleware(handlers.StopTaskHandler)))

	// Trigger management endpoints
	http.HandleFunc("/api/tasks/triggers", middleware.CorsMiddleware(middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handlers.AddTriggerHandler(w, r)
		case http.MethodPatch:
			handlers.UpdateTriggerHandler(w, r)
		case http.MethodDelete:
			handlers.DeleteTriggerHandler(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	// Action management endpoints
	http.HandleFunc("/api/tasks/actions", middleware.CorsMiddleware(middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handlers.AddActionHandler(w, r)
		case http.MethodPatch:
			handlers.UpdateActionHandler(w, r)
		case http.MethodDelete:
			handlers.DeleteActionHandler(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	// Remote computer management endpoints
	// Get method endpoints
	http.HandleFunc("/api/computers/list", middleware.CorsMiddleware(middleware.AuthMiddleware(remoteComputerHandler.HandleGetUserComputers)))
	http.HandleFunc("/api/computers/credentials", middleware.CorsMiddleware(middleware.AuthMiddleware(remoteComputerHandler.HandleGetComputerCredentials)))
	http.HandleFunc("/api/credentials/list", middleware.CorsMiddleware(middleware.AuthMiddleware(remoteComputerHandler.HandleListUserCredentials)))
	http.HandleFunc("/api/computers/credential-mappings", middleware.CorsMiddleware(middleware.AuthMiddleware(remoteComputerHandler.HandleGetComputerCredentialMappings)))

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

	fmt.Printf("Server is running on :%s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		fmt.Printf("Server failed to start: %v\n", err)
	}
}
