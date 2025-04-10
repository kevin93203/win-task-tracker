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
    [int]$Index,
    [Parameter(Mandatory=$true)]
    [string]$Execute,
    [string]$Arguments = "",
    [string]$WorkingDirectory
)

[Console]::OutputEncoding = [Text.Encoding]::UTF8

try {
    $securePassword = ConvertTo-SecureString $Password -AsPlainText -Force
    $credential = New-Object System.Management.Automation.PSCredential ($UserName, $securePassword)

    $result = Invoke-Command -ComputerName $ComputerName -Credential $credential -ScriptBlock {
        param($taskName, $index, $execute, $arguments, $workingDirectory)

        try {
            $task = Get-ScheduledTask -TaskName $taskName -ErrorAction Stop
            $actions = $task.Actions

            if ($index -lt 0 -or $index -ge $actions.Count) {
                throw "Index out of range"
            }

            $actions[$index].Execute = $execute

            $actions[$index].Arguments = $arguments

            $actions[$index].WorkingDirectory = $workingDirectory

            Set-ScheduledTask -TaskName $taskName -Action $actions -ErrorAction Stop

            return @{
                Success = $true
                Message = "Action at index $index updated successfully"
            }
        } catch {
            return @{
                Success = $false
                Error = $_.Exception.Message
            }
        }
    } -ArgumentList $TaskName, $Index, $Execute, $Arguments, $WorkingDirectory

    $result | ConvertTo-Json -Depth 5
    if ($result.Success -eq $false) { exit 1 } else { exit 0 }
} catch {
    @{
        Success = $false
        Error = $_.Exception.Message
    } | ConvertTo-Json -Depth 5
    exit 1
}
