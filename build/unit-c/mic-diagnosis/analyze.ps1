param([string]$Path)
$b=[IO.File]::ReadAllBytes($Path)
if ($b.Length -lt 44 -or [Text.Encoding]::ASCII.GetString($b,0,4) -ne 'RIFF') { 'NOT_WAV'; exit 1 }
$rate=[BitConverter]::ToInt32($b,24); $ch=[BitConverter]::ToInt16($b,22); $bits=[BitConverter]::ToInt16($b,34)
$dataLen=[BitConverter]::ToInt32($b,40); $n=[int]($dataLen/2)
$peak=0.0; $sum=0.0; $nz=0; $over=0; $clip=0
for($i=44;$i -lt 44+$dataLen;$i+=2){ $s=[BitConverter]::ToInt16($b,$i)/32768.0; $a=[Math]::Abs($s); if($a -gt $peak){$peak=$a}; $sum+=$s*$s; if($s -ne 0){$nz++}; if($a -gt 0.02){$over++}; if($a -ge 0.999){$clip++} }
$rms=[Math]::Sqrt($sum/[Math]::Max(1,$n))
'{0} rate={1} ch={2} bits={3} samples={4} dur={5:N2}s peak={6:N4} rms={7:N5} nonzero={8} over2pct={9} clip={10} bytes={11}' -f (Split-Path $Path -Leaf),$rate,$ch,$bits,$n,($n/[double]$rate),$peak,$rms,$nz,$over,$clip,$b.Length
