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
// Turn wiring
//
// Prompt assembly and the send pipeline live in the chat domain (chat.Turn and
// Session.BuildMessages), reading the conversation directly. What stays here is
// what only the host can supply: the sampler options derived from the settings
// authority, the display splitter, the event sink and the error log.
// ---------------------------------------------------------------------------

// turn returns the domain turn for the active conversation with the host's
// sampler and prompt settings attached.
func (a *App) turn() *chat.Turn {
	return &chat.Turn{
		Conversation: a.conversation,
		Client:       a.engine.Client(),
		Settings:     a.promptSettings(),
		Options:      a.buildInferenceOptions(),
	}
}

// promptSettings collects the non-conversation input prompt assembly needs.
func (a *App) promptSettings() chat.PromptSettings {
	return chat.PromptSettings{
		ModelSettings: a.settings.Model(),
		ContextSize:   a.contextSize(),
	}
}

// contextSize reports the loaded model's window, or 0 when none is loaded.
func (a *App) contextSize() int {
	if a.engine.Client() == nil {
		return 0
	}
	return a.engine.Client().GetContextSize()
}

// ---------------------------------------------------------------------------
// Streaming message handlers
// ---------------------------------------------------------------------------

// SendMessageStreamWithImage streams an LLM response, optionally with an image.
func (a *App) SendMessageStreamWithImage(message string, imageData string) error {
	if err := a.ensureAcceptingWork("send message"); err != nil {
		return err
	}
	if a.engine.Client() == nil {
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
		if a.engine.Profile() == nil || !a.engine.Profile().HasVision {
			errMsg := "This model does not support image input. Image input requires a projector whose metadata declares a vision encoder."
			if a.engine.Profile() != nil && a.engine.Profile().MMPProjPath == "" {
				errMsg = fmt.Sprintf("No projector found for %s in %s. Image input requires a projector that declares a vision encoder alongside the model.",
					a.engine.Profile().Family, filepath.Dir(a.engine.Profile().Path))
			} else if kind := a.engine.Profile().ProjectorKind; kind == "hybrid" || kind == "vision" {
				errMsg = fmt.Sprintf("The projector beside %s declares %q, but image input has not been verified for this model and is not advertised.",
					a.engine.Profile().Family, kind)
			} else {
				errMsg = fmt.Sprintf("The projector beside %s declares %q, not a vision encoder, so image input is unavailable.",
					a.engine.Profile().Family, a.engine.Profile().ProjectorKind)
			}
			a.emitEvent("chat:error", errMsg)
			return fmt.Errorf("%s", errMsg)
		}
	}

	streamCtx, cancel := context.WithCancel(context.Background())
	streamToken := a.level0().Streams.Begin(cancel)
	defer a.level0().Streams.End(streamToken)

	// Thinking goes to its own channel: it is shown (or not) but never becomes
	// the reply, never enters history and is never sent back next turn.
	turn := a.turn()
	turn.Options.OnReasoning = func(chunk string) {
		select {
		case <-streamCtx.Done():
		default:
			a.emitEvent("chat:reasoning", chunk)
		}
	}
	a.emitEvent("chat:start")

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
				a.emitEvent("chat:chunk", text)
			}
		},
	)

	response, streamErr := turn.SendStream(streamCtx, message, imageBytes, chat.Attribution{Sender: "AI"}, chat.StreamHooks{
		Delta: splitter.Write,
		Unsaved: func(err error) {
			// The reply is displayed but is not in the authoritative log. Report
			// the refusal: a reply that was not written must never look saved.
			a.logError(errors.CHATSaveFail, errors.High, "chat.append", "app-prompt.go:138", err)
			a.emitEvent("chat:unsaved", map[string]interface{}{
				"role":   "assistant",
				"reason": err.Error(),
			})
		},
	})
	splitter.Flush()

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

	a.emitEvent("chat:complete", map[string]interface{}{"message": response})

	// Autosave so assistant messages persist
	if a.conversation.Bound() {
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
	if err := a.ensureAcceptingWork("character response"); err != nil {
		return err
	}
	if a.engine.Client() == nil {
		return fmt.Errorf("llama client not initialized")
	}

	var character *chat.CharacterContext
	for _, c := range a.sessionCharacters() {
		if c.ID == characterID {
			character = c
			break
		}
	}

	// Fallback: if the character isn't in memory, try to load from the
	// currently persisted chat.  This fixes desyncs where the frontend
	// chat data has a character that never made it into a.sessionCharacters().
	if character == nil && a.conversation.Bound() {
		if chatID := a.conversation.ActiveID(); chatID != "" {
			if chatData, err := a.conversation.Store().LoadChat(chatID); err == nil {
				for _, c := range chatData.Characters {
					if c.ID == characterID {
						character = c
						// Seed back into in-memory state so next call hits cache
						a.session().UpsertCharacter(c)
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
			a.session().UpsertCharacter(character)
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

	streamCtx, cancel := context.WithCancel(context.Background())
	streamToken := a.level0().Streams.Begin(cancel)
	defer a.level0().Streams.End(streamToken)

	eventChunk := fmt.Sprintf("chat:chunk:%s", characterID)
	eventComplete := fmt.Sprintf("chat:complete:%s", characterID)
	eventError := fmt.Sprintf("chat:error:%s", characterID)
	eventAborted := fmt.Sprintf("chat:aborted:%s", characterID)

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
				a.emitEvent(eventChunk, text)
			}
		},
	)

	response, err := a.turn().StreamMessages(
		streamCtx,
		messages,
		nil,
		chat.Attribution{Sender: character.Name, CharacterID: characterID},
		chat.StreamHooks{Delta: splitter.Write},
	)
	splitter.Flush()

	if err != nil {
		select {
		case <-streamCtx.Done():
			a.emitEvent(eventAborted)
			return nil
		default:
			a.logError(errors.LLMStreamInterrupted, errors.High, "llm.infer.character", "app-prompt.go:165", err)
			a.emitEvent(eventError, err.Error())
			return err
		}
	}

	a.emitEvent("chat:complete", map[string]interface{}{
		"message":     response,
		"characterId": characterID,
	})
	a.emitEvent(eventComplete, response)
	return nil
}

// ---------------------------------------------------------------------------
// Chat history helpers
// ---------------------------------------------------------------------------

// AppendChatMessage records a message in the in-memory history and appends it
// to history.ndjson — O(1) disk write, no full file rewrite.
//
// This public binding is used for user/system turns, which have no character
// owner. Character responses carry an explicit ID.
func (a *App) AppendChatMessage(role, content, sender string) error {
	return a.appendChatMessage(role, content, sender, "")
}

func (a *App) appendChatMessage(role, content, sender, characterID string) error {
	return a.conversation.AppendTurn(chat.SessionEntry{
		Role:        role,
		Content:     content,
		Sender:      sender,
		CharacterID: characterID,
	})
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
	sessionEntries := make([]chat.SessionEntry, 0, len(entries))
	now := time.Now().Format(time.RFC3339)

	for _, entry := range entries {
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
		characterID := historyFieldString(entry, "characterId", "CharacterId")
		hasImage, _ := entry["hasImage"].(bool)

		sessionEntry := chat.SessionEntry{
			Role:        role,
			Content:     text,
			Timestamp:   timestamp,
			HasImage:    hasImage,
			CharacterID: characterID,
		}
		// An assistant turn's named speaker is the character whose voice it is;
		// any other role's is plain sender attribution.
		if role == "assistant" {
			sessionEntry.Character = sender
		} else {
			sessionEntry.Sender = sender
		}
		sessionEntries = append(sessionEntries, sessionEntry)
	}

	return a.conversation.ReplaceHistory(sessionEntries)
}

// ClearChatHistory wipes the in-memory chat history and triggered lore state, and
// truncates the on-disk log for the active chat. Truncating is the important part:
// the delete/edit path calls this and then re-appends the surviving messages, so
// leaving the old lines on disk resurrected deleted messages on the next load.
func (a *App) ClearChatHistory() error {
	return a.conversation.ClearHistory()
}

// AddSystemMessage appends a system-level narrator event to the chat history.
func (a *App) AddSystemMessage(content string) error {
	a.session().AppendEntry(chat.SessionEntry{
		Role:      "system",
		Content:   content,
		Timestamp: time.Now().Format(time.RFC3339),
	})
	return nil
}

// AbortStream cancels any in-progress generation stream.
func (a *App) AbortStream() error {
	if a.engine.Client() != nil {
		a.engine.Client().CancelRequest()
	}
	a.level0().Streams.Cancel()
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
	// After cutover the canonical files own the text; a chat's saved snapshot is
	// only the selection. Resolve against canonical content so an edited lorebook
	// reaches the prompt without waiting for a new chat.
	a.session().SetLorebooks(a.mergeCanonicalLore(parsed, true))
	return nil
}

// UpdateLorebookContext enables only the entries whose IDs are in activeIDs.
func (a *App) UpdateLorebookContext(activeIDs []string) error {
	active := make(map[string]bool, len(activeIDs))
	for _, id := range activeIDs {
		active[id] = true
	}
	lorebooks := a.sessionLorebooks()
	for i := range lorebooks {
		lorebooks[i].Enabled = active[lorebooks[i].ID]
	}

	// A lorebook created after this chat was saved is not in the chat snapshot yet.
	// Enabling it must attach the canonical entry rather than silently doing nothing.
	known := make(map[string]bool, len(lorebooks))
	for _, entry := range lorebooks {
		known[entry.ID] = true
	}
	for _, id := range activeIDs {
		if known[id] {
			continue
		}
		entry, ok := a.canonicalLorebook(id)
		if !ok {
			continue
		}
		entry.Enabled = true
		lorebooks = append(lorebooks, entry)
		known[id] = true
	}
	a.session().SetLorebooks(lorebooks)
	return nil
}

// BuildCharacterLoreBlock is the one lore policy for character turns. The
// frontend composer calls it with the turn it is about to send, so ownership,
// visibility, activation, deleted and stale-entry rules are applied by the same
// resolver as the general prompt path instead of being rebuilt in JavaScript.
func (a *App) BuildCharacterLoreBlock(characterID string, message string) string {
	return a.session().BuildCharacterLoreBlock(characterID, message, a.contextSize())
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
	if s := a.settings.Sampler(); s != nil {
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
	if ms := a.settings.Model(); ms != nil {
		if custom, ok := ms["stop_sequences"].([]string); ok && len(custom) > 0 {
			stops = custom
		}
	}
	opts.StopStrings = stops

	// A thinking model spends its allowance on reasoning first, so the response
	// cap would otherwise truncate the answer. Give the thought its own headroom
	// on top of n_predict instead of eating the visible reply.
	if ms := a.settings.Model(); ms != nil {
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
