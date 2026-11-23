package main

import "howl-chat/internal/backend/chat"

// stripThinkBlocks removes reasoning from a reply before it is stored. The rule
// lives in the chat domain (chat.StripThinkBlocks) so the text the domain appends
// and the text the host displays cannot drift apart. The incremental splitter
// (app-think-stream.go) stays here because splitting a live stream is a display
// concern.
func stripThinkBlocks(s string) string {
	return chat.StripThinkBlocks(s)
}
