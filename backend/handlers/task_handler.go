package handlers

import (
	"database/sql"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"sync"

	"github.com/kevin93203/win-task-tracker/auth"
	"github.com/kevin93203/win-task-tracker/models"
)

type ScheduledTasks struct {
	Tasks []Task `xml:"Task"`
}

type Task struct {
	RegistrationInfo RegistrationInfo `xml:"RegistrationInfo"`
	Triggers         []Triggers       `xml:"Triggers"`
	Actions          Actions          `xml:"Actions"`
	ExtraInfo        ExtraInfo        `xml:"ExtraInfo"`
}

type Triggers struct {
	TimeTriggers     []TimeTrigger     `xml:"TimeTrigger"`
	CalendarTriggers []CalendarTrigger `xml:"CalendarTrigger"`
}

type RegistrationInfo struct {
	Date        string `xml:"Date"`
	Author      string `xml:"Author"`
	Description string `xml:"Description"`
	URI         string `xml:"URI"`
}

type TimeTrigger struct {
	StartBoundary string        `xml:"StartBoundary"`
	Repetition    Repetition    `xml:"Repetition"`
	ScheduleByDay ScheduleByDay `xml:"ScheduleByDay"`
}

type CalendarTrigger struct {
	StartBoundary string        `xml:"StartBoundary"`
	Repetition    Repetition    `xml:"Repetition"`
	ScheduleByDay ScheduleByDay `xml:"ScheduleByDay"`
}

type Repetition struct {
	Interval string `xml:"Interval"`
	Duration string `xml:"Duration"`
}

type ScheduleByDay struct {
	DaysInterval string `xml:"DaysInterval"`
}

type Actions struct {
	Execs []Exec `xml:"Exec"`
}

type Exec struct {
	Command          string `xml:"Command"`
	Arguments        string `xml:"Arguments"`
	WorkingDirectory string `xml:"WorkingDirectory"`
}

type ExtraInfo struct {
	TaskName       string `xml:"TaskName"`
	ComputerName   string `xml:"ComputerName"`
	State          string `xml:"State"`
	LastRunTime    string `xml:"LastRunTime"`
	NextRunTime    string `xml:"NextRunTime"`
	LastTaskResult int    `xml:"LastTaskResult"`
}

type RemoteHost struct {
	UserName     string `json:"username"`
	Password     string `json:"password"`
	ComputerName string `json:"computer_name"`
}

type TaskError struct {
	ComputerName string `json:"computer_name"`
	UserName     string `json:"username"`
	Error        string `json:"error"`
}

type TaskResponse struct {
	Tasks  []Task      `json:"tasks"`
	Errors []TaskError `json:"errors,omitempty"`
}

type DisableTaskRequest struct {
	ComputerName string `json:"computer_name"`
	TaskName     string `json:"task_name"`
}

type DisableTaskResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

func executeCommand(command string, host RemoteHost, wg *sync.WaitGroup, results chan<- ScheduledTasks, errChan chan<- TaskError) {
	defer wg.Done()

	cmd := exec.Command("powershell", "-Command", command)
	output, err := cmd.Output()
	if err != nil {
		errChan <- TaskError{
			ComputerName: host.ComputerName,
			UserName:     host.UserName,
			Error:        fmt.Sprintf("執行命令錯誤: %v", err),
		}
		return
	}

	var scheduledTasks ScheduledTasks
	err = xml.Unmarshal([]byte(output), &scheduledTasks)
	if err != nil {
		errChan <- TaskError{
			ComputerName: host.ComputerName,
			UserName:     host.UserName,
			Error:        fmt.Sprintf("解析 XML 時發生錯誤: %v", err),
		}
		return
	}

	results <- scheduledTasks
}

func getHostsFromCredentials(db *sql.DB, userID int64) ([]RemoteHost, error) {
	mappings, err := models.GetComputerCredentialMappingsByUser(db, userID)
	if err != nil {
		return nil, fmt.Errorf("獲取電腦憑證映射失敗: %v", err)
	}

	var hosts []RemoteHost
	for _, mapping := range mappings {
		// 跳過沒有憑證的電腦
		if mapping.CredentialID == nil || mapping.CredentialUsername == nil {
			continue
		}

		// 獲取憑證密碼
		var password string
		err := db.QueryRow("SELECT password FROM credentials WHERE id = ?", *mapping.CredentialID).Scan(&password)
		if err != nil {
			return nil, fmt.Errorf("獲取憑證密碼失敗: %v", err)
		}

		hosts = append(hosts, RemoteHost{
			UserName:     *mapping.CredentialUsername,
			Password:     password,
			ComputerName: mapping.ComputerName,
		})
	}

	return hosts, nil
}

// DisableTaskHandler handles the request to disable a scheduled task on a remote computer
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

	// Get hosts from credentials
	hosts, err := getHostsFromCredentials(db, userID)
	if err != nil {
		http.Error(w, "Failed to get host credentials", http.StatusInternalServerError)
		return
	}

	// Find the target host
	var targetHost RemoteHost
	hostFound := false
	for _, host := range hosts {
		if host.ComputerName == req.ComputerName {
			targetHost = host
			hostFound = true
			break
		}
	}

	if !hostFound {
		http.Error(w, "Computer not found in user's credentials", http.StatusNotFound)
		return
	}

	// 使用 PowerShell 腳本停用任務
	cmd := fmt.Sprintf(
		"powershell -File ./disableTask.ps1 "+
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
		response.Success = true
		response.Message = fmt.Sprintf("Successfully disabled task '%s' on computer '%s'", req.TaskName, req.ComputerName)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func GetTasksHandler(w http.ResponseWriter, r *http.Request) {
	// 從 JWT 獲取用戶 ID
	cookie, err := r.Cookie("jwt")
	if err != nil {
		http.Error(w, "未授權", http.StatusUnauthorized)
		return
	}

	claims, err := auth.VerifyToken(cookie.Value)
	if err != nil {
		http.Error(w, "無效的令牌", http.StatusUnauthorized)
		return
	}

	// 獲取用戶的電腦憑證
	db := models.GetDB()
	remoteHosts, err := getHostsFromCredentials(db, int64(claims.UserID))

	if err != nil {
		http.Error(w, fmt.Sprintf("獲取電腦憑證失敗: %v", err), http.StatusInternalServerError)
		return
	}

	if len(remoteHosts) == 0 {
		// 如果用戶沒有設定任何電腦憑證，返回空回應
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(TaskResponse{Tasks: []Task{}})
		return
	}

	var wg sync.WaitGroup
	results := make(chan ScheduledTasks, len(remoteHosts))
	errChan := make(chan TaskError, len(remoteHosts))

	// 讀取 PowerShell 腳本
	script, err := os.ReadFile("taskScript.ps1")
	if err != nil {
		http.Error(w, "無法讀取 PowerShell 腳本", http.StatusInternalServerError)
		return
	}

	// 執行每個電腦的命令
	for _, host := range remoteHosts {
		command := fmt.Sprintf(string(script),
			host.Password, host.UserName, host.ComputerName, host.UserName)
		wg.Add(1)
		go executeCommand(command, host, &wg, results, errChan)
	}

	go func() {
		wg.Wait()
		close(results)
		close(errChan)
	}()

	// 收集所有任務和錯誤
	response := TaskResponse{}

	// 收集錯誤
	for err := range errChan {
		response.Errors = append(response.Errors, err)
	}

	// 收集任務
	for output := range results {
		response.Tasks = append(response.Tasks, output.Tasks...)
	}

	// 返回JSON回應
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "無法轉換為 JSON", http.StatusInternalServerError)
		return
	}
}
