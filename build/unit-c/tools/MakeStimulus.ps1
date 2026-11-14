param(
    [string]$OutFile = 'build/unit-c/stimulus.wav',
    [string]$Text = 'Hello. This is a test of the local speech recognition system.'
)
$ErrorActionPreference = 'Stop'
Add-Type -AssemblyName System.Speech

$dir = Split-Path -Parent $OutFile
if ($dir -and -not (Test-Path $dir)) { New-Item -ItemType Directory -Path $dir | Out-Null }
$full = [System.IO.Path]::GetFullPath((Join-Path (Get-Location) $OutFile))

$synth = New-Object System.Speech.Synthesis.SpeechSynthesizer
$synth.Rate = -1
$synth.Volume = 100
$voices = $synth.GetInstalledVoices() | Where-Object { $_.Enabled }
$chosen = ($voices | Select-Object -First 1)
if ($chosen) {
    $synth.SelectVoice($chosen.VoiceInfo.Name)
    Write-Output ("VOICE=" + $chosen.VoiceInfo.Name + " culture=" + $chosen.VoiceInfo.Culture.Name)
}
$fmt = New-Object System.Speech.AudioFormat.SpeechAudioFormatInfo(
    16000,
    [System.Speech.AudioFormat.AudioBitsPerSample]::Sixteen,
    [System.Speech.AudioFormat.AudioChannel]::Mono)
$synth.SetOutputToWaveFile($full, $fmt)
$synth.Speak($Text)
$synth.SetOutputToNull()
$synth.Dispose()
Write-Output "TEXT=$Text"
Write-Output "WROTE=$full"

# Report duration and peak so the stimulus level is a measured fact, not a guess.
$bytes = [System.IO.File]::ReadAllBytes($full)
$sampleRate = [BitConverter]::ToInt32($bytes, 24)
$dataBytes = $bytes.Length - 44
$seconds = $dataBytes / 2.0 / $sampleRate
$peak = 0
for ($i = 44; $i + 1 -lt $bytes.Length; $i += 2) {
    $s = [BitConverter]::ToInt16($bytes, $i)
    $a = [math]::Abs([int]$s)
    if ($a -gt $peak) { $peak = $a }
}
Write-Output ("FORMAT={0}Hz mono 16-bit bytes={1}" -f $sampleRate, $bytes.Length)
Write-Output ("DURATION={0:N2}s PEAK={1} ({2:N1} dBFS)" -f $seconds, $peak, (20 * [math]::Log10([math]::Max($peak, 1) / 32768.0)))
