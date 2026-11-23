package chat

import (
	entities "howl-chat/internal/backend/types/entities"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func newTestConversation(t *testing.T) (*ActiveConversation, *PersistenceManager) {
	t.Helper()
	pm, err := NewPersistenceManagerWithDir(t.TempDir())
	if err != nil {
		t.Fatalf("NewPersistenceManagerWithDir: %v", err)
	}
	return NewActiveConversation(pm), pm
}

func logLines(t *testing.T, pm *PersistenceManager, id string) []string {
	t.Helper()
	dir, err := pm.ChatDirectory(id)
	if err != nil {
		t.Fatalf("ChatDirectory: %v", err)
	}
	data, err := os.ReadFile(filepath.Join(dir, "history.ndjson"))
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		t.Fatalf("read log: %v", err)
	}
	lines := []string{}
	for _, line := range strings.Split(string(data), "\n") {
		if strings.TrimSpace(line) != "" {
			lines = append(lines, line)
		}
	}
	return lines
}

func TestActiveConversationAppendTurnWritesLogAndRestores(t *testing.T) {
	c, pm := newTestConversation(t)
	meta, err := c.New("Log")
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if c.Session().Len() != 0 {
		t.Fatalf("new conversation kept %d turns", c.Session().Len())
	}
	if err := c.AppendTurn(SessionEntry{Role: "user", Content: "hello", Sender: "User"}); err != nil {
		t.Fatalf("AppendTurn: %v", err)
	}
	if err := c.AppendTurn(SessionEntry{Role: "assistant", Content: "hi there", Sender: "AI"}); err != nil {
		t.Fatalf("AppendTurn: %v", err)
	}
	if got := len(logLines(t, pm, meta.ID)); got != 2 {
		t.Fatalf("log has %d lines, want 2", got)
	}
	if err := c.AppendTurn(SessionEntry{Role: "user", Content: "   "}); err != nil {
		t.Fatalf("empty turn: %v", err)
	}
	if got := len(logLines(t, pm, meta.ID)); got != 2 {
		t.Fatalf("blank turn reached the log: %d lines", got)
	}

	// A fresh conversation over the same collection restores the same turns.
	restarted := NewActiveConversation(pm)
	if _, err := restarted.Load(meta.ID); err != nil {
		t.Fatalf("Load: %v", err)
	}
	entries := restarted.Session().Entries()
	if len(entries) != 2 || entries[1].Content != "hi there" {
		t.Fatalf("restored turns = %+v", entries)
	}
}

func TestActiveConversationReplaceHistoryRewritesTheLog(t *testing.T) {
	c, pm := newTestConversation(t)
	meta, err := c.New("Rewrite")
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	for _, text := range []string{"one", "two", "three"} {
		if err := c.AppendTurn(SessionEntry{Role: "user", Content: text}); err != nil {
			t.Fatalf("AppendTurn: %v", err)
		}
	}

	// The delete path: drop the middle turn and rewrite.
	if err := c.ReplaceHistory([]SessionEntry{
		{Role: "user", Content: "one"},
		{Role: "user", Content: "three"},
	}); err != nil {
		t.Fatalf("ReplaceHistory: %v", err)
	}
	lines := logLines(t, pm, meta.ID)
	if len(lines) != 2 {
		t.Fatalf("log has %d lines after rewrite, want 2", len(lines))
	}
	for _, line := range lines {
		if strings.Contains(line, `"two"`) {
			t.Fatalf("deleted turn survived the rewrite: %s", line)
		}
	}
	loaded, err := pm.LoadChat(meta.ID)
	if err != nil {
		t.Fatalf("LoadChat: %v", err)
	}
	if len(loaded.History) != 2 {
		t.Fatalf("reloaded %d turns, want 2", len(loaded.History))
	}
}

func TestActiveConversationClearHistoryTruncatesTheLog(t *testing.T) {
	c, pm := newTestConversation(t)
	meta, err := c.New("Clear")
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if err := c.AppendTurn(SessionEntry{Role: "user", Content: "gone"}); err != nil {
		t.Fatalf("AppendTurn: %v", err)
	}
	if err := c.ClearHistory(); err != nil {
		t.Fatalf("ClearHistory: %v", err)
	}
	if c.Session().Len() != 0 {
		t.Fatalf("session kept %d turns", c.Session().Len())
	}
	if got := len(logLines(t, pm, meta.ID)); got != 0 {
		t.Fatalf("log has %d lines after clear, want 0", got)
	}
	if loaded, err := pm.LoadChat(meta.ID); err != nil || len(loaded.History) != 0 {
		t.Fatalf("cleared conversation reloaded %d turns (err %v)", len(loaded.History), err)
	}
}

func TestActiveConversationSaveRecomputesMembershipFromTheSession(t *testing.T) {
	c, pm := newTestConversation(t)
	meta, err := c.New("Membership")
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	c.Session().SetWorld(&WorldContext{World: entities.World{ID: "world-1", Name: "Ashfall"}})
	c.Session().SetScenario(&ScenarioContext{Scenario: entities.Scenario{ID: "scenario-1", Name: "The Siege"}})
	c.Session().SetCharacters([]*CharacterContext{
		{Character: entities.Character{ID: "char_aria", Name: "Aria"}},
		{Character: entities.Character{ID: "char_bram", Name: "Bram"}},
	})
	c.Session().UpsertCharacter(&CharacterContext{Character: entities.Character{ID: "char_aria", Name: "Aria"}})
	if err := c.Save(); err != nil {
		t.Fatalf("Save: %v", err)
	}

	loaded, err := pm.LoadChat(meta.ID)
	if err != nil {
		t.Fatalf("LoadChat: %v", err)
	}
	if loaded.Meta.WorldID != "world-1" || loaded.Meta.ScenarioID != "scenario-1" {
		t.Fatalf("membership ids = %q / %q", loaded.Meta.WorldID, loaded.Meta.ScenarioID)
	}
	if len(loaded.Meta.ActiveCharacters) != 2 || loaded.Meta.ActiveCharacters[0] != "char_aria" {
		t.Fatalf("active characters = %v", loaded.Meta.ActiveCharacters)
	}
	if loaded.World == nil || loaded.World.ID != "world-1" {
		t.Fatalf("world snapshot lost: %+v", loaded.World)
	}
}

func TestActiveConversationDeleteActiveResetsTheSession(t *testing.T) {
	c, _ := newTestConversation(t)
	meta, err := c.New("Delete")
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if err := c.AppendTurn(SessionEntry{Role: "user", Content: "still here"}); err != nil {
		t.Fatalf("AppendTurn: %v", err)
	}
	if err := c.Delete(meta.ID); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if c.ActiveID() != "" {
		t.Fatalf("deleted conversation stayed active: %q", c.ActiveID())
	}
	if c.Session().Len() != 0 {
		t.Fatalf("deleted conversation kept %d turns in memory", c.Session().Len())
	}
}

func TestActiveConversationRefusesForeignRecord(t *testing.T) {
	c, _ := newTestConversation(t)
	err := c.SaveRecord(&ChatData{Meta: ChatMeta{ID: "not-in-this-collection"}})
	if err == nil {
		t.Fatal("SaveRecord accepted a conversation the collection does not claim")
	}
	if _, ok := err.(*Refusal); !ok {
		t.Fatalf("refusal is not a Refusal: %T %v", err, err)
	}
}

func TestActiveConversationFinalSaveCommitsTurns(t *testing.T) {
	c, pm := newTestConversation(t)
	meta, err := c.New("Final")
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if err := c.AppendTurn(SessionEntry{Role: "user", Content: "final turn"}); err != nil {
		t.Fatalf("AppendTurn: %v", err)
	}
	count, err := c.FinalSave()
	if err != nil {
		t.Fatalf("FinalSave: %v", err)
	}
	if count != 1 {
		t.Fatalf("FinalSave committed %d entries, want 1", count)
	}
	loaded, err := pm.LoadChat(meta.ID)
	if err != nil {
		t.Fatalf("LoadChat: %v", err)
	}
	if len(loaded.History) != 1 || loaded.History[0].Text != "final turn" {
		t.Fatalf("committed history = %+v", loaded.History)
	}
}

func TestActiveConversationUnboundReportsNoCollection(t *testing.T) {
	c := NewActiveConversation(nil)
	if c.Bound() {
		t.Fatal("unbound conversation reported a collection")
	}
	if _, err := c.List(); err == nil {
		t.Fatal("List without a collection did not report an error")
	}
	if err := c.AppendTurn(SessionEntry{Role: "user", Content: "memory only"}); err != nil {
		t.Fatalf("in-memory AppendTurn: %v", err)
	}
	if c.Session().Len() != 1 {
		t.Fatalf("in-memory turn was dropped")
	}
}
