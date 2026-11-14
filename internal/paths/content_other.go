//go:build !windows

package paths

import (
	"fmt"
	"os"
	"path/filepath"
)

// Non-Windows builds use the conventional home-relative folders. Howl ships on
// Windows; this keeps tests and development builds portable without pretending
// that the Windows known-folder API exists here.
func knownFolder(name string) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolve user home: %w", err)
	}
	switch name {
	case "documents":
		return filepath.Join(home, "Documents"), nil
	case "pictures":
		return filepath.Join(home, "Pictures"), nil
	case "videos":
		return filepath.Join(home, "Videos"), nil
	default:
		return "", fmt.Errorf("unknown known folder %q", name)
	}
}
