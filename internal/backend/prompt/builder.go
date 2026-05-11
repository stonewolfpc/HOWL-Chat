/**
 * Prompt Pipeline Builder
 *
 * This package provides functionality for building prompts from messages,
 * settings, and templates. It handles the application-level prompt construction
 * that is not part of gollama.cpp.
 *
 * @package prompt
 */

package prompt

import (
	"fmt"
	"strings"

	"howl-chat/internal/backend/types"
)

// Builder handles prompt construction from messages and settings
type Builder struct {
	profile        *types.ModelProfile
	runtime        *types.RuntimeSettings
	promptSettings *types.PromptSettings
}

// NewBuilder creates a new prompt builder
func NewBuilder() *Builder {
	return &Builder{}
}

// SetProfile sets the model profile for template resolution
func (b *Builder) SetProfile(profile *types.ModelProfile) *Builder {
	b.profile = profile
	return b
}

// SetRuntimeSettings sets the runtime settings
func (b *Builder) SetRuntimeSettings(settings *types.RuntimeSettings) *Builder {
	b.runtime = settings
	return b
}

// SetPromptSettings sets the prompt settings
func (b *Builder) SetPromptSettings(settings *types.PromptSettings) *Builder {
	b.promptSettings = settings
	return b
}

// ChatMessage represents a chat message with role and content
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// BuildPrompt constructs the final prompt string from messages and settings
// Returns the prompt string and stop sequences
func (b *Builder) BuildPrompt(messages []ChatMessage, memory string) (string, []string, error) {
	if b.profile == nil {
		return "", nil, fmt.Errorf("no model profile loaded")
	}

	// Determine which template to use
	template := b.resolveTemplate()

	// Build system prompt with memory injection
	system := b.buildSystemPrompt(memory)

	// Apply prefixes to messages
	templatedMessages := b.applyPrefixes(messages)

	// Render template with system and messages
	prompt, err := b.renderTemplate(template, system, templatedMessages)
	if err != nil {
		return "", nil, err
	}

	// Resolve stop sequences
	stops := b.resolveStopSequences()

	return prompt, stops, nil
}

// resolveTemplate determines which template to use
func (b *Builder) resolveTemplate() string {
	if b.promptSettings != nil && b.promptSettings.PromptTemplate != "" {
		return b.promptSettings.PromptTemplate
	}
	if b.profile != nil && b.profile.Template != "" {
		return b.profile.Template
	}
	return "{{#each messages}}{{ role }}: {{ content }}\n{{/each}}assistant:"
}

// buildSystemPrompt constructs the system prompt with memory injection
func (b *Builder) buildSystemPrompt(memory string) string {
	var system string

	if b.promptSettings != nil && b.promptSettings.SystemPromptOverride != "" {
		system = b.promptSettings.SystemPromptOverride
	}

	// Inject memory if present
	if memory != "" {
		if system == "" {
			system = "Relevant memory:\n" + memory
		} else {
			system = system + "\n\nRelevant memory:\n" + memory
		}
	}

	return system
}

// applyPrefixes applies user and assistant prefixes to messages
func (b *Builder) applyPrefixes(messages []ChatMessage) []map[string]string {
	var result []map[string]string

	for _, msg := range messages {
		content := msg.Content

		// Apply user prefix
		if msg.Role == "user" && b.promptSettings != nil && b.promptSettings.UserPrefix != "" {
			content = b.promptSettings.UserPrefix + content
		}

		// Apply assistant prefix
		if msg.Role == "assistant" && b.promptSettings != nil && b.promptSettings.AssistantPrefix != "" {
			content = b.promptSettings.AssistantPrefix + content
		}

		result = append(result, map[string]string{
			"role":    msg.Role,
			"content": content,
		})
	}

	return result
}

// renderTemplate renders the template with system and messages
// This is a simplified template renderer that handles basic mustache-style templates
func (b *Builder) renderTemplate(template string, system string, messages []map[string]string) (string, error) {
	result := template

	// Replace system variable
	result = strings.ReplaceAll(result, "{{ system }}", system)

	// Handle messages loop (simplified)
	if strings.Contains(result, "{{#each messages}}") {
		startIdx := strings.Index(result, "{{#each messages}}")
		endIdx := strings.Index(result, "{{/each}}")
		
		if startIdx != -1 && endIdx != -1 {
			prefix := result[:startIdx]
			suffix := result[endIdx+len("{{/each}}"):]
			templateBody := result[startIdx+len("{{#each messages}}") : endIdx]
			
			var messagesStr strings.Builder
			for _, msg := range messages {
				msgBody := templateBody
				msgBody = strings.ReplaceAll(msgBody, "{{ role }}", msg["role"])
				msgBody = strings.ReplaceAll(msgBody, "{{ content }}", msg["content"])
				messagesStr.WriteString(msgBody)
			}
			
			result = prefix + messagesStr.String() + suffix
		}
	}

	return result, nil
}

// resolveStopSequences determines which stop sequences to use
func (b *Builder) resolveStopSequences() []string {
	if b.promptSettings != nil && len(b.promptSettings.StopSequences) > 0 {
		return b.promptSettings.StopSequences
	}
	if b.profile != nil && len(b.profile.StopSequences) > 0 {
		return b.profile.StopSequences
	}
	return []string{}
}
