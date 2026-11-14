param([string]$Name,[string]$Model,[string]$Mmproj='')
$exe='C:\Users\StoneWolfPC\HOWL_Chat\runtime\llama\llama-server.exe'
$out="C:\Users\StoneWolfPC\HOWL_Chat\build\unit-c\model-audit\probes\$Name.out"
$err="C:\Users\StoneWolfPC\HOWL_Chat\build\unit-c\model-audit\probes\$Name.err"
$a=@('-m',('"'+$Model+'"'),'--host','127.0.0.1','--port','18099','-c','512','-ngl','0','-t','6','--no-warmup')
if($Mmproj -ne ''){ $a += @('--mmproj',('"'+$Mmproj+'"')) }
$sw=[Diagnostics.Stopwatch]::StartNew()
$p=Start-Process -FilePath $exe -ArgumentList $a -PassThru -RedirectStandardOutput $out -RedirectStandardError $err -WindowStyle Hidden
$ok=$false
while($sw.Elapsed.TotalSeconds -lt 90){
  Start-Sleep -Milliseconds 700
  if($p.HasExited){ break }
  try { $r=Invoke-RestMethod -Uri 'http://127.0.0.1:18099/health' -TimeoutSec 2; if($r.status -eq 'ok'){ $ok=$true; break } } catch {}
}
$sw.Stop()
"$Name HEALTH=$ok SECONDS=$([math]::Round($sw.Elapsed.TotalSeconds,1)) EXITED=$($p.HasExited)"
if(-not $p.HasExited){ Stop-Process -Id $p.Id -Force; Start-Sleep -Milliseconds 500 }
"  tail: " + ((Get-Content $err -Tail 2) -join ' | ')
