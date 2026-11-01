package main

import (
	"regexp"
	"strings"
)

// thinkBlockRE matches a complete reasoning block in the shapes models actually
// emit:  thinking...</think>, <|think|>...<|/think|>, and the <thinking> variant,
// case-insensitively and across newlines.
var thinkBlockRE = regexp.MustCompile(`(?is)<\|?think(?:ing)?\|?>.*?<\|?/think(?:ing)?\|?>`)

// orphanCloseRE matches a closing reasoning tag with no opening tag before it.
// Models emit these when the opening tag was consumed by the reasoning parser, and
// leaving "</think>" in the visible answer is worse than removing it.
var orphanCloseRE = regexp.MustCompile(`(?is)<\|?/think(?:ing)?\|?>`)

// stripThinkBlocks removes reasoning from a reply before it is stored.
//
// Reasoning is for display only. If it is saved into the message body it becomes
// part of the transcript, gets counted against the context budget, and is fed
// back to the model on the next turn — three things the app must not do. An
// unterminated block means the model was cut off mid-thought, so everything from
// the opening tag on is dropped: a half-finished thought in history is worse
// than losing it.
func stripThinkBlocks(s string) string {
	if s == "" {
		return s
	}
	out := thinkBlockRE.ReplaceAllString(s, "")
	out = orphanCloseRE.ReplaceAllString(out, "")

	lower := strings.ToLower(out)
	for _, tag := range []string{"<|think", "<think"} {
		if i := strings.Index(lower, tag); i >= 0 {
			out = out[:i]
			lower = strings.ToLower(out)
		}
	}
	return strings.TrimSpace(out)
}
