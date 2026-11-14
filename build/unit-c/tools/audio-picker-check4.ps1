param()
$ErrorActionPreference = 'Continue'
$root = 'C:\Users\StoneWolfPC\HOWL_Chat'
$iso  = Join-Path $root 'build\tmp\unit-c-rerun'
$ev   = Join-Path $root 'build\unit-c\coexistence'
$uia  = Join-Path $root 'build\unit-c\tools\uia.ps1'
$shot = Join-Path $root 'build\unit-c\tools\shot.ps1'
$log  = Join-Path $ev 'audio-picker-check4.log'
$cfg  = Join-Path $iso 'config\model-scan-dirs.json'
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

Copy-Item $cfg "$cfg.audiotest.bak" -Force
New-Item -ItemType Directory -Force -Path (Join-Path $iso 'demo\tts') | Out-Null
Set-Content -Path (Join-Path $iso 'demo\tts\TestVoiceA.onnx') -Value 'dummy onnx voice model' -Encoding ASCII
Set-Content -Path (Join-Path $iso 'demo\tts\ShouldNotAppear.gguf') -Value 'GGUF' -Encoding ASCII
python -c "import json;p=r'$cfg';d=json.load(open(p,encoding='utf-8'));d[r'$(Join-Path $iso 'demo\tts')']='tts';json.dump(d,open(p,'w',encoding='utf-8'),indent=2)"

$p = Start-Process -FilePath (Join-Path $iso 'howl-chat-unit-c.exe') -WorkingDirectory $iso -PassThru
$sw = [Diagnostics.Stopwatch]::StartNew(); $q = $null
while ($sw.Elapsed.TotalSeconds -lt 90) {
    $q = Get-Process -Name 'howl-chat-unit-c' -ErrorAction SilentlyContinue | Where-Object { $_.MainWindowHandle -ne 0 } | Select-Object -First 1
    if ($null -ne $q) { break }
    Start-Sleep -Milliseconds 400
}
if ($null -eq $q) { L 'WINDOW_TIMEOUT'; exit 1 }
L "UI_READY=$([bool](WaitTree 'menu-button' 60))"

L "OPEN_SETTINGS_MENU $(InvokeWhen 'name=Settings' @('-Action','invoke','-Name','^Settings$','-Class','menu-button'))"
Start-Sleep -Milliseconds 800
L "OPEN_AUDIO_SETTINGS $(InvokeWhen 'name=Audio Settings' @('-Action','invoke','-Name','^Audio Settings$'))"
L "AUDIO_DOC_READY=$([bool](WaitTree 'Audio Settings - HOWL Chat' 30))"
L "OPEN_TTS_TAB $(InvokeWhen 'Text-to-Speech' @('-Action','invoke','-Name','Text-to-Speech'))"
Start-Sleep -Seconds 2
$t1 = Uia @('-Action','tree','-OutFile',(Join-Path $ev 'audio4-tts-tree.txt'))
L "HAS_TTS_VOICE_FIELD=$($t1 -match 'tts-voice')"
L "HAS_TESTVOICE_COLLAPSED=$($t1 -match 'TestVoiceA')"
foreach ($m in [regex]::Matches($t1, '[^\n]*tts-voice[^\n]*')) { L "FIELD $($m.Value.Trim())" }

L "OPEN_SELECT $(InvokeWhen 'tts-voice' @('-Action','invoke','-Id','tts-voice'))"
Start-Sleep -Milliseconds 1500
$t2 = Uia @('-Action','tree','-OutFile',(Join-Path $ev 'audio4-tts-open-tree.txt'))
L "OPEN_HAS_TESTVOICE=$($t2 -match 'TestVoiceA')"
L "OPEN_HAS_BUILTIN=$($t2 -match 'en_US-lessac-medium')"
L "OPEN_HAS_CHATFILE=$($t2 -match 'ShouldNotAppear')"
foreach ($m in [regex]::Matches($t2, '[^\n]*(TestVoiceA|lessac|ShouldNotAppear)[^\n]*')) { L "OPT $($m.Value.Trim())" }
& powershell -NoProfile -ExecutionPolicy Bypass -File $shot -OutFile (Join-Path $ev 'audio4-tts-open.png') | Out-Null
L 'SHOT audio4-tts-open.png'

$app = Get-Process -Name 'howl-chat-unit-c' -ErrorAction SilentlyContinue | Select-Object -First 1
if ($app) { [void]$app.CloseMainWindow(); [void]$app.WaitForExit(30000) }
Start-Sleep -Seconds 2
Copy-Item "$cfg.audiotest.bak" $cfg -Force
Remove-Item "$cfg.audiotest.bak" -Force -ErrorAction SilentlyContinue
Remove-Item (Join-Path $iso 'demo') -Recurse -Force -ErrorAction SilentlyContinue
L "RESTORED_DIRS=$((Get-Content $cfg -Raw) -replace '\s+',' ')"
L "LEFTOVER=$(@(Get-Process -Name 'howl-chat-unit-c','llama-server' -ErrorAction SilentlyContinue).Count)"
