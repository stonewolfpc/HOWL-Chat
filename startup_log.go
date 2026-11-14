package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// processStart is the reference point for every startup measurement.
var processStart = time.Now()

// startupMark records one boot phase with the time elapsed since process start.
// Startup cost has to be measured, not guessed: the gap between "NewApp done"
// and "OnStartup" is exactly how long the webview took to come up, while the
// gap from process start to "NewApp done" is backend initialisation.
func startupMark(label string) {
	nl := string([]byte{10})
	elapsed := float64(time.Since(processStart).Microseconds()) / 1000.0
	line := fmt.Sprintf("%s  +%7.1fms  %s", time.Now().Format("15:04:05.000"), elapsed, label) + nl

	fmt.Print("[startup] " + line)

	logDir := isolatedStatePath("config", "logs")
	if logDir == "" {
		logDir = filepath.Join("config", "logs")
	}
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return
	}
	f, err := os.OpenFile(filepath.Join(logDir, "startup.log"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer f.Close()
	_, _ = f.WriteString(line)
}
