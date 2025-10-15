package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"

	"win-task-tracker/backend/models"
)

// StartTaskHandler 處理啟動排程任務的請求
func StartTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get user ID from context (set by auth middleware)
	userID := r.Context().Value("user_id").(int64)

	// Parse request body
	var req StartTaskRequest
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

	// 使用 PowerShell 腳本啟動任務
	cmd := fmt.Sprintf(
		"powershell -File ./scripts/startTask.ps1 "+
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

	response := StartTaskResponse{}
	if err != nil {
		response.Success = false
		response.Error = fmt.Sprintf("Failed to start task: %s - %s", err.Error(), string(output))
	} else {
		success, message := parsePowerShellOutput(output, fmt.Sprintf("Successfully started task '%s' on computer '%s'", req.TaskName, targetHost.ComputerName))
		response.Success = success

		if success {
			response.Message = message

			// 解析狀態信息
			var result map[string]interface{}
			if json.Unmarshal(output, &result) == nil {
				if state, ok := result["State"].(string); ok {
					response.State = state
				}
			} else {
				// 嘗試解析為數組
				var resultArray []interface{}
				if json.Unmarshal(output, &resultArray) == nil && len(resultArray) > 0 {
					if obj, ok := resultArray[0].(map[string]interface{}); ok {
						if state, ok := obj["State"].(string); ok {
							response.State = state
						}
					}
				}
			}
		} else {
			response.Error = message
		}
	}

	sendAPIResponse(w, response.Success, response)
}

// StopTaskHandler 處理停止排程任務的請求
func StopTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get user ID from context (set by auth middleware)
	userID := r.Context().Value("user_id").(int64)

	// Parse request body
	var req StopTaskRequest
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

	// 使用 PowerShell 腳本停止任務
	cmd := fmt.Sprintf(
		"powershell -File ./scripts/stopTask.ps1 "+
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

	response := StopTaskResponse{}
	if err != nil {
		response.Success = false
		response.Error = fmt.Sprintf("Failed to stop task: %s - %s", err.Error(), string(output))
	} else {
		success, message := parsePowerShellOutput(output, fmt.Sprintf("Successfully stopped task '%s' on computer '%s'", req.TaskName, targetHost.ComputerName))
		response.Success = success

		if success {
			response.Message = message

			// 解析狀態信息
			var result map[string]interface{}
			if json.Unmarshal(output, &result) == nil {
				if state, ok := result["State"].(string); ok {
					response.State = state
				}
			} else {
				// 嘗試解析為數組
				var resultArray []interface{}
				if json.Unmarshal(output, &resultArray) == nil && len(resultArray) > 0 {
					if obj, ok := resultArray[0].(map[string]interface{}); ok {
						if state, ok := obj["State"].(string); ok {
							response.State = state
						}
					}
				}
			}
		} else {
			response.Error = message
		}
	}

	sendAPIResponse(w, response.Success, response)
}
