param(
    [Parameter(Mandatory=$true)]
    [string]$UserName,
    [Parameter(Mandatory=$true)]
    [string]$Password,
    [Parameter(Mandatory=$true)]
    [int64]$ComputerID,
    [Parameter(Mandatory=$true)]
    [string]$ComputerName
)

# 設定輸出編碼為 UTF8
[Console]::OutputEncoding = [Text.Encoding]::UTF8

try {
    # 建立認證物件
    $secpasswd = ConvertTo-SecureString $Password -AsPlainText -Force
    $cred = New-Object System.Management.Automation.PSCredential ($UserName, $secpasswd)
    
    # 建立遠端連線
    $session = New-PSSession -ComputerName $ComputerName -Credential $cred

    # 在遠端執行獲取排程任務的指令
    $tasks = Invoke-Command -Session $session -ScriptBlock {
        param($userName, $computerID)
        
        $hostname = hostname
        $xmlDoc = New-Object System.Xml.XmlDocument
        $root = $xmlDoc.CreateElement("ScheduledTasks")
        $xmlDoc.AppendChild($root) | Out-Null

        # 只檢查 UserId
        Get-ScheduledTask | Where-Object { $_.Principal.UserId -eq $userName } | ForEach-Object {
            $task = $_
            $taskInfo = Get-ScheduledTaskInfo -TaskName $task.TaskName -TaskPath $task.TaskPath 2>$null
            $taskXml = Export-ScheduledTask -TaskName $task.TaskName -TaskPath $task.TaskPath

            $tempDoc = New-Object System.Xml.XmlDocument
            $tempDoc.LoadXml($taskXml)

            # 檢查 Author，分割並取最後一段
            $author = $tempDoc.Task.RegistrationInfo.Author
            $authorName = if ($author) { $author.Split("\")[-1] } else { "" }
            if ($authorName -ne $userName) {
                return
            }

            # 添加額外的節點
            $extraInfo = $tempDoc.CreateElement("ExtraInfo")
            $taskName = $tempDoc.CreateElement("TaskName")
            $taskName.InnerText = $task.TaskName
            
            $computerIDElement = $tempDoc.CreateElement("ComputerID")
            $computerIDElement.InnerText = $computerID
            
            $computerName = $tempDoc.CreateElement("ComputerName")
            $computerName.InnerText = $hostname
            $state = $tempDoc.CreateElement("State")
            $state.InnerText = $task.State
            $lastRunTime = $tempDoc.CreateElement("LastRunTime")
            $lastRunTime.InnerText = if ($taskInfo.LastRunTime) { $taskInfo.LastRunTime.ToString() } else { "N/A" }
            $nextRunTime = $tempDoc.CreateElement("NextRunTime")
            $nextRunTime.InnerText = if ($taskInfo.NextRunTime) { $taskInfo.NextRunTime.ToString() } else { "N/A" }
            $lastTaskResult = $tempDoc.CreateElement("LastTaskResult")
            $lastTaskResult.InnerText = if ($taskInfo.LastTaskResult) { $taskInfo.LastTaskResult.ToString() } else { 0 }

            $extraInfo.AppendChild($taskName) | Out-Null
            $extraInfo.AppendChild($computerIDElement) | Out-Null
            $extraInfo.AppendChild($computerName) | Out-Null
            $extraInfo.AppendChild($state) | Out-Null
            $extraInfo.AppendChild($lastRunTime) | Out-Null
            $extraInfo.AppendChild($nextRunTime) | Out-Null
            $extraInfo.AppendChild($lastTaskResult) | Out-Null
            $tempDoc.DocumentElement.AppendChild($extraInfo) | Out-Null

            $importNode = $xmlDoc.ImportNode($tempDoc.DocumentElement, $true)
            $root.AppendChild($importNode) | Out-Null
        }

        return $xmlDoc.OuterXml
    } -ArgumentList $UserName, $ComputerID

    # 輸出結果
    Write-Output $tasks

    # 清理遠端連線
    Remove-PSSession $session
}
catch {
    Write-Error $_.Exception.Message
    exit 1
}
