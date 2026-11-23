package durable

import (
	"os"
	"testing"
)

// TestRecoveredIsDistinguishableFromCurrentRead pins the acceptance condition on
// reader consistency: a recovered generation must be observable to the caller and
// must not be mistakable for an ordinary current read.
func TestRecoveredIsDistinguishableFromCurrentRead(t *testing.T) {
	s := newStore(t, DefaultOptions())
	mustSave(t, s, "d1")
	mustSave(t, s, "d2")

	_, current := s.Load()
	if !current.Current() {
		t.Fatalf("a clean read reported Current()=false: %+v", current)
	}
	if current.Outcome != Saved {
		t.Fatalf("clean read outcome = %v", current.Outcome)
	}

	if err := os.WriteFile(s.Path(), []byte("torn"), 0o644); err != nil {
		t.Fatalf("corrupt live: %v", err)
	}
	got, recovered := s.Load()

	if recovered.Outcome != Recovered {
		t.Fatalf("outcome = %v, want Recovered", recovered.Outcome)
	}
	// Usable, but explicitly not current.
	if !recovered.OK() {
		t.Fatalf("recovered read reported OK()=false")
	}
	if recovered.Current() {
		t.Fatal("a recovered generation reported Current()=true; it can be mistaken for the current document")
	}
	// The caller can see why without inferring it.
	if recovered.Reason == "" || recovered.From == "" {
		t.Fatalf("recovered read did not explain itself: %+v", recovered)
	}
	if string(got) != "d1" {
		t.Fatalf("recovered %q, want d1", got)
	}
}
