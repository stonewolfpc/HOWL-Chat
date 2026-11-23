//go:build windows

package paths

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// makeJunction creates a directory junction with mklink. Junctions do not require
// elevation (unlike symbolic links or hard links), so this exercises the real
// reparse-point case rather than simulating it.
func makeJunction(t *testing.T, link, target string) {
	t.Helper()
	if _, err := exec.LookPath("cmd"); err != nil {
		t.Skipf("cmd not available: %v", err)
	}
	out, err := exec.Command("cmd", "/c", "mklink", "/J", link, target).CombinedOutput()
	if err != nil {
		t.Skipf("junction creation unavailable (privilege or filesystem): %v: %s", err, out)
	}
	if _, statErr := os.Stat(link); statErr != nil {
		t.Skipf("junction not created: %v", statErr)
	}
}

// A path that is lexically inside the root but resolves outside it must be
// refused. Until this test existed the identity path was exercised but a junction
// escape was never demonstrated.
func TestCheckContainedRefusesJunctionEscape(t *testing.T) {
	base := t.TempDir()
	root := filepath.Join(base, "root")
	outside := filepath.Join(base, "outside")
	for _, dir := range []string{filepath.Join(root, "inside"), filepath.Join(outside, "secret")} {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatalf("mkdir %q: %v", dir, err)
		}
	}

	link := filepath.Join(root, "esc")
	makeJunction(t, link, outside)

	// Lexically under the root; physically under "outside".
	candidate := filepath.Join(link, "secret", "file.txt")
	if err := CheckContained(root, candidate); !errors.Is(err, ErrNotContained) {
		t.Fatalf("junction escape accepted: %v (candidate %q)", err, candidate)
	}

	// The junction directory itself is physically outside the root too.
	if err := CheckContained(root, filepath.Join(link, "secret")); !errors.Is(err, ErrNotContained) {
		t.Errorf("existing path behind a junction accepted: %v", err)
	}

	// A genuine child of the root is still accepted, so the refusal above is about
	// the reparse point and not about the test's shape.
	real := filepath.Join(root, "inside", "file.txt")
	if err := CheckContained(root, real); err != nil {
		t.Errorf("legitimate child refused: %v", err)
	}
}

// Containment must not rely on a shared volume. A candidate on another volume
// cannot be inside the root even when the textual path is well formed.
func TestCheckContainedRefusesOtherVolume(t *testing.T) {
	root := t.TempDir()
	other := otherVolume(t, root)
	candidate := filepath.Join(other, "howl-probe-does-not-exist")

	if err := CheckContained(root, candidate); !errors.Is(err, ErrNotContained) {
		t.Fatalf("cross-volume candidate accepted: %v (candidate %q)", err, candidate)
	}
}

// otherVolume returns an existing volume different from the one holding from.
func otherVolume(t *testing.T, from string) string {
	t.Helper()
	want, err := identityOf(from)
	if err != nil {
		t.Skipf("identity of %q: %v", from, err)
	}
	for letter := 'A'; letter <= 'Z'; letter++ {
		drive := string(letter) + ":" + string(os.PathSeparator)
		if _, statErr := os.Stat(drive); statErr != nil {
			continue
		}
		id, idErr := identityOf(drive)
		if idErr != nil {
			continue
		}
		if id.Volume != want.Volume {
			return drive
		}
	}
	t.Skip("only one volume is present; cross-volume containment cannot be exercised")
	return ""
}
