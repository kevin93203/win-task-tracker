package handlers

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"os/exec"
	"strconv"
	"strings"

	"win-task-tracker/backend/auth"
	"win-task-tracker/backend/models"
)

// GetTaskHandler 處理獲取單個排程任務的請求
func GetTaskHandler(w http.ResponseWriter, r *http.Request) {
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

	// 從 URL 獲取電腦 ID 和任務名稱
	path := strings.TrimPrefix(r.URL.Path, "/api/tasks/")
	parts := strings.SplitN(path, "/", 2)
	if len(parts) != 2 {
		http.Error(w, "無效的請求路徑", http.StatusBadRequest)
		return
	}

	computerID, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		http.Error(w, "無效的電腦 ID", http.StatusBadRequest)
		return
	}

	taskName := parts[1]

	// 獲取用戶的電腦憑證
	db := models.GetDB()

	// 檢查用戶是否有權限存取該電腦
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM computer_credential_mappings WHERE computer_id = ? AND created_by_id = ?", computerID, claims.UserID).Scan(&count)
	if err != nil {
		http.Error(w, fmt.Sprintf("檢查權限失敗: %v", err), http.StatusInternalServerError)
		return
	}

	if count == 0 {
		http.Error(w, "無權存取該電腦", http.StatusForbidden)
		return
	}

	// 獲取電腦的連接資訊
	var host RemoteHost
	err = db.QueryRow(`
		SELECT 
			rc.id, rc.name, 
			c.username, c.password
		FROM 
			remote_computers rc
		JOIN 
			computer_credential_mappings cc ON rc.id = cc.computer_id
		JOIN 
			credentials c ON cc.credential_id = c.id
		WHERE 
			rc.id = ? AND cc.created_by_id = ?
	`, computerID, claims.UserID).Scan(
		&host.ComputerID, &host.ComputerName,
		&host.UserName, &host.Password,
	)

	if err != nil {
		http.Error(w, fmt.Sprintf("獲取電腦憑證失敗: %v", err), http.StatusInternalServerError)
		return
	}

	// 執行獲取特定任務的指令
	command := fmt.Sprintf(
		"powershell -File ./scripts/getTask.ps1 "+
			"-UserName '%s' "+
			"-Password '%s' "+
			"-ComputerID %d "+
			"-ComputerName '%s' "+
			"-TaskName '%s'",
		host.UserName,
		host.Password,
		host.ComputerID,
		host.ComputerName,
		taskName,
	)

	output, err := executeCommandSync(command, host)
	if err != nil {
		response := TaskResponse{
			Errors: []TaskError{{
				ComputerName: host.ComputerName,
				UserName:     host.UserName,
				Error:        fmt.Sprintf("獲取任務失敗: %v", err),
			}},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// 處理成功的情況
	response := TaskResponse{
		Tasks: output.Tasks,
	}

	// 如果沒有找到任務，返回空任務和錯誤信息
	if len(output.Tasks) == 0 {
		response.Errors = []TaskError{{
			ComputerName: host.ComputerName,
			UserName:     host.UserName,
			Error:        fmt.Sprintf("找不到指定的任務: %s", taskName),
		}}
	}

	// 為前端準備單個任務的回應格式
	singleTaskResponse := map[string]interface{}{
		"task": nil,
	}
	if len(output.Tasks) > 0 {
		singleTaskResponse["task"] = output.Tasks[0]
	}

	// 返回JSON回應
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(singleTaskResponse); err != nil {
		http.Error(w, "無法轉換為 JSON", http.StatusInternalServerError)
		return
	}
}

// executeCommandSync 同步執行命令並處理輸出
func executeCommandSync(command string, host RemoteHost) (ScheduledTasks, error) {
	cmd := exec.Command("powershell", "-Command", command)
	output, err := cmd.Output()
	if err != nil {
		return ScheduledTasks{}, err
	}

	// 解析 XML 輸出
	var scheduledTasks ScheduledTasks
	err = xml.Unmarshal([]byte(output), &scheduledTasks)
	if err != nil {
		return ScheduledTasks{}, err
	}

	return scheduledTasks, nil
}
