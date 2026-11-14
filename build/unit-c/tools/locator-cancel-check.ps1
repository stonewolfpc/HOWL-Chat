param()
$ErrorActionPreference = 'Continue'
$root = 'C:\Users\StoneWolfPC\HOWL_Chat'
$iso  = Join-Path $root 'build\tmp\unit-c-rerun'
$ev   = Join-Path $root 'build\unit-c\coexistence'
$uia  = Join-Path $root 'build\unit-c\tools\uia.ps1'
$log  = Join-Path $ev 'locator-cancel-check.log'
Set-Location $root

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

L 'CHECK_START'
$p = Start-Process -FilePath (Join-Path $iso 'howl-chat-unit-c.exe') -WorkingDirectory $iso -PassThru
$sw = [Diagnostics.Stopwatch]::StartNew()
$q = $null
while ($sw.Elapsed.TotalSeconds -lt 90) {
    $q = Get-Process -Name 'howl-chat-unit-c' -ErrorAction SilentlyContinue | Where-Object { $_.MainWindowHandle -ne 0 } | Select-Object -First 1
    if ($null -ne $q) { break }
    Start-Sleep -Milliseconds 400
}
if ($null -eq $q) { L 'WINDOW_TIMEOUT'; exit 1 }
L "WINDOW pid=$($q.Id) title=[$($q.MainWindowTitle)] after=$([int]$sw.Elapsed.TotalSeconds)s"
Start-Sleep -Seconds 6

$tree = Uia @('-Action','tree')
L "TREE_LEN=$($tree.Length)"
if ($tree.Length -lt 200) { L "TREE_DUMP >>> $tree"; exit 1 }

L "OPEN_BROWSER $(Uia @('-Action','invoke','-Name','^Browser$'))"
Start-Sleep -Seconds 3
L "OPEN_DIRS $(Uia @('-Action','invoke','-Name','^All directories$'))"
Start-Sleep -Milliseconds 1200
L "OPEN_LOCATOR $(Uia @('-Action','invoke','-Name','Browse for Directory'))"
Start-Sleep -Seconds 2
$before = Uia @('-Action','tree','-OutFile',(Join-Path $ev 'cancel-before.txt'))
L "BEFORE_LEN=$($before.Length)"
L "BEFORE_SHOWS_LOCATOR=$($before -match 'md-locator-overlay')"
L "BEFORE_HAS_MD_CANCEL=$($before -match 'id=md-cancel')"

L "CANCEL_BY_ID $(Uia @('-Action','invoke','-Id','md-cancel'))"
Start-Sleep -Milliseconds 1500
$after = Uia @('-Action','tree','-OutFile',(Join-Path $ev 'cancel-after.txt'))
L "AFTER_LEN=$($after.Length)"
L "AFTER_SHOWS_LOCATOR=$($after -match 'md-locator-overlay')"
L "AFTER_SHOWS_BROWSER=$($after -match 'model-browser-modal')"

$app = Get-Process -Name 'howl-chat-unit-c' -ErrorAction SilentlyContinue | Select-Object -First 1
if ($app) { [void]$app.CloseMainWindow(); [void]$app.WaitForExit(30000) }
Start-Sleep -Seconds 2
L "LEFTOVER=$(@(Get-Process -Name 'howl-chat-unit-c','llama-server' -ErrorAction SilentlyContinue).Count)"
