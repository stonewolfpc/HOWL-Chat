# Defender false-positive differential test
#
# Purpose: decide, with evidence, whether howl-chat.exe trips Microsoft Defender
# as Trojan:Win32/Bearfoos.A!ml (ThreatID 2147731250) because of something in
# Howl's code, or because of the generic unsigned-Go-binary ML bucket.
#
# Background: that ThreatID is a generic !ml rule that also appears in public
# false-positive reports for other unsigned Go, Rust and PyInstaller binaries.
#
# Measured on this machine, November 2026: the quarantined 5 Nov build
# (sha256 ac6f4347e5518dcc35235e906ad62dd1292c80bdae542933946cb981f6579739,
# go1.23.8, tags desktop,production,wv2runtime.download,devtools) is still
# flagged every single time it is written to disk. Rebuilding the same source
# with that same toolchain and those same flags is NOT flagged, and appending
# 16 bytes to the known-bad file also clears it. So the verdict is pinned to
# that exact file content -- not to Go binaries in general, not to the missing
# version resource, and not to the process-control code.
#
# Consequence: a scan of a fresh build cannot reproduce the original event. To
# reproduce it you must write the known-flagged bytes to disk.
#
# This test builds candidates OUTSIDE the project's runnable directory and scans
# copies only. It never executes a candidate, never restores anything from
# quarantine, never touches Defender settings or exclusions, and never deletes
# evidence.
#
# Usage (from the repo root):
#   powershell -NoProfile -ExecutionPolicy Bypass -File build\security\defender-fp-test.ps1
#
# Exit codes: 0 = all candidates clean, 1 = at least one candidate detected,
#             2 = setup or build failure.

[CmdletBinding()]
param(
    [string]$OutDir = 'C:\Howl-Security-Analysis\fp-test'
)

$ErrorActionPreference = 'Stop'
$repo = (Resolve-Path (Join-Path $PSScriptRoot '..\..')).Path

function Get-MpCmdRun {
    $platform = Join-Path $env:ProgramData 'Microsoft\Windows Defender\Platform'
    $exe = Get-ChildItem -Path $platform -Directory -ErrorAction SilentlyContinue |
        Sort-Object Name -Descending |
        ForEach-Object { Join-Path $_.FullName 'MpCmdRun.exe' } |
        Where-Object { Test-Path $_ } |
        Select-Object -First 1
    if (-not $exe) { throw "MpCmdRun.exe not found under $platform" }
    return $exe
}

function Invoke-Scan([string]$path) {
    $mp = Get-MpCmdRun
    $out = & $mp -Scan -ScanPath $path -ReturnHR 2>&1 | Out-String
    $hr = $LASTEXITCODE
    return [pscustomobject]@{
        Path   = $path
        Hr     = $hr
        Output = $out.Trim()
    }
}

function Get-VersionInfoSummary([string]$path) {
    $vi = (Get-Item $path).VersionInfo
    $fields = @($vi.CompanyName, $vi.ProductName, $vi.FileVersion) | Where-Object { $_ }
    if ($fields.Count -eq 0) { return 'EMPTY (no PE version resource)' }
    return ($fields -join ' | ')
}

New-Item -ItemType Directory -Force -Path $OutDir | Out-Null

$status = Get-MpComputerStatus
Write-Host "Defender engine   : $($status.AMEngineVersion)"
Write-Host "Signature version : $($status.AntivirusSignatureVersion)"
Write-Host "Signatures dated  : $($status.AntivirusSignatureLastUpdated)"
Write-Host "Real-time protect : $($status.RealTimeProtectionEnabled)"
Write-Host "Output directory  : $OutDir"
Write-Host ''

# --- Candidate 1: minimal Go Windows GUI binary -----------------------------
# If this alone is detected, the trigger is the Go toolchain's output shape and
# not any Howl feature.
$minimalDir = Join-Path $OutDir 'minimal'
New-Item -ItemType Directory -Force -Path $minimalDir | Out-Null
$minimalSrc = Join-Path $minimalDir 'main.go'
@'
package main

import "fmt"

func main() {
	fmt.Println("minimal windowsgui binary")
}
'@ | Set-Content -Path $minimalSrc -Encoding ASCII

$env:GOTOOLCHAIN = 'go1.27.3'
$env:CGO_ENABLED = '0'
$env:GOOS = 'windows'
$env:GOARCH = 'amd64'

$minimalExe = Join-Path $minimalDir 'minimal-gui.exe'
Write-Host 'Building minimal Go GUI binary...'
Push-Location $minimalDir
try {
    & go build -ldflags '-H windowsgui -s -w' -o $minimalExe main.go
    if ($LASTEXITCODE -ne 0) { throw "minimal build failed ($LASTEXITCODE)" }
} finally { Pop-Location }

# --- Candidate 2: current Howl source --------------------------------------
# Built to the analysis directory on purpose: never into build\bin.
$howlExe = Join-Path $OutDir 'howl-current.exe'
Write-Host 'Building current Howl source (desktop,production)...'
Push-Location $repo
try {
    & go build -tags 'desktop,production' -ldflags '-H windowsgui' -o $howlExe .
    if ($LASTEXITCODE -ne 0) { throw "howl build failed ($LASTEXITCODE)" }
} finally { Pop-Location }

# --- Scan -------------------------------------------------------------------
# Copies only. Nothing here is executed.
$results = foreach ($exe in @($minimalExe, $howlExe)) {
    $scan = Invoke-Scan $exe
    [pscustomobject]@{
        Candidate   = (Split-Path $exe -Leaf)
        Size        = (Get-Item $exe).Length
        VersionInfo = (Get-VersionInfoSummary $exe)
        Signed      = ((Get-AuthenticodeSignature $exe).Status)
        ScanHr      = $scan.Hr
        Survived    = (Test-Path $exe)
        Verdict     = if (-not (Test-Path $exe)) { 'DETECTED (removed)' } elseif ($scan.Hr -eq 0) { 'clean' } else { "hr=$($scan.Hr)" }
    }
}

Write-Host ''
$results | Format-Table Candidate, Size, VersionInfo, Signed, ScanHr, Verdict -AutoSize | Out-String | Write-Host

# Any detection after the scan start would show up here.
$since = (Get-Date).AddMinutes(-10)
$new = Get-MpThreatDetection -ErrorAction SilentlyContinue |
    Where-Object { $_.InitialDetectionTime -gt $since }
if ($new) {
    Write-Host 'New detections recorded during this test:'
    $new | Select-Object InitialDetectionTime, ThreatID, Resources | Format-List | Out-String | Write-Host
}

$detected = $results | Where-Object { $_.Verdict -ne 'clean' }
if ($detected) {
    Write-Host 'RESULT: at least one candidate was flagged. Compare which one to isolate the trigger.'
    exit 1
}
Write-Host 'RESULT: both candidates clean under current definitions.'
Write-Host ''
Write-Host 'This means current source does not reproduce the alert. The original event was'
Write-Host 'a content-pinned verdict on one specific build, so no rebuild will match it.'
Write-Host 'Settle it by submitting the known-flagged file to Microsoft as a false positive:'
Write-Host 'https://www.microsoft.com/en-us/wdsi/filesubmission'
Write-Host 'Specimen sha256: ac6f4347e5518dcc35235e906ad62dd1292c80bdae542933946cb981f6579739'
exit 0
