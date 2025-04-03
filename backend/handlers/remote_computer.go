package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/kevin93203/win-task-tracker/models"
)

type RemoteComputerHandler struct {
	db *sql.DB
}

func NewRemoteComputerHandler(db *sql.DB) *RemoteComputerHandler {
	return &RemoteComputerHandler{db: db}
}

// CreateRemoteComputerRequest represents the request body for creating a remote computer
type CreateRemoteComputerRequest struct {
	Name         string `json:"name"`
	CredentialID int64  `json:"credential_id,omitempty"`
}

// CreateCredentialRequest represents the request body for creating a credential
type CreateCredentialRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// UpdateCredentialRequest represents the request body for updating a credential
type UpdateCredentialRequest struct {
	CredentialID int64  `json:"credential_id"`
	Password     string `json:"password"`
}

// UpdateComputerCredentialMappingRequest represents the request body for updating computer-credential mapping
type UpdateComputerCredentialMappingRequest struct {
	ComputerCredentialMappingID int64 `json:"computer_credential_mapping_id"`
	CredentialID                int64 `json:"credential_id"`
}

// MapComputerCredentialRequest represents the request body for mapping a credential to a computer
type MapComputerCredentialRequest struct {
	ComputerID   int64 `json:"computer_id"`
	CredentialID int64 `json:"credential_id"`
}

// DeleteComputerRequest represents the request body for deleting a computer
type DeleteComputerRequest struct {
	ComputerID int64 `json:"computer_id"`
}

// HandleCreateRemoteComputer handles the creation of a new remote computer
func (h *RemoteComputerHandler) HandleCreateRemoteComputer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CreateRemoteComputerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		// return json response with error message
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "request body is invalid"})
		return
	}

	if req.Name == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "no name provided"})
		return
	}
	// Get user ID from context (assuming it's set by authentication middleware)
	userID, ok := r.Context().Value("user_id").(int64)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}

	if req.CredentialID != 0 {
		// Check if the credential belongs to the user
		isCredentialOwner, err := models.CheckCredentialOwnership(h.db, req.CredentialID, userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if !isCredentialOwner {
			http.Error(w, "You do not have permission to use this credential", http.StatusForbidden)
			return
		}
	}

	computer, err := models.CreateRemoteComputer(h.db, req.Name, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// If a credential ID is provided, create the mapping
	if req.CredentialID != 0 {
		err = models.MapComputerToCredential(h.db, computer.ID, req.CredentialID, userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	json.NewEncoder(w).Encode(computer)
}

// HandleCreateCredential handles the creation of a new credential
func (h *RemoteComputerHandler) HandleCreateCredential(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CreateCredentialRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.Username == "" || req.Password == "" {
		http.Error(w, "username and password are required", http.StatusBadRequest)
		return
	}

	// Get user ID from context (assuming it's set by authentication middleware)
	userID, ok := r.Context().Value("user_id").(int64)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}

	credential, err := models.CreateCredential(h.db, req.Username, req.Password, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(credential)
}

// HandleGetUserComputers handles retrieving all computers for a user
func (h *RemoteComputerHandler) HandleGetUserComputers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := r.Context().Value("user_id").(int64)

	computers, err := models.GetComputersByUserID(h.db, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(computers)
}

// HandleGetComputerCredentials handles retrieving all credentials for a computer
func (h *RemoteComputerHandler) HandleGetComputerCredentials(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	computerID, err := strconv.ParseInt(r.URL.Query().Get("computer_id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid computer ID", http.StatusBadRequest)
		return
	}

	// Get user ID from context
	userID, ok := r.Context().Value("user_id").(int64)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}

	// Check if the computer belongs to the user
	isOwner, err := models.CheckComputerOwnership(h.db, computerID, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !isOwner {
		http.Error(w, "You do not have permission to access this computer", http.StatusForbidden)
		return
	}

	credentials, err := models.GetCredentialsByComputerID(h.db, computerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(credentials)
}

// HandleMapComputerCredential handles mapping a credential to a computer
func (h *RemoteComputerHandler) HandleMapComputerCredential(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req MapComputerCredentialRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Get user ID from context
	userID, ok := r.Context().Value("user_id").(int64)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}

	// Check if the computer belongs to the user
	isComputerOwner, err := models.CheckComputerOwnership(h.db, req.ComputerID, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !isComputerOwner {
		http.Error(w, "You do not have permission to access this computer", http.StatusForbidden)
		return
	}

	// Check if the credential belongs to the user
	isCredentialOwner, err := models.CheckCredentialOwnership(h.db, req.CredentialID, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !isCredentialOwner {
		http.Error(w, "You do not have permission to use this credential", http.StatusForbidden)
		return
	}

	err = models.MapComputerToCredential(h.db, req.ComputerID, req.CredentialID, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Computer credential mapping created successfully"})
}

// HandleDeleteComputer handles deleting a computer and its mappings
func (h *RemoteComputerHandler) HandleDeleteComputer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req DeleteComputerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Get user ID from context
	userID, ok := r.Context().Value("user_id").(int64)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}

	// Check if the computer belongs to the user
	isOwner, err := models.CheckComputerOwnership(h.db, req.ComputerID, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !isOwner {
		http.Error(w, "You do not have permission to delete this computer", http.StatusForbidden)
		return
	}

	err = models.DeleteComputer(h.db, req.ComputerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Computer deleted successfully"})
}

// HandleListUserCredentials handles retrieving all credentials for a user
func (h *RemoteComputerHandler) HandleListUserCredentials(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := r.Context().Value("user_id").(int64)

	credentials, err := models.GetCredentialsByUserID(h.db, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(credentials)
}

// HandleDeleteCredential handles deleting a credential
func (h *RemoteComputerHandler) HandleDeleteCredential(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	credentialID, err := strconv.ParseInt(r.URL.Query().Get("credential_id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid credential ID", http.StatusBadRequest)
		return
	}

	userID := r.Context().Value("user_id").(int64)

	// Check if the credential belongs to the user
	isOwner, err := models.CheckCredentialOwnership(h.db, credentialID, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !isOwner {
		http.Error(w, "You do not have permission to delete this credential", http.StatusForbidden)
		return
	}

	err = models.DeleteCredential(h.db, credentialID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// HandleUpdateCredential handles updating a credential's password
func (h *RemoteComputerHandler) HandleUpdateCredential(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req UpdateCredentialRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.Password == "" {
		http.Error(w, "password is required", http.StatusBadRequest)
		return
	}

	userID := r.Context().Value("user_id").(int64)

	// Check if the credential belongs to the user
	isOwner, err := models.CheckCredentialOwnership(h.db, req.CredentialID, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !isOwner {
		http.Error(w, "You do not have permission to update this credential", http.StatusForbidden)
		return
	}

	err = models.UpdateCredentialPassword(h.db, req.CredentialID, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Credential updated successfully"})
}

func (h *RemoteComputerHandler) HandleUpdateComputerCredentialMapping(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req UpdateComputerCredentialMappingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID := r.Context().Value("user_id").(int64)

	// Check if the computer credentail mapping belongs to the user
	isOwner, err := models.CheckComputerCredentialMappingOwnership(h.db, req.ComputerCredentialMappingID, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !isOwner {
		http.Error(w, "You do not have permission to update this computer's credential mapping", http.StatusForbidden)
		return
	}

	// Check if the credential belongs to the user
	isCredentialOwner, err := models.CheckCredentialOwnership(h.db, req.CredentialID, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !isCredentialOwner {
		http.Error(w, "You do not have permission to use this credential", http.StatusForbidden)
		return
	}

	err = models.UpdateComputerCredentialMapping(h.db, req.ComputerCredentialMappingID, req.CredentialID, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
