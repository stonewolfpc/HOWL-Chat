package main

import "context"

// App represents the Wails application structure
// Contains the application context for Wails lifecycle management
type App struct {
	ctx context.Context
}

// NewApp creates a new instance of the App struct
// Initializes the application with minimal setup
func NewApp() *App {
	return &App{}
}

// startup is called when the application is starting up
// Initializes the application context
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// OnShutdown is called when the application is shutting down
// Performs cleanup operations before application termination
func (a *App) OnShutdown(ctx context.Context) {
	// Cleanup on shutdown
}
