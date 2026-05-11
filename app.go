package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"howl-chat/internal/backend/gguf"
	"howl-chat/internal/backend/llama"
	"howl-chat/internal/backend/lorebook"
	"howl-chat/internal/backend/types"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App represents the Wails application
// Contains the application context and backend services
type App struct {
	ctx             context.Context
	httpClient      *llama.HTTPClient
	serverProc      *exec.Cmd
	samplerSettings map[string]interface{}
	modelSettings   map[string]interface{}
	chatMessages    []map[string]interface{}
	lorebooks       []lorebook.Entry
	triggeredLore   map[string]bool
	streamCancel    context.CancelFunc
}

// findLlamaServer returns the path to llama-server.exe, checking next to exe and project root
func findLlamaServer() string {
	candidates := []string{
		filepath.Join(filepath.Dir(os.Args[0]), "llama-server.exe"),
		"llama-server.exe",
		`d:\Fantasy\llama-cpu\llama-server.exe`,
	}
	for _, c := range candidates {
		if _, err := os.Stat(c); err == nil {
			return c
		}
	}
	return ""
}

// intSetting extracts an int from a settings map with a fallback default
func intSetting(m map[string]interface{}, key string, def int) string {
	if v, ok := m[key]; ok {
		switch n := v.(type) {
		case float64:
			return fmt.Sprintf("%d", int(n))
		case int:
			return fmt.Sprintf("%d", n)
		}
	}
	return fmt.Sprintf("%d", def)
}

func strSetting(m map[string]interface{}, key string, def string) string {
	if v, ok := m[key].(string); ok && v != "" && v != "auto" {
		return v
	}
	return def
}

// spawnLlamaServerWithModel starts llama-server with a model and all applicable settings
func spawnLlamaServerWithModel(modelPath string, settings map[string]interface{}) (*exec.Cmd, error) {
	serverPath := findLlamaServer()
	if serverPath == "" {
		return nil, fmt.Errorf("llama-server.exe not found")
	}

	args := []string{
		"--host", "127.0.0.1",
		"--port", "8080",
		"-m", modelPath,
		"-c", intSetting(settings, "context_size", 4096),
		"--threads", intSetting(settings, "threads", 8),
		"-b", intSetting(settings, "batch_size", 512),
		"-ngl", intSetting(settings, "gpu_layers", 0),
	}

	// Auto-detect and add mmproj file for vision models
	mmprojPath := findMMProj(modelPath)
	if mmprojPath != "" {
		args = append(args, "--mmproj", mmprojPath)
		fmt.Printf("INFO: Auto-detected mmproj file: %s\n", mmprojPath)
	} else {
		fmt.Printf("INFO: No mmproj file found for model: %s\n", modelPath)
	}

	// Optional: rope scaling
	if rs := strSetting(settings, "rope_mode", ""); rs != "" {
		args = append(args, "--rope-scaling", rs)
	}
	if v, ok := settings["rope_factor"].(float64); ok && v > 0 && v != 1.0 {
		args = append(args, "--rope-scale", fmt.Sprintf("%g", v))
	}
	if v, ok := settings["rope_base"].(float64); ok && v > 0 {
		args = append(args, "--rope-freq-base", fmt.Sprintf("%g", v))
	}

	// Optional: flash attention
	if v, ok := settings["flash_attention"].(bool); ok && v {
		args = append(args, "--flash-attn")
	}

	// Optional: tensor split (multi-GPU)
	if ts := strSetting(settings, "tensor_split", ""); ts != "" && ts != "0" {
		args = append(args, "--tensor-split", ts)
	}

	// Optional: jinja chat template override
	if jinja := strSetting(settings, "custom_jinja_template", ""); jinja != "" {
		args = append(args, "--chat-template", jinja)
	}

	serverDir := filepath.Dir(serverPath)
	cmd := exec.Command(serverPath, args...)
	cmd.Dir = serverDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start llama-server: %w", err)
	}
	fmt.Printf("INFO: llama-server started (pid %d) with model %s\n", cmd.Process.Pid, modelPath)
	return cmd, nil
}

// findMMProj searches for an mmproj file in the same directory as the model
// Per llama.cpp docs: mmproj file name must start with "mmproj" (e.g., mmproj-F16.gguf)
func findMMProj(modelPath string) string {
	dir := filepath.Dir(modelPath)
	baseName := filepath.Base(modelPath)

	// Remove extension from model name
	modelName := strings.TrimSuffix(baseName, filepath.Ext(baseName))

	// Try common mmproj naming patterns
	// llama.cpp convention: file name must start with "mmproj"
	patterns := []string{
		filepath.Join(dir, "mmproj-F16.gguf"),
		filepath.Join(dir, "mmproj-Q4_0.gguf"),
		filepath.Join(dir, "mmproj-Q4_K_M.gguf"),
		filepath.Join(dir, "mmproj-Q8_0.gguf"),
		filepath.Join(dir, "mmproj-f16.gguf"),
		filepath.Join(dir, "mmproj-q4_0.gguf"),
		filepath.Join(dir, "mmproj-q4_k_m.gguf"),
		filepath.Join(dir, "mmproj-q8_0.gguf"),
		filepath.Join(dir, "mmproj.gguf"),
		filepath.Join(dir, modelName+"-mmproj.gguf"),
		filepath.Join(dir, modelName+".mmproj.gguf"),
	}

	for _, pattern := range patterns {
		if _, err := os.Stat(pattern); err == nil {
			return pattern
		}
	}

	// No mmproj file found - do NOT scan directory
	// This prevents picking up unrelated files like tokenizer.gguf
	return ""
}

// NewApp creates a new instance of the App struct
// Initializes backend services with llama-server HTTP client
func NewApp() *App {
	client := llama.NewHTTPClient("localhost", 8080, true)

	// Default sampler settings
	defaultSettings := map[string]interface{}{
		"temperature":                0.7,
		"top_p_enabled":              true,
		"top_p":                      0.9,
		"top_k":                      40,
		"min_p":                      0.05,
		"repeat_penalty":             1.1,
		"repeat_last_n":              64,
		"frequency_penalty_enabled":  true,
		"frequency_penalty":          0.0,
		"presence_penalty_enabled":   true,
		"presence_penalty":           0.0,
		"typical_p_enabled":          true,
		"typical_p":                  1.0,
		"mirostat_enabled":           true,
		"mirostat":                   0,
		"mirostat_tau":               5.0,
		"mirostat_eta":               0.1,
		"dynamic_temp_range_enabled": true,
		"dynamic_temp_range":         0.0,
		"dynamic_temp_exponent":      1.0,
		"dry_multiplier":             0.0,
		"dry_allowed_length":         2,
		"dry_base":                   1.0,
		"smoothing_factor":           0.0,
		"smoothing_curve":            1.0,
		"top_a_enabled":              true,
		"top_a":                      0.0,
		"epsilon_cutoff":             0.0,
		"eta_cutoff":                 0.0,
		"encoder_repeat_penalty":     1.0,
		"no_repeat_ngram_size":       0,
		"seed":                       -1,
	}

	return &App{
		httpClient:      client,
		serverProc:      nil,
		samplerSettings: defaultSettings,
		modelSettings:   make(map[string]interface{}),
		chatMessages:    []map[string]interface{}{},
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
	if a.serverProc != nil && a.serverProc.Process != nil {
		_ = a.serverProc.Process.Kill()
		_ = a.serverProc.Wait()
	}
}

// SendMessage sends a message to the chat service and returns the AI response
// This is a Wails binding that can be called from the frontend
func (a *App) SendMessage(message string) (string, error) {
	if a.httpClient == nil {
		return "", fmt.Errorf("llama client not initialized")
	}

	a.chatMessages = append(a.chatMessages, map[string]interface{}{
		"role":      "user",
		"content":   message,
		"timestamp": time.Now().Unix(),
	})

	// Build prompt from chat messages and triggered lore
	prompt := a.buildPrompt()

	// Convert sampler settings to InferenceOptions
	opts := llama.NewInferenceOptions()
	if s := a.samplerSettings; s != nil {
		if v, ok := s["temperature"].(float64); ok {
			opts.Temperature = v
		}
		if v, ok := s["top_p_enabled"].(bool); ok {
			opts.TopPEnabled = v
		}
		if v, ok := s["top_p"].(float64); ok {
			if opts.TopPEnabled {
				opts.TopP = v
			} else {
				opts.TopP = 1.0
			}
		}
		if v, ok := s["top_k"].(float64); ok {
			opts.TopK = int(v)
		}
		if v, ok := s["min_p"].(float64); ok {
			opts.MinP = v
		}
		if v, ok := s["typical_p_enabled"].(bool); ok {
			opts.TypicalPEnabled = v
		}
		if v, ok := s["typical_p"].(float64); ok {
			opts.TypicalP = v
		}
		if v, ok := s["repeat_penalty"].(float64); ok {
			opts.RepeatPenalty = v
		}
		if v, ok := s["repeat_last_n"].(float64); ok {
			opts.RepeatLastN = int(v)
		}
		if v, ok := s["frequency_penalty_enabled"].(bool); ok {
			opts.FrequencyPenaltyEnabled = v
		}
		if v, ok := s["frequency_penalty"].(float64); ok {
			opts.FrequencyPenalty = v
		}
		if v, ok := s["presence_penalty_enabled"].(bool); ok {
			opts.PresencePenaltyEnabled = v
		}
		if v, ok := s["presence_penalty"].(float64); ok {
			opts.PresencePenalty = v
		}
		if v, ok := s["mirostat"].(float64); ok {
			opts.Mirostat = int(v)
		}
		if v, ok := s["mirostat_tau"].(float64); ok {
			opts.MirostatTau = v
		}
		if v, ok := s["mirostat_eta"].(float64); ok {
			opts.MirostatETA = v
		}
		if v, ok := s["dynamic_temp_range"].(float64); ok {
			opts.DynamicTempRange = v
		}
		if v, ok := s["dynamic_temp_exponent"].(float64); ok {
			opts.DynamicTempExponent = v
		}
		if v, ok := s["dry_multiplier"].(float64); ok {
			opts.DRYMultiplier = v
		}
		if v, ok := s["dry_allowed_length"].(float64); ok {
			opts.DRYAllowedLength = int(v)
		}
		if v, ok := s["dry_base"].(float64); ok {
			opts.DRYBase = v
		}
		if v, ok := s["smoothing_factor"].(float64); ok {
			opts.SmoothingFactor = v
		}
		if v, ok := s["smoothing_curve"].(float64); ok {
			opts.SmoothingCurve = v
		}
		if v, ok := s["top_a"].(float64); ok {
			opts.TopA = v
		}
		if v, ok := s["epsilon_cutoff"].(float64); ok {
			opts.EpsilonCutoff = v
		}
		if v, ok := s["eta_cutoff"].(float64); ok {
			opts.EtaCutoff = v
		}
		if v, ok := s["no_repeat_ngram_size"].(float64); ok {
			opts.NoRepeatNGramSize = int(v)
		}
		if v, ok := s["seed"].(float64); ok {
			opts.Seed = int(v)
		}
	}

	// Generate response
	response, err := a.httpClient.Generate(prompt, opts)
	if err != nil {
		return "", err
	}

	a.chatMessages = append(a.chatMessages, map[string]interface{}{
		"role":      "assistant",
		"content":   response,
		"timestamp": time.Now().Unix(),
	})

	return response, nil
}

// SendMessageStream sends a message and streams the response via Wails events
// Emits: chat:start, chat:chunk (for each token), chat:complete, chat:error
func (a *App) SendMessageStream(message string) error {
	return a.SendMessageStreamWithImage(message, "")
}

// SendMessageStreamWithImage sends a message with optional image and streams the response
func (a *App) SendMessageStreamWithImage(message string, imageData string) error {
	if a.httpClient == nil {
		return fmt.Errorf("llama client not initialized")
	}

	// Decode base64 image data if provided
	var imageBytes []byte
	var err error
	if imageData != "" {
		imageBytes, err = base64.StdEncoding.DecodeString(imageData)
		if err != nil {
			return fmt.Errorf("failed to decode image data: %w", err)
		}
	}

	// Add user message to history first
	a.chatMessages = append(a.chatMessages, map[string]interface{}{
		"role":      "user",
		"content":   message,
		"timestamp": time.Now().Unix(),
	})

	// Build prompt from chat messages
	prompt := a.buildPrompt()

	// Convert sampler settings to InferenceOptions
	opts := llama.NewInferenceOptions()
	if s := a.samplerSettings; s != nil {
		if v, ok := s["temperature"].(float64); ok {
			opts.Temperature = v
		}
		if v, ok := s["top_p_enabled"].(bool); ok {
			opts.TopPEnabled = v
		}
		if v, ok := s["top_p"].(float64); ok {
			if opts.TopPEnabled {
				opts.TopP = v
			} else {
				opts.TopP = 1.0
			}
		}
		if v, ok := s["top_k"].(float64); ok {
			opts.TopK = int(v)
		}
		if v, ok := s["min_p"].(float64); ok {
			opts.MinP = v
		}
		if v, ok := s["typical_p_enabled"].(bool); ok {
			opts.TypicalPEnabled = v
		}
		if v, ok := s["typical_p"].(float64); ok {
			opts.TypicalP = v
		}
		if v, ok := s["repeat_penalty"].(float64); ok {
			opts.RepeatPenalty = v
		}
		if v, ok := s["repeat_last_n"].(float64); ok {
			opts.RepeatLastN = int(v)
		}
		if v, ok := s["mirostat"].(float64); ok {
			opts.Mirostat = int(v)
		}
		if v, ok := s["mirostat_tau"].(float64); ok {
			opts.MirostatTau = v
		}
		if v, ok := s["dynamic_temp_range"].(float64); ok {
			opts.DynamicTempRange = v
		}
		if v, ok := s["dynamic_temp_exponent"].(float64); ok {
			opts.DynamicTempExponent = v
		}
		if v, ok := s["dry_multiplier"].(float64); ok {
			opts.DRYMultiplier = v
		}
		if v, ok := s["dry_allowed_length"].(float64); ok {
			opts.DRYAllowedLength = int(v)
		}
		if v, ok := s["dry_base"].(float64); ok {
			opts.DRYBase = v
		}
		if v, ok := s["smoothing_factor"].(float64); ok {
			opts.SmoothingFactor = v
		}
		if v, ok := s["smoothing_curve"].(float64); ok {
			opts.SmoothingCurve = v
		}
		if v, ok := s["top_a"].(float64); ok {
			opts.TopA = v
		}
		if v, ok := s["epsilon_cutoff"].(float64); ok {
			opts.EpsilonCutoff = v
		}
		if v, ok := s["eta_cutoff"].(float64); ok {
			opts.EtaCutoff = v
		}
		if v, ok := s["no_repeat_ngram_size"].(float64); ok {
			opts.NoRepeatNGramSize = int(v)
		}
		if v, ok := s["seed"].(float64); ok {
			opts.Seed = int(v)
		}
	}

	// Create cancellable context for this stream
	streamCtx, cancel := context.WithCancel(context.Background())
	a.streamCancel = cancel

	// Emit start event
	runtime.EventsEmit(a.ctx, "chat:start")

	var response strings.Builder

	// Stream the response with or without image
	chunkCount := 0
	if len(imageBytes) > 0 {
		fmt.Printf("DEBUG: Starting image stream with prompt len=%d, image size=%d bytes\n", len(prompt), len(imageBytes))
		err := a.httpClient.GenerateStreamWithImage(prompt, imageBytes, opts, func(chunk string, done bool) {
			select {
			case <-streamCtx.Done():
				// Stream was cancelled
				fmt.Printf("DEBUG: Stream cancelled by context\n")
				return
			default:
				if done {
					fmt.Printf("DEBUG: Image stream completed, total chunks emitted=%d, response len=%d\n", chunkCount, response.Len())
					return
				}
				response.WriteString(chunk)
				runtime.EventsEmit(a.ctx, "chat:chunk", chunk)
				chunkCount++
				if chunkCount <= 5 || chunkCount%50 == 0 {
					fmt.Printf("DEBUG: Emitting chat:chunk #%d, len=%d: %q\n", chunkCount, len(chunk), chunk[:min(50, len(chunk))])
				}
			}
		})
		if err != nil {
			fmt.Printf("DEBUG: Image stream error: %v\n", err)
			runtime.EventsEmit(a.ctx, "chat:error", err.Error())
			return err
		}
		fmt.Printf("DEBUG: Image stream finished successfully, chunks=%d, response len=%d\n", chunkCount, response.Len())
	} else {
		err := a.httpClient.GenerateStream(prompt, opts, func(chunk string, done bool) {
			select {
			case <-streamCtx.Done():
				return
			default:
				if done {
					return
				}
				response.WriteString(chunk)
				runtime.EventsEmit(a.ctx, "chat:chunk", chunk)
			}
		})
		if err != nil {
			runtime.EventsEmit(a.ctx, "chat:error", err.Error())
			return err
		}
	}

	// Clear the cancel function when done
	a.streamCancel = nil

	// Add assistant message to history
	a.chatMessages = append(a.chatMessages, map[string]interface{}{
		"role":      "assistant",
		"content":   response.String(),
		"timestamp": time.Now().Unix(),
	})

	// Emit complete event
	runtime.EventsEmit(a.ctx, "chat:complete", response.String())
	return nil
}

// buildPrompt builds a prompt from the chat history
func (a *App) buildPrompt() string {
	var prompt string
	if loreBlock := a.resolveLorePromptBlock(); loreBlock != "" {
		prompt += "System: " + loreBlock + "\n"
	}
	for _, msg := range a.chatMessages {
		role, _ := msg["role"].(string)
		content, _ := msg["content"].(string)
		if role == "user" {
			prompt += "User: " + content + "\n"
		} else if role == "assistant" {
			prompt += "Assistant: " + content + "\n"
		}
	}
	prompt += "Assistant: "
	return prompt
}

func (a *App) resolveLorePromptBlock() string {
	if len(a.lorebooks) == 0 || len(a.chatMessages) == 0 {
		return ""
	}

	lastMessage := ""
	history := make([]lorebook.Message, 0, min(len(a.chatMessages), 6))
	start := len(a.chatMessages) - 6
	if start < 0 {
		start = 0
	}

	for i, msg := range a.chatMessages[start:] {
		role, _ := msg["role"].(string)
		content, _ := msg["content"].(string)
		if strings.TrimSpace(content) == "" {
			continue
		}
		if start+i == len(a.chatMessages)-1 && role == "user" {
			lastMessage = content
		} else {
			history = append(history, lorebook.Message{Role: role, Content: content})
		}
	}

	if lastMessage == "" {
		return ""
	}

	resolved := lorebook.Resolve(a.lorebooks, lorebook.ResolveRequest{
		Message:       lastMessage,
		History:       history,
		Actor:         lorebook.TriggerUser,
		Triggered:     a.triggeredLore,
		MaxEntries:    8,
		MaxCharacters: 2400,
	})

	if len(resolved) == 0 {
		return ""
	}

	if a.triggeredLore == nil {
		a.triggeredLore = map[string]bool{}
	}
	for _, entry := range resolved {
		if entry.TriggerFrequency != lorebook.FrequencyAlways {
			a.triggeredLore[entry.ID] = true
		}
	}

	return lorebook.BuildPromptBlock(resolved)
}

// ClearChatHistory clears the backend chat history
// This is a Wails binding that can be called from the frontend
func (a *App) ClearChatHistory() error {
	a.chatMessages = []map[string]interface{}{}
	a.triggeredLore = map[string]bool{}
	return nil
}

// SetLorebooks updates the active lorebook entries used by the prompt resolver.
func (a *App) SetLorebooks(entries []lorebook.Entry) error {
	a.lorebooks = entries
	return nil
}

// AbortStream aborts the current streaming response
// This is a Wails binding that can be called from the frontend
func (a *App) AbortStream() error {
	if a.streamCancel != nil {
		a.streamCancel()
		a.streamCancel = nil
		runtime.EventsEmit(a.ctx, "chat:aborted")
	}
	return nil
}

// LoadModel loads a model from the given path
// This is a Wails binding that can be called from the frontend
func (a *App) LoadModel(modelPath string) error {
	if modelPath == "" {
		return fmt.Errorf("model path cannot be empty")
	}

	// Stop any existing llama-server
	if a.serverProc != nil && a.serverProc.Process != nil {
		_ = a.serverProc.Process.Kill()
		_ = a.serverProc.Wait()
		fmt.Printf("INFO: stopped existing llama-server (pid %d)\n", a.serverProc.Process.Pid)
	}

	// Extract GGUF metadata to get chat template
	fmt.Printf("INFO: Reading GGUF metadata from: %s\n", modelPath)
	profile, err := gguf.ExtractModelProfile(modelPath)
	if err != nil {
		fmt.Printf("WARN: Failed to extract GGUF metadata: %v\n", err)
	} else {
		template := profile.Template
		// For Gemma 4 models, use known-good Jinja2 template
		// Handlebars templates from GGUF don't work well with llama.cpp
		if strings.Contains(template, "{{#each") || strings.Contains(profile.Family, "gemma4") {
			template = "{%- for message in messages %}{{ message.role }}: {{ message.content }}\n{% endfor %}assistant:\n"
			fmt.Printf("INFO: Using built-in Gemma 4 Jinja2 template\n")
		}
		fmt.Printf("INFO: Extracted model profile - Family: %s\n", profile.Family)
		// If we have a chat template and user hasn't set a custom one, use it
		if template != "" {
			if a.modelSettings == nil {
				a.modelSettings = make(map[string]interface{})
			}
			// Only set if not already configured by user
			if _, exists := a.modelSettings["custom_jinja_template"]; !exists {
				a.modelSettings["custom_jinja_template"] = template
				fmt.Printf("INFO: Using chat template\n")
			}
		}
	}

	// Spawn new llama-server process with current settings
	cmd, err := spawnLlamaServerWithModel(modelPath, a.modelSettings)
	if err != nil {
		return fmt.Errorf("failed to spawn llama-server: %w", err)
	}
	a.serverProc = cmd
	go func() {
		for i := 0; i < 120; i++ {
			time.Sleep(500 * time.Millisecond)
			// Animate progress 0→90% over the wait period
			progress := int(float64(i) / 120.0 * 90.0)
			runtime.EventsEmit(a.ctx, "model:progress", progress)
			if err := a.httpClient.Health(); err == nil {
				fmt.Println("INFO: model loaded and server ready")
				runtime.EventsEmit(a.ctx, "model:progress", 100)
				runtime.EventsEmit(a.ctx, "model:loaded", modelPath)
				return
			}
		}
		fmt.Println("ERROR: llama-server did not become ready after 60s")
		runtime.EventsEmit(a.ctx, "model:error", "server did not become ready after 60s")
	}()
	return nil
}

// UnloadModel kills the llama-server process, freeing all model memory
// This is a Wails binding that can be called from the frontend
func (a *App) UnloadModel() error {
	if a.serverProc != nil && a.serverProc.Process != nil {
		_ = a.serverProc.Process.Kill()
		_ = a.serverProc.Wait()
		a.serverProc = nil
		fmt.Println("INFO: llama-server killed, model memory freed")
	}
	runtime.EventsEmit(a.ctx, "model:unloaded")
	return nil
}

// IsModelLoaded returns true if a model is currently loaded
// This is a Wails binding that can be called from the frontend
func (a *App) IsModelLoaded() bool {
	if a.httpClient == nil {
		return false
	}

	return a.httpClient.IsLoaded()
}

// GetLoadingProgress returns the current model loading progress (0.0-1.0)
// This is a Wails binding that can be called from the frontend
func (a *App) GetLoadingProgress() float64 {
	// HTTP client doesn't provide progress tracking
	return 0.0
}

// GetLoadingStage returns the current loading stage
// This is a Wails binding that can be called from the frontend
func (a *App) GetLoadingStage() string {
	// HTTP client doesn't provide stage tracking
	return ""
}

// GetLoadedModelName returns the name of the currently loaded model
// This is a Wails binding that can be called from the frontend
func (a *App) GetLoadedModelName() string {
	if a.httpClient == nil {
		return "No model loaded"
	}

	info, err := a.httpClient.GetLoadedModel()
	if err != nil || info == nil {
		return "No model loaded"
	}

	if data, ok := info["data"].([]interface{}); ok && len(data) > 0 {
		if m, ok := data[0].(map[string]interface{}); ok {
			if id, ok := m["id"].(string); ok {
				return filepath.Base(id)
			}
		}
	}
	return "No model loaded"
}

// GetChatMessages returns all messages in the current chat context
// This is a Wails binding that can be called from the frontend
func (a *App) GetChatMessages() []map[string]interface{} {
	return a.chatMessages
}

// ClearChat clears the current chat context
// This is a Wails binding that can be called from the frontend
func (a *App) ClearChat() error {
	a.chatMessages = []map[string]interface{}{}
	a.triggeredLore = map[string]bool{}
	return nil
}

// GetSamplerSettings returns the current sampler settings
// This is a Wails binding that can be called from the frontend
func (a *App) GetSamplerSettings() map[string]interface{} {
	if a.samplerSettings == nil {
		return map[string]interface{}{}
	}
	return a.samplerSettings
}

// SetSamplerSettings updates the sampler settings
// This is a Wails binding that can be called from the frontend
func (a *App) SetSamplerSettings(settings map[string]interface{}) error {
	a.samplerSettings = settings
	return nil
}

// SetModelOptions updates the model loading options
// This is a Wails binding that can be called from the frontend
func (a *App) SetModelOptions(options map[string]interface{}) error {
	// Store options for future use (llama-server handles most of these)
	a.modelSettings = options
	return nil
}

// GetModelOptions returns the current model options
// This is a Wails binding that can be called from the frontend
func (a *App) GetModelOptions() map[string]interface{} {
	if a.modelSettings == nil {
		return map[string]interface{}{}
	}
	return a.modelSettings
}

// GetModelsDir returns the absolute path to the models directory
// This is a Wails binding that can be called from the frontend
func (a *App) GetModelsDir() string {
	cwd, _ := os.Getwd()
	modelsDir := filepath.Join(cwd, "models")
	if _, err := os.Stat(modelsDir); err != nil {
		exeDir := filepath.Dir(os.Args[0])
		modelsDir = filepath.Join(exeDir, "models")
	}
	_ = os.MkdirAll(modelsDir, 0755)
	return modelsDir
}

// OpenModelFilePicker opens a native file dialog and returns the selected model path
// This is a Wails binding that can be called from the frontend
func (a *App) OpenModelFilePicker() (string, error) {
	modelsDir := a.GetModelsDir()
	path, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title:            "Select Model File",
		DefaultDirectory: modelsDir,
		Filters: []runtime.FileFilter{
			{DisplayName: "GGUF Models (*.gguf)", Pattern: "*.gguf"},
			{DisplayName: "All Model Files (*.gguf;*.bin)", Pattern: "*.gguf;*.bin"},
		},
	})
	if err != nil {
		return "", fmt.Errorf("file dialog error: %w", err)
	}
	return path, nil
}

// LoadModelMetadata reads GGUF metadata from a model file and returns a ModelProfile
// This is a Wails binding that can be called from the frontend
func (a *App) LoadModelMetadata(path string) (*types.ModelProfile, error) {
	reader := gguf.NewReader()
	profile, err := reader.ReadProfile(path)
	if err != nil {
		return nil, err
	}
	return profile, nil
}

// UpdateRuntimeSettings updates the runtime settings for model loading
// This is a Wails binding that can be called from the frontend
func (a *App) UpdateRuntimeSettings(settings *types.RuntimeSettings) error {
	// Store settings for future use (llama-server handles most of these)
	// Convert to map for storage
	a.modelSettings = map[string]interface{}{
		"threads":         settings.Threads,
		"batch_size":      settings.BatchSize,
		"context_size":    settings.ContextSize,
		"gpu_layers":      settings.GPULayers,
		"rope_mode":       settings.RopeMode,
		"rope_factor":     settings.RopeFactor,
		"rope_base":       settings.RopeBase,
		"flash_attention": settings.FlashAttention,
		"tensor_split":    settings.TensorSplit,
		"main_gpu":        settings.MainGPU,
		"offload_kqv":     settings.OffloadKQV,
		"use_mmap":        settings.UseMMap,
		"use_mlock":       settings.UseMLock,
		"vocab_only":      settings.VocabOnly,
	}
	return nil
}

// UpdatePromptSettings updates the prompt settings for the prompt pipeline
// This is a Wails binding that can be called from the frontend
func (a *App) UpdatePromptSettings(settings *types.PromptSettings) error {
	// Store prompt settings for the prompt pipeline
	// Convert to map for storage
	if a.modelSettings == nil {
		a.modelSettings = make(map[string]interface{})
	}

	a.modelSettings["prompt_template"] = settings.PromptTemplate
	a.modelSettings["custom_jinja_template"] = settings.CustomJinjaTemplate
	a.modelSettings["system_prompt_override"] = settings.SystemPromptOverride
	a.modelSettings["user_prefix"] = settings.UserPrefix
	a.modelSettings["assistant_prefix"] = settings.AssistantPrefix
	a.modelSettings["stop_sequences"] = settings.StopSequences

	return nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
