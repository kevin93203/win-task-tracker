package handlers

// XML 解析相關的結構體
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
	TaskName       string `xml:"TaskName" json:"TaskName"`
	ComputerName   string `xml:"ComputerName" json:"ComputerName"`
	ComputerID     int64  `xml:"ComputerID" json:"ComputerID"`
	State          string `xml:"State" json:"State"`
	LastRunTime    string `xml:"LastRunTime" json:"LastRunTime"`
	NextRunTime    string `xml:"NextRunTime" json:"NextRunTime"`
	LastTaskResult int    `xml:"LastTaskResult" json:"LastTaskResult"`
}

// 遠端主機相關結構體
type RemoteHost struct {
	UserName     string `json:"username"`
	Password     string `json:"password"`
	ComputerID   int64  `json:"computer_id"`
	ComputerName string `json:"computer_name"`
}

// 錯誤和回應相關結構體
type TaskError struct {
	ComputerName string `json:"computer_name"`
	UserName     string `json:"username"`
	Error        string `json:"error"`
}

type TaskResponse struct {
	Tasks  []Task      `json:"tasks"`
	Errors []TaskError `json:"errors,omitempty"`
}

// 啟用/停用任務相關結構體
type DisableTaskRequest struct {
	ComputerID int64  `json:"computer_id"`
	TaskName   string `json:"task_name"`
}

type DisableTaskResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

type EnableTaskRequest struct {
	ComputerID int64  `json:"computer_id"`
	TaskName   string `json:"task_name"`
}

type EnableTaskResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

// 啟動/停止任務相關結構體
type StartTaskRequest struct {
	ComputerID int64  `json:"computer_id"`
	TaskName   string `json:"task_name"`
}

type StartTaskResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	State   string `json:"state,omitempty"`
	Error   string `json:"error,omitempty"`
}

type StopTaskRequest struct {
	ComputerID int64  `json:"computer_id"`
	TaskName   string `json:"task_name"`
}

type StopTaskResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	State   string `json:"state,omitempty"`
	Error   string `json:"error,omitempty"`
}
