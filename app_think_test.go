package main

import "testing"

// Reasoning must never survive into a stored message: if it does it enters the
// transcript, counts against the context budget, and is fed back next turn.
func TestStripThinkBlocks(t *testing.T) {
	open := "<" + "think" + ">"
	closing := "<" + "/" + "think" + ">"
	pipeOpen := "<|" + "think" + "|>"
	pipeClose := "<|/" + "think" + "|>"

	cases := []struct {
		name string
		in   string
		want string
	}{
		{"plain answer untouched", "The door opens.", "The door opens."},
		{"empty", "", ""},
		{"complete block removed", "answer " + open + "reasoning here" + closing + " more", "answer  more"},
		{"block only", open + "just thinking" + closing, ""},
		{"multiline block", "Before.\n" + open + "\nline one\nline two\n" + closing + "\nAfter.", "Before.\n\nAfter."},
		{"pipe tags", "a " + pipeOpen + "hmm" + pipeClose + " b", "a  b"},
		{"unterminated block loses tail", "hello " + open + "still going and going", "hello"},
	}

	for _, tc := range cases {
		if got := stripThinkBlocks(tc.in); got != tc.want {
			t.Errorf("%s: got %q, want %q", tc.name, got, tc.want)
		}
	}
}
