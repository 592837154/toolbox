Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"

$AppConfig = [ordered]@{
    TargetProcessNames = @("SunloginClient", "sunloginclient", "sunlogin_guard")
    TargetTitleKeywords = @("向日葵", "Sunlogin", "Sunflower")
    PollIntervalMs = 50
    StartupGraceMs = 300
    MinimizeCommand = 6
    ShowWindowAsyncSource = @"
using System;
using System.Runtime.InteropServices;

public static class Win32WindowApi
{
    [DllImport("user32.dll")]
    public static extern bool ShowWindowAsync(IntPtr hWnd, int nCmdShow);

    [DllImport("user32.dll")]
    public static extern bool IsWindowVisible(IntPtr hWnd);
}
"@
    CursorSource = @"
using System;
using System.Runtime.InteropServices;

public static class Win32CursorApi
{
    [StructLayout(LayoutKind.Sequential)]
    public struct POINT
    {
        public int X;
        public int Y;
    }

    [DllImport("user32.dll")]
    public static extern bool GetCursorPos(out POINT lpPoint);
}
"@
}

$AppText = [ordered]@{
    CursorPositionReadFailed = "无法读取当前鼠标位置。"
    MonitorStarted = "正在持续监听鼠标移动。检测到移动后会立即最小化向日葵远程窗口。"
    MonitorStopHint = "按 Ctrl+C 可停止监听。"
    NoVisibleTargetWindow = "检测到鼠标移动，但没有找到可见的向日葵远程窗口。"
    MinimizedWindow = "已最小化：{0} {1} {2}"
}

Add-Type -TypeDefinition $AppConfig.ShowWindowAsyncSource
Add-Type -TypeDefinition $AppConfig.CursorSource

function Get-CursorPosition {
    <#
    .SYNOPSIS
    获取当前全局鼠标坐标。

    .DESCRIPTION
    通过 user32.dll 的 GetCursorPos 读取 Windows 桌面级鼠标位置。

    .PARAMETER 无
    本函数不接收参数。

    .OUTPUTS
    返回包含 X、Y 两个整数属性的 PSCustomObject。若系统调用失败则抛出异常。

    .NOTES
    该函数只读取系统状态，不移动鼠标，也不修改窗口状态。
    #>
    $point = New-Object Win32CursorApi+POINT
    $success = [Win32CursorApi]::GetCursorPos([ref]$point)
    if (-not $success) {
        throw $AppText.CursorPositionReadFailed
    }

    [pscustomobject]@{
        X = $point.X
        Y = $point.Y
    }
}

function Test-TargetWindow {
    <#
    .SYNOPSIS
    判断进程窗口是否属于需要最小化的向日葵窗口。

    .DESCRIPTION
    使用集中配置的进程名和标题关键字匹配窗口，并要求窗口句柄有效且窗口可见。

    .PARAMETER Process
    由 Get-Process 返回的进程对象，函数会读取 ProcessName、MainWindowTitle 和 MainWindowHandle。

    .OUTPUTS
    返回 Boolean。匹配目标窗口时返回 True，否则返回 False。

    .NOTES
    本函数只做判断，不会改变窗口状态。
    #>
    param(
        [Parameter(Mandatory = $true)]
        [System.Diagnostics.Process]$Process
    )

    if ($Process.MainWindowHandle -eq [IntPtr]::Zero) {
        return $false
    }

    if (-not [Win32WindowApi]::IsWindowVisible($Process.MainWindowHandle)) {
        return $false
    }

    $nameMatched = $AppConfig.TargetProcessNames -contains $Process.ProcessName
    $titleMatched = $false
    foreach ($keyword in $AppConfig.TargetTitleKeywords) {
        if ($Process.MainWindowTitle -like "*$keyword*") {
            $titleMatched = $true
            break
        }
    }

    return ($nameMatched -or $titleMatched)
}

function Get-TargetWindows {
    <#
    .SYNOPSIS
    查找当前可最小化的向日葵窗口。

    .DESCRIPTION
    扫描当前进程列表，并通过 Test-TargetWindow 过滤出匹配向日葵进程名或标题关键字的可见窗口。

    .PARAMETER 无
    本函数不接收参数。

    .OUTPUTS
    返回零个或多个 System.Diagnostics.Process 对象。

    .NOTES
    该函数会读取进程列表，但不会启动、关闭或修改任何进程。
    #>
    Get-Process | Where-Object { Test-TargetWindow -Process $_ }
}

function Minimize-TargetWindows {
    <#
    .SYNOPSIS
    最小化所有匹配到的向日葵窗口。

    .DESCRIPTION
    调用 ShowWindowAsync，对每个目标窗口发送最小化命令。

    .PARAMETER Targets
    需要最小化的进程对象集合，每个对象必须包含有效的 MainWindowHandle。

    .OUTPUTS
    不返回业务数据。执行过程中会向控制台输出最小化结果。

    .NOTES
    该函数会改变窗口显示状态：匹配窗口会被最小化，但进程不会被结束。
    #>
    param(
        [Parameter(Mandatory = $true)]
        [System.Diagnostics.Process[]]$Targets
    )

    foreach ($target in $Targets) {
        [void][Win32WindowApi]::ShowWindowAsync($target.MainWindowHandle, $AppConfig.MinimizeCommand)
        Write-Host ($AppText.MinimizedWindow -f $target.ProcessName, $target.Id, $target.MainWindowTitle)
    }
}

function Start-MouseMoveMonitor {
    <#
    .SYNOPSIS
    启动持续鼠标移动监听，并在每次移动后最小化向日葵窗口。

    .DESCRIPTION
    记录启动时的鼠标坐标，短暂宽限后按固定间隔轮询鼠标位置；只要坐标发生变化，就查找并最小化目标窗口。
    每次处理移动事件后都会把当前坐标保存为新的基准点，使脚本在窗口被重新放大后仍可继续监听下一次移动。

    .PARAMETER 无
    本函数不接收参数。

    .OUTPUTS
    不返回业务数据。监听状态和动作结果会输出到控制台。

    .NOTES
    该函数会持续占用当前 PowerShell 会话并维护最近一次鼠标坐标状态，直到用户按 Ctrl+C 终止。
    #>
    $lastPosition = Get-CursorPosition
    Start-Sleep -Milliseconds $AppConfig.StartupGraceMs

    Write-Host $AppText.MonitorStarted
    Write-Host $AppText.MonitorStopHint

    while ($true) {
        Start-Sleep -Milliseconds $AppConfig.PollIntervalMs
        $current = Get-CursorPosition
        $moved = ($current.X -ne $lastPosition.X) -or ($current.Y -ne $lastPosition.Y)

        if ($moved) {
            $targets = @(Get-TargetWindows)
            if ($targets.Count -eq 0) {
                Write-Host $AppText.NoVisibleTargetWindow
            }
            else {
                Minimize-TargetWindows -Targets $targets
            }

            $lastPosition = $current
        }
    }
}

Start-MouseMoveMonitor
