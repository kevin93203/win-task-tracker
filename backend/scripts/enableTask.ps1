param(
    [Parameter(Mandatory=$true)]
    [string]$TaskName,
    [Parameter(Mandatory=$true)]
    [string]$UserName,
    [Parameter(Mandatory=$true)]
    [string]$Password,
    [Parameter(Mandatory=$true)]
    [string]$ComputerName
)

# 設定輸出編碼為 UTF8
[Console]::OutputEncoding = [Text.Encoding]::UTF8

try {
    # 建立認證物件
    $securePassword = ConvertTo-SecureString $Password -AsPlainText -Force
    $credential = New-Object System.Management.Automation.PSCredential ($UserName, $securePassword)

    # 在遠端電腦上執行啟用任務指令
    $result = Invoke-Command -ComputerName $ComputerName -Credential $credential -ScriptBlock {
        param($taskName)
        
        try {
            $task = Get-ScheduledTask -TaskName $taskName -ErrorAction Stop
            if ($task) {
                Enable-ScheduledTask -TaskName $taskName -ErrorAction Stop
                return @{
                    Success = $true
                    Message = "Successfully enabled task '$taskName'"
                }
            }
        }
        catch {
            return @{
                Success = $false
                Error = $_.Exception.Message
            }
        }
    } -ArgumentList $TaskName

    # 輸出結果為 JSON 格式
    $result | ConvertTo-Json
}
catch {
    @{
        Success = $false
        Error = $_.Exception.Message
    } | ConvertTo-Json
}
