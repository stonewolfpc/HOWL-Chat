param(
    [string]$OutFile = 'build/unit-c/shot.png',
    # Defaults to the historical unit-c executable; pass howl-chat to capture the
    # packaged production build on an isolated root.
    [string]$ProcessName = 'howl-chat-unit-c'
)
$ErrorActionPreference = 'Stop'
Add-Type -AssemblyName System.Drawing
$w32 = @'
using System;
using System.Runtime.InteropServices;
public static class SW32 {
    [DllImport("user32.dll")] public static extern bool GetWindowRect(IntPtr h, out RECT r);
    [DllImport("user32.dll")] public static extern bool SetForegroundWindow(IntPtr h);
    [DllImport("user32.dll")] public static extern bool SetProcessDPIAware();
    [StructLayout(LayoutKind.Sequential)] public struct RECT { public int L, T, R, B; }
}
'@
Add-Type -TypeDefinition $w32
[void][SW32]::SetProcessDPIAware()
$p = Get-Process -Name $ProcessName -ErrorAction SilentlyContinue | Where-Object { $_.MainWindowHandle -ne 0 } | Select-Object -First 1
if (-not $p) { throw 'app window not found' }
[void][SW32]::SetForegroundWindow($p.MainWindowHandle)
Start-Sleep -Milliseconds 400
$r = New-Object SW32+RECT
[void][SW32]::GetWindowRect($p.MainWindowHandle, [ref]$r)
$w = $r.R - $r.L; $h = $r.B - $r.T
if ($w -le 0 -or $h -le 0) { throw 'bad window rect' }
$bmp = New-Object System.Drawing.Bitmap $w, $h
$g = [System.Drawing.Graphics]::FromImage($bmp)
$g.CopyFromScreen($r.L, $r.T, 0, 0, (New-Object System.Drawing.Size $w, $h))
$dir = Split-Path -Parent $OutFile
if ($dir -and -not (Test-Path $dir)) { New-Item -ItemType Directory -Path $dir | Out-Null }
$bmp.Save((Resolve-Path -LiteralPath $dir).Path + '\' + (Split-Path -Leaf $OutFile), [System.Drawing.Imaging.ImageFormat]::Png)
$g.Dispose(); $bmp.Dispose()
Write-Output "SAVED $OutFile ${w}x${h} window=$($r.L),$($r.T)"
