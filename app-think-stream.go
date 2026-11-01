package main

// app-think-stream.go — splits inline reasoning out of a streaming reply.
//
// llama-server only delivers thought and answer separately when the model's
// template and --reasoning-format agree. Some models (Qwen3.5 in this build)
// emit " thinking...</think>" inside the content stream instead, so the thought
// and its tags land in the visible message, in the saved history, and in the
// next prompt. The splitter routes the thought to the reasoning callback and
// only the answer onward, and it copes with a tag split across two deltas.

import "strings"

var thinkOpenTags = []string{"\x3cthink\x3e", "\x3cthinking\x3e", "\x3c|think|\x3e", "\x3c|thinking|\x3e"}
var thinkCloseTags = []string{"\x3c/think\x3e", "\x3c/thinking\x3e", "\x3c|/think|\x3e", "\x3c|/thinking|\x3e"}

type thinkSplitter struct {
	buf     string
	inThink bool
	onThink func(string)
	onText  func(string)
}

func newThinkSplitter(onThink, onText func(string)) *thinkSplitter {
	return &thinkSplitter{onThink: onThink, onText: onText}
}

func longestTag(tags []string) int {
	longest := 0
	for _, t := range tags {
		if len(t) > longest {
			longest = len(t)
		}
	}
	return longest
}

// earliestTag reports the first tag of any kind and its offset, or -1.
func earliestTag(s string, tags []string) (int, string) {
	best, match := -1, ""
	for _, t := range tags {
		if i := strings.Index(s, t); i >= 0 && (best < 0 || i < best) {
			best, match = i, t
		}
	}
	return best, match
}

// Write consumes one streamed delta, emitting thought and answer text as they
// become unambiguous.
func (s *thinkSplitter) Write(delta string) {
	if delta == "" {
		return
	}
	s.buf += delta

	for {
		if s.inThink {
			i, tag := earliestTag(s.buf, thinkCloseTags)
			if i < 0 {
				// Hold back a possible partial closing tag.
				keep := longestTag(thinkCloseTags) - 1
				if len(s.buf) > keep {
					if s.onThink != nil {
						s.onThink(s.buf[:len(s.buf)-keep])
					}
					s.buf = s.buf[len(s.buf)-keep:]
				}
				return
			}
			if i > 0 && s.onThink != nil {
				s.onThink(s.buf[:i])
			}
			s.buf = s.buf[i+len(tag):]
			s.inThink = false
			continue
		}

		i, tag := earliestTag(s.buf, thinkOpenTags)
		if i < 0 {
			keep := longestTag(thinkOpenTags) - 1
			if len(s.buf) > keep {
				if s.onText != nil {
					s.onText(s.buf[:len(s.buf)-keep])
				}
				s.buf = s.buf[len(s.buf)-keep:]
			}
			return
		}
		if i > 0 && s.onText != nil {
			s.onText(s.buf[:i])
		}
		s.buf = s.buf[i+len(tag):]
		s.inThink = true
	}
}

// Flush emits whatever is still buffered once the stream ends.
func (s *thinkSplitter) Flush() {
	if s.buf == "" {
		return
	}
	if s.inThink {
		if s.onThink != nil {
			s.onThink(s.buf)
		}
	} else if s.onText != nil {
		s.onText(s.buf)
	}
	s.buf = ""
}
