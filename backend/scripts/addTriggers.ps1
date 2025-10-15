param (
    [string]$TaskName,
    [string]$StartBoundary,
    [string]$RepetitionInterval,
    [string]$RepetitionDuration,
    [int]$DaysInterval,
    [int]$WeeksInterval,
    [string[]]$DaysOfWeek,
    [int[]]$DaysOfMonth,
    [string[]]$Months,
    [string]$UserName,
    [string]$Password,
    [string]$ComputerName
)

Write-Output $DaysOfWeek[0]
    

# # 創建安全密碼
# $securePassword = ConvertTo-SecureString $Password -AsPlainText -Force
# $credential = New-Object System.Management.Automation.PSCredential ($UserName, $securePassword)

# # 定義返回結果
# $resultObj = @{ Success = $false; Message = ""; Error = "" }

# # 嘗試遠程連線
# try {
#     $session = New-PSSession -ComputerName $ComputerName -Credential $credential -ErrorAction Stop
    
#     # 在遠程機器上執行命令
#     $result = Invoke-Command -Session $session -ScriptBlock {
#         param($taskName, $startBoundary, $repetitionInterval, $repetitionDuration, $daysInterval, $weeksInterval, $daysOfWeek, $daysOfMonth, $months)
        
#         # 定義返回結果
#         $resultObj = @{ Success = $false; Message = ""; Error = "" }
        
#         # 嘗試獲取任務
#         try {
#             $task = Get-ScheduledTask -TaskName $taskName -ErrorAction Stop
#             $existingTriggers = $task.Triggers
            
#             # 解析開始時間
#             $startTime = [DateTime]::Parse($startBoundary)
#             $timeOnly = $startTime.ToString("HH:mm")

#             # 創建新觸發器
#             $newTrigger = $null
            
#             # 根據參數創建適當類型的觸發器
#             if ($daysInterval -gt 0) {
#                 # 每日觸發器
#                 $newTrigger = New-ScheduledTaskTrigger -Daily -DaysInterval $daysInterval -At $timeOnly
#             } elseif ($weeksInterval -gt 0) {
#                 # 每週觸發器
#                 $newTrigger = New-ScheduledTaskTrigger -Weekly -WeeksInterval $weeksInterval -At $timeOnly
#             } else {
#                 # 一次性觸發器
#                 $newTrigger = New-ScheduledTaskTrigger -Once -At $startTime
#             }
            
#             if ($existingTriggers.Count -eq 0) {
#                 $updatedTriggers = @($newTrigger)
#             } else {
#                 $updatedTriggers = $existingTriggers + $newTrigger
#             }
            
#             # 更新任務
#             Set-ScheduledTask -TaskName $taskName -Trigger $updatedTriggers -ErrorAction Stop
            
#             $resultObj.Success = $true
#             $resultObj.Message = "Trigger added successfully"
#         } catch {
#             $resultObj.Success = $false
#             $resultObj.Error = $_.Exception.Message
#         }
        
#         return $resultObj
#     } -ArgumentList $TaskName, $StartBoundary, $RepetitionInterval, $RepetitionDuration, $DaysInterval, $WeeksInterval, $DaysOfWeek, $DaysOfMonth, $Months
    
#     # 關閉遠程會話
#     Remove-PSSession $session
    
#     $resultObj = $result
# } catch {
#     $resultObj.Success = $false
#     $resultObj.Error = $_.Exception.Message
# }

# # 輸出結果
# $resultObj | ConvertTo-Json -Depth 5

# # 設置退出代碼
# if ($resultObj.Success -eq $false) { exit 1 } else { exit 0 }
