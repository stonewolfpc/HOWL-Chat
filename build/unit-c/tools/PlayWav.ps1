param([string]$Path = 'build/unit-c/stimulus.wav')
$ErrorActionPreference = 'Stop'
$full = [System.IO.Path]::GetFullPath((Join-Path (Get-Location) $Path))
if (-not (Test-Path $full)) { throw "missing wav: $full" }
$player = New-Object System.Media.SoundPlayer $full
$sw = [System.Diagnostics.Stopwatch]::StartNew()
$player.PlaySync()
$sw.Stop()
Write-Output ("PLAYED {0} in {1:N2}s" -f $Path, $sw.Elapsed.TotalSeconds)
