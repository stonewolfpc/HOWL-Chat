param(
    [ValidateSet('precheck','official','restore','pick')]
    [string]$Phase = 'precheck',
    [string]$RenderMatch = 'Speakers'
)
$ErrorActionPreference = 'Stop'

$root = 'C:\Users\StoneWolfPC\HOWL_Chat'
$iso  = Join-Path $root 'build\tmp\unit-c-rerun'
$prof = Join-Path $iso 'profile'
$ev   = Join-Path $root 'build\unit-c\coexistence'
$uia  = Join-Path $root 'build\unit-c\tools\uia.ps1'
$stimWav = Join-Path $ev 'stimulus.wav'
$stimRel = 'build/unit-c/coexistence/stimulus.wav'
$log  = Join-Path $ev "$Phase-run.log"
New-Item -ItemType Directory -Force -Path $ev | Out-Null
Set-Location $root

function L([string]$m) {
    $line = '{0}  {1}' -f (Get-Date).ToString('HH:mm:ss.fff'), $m
    Write-Host $line
    Add-Content -Path $log -Value $line -Encoding UTF8
}

function Uia([string[]]$a) {
    $prev = $ErrorActionPreference
    $ErrorActionPreference = 'Continue'
    try {
        $out = & powershell -NoProfile -ExecutionPolicy Bypass -File $uia @a 2>&1 | Out-String
        return $out.Trim()
    } catch {
        return "UIA_ERROR $($_.Exception.Message)"
    } finally {
        $ErrorActionPreference = $prev
    }
}

function ParseRows([string]$text) {
    $rows = @()
    foreach ($ln in ($text -split "`r?`n")) {
        $m = [regex]::Match($ln, '^\[(\d+)\] (\S+) cls="(.*?)" id="(.*?)" name="(.*?)" rect=(-?\d+),(-?\d+),(-?\d+),(-?\d+)')
        if (-not $m.Success) { continue }
        $rows += [pscustomobject]@{
            Index = [int]$m.Groups[1].Value
            Type  = $m.Groups[2].Value
            Class = $m.Groups[3].Value
            Id    = $m.Groups[4].Value
            Name  = $m.Groups[5].Value
            X     = [int]$m.Groups[6].Value
            Y     = [int]$m.Groups[7].Value
            W     = [int]$m.Groups[8].Value
            H     = [int]$m.Groups[9].Value
            Off   = ($ln -match 'offscreen=True')
        }
    }
    return $rows
}

function Health([int]$port) {
    try { $r = Invoke-RestMethod -Uri ("http://127.0.0.1:$port/health") -TimeoutSec 2; return ($r.status -eq 'ok') } catch { return $false }
}

function Gpu {
    try { return ((& nvidia-smi --query-gpu=memory.used,memory.total,utilization.gpu --format=csv,noheader) -join ' | ').Trim() } catch { return 'nvidia-smi unavailable' }
}

Add-Type -TypeDefinition @'
using System; using System.Runtime.InteropServices;
public static class CX {
  [DllImport("user32.dll")] public static extern bool SetCursorPos(int x,int y);
  [DllImport("user32.dll")] public static extern void mouse_event(uint f,uint dx,uint dy,int d,UIntPtr e);
  [DllImport("user32.dll")] public static extern bool GetWindowRect(IntPtr h, out RECT r);
  public struct RECT { public int L,T,R,B; }
}
'@

function Shot([string]$path) {
    $p = Get-Process -Name 'howl-chat-unit-c' -ErrorAction SilentlyContinue | Where-Object { $_.MainWindowHandle -ne 0 } | Select-Object -First 1
    if ($null -eq $p) { L "SHOT_FAIL no window"; return }
    $r = New-Object CX+RECT
    [void][CX]::GetWindowRect($p.MainWindowHandle, [ref]$r)
    Add-Type -AssemblyName System.Drawing
    $w = $r.R - $r.L; $h = $r.B - $r.T
    $bmp = New-Object System.Drawing.Bitmap $w, $h
    $g = [System.Drawing.Graphics]::FromImage($bmp)
    $g.CopyFromScreen($r.L, $r.T, 0, 0, (New-Object System.Drawing.Size $w, $h))
    $bmp.Save($path, [System.Drawing.Imaging.ImageFormat]::Png)
    $g.Dispose(); $bmp.Dispose()
    L "SHOT=$path"
}

# ---- isolation -------------------------------------------------------------
$env:USERPROFILE  = $prof
$env:APPDATA      = Join-Path $prof 'AppData\Roaming'
$env:LOCALAPPDATA = Join-Path $prof 'AppData\Local'
$env:TEMP         = Join-Path $prof 'Temp'
$env:TMP          = $env:TEMP
foreach ($d in @($env:APPDATA, $env:LOCALAPPDATA, $env:TEMP)) { New-Item -ItemType Directory -Force -Path $d | Out-Null }
$startupLogs = @(
    (Join-Path $env:APPDATA 'HOWLChat\config\logs\startup.log'),
    (Join-Path $iso 'config\logs\startup.log')
) | Select-Object -Unique
L "PHASE=$Phase ISOLATED_ROOT=$iso PROFILE=$prof"

$priorVolume = (& powershell -NoProfile -ExecutionPolicy Bypass -File (Join-Path $root 'build\unit-c\tools\AudioCtl.ps1') -Action getvolume 2>&1 | Out-String)
L "RENDER_VOLUME_BEFORE=$($priorVolume.Trim())"

# ---- locate ASR / chat servers --------------------------------------------
function ServerInfo {
    $apps = Get-CimInstance Win32_Process -Filter "Name='howl-chat-unit-c.exe'"
    if (-not $apps) { return @() }
    $ids = $apps.ProcessId
    $kids = Get-CimInstance Win32_Process -Filter "Name='llama-server.exe'" | Where-Object { $ids -contains $_.ParentProcessId }
    return @($kids | ForEach-Object {
        $port = ''
        $m = [regex]::Match($_.CommandLine, '--port\s+(\d+)'); if ($m.Success) { $port = $m.Groups[1].Value }
        $mm = [regex]::Match($_.CommandLine, '--mmproj\s+"?([^"]+?)"?\s')
        [pscustomobject]@{ Pid = $_.ProcessId; Port = $port; Mmproj = $mm.Groups[1].Value; Cmd = $_.CommandLine }
    })
}

function DumpServers([string]$tag) {
    $srvs = ServerInfo
    L "SERVERS[$tag] count=$($srvs.Count)"
    foreach ($s in $srvs) { L "  pid=$($s.Pid) port=$($s.Port) mmproj=$($s.Mmproj)" }
    return $srvs
}

function WaitRescan([int]$timeoutSec = 240) {
    $sw = [Diagnostics.Stopwatch]::StartNew()
    while ($sw.Elapsed.TotalSeconds -lt $timeoutSec) {
        foreach ($sl in $startupLogs) {
            if (-not (Test-Path $sl)) { continue }
            $txt = Get-Content $sl -Raw -ErrorAction SilentlyContinue
            if ($null -eq $txt) { continue }
            # only honour lines produced by THIS launch; the log persists across runs
            $off = 0
            if ($script:logOffsets -and $script:logOffsets.ContainsKey($sl)) { $off = [int]$script:logOffsets[$sl] }
            if ($txt.Length -lt $off) { $off = 0 }
            $fresh = $txt.Substring($off)
            if ($fresh -match 'rescan: done G:\\AI GGUF\\voice') {
                L ("RESCAN_DONE after {0:N1}s in {1}" -f $sw.Elapsed.TotalSeconds, $sl); return $true
            }
        }
        Start-Sleep -Milliseconds 1500
    }
    L "RESCAN_TIMEOUT after $timeoutSec s"; return $false
}

function WaitWindow([int]$timeoutSec = 60) {
    $sw = [Diagnostics.Stopwatch]::StartNew()
    while ($sw.Elapsed.TotalSeconds -lt $timeoutSec) {
        $t = Uia @('-Action','rect')
        if ($t -match 'WINDOW rect=') { L "WINDOW_READY $($t -replace "`r?`n",' | ')"; return $true }
        Start-Sleep -Milliseconds 1000
    }
    L "WINDOW_TIMEOUT"; return $false
}

function OpenPicker([int]$timeoutSec = 30) {
    $sw = [Diagnostics.Stopwatch]::StartNew()
    while ($sw.Elapsed.TotalSeconds -lt $timeoutSec) {
        $r = Uia @('-Action','click','-Class','model-select-trigger','-Index','0')
        if ($r -match 'CLICKED|INVOKED') {
            L "PICKER_OPEN $r after $([int]$sw.Elapsed.TotalSeconds)s"
            Start-Sleep -Milliseconds 1500
            return $true
        }
        Start-Sleep -Milliseconds 1000
    }
    L "PICKER_TRIGGER_NOT_FOUND after $timeoutSec s"
    return $false
}

function WindowRect {
    $t = Uia @('-Action','rect')
    $m = [regex]::Match($t, 'WINDOW rect=(-?\d+),(-?\d+),(\d+),(\d+)')
    if (-not $m.Success) { throw "window rect unreadable: $t" }
    $l = [int]$m.Groups[1].Value; $tp = [int]$m.Groups[2].Value
    $w = [int]$m.Groups[3].Value; $hh = [int]$m.Groups[4].Value
    return [pscustomobject]@{ L = $l; T = $tp; R = $l + $w; B = $tp + $hh; W = $w; H = $hh }
}

function PickModel([string]$pattern) {
    $wr = WindowRect
    L "WINDOW_BOUNDS $($wr.L),$($wr.T)..$($wr.R),$($wr.B)"
    for ($step = 0; $step -lt 25; $step++) {
        $rows = @(ParseRows (Uia @('-Action','find','-Name',$pattern,'-Max','40')))
        L "SCAN step=$step matches=$($rows.Count)"
        foreach ($r in $rows) { L ("  ROW cls={0} off={1} name={2} rect={3},{4},{5},{6}" -f $r.Class, $r.Off, $r.Name, $r.X, $r.Y, $r.W, $r.H) }
        $hit = $rows | Where-Object { $_.H -gt 0 -and $_.W -gt 10 -and $_.Y -ge $wr.T -and ($_.Y + $_.H) -le $wr.B -and $_.X -ge $wr.L -and $_.X -le $wr.R } | Select-Object -First 1
        if ($hit) {
            $cx = $hit.X + [int]($hit.W / 2); $cy = $hit.Y + [int]($hit.H / 2)
            L "PICK_TARGET name=$($hit.Name) rect=$($hit.X),$($hit.Y),$($hit.W),$($hit.H) click=$cx,$cy"
            L "PICK_CLICK $(Uia @('-Action','click','-X',[string]$cx,'-Y',[string]$cy))"; Start-Sleep -Milliseconds 900
            $trig = Uia @('-Action','find','-Class','model-select-trigger','-Max','5')
            L "TRIGGER_AFTER_PICK $($trig -replace "`r?`n",' | ')"
            return $true
        }
        if ($rows.Count -eq 0) { L 'MODEL_ROW_ABSENT'; return $false }
        $anchor = $rows | Sort-Object Y | Select-Object -First 1
        $dir = if ($anchor.Y -lt $wr.T) { 120 } else { -120 }
        $wx = [Math]::Min([Math]::Max($anchor.X + 20, $wr.L + 20), $wr.R - 20)
        $wy = [Math]::Min([Math]::Max($wr.T + 400, $wr.T + 40), $wr.B - 40)
        L "WHEEL step=$step dir=$dir at $wx,$wy anchorY=$($anchor.Y)"
        [void][CX]::SetCursorPos($wx, $wy); Start-Sleep -Milliseconds 150
        for ($k = 0; $k -lt 3; $k++) { [CX]::mouse_event(0x0800, 0, 0, $dir, [UIntPtr]::Zero); Start-Sleep -Milliseconds 120 }
        Start-Sleep -Milliseconds 400
    }
    return $false
}

# ---- dictation -------------------------------------------------------------
function Analyze([string]$path) {
    return (& powershell -NoProfile -ExecutionPolicy Bypass -File (Join-Path $root 'build\unit-c\mic-diagnosis\analyze.ps1') -Path $path 2>&1 | Out-String).Trim()
}

function Stimulus([string]$text) {
    $out = (& powershell -NoProfile -ExecutionPolicy Bypass -File (Join-Path $root 'build\unit-c\tools\MakeStimulus.ps1') -OutFile $stimRel -Text $text 2>&1 | Out-String)
    L ("STIMULUS " + ($out -replace "`r?`n", ' | '))
    $m = [regex]::Match($out, 'DURATION=([0-9.]+)s')
    if (-not $m.Success) { throw 'stimulus duration not parsed' }
    return [double]$m.Groups[1].Value
}

function RenderGet { return (((& powershell -NoProfile -ExecutionPolicy Bypass -File (Join-Path $root 'build\unit-c\tools\AudioCtl.ps1') -Action getdefault 2>&1 | Out-String).Trim()) -replace '^DEFAULT_RENDER=','') }
function RenderSet([string]$match, [string]$roles = '0,1,2') { return ((& powershell -NoProfile -ExecutionPolicy Bypass -File (Join-Path $root 'build\unit-c\tools\AudioCtl.ps1') -Action setdefault -Match $match -Roles $roles 2>&1 | Out-String).Trim()) }
function VolGet { return ((& powershell -NoProfile -ExecutionPolicy Bypass -File (Join-Path $root 'build\unit-c\tools\AudioCtl.ps1') -Action getvolume 2>&1 | Out-String).Trim()) }
function VolSet([int]$pct) { return ((& powershell -NoProfile -ExecutionPolicy Bypass -File (Join-Path $root 'build\unit-c\tools\AudioCtl.ps1') -Action setvolume -Percent $pct 2>&1 | Out-String).Trim()) }

function Dictate([string]$tag, [string]$renderMatch = 'Headset (WH-1000XM4)') {
    $text = 'The blue lantern guild meets at seven in the copper harbour.'
    $dur = Stimulus $text

    # The microphone runs on the Bluetooth HFP endpoint. Playing the stimulus on
    # the A2DP headphones endpoint leaves the headset silent while it records, so
    # target the HFP render endpoint that shares the capture link.
    $beforeRender = RenderGet
    L "RENDER_BEFORE=$beforeRender"
    $switched = $false
    if ($beforeRender -notmatch [regex]::Escape($renderMatch)) {
        L "RENDER_SWITCH $(RenderSet $renderMatch)"
        $switched = $true
    }
    $volTxt = VolGet
    L "VOLUME_BEFORE=$volTxt"
    $priorPct = 100
    $vm = [regex]::Match($volTxt, 'VOLUME=(\d+)%')
    if ($vm.Success) { $priorPct = [int]$vm.Groups[1].Value }
    L "VOLUME_SET $(VolSet 100)"
    # the composer toolbar only exposes the mic once the view is interactive
    $swMic = [Diagnostics.Stopwatch]::StartNew()
    $micRows = @()
    while ($swMic.Elapsed.TotalSeconds -lt 30) {
        $micRows = @(ParseRows (Uia @('-Action','find','-Class','mic-button','-Max','20')))
        if ($micRows.Count -gt 0) { break }
        Start-Sleep -Milliseconds 700
    }
    L "MIC_ELEMENT_FOUND=$($micRows.Count) after $([int]$swMic.Elapsed.TotalSeconds)s"
    if ($micRows.Count -eq 0) { throw 'mic control never appeared' }

    $on = ''
    for ($i = 0; $i -lt 4; $i++) {
        $on = Uia @('-Action','invoke','-Class','mic-button','-Index','0')
        Start-Sleep -Milliseconds 600
        if (@(ParseRows (Uia @('-Action','find','-Class','mic-button recording','-Max','20'))).Count -gt 0) { break }
    }
    L "MIC_ON $on"
    L "MIC_RECORDING_STATE=$(@(ParseRows (Uia @('-Action','find','-Class','mic-button recording','-Max','20'))).Count)"

    $player = Start-Process -FilePath 'powershell' -ArgumentList @('-NoProfile','-ExecutionPolicy','Bypass','-File',(Join-Path $root 'build\unit-c\tools\PlayWav.ps1'),'-Path',$stimRel) -PassThru -WindowStyle Hidden
    L 'PLAYBACK_STARTED'
    Start-Sleep -Milliseconds ([int](($dur + 1.2) * 1000))
    L "PLAYBACK_EXITED=$($player.HasExited)"

    $off = ''
    for ($i = 0; $i -lt 4; $i++) {
        $off = Uia @('-Action','invoke','-Class','mic-button','-Index','0')
        Start-Sleep -Milliseconds 700
        if (@(ParseRows (Uia @('-Action','find','-Class','mic-button recording','-Max','20'))).Count -eq 0) { break }
    }
    L "MIC_OFF $off"
    L "VOLUME_RESTORE $(VolSet $priorPct)"
    if ($switched) { L "RENDER_RESTORE $(RenderSet (($beforeRender -split '\|')[0]))" }

    # The app deletes its temp WAV once the recognizer answers, so sample the ASR
    # temp directory while the request is still in flight.
    $captureDst = Join-Path $ev "$tag-capture.wav"
    if (Test-Path $captureDst) { Remove-Item $captureDst -Force }
    $asrDirs = @((Join-Path $env:TEMP 'howl-asr'), (Join-Path $env:LOCALAPPDATA 'Temp\howl-asr'), (Join-Path $env:USERPROFILE 'Temp\howl-asr')) | Select-Object -Unique
    L "ASR_DIR_EXISTS=" + (($asrDirs | Where-Object { Test-Path $_ }) -join ';')
    $swC = [Diagnostics.Stopwatch]::StartNew()
    while ($swC.Elapsed.TotalSeconds -lt 45 -and -not (Test-Path $captureDst)) {
        foreach ($d in $asrDirs) {
            if (-not (Test-Path $d)) { continue }
            $w = Get-ChildItem $d -Filter '*.wav' -File -ErrorAction SilentlyContinue | Sort-Object LastWriteTime -Descending | Select-Object -First 1
            if ($null -ne $w) { try { Copy-Item $w.FullName $captureDst -Force } catch {} ; break }
        }
        if (-not (Test-Path $captureDst)) {
            $w = Get-ChildItem (Join-Path $prof 'Temp') -Recurse -Filter 'recording-*.wav' -File -ErrorAction SilentlyContinue | Sort-Object LastWriteTime -Descending | Select-Object -First 1
            if ($null -ne $w) { try { Copy-Item $w.FullName $captureDst -Force } catch {} }
        }
        Start-Sleep -Milliseconds 100
    }
    if (Test-Path $captureDst) {
        L "CAPTURE sha256=$((Get-FileHash $captureDst -Algorithm SHA256).Hash)"
        L "CAPTURE_ANALYSIS $(Analyze $captureDst)"
    } else { L 'CAPTURE_NOT_OBSERVED during transcription' }

    $sw = [Diagnostics.Stopwatch]::StartNew()
    $transcript = ''
    while ($sw.Elapsed.TotalSeconds -lt 60) {
        Start-Sleep -Milliseconds 1500
        $v = Uia @('-Action','value','-Id','message-input')
        $m = [regex]::Match($v, 'VALUE=(.*)')
        if ($m.Success -and $m.Groups[1].Value.Trim().Length -gt 0) { $transcript = $m.Groups[1].Value.Trim(); break }
        if ($sw.Elapsed.TotalSeconds -gt 25) { L "COMPOSER_WAIT t=$([int]$sw.Elapsed.TotalSeconds)s still empty" }
    }
    L "TRANSCRIPT[$tag]=$transcript"
    Set-Content -Path (Join-Path $ev "$tag-transcript.txt") -Value $transcript -Encoding UTF8

    return $transcript
}

# ---- send / reply ----------------------------------------------------------
function MessageRows {
    return ParseRows (Uia @('-Action','find','-Class','^chat-message','-Max','500'))
}

function SendAndWait([string]$tag, [int]$timeoutSec = 150) {
    $before = (MessageRows).Count
    L "SEND[$tag] messages_before=$before"
    L "SEND_INVOKE $(Uia @('-Action','invoke','-Class','send-button','-Index','0'))"
    $sw = [Diagnostics.Stopwatch]::StartNew()
    while ($sw.Elapsed.TotalSeconds -lt $timeoutSec) {
        Start-Sleep -Milliseconds 2000
        $rows = MessageRows
        $streaming = @($rows | Where-Object { $_.Class -match 'streaming' })
        if ($rows.Count -ge ($before + 2) -and $streaming.Count -eq 0) {
            L ("REPLY_DONE[{0}] after {1:N0}s messages={2}" -f $tag, $sw.Elapsed.TotalSeconds, $rows.Count)
            return $rows
        }
    }
    L "REPLY_TIMEOUT[$tag] after $timeoutSec s"
    return (MessageRows)
}

function DumpMessages([string]$tag, $rows) {
    $lines = $rows | ForEach-Object { "  [{0}] cls={1} name={2}" -f $_.Index, $_.Class, $_.Name }
    $path = Join-Path $ev "$tag-messages.txt"
    Set-Content -Path $path -Value ($lines -join "`n") -Encoding UTF8
    L "MESSAGES[$tag]"
    foreach ($l in $lines) { L $l }
}

# ---- launch ----------------------------------------------------------------
function Launch {
    $script:logOffsets = @{}
    foreach ($sl in $startupLogs) {
        if (Test-Path $sl) { $script:logOffsets[$sl] = (Get-Item $sl).Length } else { $script:logOffsets[$sl] = 0 }
    }
    $p = Start-Process -FilePath (Join-Path $iso 'howl-chat-unit-c.exe') -WorkingDirectory $iso -PassThru
    L "APP_PID=$($p.Id) LOG_OFFSETS=$((($script:logOffsets.GetEnumerator() | ForEach-Object { "$($_.Key.Substring($_.Key.Length-40))=$($_.Value)" }) -join ' '))"
    if (-not (WaitWindow 90)) { throw 'window never appeared' }
    Start-Sleep -Milliseconds 1500
    $srvs = @()
    $asr = $null
    $sw = [Diagnostics.Stopwatch]::StartNew()
    while ($sw.Elapsed.TotalSeconds -lt 120) {
        $srvs = ServerInfo
        $asr = $srvs | Where-Object { $_.Cmd -match 'Qwen3-ASR' } | Select-Object -First 1
        if ($null -ne $asr) { break }
        Start-Sleep -Milliseconds 1000
    }
    $srvs = DumpServers 'startup'
    if ($null -eq $asr) { throw ('ASR server not found within 120s; servers=' + $srvs.Count) }
    L ("ASR_DISCOVERED pid={0} port={1} after {2:N1}s" -f $asr.Pid, $asr.Port, $sw.Elapsed.TotalSeconds)
    $sw = [Diagnostics.Stopwatch]::StartNew()
    while ($sw.Elapsed.TotalSeconds -lt 90) { if (Health ([int]$asr.Port)) { break }; Start-Sleep -Milliseconds 1000 }
    L ("ASR_READY port={0} after {1:N1}s health={2}" -f $asr.Port, $sw.Elapsed.TotalSeconds, (Health ([int]$asr.Port)))
    L "GPU_AT_STARTUP $(Gpu)"
    return $p
}

function CloseApp {
    $p = Get-Process -Name 'howl-chat-unit-c' -ErrorAction SilentlyContinue | Select-Object -First 1
    if ($null -eq $p) { L 'CLOSE no app process'; return }
    $rec = @(ParseRows (Uia @('-Action','find','-Class','mic-button recording','-Max','20')))
    if ($rec.Count -gt 0) { L "CLOSE_MIC_STILL_RECORDING - toggling off: $(Uia @('-Action','invoke','-Class','mic-button','-Index','0'))"; Start-Sleep -Milliseconds 1500 }
    L "CLOSE_REQUEST pid=$($p.Id) main=$($p.MainWindowHandle)"
    [void]$p.CloseMainWindow()
    if (-not $p.WaitForExit(40000)) { L 'CLOSE_TIMEOUT still running'; return }
    L "EXITED code=$($p.ExitCode)"
    Start-Sleep -Milliseconds 3000
    $left = @(Get-Process -Name 'howl-chat-unit-c','llama-server' -ErrorAction SilentlyContinue)
    L "LEFTOVER=$($left.Count)"
    foreach ($l in $left) { L "  leftover pid=$($l.Id) name=$($l.ProcessName)" }
}

# ---- restore phase ---------------------------------------------------------
function RestorePhase {
    $p = Launch
    if (-not (WaitRescan 200)) { throw 'rescan timeout' }
    Start-Sleep -Milliseconds 3000
    $tree = Uia @('-Action','tree','-Depth','60')
    $treePath = Join-Path $ev 'restore-tree.txt'
    Set-Content -Path $treePath -Value $tree -Encoding UTF8
    L "RESTORE_TREE=$treePath elements=$((($tree -split "`r?`n") | Where-Object { $_ -match '^\s*\|' }).Count)"
    Set-Content -Path (Join-Path $ev 'restore-tree-lower.txt') -Value (Uia @('-Action','find','-Id','message-input')) -Encoding UTF8
    $rows = MessageRows
    DumpMessages 'restore' $rows
    $srv = DumpServers 'restore'
    $chat = $srv | Where-Object { $_.Port -eq '8080' }
    L "RESTORE_CHAT_SERVER_PRESENT=$([bool]$chat)"
    L "RESTORE_GPU $(Gpu)"
    Shot (Join-Path $ev 'restore.png')
    CloseApp
}

# ---- main ------------------------------------------------------------------
try {
    if ($Phase -eq 'restore') {
        RestorePhase
        exit 0
    }

    $app = Launch

    if ($Phase -eq 'pick') {
        if (-not (WaitRescan 240)) { throw 'rescan timeout' }
        Start-Sleep -Milliseconds 2500
        if (-not (OpenPicker)) { throw 'picker trigger not found' }
        $ok = PickModel 'E4B'
        L "PICK_RESULT=$ok"
        Shot (Join-Path $ev 'pick-probe.png')
        CloseApp
        exit 0
    }

    if ($Phase -eq 'official') {
        if (-not (WaitRescan 200)) { throw 'rescan timeout' }
        Start-Sleep -Milliseconds 2000
        if (-not (OpenPicker)) { throw 'picker trigger not found' }
        if (-not (PickModel 'E4B')) { throw 'E4B option not pointer-selectable' }
        Start-Sleep -Milliseconds 800
        Set-Content -Path (Join-Path $ev 'picker-selected.txt') -Value (Uia @('-Action','find','-Class','model-select','-Max','20')) -Encoding UTF8
        Shot (Join-Path $ev 'picked.png')
        L "LOAD_CLICK $(Uia @('-Action','invoke','-Name','^Load$','-Class','load-button','-Index','0'))"
        $sw = [Diagnostics.Stopwatch]::StartNew()
        $ok = $false
        while ($sw.Elapsed.TotalSeconds -lt 150) { if (Health 8080) { $ok = $true; break }; Start-Sleep -Milliseconds 1500 }
        L ("CHAT_READY={0} after {1:N1}s" -f $ok, $sw.Elapsed.TotalSeconds)
        if (-not $ok) { throw 'chat server never became healthy' }
        Start-Sleep -Milliseconds 2000
        $srv = DumpServers 'loaded'
        foreach ($s in $srv) { L "  CMD pid=$($s.Pid): $($s.Cmd)" }
        $srv | ForEach-Object { L "CHAT_MODEL_ARG pid=$($_.Pid) $(($_.Cmd -split ' ' | Where-Object { $_ -match '\.gguf' }) -join ' ; ')" }
        L "GPU_BOTH_LOADED $(Gpu)"
        L "RAM_APP_MB=$([int]((Get-Process -Name 'howl-chat-unit-c').WorkingSet64 / 1MB))"
        Set-Content -Path (Join-Path $ev 'loaded-tree.txt') -Value (Uia @('-Action','find','-Class','load-button|model-select-trigger','-Max','20')) -Encoding UTF8
        Shot (Join-Path $ev 'both-loaded.png')

        $t1 = Dictate 'dictated' $RenderMatch
        Shot (Join-Path $ev 'dictated.png')
        if ($t1.Length -gt 0) {
            $rows = SendAndWait 'dictated'
            DumpMessages 'dictated' $rows
            Shot (Join-Path $ev 'dictated-reply.png')
        } else {
            L 'DICTATION_EMPTY - official acoustic step failed; continuing to typed regression'
        }

        L "GPU_AFTER_DICTATION $(Gpu)"
        $srv2 = DumpServers 'after-dictation'
        L "BOTH_RESIDENT_AFTER_DICTATION=$( ($srv2 | Where-Object { $_.Port -eq '8080' }).Count -eq 1 -and ($srv2 | Where-Object { $_.Cmd -match 'Qwen3-ASR' }).Count -eq 1 )"

        $typed = 'Reply with exactly this word and nothing else: LEDGER'
        L "TYPED_SET $(Uia @('-Action','text','-Id','message-input','-Text',$typed))"
        L "TYPED_VALUE $(Uia @('-Action','value','-Id','message-input'))"
        $rows2 = SendAndWait 'typed'
        DumpMessages 'typed' $rows2
        Shot (Join-Path $ev 'typed-reply.png')
        L "GPU_AFTER_TYPED $(Gpu)"
        L "SERVERS_FINAL_COUNT=$((ServerInfo).Count)"

        $chats = Join-Path $env:APPDATA 'HOWLChat\chats'
        $newest = Get-ChildItem $chats -Directory | Sort-Object LastWriteTime -Descending | Select-Object -First 1
        L "CHAT_DIR=$($newest.FullName)"
        Set-Content -Path (Join-Path $ev 'chat-id.txt') -Value $newest.Name -Encoding UTF8
        Copy-Item $newest.FullName (Join-Path $ev 'chat-copy') -Recurse -Force
        foreach ($f in @('meta.json','history.ndjson','history.json')) {
            $fp = Join-Path $newest.FullName $f
            if (Test-Path $fp) { L "HASH $f $((Get-FileHash $fp -Algorithm SHA256).Hash)" }
        }
        $cfg = Join-Path $env:APPDATA 'HOWLChat\config.json'
        if (Test-Path $cfg) { L "CONFIG $((Get-Content $cfg -Raw).Trim())"; L "HASH config.json $((Get-FileHash $cfg -Algorithm SHA256).Hash)" }
    }
    else {
        # precheck: dictation only, no model load (ASR is independent of the picker)
        $t = Dictate 'precheck' $RenderMatch
        Shot (Join-Path $ev 'precheck.png')
        L "PRECHECK_RESULT_LENGTH=$($t.Length)"
        L "GPU_PRECHECK $(Gpu)"
    }

    CloseApp
    L 'PHASE_COMPLETE'
}
catch {
    L "EXCEPTION $($_.Exception.Message)"
    L "STACK $($_.ScriptStackTrace)"
    CloseApp
    exit 1
}
