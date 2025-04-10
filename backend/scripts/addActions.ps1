param (
    [Parameter(Mandatory=$true)]
    [string]$TaskName,
    [Parameter(Mandatory=$true)]
    [string]$UserName,
    [Parameter(Mandatory=$true)]
    [string]$Password,
    [Parameter(Mandatory=$true)]
    [string]$ComputerName,
    [Parameter(Mandatory=$true)]
    [string]$Execute,
    [string]$Arguments = "",
    [string]$WorkingDirectory = ""
)

[Console]::OutputEncoding = [Text.Encoding]::UTF8

try {
    $securePassword = ConvertTo-SecureString $Password -AsPlainText -Force
    $credential = New-Object System.Management.Automation.PSCredential ($UserName, $securePassword)

    $result = Invoke-Command -ComputerName $ComputerName -Credential $credential -ScriptBlock {
        param($taskName, $execute, $arguments, $workingDirectory)

        try {
            $task = Get-ScheduledTask -TaskName $taskName -ErrorAction Stop

            $params = @{
                Execute = $execute
            }
            if ($arguments -ne "") {
                $params["Argument"] = $arguments
            }
            if ($workingDirectory -ne "") {
                $params["WorkingDirectory"] = $workingDirectory
            }

            $newAction = New-ScheduledTaskAction @params

            $actions = $task.Actions + $newAction

            Set-ScheduledTask -TaskName $taskName -Action $actions -ErrorAction Stop

            return @{
                Success = $true
                Message = "Action added successfully"
            }
        } catch {
            return @{
                Success = $false
                Error = $_.Exception.Message
            }
        }
    } -ArgumentList $TaskName, $Execute, $Arguments, $WorkingDirectory

    $result | ConvertTo-Json -Depth 5
    if ($result.Success -eq $false) { exit 1 } else { exit 0 }
} catch {
    @{
        Success = $false
        Error = $_.Exception.Message
    } | ConvertTo-Json -Depth 5
    exit 1
}
