param (
    [Parameter(Mandatory=$true)]
    [string]$TaskName,
    [Parameter(Mandatory=$true)]
    [int]$Index,
    [Parameter(Mandatory=$true)]
    [string]$TriggerType,
    [Parameter(Mandatory=$true)]
    [string]$TriggerTime,
    [Parameter(Mandatory=$true)]
    [string]$UserName,
    [Parameter(Mandatory=$true)]
    [string]$Password,
    [Parameter(Mandatory=$true)]
    [string]$ComputerName
)

[Console]::OutputEncoding = [Text.Encoding]::UTF8

try {
    $securePassword = ConvertTo-SecureString $Password -AsPlainText -Force
    $credential = New-Object System.Management.Automation.PSCredential ($UserName, $securePassword)

    $result = Invoke-Command -ComputerName $ComputerName -Credential $credential -ScriptBlock {
        param($taskName, $index, $triggerType, $triggerTime)

        try {
            $triggerTypeParts = $triggerType -split ' '
            $baseTriggerType = $triggerTypeParts[0]

            # 處理時間格式
            $timeParts = $triggerTime -split ':'
            $hour = $timeParts[0].PadLeft(2, '0')
            $minute = $timeParts[1].PadLeft(2, '0')
            $timeOnly = "$hour`:$minute"

            # 預設 null
            $newTrigger = $null

            switch ($baseTriggerType) {
                "Daily" {
                    $newTrigger = New-ScheduledTaskTrigger -Daily -At $timeOnly
                }
                "Weekly" {
                    $daysOfWeek = @()
                    for ($i = 1; $i -lt $triggerTypeParts.Length; $i += 2) {
                        if ($triggerTypeParts[$i] -eq "-DaysOfWeek" -and ($i + 1) -lt $triggerTypeParts.Length) {
                            $daysOfWeek = $triggerTypeParts[$i + 1] -split ','
                        }
                    }
                    $newTrigger = New-ScheduledTaskTrigger -Weekly -DaysOfWeek $daysOfWeek -At $timeOnly
                }
                "Monthly" {
                    throw "Monthly trigger is not supported by PowerShell cmdlets"
                }
                default {
                    throw "Unsupported trigger type: $baseTriggerType"
                }
            }

            $task = Get-ScheduledTask -TaskName $taskName -ErrorAction Stop
            $triggers = $task.Triggers

            if ($index -lt 0 -or $index -ge $triggers.Count) {
                throw "Index out of range"
            }

            $triggers[$index] = $newTrigger

            Set-ScheduledTask -TaskName $taskName -Trigger $triggers -ErrorAction Stop

            return @{
                Success = $true
                Message = "Trigger at index $index updated successfully"
            }
        } catch {
            return @{
                Success = $false
                Error = $_.Exception.Message
            }
        }
    } -ArgumentList $TaskName, $Index, $TriggerType, $TriggerTime

    $result | ConvertTo-Json -Depth 5
    if ($result.Success -eq $false) { exit 1 } else { exit 0 }
} catch {
    @{
        Success = $false
        Error = $_.Exception.Message
    } | ConvertTo-Json -Depth 5
    exit 1
}
