package durable

import (
	"os"
	"testing"
	"time"
)

// TestFailedCommitPreservesPreImageChain audits the commit protocol's failure
// path: a write that cannot replace the live document must leave the recovery
// chain exactly as it was.
//
// The scenario that matters is a reader holding the document while a writer
// tries to commit. Afterwards the document is still readable, but if it is lost
// the pre-image chain must still let recovery find a real generation.
func TestFailedCommitPreservesPreImageChain(t *testing.T) {
	s := newStore(t, Options{
		LockTimeout:     time.Second,
		RenameAttempts:  1,
		RenameWait:      10 * time.Millisecond,
		PreImages:       2,
		SupportedSchema: 1,
	})
	for _, doc := range []string{"d1", "d2", "d3"} {
		mustSave(t, s, doc)
	}
	// Chain as committed: v1 = d2, v2 = d1.
	if got := readFileOrEmpty(t, s.preImagePath(1)); got != "d2" {
		t.Fatalf("precondition: v1 = %q, want d2", got)
	}

	reader, err := os.Open(s.Path())
	if err != nil {
		t.Fatalf("open reader: %v", err)
	}
	res := s.Save([]byte("d4"), 1)
	reader.Close()
	if res.Outcome != NotSaved || res.Reason != "rename_denied_by_open_reader" {
		t.Fatalf("failed commit = %+v, want NotSaved(rename_denied_by_open_reader)", res)
	}

	// The failed commit must not have rewritten the retained generations.
	if got := readFileOrEmpty(t, s.preImagePath(1)); got != "d2" {
		t.Fatalf("after failed commit, v1 = %q, want the untouched d2", got)
	}
	if got := readFileOrEmpty(t, s.preImagePath(2)); got != "d1" {
		t.Fatalf("after failed commit, v2 = %q, want the untouched d1", got)
	}

	// Recovery must still work on the surviving document.
	if got, lres := s.Load(); string(got) != "d3" || !lres.OK() {
		t.Fatalf("after failed commit, Load = %q %+v, want d3/Saved", got, lres)
	}

	// And if the live document is later unusable, the chain must still yield a
	// real generation rather than reporting corruption.
	if err := os.WriteFile(s.Path(), []byte("torn"), 0o644); err != nil {
		t.Fatalf("corrupt live: %v", err)
	}
	got, lres := s.Load()
	if lres.Outcome != Recovered {
		t.Fatalf("after failed commit + corruption: %q %+v, want Recovered from a retained generation", got, lres)
	}
	if string(got) != "d2" {
		t.Fatalf("recovered %q, want d2 (the newest retained generation)", got)
	}
}

// TestFailedCommitLeavesNoStagedArtifacts audits that a failed commit cleans up
// anything it wrote while preparing, so the next write is not confused by it.
func TestFailedCommitLeavesNoStagedArtifacts(t *testing.T) {
	s := newStore(t, Options{
		LockTimeout:     time.Second,
		RenameAttempts:  1,
		RenameWait:      10 * time.Millisecond,
		PreImages:       2,
		SupportedSchema: 1,
	})
	mustSave(t, s, "d1")

	reader, err := os.Open(s.Path())
	if err != nil {
		t.Fatalf("open reader: %v", err)
	}
	res := s.Save([]byte("d2"), 1)
	reader.Close()
	if res.Outcome != NotSaved {
		t.Fatalf("failed commit = %+v", res)
	}

	if _, err := os.Stat(s.preImagePath(1) + ".staged"); err == nil {
		t.Fatal("failed commit left a staged pre-image behind")
	}
	mustSave(t, s, "d2-after")
}
