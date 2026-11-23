package chat

import (
	"context"
	"fmt"
	"strings"
	"time"

	"howl-chat/internal/backend/llama"
)

// ChatClient is the completions surface a turn needs. The inference client the
// host owns satisfies it; the domain never reaches into process or settings state.
type ChatClient interface {
	GenerateChat(messages []llama.ChatMessage, options *llama.InferenceOptions) (string, error)
	GenerateChatStream(messages []llama.ChatMessage, options *llama.InferenceOptions, callback llama.TokenCallback) error
	GenerateChatStreamWithImage(messages []llama.ChatMessage, imageData []byte, options *llama.InferenceOptions, callback llama.TokenCallback) error
}

// Attribution says who produced an assistant turn. A named character turn keeps
// its identifier so the transcript can tell whose voice it was after a reload.
type Attribution struct {
	Sender      string
	Character   string
	CharacterID string
}

// StreamHooks are the host's rendering callbacks. Delta receives raw model text:
// splitting reasoning from visible text is a display concern and stays in the
// host. Unsaved reports a reply that is on screen but could not be written to the
// log, so the host can surface the refusal.
type StreamHooks struct {
	Delta   func(string)
	Unsaved func(error)
}

// Turn is one generation over a conversation: it assembles the prompt from the
// session, asks the client for the reply, and records the turn through the
// conversation. The host supplies only the sampler options and the display
// callbacks.
type Turn struct {
	Conversation *ActiveConversation
	Client       ChatClient
	Settings     PromptSettings
	Options      *llama.InferenceOptions
}

// SendStream records the user's turn, assembles the prompt and streams a reply.
//
// The user's turn reaches the session without a log line: the renderer records the
// user's message in the log before inference starts, and duplicating it here would
// replay every message twice on reload.
func (t *Turn) SendStream(ctx context.Context, message string, image []byte, attr Attribution, hooks StreamHooks) (string, error) {
	if t.Conversation == nil {
		return "", fmt.Errorf("no conversation")
	}
	t.Conversation.RecordUserTurn(message)
	messages := t.Conversation.Session().BuildMessages(t.Settings)
	return t.StreamMessages(ctx, messages, image, attr, hooks)
}

// StreamMessages streams a fully formed message list and records the reply. The
// character path uses it because the renderer already assembled that prompt.
func (t *Turn) StreamMessages(ctx context.Context, messages []llama.ChatMessage, image []byte, attr Attribution, hooks StreamHooks) (string, error) {
	if t.Client == nil {
		return "", fmt.Errorf("llama client not initialized")
	}

	var raw strings.Builder
	callback := func(chunk string, done bool) {
		if done {
			return
		}
		if ctx != nil && ctx.Err() != nil {
			return
		}
		raw.WriteString(chunk)
		if hooks.Delta != nil {
			hooks.Delta(chunk)
		}
	}

	var err error
	if len(image) > 0 {
		err = t.Client.GenerateChatStreamWithImage(messages, image, t.Options, callback)
	} else {
		err = t.Client.GenerateChatStream(messages, t.Options, callback)
	}
	if err != nil {
		return "", err
	}

	text := StripThinkBlocks(raw.String())
	if t.Conversation != nil {
		saveErr := t.Conversation.AppendTurn(SessionEntry{
			Role:        "assistant",
			Content:     text,
			Sender:      attr.Sender,
			Character:   attr.Character,
			CharacterID: attr.CharacterID,
			Timestamp:   time.Now().Format(time.RFC3339),
		})
		if saveErr != nil && hooks.Unsaved != nil {
			hooks.Unsaved(saveErr)
		}
	}
	return text, nil
}

// Complete asks for a whole reply at once and records it in the session.
//
// This path never wrote a log line: the caller's renderer records both turns
// itself. Persisting stays with Save, so the behaviour is unchanged.
func (t *Turn) Complete(ctx context.Context, message string, attr Attribution) (string, error) {
	if t.Conversation == nil {
		return "", fmt.Errorf("no conversation")
	}
	if t.Client == nil {
		return "", fmt.Errorf("llama client not initialized")
	}
	if ctx != nil && ctx.Err() != nil {
		return "", ctx.Err()
	}

	t.Conversation.RecordUserTurn(message)
	messages := t.Conversation.Session().BuildMessages(t.Settings)

	response, err := t.Client.GenerateChat(messages, t.Options)
	if err != nil {
		return "", err
	}
	response = StripThinkBlocks(response)
	t.Conversation.Session().AppendEntry(SessionEntry{
		Role:        "assistant",
		Content:     response,
		Sender:      attr.Sender,
		Character:   attr.Character,
		CharacterID: attr.CharacterID,
		Timestamp:   time.Now().Format(time.RFC3339),
	})
	return response, nil
}
