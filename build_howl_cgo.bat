@echo off
setlocal
rem Build the release binary from wherever this folder lives.
set "ROOT=%~dp0"
if not exist "%ROOT%go.mod" (
    echo [ERROR] go.mod not found next to this script.
    echo Start this from the HOWL Chat project folder.
    pause
    exit /b 1
)
cd /d "%ROOT%"

rem Pure Go: llama.cpp runs as an external llama-server.exe, so no cgo,
rem no MSVC/clang toolchain and no DLL bundling are required.
set CGO_ENABLED=0
set GOFLAGS=-buildvcs=false

echo Building build\bin\howl-chat.exe ...
go build -tags "desktop,production" -ldflags "-w -s -H windowsgui" -o "%ROOT%build\bin\howl-chat.exe" .
if errorlevel 1 (
    echo [ERROR] Build failed.
    pause
    exit /b 1
)
echo Build OK: build\bin\howl-chat.exe
pause
