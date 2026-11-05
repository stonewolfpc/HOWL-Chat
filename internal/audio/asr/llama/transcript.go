package llama

import (
	"regexp"
	"strings"
)

// transcriptPrefix matches the wrapper llama-server returns for Qwen3-ASR:
// "language English<asr_text>...". The text after the marker is the transcript.
var transcriptPrefix = regexp.MustCompile(`(?is)^\s*language\s+([a-z-]+)\s*<asr_text>`)

// CleanTranscript removes the language/asr_text wrapper and reports the
// detected language. Text without a wrapper is returned unchanged.
func CleanTranscript(raw string) (text, language string) {
	raw = strings.TrimSpace(raw)
	if match := transcriptPrefix.FindStringSubmatch(raw); match != nil {
		return strings.TrimSpace(raw[len(match[0]):]), strings.ToLower(match[1])
	}
	if idx := strings.Index(strings.ToLower(raw), "<asr_text>"); idx >= 0 {
		return strings.TrimSpace(raw[idx+len("<asr_text>"):]), ""
	}
	return raw, ""
}
