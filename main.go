//go:build !uiharness

package main

import (
	"embed"
	"fmt"
	"os"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/windows"

	"howl-chat/internal/paths"
)

//go:embed all:frontend/dist
var assets embed.FS

// main is the entry point for the HOWL Chat application.
// Initializes the Wails application window, state, and UI components.
func main() {
	startupMark("process start")

	// Isolation is opt-in (HOWL_WEBVIEW_DATA_DIR) and validated before anything
	// can read or write application state. A configured-but-unusable root aborts
	// here: falling back to the real profile would defeat the point of the switch.
	isolation, isoErr := paths.ResolveIsolation()
	if isoErr != nil {
		fmt.Fprintf(os.Stderr, "[isolation] %s is set but unusable: %v\n", paths.IsolationEnvVar, isoErr)
		fmt.Fprintf(os.Stderr, "[isolation] refusing to start. Unset %s to run normally.\n", paths.IsolationEnvVar)
		os.Exit(2)
	}

	// Create an instance of the app structure
	app := NewApp()
	startupMark("NewApp done")

	// Create application with options
	appOptions := &options.App{
		Title:  "HOWL Chat",
		Width:  1400,
		Height: 900,
		AssetServer: &assetserver.Options{
			Assets: assets,
			// Pages navigated to directly (lorebooks.html, characters.html, ...) are
			// not the AssetServer's "main page", so they would otherwise run without
			// window.go. One seam gives every page the same backend bridge.
			Middleware: bridgeMiddleware,
		},
		OnStartup:  app.startup,
		OnShutdown: app.OnShutdown,
		Bind: []interface{}{
			app,
		},
	}

	if isolation != nil {
		// WebView2 is handed its own user-data folder so the browser profile and
		// the application's state root are both inside the isolation root.
		fmt.Printf("[isolation] root=%s webview=%s state=%s\n",
			isolation.Root, isolation.WebViewDataDir, isolation.StateRoot)
		appOptions.Windows = &windows.Options{WebviewUserDataPath: isolation.WebViewDataDir}
	}

	err := wails.Run(appOptions)
	startupMark("wails.Run returned")

	if err != nil {
		println("Error:", err.Error())
	}
}
