package handlers

import "encoding/xml"

// XML 解析相關的結構體
type ScheduledTasks struct {
	Tasks []Task `xml:"Task"`
}

type Task struct {
	RegistrationInfo RegistrationInfo `xml:"RegistrationInfo"`
	Triggers         Triggers         `xml:"Triggers"`
	Actions          Actions          `xml:"Actions"`
	ExtraInfo        ExtraInfo        `xml:"ExtraInfo"`
}

type Triggers struct {
	TimeTriggers     *[]TimeTrigger     `xml:"TimeTrigger,omitempty" json:"TimeTriggers,omitempty"`
	CalendarTriggers *[]CalendarTrigger `xml:"CalendarTrigger,omitempty" json:"CalendarTriggers,omitempty"`
}

type RegistrationInfo struct {
	Date        string `xml:"Date"`
	Author      string `xml:"Author"`
	Description string `xml:"Description,omitempty" json:"Description,omitempty"`
	URI         string `xml:"URI"`
}

type TimeTrigger struct {
	StartBoundary string         `xml:"StartBoundary"`
	Repetition    *Repetition    `xml:"Repetition,omitempty" json:"Repetition,omitempty"`
	ScheduleByDay *ScheduleByDay `xml:"ScheduleByDay,omitempty" json:"ScheduleByDay,omitempty"`
}

type CalendarTrigger struct {
	StartBoundary   string           `xml:"StartBoundary"`
	Repetition      *Repetition      `xml:"Repetition,omitempty" json:"Repetition,omitempty"`
	ScheduleByDay   *ScheduleByDay   `xml:"ScheduleByDay,omitempty" json:"ScheduleByDay,omitempty"`
	ScheduleByWeek  *ScheduleByWeek  `xml:"ScheduleByWeek,omitempty" json:"ScheduleByWeek,omitempty"`
	ScheduleByMonth *ScheduleByMonth `xml:"ScheduleByMonth,omitempty" json:"ScheduleByMonth,omitempty"`
}

type Repetition struct {
	Interval string `xml:"Interval,omitempty" json:"Interval,omitempty"`
	Duration string `xml:"Duration,omitempty" json:"Duration,omitempty"`
}

type ScheduleByDay struct {
	DaysInterval int `xml:"DaysInterval,omitempty"`
}

type ScheduleByWeek struct {
	WeeksInterval int        `xml:"WeeksInterval,omitempty"`
	DaysOfWeek    DaysOfWeek `xml:"DaysOfWeek,omitempty"`
}

type ScheduleByMonth struct {
	Months      Months      `xml:"Months,omitempty"`
	DaysOfMonth DaysOfMonth `xml:"DaysOfMonth,omitempty"`
}

type BoolElement bool

func (b *BoolElement) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	*b = true
	return d.Skip()
}

type DaysOfWeek struct {
	Sunday    BoolElement `xml:"Sunday"`
	Monday    BoolElement `xml:"Monday"`
	Tuesday   BoolElement `xml:"Tuesday"`
	Wednesday BoolElement `xml:"Wednesday"`
	Thursday  BoolElement `xml:"Thursday"`
	Friday    BoolElement `xml:"Friday"`
	Saturday  BoolElement `xml:"Saturday"`
}

// 定義 DaysOfMonth 結構體
type DaysOfMonth struct {
	Days []int `xml:"Day"` // 使用 xml 標籤來映射 Day 元素
}

type Months struct {
	January   BoolElement `xml:"January"`
	February  BoolElement `xml:"February"`
	March     BoolElement `xml:"March"`
	April     BoolElement `xml:"April"`
	May       BoolElement `xml:"May"`
	June      BoolElement `xml:"June"`
	July      BoolElement `xml:"July"`
	August    BoolElement `xml:"August"`
	September BoolElement `xml:"September"`
	October   BoolElement `xml:"October"`
	November  BoolElement `xml:"November"`
	December  BoolElement `xml:"December"`
}

type Actions struct {
	Execs []Exec `xml:"Exec"`
}

type Exec struct {
	Command          string `xml:"Command"`
	Arguments        string `xml:"Arguments,omitempty" json:"Arguments,omitempty"`
	WorkingDirectory string `xml:"WorkingDirectory,omitempty" json:"WorkingDirectory,omitempty"`
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
