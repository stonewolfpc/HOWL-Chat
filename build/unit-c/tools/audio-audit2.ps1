param()
$ErrorActionPreference = 'Continue'
$root = 'C:\Users\StoneWolfPC\HOWL_Chat'
$iso  = Join-Path $root 'build\tmp\unit-c-rerun'
$ev   = Join-Path $root 'build\unit-c\audio-audit'
$uia  = Join-Path $root 'build\unit-c\tools\uia.ps1'
$shot = Join-Path $root 'build\unit-c\tools\shot.ps1'
$log  = Join-Path $ev 'audio-audit2.log'
Set-Location $root
Remove-Item $log -ErrorAction SilentlyContinue

function L([string]$m) {
    $line = '{0}  {1}' -f (Get-Date).ToString('HH:mm:ss.fff'), $m
    Write-Host $line
    Add-Content -Path $log -Value $line -Encoding UTF8
}
function Uia([string[]]$a) {
    $prev = $ErrorActionPreference
    $ErrorActionPreference = 'Continue'
    try { return ((& powershell -NoProfile -ExecutionPolicy Bypass -File $uia @a 2>&1 | Out-String).Trim()) }
    catch { return "UIA_ERROR $($_.Exception.Message)" }
    finally { $ErrorActionPreference = $prev }
}
function WaitTree([string]$pattern, [int]$seconds = 45) {
    $sw = [Diagnostics.Stopwatch]::StartNew()
    while ($sw.Elapsed.TotalSeconds -lt $seconds) {
        $t = Uia @('-Action','tree')
        if ($t -match $pattern) { return $t }
        Start-Sleep -Milliseconds 600
    }
    return ''
}

Add-Type -TypeDefinition @'
using System;
using System.Runtime.InteropServices;
public static class WinSize2 {
    [DllImport("user32.dll")] public static extern bool MoveWindow(IntPtr h, int x, int y, int w, int hh, bool repaint);
}
'@

$p = Start-Process -FilePath (Join-Path $iso 'howl-chat-unit-c.exe') -WorkingDirectory $iso -PassThru
$sw = [Diagnostics.Stopwatch]::StartNew(); $q = $null
while ($sw.Elapsed.TotalSeconds -lt 90) {
    $q = Get-Process -Name 'howl-chat-unit-c' -ErrorAction SilentlyContinue | Where-Object { $_.MainWindowHandle -ne 0 } | Select-Object -First 1
    if ($null -ne $q) { break }
    Start-Sleep -Milliseconds 400
}
if ($null -eq $q) { L 'WINDOW_TIMEOUT'; exit 1 }
L "WINDOW after=$([int]$sw.Elapsed.TotalSeconds)s"
L "UI_READY=$([bool](WaitTree 'menu-button' 60))"
[void][WinSize2]::MoveWindow($q.MainWindowHandle, 0, 0, 1280, 800, $true)
Start-Sleep -Seconds 2
L 'RESIZED 1280x800'

L "SETTINGS_MENU $(Uia @('-Action','invoke','-Name','^Settings$','-Class','menu-button'))"
Start-Sleep -Milliseconds 800
L "AUDIO_ITEM $(Uia @('-Action','invoke','-Name','^Audio Settings$'))"
L "AUDIO_DOC=$([bool](WaitTree 'Audio Settings - HOWL Chat' 30))"
Start-Sleep -Milliseconds 1200

$g = Uia @('-Action','tree','-OutFile',(Join-Path $ev '1280-general-tree.txt'))
foreach ($pat in @('cls=real-settings-overlay[^|]*\|id=audio-settings-overlay[^|]*\|name=[^|]*\|rect=([0-9,-]+)','cls=real-settings-container[^|]*\|name=[^|]*\|rect=([0-9,-]+)','name=Audio Settings - HOWL Chat\|rect=([0-9,-]+)','Window\|cls=wailsWindow[^|]*\|name=HOWL Chat\|rect=([0-9,-]+)')) {
    $m = [regex]::Match($g, $pat)
    if ($m.Success) { L "RECT $($pat.Split('|')[0]) => $($m.Groups[1].Value)" }
}
L "TABS_IN_TREE=$(([regex]::Matches($g,'cls=audio-tab') | Measure-Object).Count)"
L "TRUNCATED_TAB=$($g -match 'Text-to-Spee')"
& powershell -NoProfile -ExecutionPolicy Bypass -File $shot -OutFile (Join-Path $ev '1280-general.png') | Out-Null

L "OPEN_ASR $(Uia @('-Action','invoke','-Name','Speech Recognition'))"
Start-Sleep -Milliseconds 900
$asr = Uia @('-Action','tree','-OutFile',(Join-Path $ev '1280-asr-tree.txt'))
L "ASR_ENGINE_FIELD=$($asr -match 'id=asr-engine')"
L "ASR_CARDS=$(([regex]::Matches($asr,'cls=setting-card') | Measure-Object).Count)"
& powershell -NoProfile -ExecutionPolicy Bypass -File $shot -OutFile (Join-Path $ev '1280-asr.png') | Out-Null

L "OPEN_HW $(Uia @('-Action','invoke','-Name','Hardware'))"
Start-Sleep -Milliseconds 900
$hw = Uia @('-Action','tree','-OutFile',(Join-Path $ev '1280-hardware-tree.txt'))
L "HW_CARDS=$(([regex]::Matches($hw,'cls=setting-card') | Measure-Object).Count)"
L "HW_INPUT_FIELD=$($hw -match 'audio-input-device')"
L "HW_OFFSCR_COUNT=$(([regex]::Matches($hw,'OFFSCR') | Measure-Object).Count)"
& powershell -NoProfile -ExecutionPolicy Bypass -File $shot -OutFile (Join-Path $ev '1280-hardware.png') | Out-Null

L "OPEN_VAD $(Uia @('-Action','invoke','-Name','^VAD$'))"
Start-Sleep -Milliseconds 800
& powershell -NoProfile -ExecutionPolicy Bypass -File $shot -OutFile (Join-Path $ev '1280-vad.png') | Out-Null

$app = Get-Process -Name 'howl-chat-unit-c' -ErrorAction SilentlyContinue | Select-Object -First 1
if ($app) { [void]$app.CloseMainWindow(); [void]$app.WaitForExit(30000) }
Start-Sleep -Seconds 2
L "LEFTOVER=$(@(Get-Process -Name 'howl-chat-unit-c','llama-server' -ErrorAction SilentlyContinue).Count)"
