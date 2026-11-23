package chat

import (
	"fmt"
	"regexp"
	"strings"

	"howl-chat/internal/backend/llama"
	"howl-chat/internal/backend/lorebook"
)

// PromptSettings is the non-conversation input prompt assembly needs. Everything
// else it reads comes from the session itself.
type PromptSettings struct {
	// ModelSettings carries the prompt shaping keys the settings authority owns:
	// user_prefix, assistant_prefix, system_prompt_override.
	ModelSettings map[string]interface{}
	// ContextSize is the loaded model's window; 0 means unknown.
	ContextSize int
}

// BuildMessages assembles the chat-completions request for the next turn.
//
// It reads the session directly — turns, world, scenario, characters, lorebooks —
// so there is no second message representation between what is stored and what is
// sent.
func (s *Session) BuildMessages(settings PromptSettings) []llama.ChatMessage {
	userPrefix, assistantPrefix := "", ""
	if settings.ModelSettings != nil {
		if v, ok := settings.ModelSettings["user_prefix"].(string); ok {
			userPrefix = v
		}
		if v, ok := settings.ModelSettings["assistant_prefix"].(string); ok {
			assistantPrefix = v
		}
	}

	entries := s.Entries()
	world := s.World()
	characters := s.Characters()
	scenario := s.Scenario()

	hasRoleplayContext := world != nil || scenario != nil || len(characters) > 0

	var sys strings.Builder
	if settings.ModelSettings != nil {
		if spo, ok := settings.ModelSettings["system_prompt_override"].(string); ok && spo != "" {
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

	if world != nil {
		if world.Name != "" {
			fmt.Fprintf(&sys, "[WORLD: %s]\n", world.Name)
		}
		if world.Description != "" {
			fmt.Fprintf(&sys, "Description: %s\n", world.Description)
		}
		if len(world.KeyFacts) > 0 {
			fmt.Fprintf(&sys, "Key Facts: %s\n", strings.Join(world.KeyFacts, ", "))
		}
		sys.WriteString("\n")
	}

	if len(characters) > 0 {
		sys.WriteString("[CHARACTER SHEETS]\n")
		for _, char := range characters {
			if char == nil {
				continue
			}
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
	if scenario != nil {
		if scenario.Name != "" {
			fmt.Fprintf(&sys, "[SCENARIO: %s]\n", scenario.Name)
		}
		if scenario.Description != "" {
			fmt.Fprintf(&sys, "Description: %s\n", scenario.Description)
		}
		if len(scenario.KeyDetails) > 0 {
			fmt.Fprintf(&sys, "Key Details: %s\n", strings.Join(scenario.KeyDetails, ", "))
		}
		sys.WriteString("\n")
	}

	if loreBlock := s.lorePromptBlock("", settings.ContextSize); loreBlock != "" {
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

	for _, entry := range entries {
		if entry.Role == "" || entry.Content == "" {
			continue
		}
		content := entry.Content
		if entry.Role == "user" && userPrefix != "" {
			content = userPrefix + content
		} else if entry.Role == "assistant" && assistantPrefix != "" {
			content = assistantPrefix + content
		}
		messages = append(messages, llama.ChatMessage{Role: entry.Role, Content: content})
	}

	// Depth-0 instructions. Everything above sits at the head of the prompt, which
	// small models attend to least once the transcript grows, so the rules that
	// must survive a long chat are restated immediately before generation. They
	// ride on the final user turn rather than as a trailing system message: several
	// chat templates reject a system message after the first turn, and a rejected
	// template breaks generation outright. Prompt-only — the session's turns are
	// untouched, so nothing here reaches saved history.
	if n := len(messages); n > 0 && messages[n-1].Role == "user" {
		if content, ok := messages[n-1].Content.(string); ok {
			messages[n-1].Content = content + "\n\n[System note: stay in character. Do not describe your reasoning and do not show thinking. Do not write the user's words, actions, or thoughts. Reply only with the character's next response.]"
		}
	}

	return truncateMessagesToFitContext(messages, settings.ContextSize)
}

// BuildCharacterLoreBlock is the one lore policy for character turns. The
// composer calls it with the turn it is about to send, so ownership, visibility,
// activation, deleted and stale-entry rules are applied by the same resolver as
// the general prompt path instead of being rebuilt in the renderer.
func (s *Session) BuildCharacterLoreBlock(characterID, message string, contextSize int) string {
	if strings.TrimSpace(message) == "" {
		return ""
	}

	entries := s.Entries()
	history := make([]lorebook.Message, 0, 6)
	start := len(entries) - 6
	if start < 0 {
		start = 0
	}
	for _, entry := range entries[start:] {
		if strings.TrimSpace(entry.Content) == "" {
			continue
		}
		history = append(history, lorebook.Message{Role: entry.Role, Content: entry.Content})
	}

	return s.ResolveLoreBlock(characterID, message, history, contextSize)
}

// LoreBlockFor resolves the lore block for the latest user turn, optionally
// scoped to one character. It is the same route BuildMessages uses.
func (s *Session) LoreBlockFor(characterID string, contextSize int) string {
	return s.lorePromptBlock(characterID, contextSize)
}

// lorePromptBlock resolves active lorebooks against the latest user turn.
func (s *Session) lorePromptBlock(characterID string, contextSize int) string {
	entries := s.Entries()

	// Anchor on the most recent user turn, not on the final message: a regenerate
	// leaves the assistant reply last, and demanding a trailing user message
	// silently dropped every lore entry for that reply.
	lastUserIdx := -1
	for i := len(entries) - 1; i >= 0; i-- {
		if entries[i].Role == "user" {
			lastUserIdx = i
			break
		}
	}
	if lastUserIdx < 0 {
		return ""
	}

	lastMessage := entries[lastUserIdx].Content
	if strings.TrimSpace(lastMessage) == "" {
		return ""
	}

	// Six messages: the trigger turn plus the five before it.
	start := lastUserIdx - 5
	if start < 0 {
		start = 0
	}

	history := make([]lorebook.Message, 0, 6)
	for _, entry := range entries[start:lastUserIdx] {
		if strings.TrimSpace(entry.Content) == "" {
			continue
		}
		history = append(history, lorebook.Message{Role: entry.Role, Content: entry.Content})
	}

	return s.ResolveLoreBlock(characterID, lastMessage, history, contextSize)
}

// ResolveLoreBlock applies the budget, resolved-entry bookkeeping and scope
// routing. Callers differ only in where the turn text comes from.
func (s *Session) ResolveLoreBlock(characterID, lastMessage string, history []lorebook.Message, contextSize int) string {
	if len(s.Lorebooks()) == 0 {
		return ""
	}

	// Lore budget scales with the window instead of a fixed byte count: roughly a
	// quarter of the context, which is the share the SillyTavern defaults and the
	// community guides converge on. Characters are ~4x tokens, so a cap of `ctx`
	// characters is that quarter. The 2400 floor keeps small windows usable.
	maxLoreChars := 2400
	if contextSize > maxLoreChars {
		maxLoreChars = contextSize
	}

	// Trigger bookkeeping is session state: ResolveLore records which once-per-X
	// entries fired, so no caller holds that map.
	resolved := s.ResolveLore(lorebook.ResolveRequest{
		Message:       lastMessage,
		History:       history,
		Actor:         lorebook.TriggerUser,
		MaxEntries:    8,
		MaxCharacters: maxLoreChars,
	})
	if len(resolved) == 0 {
		return ""
	}

	scenarioName := ""
	if scenario := s.Scenario(); scenario != nil {
		scenarioName = scenario.Name
	}
	worldName := ""
	if world := s.World(); world != nil {
		worldName = world.Name
	}
	routed := lorebook.RouteResolvedEntries(resolved, characterID, scenarioName, worldName)
	return lorebook.BuildPromptBlock(routed)
}

// truncateMessagesToFitContext truncates the messages array to fit within the
// context window. It keeps the system prompt and the most recent messages,
// dropping older messages from the middle.
func truncateMessagesToFitContext(messages []llama.ChatMessage, contextSize int) []llama.ChatMessage {
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

	result = append(result, recentMessages...)

	fmt.Printf("[Context] Truncated from %d to %d messages to fit %d token context\n",
		len(messages), len(result), contextSize)

	return result
}

// thinkBlockRE matches a complete reasoning block in the shapes models actually
// emit:  thinking...</think>, <|think|>...<|/think|>, and the <thinking> variant,
// case-insensitively and across newlines.
var thinkBlockRE = regexp.MustCompile(`(?is)<\|?think(?:ing)?\|?>.*?<\|?/think(?:ing)?\|?>`)

// orphanCloseRE matches a closing reasoning tag with no opening tag before it.
// Models emit these when the opening tag was consumed by the reasoning parser, and
// leaving "</think>" in the visible answer is worse than removing it.
var orphanCloseRE = regexp.MustCompile(`(?is)<\|?/think(?:ing)?\|?>`)

// StripThinkBlocks removes reasoning from a reply before it is stored.
//
// Reasoning is for display only. If it is saved into the message body it becomes
// part of the transcript, gets counted against the context budget, and is fed back
// to the model on the next turn — three things the app must not do. An
// unterminated block means the model was cut off mid-thought, so everything from
// the opening tag on is dropped: a half-finished thought in history is worse than
// losing it.
func StripThinkBlocks(s string) string {
	if s == "" {
		return s
	}
	out := thinkBlockRE.ReplaceAllString(s, "")
	out = orphanCloseRE.ReplaceAllString(out, "")

	lower := strings.ToLower(out)
	for _, tag := range []string{"<|think", "<think"} {
		if i := strings.Index(lower, tag); i >= 0 {
			out = out[:i]
			lower = strings.ToLower(out)
		}
	}
	return strings.TrimSpace(out)
}
