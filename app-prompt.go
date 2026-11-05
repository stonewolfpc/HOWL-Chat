package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	stderrors "errors"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"howl-chat/internal/backend/chat"
	"howl-chat/internal/backend/llama"
	"howl-chat/internal/backend/lorebook"
	"howl-chat/internal/errors"
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
			errMsg := "This model does not support image input. Image input requires a projector whose metadata declares a vision encoder."
			if a.currentModelProfile != nil && a.currentModelProfile.MMPProjPath == "" {
				errMsg = fmt.Sprintf("No projector found for %s in %s. Image input requires a projector that declares a vision encoder alongside the model.",
					a.currentModelProfile.Family, filepath.Dir(a.currentModelProfile.Path))
			} else if kind := a.currentModelProfile.ProjectorKind; kind == "hybrid" || kind == "vision" {
				errMsg = fmt.Sprintf("The projector beside %s declares %q, but image input has not been verified for this model and is not advertised.",
					a.currentModelProfile.Family, kind)
			} else {
				errMsg = fmt.Sprintf("The projector beside %s declares %q, not a vision encoder, so image input is unavailable.",
					a.currentModelProfile.Family, a.currentModelProfile.ProjectorKind)
			}
			a.emitEvent("chat:error", errMsg)
			return fmt.Errorf("%s", errMsg)
		}
	}

	a.recordUserTurn(message)

	messages := a.buildPromptMessages()
	opts := a.buildInferenceOptions()

	streamCtx, cancel := context.WithCancel(context.Background())
	a.streamCancel = cancel

	// Thinking goes to its own channel: it is shown (or not) but never becomes
	// the reply, never enters history and is never sent back next turn.
	opts.OnReasoning = func(chunk string) {
		select {
		case <-streamCtx.Done():
		default:
			a.emitEvent("chat:reasoning", chunk)
		}
	}
	a.emitEvent("chat:start")

	var response strings.Builder

	// Route deltas through the splitter: a model that emits its thought inline
	// (tags inside the content stream) must still keep that thought out of the
	// visible reply, out of history and out of the next prompt.
	splitter := newThinkSplitter(
		func(t string) {
			select {
			case <-streamCtx.Done():
			default:
				a.emitEvent("chat:reasoning", t)
			}
		},
		func(text string) {
			select {
			case <-streamCtx.Done():
			default:
				response.WriteString(text)
				a.emitEvent("chat:chunk", text)
			}
		},
	)

	var streamErr error
	if len(imageBytes) > 0 {
		streamErr = a.httpClient.GenerateChatStreamWithImage(messages, imageBytes, opts, func(chunk string, done bool) {
			if done {
				return
			}
			splitter.Write(chunk)
		})
	} else {
		streamErr = a.httpClient.GenerateChatStream(messages, opts, func(chunk string, done bool) {
			if done {
				return
			}
			splitter.Write(chunk)
		})
	}
	splitter.Flush()

	a.streamCancel = nil

	if streamErr != nil {
		// A user abort is normal and must never surface as an error. AbortStream
		// cancels the stream context we own, and the HTTP client reports the
		// cancellation wrapped ("scanner error: context canceled"), so a direct
		// == comparison misses it. Trust our own context and use errors.Is.
		if streamCtx.Err() != nil || stderrors.Is(streamErr, context.Canceled) {
			fmt.Println("[SendMessageStream] stream canceled by user")
			a.emitEvent("chat:aborted")
			return nil
		}
		a.logError(errors.LLMStreamInterrupted, errors.High, "llm.infer", "app-prompt.go:82", streamErr)
		a.emitEvent("chat:error", streamErr.Error())
		return streamErr
	}

	// Route through AppendChatMessage so the reply reaches history.ndjson.
	// Appending only to a.chatMessages wrote it to history.json, which the
	// loader ignores whenever the append log is non-empty — so assistant
	// turns vanished on reload while user turns survived.
	_ = a.AppendChatMessage("assistant", stripThinkBlocks(response.String()), "AI")
	a.emitEvent("chat:complete", map[string]interface{}{"message": response.String()})

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

	// Fallback: if the character isn't in memory, try to load from the
	// currently persisted chat.  This fixes desyncs where the frontend
	// chat data has a character that never made it into a.characters.
	if character == nil && a.chatPersistence != nil {
		chatID := a.chatPersistence.GetLastChatID()
		if chatID != "" {
			if chatData, err := a.chatPersistence.LoadChat(chatID); err == nil {
				for _, c := range chatData.Characters {
					if c.ID == characterID {
						character = c
						// Seed back into in-memory state so next call hits cache
						a.characters = append(a.characters, c)
						break
					}
				}
			}
		}
	}

	// Fallback 2: scan the Characters/ directory for the matching ID.
	if character == nil {
		character = loadCharacterFromDiskByID(characterID)
		if character != nil {
			a.characters = append(a.characters, character)
		}
	}

	if character == nil {
		return fmt.Errorf("character not found: %s", characterID)
	}

	// The frontend already built a complete character-specific message including
	// character info, lore, memories, group context, and the user message.
	// Do NOT prepend backend character info again — that would double-prompt.
	messages := []llama.ChatMessage{
		{Role: "user", Content: message},
	}

	opts := a.buildInferenceOptions()
	streamCtx, cancel := context.WithCancel(context.Background())
	a.streamCancel = cancel

	eventChunk := fmt.Sprintf("chat:chunk:%s", characterID)
	eventComplete := fmt.Sprintf("chat:complete:%s", characterID)
	eventError := fmt.Sprintf("chat:error:%s", characterID)
	eventAborted := fmt.Sprintf("chat:aborted:%s", characterID)

	var response strings.Builder
	// Same inline-thought handling as the generic path: the character's visible
	// reply must never contain the thought or its tags.
	splitter := newThinkSplitter(
		func(t string) {
			select {
			case <-streamCtx.Done():
			default:
				a.emitEvent("chat:reasoning", t)
			}
		},
		func(text string) {
			select {
			case <-streamCtx.Done():
			default:
				response.WriteString(text)
				a.emitEvent(eventChunk, text)
			}
		},
	)

	err := a.httpClient.GenerateChatStream(messages, opts, func(chunk string, done bool) {
		if done {
			return
		}
		splitter.Write(chunk)
	})
	splitter.Flush()

	if err != nil {
		select {
		case <-streamCtx.Done():
			a.emitEvent(eventAborted)
			a.streamCancel = nil
			return nil
		default:
			a.logError(errors.LLMStreamInterrupted, errors.High, "llm.infer.character", "app-prompt.go:165", err)
			a.emitEvent(eventError, err.Error())
			a.streamCancel = nil
			return err
		}
	}

	a.streamCancel = nil
	_ = a.AppendChatMessage("assistant", stripThinkBlocks(response.String()), character.Name)
	a.emitEvent("chat:complete", map[string]interface{}{
		"message":     response.String(),
		"characterId": characterID,
	})
	a.emitEvent(eventComplete, response.String())
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

// historyFieldString returns the first non-empty string value among keys.
func historyFieldString(m map[string]interface{}, keys ...string) string {
	for _, k := range keys {
		if v, ok := m[k].(string); ok && v != "" {
			return v
		}
	}
	return ""
}

// ReplaceChatHistory replaces the whole chat history with the supplied entries and
// rewrites history.ndjson in place. This is the delete/edit path: appending the
// surviving messages onto the log would leave deleted messages on disk and they
// would be replayed on the next load.
func (a *App) ReplaceChatHistory(entries []map[string]interface{}) error {
	history := make([]chat.ChatHistory, 0, len(entries))
	messages := make([]map[string]interface{}, 0, len(entries))
	now := time.Now().Format(time.RFC3339)

	for i, entry := range entries {
		role := strings.ToLower(historyFieldString(entry, "role", "Role"))
		text := historyFieldString(entry, "text", "Text")
		if strings.TrimSpace(text) == "" {
			continue
		}
		sender := historyFieldString(entry, "sender", "Sender")
		if sender == "" {
			switch role {
			case "user":
				sender = "User"
			case "assistant":
				sender = "Assistant"
			default:
				sender = "System"
			}
		}
		timestamp := historyFieldString(entry, "timestamp", "Timestamp")
		if timestamp == "" {
			timestamp = now
		}

		history = append(history, chat.ChatHistory{
			ID:        fmt.Sprintf("%d", i),
			Sender:    sender,
			Text:      text,
			Role:      role,
			Timestamp: timestamp,
		})
		messages = append(messages, map[string]interface{}{
			"role":      role,
			"content":   text,
			"sender":    sender,
			"timestamp": time.Now().Unix(),
		})
	}

	a.chatMessages = messages

	if a.chatPersistence == nil {
		return nil
	}
	chatID := a.chatPersistence.GetLastChatID()
	if chatID == "" {
		return nil
	}
	return a.chatPersistence.RewriteHistory(chatID, history)
}

// ClearChatHistory wipes the in-memory chat history and triggered lore state, and
// truncates the on-disk log for the active chat. Truncating is the important part:
// the delete/edit path calls this and then re-appends the surviving messages, so
// leaving the old lines on disk resurrected deleted messages on the next load.
func (a *App) ClearChatHistory() error {
	a.chatMessages = []map[string]interface{}{}
	a.triggeredLore = map[string]bool{}
	if a.chatPersistence == nil {
		return nil
	}
	chatID := a.chatPersistence.GetLastChatID()
	if chatID == "" {
		return nil
	}
	return a.chatPersistence.RewriteHistory(chatID, []chat.ChatHistory{})
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
	a.emitEvent("chat:aborted")
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

// truncateMessagesToFitContext truncates the messages array to fit within the context window.
// It keeps the system prompt and the most recent messages, dropping older messages from the middle.
func (a *App) truncateMessagesToFitContext(messages []llama.ChatMessage) []llama.ChatMessage {
	if a.httpClient == nil {
		return messages
	}

	// Get the context size from the model
	contextSize := a.httpClient.GetContextSize()
	if contextSize <= 0 {
		contextSize = 2048 // fallback
	}

	// Rough estimate: 4 characters per token (conservative estimate)
	maxChars := contextSize * 4

	// Always keep the system prompt
	var result []llama.ChatMessage
	if len(messages) > 0 && messages[0].Role == "system" {
		result = append(result, messages[0])
	}

	// Calculate total characters
	totalChars := 0
	for _, msg := range messages {
		if content, ok := msg.Content.(string); ok {
			totalChars += len(content)
		}
	}

	// If we're within the context limit, return as-is
	if totalChars <= maxChars {
		return messages
	}

	// We need to truncate. Keep system prompt, last few messages, and drop from the middle
	// Reserve space for system prompt
	systemChars := 0
	if len(result) > 0 {
		if content, ok := result[0].Content.(string); ok {
			systemChars = len(content)
		}
	}
	availableChars := maxChars - systemChars

	// Keep the last few messages that fit
	var recentMessages []llama.ChatMessage
	usedChars := 0

	// Process messages in reverse order, keeping as many as will fit
	for i := len(messages) - 1; i >= 0; i-- {
		if i == 0 && messages[i].Role == "system" {
			continue // Already handled system prompt
		}

		var msgChars int
		if content, ok := messages[i].Content.(string); ok {
			msgChars = len(content)
		} else {
			continue // Skip messages without string content
		}
		if usedChars+msgChars > availableChars {
			break
		}

		recentMessages = append([]llama.ChatMessage{messages[i]}, recentMessages...)
		usedChars += msgChars
	}

	// Combine system prompt with recent messages
	result = append(result, recentMessages...)

	fmt.Printf("[Context] Truncated from %d to %d messages to fit %d token context\n",
		len(messages), len(result), contextSize)

	return result
}

// buildPromptMessages builds a proper messages array for the chat completions endpoint.
func (a *App) buildPromptMessages() []llama.ChatMessage {
	userPrefix, assistantPrefix := "", ""
	if a.modelSettings != nil {
		if v, ok := a.modelSettings["user_prefix"].(string); ok {
			userPrefix = v
		}
		if v, ok := a.modelSettings["assistant_prefix"].(string); ok {
			assistantPrefix = v
		}
	}

	hasRoleplayContext := a.world != nil || a.scenario != nil || len(a.characters) > 0

	var sys strings.Builder
	if a.modelSettings != nil {
		if spo, ok := a.modelSettings["system_prompt_override"].(string); ok && spo != "" {
			fmt.Fprintf(&sys, "%s\n\n", spo)
		} else if hasRoleplayContext {
			sys.WriteString("You are in an immersive roleplay. Stay in character at all times. Speak naturally and concisely unless the scene calls for detail.\n\n")
		} else {
			sys.WriteString("You are a helpful, knowledgeable assistant.\n\n")
		}
	} else if hasRoleplayContext {
		sys.WriteString("You are in an immersive roleplay. Stay in character at all times. Speak naturally and concisely unless the scene calls for detail.\n\n")
	} else {
		sys.WriteString("You are a helpful, knowledgeable assistant.\n\n")
	}

	if a.world != nil {
		if a.world.Name != "" {
			fmt.Fprintf(&sys, "[WORLD: %s]\n", a.world.Name)
		}
		if a.world.Description != "" {
			fmt.Fprintf(&sys, "Description: %s\n", a.world.Description)
		}
		if len(a.world.KeyFacts) > 0 {
			fmt.Fprintf(&sys, "Key Facts: %s\n", strings.Join(a.world.KeyFacts, ", "))
		}
		sys.WriteString("\n")
	}

	if len(a.characters) > 0 {
		sys.WriteString("[CHARACTER SHEETS]\n")
		for _, char := range a.characters {
			if char.Name != "" {
				fmt.Fprintf(&sys, "Character: %s\n", char.Name)
			}
			if char.Personality != "" {
				fmt.Fprintf(&sys, "Personality: %s\n", char.Personality)
			}
			if char.Backstory != "" {
				fmt.Fprintf(&sys, "Backstory: %s\n", char.Backstory)
			}
			if len(char.Goals) > 0 {
				fmt.Fprintf(&sys, "Goals: %s\n", strings.Join(char.Goals, ", "))
			}
			if len(char.Relationships) > 0 {
				fmt.Fprintf(&sys, "Relationships: %s\n", strings.Join(char.Relationships, ", "))
			}
			sys.WriteString("\n")
		}
		sys.WriteString("When replying, use the character information above to keep each response aligned to the character sheet and the scenario context.\n\n")
	}

	// The scene follows the character sheets on purpose: a card reads description,
	// personality, then scenario, so the scene lands as the setting those
	// characters are standing in rather than as a document ahead of them.
	if a.scenario != nil {
		if a.scenario.Name != "" {
			fmt.Fprintf(&sys, "[SCENARIO: %s]\n", a.scenario.Name)
		}
		if a.scenario.Description != "" {
			fmt.Fprintf(&sys, "Description: %s\n", a.scenario.Description)
		}
		if len(a.scenario.KeyDetails) > 0 {
			fmt.Fprintf(&sys, "Key Details: %s\n", strings.Join(a.scenario.KeyDetails, ", "))
		}
		sys.WriteString("\n")
	}

	if loreBlock := a.resolveLorePromptBlock(); loreBlock != "" {
		fmt.Fprintf(&sys, "%s\n", loreBlock)
	}

	if hasRoleplayContext {
		sys.WriteString("[RULE: Do not show internal thinking. Do not use 'Here is my thought process'. Start your response immediately as the character.]")
	} else {
		sys.WriteString("[RULE: Do not show internal thinking. Do not use 'Here is my thought process'. Start your response immediately.]")
	}

	messages := []llama.ChatMessage{
		{Role: "system", Content: sys.String()},
	}

	for _, msg := range a.chatMessages {
		role, _ := msg["role"].(string)
		content, _ := msg["content"].(string)
		if role == "" || content == "" {
			continue
		}
		if role == "user" && userPrefix != "" {
			content = userPrefix + content
		} else if role == "assistant" && assistantPrefix != "" {
			content = assistantPrefix + content
		}
		messages = append(messages, llama.ChatMessage{Role: role, Content: content})
	}

	// Depth-0 instructions. Everything above sits at the head of the prompt,
	// which small models attend to least once the transcript grows, so the rules
	// that must survive a long chat are restated immediately before generation.
	// They ride on the final user turn rather than as a trailing system message:
	// several chat templates reject a system message after the first turn, and a
	// rejected template breaks generation outright. Prompt-only — a.chatMessages
	// is untouched, so nothing here reaches saved history.
	if n := len(messages); n > 0 && messages[n-1].Role == "user" {
		if content, ok := messages[n-1].Content.(string); ok {
			messages[n-1].Content = content + "\n\n[System note: stay in character. Do not describe your reasoning and do not show thinking. Do not write the user's words, actions, or thoughts. Reply only with the character's next response.]"
		}
	}

	// Apply context truncation to prevent exceeding model's context window
	messages = a.truncateMessagesToFitContext(messages)

	return messages
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

	// Anchor on the most recent user turn, not on the final message: a
	// regenerate leaves the assistant reply last, and demanding a trailing user
	// message silently dropped every lore entry for that reply.
	lastUserIdx := -1
	for i := len(a.chatMessages) - 1; i >= 0; i-- {
		if role, _ := a.chatMessages[i]["role"].(string); role == "user" {
			lastUserIdx = i
			break
		}
	}
	if lastUserIdx < 0 {
		return ""
	}

	lastMessage, _ := a.chatMessages[lastUserIdx]["content"].(string)
	if strings.TrimSpace(lastMessage) == "" {
		return ""
	}

	// Six messages: the trigger turn plus the five before it.
	start := lastUserIdx - 5
	if start < 0 {
		start = 0
	}

	history := make([]lorebook.Message, 0, 6)
	for _, msg := range a.chatMessages[start:lastUserIdx] {
		role, _ := msg["role"].(string)
		content, _ := msg["content"].(string)
		if strings.TrimSpace(content) == "" {
			continue
		}
		history = append(history, lorebook.Message{Role: role, Content: content})
	}

	// Lore budget scales with the window instead of a fixed byte count: roughly a
	// quarter of the context, which is the share the SillyTavern defaults and the
	// community guides converge on. Characters are ~4x tokens, so a cap of `ctx`
	// characters is that quarter. The 2400 floor keeps small windows usable.
	maxLoreChars := 2400
	if a.httpClient != nil {
		if ctx := a.httpClient.GetContextSize(); ctx > maxLoreChars {
			maxLoreChars = ctx
		}
	}

	resolved := lorebook.Resolve(a.lorebooks, lorebook.ResolveRequest{
		Message:       lastMessage,
		History:       history,
		Actor:         lorebook.TriggerUser,
		Triggered:     a.triggeredLore,
		MaxEntries:    8,
		MaxCharacters: maxLoreChars,
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

	// Per-family defaults first: the publishers' own recommendations differ too
	// much for one set to serve them all. Anything the user has set below still
	// wins, because the sampler overrides are applied after this.
	familySource := familySamplingDefaults(opts, a.familyName())
	if s := a.samplerSettings; s != nil {
		setFloat := func(dst *float64, key string) {
			if v, ok := s[key].(float64); ok {
				*dst = v
			} else if v, ok := s[key].(int); ok {
				*dst = float64(v)
			} else if v, ok := s[key].(json.Number); ok {
				if fv, err := v.Float64(); err == nil {
					*dst = fv
				}
			}
		}
		setInt := func(dst *int, key string) {
			if v, ok := s[key].(float64); ok {
				*dst = int(v)
			} else if v, ok := s[key].(int); ok {
				*dst = v
			} else if v, ok := s[key].(json.Number); ok {
				if iv, err := v.Int64(); err == nil {
					*dst = int(iv)
				}
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
		setInt(&opts.NoRepeatNGramSize, "no_repeat_ngram_size")
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

	// A thinking model spends its allowance on reasoning first, so the response
	// cap would otherwise truncate the answer. Give the thought its own headroom
	// on top of n_predict instead of eating the visible reply.
	if ms := a.modelSettings; ms != nil {
		if mode := strSetting(ms, "think_mode", "off"); mode != "off" {
			if extra, err := strconv.Atoi(intSetting(ms, "think_budget", 2048)); err == nil && extra > 0 {
				opts.NPredict += extra
			}
		}
	}
	fmt.Printf("[inference] family=%q source=%q temp=%.2f top_p=%.2f top_k=%d min_p=%.2f repeat=%.2f n_predict=%d\n",
		a.familyName(), familySource, opts.Temperature, opts.TopP, opts.TopK, opts.MinP, opts.RepeatPenalty, opts.NPredict)

	return opts
}
