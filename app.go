package main

import (
	"context"
	"howl-chat/internal/backend/chat"
	"howl-chat/internal/backend/llama"
	"howl-chat/internal/backend/model"
)

// App represents the Wails application structure
// Contains the application context and backend services
type App struct {
	ctx          context.Context
	chatService  *chat.Service
	modelManager *model.Manager
	llamaClient  llama.Client
}

// NewApp creates a new instance of the App struct
// Initializes backend services with stub client for testing
func NewApp() *App {
	// Create stub llama client for testing (will be replaced with actual implementation)
	client := llama.NewStubClient()

	// Create model manager
	manager := model.NewManager(client)

	// Create chat service
	service := chat.NewService(manager, client)

	return &App{
		chatService:  service,
		modelManager: manager,
		llamaClient:  client,
	}
}

// startup is called when the application is starting up
// Initializes the application context
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// OnShutdown is called when the application is shutting down
// Performs cleanup operations before application termination
func (a *App) OnShutdown(ctx context.Context) {
	// Cleanup backend services
	if a.chatService != nil {
		a.chatService.Close()
	}
	if a.modelManager != nil {
		a.modelManager.Close()
	}
}

// SendMessage sends a message to the chat service and returns the AI response
// This is a Wails binding that can be called from the frontend
func (a *App) SendMessage(message string) (string, error) {
	if a.chatService == nil {
		return "", nil
	}

	responseMsg, err := a.chatService.SendMessage(message)
	if err != nil {
		return "", err
	}

	return responseMsg.Content, nil
}

// LoadModel loads a model from the given path
// This is a Wails binding that can be called from the frontend
func (a *App) LoadModel(modelPath string) error {
	if a.modelManager == nil {
		return nil
	}

	_, err := a.modelManager.LoadModel(modelPath, nil)
	return err
}

// UnloadModel unloads the currently loaded model
// This is a Wails binding that can be called from the frontend
func (a *App) UnloadModel() error {
	if a.modelManager == nil {
		return nil
	}

	return a.modelManager.UnloadModel()
}

// IsModelLoaded returns true if a model is currently loaded
// This is a Wails binding that can be called from the frontend
func (a *App) IsModelLoaded() bool {
	if a.modelManager == nil {
		return false
	}

	return a.modelManager.IsModelLoaded()
}

// GetLoadingProgress returns the current model loading progress (0.0-1.0)
// This is a Wails binding that can be called from the frontend
func (a *App) GetLoadingProgress() float64 {
	if a.modelManager == nil {
		return 0.0
	}

	return a.modelManager.GetLoadingProgress()
}

// GetLoadingStage returns the current loading stage
// This is a Wails binding that can be called from the frontend
func (a *App) GetLoadingStage() string {
	if a.modelManager == nil {
		return ""
	}

	stage := a.modelManager.GetLoadingStage()
	return string(stage)
}

// GetChatMessages returns all messages in the current chat context
// This is a Wails binding that can be called from the frontend
func (a *App) GetChatMessages() []map[string]interface{} {
	if a.chatService == nil {
		return []map[string]interface{}{}
	}

	messages := a.chatService.GetMessages()
	result := make([]map[string]interface{}, len(messages))

	for i, msg := range messages {
		result[i] = map[string]interface{}{
			"role":      string(msg.Role),
			"content":   msg.Content,
			"timestamp": msg.Timestamp,
		}
	}

	return result
}

// ClearChat clears the current chat context
// This is a Wails binding that can be called from the frontend
func (a *App) ClearChat() error {
	if a.chatService == nil {
		return nil
	}

	return a.chatService.ClearContext()
}
