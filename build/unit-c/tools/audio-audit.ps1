param()
$ErrorActionPreference = 'Continue'
$root = 'C:\Users\StoneWolfPC\HOWL_Chat'
$iso  = Join-Path $root 'build\tmp\unit-c-rerun'
$ev   = Join-Path $root 'build\unit-c\audio-audit'
$uia  = Join-Path $root 'build\unit-c\tools\uia.ps1'
$shot = Join-Path $root 'build\unit-c\tools\shot.ps1'
$log  = Join-Path $ev 'audio-audit.log'
$cfg  = Join-Path $iso 'config\settings.json'
$startupLog = Join-Path $iso 'config\logs\startup.log'
New-Item -ItemType Directory -Force -Path $ev | Out-Null
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
function InvokeWhen([string]$waitPattern, [string[]]$invokeArgs) {
    if (-not (WaitTree $waitPattern 45)) { return "WAIT_MISS $waitPattern" }
    return (Uia $invokeArgs)
}
function AudioJson([string]$tag) {
    $j = (python -c "import json;d=json.load(open(r'$cfg',encoding='utf-8'));a=d.get('audio',{});print(json.dumps({k:a.get(k) for k in ['asrEnabled','asrEngine','ttsEnabled','ttsEngine','ttsSpeed','ttsVolume','asrLanguage','vadEnabled','pushToTalk','micSensitivity','audioInputDevice','audioOutputDevice']}))" 2>&1 | Out-String).Trim()
    L "SETTINGS[$tag] $j"
}
function Launch {
    $p = Start-Process -FilePath (Join-Path $iso 'howl-chat-unit-c.exe') -WorkingDirectory $iso -PassThru
    $sw = [Diagnostics.Stopwatch]::StartNew(); $q = $null
    while ($sw.Elapsed.TotalSeconds -lt 90) {
        $q = Get-Process -Name 'howl-chat-unit-c' -ErrorAction SilentlyContinue | Where-Object { $_.MainWindowHandle -ne 0 } | Select-Object -First 1
        if ($null -ne $q) { break }
        Start-Sleep -Milliseconds 400
    }
    if ($null -eq $q) { L 'WINDOW_TIMEOUT'; return $null }
    L "WINDOW pid=$($q.Id) after=$([int]$sw.Elapsed.TotalSeconds)s"
    return $q
}
function OpenAudio {
    L "SETTINGS_MENU $(InvokeWhen 'name=Settings' @('-Action','invoke','-Name','^Settings$','-Class','menu-button'))"
    Start-Sleep -Milliseconds 800
    L "AUDIO_ITEM $(InvokeWhen 'name=Audio Settings' @('-Action','invoke','-Name','^Audio Settings$'))"
    L "AUDIO_DOC=$([bool](WaitTree 'Audio Settings - HOWL Chat' 30))"
    Start-Sleep -Milliseconds 1200
}
function CloseApp {
    $p = Get-Process -Name 'howl-chat-unit-c' -ErrorAction SilentlyContinue | Select-Object -First 1
    if ($p) { [void]$p.CloseMainWindow(); [void]$p.WaitForExit(30000) }
    Start-Sleep -Seconds 2
}

# --- backup + seed known values -------------------------------------------
Copy-Item $cfg "$cfg.audit.bak" -Force
$logLinesBefore = 0
if (Test-Path $startupLog) { $logLinesBefore = @(Get-Content $startupLog).Count }
python -c "
import json
p=r'$cfg'
d=json.load(open(p,encoding='utf-8'))
a=d.setdefault('audio',{})
a.update({'asrEnabled':True,'asrEngine':'qwen3-asr-1.7b-local','ttsEnabled':False,'ttsEngine':'piper','ttsSpeed':1.0,'asrLanguage':'de','vadEnabled':True,'pushToTalk':False,'micSensitivity':0.25,'audioInputDevice':'none','audioOutputDevice':'default'})
json.dump(d,open(p,'w',encoding='utf-8'),indent=2)
"
AudioJson 'seeded'

# --- pass 1: launch, inspect, change, save --------------------------------
L 'PASS1_START'
if (-not (Launch)) { exit 1 }
Start-Sleep -Seconds 6
$t = WaitTree 'menu-button' 60
L "UI_READY=$([bool]$t)"
L "STARTUP_STATUS_ASR=$([regex]::Matches($t,'name=Status:[^|]*') | ForEach-Object { $_.Value } | Select-Object -First 3)"
OpenAudio
$tree = Uia @('-Action','tree','-OutFile',(Join-Path $ev 'general-tree.txt'))
L "GENERAL_HAS_ASR_LABEL=$($tree -match 'Enable Speech Recognition')"
L "GENERAL_TABS=$((([regex]::Matches($tree,'\|\w+\|cls=audio-tab[^|]*\|[^|]*\|name=([^|]*)') | ForEach-Object { $_.Groups[1].Value.Trim() }) -join ',') )"
& powershell -NoProfile -ExecutionPolicy Bypass -File $shot -OutFile (Join-Path $ev 'tab-general.png') | Out-Null
foreach ($tab in @('Speech Recognition','Text-to-Speech','VAD','Hardware')) {
    $r = Uia @('-Action','invoke','-Name',[regex]::Escape($tab))
    Start-Sleep -Milliseconds 900
    $nm = ($tab -replace '\s+','').ToLower()
    & powershell -NoProfile -ExecutionPolicy Bypass -File $shot -OutFile (Join-Path $ev ("tab-$nm.png")) | Out-Null
    L "TAB[$tab] $r"
}
# control value checks from the rendered tree
$tt = Uia @('-Action','tree','-OutFile',(Join-Path $ev 'hardware-tree.txt'))
L "HARDWARE_INPUT_PRESENT=$($tt -match 'Input Device')"
L "TTS_SPEED_FIELD=$($tt -match 'tts-speed')"

# change TTS speed via number input, disable ASR, then save
L "SET_SPEED $(Uia @('-Action','invoke','-Name','Text-to-Speech'))"
Start-Sleep -Milliseconds 800
L "SPEED_SET $(Uia @('-Action','text','-Id','tts-speed-value','-Text','1.35'))"
L "ASR_OFF $(Uia @('-Action','toggle','-Id','asr-enabled'))"
Start-Sleep -Milliseconds 400
L "SAVE_CLICK $(Uia @('-Action','invoke','-Id','audio-settings-save'))"
Start-Sleep -Seconds 6
$after = Uia @('-Action','tree')
L "OUTCOME_SAVED=$($after -match 'name=Saved')"
L "OUTCOME_APPLIED=$($after -match 'name=Applied')"
L "MODAL_CLOSED=$(-not ($after -match 'Audio Settings - HOWL Chat'))"
AudioJson 'after-save'
# startup audio activity evidence
if (Test-Path $startupLog) {
    $all = @(Get-Content $startupLog)
    $new = @($all | Select-Object -Skip $logLinesBefore)
    $hits = @($new | Where-Object { $_ -match 'asr|ASR|audio|llama|whisper|tts|TTS' })
    L "STARTUP_AUDIO_LINES=$($hits.Count)"
    foreach ($h in $hits | Select-Object -First 12) { L "  STARTUP $h" }
}
CloseApp

# --- pass 2: relaunch restore --------------------------------------------
L 'PASS2_START'
if (-not (Launch)) { exit 1 }
Start-Sleep -Seconds 6
OpenAudio
$r2 = Uia @('-Action','invoke','-Name','Text-to-Speech')
Start-Sleep -Milliseconds 900
$tree2 = Uia @('-Action','tree','-OutFile',(Join-Path $ev 'restore-tts-tree.txt'))
& powershell -NoProfile -ExecutionPolicy Bypass -File $shot -OutFile (Join-Path $ev 'restore-tts.png') | Out-Null
L "RESTORE_TTS_TAB $r2"
L "RESTORE_SPEED_FIELD=$($tree2 -match 'tts-speed')"
AudioJson 'after-relaunch'

# --- pass 3: UI audit at 1920x1080 and light theme ------------------------
Add-Type -TypeDefinition @'
using System;
using System.Runtime.InteropServices;
public static class WinSize {
    [DllImport("user32.dll")] public static extern bool MoveWindow(IntPtr h, int x, int y, int w, int hh, bool repaint);
}
'@
$proc = Get-Process -Name 'howl-chat-unit-c' -ErrorAction SilentlyContinue | Where-Object { $_.MainWindowHandle -ne 0 } | Select-Object -First 1
if ($proc) {
    [void][WinSize]::MoveWindow($proc.MainWindowHandle, 0, 0, 1920, 1080, $true)
    Start-Sleep -Seconds 2
    L 'RESIZED 1920x1080'
}
& powershell -NoProfile -ExecutionPolicy Bypass -File $shot -OutFile (Join-Path $ev 'tts-1920.png') | Out-Null
$big = Uia @('-Action','tree','-OutFile',(Join-Path $ev 'tts-1920-tree.txt'))
$modal = [regex]::Match($big, 'name=Audio Settings - HOWL Chat\|rect=([0-9-]+),([0-9-]+),([0-9]+),([0-9]+)')
if ($modal.Success) { L "MODAL_RECT_1920 $($modal.Groups[1].Value),$($modal.Groups[2].Value) $($modal.Groups[3].Value)x$($modal.Groups[4].Value)" }
L 'LIGHT_THEME_START'
L "VIEW_MENU $(Uia @('-Action','invoke','-Name','^View$','-Class','menu-button'))"
Start-Sleep -Milliseconds 700
L "LIGHT_CLICK $(Uia @('-Action','invoke','-Name','^Light Theme$'))"
Start-Sleep -Seconds 2
& powershell -NoProfile -ExecutionPolicy Bypass -File $shot -OutFile (Join-Path $ev 'tts-1920-light.png') | Out-Null
L "LIGHT_THEME_SHOT=tts-1920-light.png"
CloseApp

# --- restore --------------------------------------------------------------
Copy-Item "$cfg.audit.bak" $cfg -Force
Remove-Item "$cfg.audit.bak" -Force -ErrorAction SilentlyContinue
AudioJson 'restored'
L "LEFTOVER=$(@(Get-Process -Name 'howl-chat-unit-c','llama-server' -ErrorAction SilentlyContinue).Count)"
