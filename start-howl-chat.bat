@echo off
setlocal
rem Launch HOWL Chat from wherever this folder lives (drive letters change).
set "ROOT=%~dp0"
if not exist "%ROOT%wails.json" (
    echo [ERROR] wails.json not found next to this script.
    echo Start this from the HOWL Chat project folder.
    pause
    exit /b 1
)
cd /d "%ROOT%"

rem Pass "dev" to use live-reload mode instead of the built binary.
if /i "%~1"=="dev" goto :devmode

if exist "%ROOT%build\bin\howl-chat.exe" goto :launch

echo No built app yet - building it now (first run takes a minute)...
set GOFLAGS=-buildvcs=false
go build -tags "desktop,production" -ldflags "-w -s -H windowsgui" -o "%ROOT%build\bin\howl-chat.exe" .
if errorlevel 1 (
    echo [ERROR] Build failed.
    pause
    exit /b 1
)

:launch
echo Starting HOWL Chat...
start "" "%ROOT%build\bin\howl-chat.exe"
exit /b 0

:devmode
rem Live-reload mode. The Wails CLI installs into GOPATH\bin.
for /f "usebackq delims=" %%G in (`go env GOPATH`) do set "GOPATH=%%G"
if exist "%GOPATH%\bin\wails.exe" set "PATH=%GOPATH%\bin;%PATH%"
set GOFLAGS=-buildvcs=false
wails dev
if errorlevel 1 pause
