param()
$ErrorActionPreference = 'Continue'

$root = 'C:\Users\StoneWolfPC\HOWL_Chat'
$iso  = Join-Path $root 'build\tmp\unit-c-rerun'
$ev   = Join-Path $root 'build\unit-c\coexistence'
$uia  = Join-Path $root 'build\unit-c\tools\uia.ps1'
$shot = Join-Path $root 'build\unit-c\tools\shot.ps1'
$log  = Join-Path $ev 'discovery-ui.log'
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

function Shot([string]$name) {
    $out = (& powershell -NoProfile -ExecutionPolicy Bypass -File $shot -OutFile (Join-Path $ev $name) 2>&1 | Out-String).Trim()
    L "SHOT $name -> $out"
}

function WaitWindow([int]$seconds = 60) {
    $sw = [Diagnostics.Stopwatch]::StartNew()
    while ($sw.Elapsed.TotalSeconds -lt $seconds) {
        $p = Get-Process -Name 'howl-chat-unit-c' -ErrorAction SilentlyContinue | Where-Object { $_.MainWindowHandle -ne 0 } | Select-Object -First 1
        if ($null -ne $p) { return $p }
        Start-Sleep -Milliseconds 500
    }
    return $null
}

L 'PHASE discovery-ui'
$p = Start-Process -FilePath (Join-Path $iso 'howl-chat-unit-c.exe') -WorkingDirectory $iso -PassThru
L "APP_PID=$($p.Id)"
if (-not (WaitWindow 90)) { L 'WINDOW_TIMEOUT'; exit 1 }
L 'WINDOW_READY'
Start-Sleep -Seconds 6

# 1. Discovery prompt (seeded sibling folder under demo\).
$tree = Uia @('-Action', 'tree', '-OutFile', (Join-Path $ev 'ui-prompt-tree.txt'))
$promptSeen = (Get-Content (Join-Path $ev 'ui-prompt-tree.txt') -Raw) -match 'New model folder detected'
L "PROMPT_VISIBLE=$promptSeen"
if ($promptSeen) {
    $path = (Get-Content (Join-Path $ev 'ui-prompt-tree.txt') -Raw | Select-String -Pattern 'demo|New Folder' -AllMatches).Matches.Count
    L "PROMPT_PATH_MENTIONS=$path"
    Shot 'discovery-prompt.png'
    $btn = Uia @('-Action', 'find', '-Name', '^Not now$')
    L "PROMPT_BUTTONS $(($btn -split "`r?`n").Count)"
    L "PROMPT_NOT_NOW $(Uia @('-Action', 'invoke', '-Name', '^Not now$'))"
    Start-Sleep -Milliseconds 1200
    $after = Uia @('-Action', 'tree')
    L "PROMPT_AFTER_DISMISS_VISIBLE=$($after -match 'New model folder detected')"
    $recPath = Join-Path $iso 'config\model-discovery.json'
    if (Test-Path $recPath) { L "DEFERRED_RECORD $((Get-Content $recPath -Raw).Replace("`r",'').Replace("`n",' '))" }
    else { L 'DEFERRED_RECORD missing' }
}

# 2. Model browser: rendered chat picker contents.
L "OPEN_BROWSER $(Uia @('-Action', 'invoke', '-Name', '^Browser$'))"
Start-Sleep -Seconds 4
Uia @('-Action', 'tree', '-OutFile', (Join-Path $ev 'ui-browser-tree.txt')) | Out-Null
$browser = Get-Content (Join-Path $ev 'ui-browser-tree.txt') -Raw
L "BROWSER_HAS_GEMMA=$($browser -match 'Gemma-4-E4B')"
L "BROWSER_HAS_ASR=$($browser -match 'Qwen3-ASR')"
L "BROWSER_HAS_MMPROJ=$($browser -match 'mmproj')"
Shot 'model-browser-rendered.png'

# 3. The new folder locator.
L "OPEN_LOCATOR $(Uia @('-Action', 'invoke', '-Name', 'Browse for Directory'))"
Start-Sleep -Seconds 2
Uia @('-Action', 'tree', '-OutFile', (Join-Path $ev 'ui-locator-tree.txt')) | Out-Null
$loc = Get-Content (Join-Path $ev 'ui-locator-tree.txt') -Raw
L "LOCATOR_VISIBLE=$($loc -match 'Choose a model folder')"
L "LOCATOR_HAS_DRIVES=$($loc -match 'C:\\\\|G:\\\\')"
Shot 'directory-locator.png'
L "LOCATOR_CANCEL $(Uia @('-Action', 'invoke', '-Name', '^Cancel$'))"
Start-Sleep -Milliseconds 800
Uia @('-Action', 'tree', '-OutFile', (Join-Path $ev 'ui-after-cancel.txt')) | Out-Null
L "LOCATOR_CLOSED=$(-not ((Get-Content (Join-Path $ev 'ui-after-cancel.txt') -Raw) -match 'Choose a model folder'))"

# 4. Close normally.
$proc = Get-Process -Name 'howl-chat-unit-c' -ErrorAction SilentlyContinue | Select-Object -First 1
if ($proc) {
    [void]$proc.CloseMainWindow()
    if (-not $proc.WaitForExit(40000)) { L 'CLOSE_TIMEOUT' } else { L "EXITED code=$($proc.ExitCode)" }
}
Start-Sleep -Seconds 3
$left = @(Get-Process -Name 'howl-chat-unit-c', 'llama-server' -ErrorAction SilentlyContinue)
L "LEFTOVER=$($left.Count)"
L 'PHASE_DONE'
