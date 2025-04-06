package handlers

import (
	"database/sql"
	"encoding/xml"
	"fmt"
	"os/exec"
	"sync"

	"github.com/kevin93203/win-task-tracker/models"
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
