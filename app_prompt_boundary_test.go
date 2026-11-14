package main

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"howl-chat/internal/backend/chat"
	"howl-chat/internal/backend/llama"
	"howl-chat/internal/backend/lorebook"
	entities "howl-chat/internal/backend/types/entities"
)

// recordingInference sits exactly where the HTTP client sits: it receives the
// assembled messages the app would hand to llama-server. It records them in
// memory and answers with a fixed reply, so the real send path can be verified
// without a model, a prompt log on disk, or a second assembly route.
type recordingInference struct {
	mu    sync.Mutex
	last  []llama.ChatMessage
	calls int
	reply string
}

func (r *recordingInference) record(messages []llama.ChatMessage) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.last = append([]llama.ChatMessage(nil), messages...)
	r.calls++
}

func (r *recordingInference) GenerateChat(messages []llama.ChatMessage, options *llama.InferenceOptions) (string, error) {
	r.record(messages)
	return r.reply, nil
}

func (r *recordingInference) GenerateChatStream(messages []llama.ChatMessage, options *llama.InferenceOptions, callback llama.TokenCallback) error {
	r.record(messages)
	callback(r.reply, false)
	callback("", true)
	return nil
}

func (r *recordingInference) GenerateChatStreamWithImage(messages []llama.ChatMessage, imageData []byte, options *llama.InferenceOptions, callback llama.TokenCallback) error {
	r.record(messages)
	callback(r.reply, false)
	callback("", true)
	return nil
}

func (r *recordingInference) CancelRequest() {}

func (r *recordingInference) GetContextSize() int { return 4096 }

func (r *recordingInference) GetLoadedModel() (map[string]interface{}, error) { return nil, nil }

func (r *recordingInference) Health() error { return nil }

func (r *recordingInference) IsLoaded() bool { return true }

func (r *recordingInference) messages() []llama.ChatMessage {
	r.mu.Lock()
	defer r.mu.Unlock()
	return append([]llama.ChatMessage(nil), r.last...)
}

func (r *recordingInference) system() string {
	messages := r.messages()
	if len(messages) == 0 {
		return ""
	}
	content, _ := messages[0].Content.(string)
	return content
}

func boundaryApp(t *testing.T) (*App, *recordingInference) {
	t.Helper()
	app, _ := newLorebookApp(t)
	recorder := &recordingInference{reply: "The hall answers."}
	app.httpClient = recorder
	app.characters = []*chat.CharacterContext{
		{Character: entities.Character{Name: "Aria", ID: "char_aria"}},
	}
	return app, recorder
}

// W13 + W14 + W15 + W16 at the inference boundary: canonical edit wins over the
// chat snapshot, an entry created after the chat save is activated, and deleted,
// hidden and private character lore never reach the model-facing request.
func TestInferenceBoundaryCarriesOnlyCurrentCanonicalLore(t *testing.T) {
	app, recorder := boundaryApp(t)

	if _, err := app.ImportLorebooks([]map[string]interface{}{
		{"id": "canon-1", "title": "Canon", "content": "CANON-ORIGINAL", "triggerPhrases": []string{"sigil"}},
		{"id": "hidden-1", "title": "Hidden", "content": "HIDDEN-TEXT", "triggerPhrases": []string{"sigil"}, "visibility": "hidden"},
		{"id": "deleted-1", "title": "Gone", "content": "DELETED-TEXT", "triggerPhrases": []string{"sigil"}},
	}); err != nil {
		t.Fatalf("import: %v", err)
	}
	if _, err := app.SaveLorebook(map[string]interface{}{
		"id": "canon-1", "title": "Canon", "content": "CANON-EDITED", "triggerPhrases": []string{"sigil"},
	}, true); err != nil {
		t.Fatalf("canonical edit: %v", err)
	}
	if _, err := app.DeleteLorebook("deleted-1"); err != nil {
		t.Fatalf("canonical delete: %v", err)
	}
	activated, err := app.SaveLorebook(map[string]interface{}{
		"title": "Activated", "content": "ACTIVATED-TEXT", "triggerPhrases": []string{"sigil"},
	}, false)
	if err != nil {
		t.Fatalf("create after chat save: %v", err)
	}

	// A chat persisted before those edits: stale copies plus a private entry.
	snapshot := []map[string]interface{}{
		{"id": "canon-1", "title": "Canon", "content": "SNAPSHOT-TEXT", "triggerPhrases": []string{"sigil"}, "enabled": true},
		{"id": "hidden-1", "title": "Hidden", "content": "SNAPSHOT-HIDDEN", "triggerPhrases": []string{"sigil"}, "enabled": true},
		{"id": "deleted-1", "title": "Gone", "content": "SNAPSHOT-DELETED", "triggerPhrases": []string{"sigil"}, "enabled": true},
		{"id": "aria-private", "title": "Aria Secret", "content": "PRIVATE-TEXT", "triggerPhrases": []string{"sigil"}, "enabled": true,
			"scope": "character", "ownerId": "Aria", "visibility": "private"},
	}
	if err := app.SetLorebooks(snapshot); err != nil {
		t.Fatalf("set: %v", err)
	}
	if err := app.UpdateLorebookContext([]string{"canon-1", "hidden-1", "deleted-1", "aria-private", activated.ID}); err != nil {
		t.Fatalf("activate: %v", err)
	}

	if err := app.SendMessageStreamWithImage("I whisper the sigil", ""); err != nil {
		t.Fatalf("send: %v", err)
	}
	if recorder.calls == 0 {
		t.Fatal("the inference boundary was never reached")
	}

	messages := recorder.messages()
	if len(messages) == 0 || messages[0].Role != "system" {
		t.Fatalf("expected a leading system message, got %#v", messages)
	}
	if last := messages[len(messages)-1]; last.Role != "user" {
		t.Fatalf("expected the user turn last, got %#v", last)
	}

	system := recorder.system()
	for _, want := range []string{"CANON-EDITED", "ACTIVATED-TEXT"} {
		if !strings.Contains(system, want) {
			t.Errorf("%s missing from the inference-boundary prompt:\n%s", want, system)
		}
	}
	for _, banned := range []string{"SNAPSHOT-TEXT", "SNAPSHOT-HIDDEN", "SNAPSHOT-DELETED", "HIDDEN-TEXT", "DELETED-TEXT", "PRIVATE-TEXT"} {
		if strings.Contains(system, banned) {
			t.Errorf("%s reached the inference boundary:\n%s", banned, system)
		}
	}

	// The reply still lands in history through the normal append path.
	found := false
	for _, msg := range app.chatMessages {
		if content, _ := msg["content"].(string); strings.Contains(content, "The hall answers.") {
			found = true
		}
	}
	if !found {
		t.Error("recorder reply did not reach chat history")
	}
}

// An unavailable root keeps its last-known entry out of the request entirely.
func TestInferenceBoundaryExcludesStaleLore(t *testing.T) {
	app, recorder := boundaryApp(t)
	if _, err := app.ImportLorebooks([]map[string]interface{}{
		{"id": "canon-1", "title": "Canon", "content": "OFFLINE-CANON", "triggerPhrases": []string{"sigil"}},
	}); err != nil {
		t.Fatalf("import: %v", err)
	}
	list, err := app.ListLorebooks()
	if err != nil || len(list.Items) == 0 {
		t.Fatalf("list: %v %+v", err, list)
	}
	rootPath := filepath.Dir(list.Items[0].Path)
	offline := rootPath + "-offline"
	if err := os.Rename(rootPath, offline); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Rename(offline, rootPath) }()

	if err := app.SetLorebooks([]map[string]interface{}{
		{"id": "canon-1", "title": "Canon", "content": "SNAPSHOT-OFFLINE", "triggerPhrases": []string{"sigil"}, "enabled": true},
	}); err != nil {
		t.Fatalf("set: %v", err)
	}
	if len(app.lorebooks) != 1 || !app.lorebooks[0].Stale {
		t.Fatalf("pool = %+v, want one stale entry", app.lorebooks)
	}
	if err := app.SendMessageStreamWithImage("I whisper the sigil", ""); err != nil {
		t.Fatalf("send: %v", err)
	}
	system := recorder.system()
	for _, banned := range []string{"OFFLINE-CANON", "SNAPSHOT-OFFLINE"} {
		if strings.Contains(system, banned) {
			t.Errorf("%s reached the inference boundary while the root was offline:\n%s", banned, system)
		}
	}
}

// The character-turn path sends a frontend-composed prompt verbatim. Privacy
// therefore depends on the composer using the shared backend policy; this test
// pins both the unchanged boundary shape and the policy it must consume.
func TestCharacterResponseBoundaryUsesSharedLorePolicy(t *testing.T) {
	app, recorder := boundaryApp(t)
	app.lorebooks = []lorebook.Entry{
		privacyEntry("aria-private", "Aria Secret", "ARIA-SECRET", lorebook.ScopeCharacter, "char_aria", lorebook.VisibilityPrivate),
		privacyEntry("bram-private", "Bram Secret", "BRAM-SECRET", lorebook.ScopeCharacter, "char_bram", lorebook.VisibilityPrivate),
		privacyEntry("unowned-private", "Unowned Secret", "UNOWNED-SECRET", lorebook.ScopeCharacter, "", lorebook.VisibilityPrivate),
		privacyEntry("hidden-world", "Hidden World", "HIDDEN-SECRET", lorebook.ScopeWorld, "", lorebook.VisibilityHidden),
	}

	block := app.BuildCharacterLoreBlock("char_aria", "I whisper the sigil")
	if !strings.Contains(block, "ARIA-SECRET") {
		t.Fatalf("Aria's own lore missing from the character block:\n%s", block)
	}
	for _, leaked := range []string{"BRAM-SECRET", "UNOWNED-SECRET", "HIDDEN-SECRET"} {
		if strings.Contains(block, leaked) {
			t.Fatalf("%s leaked into Aria's character block:\n%s", leaked, block)
		}
	}

	composed := block + "\n\nUser: I whisper the sigil\nRespond as Aria only."
	if err := app.GenerateCharacterResponse(composed, "char_aria"); err != nil {
		t.Fatalf("generate: %v", err)
	}
	messages := recorder.messages()
	if len(messages) != 1 {
		t.Fatalf("character turn carried %d messages, want exactly the composed prompt: %#v", len(messages), messages)
	}
	content, _ := messages[0].Content.(string)
	if content != composed {
		t.Fatalf("composed prompt was altered at the boundary")
	}
}

func TestCharacterResponseBoundarySendsFrontendPromptOnly(t *testing.T) {
	app, recorder := boundaryApp(t)
	app.lorebooks = []lorebook.Entry{
		privacyEntry("aria-private", "Aria Secret", "PRIVATE-TEXT", lorebook.ScopeCharacter, "char_aria", lorebook.VisibilityPrivate),
	}

	if err := app.GenerateCharacterResponse("FRONTEND-COMPOSED PROMPT", "char_aria"); err != nil {
		t.Fatalf("generate: %v", err)
	}
	messages := recorder.messages()
	if len(messages) != 1 {
		t.Fatalf("character turn carried %d messages, want exactly the frontend prompt: %#v", len(messages), messages)
	}
	if messages[0].Role != "user" {
		t.Fatalf("role = %q", messages[0].Role)
	}
	content, _ := messages[0].Content.(string)
	if content != "FRONTEND-COMPOSED PROMPT" {
		t.Fatalf("content = %q, want the frontend prompt unchanged", content)
	}
	if strings.Contains(content, "PRIVATE-TEXT") {
		t.Fatal("backend lore was appended to a character turn")
	}
}
