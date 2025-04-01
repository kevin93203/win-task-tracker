package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"sync"
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

// 定義 Triggers 結構
type Triggers struct {
	TimeTriggers     []TimeTrigger     `xml:"TimeTrigger"`
	CalendarTriggers []CalendarTrigger `xml:"CalendarTrigger"`
}

type RegistrationInfo struct {
	Date   string `xml:"Date"`
	Author string `xml:"Author"`
	URI    string `xml:"URI"`
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
	ComputerName   string `xml:"ComputerName"`
	State          string `xml:"State"`
	LastRunTime    string `xml:"LastRunTime"`
	NextRunTime    string `xml:"NextRunTime"`
	LastTaskResult int    `xml:"LastTaskResult"`
}

type RemoteHost struct {
	UserName     string `json:"UserName"`
	Password     string `json:"Password"`
	ComputerName string `json:"ComputerName"`
}

func executeCommand(command string, wg *sync.WaitGroup, results chan<- ScheduledTasks) {
	defer wg.Done()

	cmd := exec.Command("powershell", "-Command", command)
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("執行命令錯誤:", err)
		return
	}

	var scheduledTasks ScheduledTasks
	err = xml.Unmarshal([]byte(output), &scheduledTasks)
	if err != nil {
		fmt.Println("解析 XML 時發生錯誤:", err)
		os.Exit(1)
	}

	results <- scheduledTasks
}

func getHostConfig(config_file string) ([]RemoteHost, error) {
	// 解析 JSON 為結構體切片
	var remoteHost []RemoteHost

	// 打開 JSON 檔案
	file, err := os.Open(config_file)
	if err != nil {
		return remoteHost, err
	}
	defer file.Close()

	// 讀取檔案內容
	bytes, err := io.ReadAll(file)
	if err != nil {
		return remoteHost, err
	}

	err = json.Unmarshal(bytes, &remoteHost)
	if err != nil {
		return remoteHost, err
	}

	return remoteHost, nil
}

func getTasks(w http.ResponseWriter, r *http.Request) {
	// 設置 CORS 標頭
	w.Header().Set("Access-Control-Allow-Origin", "*")                            // 允許所有來源
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")          // 允許的 HTTP 方法
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization") // 允許的請求標頭

	// 處理 OPTIONS 請求
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent) // 返回 204 No Content
		return
	}

	remoteHosts, err := getHostConfig("config.json")
	if err != nil {
		fmt.Println(err)
		return
	}

	var wg sync.WaitGroup
	results := make(chan ScheduledTasks, len(remoteHosts))

	// 讀取 PowerShell 腳本
	script, err := os.ReadFile("taskScript.ps1")
	if err != nil {
		http.Error(w, "無法讀取 PowerShell 腳本", http.StatusInternalServerError)
		return
	}

	var commands []string

	for _, host := range remoteHosts {
		// 將腳本中的佔位符替換為具體值
		command := fmt.Sprintf(string(script),
			host.Password, host.UserName, host.ComputerName, host.UserName)

		commands = append(commands, command)
	}

	for _, command := range commands {
		wg.Add(1)
		go executeCommand(command, &wg, results)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	var allTasks []Task // 用來存放所有的 Task

	// 設置 HTTP 標頭為 JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// 收集每個命令的輸出並附加到 allTasks 中
	for output := range results {
		allTasks = append(allTasks, output.Tasks...) // 將每個 output 的 Tasks 附加到 allTasks
	}

	// 將所有的 Task 轉換為 JSON
	jsonOutput, err := json.Marshal(allTasks)
	if err != nil {
		http.Error(w, "無法轉換為 JSON", http.StatusInternalServerError)
		return
	}
	w.Write(jsonOutput)
}

func main() {
	http.HandleFunc("/api/tasks", getTasks)
	fmt.Println("伺服器啟動於 http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
