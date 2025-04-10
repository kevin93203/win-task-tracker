param(
    [Parameter(Mandatory=$true)][string]$Mode,
    [Parameter(Mandatory=$true)][string]$TaskName,
    [string]$TriggerType,
    [string]$TriggerTime,
    [string]$Command,
    [string]$Args,
    [int]$Index = 0,
    [Parameter(Mandatory=$true)][string]$UserName,
    [Parameter(Mandatory=$true)][string]$Password,
    [Parameter(Mandatory=$true)][string]$ComputerName
)

# Set output encoding
[Console]::OutputEncoding = [Text.Encoding]::UTF8

try {
    $securePassword = ConvertTo-SecureString $Password -AsPlainText -Force
    $credential = New-Object System.Management.Automation.PSCredential ($UserName, $securePassword)

    $result = Invoke-Command -ComputerName $ComputerName -Credential $credential -ScriptBlock {
        param($mode, $taskName, $triggerType, $triggerTime, $command, $args, $index)

        try {
            switch ($mode) {
                "AddTrigger" {
                    # 解析觸發器類型和參數
                    $triggerTypeParts = $triggerType -split ' '
                    $baseTriggerType = $triggerTypeParts[0]
                    
                    # 處理時間格式 - 確保是 HH:MM 格式（兩位數小時和分鐘）
                    $timeParts = $triggerTime -split ':'
                    $hour = $timeParts[0].PadLeft(2, '0')
                    $minute = $timeParts[1].PadLeft(2, '0')
                    $timeOnly = "$hour`:$minute"
                    
                    # 轉換 PowerShell 的觸發器類型到 SCHTASKS 格式
                    $schTaskType = ""
                    $additionalParams = ""
                    
                    switch ($baseTriggerType) {
                        "Daily" {
                            $schTaskType = "DAILY"
                        }
                        "Weekly" {
                            $schTaskType = "WEEKLY"
                            # 尋找星期幾參數
                            for ($i = 1; $i -lt $triggerTypeParts.Length; $i += 2) {
                                if ($triggerTypeParts[$i] -eq "-DaysOfWeek" -and ($i + 1) -lt $triggerTypeParts.Length) {
                                    $day = $triggerTypeParts[$i + 1]
                                    $additionalParams = "/D $day"
                                    break
                                }
                            }
                        }
                        "Monthly" {
                            $schTaskType = "MONTHLY"
                            # 尋找日期參數
                            $daysOfMonth = ""
                            $monthsOfYear = ""
                            
                            for ($i = 1; $i -lt $triggerTypeParts.Length; $i += 2) {
                                if ($triggerTypeParts[$i] -eq "-DaysOfMonth" -and ($i + 1) -lt $triggerTypeParts.Length) {
                                    $daysOfMonth = $triggerTypeParts[$i + 1]
                                }
                                if ($triggerTypeParts[$i] -eq "-MonthsOfYear" -and ($i + 1) -lt $triggerTypeParts.Length) {
                                    $monthsOfYear = $triggerTypeParts[$i + 1]
                                }
                            }
                            
                            if ($daysOfMonth -ne "") {
                                $additionalParams += "/D $daysOfMonth"
                            }
                            if ($monthsOfYear -ne "") {
                                $additionalParams += " /M $monthsOfYear"
                            }
                        }
                    }
                    
                    # 獲取現有任務的詳細信息
                    $taskInfo = SCHTASKS /Query /TN `"$taskName`" /FO LIST /V
                    if ($LASTEXITCODE -ne 0) {
                        throw "Failed to get task information"
                    }

                    # 調試輸出
                    Write-Host "Task Info Output:"
                    $taskInfo | ForEach-Object { Write-Host $_ }

                    # 解析任務信息
                    $taskDetails = @{}
                    $taskAction = $null
                    $runAsUser = $null
                    
                    # 直接使用 Get-ScheduledTask 獲取任務信息
                    try {
                        $task = Get-ScheduledTask -TaskName $taskName -ErrorAction Stop
                        $taskAction = $task.Actions[0].Execute
                        if ($task.Actions[0].Arguments) {
                            $taskAction += " " + $task.Actions[0].Arguments
                        }
                        $runAsUser = $task.Principal.UserId
                        
                        $taskDetails['TaskToRun'] = $taskAction
                        $taskDetails['RunAsUser'] = $runAsUser
                    }
                    catch {
                        # 如果 Get-ScheduledTask 失敗，嘗試解析 SCHTASKS 輸出
                        $taskInfo | ForEach-Object {
                            if ($_ -match "Task To Run:(.*)") {
                                $taskDetails['TaskToRun'] = $matches[1].Trim()
                            }
                            elseif ($_ -match "Run As User:(.*)") {
                                $taskDetails['RunAsUser'] = $matches[1].Trim()
                            }
                            elseif ($_ -match "Start In:(.*)") {
                                $taskDetails['StartIn'] = $matches[1].Trim()
                            }
                        }
                    }
                    
                    # 檢查必要的任務信息是否存在
                    if (-not $taskDetails['TaskToRun']) {
                        throw "Failed to get task command information. Please check if the task exists and has a valid action."
                    }
                    
                    Write-Host "Task to Run: $($taskDetails['TaskToRun'])"
                    Write-Host "Run As User: $($taskDetails['RunAsUser'])"
                    
                    # 準備 SCHTASKS 命令
                    $schTasksCommand = "SCHTASKS /Create /F /TN `"$taskName`" /TR `"$($taskDetails['TaskToRun'])`" /ST $timeOnly /SC $schTaskType"
                    if ($additionalParams -ne "") {
                        $schTasksCommand += " $additionalParams"
                    }
                    
                    # 添加其他必要的參數
                    if ($taskDetails['RunAsUser']) {
                        $schTasksCommand += " /RU `"$($taskDetails['RunAsUser'])`""
                    }
                    if ($taskDetails['StartIn']) {
                        $schTasksCommand += " /SD `"$($taskDetails['StartIn'])`""
                    }
                    
                    # 執行命令
                    Write-Host "Executing: $schTasksCommand"
                    $output = Invoke-Expression $schTasksCommand
                    
                    # 檢查命令是否成功執行
                    if ($LASTEXITCODE -ne 0) {
                        throw "SCHTASKS command failed with exit code $LASTEXITCODE. Output: $output"
                    }
                    
                    return @{ 
                        Success = $true
                        Message = "Trigger added successfully"
                        Details = $output
                        Command = $schTasksCommand
                    }
                }
                "UpdateTrigger" {
                    # 使用與 AddTrigger 相同的邏輯
                    $triggerTypeParts = $triggerType -split ' '
                    $baseTriggerType = $triggerTypeParts[0]
                    
                    # 處理時間格式 - 確保是 HH:MM 格式（兩位數小時和分鐘）
                    $timeParts = $triggerTime -split ':'
                    $hour = $timeParts[0].PadLeft(2, '0')
                    $minute = $timeParts[1].PadLeft(2, '0')
                    $timeOnly = "$hour`:$minute"
                    
                    # 轉換 PowerShell 的觸發器類型到 SCHTASKS 格式
                    $schTaskType = ""
                    $additionalParams = ""
                    
                    switch ($baseTriggerType) {
                        "Daily" {
                            $schTaskType = "DAILY"
                        }
                        "Weekly" {
                            $schTaskType = "WEEKLY"
                            # 尋找星期幾參數
                            for ($i = 1; $i -lt $triggerTypeParts.Length; $i += 2) {
                                if ($triggerTypeParts[$i] -eq "-DaysOfWeek" -and ($i + 1) -lt $triggerTypeParts.Length) {
                                    $day = $triggerTypeParts[$i + 1]
                                    $additionalParams = "/D $day"
                                    break
                                }
                            }
                        }
                        "Monthly" {
                            $schTaskType = "MONTHLY"
                            # 尋找日期參數
                            $daysOfMonth = ""
                            $monthsOfYear = ""
                            
                            for ($i = 1; $i -lt $triggerTypeParts.Length; $i += 2) {
                                if ($triggerTypeParts[$i] -eq "-DaysOfMonth" -and ($i + 1) -lt $triggerTypeParts.Length) {
                                    $daysOfMonth = $triggerTypeParts[$i + 1]
                                }
                                if ($triggerTypeParts[$i] -eq "-MonthsOfYear" -and ($i + 1) -lt $triggerTypeParts.Length) {
                                    $monthsOfYear = $triggerTypeParts[$i + 1]
                                }
                            }
                            
                            if ($daysOfMonth -ne "") {
                                $additionalParams += "/D $daysOfMonth"
                            }
                            if ($monthsOfYear -ne "") {
                                $additionalParams += " /M $monthsOfYear"
                            }
                        }
                    }
                    
                    # 獲取現有任務的詳細信息
                    $taskInfo = SCHTASKS /Query /TN `"$taskName`" /FO LIST /V
                    if ($LASTEXITCODE -ne 0) {
                        throw "Failed to get task information"
                    }
                    
                    # 調試輸出
                    Write-Host "Task Info Output:"
                    $taskInfo | ForEach-Object { Write-Host $_ }

                    # 解析任務信息
                    $taskDetails = @{}
                    $taskAction = $null
                    $runAsUser = $null
                    
                    # 直接使用 Get-ScheduledTask 獲取任務信息
                    try {
                        $task = Get-ScheduledTask -TaskName $taskName -ErrorAction Stop
                        $taskAction = $task.Actions[0].Execute
                        if ($task.Actions[0].Arguments) {
                            $taskAction += " " + $task.Actions[0].Arguments
                        }
                        $runAsUser = $task.Principal.UserId
                        
                        $taskDetails['TaskToRun'] = $taskAction
                        $taskDetails['RunAsUser'] = $runAsUser
                    }
                    catch {
                        # 如果 Get-ScheduledTask 失敗，嘗試解析 SCHTASKS 輸出
                        $taskInfo | ForEach-Object {
                            if ($_ -match "Task To Run:(.*)") {
                                $taskDetails['TaskToRun'] = $matches[1].Trim()
                            }
                            elseif ($_ -match "Run As User:(.*)") {
                                $taskDetails['RunAsUser'] = $matches[1].Trim()
                            }
                            elseif ($_ -match "Start In:(.*)") {
                                $taskDetails['StartIn'] = $matches[1].Trim()
                            }
                        }
                    }
                    
                    # 檢查必要的任務信息是否存在
                    if (-not $taskDetails['TaskToRun']) {
                        throw "Failed to get task command information. Please check if the task exists and has a valid action."
                    }
                    
                    Write-Host "Task to Run: $($taskDetails['TaskToRun'])"
                    Write-Host "Run As User: $($taskDetails['RunAsUser'])"
                    
                    # 準備 SCHTASKS 命令
                    $schTasksCommand = "SCHTASKS /Create /F /TN `"$taskName`" /TR `"$($taskDetails['TaskToRun'])`" /ST $timeOnly /SC $schTaskType"
                    if ($additionalParams -ne "") {
                        $schTasksCommand += " $additionalParams"
                    }
                    
                    # 添加其他必要的參數
                    if ($taskDetails['RunAsUser']) {
                        $schTasksCommand += " /RU `"$($taskDetails['RunAsUser'])`""
                    }
                    if ($taskDetails['StartIn']) {
                        $schTasksCommand += " /SD `"$($taskDetails['StartIn'])`""
                    }
                    
                    # 執行命令
                    Write-Host "Executing: $schTasksCommand"
                    $output = Invoke-Expression $schTasksCommand
                    
                    # 檢查命令是否成功執行
                    if ($LASTEXITCODE -ne 0) {
                        throw "SCHTASKS command failed with exit code $LASTEXITCODE. Output: $output"
                    }
                    
                    return @{ 
                        Success = $true
                        Message = "Trigger updated successfully"
                        Details = $output
                        Command = $schTasksCommand
                    }
                }
                "DeleteTrigger" {
                    # 刪除觸發器需要重新創建任務，但先刪除指定的觸發器
                    return @{ Success = $true; Message = "Trigger deleted (not implemented)" }
                }
                "AddAction" {
                    # 使用 SCHTASKS 添加動作
                    $command = "SCHTASKS /Change /TN `"$taskName`" /TR `"$command $args`""
                    $output = Invoke-Expression $command
                    
                    return @{ 
                        Success = $true
                        Message = "Action added successfully"
                        Details = $output
                    }
                }
                "UpdateAction" {
                    # 更新動作
                    $command = "SCHTASKS /Change /TN `"$taskName`" /TR `"$command $args`""
                    $output = Invoke-Expression $command
                    
                    return @{ 
                        Success = $true
                        Message = "Action updated successfully"
                        Details = $output
                    }
                }
                "DeleteAction" {
                    # 刪除動作 (實際上是設置為空)
                    $command = "SCHTASKS /Change /TN `"$taskName`" /TR `"`""
                    $output = Invoke-Expression $command
                    
                    return @{ 
                        Success = $true
                        Message = "Action deleted successfully"
                        Details = $output
                    }
                }
                default {
                    return @{ Success = $false; Error = "Invalid mode: $mode" }
                }
            }
        } catch {
            return @{ 
                Success = $false
                Error = $_.Exception.Message
                StackTrace = $_.ScriptStackTrace
            }
        }
    } -ArgumentList $Mode, $TaskName, $TriggerType, $TriggerTime, $Command, $Args, $Index

    $result | ConvertTo-Json -Depth 10
} catch {
    @{ 
        Success = $false
        Error = $_.Exception.Message
        StackTrace = $_.ScriptStackTrace
    } | ConvertTo-Json -Depth 10
}