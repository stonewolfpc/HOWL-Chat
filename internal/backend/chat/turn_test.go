package chat

import (
	"context"
	"errors"
	entities "howl-chat/internal/backend/types/entities"
	"strings"
	"testing"

	"howl-chat/internal/backend/llama"
	"howl-chat/internal/backend/lorebook"
)

type fakeChatClient struct {
	chunks       []string
	reply        string
	err          error
	lastMessages []llama.ChatMessage
	sawImage     bool
}

func (f *fakeChatClient) GenerateChat(messages []llama.ChatMessage, _ *llama.InferenceOptions) (string, error) {
	f.lastMessages = messages
	if f.err != nil {
		return "", f.err
	}
	return f.reply, nil
}

func (f *fakeChatClient) GenerateChatStream(messages []llama.ChatMessage, _ *llama.InferenceOptions, callback llama.TokenCallback) error {
	f.lastMessages = messages
	if f.err != nil {
		return f.err
	}
	for _, chunk := range f.chunks {
		callback(chunk, false)
	}
	callback("", true)
	return nil
}

func (f *fakeChatClient) GenerateChatStreamWithImage(messages []llama.ChatMessage, _ []byte, _ *llama.InferenceOptions, callback llama.TokenCallback) error {
	f.sawImage = true
	return f.GenerateChatStream(messages, nil, callback)
}

func TestTurnSendStreamAssemblesFromTheSessionAndAppendsTheReply(t *testing.T) {
	c, pm := newTestConversation(t)
	meta, err := c.New("Turn")
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	c.Session().SetCharacters([]*CharacterContext{{Character: entities.Character{ID: "char_aria", Name: "Aria", Personality: "wry"}}})
	c.Session().SetLorebooks([]lorebook.Entry{{
		ID:               "lb1",
		Title:            "Sigil",
		Content:          "The sigil answers.",
		Enabled:          true,
		Scope:            lorebook.ScopeWorld,
		TriggerPhrases:   []string{"sigil"},
		TriggerMode:      lorebook.TriggerLoose,
		TriggerDirection: lorebook.TriggerBoth,
		PriorityLevel:    3,
	}})

	client := &fakeChatClient{chunks: []string{"The sigil ", "answers. ", "</think>only the answer"}}
	turn := &Turn{
		Conversation: c,
		Client:       client,
		Settings:     PromptSettings{ContextSize: 4096},
		Options:      llama.NewInferenceOptions(),
	}

	var streamed strings.Builder
	reply, err := turn.SendStream(context.Background(), "I touch the sigil", nil, Attribution{Sender: "AI"}, StreamHooks{
		Delta: func(s string) { streamed.WriteString(s) },
	})
	if err != nil {
		t.Fatalf("SendStream: %v", err)
	}
	if reply != "The sigil answers. only the answer" {
		t.Fatalf("reply = %q", reply)
	}
	if streamed.String() != client.chunks[0]+client.chunks[1]+client.chunks[2] {
		t.Fatalf("display deltas lost raw text: %q", streamed.String())
	}

	// The prompt must come from the session: system content, the character sheet
	// and the resolved lore, followed by the user's turn.
	if len(client.lastMessages) < 2 {
		t.Fatalf("assembled %d messages", len(client.lastMessages))
	}
	system, _ := client.lastMessages[0].Content.(string)
	for _, want := range []string{"[CHARACTER SHEETS]", "Aria", "The sigil answers.", "I touch the sigil"} {
		if !strings.Contains(system+joinContents(client.lastMessages), want) {
			t.Fatalf("prompt missing %q:\n%s", want, system)
		}
	}

	entries := c.Session().Entries()
	if len(entries) != 2 {
		t.Fatalf("session has %d turns, want user + assistant", len(entries))
	}
	if entries[0].Role != "user" || entries[0].Content != "I touch the sigil" {
		t.Fatalf("user turn = %+v", entries[0])
	}
	if entries[1].Role != "assistant" || entries[1].Content != reply {
		t.Fatalf("assistant turn = %+v", entries[1])
	}
	lines := logLines(t, pm, meta.ID)
	if len(lines) != 1 || !strings.Contains(lines[0], "only the answer") {
		t.Fatalf("log lines = %v", lines)
	}
	if strings.Contains(lines[0], "think") {
		t.Fatalf("reasoning reached the log: %s", lines[0])
	}

	// Records the same history a relaunch would replay.
	restarted := NewActiveConversation(pm)
	if _, err := restarted.Load(meta.ID); err != nil {
		t.Fatalf("Load: %v", err)
	}
	restored := restarted.Session().Entries()
	if len(restored) != 1 || restored[0].Content != reply {
		t.Fatalf("restored history = %+v", restored)
	}
}

func joinContents(messages []llama.ChatMessage) string {
	var b strings.Builder
	for _, m := range messages {
		if s, ok := m.Content.(string); ok {
			b.WriteString(s)
		}
	}
	return b.String()
}

func TestTurnStreamMessagesCharacterAttribution(t *testing.T) {
	c, pm := newTestConversation(t)
	meta, err := c.New("Character")
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	client := &fakeChatClient{chunks: []string{"Aria speaks."}}
	turn := &Turn{Conversation: c, Client: client, Options: llama.NewInferenceOptions()}
	reply, err := turn.StreamMessages(context.Background(),
		[]llama.ChatMessage{{Role: "user", Content: "Aria, say something"}},
		nil,
		Attribution{Sender: "Aria", CharacterID: "char_aria"},
		StreamHooks{})
	if err != nil {
		t.Fatalf("StreamMessages: %v", err)
	}
	if reply != "Aria speaks." {
		t.Fatalf("reply = %q", reply)
	}
	entries := c.Session().Entries()
	if len(entries) != 1 || entries[0].CharacterID != "char_aria" || entries[0].Sender != "Aria" {
		t.Fatalf("character attribution = %+v", entries)
	}
	lines := logLines(t, pm, meta.ID)
	if len(lines) != 1 || !strings.Contains(lines[0], "char_aria") {
		t.Fatalf("character attribution did not reach the log: %v", lines)
	}

	loaded, err := pm.LoadChat(meta.ID)
	if err != nil {
		t.Fatalf("LoadChat: %v", err)
	}
	if len(loaded.History) != 1 || loaded.History[0].CharacterID != "char_aria" {
		t.Fatalf("saved history = %+v", loaded.History)
	}
}

func TestTurnCompleteRecordsWithoutALogLine(t *testing.T) {
	c, pm := newTestConversation(t)
	meta, err := c.New("Complete")
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	client := &fakeChatClient{reply: "complete reply"}
	turn := &Turn{Conversation: c, Client: client, Options: llama.NewInferenceOptions()}
	reply, err := turn.Complete(nil, "ping", Attribution{Sender: "AI"})
	if err != nil {
		t.Fatalf("Complete: %v", err)
	}
	if reply != "complete reply" {
		t.Fatalf("reply = %q", reply)
	}
	if got := len(c.Session().Entries()); got != 2 {
		t.Fatalf("session has %d turns, want 2", got)
	}
	if got := len(logLines(t, pm, meta.ID)); got != 0 {
		t.Fatalf("Complete wrote %d log lines; the renderer owns that log", got)
	}
}

func TestTurnStreamErrorLeavesTheSessionWithoutAnAssistantTurn(t *testing.T) {
	c, _ := newTestConversation(t)
	if _, err := c.New("Error"); err != nil {
		t.Fatalf("New: %v", err)
	}
	client := &fakeChatClient{err: errors.New("upstream failed")}
	turn := &Turn{Conversation: c, Client: client, Options: llama.NewInferenceOptions()}
	if _, err := turn.SendStream(context.Background(), "hello", nil, Attribution{Sender: "AI"}, StreamHooks{}); err == nil {
		t.Fatal("stream error was not returned")
	}
	entries := c.Session().Entries()
	if len(entries) != 1 || entries[0].Role != "user" {
		t.Fatalf("failed generation left %+v", entries)
	}
}
