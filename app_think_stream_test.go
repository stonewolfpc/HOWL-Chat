package main

import "testing"

// Tag literals use Go escapes so the angle brackets are built by the
// compiler instead of travelling through anything that strips tag-like text.
const (
	tOpen  = "\x3cthink\x3e"
	tClose = "\x3c/think\x3e"
)

type splitResult struct {
	think string
	text  string
}

func runSplitter(deltas []string) splitResult {
	var r splitResult
	s := newThinkSplitter(
		func(t string) { r.think += t },
		func(t string) { r.text += t },
	)
	for _, d := range deltas {
		s.Write(d)
	}
	s.Flush()
	return r
}

func TestThinkSplitter(t *testing.T) {
	cases := []struct {
		name   string
		deltas []string
		think  string
		text   string
	}{
		{
			name:   "inline block",
			deltas: []string{"Hi. " + tOpen + " plan the reply. " + tClose + " Hello!"},
			think:  " plan the reply. ",
			text:   "Hi.  Hello!",
		},
		{
			name:   "tag split across deltas",
			deltas: []string{"Hi. \x3cth", "ink\x3ethink", "ing\x3c/thi", "nk\x3e there"},
			think:  "thinking",
			text:   "Hi.  there",
		},
		{
			name:   "no tags",
			deltas: []string{"plain ", "answer"},
			think:  "",
			text:   "plain answer",
		},
		{
			name:   "unterminated thought",
			deltas: []string{"Answer. " + tOpen + " still thinking"},
			think:  " still thinking",
			text:   "Answer. ",
		},
		{
			name:   "no reasoning at all",
			deltas: []string{"Name", "'s Yara."},
			think:  "",
			text:   "Name's Yara.",
		},
		{
			name:   "tags arriving one fragment at a time",
			deltas: []string{"Hello", " \x3ct", "hink", "\x3ex\x3c/think", "\x3e world"},
			think:  "x",
			text:   "Hello  world",
		},
	}

	for _, tc := range cases {
		got := runSplitter(tc.deltas)
		if got.think != tc.think {
			t.Errorf("%s: think = %q, want %q", tc.name, got.think, tc.think)
		}
		if got.text != tc.text {
			t.Errorf("%s: text = %q, want %q", tc.name, got.text, tc.text)
		}
		if hasThinkTag(got.text) {
			t.Errorf("%s: answer still carries a tag: %q", tc.name, got.text)
		}
	}
}

func hasThinkTag(s string) bool {
	for _, tag := range append(append([]string{}, thinkOpenTags...), thinkCloseTags...) {
		if tag == "" {
			continue
		}
		if stringsContains(s, tag) {
			return true
		}
	}
	return false
}

// stringsContains keeps the test free of extra imports.
func stringsContains(s, sub string) bool {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
