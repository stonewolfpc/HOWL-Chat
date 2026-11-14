package paths

import (
	"fmt"
	"path/filepath"
)

// Visible creative content belongs in user-owned, discoverable folders. These
// categories are deliberately separate from application state and downloaded
// models so a settings reset or uninstall cannot mistake user work for cache.
const (
	ContentChats      = "chats"
	ContentCharacters = "characters"
	ContentWorlds     = "worlds"
	ContentScenarios  = "scenarios"
	ContentLorebooks  = "lorebooks"
	ContentImages     = "images"
	ContentVideos     = "videos"
)

// DefaultContentDir returns the preferred visible location for one category.
// Windows known folders are resolved through the shell, not by assuming that
// Documents is always under the literal user home (OneDrive and redirected
// profiles are common).
func DefaultContentDir(category string) (string, error) {
	var folder, known string

	switch category {
	case ContentChats:
		folder, known = "Chats", "documents"
	case ContentCharacters:
		folder, known = "Characters", "documents"
	case ContentWorlds:
		folder, known = "Worlds", "documents"
	case ContentScenarios:
		folder, known = "Scenarios", "documents"
	case ContentLorebooks:
		folder, known = "Lorebooks", "documents"
	case ContentImages:
		known = "pictures"
	case ContentVideos:
		known = "videos"
	default:
		return "", fmt.Errorf("unknown content category %q", category)
	}

	// An isolated native test run must never point a content destination at the
	// user's real Documents/Pictures/Videos folders.
	if IsolationActive() {
		parts := []string{"content"}
		if folder != "" {
			parts = append(parts, folder)
		}
		return StatePath(parts...), nil
	}

	base, err := knownFolder(known)
	if err != nil {
		return "", err
	}
	if folder == "" {
		return filepath.Join(base, "Howl"), nil
	}
	return filepath.Join(base, "Howl", folder), nil
}

// InternalConfigDir is the existing application-owned settings and metadata
// location. F2 intentionally does not move it: changing that is a migration
// decision, and existing settings must remain reachable. New state files use
// this directory rather than the installation directory.
//
// It follows StateRoot, so an isolated native test run keeps its settings,
// registries and journals inside the isolation root.
func InternalConfigDir() string {
	return filepath.Join(StateRoot(), "config")
}
