package chat

import (
	"os"
	"path/filepath"
	"testing"
)

// These tests treat the lifecycle log as the authority and the directory as the
// payload: present, missing and orphaned are all read-only conclusions.

func newMembershipFixture(t *testing.T) (root, chatsDir string) {
	t.Helper()
	root = t.TempDir()
	chatsDir = filepath.Join(root, "chats")
	if err := os.MkdirAll(chatsDir, 0o755); err != nil {
		t.Fatalf("create chats dir: %v", err)
	}
	return root, chatsDir
}

func writeDir(t *testing.T, chatsDir, id, meta string) {
	t.Helper()
	dir := filepath.Join(chatsDir, id)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("create dir %s: %v", id, err)
	}
	if meta != "" {
		if err := os.WriteFile(filepath.Join(dir, "meta.json"), []byte(meta), 0o644); err != nil {
			t.Fatalf("write meta %s: %v", id, err)
		}
	}
}

func stateOf(convs []Conversation, id string) ConversationState {
	for _, conv := range convs {
		if conv.ID == id {
			return conv.State
		}
	}
	return ""
}

func TestMembershipReportsPresentMissingAndOrphaned(t *testing.T) {
	root, chatsDir := newMembershipFixture(t)
	m, err := LoadMembership(root)
	if err != nil {
		t.Fatalf("load membership: %v", err)
	}

	// Present: created, directory there.
	if err := m.Create("2026-11-22T10-00-00Z", "kept"); err != nil {
		t.Fatalf("create member: %v", err)
	}
	writeDir(t, chatsDir, "2026-11-22T10-00-00Z", `{"name":"kept"}`)

	// Missing: created, directory gone.
	if err := m.Create("2026-11-22T10-01-00Z", "vanished"); err != nil {
		t.Fatalf("create member: %v", err)
	}

	// Orphaned: a directory nobody claimed.
	writeDir(t, chatsDir, "2026-11-22T10-02-00Z", `{"name":"stranger"}`)

	convs, err := m.Reconcile(chatsDir)
	if err != nil {
		t.Fatalf("reconcile: %v", err)
	}

	if got := stateOf(convs, "2026-11-22T10-00-00Z"); got != StatePresent {
		t.Fatalf("kept conversation state = %q, want %q", got, StatePresent)
	}
	if got := stateOf(convs, "2026-11-22T10-01-00Z"); got != StateMissing {
		t.Fatalf("vanished conversation state = %q, want %q", got, StateMissing)
	}
	if got := stateOf(convs, "2026-11-22T10-02-00Z"); got != StateOrphaned {
		t.Fatalf("stranger directory state = %q, want %q", got, StateOrphaned)
	}
	for _, conv := range convs {
		if conv.ID != "2026-11-22T10-02-00Z" {
			continue
		}
		if conv.Name != "stranger" {
			t.Fatalf("orphan name = %q, want the label the folder still carries", conv.Name)
		}
	}
}

func TestMembershipDeletionSurvivesRestartAndRefusesResurrection(t *testing.T) {
	root, chatsDir := newMembershipFixture(t)
	m, err := LoadMembership(root)
	if err != nil {
		t.Fatalf("load membership: %v", err)
	}

	id := "2026-11-22T11-00-00Z"
	if err := m.Create(id, "doomed"); err != nil {
		t.Fatalf("create: %v", err)
	}
	writeDir(t, chatsDir, id, `{"name":"doomed"}`)
	if err := m.Delete(id, filepath.Join(chatsDir, id)); err != nil {
		t.Fatalf("delete: %v", err)
	}

	// Restart: the log is the authority, not memory.
	reloaded, err := LoadMembership(root)
	if err != nil {
		t.Fatalf("reload membership: %v", err)
	}
	if !reloaded.Deleted(id) {
		t.Fatalf("deletion did not survive a reload")
	}
	if err := reloaded.Create(id, "resurrected"); err == nil {
		t.Fatalf("creating a deleted identifier was accepted")
	}

	// A leftover directory after a failed removal is orphaned, never present.
	writeDir(t, chatsDir, id, `{"name":"doomed"}`)
	convs, err := reloaded.Reconcile(chatsDir)
	if err != nil {
		t.Fatalf("reconcile: %v", err)
	}
	if got := stateOf(convs, id); got != StateOrphaned {
		t.Fatalf("leftover of a deleted conversation state = %q, want %q", got, StateOrphaned)
	}
}

func TestMembershipSurvivesATornTail(t *testing.T) {
	root, chatsDir := newMembershipFixture(t)
	m, err := LoadMembership(root)
	if err != nil {
		t.Fatalf("load membership: %v", err)
	}
	if err := m.Create("2026-11-22T12-00-00Z", "complete"); err != nil {
		t.Fatalf("create: %v", err)
	}

	// Simulate a crash mid-append: a record with no newline and no closing brace.
	logPath := filepath.Join(root, lifecycleLogName)
	f, err := os.OpenFile(logPath, os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		t.Fatalf("open log: %v", err)
	}
	if _, err := f.WriteString(`{"event":"created","id":"2026-11-22T12-00-01Z"`); err != nil {
		t.Fatalf("write torn record: %v", err)
	}
	f.Close()

	reloaded, err := LoadMembership(root)
	if err != nil {
		t.Fatalf("reload after torn tail: %v", err)
	}
	if _, ok := reloaded.Member("2026-11-22T12-00-00Z"); !ok {
		t.Fatalf("the complete record was lost")
	}
	if _, ok := reloaded.Member("2026-11-22T12-00-01Z"); ok {
		t.Fatalf("an incomplete record was adopted as a member")
	}
	_ = chatsDir
}

func TestMembershipAbsorbsLegacyDeletionLogOnce(t *testing.T) {
	root, chatsDir := newMembershipFixture(t)
	id := "2026-11-22T13-00-00Z"
	legacy := filepath.Join(root, legacyDeletedLog)
	line := `{"id":"` + id + `","at":"2026-11-22T13:01:00Z","directory":"` + filepath.ToSlash(filepath.Join(chatsDir, id)) + `"}` + "\n"
	if err := os.WriteFile(legacy, []byte(line), 0o644); err != nil {
		t.Fatalf("write legacy log: %v", err)
	}

	m, err := LoadMembership(root)
	if err != nil {
		t.Fatalf("load membership: %v", err)
	}
	if !m.Deleted(id) {
		t.Fatalf("legacy deletion was not absorbed")
	}
	if _, err := os.Stat(filepath.Join(root, lifecycleLogName)); err != nil {
		t.Fatalf("lifecycle log was not created by the migration: %v", err)
	}

	// The migration does not run twice: a second load must not duplicate records
	// or lose the deletion.
	again, err := LoadMembership(root)
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	if !again.Deleted(id) {
		t.Fatalf("deletion lost on reload")
	}

	data, err := os.ReadFile(filepath.Join(root, lifecycleLogName))
	if err != nil {
		t.Fatalf("read lifecycle log: %v", err)
	}
	if count := countLines(string(data)); count != 1 {
		t.Fatalf("lifecycle log has %d records, want exactly the migrated one", count)
	}
}

func countLines(s string) int {
	n := 0
	for _, ch := range s {
		if ch == '\n' {
			n++
		}
	}
	return n
}

func TestMembershipBootstrapsAnExistingInstallOnce(t *testing.T) {
	root, chatsDir := newMembershipFixture(t)
	// An install that predates this authority: a directory with metadata, and one
	// directory that does not identify itself.
	writeDir(t, chatsDir, "2026-11-22T14-00-00Z", `{"name":"old chat"}`)
	writeDir(t, chatsDir, "not-a-chat", "")

	m, err := LoadMembership(root)
	if err != nil {
		t.Fatalf("load membership: %v", err)
	}
	convs, err := m.Reconcile(chatsDir)
	if err != nil {
		t.Fatalf("reconcile: %v", err)
	}
	if got := stateOf(convs, "2026-11-22T14-00-00Z"); got != StatePresent {
		t.Fatalf("existing conversation state = %q, want %q", got, StatePresent)
	}
	if got := stateOf(convs, "not-a-chat"); got != StateOrphaned {
		t.Fatalf("unidentifiable directory state = %q, want %q", got, StateOrphaned)
	}

	// A directory created after the log exists is never claimed by a later load.
	writeDir(t, chatsDir, "2026-11-22T14-05-00Z", `{"name":"late arrival"}`)
	again, err := LoadMembership(root)
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	convs, err = again.Reconcile(chatsDir)
	if err != nil {
		t.Fatalf("reconcile after reload: %v", err)
	}
	if got := stateOf(convs, "2026-11-22T14-05-00Z"); got != StateOrphaned {
		t.Fatalf("a directory added after bootstrap has state %q, want %q", got, StateOrphaned)
	}
}
