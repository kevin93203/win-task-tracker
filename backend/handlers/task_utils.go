package handlers

import (
	"database/sql"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"os/exec"
	"sync"

	"win-task-tracker/backend/models"
)

// executeCommand 執行 PowerShell 命令並解析 XML 結果
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

// getHostsFromCredentials 從資料庫獲取用戶的遠端主機資訊
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
			ComputerID:   mapping.ComputerID,
			ComputerName: mapping.ComputerName,
		})
	}

	return hosts, nil
}

// getTaskCredentials 獲取執行任務所需的憑證和電腦資訊
func getTaskCredentials(db *sql.DB, userID int64, computerID int64) (*RemoteHost, error) {
	// 從資料庫中獲取電腦資訊
	computer, err := models.GetComputerByID(db, computerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get computer info: %v", err)
	}

	// 檢查用戶是否有權限訪問該電腦
	hasAccess, err := models.CheckUserComputerAccess(db, userID, computerID)
	if err != nil {
		return nil, fmt.Errorf("failed to check computer access: %v", err)
	}

	if !hasAccess {
		return nil, fmt.Errorf("unauthorized access to computer")
	}

	// 獲取電腦的認證資訊
	credential, err := models.GetComputerCredential(db, computerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get computer credentials: %v", err)
	}

	// 檢查是否有認證資訊
	if credential == nil || credential.Password == nil {
		return nil, fmt.Errorf("no valid credentials found for this computer")
	}

	// 準備遠端主機資訊
	targetHost := &RemoteHost{
		UserName:     credential.Username,
		Password:     *credential.Password,
		ComputerID:   computerID,
		ComputerName: computer.Name,
	}

	return targetHost, nil
}

// parsePowerShellOutput 解析 PowerShell 腳本輸出並檢查操作是否成功
func parsePowerShellOutput(output []byte, defaultSuccessMessage string) (bool, string) {
	// 先嘗試將輸出解析為 JSON 對象
	var result map[string]interface{}
	if err := json.Unmarshal(output, &result); err != nil {
		// 如果解析失敗，嘗試解析為 JSON 數組
		var resultArray []interface{}
		if err := json.Unmarshal(output, &resultArray); err != nil {
			// 如果兩種方式都解析失敗，返回錯誤
			return false, fmt.Sprintf("Failed to parse script output: %s\nOutput: %s", err.Error(), string(output))
		}

		// 如果成功解析為數組但數組為空，返回默認成功消息
		if len(resultArray) == 0 {
			return true, defaultSuccessMessage
		}

		// 如果數組不為空，檢查第一個元素是否為對象
		if obj, ok := resultArray[0].(map[string]interface{}); ok {
			// 檢查成功狀態
			if success, ok := obj["Success"].(bool); ok && !success {
				if errMsg, ok := obj["Error"].(string); ok {
					return false, fmt.Sprintf("PowerShell error: %s", errMsg)
				}
				return false, "Unknown PowerShell error"
			}

			// 檢查消息
			if message, ok := obj["Message"].(string); ok && message != "" {
				return true, message
			}
		}

		// 默認認為是成功的
		return true, defaultSuccessMessage
	}

	// 檢查 PowerShell 腳本是否回報錯誤
	if success, ok := result["Success"].(bool); ok && !success {
		if errMsg, ok := result["Error"].(string); ok {
			return false, fmt.Sprintf("PowerShell error: %s", errMsg)
		}
		return false, "Unknown PowerShell error"
	}

	// 如果有自定義成功消息，使用它
	if message, ok := result["Message"].(string); ok && message != "" {
		return true, message
	}

	// 否則使用默認成功消息
	return true, defaultSuccessMessage
}

// sendAPIResponse 根據操作結果發送適當的 HTTP 響應
func sendAPIResponse(w http.ResponseWriter, success bool, data interface{}) {
	w.Header().Set("Content-Type", "application/json")

	// 如果操作失敗，設置 HTTP 狀態碼為 400 Bad Request
	if !success {
		w.WriteHeader(http.StatusBadRequest)
	}

	// 編碼並發送 JSON 響應
	json.NewEncoder(w).Encode(data)
}
