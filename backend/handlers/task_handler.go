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

func executeCommand(command string, wg *sync.WaitGroup, results chan<- ScheduledTasks, errChan chan<- error) {
	defer wg.Done()

	cmd := exec.Command("powershell", "-Command", command)
	output, err := cmd.Output()
	if err != nil {
		errChan <- fmt.Errorf("執行命令錯誤: %v", err)
		return
	}

	var scheduledTasks ScheduledTasks
	err = xml.Unmarshal([]byte(output), &scheduledTasks)
	if err != nil {
		errChan <- fmt.Errorf("解析 XML 時發生錯誤: %v", err)
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
		// 如果用戶沒有設定任何電腦憑證，返回空陣列
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("[]"))
		return
	}

	var wg sync.WaitGroup
	results := make(chan ScheduledTasks, len(remoteHosts))
	errChan := make(chan error, len(remoteHosts))

	// 讀取 PowerShell 腳本
	script, err := os.ReadFile("taskScript.ps1")
	if err != nil {
		http.Error(w, "無法讀取 PowerShell 腳本", http.StatusInternalServerError)
		return
	}

	var commands []string
	for _, host := range remoteHosts {
		command := fmt.Sprintf(string(script),
			host.Password, host.UserName, host.ComputerName, host.UserName)
		commands = append(commands, command)
	}

	for _, command := range commands {
		wg.Add(1)
		go executeCommand(command, &wg, results, errChan)
	}

	go func() {
		wg.Wait()
		close(results)
		close(errChan)
	}()

	// 檢查是否有錯誤
	var errors []error
	for err := range errChan {
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		errorMsg := "執行任務時發生錯誤:\n"
		for _, err := range errors {
			errorMsg += err.Error() + "\n"
		}
		http.Error(w, errorMsg, http.StatusInternalServerError)
		return
	}

	var allTasks []Task
	for output := range results {
		allTasks = append(allTasks, output.Tasks...)
	}

	w.Header().Set("Content-Type", "application/json")
	jsonOutput, err := json.Marshal(allTasks)
	if err != nil {
		http.Error(w, "無法轉換為 JSON", http.StatusInternalServerError)
		return
	}
	w.Write(jsonOutput)
}
