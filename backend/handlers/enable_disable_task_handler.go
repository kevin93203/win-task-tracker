package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"

	"win-task-tracker/backend/models"
)

// DisableTaskHandler 處理停用排程任務的請求
func DisableTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get user ID from context (set by auth middleware)
	userID := r.Context().Value("user_id").(int64)

	// Parse request body
	var req DisableTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Get database connection
	db := models.GetDB()

	// 獲取任務憑證
	targetHost, err := getTaskCredentials(db, userID, req.ComputerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 使用 PowerShell 腳本停用任務
	cmd := fmt.Sprintf(
		"powershell -File ./scripts/disableTask.ps1 "+
			"-TaskName '%s' "+
			"-UserName '%s' "+
			"-Password '%s' "+
			"-ComputerName '%s'",
		req.TaskName,
		targetHost.UserName,
		targetHost.Password,
		targetHost.ComputerName,
	)

	// Execute the command
	powershell := exec.Command("powershell", "-Command", cmd)
	output, err := powershell.CombinedOutput()

	response := DisableTaskResponse{}
	if err != nil {
		response.Success = false
		response.Error = fmt.Sprintf("Failed to disable task: %s - %s", err.Error(), string(output))
	} else {
		success, message := parsePowerShellOutput(output, fmt.Sprintf("Successfully disabled task '%s' on computer '%s'", req.TaskName, targetHost.ComputerName))
		response.Success = success
		if success {
			response.Message = message
		} else {
			response.Error = message
		}
	}

	sendAPIResponse(w, response.Success, response)
}

// EnableTaskHandler 處理啟用排程任務的請求
func EnableTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get user ID from context (set by auth middleware)
	userID := r.Context().Value("user_id").(int64)

	// Parse request body
	var req EnableTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Get database connection
	db := models.GetDB()

	// 獲取任務憑證
	targetHost, err := getTaskCredentials(db, userID, req.ComputerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 使用 PowerShell 腳本啟用任務
	cmd := fmt.Sprintf(
		"powershell -File ./scripts/enableTask.ps1 "+
			"-TaskName '%s' "+
			"-UserName '%s' "+
			"-Password '%s' "+
			"-ComputerName '%s'",
		req.TaskName,
		targetHost.UserName,
		targetHost.Password,
		targetHost.ComputerName,
	)

	// Execute the command
	powershell := exec.Command("powershell", "-Command", cmd)
	output, err := powershell.CombinedOutput()

	response := EnableTaskResponse{}
	if err != nil {
		response.Success = false
		response.Error = fmt.Sprintf("Failed to enable task: %s - %s", err.Error(), string(output))
	} else {
		success, message := parsePowerShellOutput(output, fmt.Sprintf("Successfully enabled task '%s' on computer '%s'", req.TaskName, targetHost.ComputerName))
		response.Success = success
		if success {
			response.Message = message
		} else {
			response.Error = message
		}
	}

	sendAPIResponse(w, response.Success, response)
}
