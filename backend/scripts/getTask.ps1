param(
    [Parameter(Mandatory=$true)]
    [string]$UserName,
    [Parameter(Mandatory=$true)]
    [string]$Password,
    [Parameter(Mandatory=$true)]
    [int64]$ComputerID,
    [Parameter(Mandatory=$true)]
    [string]$ComputerName,
    [Parameter(Mandatory=$true)]
    [string]$TaskName
)

# 設定輸出編碼為 UTF8
[Console]::OutputEncoding = [Text.Encoding]::UTF8

try {
    # 建立認證物件
    $secpasswd = ConvertTo-SecureString $Password -AsPlainText -Force
    $cred = New-Object System.Management.Automation.PSCredential ($UserName, $secpasswd)
    
    # 建立遠端連線
    $session = New-PSSession -ComputerName $ComputerName -Credential $cred

    # 在遠端執行獲取特定排程任務的指令
    $task = Invoke-Command -Session $session -ScriptBlock {
        param($userName, $computerID, $computerName, $taskName)
        
        $hostname = $computerName
        $xmlDoc = New-Object System.Xml.XmlDocument
        $root = $xmlDoc.CreateElement("ScheduledTasks")
        $xmlDoc.AppendChild($root) | Out-Null

        # 根據任務名稱獲取特定任務
        try {
            $task = Get-ScheduledTask -TaskName $taskName 2>$null
            if (!$task) {
                throw "找不到指定的任務: $taskName"
            }
            
            $taskInfo = Get-ScheduledTaskInfo -TaskName $task.TaskName -TaskPath $task.TaskPath 2>$null
            $taskXml = Export-ScheduledTask -TaskName $task.TaskName -TaskPath $task.TaskPath

            $tempDoc = New-Object System.Xml.XmlDocument
            $tempDoc.LoadXml($taskXml)

            # 添加額外的節點
            $extraInfo = $tempDoc.CreateElement("ExtraInfo")
            $taskNameElement = $tempDoc.CreateElement("TaskName")
            $taskNameElement.InnerText = $task.TaskName
            
            $computerIDElement = $tempDoc.CreateElement("ComputerID")
            $computerIDElement.InnerText = $computerID
            
            $computerNameElement = $tempDoc.CreateElement("ComputerName")
            $computerNameElement.InnerText = $hostname
            $state = $tempDoc.CreateElement("State")
            $state.InnerText = $task.State
            $lastRunTime = $tempDoc.CreateElement("LastRunTime")
            $lastRunTime.InnerText = if ($taskInfo.LastRunTime) { $taskInfo.LastRunTime.ToString() } else { "N/A" }
            $nextRunTime = $tempDoc.CreateElement("NextRunTime")
            $nextRunTime.InnerText = if ($taskInfo.NextRunTime) { $taskInfo.NextRunTime.ToString() } else { "N/A" }
            $lastTaskResult = $tempDoc.CreateElement("LastTaskResult")
            $lastTaskResult.InnerText = if ($taskInfo.LastTaskResult) { $taskInfo.LastTaskResult.ToString() } else { 0 }

            $extraInfo.AppendChild($taskNameElement) | Out-Null
            $extraInfo.AppendChild($computerIDElement) | Out-Null
            $extraInfo.AppendChild($computerNameElement) | Out-Null
            $extraInfo.AppendChild($state) | Out-Null
            $extraInfo.AppendChild($lastRunTime) | Out-Null
            $extraInfo.AppendChild($nextRunTime) | Out-Null
            $extraInfo.AppendChild($lastTaskResult) | Out-Null
            $tempDoc.DocumentElement.AppendChild($extraInfo) | Out-Null

            $importNode = $xmlDoc.ImportNode($tempDoc.DocumentElement, $true)
            $root.AppendChild($importNode) | Out-Null
        } catch {
            Write-Error "無法獲取任務 '$taskName': $_"
            return $null
        }

        return $xmlDoc.OuterXml
    } -ArgumentList $UserName, $ComputerID, $ComputerName, $TaskName

    # 輸出結果
    Write-Output $task

    # 清理遠端連線
    Remove-PSSession $session
}
catch {
    Write-Error $_.Exception.Message
    exit 1
} 