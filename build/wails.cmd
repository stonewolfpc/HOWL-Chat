@echo off
setlocal EnableExtensions
rem ===========================================================================
rem HOWL Chat - project-local Wails CLI wrapper
rem ===========================================================================
rem
rem Why this exists
rem   The Wails CLI is installed inside the project (build\tools\wails.exe) and
rem   must be driven with the Go toolchain the module pins. The default `go` on
rem   this machine is newer than the module's, and building with it fails with
rem   "internal error: package ... without types was imported from ...".
rem   Nothing global is changed: GOTOOLCHAIN is set for this invocation only.
rem
rem Commands (any other arguments are forwarded unchanged to the CLI)
rem   build\wails.cmd build             native build -> build\bin\howl-chat.exe
rem   build\wails.cmd dev               live-reload development session
rem   build\wails.cmd doctor            environment report
rem   build\wails.cmd generate module   regenerate frontend\wailsjs bindings
rem
rem Behaviour
rem   * Works from any working directory; it switches to the project root.
rem   * Exits with the Wails CLI exit code, so failures are never swallowed.
rem   * The toolchain is taken from go.mod (toolchain directive, else the go
rem     directive) and verified before the CLI runs. If that is not the version
rem     `go` actually resolves to, the script aborts instead of building with
rem     whatever Go happens to be installed.
rem ===========================================================================

set "WAILS_CLI=%~dp0tools\wails.exe"

rem --- project root ---------------------------------------------------------
pushd "%~dp0.." >nul 2>&1
if errorlevel 1 (
    echo [ERROR] Cannot enter the project root next to "%~f0".
    exit /b 1
)

if not exist "wails.json" goto :no_project
if not exist "go.mod" goto :no_project

rem --- CLI preflight --------------------------------------------------------
if not exist "%WAILS_CLI%" (
    echo [ERROR] Project-local Wails CLI not found: %WAILS_CLI%
    echo         Install it inside the project, for example:
    echo             set "GOBIN=%CD%\build\tools"
    echo             go install github.com/wailsapp/wails/v2/cmd/wails@v2.12.0
    set "RC=1"
    goto :done
)

rem --- toolchain selection --------------------------------------------------
where go >nul 2>&1
if errorlevel 1 (
    echo [ERROR] No `go` executable on PATH; cannot select the module toolchain.
    set "RC=1"
    goto :done
)

set "GOTOOLCHAIN_WANT="
for /f "usebackq tokens=1,2" %%A in (`findstr /b /c:"toolchain go" go.mod`) do set "GOTOOLCHAIN_WANT=%%B"
if not defined GOTOOLCHAIN_WANT for /f "usebackq tokens=1,2" %%A in (`findstr /b /c:"go 1." go.mod`) do set "GOTOOLCHAIN_WANT=go%%B"
if not defined GOTOOLCHAIN_WANT (
    echo [ERROR] Could not read a Go version from go.mod.
    set "RC=1"
    goto :done
)
set "GOTOOLCHAIN=%GOTOOLCHAIN_WANT%"

rem Prove the pinned toolchain is the one this invocation would use.
set "GO_RESOLVED="
for /f "usebackq tokens=3" %%V in (`go version 2^>nul`) do set "GO_RESOLVED=%%V"
if not defined GO_RESOLVED (
    echo [ERROR] Go could not run with GOTOOLCHAIN=%GOTOOLCHAIN% as required by go.mod.
    echo         Refusing to fall back to a different toolchain.
    set "RC=1"
    goto :done
)
echo %GO_RESOLVED% | findstr /c:"%GOTOOLCHAIN%" >nul
if errorlevel 1 (
    echo [ERROR] go.mod requires %GOTOOLCHAIN% but `go` resolved to %GO_RESOLVED%.
    echo         Refusing to build with a different toolchain.
    set "RC=1"
    goto :done
)

rem --- run ------------------------------------------------------------------
"%WAILS_CLI%" %*
set "RC=%ERRORLEVEL%"
goto :done

:no_project
echo [ERROR] wails.json and go.mod were not found in "%CD%".
echo         Run this script from inside the HOWL Chat project checkout.
set "RC=1"

:done
popd
exit /b %RC%
