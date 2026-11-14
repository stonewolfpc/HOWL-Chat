param()
$ErrorActionPreference = 'Continue'

$root = 'C:\Users\StoneWolfPC\HOWL_Chat'
$iso  = Join-Path $root 'build\tmp\unit-c-rerun'
$ev   = Join-Path $root 'build\unit-c\coexistence'
$uia  = Join-Path $root 'build\unit-c\tools\uia.ps1'
$shot = Join-Path $root 'build\unit-c\tools\shot.ps1'
$log  = Join-Path $ev 'discovery-ui2.log'
Set-Location $root
New-Item -ItemType Directory -Force -Path $ev | Out-Null

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
function WaitWindow([int]$seconds = 60) {
    $sw = [Diagnostics.Stopwatch]::StartNew()
    while ($sw.Elapsed.TotalSeconds -lt $seconds) {
        $p = Get-Process -Name 'howl-chat-unit-c' -ErrorAction SilentlyContinue | Where-Object { $_.MainWindowHandle -ne 0 } | Select-Object -First 1
        if ($null -ne $p) { return $p }
        Start-Sleep -Milliseconds 400
    }
    return $null
}
function RegistryDump([string]$tag) {
    $reg = Join-Path $iso 'config\model-registry.json'
    $dirs = Join-Path $iso 'config\model-scan-dirs.json'
    if (Test-Path $dirs) {
        $d = Get-Content $dirs -Raw
        $hit = $d -match 'demo'
        L "DIRS[$tag] has_demo=$hit"
        foreach ($m in [regex]::Matches($d, '"([^"]*demo[^"]*)": "(\w+)"')) { L "  DIR $($m.Groups[1].Value) -> $($m.Groups[2].Value)" }
    }
    if (Test-Path $reg) {
        $r = Get-Content $reg -Raw
        L "REG[$tag] demo_entries=$(([regex]::Matches($r, 'unit-c-rerun\\\\demo')).Count)"
        foreach ($m in [regex]::Matches($r, '"name": "([^"]*)",\s*"path": "([^"]*demo[^"]*)"')) { L "  ENTRY $($m.Groups[1].Value) -> $($m.Groups[2].Value)" }
    }
}
function CloseApp {
    $p = Get-Process -Name 'howl-chat-unit-c' -ErrorAction SilentlyContinue | Select-Object -First 1
    if ($null -eq $p) { L 'CLOSE no app'; return }
    [void]$p.CloseMainWindow()
    if (-not $p.WaitForExit(30000)) { L 'CLOSE_TIMEOUT' } else { L "EXITED code=$($p.ExitCode)" }
    Start-Sleep -Seconds 2
    $left = @(Get-Process -Name 'howl-chat-unit-c', 'llama-server' -ErrorAction SilentlyContinue)
    L "LEFTOVER=$($left.Count)"
}

# --- pass A: scan the offered folder, then open the locator -------------------
Remove-Item (Join-Path $iso 'config\model-discovery.json') -ErrorAction SilentlyContinue
L 'PASS_A_START'
$p = Start-Process -FilePath (Join-Path $iso 'howl-chat-unit-c.exe') -WorkingDirectory $iso -PassThru
L "APP_PID=$($p.Id)"
if (-not (WaitWindow 90)) { L 'WINDOW_TIMEOUT'; exit 1 }
Start-Sleep -Seconds 6
Uia @('-Action', 'tree', '-OutFile', (Join-Path $ev 'ui2-prompt-tree.txt')) | Out-Null
$tree = Get-Content (Join-Path $ev 'ui2-prompt-tree.txt') -Raw
L "PROMPT_VISIBLE=$($tree -match 'New model folder detected')"
L "SCAN_CLICK $(Uia @('-Action', 'invoke', '-Name', '^Scan directory$'))"
Start-Sleep -Seconds 4
L "PROMPT_AFTER_SCAN=$((Uia @('-Action','tree')) -match 'New model folder detected')"
RegistryDump 'after-scan'

L "OPEN_BROWSER $(Uia @('-Action', 'invoke', '-Name', '^Browser$'))"
Start-Sleep -Seconds 3
L "OPEN_DIRS_DROPDOWN $(Uia @('-Action', 'invoke', '-Name', '^All directories$'))"
Start-Sleep -Milliseconds 1200
L "OPEN_LOCATOR $(Uia @('-Action', 'invoke', '-Name', 'Browse for Directory'))"
Start-Sleep -Seconds 2
Uia @('-Action', 'tree', '-OutFile', (Join-Path $ev 'ui2-locator-tree.txt')) | Out-Null
$loc = Get-Content (Join-Path $ev 'ui2-locator-tree.txt') -Raw
L "LOCATOR_VISIBLE=$($loc -match 'Choose a model folder')"
L "LOCATOR_HAS_DRIVE=$($loc -match 'G:\\\\')"
L "LOCATOR_HAS_UP=$($loc -match 'Parent folder')"
& powershell -NoProfile -ExecutionPolicy Bypass -File $shot -OutFile (Join-Path $ev 'directory-locator-app.png') | Out-Null
L 'LOCATOR_SHOT directory-locator-app.png'
L "LOCATOR_CANCEL $(Uia @('-Action', 'invoke', '-Name', '^Cancel$'))"
Start-Sleep -Milliseconds 800
L "LOCATOR_CLOSED=$(-not ((Uia @('-Action','tree')) -match 'Choose a model folder'))"
CloseApp

# --- pass B: delete the scanned file and rescan ------------------------------
Remove-Item (Join-Path $iso 'demo\New Folder\voice.onnx') -ErrorAction SilentlyContinue
L 'PASS_B_START'
$p2 = Start-Process -FilePath (Join-Path $iso 'howl-chat-unit-c.exe') -WorkingDirectory $iso -PassThru
L "APP_PID=$($p2.Id)"
if (-not (WaitWindow 90)) { L 'WINDOW_TIMEOUT'; exit 1 }
Start-Sleep -Seconds 5
L "OPEN_BROWSER $(Uia @('-Action', 'invoke', '-Name', '^Browser$'))"
Start-Sleep -Seconds 3
L "RESCAN $(Uia @('-Action', 'invoke', '-Name', '^Rescan$'))"
Start-Sleep -Seconds 3
RegistryDump 'after-delete-rescan'
CloseApp
L 'PHASE_DONE'
