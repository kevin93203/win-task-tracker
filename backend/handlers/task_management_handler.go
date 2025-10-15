package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"strings"
	"strconv"

	"github.com/kevin93203/win-task-tracker/models"
)

// TriggerRequest represents the payload for trigger operations
type TriggerRequest struct {
	ComputerID     int64  `json:"computer_id"`
	TaskName       string `json:"task_name"`
	CronExpression string `json:"cron_expression"` // Cron expression string
	Index          int    `json:"index,omitempty"` // Trigger index to update/delete
}

// ActionRequest represents the payload for action operations
type ActionRequest struct {
	ComputerID       int64  `json:"computer_id"`
	TaskName         string `json:"task_name"`
	Execute          string `json:"execute,omitempty"`
	Args             string `json:"args,omitempty"`
	WorkingDirectory string `json:"working_directory,omitempty"`
	Index            int    `json:"index,omitempty"` // Action index to update/delete
}

// GenericResponse standard API response
type GenericResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

// DeleteTriggerHandler deletes a trigger from a task
func DeleteTriggerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req TriggerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	db := models.GetDB()
	userID := r.Context().Value("user_id").(int64)
	host, err := getTaskCredentials(db, userID, req.ComputerID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get remote host info: %v", err), http.StatusInternalServerError)
		return
	}
	cmd := exec.Command("powershell", "-File", "./scripts/deleteTriggers.ps1",
		"-TaskName", req.TaskName,
		"-Index", fmt.Sprintf("%d", req.Index),
		"-UserName", host.UserName,
		"-Password", host.Password,
		"-ComputerName", host.ComputerName,
	)
	output, err := cmd.CombinedOutput()
	var resp GenericResponse
	if err != nil {
		resp = GenericResponse{Success: false, Message: fmt.Sprintf("Error: %v, Output: %s", err, string(output))}
	} else {
		success, message := parsePowerShellOutput(output, "Trigger deleted successfully")
		resp = GenericResponse{Success: success, Message: message}
	}
	sendAPIResponse(w, resp.Success, resp)
}

// AddActionHandler adds a new action to a task
func AddActionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req ActionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	db := models.GetDB()
	userID := r.Context().Value("user_id").(int64)
	host, err := getTaskCredentials(db, userID, req.ComputerID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get remote host info: %v", err), http.StatusInternalServerError)
		return
	}
	cmd := exec.Command("powershell", "-File", "./scripts/addActions.ps1",
		"-TaskName", req.TaskName,
		"-Execute", req.Execute,
		"-Arguments", req.Args,
		"-WorkingDirectory", req.WorkingDirectory,
		"-UserName", host.UserName,
		"-Password", host.Password,
		"-ComputerName", host.ComputerName,
	)
	output, err := cmd.CombinedOutput()
	var resp GenericResponse
	if err != nil {
		resp = GenericResponse{Success: false, Message: fmt.Sprintf("Error: %v, Output: %s", err, string(output))}
	} else {
		success, message := parsePowerShellOutput(output, "Action added successfully")
		resp = GenericResponse{Success: success, Message: message}
	}
	sendAPIResponse(w, resp.Success, resp)
}

// UpdateActionHandler updates an existing action
func UpdateActionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req ActionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	db := models.GetDB()
	userID := r.Context().Value("user_id").(int64)
	host, err := getTaskCredentials(db, userID, req.ComputerID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get remote host info: %v", err), http.StatusInternalServerError)
		return
	}

	cmd := exec.Command("powershell", "-File", "./scripts/changeActions.ps1",
		"-TaskName", req.TaskName,
		"-Index", fmt.Sprintf("%d", req.Index),
		"-Execute", req.Execute,
		"-Arguments", req.Args,
		"-WorkingDirectory", req.WorkingDirectory,
		"-UserName", host.UserName,
		"-Password", host.Password,
		"-ComputerName", host.ComputerName,
	)
	output, err := cmd.CombinedOutput()
	var resp GenericResponse
	if err != nil {
		resp = GenericResponse{Success: false, Message: fmt.Sprintf("Error: %v, Output: %s", err, string(output))}
	} else {
		success, message := parsePowerShellOutput(output, "Action updated successfully")
		resp = GenericResponse{Success: success, Message: message}
	}
	sendAPIResponse(w, resp.Success, resp)
}

// DeleteActionHandler deletes an action from a task
func DeleteActionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req ActionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	db := models.GetDB()
	userID := r.Context().Value("user_id").(int64)
	host, err := getTaskCredentials(db, userID, req.ComputerID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get remote host info: %v", err), http.StatusInternalServerError)
		return
	}
	cmd := exec.Command("powershell", "-File", "./scripts/deleteActions.ps1",
		"-TaskName", req.TaskName,
		"-Index", fmt.Sprintf("%d", req.Index),
		"-UserName", host.UserName,
		"-Password", host.Password,
		"-ComputerName", host.ComputerName,
	)
	output, err := cmd.CombinedOutput()
	var resp GenericResponse
	if err != nil {
		resp = GenericResponse{Success: false, Message: fmt.Sprintf("Error: %v, Output: %s", err, string(output))}
	} else {
		success, message := parsePowerShellOutput(output, "Action deleted successfully")
		resp = GenericResponse{Success: success, Message: message}
	}
	sendAPIResponse(w, resp.Success, resp)
}

// AddTriggerHandler adds a new trigger to a task
func AddTriggerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req WindowsTriggerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	db := models.GetDB()
	userID := r.Context().Value("user_id").(int64)
	host, err := getTaskCredentials(db, userID, req.ComputerID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get remote host info: %v", err), http.StatusInternalServerError)
		return
	}

	// Prepare parameters for PowerShell script
	var args []string
	args = append(args, "-File", "./scripts/addTriggers.ps1")
	args = append(args, "-TaskName", req.TaskName)
	args = append(args, "-UserName", host.UserName)
	args = append(args, "-Password", host.Password)
	args = append(args, "-ComputerName", host.ComputerName)
	args = append(args, "-StartBoundary", req.Trigger.StartBoundary)

	if req.Trigger.Repetition != nil {
		args = append(args, "-RepetitionInterval", req.Trigger.Repetition.Interval)
		if req.Trigger.Repetition.Duration != "" {
			args = append(args, "-RepetitionDuration", req.Trigger.Repetition.Duration)
		}
	}

	if req.Trigger.ScheduleByDay != nil {
		args = append(args, "-DaysInterval", fmt.Sprintf("%d", req.Trigger.ScheduleByDay.DaysInterval))
	} else if req.Trigger.ScheduleByWeek != nil {
		args = append(args, "-WeeksInterval", fmt.Sprintf("%d", req.Trigger.ScheduleByWeek.WeeksInterval))

		// Convert DaysOfWeek struct to array of strings
		daysOfWeek := []string{}
		if req.Trigger.ScheduleByWeek.DaysOfWeek.Sunday {
			daysOfWeek = append(daysOfWeek, "Sunday")
		}
		if req.Trigger.ScheduleByWeek.DaysOfWeek.Monday {
			daysOfWeek = append(daysOfWeek, "Monday")
		}
		if req.Trigger.ScheduleByWeek.DaysOfWeek.Tuesday {
			daysOfWeek = append(daysOfWeek, "Tuesday")
		}
		if req.Trigger.ScheduleByWeek.DaysOfWeek.Wednesday {
			daysOfWeek = append(daysOfWeek, "Wednesday")
		}
		if req.Trigger.ScheduleByWeek.DaysOfWeek.Thursday {
			daysOfWeek = append(daysOfWeek, "Thursday")
		}
		if req.Trigger.ScheduleByWeek.DaysOfWeek.Friday {
			daysOfWeek = append(daysOfWeek, "Friday")
		}
		if req.Trigger.ScheduleByWeek.DaysOfWeek.Saturday {
			daysOfWeek = append(daysOfWeek, "Saturday")
		}

		if len(daysOfWeek) > 0 {
			joined := "'" + strings.Join(daysOfWeek, "','") + "'"
			args = append(args, "-DaysOfWeek", fmt.Sprintf("@(%s)", joined))
		}
	} else if req.Trigger.ScheduleByMonth != nil {
		// Convert DaysOfMonth to array of days
		if len(req.Trigger.ScheduleByMonth.DaysOfMonth.Days) > 0 {
			days := make([]string, len(req.Trigger.ScheduleByMonth.DaysOfMonth.Days))
			for i, d := range req.Trigger.ScheduleByMonth.DaysOfMonth.Days {
				days[i] = strconv.Itoa(d)
			}
			args = append(args, "-DaysOfMonth", fmt.Sprintf("@(%s)", strings.Join(days, ",")))
		}

		// Convert Months struct to array of strings
		months := []string{}
		if req.Trigger.ScheduleByMonth.Months.January {
			months = append(months, "January")
		}
		if req.Trigger.ScheduleByMonth.Months.February {
			months = append(months, "February")
		}
		if req.Trigger.ScheduleByMonth.Months.March {
			months = append(months, "March")
		}
		if req.Trigger.ScheduleByMonth.Months.April {
			months = append(months, "April")
		}
		if req.Trigger.ScheduleByMonth.Months.May {
			months = append(months, "May")
		}
		if req.Trigger.ScheduleByMonth.Months.June {
			months = append(months, "June")
		}
		if req.Trigger.ScheduleByMonth.Months.July {
			months = append(months, "July")
		}
		if req.Trigger.ScheduleByMonth.Months.August {
			months = append(months, "August")
		}
		if req.Trigger.ScheduleByMonth.Months.September {
			months = append(months, "September")
		}
		if req.Trigger.ScheduleByMonth.Months.October {
			months = append(months, "October")
		}
		if req.Trigger.ScheduleByMonth.Months.November {
			months = append(months, "November")
		}
		if req.Trigger.ScheduleByMonth.Months.December {
			months = append(months, "December")
		}

		for _, month := range months {
			args = append(args, "-Months", month)
		}
	}

	// Use PowerShell script
	cmd := exec.Command("powershell", args...)
	fmt.Println(cmd.String())

	output, err := cmd.CombinedOutput()
	if err != nil {
		resp := GenericResponse{Success: false, Message: fmt.Sprintf("Error: %v, Output: %s", err, string(output))}
		sendAPIResponse(w, resp.Success, resp)
		return
	}

	success, message := parsePowerShellOutput(output, "Trigger added successfully")
	resp := GenericResponse{Success: success, Message: message}
	sendAPIResponse(w, resp.Success, resp)
}

// UpdateTriggerHandler updates a Windows Task Scheduler trigger in a task
func UpdateTriggerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req WindowsTriggerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	db := models.GetDB()
	userID := r.Context().Value("user_id").(int64)
	host, err := getTaskCredentials(db, userID, req.ComputerID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get remote host info: %v", err), http.StatusInternalServerError)
		return
	}

	// Prepare parameters for PowerShell script
	var args []string
	args = append(args, "-File", "./scripts/changeTriggers.ps1")
	args = append(args, "-TaskName", req.TaskName)
	args = append(args, "-Index", fmt.Sprintf("%d", req.Index))
	args = append(args, "-UserName", host.UserName)
	args = append(args, "-Password", host.Password)
	args = append(args, "-ComputerName", host.ComputerName)
	args = append(args, "-StartBoundary", req.Trigger.StartBoundary)

	if req.Trigger.Repetition != nil {
		args = append(args, "-RepetitionInterval", req.Trigger.Repetition.Interval)
		if req.Trigger.Repetition.Duration != "" {
			args = append(args, "-RepetitionDuration", req.Trigger.Repetition.Duration)
		}
	}

	if req.Trigger.ScheduleByDay != nil {
		args = append(args, "-DaysInterval", fmt.Sprintf("%d", req.Trigger.ScheduleByDay.DaysInterval))
	} else if req.Trigger.ScheduleByWeek != nil {
		args = append(args, "-WeeksInterval", fmt.Sprintf("%d", req.Trigger.ScheduleByWeek.WeeksInterval))

		// Convert DaysOfWeek struct to array of strings
		daysOfWeek := []string{}
		if req.Trigger.ScheduleByWeek.DaysOfWeek.Sunday {
			daysOfWeek = append(daysOfWeek, "Sunday")
		}
		if req.Trigger.ScheduleByWeek.DaysOfWeek.Monday {
			daysOfWeek = append(daysOfWeek, "Monday")
		}
		if req.Trigger.ScheduleByWeek.DaysOfWeek.Tuesday {
			daysOfWeek = append(daysOfWeek, "Tuesday")
		}
		if req.Trigger.ScheduleByWeek.DaysOfWeek.Wednesday {
			daysOfWeek = append(daysOfWeek, "Wednesday")
		}
		if req.Trigger.ScheduleByWeek.DaysOfWeek.Thursday {
			daysOfWeek = append(daysOfWeek, "Thursday")
		}
		if req.Trigger.ScheduleByWeek.DaysOfWeek.Friday {
			daysOfWeek = append(daysOfWeek, "Friday")
		}
		if req.Trigger.ScheduleByWeek.DaysOfWeek.Saturday {
			daysOfWeek = append(daysOfWeek, "Saturday")
		}

		if len(daysOfWeek) > 0 {
			joined := "'" + strings.Join(daysOfWeek, "','") + "'"
			args = append(args, "-DaysOfWeek", fmt.Sprintf("@(%s)", joined))
		}
	} else if req.Trigger.ScheduleByMonth != nil {
		// Convert DaysOfMonth to array of days
		if len(req.Trigger.ScheduleByMonth.DaysOfMonth.Days) > 0 {
			days := make([]string, len(req.Trigger.ScheduleByMonth.DaysOfMonth.Days))
			for i, d := range req.Trigger.ScheduleByMonth.DaysOfMonth.Days {
				days[i] = strconv.Itoa(d)
			}
			args = append(args, "-DaysOfMonth", fmt.Sprintf("@(%s)", strings.Join(days, ",")))
		}

		// Convert Months struct to array of strings
		months := []string{}
		if req.Trigger.ScheduleByMonth.Months.January {
			months = append(months, "January")
		}
		if req.Trigger.ScheduleByMonth.Months.February {
			months = append(months, "February")
		}
		if req.Trigger.ScheduleByMonth.Months.March {
			months = append(months, "March")
		}
		if req.Trigger.ScheduleByMonth.Months.April {
			months = append(months, "April")
		}
		if req.Trigger.ScheduleByMonth.Months.May {
			months = append(months, "May")
		}
		if req.Trigger.ScheduleByMonth.Months.June {
			months = append(months, "June")
		}
		if req.Trigger.ScheduleByMonth.Months.July {
			months = append(months, "July")
		}
		if req.Trigger.ScheduleByMonth.Months.August {
			months = append(months, "August")
		}
		if req.Trigger.ScheduleByMonth.Months.September {
			months = append(months, "September")
		}
		if req.Trigger.ScheduleByMonth.Months.October {
			months = append(months, "October")
		}
		if req.Trigger.ScheduleByMonth.Months.November {
			months = append(months, "November")
		}
		if req.Trigger.ScheduleByMonth.Months.December {
			months = append(months, "December")
		}

		for _, month := range months {
			args = append(args, "-Months", month)
		}
	}

	fmt.Println(args)
	// Use PowerShell script
	cmd := exec.Command("powershell", args...)
	fmt.Println(cmd.String())

	output, err := cmd.CombinedOutput()
	if err != nil {
		resp := GenericResponse{Success: false, Message: fmt.Sprintf("Error: %v, Output: %s", err, string(output))}
		sendAPIResponse(w, resp.Success, resp)
		return
	}

	success, message := parsePowerShellOutput(output, "Trigger updated successfully")
	resp := GenericResponse{Success: success, Message: message}
	sendAPIResponse(w, resp.Success, resp)
}
