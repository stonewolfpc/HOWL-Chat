//go:build !uiharness

package main

import (
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

// main is the entry point for the HOWL Chat application.
// Initializes the Wails application window, state, and UI components.
func main() {
	startupMark("process start")

	// Create an instance of the app structure
	app := NewApp()
	startupMark("NewApp done")

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "HOWL Chat",
		Width:  1400,
		Height: 900,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		OnStartup: app.startup,
		Bind: []interface{}{
			app,
		},
	})
	startupMark("wails.Run returned")

	if err != nil {
		println("Error:", err.Error())
	}
}
