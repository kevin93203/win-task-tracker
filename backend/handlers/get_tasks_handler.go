package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/kevin93203/win-task-tracker/auth"
	"github.com/kevin93203/win-task-tracker/models"
)

// GetTasksHandler 處理獲取所有排程任務的請求
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

	// 執行每個電腦的命令
	for _, host := range remoteHosts {
		command := fmt.Sprintf(
			"powershell -File ./scripts/getTasks.ps1 "+
				"-UserName '%s' "+
				"-Password '%s' "+
				"-ComputerID %d "+
				"-ComputerName '%s'",
			host.UserName,
			host.Password,
			host.ComputerID,
			host.ComputerName,
		)
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
