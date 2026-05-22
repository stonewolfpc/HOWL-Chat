package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"howl-chat/internal/backend/chat"
	"howl-chat/internal/backend/llama"
	"howl-chat/internal/backend/lorebook"
	"howl-chat/internal/errors"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// ---------------------------------------------------------------------------
// Streaming message handlers
// ---------------------------------------------------------------------------

// SendMessageStreamWithImage streams an LLM response, optionally with an image.
func (a *App) SendMessageStreamWithImage(message string, imageData string) error {
	if a.httpClient == nil {
		return fmt.Errorf("llama client not initialized")
	}

	var imageBytes []byte
	if imageData != "" {
		var err error
		imageBytes, err = base64.StdEncoding.DecodeString(imageData)
		if err != nil {
			return fmt.Errorf("failed to decode image data: %w", err)
		}
		// Check if the loaded model supports vision
		if a.currentModelProfile == nil || !a.currentModelProfile.HasVision {
			errMsg := "This model does not support image input. Vision models require an mmproj file in the same directory."
			if a.currentModelProfile != nil && a.currentModelProfile.MMPProjPath == "" {
				errMsg = fmt.Sprintf("No mmproj file found for %s in %s. Vision models require an mmproj projection file alongside the model.",
					a.currentModelProfile.Family, filepath.Dir(a.currentModelProfile.Path))
			}
			runtime.EventsEmit(a.ctx, "chat:error", errMsg)
			return fmt.Errorf("%s", errMsg)
		}
	}

	a.chatMessages = append(a.chatMessages, map[string]interface{}{
		"role":      "user",
		"content":   message,
		"timestamp": time.Now().Unix(),
	})

	prompt := a.buildPrompt()
	opts := a.buildInferenceOptions()

	streamCtx, cancel := context.WithCancel(context.Background())
	a.streamCancel = cancel
	runtime.EventsEmit(a.ctx, "chat:start")

	var response strings.Builder

	var streamErr error
	if len(imageBytes) > 0 {
		streamErr = a.httpClient.GenerateStreamWithImage(prompt, imageBytes, opts, func(chunk string, done bool) {
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
	} else {
		streamErr = a.httpClient.GenerateStream(prompt, opts, func(chunk string, done bool) {
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
	}

	a.streamCancel = nil

	if streamErr != nil {
		a.logError(errors.LLMStreamInterrupted, errors.High, "llm.infer", "app-prompt.go:82", streamErr)
		runtime.EventsEmit(a.ctx, "chat:error", streamErr.Error())
		return streamErr
	}

	a.chatMessages = append(a.chatMessages, map[string]interface{}{
		"role":      "assistant",
		"content":   response.String(),
		"timestamp": time.Now().Unix(),
	})
	runtime.EventsEmit(a.ctx, "chat:complete", map[string]interface{}{"message": response.String()})

	// Autosave so assistant messages persist
	if a.chatPersistence != nil {
		go func() {
			if err := a.performSave(); err != nil {
				fmt.Printf("Autosave failed: %v\n", err)
			}
		}()
	}
	return nil
}

// GenerateCharacterResponse generates a response for a specific character and emits
// character-scoped events: chat:chunk:<id>, chat:complete:<id>, chat:error:<id>.
func (a *App) GenerateCharacterResponse(message string, characterID string) error {
	if a.httpClient == nil {
		return fmt.Errorf("llama client not initialized")
	}

	var character *chat.CharacterContext
	for _, c := range a.characters {
		if c.ID == characterID {
			character = c
			break
		}
	}
	if character == nil {
		return fmt.Errorf("character not found: %s", characterID)
	}

	charName := character.Name
	var pb strings.Builder
	fmt.Fprintf(&pb, "You are %s, a distinct participant in a group roleplay conversation.\n", charName)
	pb.WriteString("You speak only when someone addresses you directly by name, role, or context, or when it is naturally your turn to respond.\n")
	pb.WriteString("If others are speaking to each other, you listen quietly and stay in character, reacting only when your involvement makes sense.\n\n")
	if character.Age != "" {
		fmt.Fprintf(&pb, "Age: %s.\n", character.Age)
	}
	if character.Role != "" {
		fmt.Fprintf(&pb, "Role/Relationship: %s.\n", character.Role)
	}
	if character.Description != "" {
		fmt.Fprintf(&pb, "Appearance: %s\n", character.Description)
	}
	if character.Personality != "" {
		fmt.Fprintf(&pb, "Personality: %s\n", character.Personality)
	}
	if character.Backstory != "" {
		fmt.Fprintf(&pb, "Backstory: %s\n", character.Backstory)
	}
	if character.Background != "" {
		fmt.Fprintf(&pb, "Background: %s\n", character.Background)
	}
	if loreBlock := a.resolveCharacterLorePromptBlock(charName); loreBlock != "" {
		fmt.Fprintf(&pb, "\n%s\n", loreBlock)
	}
	fmt.Fprintf(&pb, "\nRespond to: %s", message)

	opts := a.buildInferenceOptions()
	streamCtx, cancel := context.WithCancel(context.Background())
	a.streamCancel = cancel

	eventChunk := fmt.Sprintf("chat:chunk:%s", characterID)
	eventComplete := fmt.Sprintf("chat:complete:%s", characterID)
	eventError := fmt.Sprintf("chat:error:%s", characterID)
	eventAborted := fmt.Sprintf("chat:aborted:%s", characterID)

	var response strings.Builder
	err := a.httpClient.GenerateStream(pb.String(), opts, func(chunk string, done bool) {
		select {
		case <-streamCtx.Done():
			return
		default:
			if done {
				return
			}
			response.WriteString(chunk)
			runtime.EventsEmit(a.ctx, eventChunk, chunk)
		}
	})

	if err != nil {
		select {
		case <-streamCtx.Done():
			runtime.EventsEmit(a.ctx, eventAborted)
			a.streamCancel = nil
			return nil
		default:
			a.logError(errors.LLMStreamInterrupted, errors.High, "llm.infer.character", "app-prompt.go:165", err)
			runtime.EventsEmit(a.ctx, eventError, err.Error())
			a.streamCancel = nil
			return err
		}
	}

	a.streamCancel = nil
	_ = a.AppendChatMessage("assistant", response.String(), charName)
	runtime.EventsEmit(a.ctx, "chat:complete", map[string]interface{}{
		"message":     response.String(),
		"characterId": characterID,
	})
	runtime.EventsEmit(a.ctx, eventComplete, response.String())
	return nil
}

// ---------------------------------------------------------------------------
// Chat history helpers
// ---------------------------------------------------------------------------

// AppendChatMessage records a message in the in-memory history and appends it
// to history.ndjson — O(1) disk write, no full file rewrite.
func (a *App) AppendChatMessage(role, content, sender string) error {
	if strings.TrimSpace(content) == "" {
		return nil
	}
	a.chatMessages = append(a.chatMessages, map[string]interface{}{
		"role":      role,
		"content":   content,
		"sender":    sender,
		"timestamp": time.Now().Unix(),
	})
	chatID := ""
	if a.chatPersistence != nil {
		chatID = a.chatPersistence.GetLastChatID()
	}
	if chatID == "" {
		return nil
	}
	entry := chat.ChatHistory{
		ID:        fmt.Sprintf("%d", len(a.chatMessages)-1),
		Sender:    sender,
		Text:      content,
		Role:      role,
		Timestamp: time.Now().Format(time.RFC3339),
	}
	return a.chatPersistence.AppendHistoryEntry(chatID, entry)
}

// ClearChatHistory wipes the in-memory chat history and triggered lore state.
func (a *App) ClearChatHistory() error {
	a.chatMessages = []map[string]interface{}{}
	a.triggeredLore = map[string]bool{}
	return nil
}

// AddSystemMessage appends a system-level narrator event to the chat history.
func (a *App) AddSystemMessage(content string) error {
	a.chatMessages = append(a.chatMessages, map[string]interface{}{
		"role":      "system",
		"content":   content,
		"timestamp": time.Now().Unix(),
	})
	return nil
}

// AbortStream cancels any in-progress generation stream.
func (a *App) AbortStream() error {
	if a.httpClient != nil {
		a.httpClient.CancelRequest()
	}
	if a.streamCancel != nil {
		a.streamCancel()
		a.streamCancel = nil
	}
	runtime.EventsEmit(a.ctx, "chat:aborted")
	return nil
}

// ---------------------------------------------------------------------------
// Lorebook bindings
// ---------------------------------------------------------------------------

// SetLorebooks replaces the active lorebook pool from frontend-supplied maps.
func (a *App) SetLorebooks(entries []map[string]interface{}) error {
	parsed := make([]lorebook.Entry, 0, len(entries))
	for _, m := range entries {
		entry, err := lorebook.ParseLorebookEntryFromMap(m)
		if err != nil {
			fmt.Printf("SetLorebooks: failed to parse entry: %v\n", err)
			continue
		}
		parsed = append(parsed, entry)
	}
	a.lorebooks = parsed
	return nil
}

// UpdateLorebookContext enables only the entries whose IDs are in activeIDs.
func (a *App) UpdateLorebookContext(activeIDs []string) error {
	active := make(map[string]bool, len(activeIDs))
	for _, id := range activeIDs {
		active[id] = true
	}
	for i := range a.lorebooks {
		a.lorebooks[i].Enabled = active[a.lorebooks[i].ID]
	}
	return nil
}

// ---------------------------------------------------------------------------
// Prompt construction
// ---------------------------------------------------------------------------

// buildPrompt builds a ChatML-formatted prompt from the current app state.
func (a *App) buildPrompt() string {
	var p strings.Builder

	userPrefix, assistantPrefix := "", ""
	if a.modelSettings != nil {
		if v, ok := a.modelSettings["user_prefix"].(string); ok {
			userPrefix = v
		}
		if v, ok := a.modelSettings["assistant_prefix"].(string); ok {
			assistantPrefix = v
		}
	}

	// Determine if we have any roleplay context loaded
	hasRoleplayContext := a.world != nil || a.scenario != nil || len(a.characters) > 0

	p.WriteString("<|im_start|>system\n")
	if a.modelSettings != nil {
		if spo, ok := a.modelSettings["system_prompt_override"].(string); ok && spo != "" {
			fmt.Fprintf(&p, "%s\n\n", spo)
		} else if hasRoleplayContext {
			p.WriteString("You are in an immersive roleplay. Stay in character at all times.\n\n")
		} else {
			p.WriteString("You are a helpful, knowledgeable assistant.\n\n")
		}
	} else if hasRoleplayContext {
		p.WriteString("You are in an immersive roleplay. Stay in character at all times.\n\n")
	} else {
		p.WriteString("You are a helpful, knowledgeable assistant.\n\n")
	}

	if a.world != nil {
		if a.world.Name != "" {
			fmt.Fprintf(&p, "[WORLD: %s]\n", a.world.Name)
		}
		if a.world.Description != "" {
			fmt.Fprintf(&p, "Description: %s\n", a.world.Description)
		}
		if len(a.world.KeyFacts) > 0 {
			fmt.Fprintf(&p, "Key Facts: %s\n", strings.Join(a.world.KeyFacts, ", "))
		}
		p.WriteString("\n")
	}

	if a.scenario != nil {
		if a.scenario.Name != "" {
			fmt.Fprintf(&p, "[SCENARIO: %s]\n", a.scenario.Name)
		}
		if a.scenario.Description != "" {
			fmt.Fprintf(&p, "Description: %s\n", a.scenario.Description)
		}
		if len(a.scenario.KeyDetails) > 0 {
			fmt.Fprintf(&p, "Key Details: %s\n", strings.Join(a.scenario.KeyDetails, ", "))
		}
		p.WriteString("\n")
	}

	if len(a.characters) > 0 {
		p.WriteString("[CHARACTER SHEETS]\n")
		for _, char := range a.characters {
			if char.Name != "" {
				fmt.Fprintf(&p, "Character: %s\n", char.Name)
			}
			if char.Personality != "" {
				fmt.Fprintf(&p, "Personality: %s\n", char.Personality)
			}
			if char.Backstory != "" {
				fmt.Fprintf(&p, "Backstory: %s\n", char.Backstory)
			}
			if len(char.Goals) > 0 {
				fmt.Fprintf(&p, "Goals: %s\n", strings.Join(char.Goals, ", "))
			}
			if len(char.Relationships) > 0 {
				fmt.Fprintf(&p, "Relationships: %s\n", strings.Join(char.Relationships, ", "))
			}
			p.WriteString("\n")
		}
		p.WriteString("When replying, use the character information above to keep each response aligned to the character sheet and the scenario context.\n\n")
	}

	if loreBlock := a.resolveLorePromptBlock(); loreBlock != "" {
		fmt.Fprintf(&p, "%s\n", loreBlock)
	}

	if hasRoleplayContext {
		p.WriteString("[RULE: Do not show internal thinking. Do not use 'Here is my thought process'. Start your response immediately as the character.]\n")
	} else {
		p.WriteString("[RULE: Do not show internal thinking. Do not use 'Here is my thought process'. Start your response immediately.]\n")
	}
	p.WriteString("```\n")

	for _, msg := range a.chatMessages {
		role, _ := msg["role"].(string)
		content, _ := msg["content"].(string)
		if role == "user" && userPrefix != "" {
			content = userPrefix + content
		} else if role == "assistant" && assistantPrefix != "" {
			content = assistantPrefix + content
		}
		fmt.Fprintf(&p, "<|im_start|>%s\n%s<|im_end|>\n", role, content)
	}
	p.WriteString("<|im_start|>assistant\n")
	return p.String()
}

// resolveLorePromptBlock resolves active lorebooks against the latest user message.
func (a *App) resolveLorePromptBlock() string {
	return a.resolveLore("")
}

// resolveCharacterLorePromptBlock resolves lorebooks scoped to a specific character.
func (a *App) resolveCharacterLorePromptBlock(characterName string) string {
	return a.resolveLore(characterName)
}

// resolveLore is the shared implementation behind both lore resolution paths.
func (a *App) resolveLore(characterName string) string {
	if len(a.lorebooks) == 0 || len(a.chatMessages) == 0 {
		return ""
	}

	start := len(a.chatMessages) - 6
	if start < 0 {
		start = 0
	}

	lastMessage := ""
	history := make([]lorebook.Message, 0, 6)
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

	scenarioName := ""
	if a.scenario != nil {
		scenarioName = a.scenario.Name
	}
	worldName := ""
	if a.world != nil {
		worldName = a.world.Name
	}
	routed := lorebook.RouteResolvedEntries(resolved, characterName, scenarioName, worldName)
	return lorebook.BuildPromptBlock(routed)
}

// ---------------------------------------------------------------------------
// Sampler → InferenceOptions conversion
// ---------------------------------------------------------------------------

// buildInferenceOptions maps the current samplerSettings map to a typed
// llama.InferenceOptions struct, applying stop-sequence overrides.
func (a *App) buildInferenceOptions() *llama.InferenceOptions {
	opts := llama.NewInferenceOptions()
	if s := a.samplerSettings; s != nil {
		setFloat := func(dst *float64, key string) {
			if v, ok := s[key].(float64); ok {
				*dst = v
			}
		}
		setInt := func(dst *int, key string) {
			if v, ok := s[key].(float64); ok {
				*dst = int(v)
			}
		}
		setBool := func(dst *bool, key string) {
			if v, ok := s[key].(bool); ok {
				*dst = v
			}
		}
		setFloat(&opts.Temperature, "temperature")
		setBool(&opts.TopPEnabled, "top_p_enabled")
		setFloat(&opts.TopP, "top_p")
		if !opts.TopPEnabled {
			opts.TopP = 1.0
		}
		setInt(&opts.TopK, "top_k")
		setFloat(&opts.MinP, "min_p")
		setBool(&opts.TypicalPEnabled, "typical_p_enabled")
		setFloat(&opts.TypicalP, "typical_p")
		setFloat(&opts.RepeatPenalty, "repeat_penalty")
		setInt(&opts.RepeatLastN, "repeat_last_n")
		setBool(&opts.FrequencyPenaltyEnabled, "frequency_penalty_enabled")
		setFloat(&opts.FrequencyPenalty, "frequency_penalty")
		setBool(&opts.PresencePenaltyEnabled, "presence_penalty_enabled")
		setFloat(&opts.PresencePenalty, "presence_penalty")
		setBool(&opts.MirostatEnabled, "mirostat_enabled")
		setInt(&opts.Mirostat, "mirostat")
		setFloat(&opts.MirostatTau, "mirostat_tau")
		setFloat(&opts.MirostatETA, "mirostat_eta")
		setBool(&opts.DynamicTempRangeEnabled, "dynamic_temp_range_enabled")
		setFloat(&opts.DynamicTempRange, "dynamic_temp_range")
		setFloat(&opts.DynamicTempExponent, "dynamic_temp_exponent")
		setFloat(&opts.DRYMultiplier, "dry_multiplier")
		setInt(&opts.DRYAllowedLength, "dry_allowed_length")
		setFloat(&opts.DRYBase, "dry_base")
		setFloat(&opts.SmoothingFactor, "smoothing_factor")
		setFloat(&opts.SmoothingCurve, "smoothing_curve")
		setBool(&opts.TopAEnabled, "top_a_enabled")
		setFloat(&opts.TopA, "top_a")
		setFloat(&opts.EpsilonCutoff, "epsilon_cutoff")
		setFloat(&opts.EtaCutoff, "eta_cutoff")
		setInt(&opts.NoRepeatNGramSize, "no_repeat_ngram")
		setFloat(&opts.EncoderRepeatPenalty, "encoder_repeat_penalty")
		setInt(&opts.Seed, "seed")
	}

	stops := []string{"<|im_end|>", "\nUser:", "System:", "\n<|think>"}
	if ms := a.modelSettings; ms != nil {
		if custom, ok := ms["stop_sequences"].([]string); ok && len(custom) > 0 {
			stops = custom
		}
	}
	opts.StopStrings = stops
	return opts
}
