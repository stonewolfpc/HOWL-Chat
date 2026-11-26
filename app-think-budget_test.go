package main

import (
	"strings"
	"testing"
)

// The reasoning mode is one setting read by two paths: the llama-server flags
// and the per-request token budget. When the user has not chosen a mode, the
// server reasons with the default budget, so the request must pay for that
// budget too. When the two disagreed, the whole reply cap was spent inside an
// unclosed thought and the assistant bubble stayed empty.
func TestAbsentThinkModePaysForTheServersReasoningBudget(t *testing.T) {
	settings := map[string]interface{}{}

	args := strings.Join(buildLlamaServerArgs("model.gguf", "", settings, false, llamaCapsForTests()), " ")
	if !strings.Contains(args, "--reasoning on") {
		t.Fatalf("the default mode must ask the server to reason: %q", args)
	}
	if !strings.Contains(args, "--reasoning-budget 2048") {
		t.Fatalf("the server runs a 2048-token thought budget; the test expects it: %q", args)
	}
	if got := thinkHeadroom(settings); got != 2048 {
		t.Fatalf("thought headroom = %d, want the server's 2048-token budget", got)
	}

	// Reasoning off stays off: no server reasoning and no reserved allowance.
	off := map[string]interface{}{"think_mode": "off"}
	if got := thinkHeadroom(off); got != 0 {
		t.Fatalf("reasoning off still reserved %d tokens", got)
	}
	if !strings.Contains(strings.Join(buildLlamaServerArgs("model.gguf", "", off, false, llamaCapsForTests()), " "), "--reasoning off") {
		t.Fatal("reasoning off did not reach the server flags")
	}

	// A budget the user chose is the budget both paths use.
	custom := map[string]interface{}{"think_mode": "on", "think_budget": float64(512)}
	if got := thinkHeadroom(custom); got != 512 {
		t.Fatalf("thought headroom = %d, want the configured 512", got)
	}
}
