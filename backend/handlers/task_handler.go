package handlers

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
	var remoteHost []RemoteHost

	file, err := os.Open(config_file)
	if err != nil {
		return remoteHost, err
	}
	defer file.Close()

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

func GetTasksHandler(w http.ResponseWriter, r *http.Request) {
	// 設置 CORS 標頭
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	// 處理 OPTIONS 請求
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
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

	var allTasks []Task

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	for output := range results {
		allTasks = append(allTasks, output.Tasks...)
	}

	jsonOutput, err := json.Marshal(allTasks)
	if err != nil {
		http.Error(w, "無法轉換為 JSON", http.StatusInternalServerError)
		return
	}
	w.Write(jsonOutput)
}
