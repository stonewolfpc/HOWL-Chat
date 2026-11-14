# UIA driver for the real Wails window (WebView2 content is exposed through UIA).
param(
    [Parameter(Mandatory = $true)][ValidateSet('tree','find','invoke','toggle','click','focus','value','key','text','rect')]
    [string]$Action,
    [string]$Name = '',
    [string]$Class = '',
    [string]$Id = '',
    [int]$Depth = 99,
    [int]$Index = 0,
    [int]$Max = 600,
    [string]$OutFile = '',
    [string]$Text = '',
    [int]$X = -1,
    [int]$Y = -1
)

$ErrorActionPreference = 'Stop'
Add-Type -AssemblyName UIAutomationClient
Add-Type -AssemblyName UIAutomationTypes
Add-Type -AssemblyName System.Windows.Forms

$win32 = @'
using System;
using System.Runtime.InteropServices;
public static class W32 {
    [DllImport("user32.dll")] public static extern bool SetCursorPos(int x, int y);
    [DllImport("user32.dll")] public static extern void mouse_event(uint f, uint dx, uint dy, int d, UIntPtr e);
    [DllImport("user32.dll")] public static extern bool SetForegroundWindow(IntPtr h);
    [DllImport("user32.dll")] public static extern bool ShowWindow(IntPtr h, int c);
    [DllImport("user32.dll")] public static extern bool IsIconic(IntPtr h);
    [DllImport("user32.dll")] public static extern IntPtr GetForegroundWindow();
    [DllImport("user32.dll")] public static extern bool GetWindowRect(IntPtr h, out RECT r);
    [StructLayout(LayoutKind.Sequential)] public struct RECT { public int L, T, R, B; }
    public static void ClickAt(int x, int y) {
        SetCursorPos(x, y);
        System.Threading.Thread.Sleep(60);
        mouse_event(0x0002, 0, 0, 0, UIntPtr.Zero);
        System.Threading.Thread.Sleep(40);
        mouse_event(0x0004, 0, 0, 0, UIntPtr.Zero);
    }
}
'@
Add-Type -TypeDefinition $win32

function Get-AppWindow {
    $p = Get-Process -Name 'howl-chat-unit-c' -ErrorAction SilentlyContinue | Where-Object { $_.MainWindowHandle -ne 0 } | Select-Object -First 1
    if (-not $p) { throw 'howl-chat-unit-c window not found' }
    if ([W32]::IsIconic($p.MainWindowHandle)) { [void][W32]::ShowWindow($p.MainWindowHandle, 9) }
    [void][W32]::SetForegroundWindow($p.MainWindowHandle)
    Start-Sleep -Milliseconds 250
    $el = [System.Windows.Automation.AutomationElement]::FromHandle($p.MainWindowHandle)
    return [pscustomobject]@{ Element = $el; Handle = $p.MainWindowHandle; Proc = $p }
}

function SafeInt($v) {
    if ($null -eq $v) { return -1 }
    $d = 0.0
    try { $d = [double]$v } catch { return -1 }
    if ([double]::IsNaN($d) -or [double]::IsInfinity($d)) { return -1 }
    if ($d -gt 2147483000.0 -or $d -lt -2147483000.0) { return -1 }
    return [int]$d
}

function Get-RectStr($r) {
    return ('{0},{1},{2},{3}' -f (SafeInt $r.X), (SafeInt $r.Y), (SafeInt $r.Width), (SafeInt $r.Height))
}

function Get-Center($r) {
    $x = SafeInt $r.X; $y = SafeInt $r.Y; $w = SafeInt $r.Width; $h = SafeInt $r.Height
    if ($x -lt 0 -or $y -lt 0 -or $w -le 0 -or $h -le 0) { return $null }
    return @([int]($x + $w / 2), [int]($y + $h / 2))
}

function Dump-Tree($root, [int]$maxDepth) {
    $walker = [System.Windows.Automation.TreeWalker]::ControlViewWalker
    $script:lines = New-Object System.Collections.Generic.List[string]
    function Recurse($el, $d) {
        if ($script:lines.Count -ge $Max) { return }
        $ct = ''
        try { $ct = $el.Current.ControlType.ProgrammaticName -replace 'ControlType\.', '' } catch {}
        $nm = ''
        try { $nm = [string]$el.Current.Name } catch {}
        if ($nm.Length -gt 90) { $nm = $nm.Substring(0, 90) + '~' }
        $cl = ''
        try { $cl = [string]$el.Current.ClassName } catch {}
        $ai = ''
        try { $ai = [string]$el.Current.AutomationId } catch {}
        $off = ''
        try { if ($el.Current.IsOffscreen) { $off = ' OFFSCR' } } catch {}
        $script:lines.Add(('{0}|{1}|cls={2}|id={3}|name={4}|rect={5}{6}' -f
            ('  ' * $d), $ct, $cl, $ai, $nm, (Get-RectStr $el.Current.BoundingRectangle), $off))
        if ($d -ge $maxDepth) { return }
        try { $child = $walker.GetFirstChild($el) } catch { return }
        while ($child -ne $null) {
            Recurse $child ($d + 1)
            try { $child = $walker.GetNextSibling($child) } catch { break }
        }
    }
    Recurse $root 0
    return $script:lines
}

function Find-Elements($root) {
    $cond = [System.Windows.Automation.Condition]::TrueCondition
    $all = $root.FindAll([System.Windows.Automation.TreeScope]::Descendants, $cond)
    $res = New-Object System.Collections.Generic.List[object]
    foreach ($e in $all) {
        if ($Max -gt 0 -and $res.Count -ge $Max) { break }
        $nm = ''; $cl = ''; $ai = ''
        try { $nm = [string]$e.Current.Name } catch { continue }
        try { $cl = [string]$e.Current.ClassName } catch {}
        try { $ai = [string]$e.Current.AutomationId } catch {}
        $ok = $true
        if ($Name -ne '' -and $nm -notmatch $Name) { $ok = $false }
        if ($Class -ne '' -and $cl -notmatch $Class) { $ok = $false }
        if ($Id -ne '' -and $ai -ne $Id) { $ok = $false }
        if ($ok) { $res.Add($e) }
    }
    return , $res
}

$app = Get-AppWindow
$root = $app.Element

switch ($Action) {
    'tree' {
        $lines = Dump-Tree $root $Depth
        $out = $lines -join "`n"
        if ($OutFile -ne '') { $out | Set-Content -Encoding UTF8 $OutFile }
        Write-Output $out
        Write-Output ("---ELEMENTS={0}---" -f $lines.Count)
    }
    'find' {
        $found = Find-Elements $root
        $i = 0
        foreach ($e in $found) {
            $ct = ''; try { $ct = $e.Current.ControlType.ProgrammaticName -replace 'ControlType\.', '' } catch {}
            Write-Output ('[{0}] {1} cls="{2}" id="{3}" name="{4}" rect={5} enabled={6} offscreen={7}' -f `
                $i, $ct, $e.Current.ClassName, $e.Current.AutomationId, $e.Current.Name, `
                (Get-RectStr $e.Current.BoundingRectangle), $e.Current.IsEnabled, $e.Current.IsOffscreen)
            $i++
        }
        Write-Output ("---MATCHES={0}---" -f $found.Count)
    }
    'invoke' {
        $found = Find-Elements $root
        if ($found.Count -le $Index) { throw ("index $Index out of range; matches=" + $found.Count) }
        $e = $found[$Index]
        $pat = $null
        if ($e.TryGetCurrentPattern([System.Windows.Automation.InvokePattern]::Pattern, [ref]$pat)) {
            $pat.Invoke(); Write-Output "INVOKED invoke-pattern index=$Index"
        } else {
            $c = Get-Center $e.Current.BoundingRectangle
            if ($null -eq $c) { throw 'element has no usable bounds and no InvokePattern' }
            [W32]::ClickAt($c[0], $c[1])
            Write-Output "INVOKED mouse-click index=$Index at $($c[0]),$($c[1])"
        }
    }
    'toggle' {
        $found = Find-Elements $root
        if ($found.Count -le $Index) { throw ("index $Index out of range; matches=" + $found.Count) }
        $e = $found[$Index]
        $pat = $null
        if ($e.TryGetCurrentPattern([System.Windows.Automation.TogglePattern]::Pattern, [ref]$pat)) {
            $before = $pat.Current.ToggleState
            $pat.Toggle()
            Start-Sleep -Milliseconds 200
            Write-Output "TOGGLED before=$before after=$($pat.Current.ToggleState)"
        } else {
            $c = Get-Center $e.Current.BoundingRectangle
            if ($null -eq $c) { throw 'element has no usable bounds and no TogglePattern' }
            [W32]::ClickAt($c[0], $c[1])
            Write-Output "TOGGLED via click (no TogglePattern) at $($c[0]),$($c[1])"
        }
    }
    'click' {
        if ($X -ge 0 -and $Y -ge 0) {
            [W32]::ClickAt($X, $Y); Write-Output "CLICKED at $X,$Y"
        } else {
            $found = Find-Elements $root
            if ($found.Count -le $Index) { throw ("index $Index out of range; matches=" + $found.Count) }
            $c = Get-Center $found[$Index].Current.BoundingRectangle
            if ($null -eq $c) { throw 'element has no usable bounds' }
            [W32]::ClickAt($c[0], $c[1])
            Write-Output "CLICKED index=$Index at $($c[0]),$($c[1])"
        }
    }
    'focus' {
        $found = Find-Elements $root
        if ($found.Count -le $Index) { throw ("index $Index out of range; matches=" + $found.Count) }
        $found[$Index].SetFocus()
        Write-Output "FOCUSED index=$Index"
    }
    'value' {
        $found = Find-Elements $root
        if ($found.Count -le $Index) { throw ("index $Index out of range; matches=" + $found.Count) }
        $e = $found[$Index]
        $pat = $null
        if ($e.TryGetCurrentPattern([System.Windows.Automation.ValuePattern]::Pattern, [ref]$pat)) {
            Write-Output ("VALUE=" + $pat.Current.Value)
        } else { Write-Output 'NOVALUEPATTERN' }
    }
    'text' {
        $found = Find-Elements $root
        if ($found.Count -le $Index) { throw ("index $Index out of range; matches=" + $found.Count) }
        $e = $found[$Index]
        $e.SetFocus()
        Start-Sleep -Milliseconds 150
        $pat = $null
        if ($e.TryGetCurrentPattern([System.Windows.Automation.ValuePattern]::Pattern, [ref]$pat)) {
            $pat.SetValue($Text)
            Write-Output "SETVALUE ok"
        } else {
            [System.Windows.Forms.SendKeys]::SendWait($Text)
            Write-Output "SENDKEYS ok"
        }
    }
    'key' {
        [System.Windows.Forms.SendKeys]::SendWait($Text)
        Write-Output "KEY sent: $Text"
    }
    'rect' {
        $r = New-Object W32+RECT
        [void][W32]::GetWindowRect($app.Handle, [ref]$r)
        Write-Output ("WINDOW rect={0},{1},{2},{3} handle={4} pid={5} title={6}" -f $r.L, $r.T, ($r.R - $r.L), ($r.B - $r.T), $app.Handle, $app.Proc.Id, $app.Proc.MainWindowTitle)
    }
}
