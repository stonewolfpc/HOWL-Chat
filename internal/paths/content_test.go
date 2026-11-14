package paths

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestDefaultContentDirectoriesUseApprovedShape(t *testing.T) {
	cases := map[string]string{
		ContentChats:      "Chats",
		ContentCharacters: "Characters",
		ContentWorlds:     "Worlds",
		ContentScenarios:  "Scenarios",
		ContentLorebooks:  "Lorebooks",
		ContentImages:     "Howl",
		ContentVideos:     "Howl",
	}
	for category, leaf := range cases {
		dir, err := DefaultContentDir(category)
		if err != nil {
			t.Fatalf("%s: %v", category, err)
		}
		if !strings.HasSuffix(strings.ToLower(dir), strings.ToLower(filepath.Join(string(filepath.Separator), leaf))) &&
			!strings.HasSuffix(strings.ToLower(dir), strings.ToLower(leaf)) {
			t.Fatalf("%s default = %q, want leaf %q", category, dir, leaf)
		}
	}
	if _, err := DefaultContentDir("music"); err == nil {
		t.Fatal("music must remain unimplemented")
	}
}
