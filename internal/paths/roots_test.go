package paths

import (
	"errors"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

// useIsolatedRoot prepares a valid isolation root and clears the package cache.
// The caller must restore the environment; t.Setenv does that automatically.
func useIsolatedRoot(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, IsolationMarkerName), []byte(IsolationMarkerContent), 0o644); err != nil {
		t.Fatalf("write marker: %v", err)
	}
	t.Setenv(IsolationEnvVar, root)
	resetIsolationCacheForTest()
	t.Cleanup(resetIsolationCacheForTest)
	return root
}

func TestRootsInstalledMode(t *testing.T) {
	t.Setenv(IsolationEnvVar, "")
	resetIsolationCacheForTest()
	t.Cleanup(resetIsolationCacheForTest)

	rs, err := Roots()
	if err != nil {
		t.Fatalf("Roots: %v", err)
	}
	if rs.Mode != ModeInstalled {
		t.Fatalf("mode = %q, want %q", rs.Mode, ModeInstalled)
	}

	cfg, err := os.UserConfigDir()
	if err != nil {
		t.Skipf("no user config dir on this platform: %v", err)
	}
	wantApp := filepath.Join(cfg, "Howl")

	app, ok := rs.Get(RootApplicationData)
	if !ok {
		t.Fatal("application-data root missing")
	}
	if app.Path != wantApp {
		t.Errorf("application data = %q, want %q", app.Path, wantApp)
	}
	if app.Class != ClassAppData || !app.Writable {
		t.Errorf("application data class/writable = %q/%v", app.Class, app.Writable)
	}

	inst, ok := rs.Get(RootInstallation)
	if !ok {
		t.Fatal("installation root missing")
	}
	if inst.Writable {
		t.Error("installation root must be read-only")
	}
	if inst.Class != ClassInstallation {
		t.Errorf("installation class = %q", inst.Class)
	}

	tmp, ok := rs.Get(RootTemporary)
	if !ok || tmp.Class != ClassTemporary {
		t.Fatalf("temporary root missing or misclassified: %+v", tmp)
	}
}

func TestRootsCoversEveryContentCategory(t *testing.T) {
	t.Setenv(IsolationEnvVar, "")
	resetIsolationCacheForTest()
	t.Cleanup(resetIsolationCacheForTest)

	rs, err := Roots()
	if err != nil {
		t.Fatalf("Roots: %v", err)
	}
	for _, category := range contentCategories() {
		root, ok := rs.Category(category)
		if !ok {
			t.Errorf("category %q not resolved", category)
			continue
		}
		if root.Class != ClassUserContent {
			t.Errorf("category %q class = %q, want %q", category, root.Class, ClassUserContent)
		}
		if !root.Writable {
			t.Errorf("category %q must be writable by the application", category)
		}
		if !filepath.IsAbs(root.Path) {
			t.Errorf("category %q path %q is not absolute", category, root.Path)
		}
	}
	if _, ok := rs.Category("nonsense"); ok {
		t.Error("unknown category resolved unexpectedly")
	}
}

// The working directory is never a storage authority (ADR-001 A2).
func TestRootsAreWorkingDirectoryIndependent(t *testing.T) {
	t.Setenv(IsolationEnvVar, "")
	resetIsolationCacheForTest()
	t.Cleanup(resetIsolationCacheForTest)

	first, err := Roots()
	if err != nil {
		t.Fatalf("Roots: %v", err)
	}
	before := first.Describe()

	// Create the directories first so their cleanup runs after the chdir restore.
	dirs := []string{t.TempDir(), t.TempDir()}

	original, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	t.Cleanup(func() { _ = os.Chdir(original) })

	for _, dir := range dirs {
		if err := os.Chdir(dir); err != nil {
			t.Fatalf("chdir %q: %v", dir, err)
		}
		got, err := Roots()
		if err != nil {
			t.Fatalf("Roots after chdir: %v", err)
		}
		if got.Describe() != before {
			t.Fatalf("layout changed with working directory:\n--- before ---\n%s\n--- after (%s) ---\n%s", before, dir, got.Describe())
		}
	}

	if err := os.Chdir(original); err != nil {
		t.Fatalf("restore cwd: %v", err)
	}
}

func TestRootsIsolatedMode(t *testing.T) {
	root := useIsolatedRoot(t)

	rs, err := Roots()
	if err != nil {
		t.Fatalf("Roots: %v", err)
	}
	if rs.Mode != ModeIsolated {
		t.Fatalf("mode = %q, want %q", rs.Mode, ModeIsolated)
	}

	app, _ := rs.Get(RootApplicationData)
	if app.Path != filepath.Join(root, "data") {
		t.Errorf("isolated application data = %q, want data under %q", app.Path, root)
	}
	if err := CheckContained(root, app.Path); err != nil {
		t.Errorf("isolated application data escapes the isolation root: %v", err)
	}

	for _, category := range contentCategories() {
		cat, ok := rs.Category(category)
		if !ok {
			t.Fatalf("category %q missing under isolation", category)
		}
		if err := CheckContained(root, cat.Path); err != nil {
			t.Errorf("isolated category %q escapes the isolation root: %v", category, err)
		}
	}
}

// Resolution must not create or modify anything. Isolation setup is excluded:
// it creates its own directories before any call to Roots.
func TestRootsPerformsNoWrites(t *testing.T) {
	root := useIsolatedRoot(t)

	if _, err := Roots(); err != nil {
		t.Fatalf("Roots: %v", err)
	}
	before := listDir(t, root)

	for i := 0; i < 3; i++ {
		if _, err := Roots(); err != nil {
			t.Fatalf("Roots: %v", err)
		}
	}
	after := listDir(t, root)

	if strings.Join(before, "\n") != strings.Join(after, "\n") {
		t.Fatalf("Roots modified the filesystem:\nbefore: %v\nafter:  %v", before, after)
	}
}

// I-02, pinned rather than left implied: under isolation the existing resolver
// creates its layout on the first resolution. Repeated resolutions create nothing.
// This documents current behaviour so that changing it is a deliberate act.
func TestIsolationResolutionCreatesLayoutOnce(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, IsolationMarkerName), []byte(IsolationMarkerContent), 0o644); err != nil {
		t.Fatalf("write marker: %v", err)
	}
	t.Setenv(IsolationEnvVar, root)
	resetIsolationCacheForTest()
	t.Cleanup(resetIsolationCacheForTest)

	if _, err := os.Stat(filepath.Join(root, "data")); !os.IsNotExist(err) {
		t.Fatalf("precondition failed: isolated data directory already exists")
	}

	if _, err := Roots(); err != nil {
		t.Fatalf("Roots: %v", err)
	}
	for _, name := range []string{"data", "temp", "webview2"} {
		if _, err := os.Stat(filepath.Join(root, name)); err != nil {
			t.Errorf("isolation layout missing %q after first resolution: %v", name, err)
		}
	}

	before := listDir(t, root)
	if _, err := Roots(); err != nil {
		t.Fatalf("Roots again: %v", err)
	}
	if after := listDir(t, root); strings.Join(before, "\n") != strings.Join(after, "\n") {
		t.Errorf("repeated resolution changed the layout: %v -> %v", before, after)
	}
}

func listDir(t *testing.T, dir string) []string {
	t.Helper()
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("read %q: %v", dir, err)
	}
	out := make([]string, 0, len(entries))
	for _, e := range entries {
		out = append(out, e.Name())
	}
	sort.Strings(out)
	return out
}

func TestCheckContainedRejectsRelativeAndEmpty(t *testing.T) {
	root := t.TempDir()
	for _, candidate := range []string{"", "relative/path", "."} {
		if err := CheckContained(root, candidate); !errors.Is(err, ErrNotContained) {
			t.Errorf("CheckContained(%q) = %v, want ErrNotContained", candidate, err)
		}
	}
	if err := CheckContained("relative", filepath.Join(root, "x")); !errors.Is(err, ErrNotContained) {
		t.Errorf("relative root accepted: %v", err)
	}
}

func TestCheckContainedRejectsTraversalAndSiblings(t *testing.T) {
	base := t.TempDir()
	root := filepath.Join(base, "root")
	if err := os.MkdirAll(filepath.Join(root, "nested"), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	if err := CheckContained(root, filepath.Join(root, "nested")); err != nil {
		t.Fatalf("nested path refused: %v", err)
	}

	escape := filepath.Join(root, "..", "outside")
	if err := CheckContained(root, escape); !errors.Is(err, ErrNotContained) {
		t.Errorf("traversal accepted: %v", err)
	}

	sibling := filepath.Join(base, "root-evil")
	if err := os.MkdirAll(sibling, 0o755); err != nil {
		t.Fatalf("mkdir sibling: %v", err)
	}
	if err := CheckContained(root, sibling); !errors.Is(err, ErrNotContained) {
		t.Errorf("prefix-sharing sibling accepted: %v", err)
	}
}

func TestCheckContainedAcceptsNonExistentChild(t *testing.T) {
	root := t.TempDir()
	candidate := filepath.Join(root, "not", "created", "yet.txt")
	if err := CheckContained(root, candidate); err != nil {
		t.Fatalf("non-existent child inside root refused: %v", err)
	}
}

func TestFileIdentityIsStableAndDistinct(t *testing.T) {
	dir := t.TempDir()
	a := filepath.Join(dir, "a.txt")
	b := filepath.Join(dir, "b.txt")
	for _, p := range []string{a, b} {
		if err := os.WriteFile(p, []byte("x"), 0o644); err != nil {
			t.Fatalf("write %q: %v", p, err)
		}
	}
	idA, err := identityOf(a)
	if err != nil {
		t.Fatalf("identity a: %v", err)
	}
	idA2, err := identityOf(a)
	if err != nil {
		t.Fatalf("identity a again: %v", err)
	}
	idB, err := identityOf(b)
	if err != nil {
		t.Fatalf("identity b: %v", err)
	}
	if idA != idA2 {
		t.Errorf("identity is not stable: %+v vs %+v", idA, idA2)
	}
	if idA == idB {
		t.Errorf("distinct files share identity: %+v", idA)
	}
}

func TestRootSetCheckReportsUnknownCategory(t *testing.T) {
	t.Setenv(IsolationEnvVar, "")
	resetIsolationCacheForTest()
	t.Cleanup(resetIsolationCacheForTest)

	rs, err := Roots()
	if err != nil {
		t.Fatalf("Roots: %v", err)
	}
	if err := rs.Check("nonsense", filepath.Join(t.TempDir(), "x")); !errors.Is(err, ErrUnknownCategory) {
		t.Errorf("unknown category check = %v, want ErrUnknownCategory", err)
	}
	chatRoot, _ := rs.Category(ContentChats)
	if err := rs.Check(ContentChats, filepath.Join(chatRoot.Path, "conversation-1")); err != nil {
		t.Errorf("nested chat path refused: %v", err)
	}
}

func TestExternalRootIsReadOnlyAndNeverGrants(t *testing.T) {
	dir := t.TempDir()
	root, err := ExternalRoot(dir)
	if err != nil {
		t.Fatalf("ExternalRoot: %v", err)
	}
	if root.Class != ClassExternal {
		t.Errorf("class = %q, want %q", root.Class, ClassExternal)
	}
	if root.Writable {
		t.Error("an external root must never be reported writable")
	}

	if _, err := ExternalRoot("relative/models"); !errors.Is(err, ErrNotContained) {
		t.Errorf("relative external root accepted: %v", err)
	}
	missing := filepath.Join(t.TempDir(), "does-not-exist")
	if _, err := ExternalRoot(missing); !errors.Is(err, ErrUnavailable) {
		t.Errorf("missing external root accepted: %v", err)
	}
}

func TestCheckGrantRequiresStrictlyInside(t *testing.T) {
	dest := t.TempDir()
	operation := filepath.Join(dest, ".howl-import-staging", "op-1")
	if err := os.MkdirAll(operation, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	if err := CheckGrant(dest, filepath.Join(operation, "item", "file.bin")); err != nil {
		t.Errorf("item inside the operation subtree refused: %v", err)
	}
	if err := CheckGrant(dest, dest); !errors.Is(err, ErrNotContained) {
		t.Errorf("destination root itself was granted: %v", err)
	}
	if err := CheckGrant(dest, filepath.Join(dest, ".howl-import-staging")); err != nil {
		// The namespace directory itself is inside the root but is not the root;
		// only equality with the root is refused.
		t.Errorf("namespace directory refused: %v", err)
	}
	if err := CheckGrant(dest, filepath.Join(t.TempDir(), "elsewhere")); !errors.Is(err, ErrNotContained) {
		t.Errorf("path outside the destination accepted: %v", err)
	}
}

func TestDescribeNamesModeAndEveryRoot(t *testing.T) {
	t.Setenv(IsolationEnvVar, "")
	resetIsolationCacheForTest()
	t.Cleanup(resetIsolationCacheForTest)

	rs, err := Roots()
	if err != nil {
		t.Fatalf("Roots: %v", err)
	}
	desc := rs.Describe()
	if !strings.Contains(desc, "mode: installed") {
		t.Errorf("describe missing mode: %q", desc)
	}
	for _, root := range rs.All() {
		if !strings.Contains(desc, root.Path) {
			t.Errorf("describe missing %s (%s)", root.ID, root.Path)
		}
	}
}
