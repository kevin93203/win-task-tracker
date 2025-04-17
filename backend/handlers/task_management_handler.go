package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"strings"

	"github.com/kevin93203/win-task-tracker/models"
)

// TriggerRequest represents the payload for trigger operations
type TriggerRequest struct {
	ComputerID     int64  `json:"computer_id"`
	TaskName       string `json:"task_name"`
	CronExpression string `json:"cron_expression"` // Cron expression string
	Index          int    `json:"index,omitempty"` // Trigger index to update/delete
}

// ActionRequest represents the payload for action operations
type ActionRequest struct {
	ComputerID       int64  `json:"computer_id"`
	TaskName         string `json:"task_name"`
	Execute          string `json:"execute,omitempty"`
	Args             string `json:"args,omitempty"`
	WorkingDirectory string `json:"working_directory,omitempty"`
	Index            int    `json:"index,omitempty"` // Action index to update/delete
}

// GenericResponse standard API response
type GenericResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

// AddTriggerHandler adds a new trigger to a task
func AddTriggerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req TriggerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	db := models.GetDB()
	// Get user ID from context (set by auth middleware)
	userID := r.Context().Value("user_id").(int64)

	host, err := getTaskCredentials(db, userID, req.ComputerID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get remote host info: %v", err), http.StatusInternalServerError)
		return
	}

	// Convert Cron expression to PowerShell trigger params
	var triggerTypeAdd, triggerTimeAdd string
	partsAdd := strings.Fields(req.CronExpression)
	if len(partsAdd) != 5 {
		http.Error(w, "Invalid cron expression format", http.StatusBadRequest)
		return
	}

	minute := partsAdd[0]
	hour := partsAdd[1]
	dayOfMonth := partsAdd[2]
	month := partsAdd[3]
	dayOfWeek := partsAdd[4]

	// 時間格式為 HH:MM (24小時制，確保兩位數)
	// 填充前導零以確保格式正確
	hourFormatted := fmt.Sprintf("%02s", hour)
	minuteFormatted := fmt.Sprintf("%02s", minute)
	triggerTimeAdd = fmt.Sprintf("%s:%s", hourFormatted, minuteFormatted)

	// 根據 cron 表達式轉換為 SCHTASKS 支持的格式
	switch {
	case dayOfMonth == "*" && month == "*" && dayOfWeek == "*":
		triggerTypeAdd = "Daily"
	case dayOfMonth == "*" && month == "*" && dayOfWeek != "*":
		triggerTypeAdd = "Weekly"
		// Map dayOfWeek number/string to day name
		days := map[string]string{
			"0": "SUN", "1": "MON", "2": "TUE", "3": "WED",
			"4": "THU", "5": "FRI", "6": "SAT",
			"7": "SUN",
		}
		triggerTypeAdd += fmt.Sprintf(" -DaysOfWeek %s", days[dayOfWeek])
	case dayOfMonth != "*" && month == "*" && dayOfWeek == "*":
		triggerTypeAdd = fmt.Sprintf("Monthly -DaysOfMonth %s", dayOfMonth)
	case dayOfMonth != "*" && month != "*" && dayOfWeek == "*":
		triggerTypeAdd = fmt.Sprintf("Monthly -DaysOfMonth %s -MonthsOfYear %s", dayOfMonth, month)
	default:
		http.Error(w, "Unsupported cron pattern for trigger conversion", http.StatusBadRequest)
		return
	}

	// 使用 PowerShell 腳本
	cmd := exec.Command("powershell", "-File", "./scripts/addTriggers.ps1",
		"-TaskName", req.TaskName,
		"-TriggerType", triggerTypeAdd,
		"-TriggerTime", triggerTimeAdd,
		"-UserName", host.UserName,
		"-Password", host.Password,
		"-ComputerName", host.ComputerName,
	)

	// Execute the command
	output, err := cmd.CombinedOutput()
	var resp GenericResponse
	if err != nil {
		resp = GenericResponse{Success: false, Message: fmt.Sprintf("Error: %v, Output: %s", err, string(output))}
	} else {
		success, message := parsePowerShellOutput(output, "Trigger added successfully")
		resp = GenericResponse{Success: success, Message: message}
	}
	sendAPIResponse(w, resp.Success, resp)
}

// UpdateTriggerHandler updates an existing trigger
func UpdateTriggerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req TriggerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	db := models.GetDB()
	userID := r.Context().Value("user_id").(int64)
	host, err := getTaskCredentials(db, userID, req.ComputerID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get remote host info: %v", err), http.StatusInternalServerError)
		return
	}
	// Convert Cron expression to PowerShell trigger params
	var triggerType, triggerTime string
	parts := strings.Fields(req.CronExpression)
	if len(parts) != 5 {
		http.Error(w, "Invalid cron expression format", http.StatusBadRequest)
		return
	}

	minute := parts[0]
	hour := parts[1]
	dayOfMonth := parts[2]
	month := parts[3]
	dayOfWeek := parts[4]

	// 時間格式為 HH:MM (24小時制，確保兩位數)
	// 填充前導零以確保格式正確
	hourFormatted := fmt.Sprintf("%02s", hour)
	minuteFormatted := fmt.Sprintf("%02s", minute)
	triggerTime = fmt.Sprintf("%s:%s", hourFormatted, minuteFormatted)

	switch {
	case dayOfMonth == "*" && month == "*" && dayOfWeek == "*":
		triggerType = "Daily"
	case dayOfMonth == "*" && month == "*" && dayOfWeek != "*":
		triggerType = "Weekly"
		days := map[string]string{
			"0": "Sunday", "1": "Monday", "2": "Tuesday", "3": "Wednesday",
			"4": "Thursday", "5": "Friday", "6": "Saturday",
			"7": "Sunday",
		}
		triggerType += fmt.Sprintf(" -DaysOfWeek %s", days[dayOfWeek])
	case dayOfMonth != "*" && month == "*" && dayOfWeek == "*":
		triggerType = fmt.Sprintf("Monthly -DaysOfMonth %s", dayOfMonth)
	case dayOfMonth != "*" && month != "*" && dayOfWeek == "*":
		triggerType = fmt.Sprintf("Monthly -MonthsOfYear %s -DaysOfMonth %s", month, dayOfMonth)
	default:
		http.Error(w, "Unsupported cron pattern for trigger conversion", http.StatusBadRequest)
		return
	}

	cmd := exec.Command("powershell", "-File", "./scripts/changeTriggers.ps1",
		"-TaskName", req.TaskName,
		"-TriggerType", triggerType,
		"-TriggerTime", triggerTime,
		"-Index", fmt.Sprintf("%d", req.Index),
		"-UserName", host.UserName,
		"-Password", host.Password,
		"-ComputerName", host.ComputerName,
	)
	output, err := cmd.CombinedOutput()
	var resp GenericResponse
	if err != nil {
		resp = GenericResponse{Success: false, Message: fmt.Sprintf("Error: %v, Output: %s", err, string(output))}
	} else {
		success, message := parsePowerShellOutput(output, "Trigger updated successfully")
		resp = GenericResponse{Success: success, Message: message}
	}
	sendAPIResponse(w, resp.Success, resp)
}

// DeleteTriggerHandler deletes a trigger from a task
func DeleteTriggerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req TriggerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	db := models.GetDB()
	userID := r.Context().Value("user_id").(int64)
	host, err := getTaskCredentials(db, userID, req.ComputerID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get remote host info: %v", err), http.StatusInternalServerError)
		return
	}
	cmd := exec.Command("powershell", "-File", "./scripts/deleteTriggers.ps1",
		"-TaskName", req.TaskName,
		"-Index", fmt.Sprintf("%d", req.Index),
		"-UserName", host.UserName,
		"-Password", host.Password,
		"-ComputerName", host.ComputerName,
	)
	output, err := cmd.CombinedOutput()
	var resp GenericResponse
	if err != nil {
		resp = GenericResponse{Success: false, Message: fmt.Sprintf("Error: %v, Output: %s", err, string(output))}
	} else {
		success, message := parsePowerShellOutput(output, "Trigger deleted successfully")
		resp = GenericResponse{Success: success, Message: message}
	}
	sendAPIResponse(w, resp.Success, resp)
}

// AddActionHandler adds a new action to a task
func AddActionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req ActionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	db := models.GetDB()
	userID := r.Context().Value("user_id").(int64)
	host, err := getTaskCredentials(db, userID, req.ComputerID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get remote host info: %v", err), http.StatusInternalServerError)
		return
	}
	cmd := exec.Command("powershell", "-File", "./scripts/addActions.ps1",
		"-TaskName", req.TaskName,
		"-Execute", req.Execute,
		"-Arguments", req.Args,
		"-WorkingDirectory", req.WorkingDirectory,
		"-UserName", host.UserName,
		"-Password", host.Password,
		"-ComputerName", host.ComputerName,
	)
	output, err := cmd.CombinedOutput()
	var resp GenericResponse
	if err != nil {
		resp = GenericResponse{Success: false, Message: fmt.Sprintf("Error: %v, Output: %s", err, string(output))}
	} else {
		success, message := parsePowerShellOutput(output, "Action added successfully")
		resp = GenericResponse{Success: success, Message: message}
	}
	sendAPIResponse(w, resp.Success, resp)
}

// UpdateActionHandler updates an existing action
func UpdateActionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req ActionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	db := models.GetDB()
	userID := r.Context().Value("user_id").(int64)
	host, err := getTaskCredentials(db, userID, req.ComputerID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get remote host info: %v", err), http.StatusInternalServerError)
		return
	}

	cmd := exec.Command("powershell", "-File", "./scripts/changeActions.ps1",
		"-TaskName", req.TaskName,
		"-Index", fmt.Sprintf("%d", req.Index),
		"-Execute", req.Execute,
		"-Arguments", req.Args,
		"-WorkingDirectory", req.WorkingDirectory,
		"-UserName", host.UserName,
		"-Password", host.Password,
		"-ComputerName", host.ComputerName,
	)
	output, err := cmd.CombinedOutput()
	var resp GenericResponse
	if err != nil {
		resp = GenericResponse{Success: false, Message: fmt.Sprintf("Error: %v, Output: %s", err, string(output))}
	} else {
		success, message := parsePowerShellOutput(output, "Action updated successfully")
		resp = GenericResponse{Success: success, Message: message}
	}
	sendAPIResponse(w, resp.Success, resp)
}

// DeleteActionHandler deletes an action from a task
func DeleteActionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req ActionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	db := models.GetDB()
	userID := r.Context().Value("user_id").(int64)
	host, err := getTaskCredentials(db, userID, req.ComputerID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get remote host info: %v", err), http.StatusInternalServerError)
		return
	}
	cmd := exec.Command("powershell", "-File", "./scripts/deleteActions.ps1",
		"-TaskName", req.TaskName,
		"-Index", fmt.Sprintf("%d", req.Index),
		"-UserName", host.UserName,
		"-Password", host.Password,
		"-ComputerName", host.ComputerName,
	)
	output, err := cmd.CombinedOutput()
	var resp GenericResponse
	if err != nil {
		resp = GenericResponse{Success: false, Message: fmt.Sprintf("Error: %v, Output: %s", err, string(output))}
	} else {
		success, message := parsePowerShellOutput(output, "Action deleted successfully")
		resp = GenericResponse{Success: success, Message: message}
	}
	sendAPIResponse(w, resp.Success, resp)
}
