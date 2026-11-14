param()
$ErrorActionPreference = 'Continue'
$root = 'C:\Users\StoneWolfPC\HOWL_Chat'
$iso  = Join-Path $root 'build\tmp\unit-c-rerun'
$ev   = Join-Path $root 'build\unit-c\coexistence'
$uia  = Join-Path $root 'build\unit-c\tools\uia.ps1'
$shot = Join-Path $root 'build\unit-c\tools\shot.ps1'
$log  = Join-Path $ev 'audio-picker-check.log'
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
function Types {
    $reg = Join-Path $iso 'config\model-registry.json'
    if (-not (Test-Path $reg)) { return 'no-registry' }
    return (python -c "import json,collections;d=json.load(open(r'$reg',encoding='utf-8'));c=collections.Counter(str(e.get('type')) for e in d);t=[e['name'] for e in d if e.get('type')=='tts'];print(dict(c));print('TTS='+','.join(t))" 2>&1 | Out-String).Trim()
}

# --- seed a tts folder with one voice file and one chat file -----------------
Copy-Item $cfg "$cfg.audiotest.bak" -Force
New-Item -ItemType Directory -Force -Path (Join-Path $iso 'demo\tts') | Out-Null
Set-Content -Path (Join-Path $iso 'demo\tts\TestVoiceA.onnx') -Value 'dummy onnx voice model' -Encoding ASCII
Set-Content -Path (Join-Path $iso 'demo\tts\ShouldNotAppear.gguf') -Value 'GGUF' -Encoding ASCII
python -c "import json,sys;p=r'$cfg';d=json.load(open(p,encoding='utf-8'));d[r'$(Join-Path $iso 'demo\tts')']='tts';json.dump(d,open(p,'w',encoding='utf-8'),indent=2)"
L "SEEDED_DIRS=$((Get-Content $cfg -Raw) -replace '\s+',' ')"

$p = Start-Process -FilePath (Join-Path $iso 'howl-chat-unit-c.exe') -WorkingDirectory $iso -PassThru
$sw = [Diagnostics.Stopwatch]::StartNew(); $q = $null
while ($sw.Elapsed.TotalSeconds -lt 90) {
    $q = Get-Process -Name 'howl-chat-unit-c' -ErrorAction SilentlyContinue | Where-Object { $_.MainWindowHandle -ne 0 } | Select-Object -First 1
    if ($null -ne $q) { break }
    Start-Sleep -Milliseconds 400
}
if ($null -eq $q) { L 'WINDOW_TIMEOUT'; exit 1 }
L "WINDOW pid=$($q.Id) after=$([int]$sw.Elapsed.TotalSeconds)s"
Start-Sleep -Seconds 7
L "REGISTRY_AFTER_START $((Types) -replace "`n",' | ')"

# --- UI rescan path ----------------------------------------------------------
L "OPEN_BROWSER $(Uia @('-Action','invoke','-Name','^Browser$'))"
Start-Sleep -Seconds 3
$bt = Uia @('-Action','tree','-OutFile',(Join-Path $ev 'audio-browser-tree.txt'))
L "BROWSER_OPEN=$($bt -match 'model-browser-modal')"
L "RESCAN_BTN=$($bt -match 'mb-dirs-rescan')"
L "RESCAN_CLICK $(Uia @('-Action','invoke','-Class','mb-dirs-rescan'))"
Start-Sleep -Seconds 6
$rt = Uia @('-Action','tree')
L "RESCAN_TOAST=$($rt -match 'Model directories rescanned')"
L "RESCAN_NO_ERROR_TOAST=$(-not ($rt -match 'Rescan failed'))"
L "REGISTRY_AFTER_UI_RESCAN $((Types) -replace "`n",' | ')"
L "CLOSE_BROWSER $(Uia @('-Action','invoke','-Id','mb-close'))"
Start-Sleep -Seconds 1

# --- audio settings -> TTS voice picker --------------------------------------
L "OPEN_SETTINGS_MENU $(Uia @('-Action','invoke','-Name','^Settings$','-Class','menu-button'))"
Start-Sleep -Milliseconds 900
L "OPEN_AUDIO_SETTINGS $(Uia @('-Action','invoke','-Name','^Audio Settings$'))"
Start-Sleep -Seconds 4
$at = Uia @('-Action','tree','-OutFile',(Join-Path $ev 'audio-tts-tree.txt'))
L "AUDIO_IFRAME_PRESENT=$($at -match 'TTS Voice Model|General|Audio')"
L "HAS_TESTVOICE=$($at -match 'TestVoiceA')"
L "HAS_SHOULDNOTAPPEAR=$($at -match 'ShouldNotAppear')"
foreach ($m in [regex]::Matches($at, 'name=([^|]*TestVoiceA[^|]*)')) { L "NODE $($m.Groups[1].Value.Trim())" }
& powershell -NoProfile -ExecutionPolicy Bypass -File $shot -OutFile (Join-Path $ev 'audio-tts-picker-app.png') | Out-Null
L 'SHOT audio-tts-picker-app.png'

# --- restore ----------------------------------------------------------------
$app = Get-Process -Name 'howl-chat-unit-c' -ErrorAction SilentlyContinue | Select-Object -First 1
if ($app) { [void]$app.CloseMainWindow(); [void]$app.WaitForExit(30000) }
Start-Sleep -Seconds 2
Copy-Item "$cfg.audiotest.bak" $cfg -Force
Remove-Item "$cfg.audiotest.bak" -Force -ErrorAction SilentlyContinue
Remove-Item (Join-Path $iso 'demo') -Recurse -Force -ErrorAction SilentlyContinue
L "RESTORED_DIRS=$((Get-Content $cfg -Raw) -replace '\s+',' ')"
L "LEFTOVER=$(@(Get-Process -Name 'howl-chat-unit-c','llama-server' -ErrorAction SilentlyContinue).Count)"
