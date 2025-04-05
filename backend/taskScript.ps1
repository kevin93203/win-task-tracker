[Console]::OutputEncoding = [Text.Encoding]::UTF8

$secpasswd = ConvertTo-SecureString '%s' -AsPlainText -Force
$cred = New-Object System.Management.Automation.PSCredential ('%s', $secpasswd)
$session = New-PSSession -ComputerName "%s" -Credential $cred

$tasks = Invoke-Command -Session $session -ScriptBlock {
    $hostname = hostname
    $userName = "%s"  # 使用者名稱變數
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
        $extraInfo.AppendChild($computerName) | Out-Null
        $extraInfo.AppendChild($state) | Out-Null
        $extraInfo.AppendChild($lastRunTime) | Out-Null
        $extraInfo.AppendChild($nextRunTime) | Out-Null
        $extraInfo.AppendChild($lastTaskResult) | Out-Null
        $tempDoc.DocumentElement.AppendChild($extraInfo) | Out-Null

        # 將任務的根節點直接導入到主 XML 文檔
        $importNode = $xmlDoc.ImportNode($tempDoc.DocumentElement, $true)
        $root.AppendChild($importNode) | Out-Null
    }

    return $xmlDoc.OuterXml
}

# 將 XML 字符串轉換為 UTF-8 編碼的字節，然後再轉回字符串
$utf8Encoding = [System.Text.Encoding]::UTF8
$utf8Bytes = $utf8Encoding.GetBytes($tasks)
$tasksUtf8 = $utf8Encoding.GetString($utf8Bytes)

# 顯示 UTF-8 編碼的 XML
$tasksUtf8

Remove-PSSession $session