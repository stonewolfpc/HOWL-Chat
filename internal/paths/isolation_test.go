package paths

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeIsolationMarker(t *testing.T, root string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(root, IsolationMarkerName), []byte(IsolationMarkerContent), 0o644); err != nil {
		t.Fatalf("write marker: %v", err)
	}
}

func TestContainsPathUsesComponentsNotPrefixes(t *testing.T) {
	base := t.TempDir()
	root := filepath.Join(base, "root")
	sibling := filepath.Join(base, "root-evil")

	if !containsPath(root, root) {
		t.Fatalf("root must contain itself")
	}
	if !containsPath(root, filepath.Join(root, "data", "config")) {
		t.Fatalf("nested path must be contained")
	}
	if containsPath(root, sibling) {
		t.Fatalf("sibling sharing a name prefix must not be contained: %q vs %q", root, sibling)
	}
	if containsPath(root, base) {
		t.Fatalf("parent must not be contained")
	}
}

func TestNewIsolationDerivesEverythingInsideRoot(t *testing.T) {
	root := t.TempDir()
	writeIsolationMarker(t, root)
	project := t.TempDir()

	cfg, err := newIsolation(root, project)
	if err != nil {
		t.Fatalf("newIsolation: %v", err)
	}
	if cfg.Root != root {
		t.Fatalf("root = %q, want %q", cfg.Root, root)
	}

	derived := map[string]string{
		"webview":          cfg.WebViewDataDir,
		"state":            cfg.StateRoot,
		"temp":             cfg.TempDir,
		"config":           filepath.Join(cfg.StateRoot, "config"),
		"logs":             filepath.Join(cfg.StateRoot, "logs"),
		"chats":            filepath.Join(cfg.StateRoot, "Chats"),
		"worlds":           filepath.Join(cfg.StateRoot, "Worlds"),
		"characters":       filepath.Join(cfg.StateRoot, "Characters"),
		"scenarios":        filepath.Join(cfg.StateRoot, "Scenarios"),
		"gallery":          filepath.Join(cfg.StateRoot, "images", "gallery"),
		"workflows":        filepath.Join(cfg.StateRoot, "workflows"),
		"persisted chats":  filepath.Join(cfg.StateRoot, "persisted-chats"),
		"cache":            filepath.Join(cfg.StateRoot, "cache"),
		"temp audio":       filepath.Join(cfg.TempDir, "howl-asr"),
		"webview cache":    filepath.Join(cfg.WebViewDataDir, "EBWebView"),
		"state root alone": cfg.StateRoot,
	}
	for name, p := range derived {
		if !containsPath(root, p) {
			t.Errorf("%s path %q is not inside the isolation root %q", name, p, root)
		}
		if containsPath(project, p) {
			t.Errorf("%s path %q resolves inside the real project %q", name, p, project)
		}
	}

	for _, dir := range []string{cfg.WebViewDataDir, cfg.StateRoot, cfg.TempDir, filepath.Join(cfg.StateRoot, "config")} {
		if info, err := os.Stat(dir); err != nil || !info.IsDir() {
			t.Errorf("isolated directory %q was not created: %v", dir, err)
		}
	}
}

func TestNewIsolationRejectsUnsafeRoots(t *testing.T) {
	project := t.TempDir()
	markerRoot := t.TempDir()
	writeIsolationMarker(t, markerRoot)

	nestedInProject := filepath.Join(project, "nested")
	if err := os.MkdirAll(nestedInProject, 0o755); err != nil {
		t.Fatalf("create nested project dir: %v", err)
	}

	fileAsRoot := filepath.Join(t.TempDir(), "not-a-dir")
	if err := os.WriteFile(fileAsRoot, []byte("x"), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	badMarkerRoot := t.TempDir()
	if err := os.WriteFile(filepath.Join(badMarkerRoot, IsolationMarkerName), []byte("production"), 0o644); err != nil {
		t.Fatalf("write marker: %v", err)
	}

	cases := []struct {
		name    string
		raw     string
		project string
		want    string
	}{
		{"relative path", filepath.Join("relative", "root"), project, "absolute path"},
		{"missing directory", filepath.Join(t.TempDir(), "absent"), project, "not usable"},
		{"root is a file", fileAsRoot, project, "not a directory"},
		{"missing marker", t.TempDir(), project, "no .howl-isolated-root marker"},
		{"wrong marker content", badMarkerRoot, project, "must contain"},
		{"inside the project", nestedInProject, project, "outside the project tree"},
		{"contains the project", filepath.Dir(project), project, "outside the project tree"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cfg, err := newIsolation(tc.raw, tc.project)
			if err == nil {
				t.Fatalf("expected failure, got %+v", cfg)
			}
			if !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("error %q does not mention %q", err.Error(), tc.want)
			}
		})
	}

	if _, err := newIsolation(markerRoot, project); err != nil {
		t.Fatalf("valid root rejected: %v", err)
	}
}

func TestIsolationDisabledKeepsProductionResolution(t *testing.T) {
	t.Setenv(IsolationEnvVar, "")
	resetIsolationCacheForTest()
	t.Cleanup(resetIsolationCacheForTest)

	cfg, err := ResolveIsolation()
	if err != nil {
		t.Fatalf("unexpected error with isolation unset: %v", err)
	}
	if cfg != nil {
		t.Fatalf("isolation must be off when the variable is empty, got %+v", cfg)
	}
	if got, want := StateRoot(), ProjectRoot(); got != want {
		t.Fatalf("StateRoot = %q, want the installation root %q", got, want)
	}
	if got, want := InternalConfigDir(), filepath.Join(ProjectRoot(), "config"); got != want {
		t.Fatalf("InternalConfigDir = %q, want %q", got, want)
	}
	if got, want := TempDir(), os.TempDir(); got != want {
		t.Fatalf("TempDir = %q, want %q", got, want)
	}
	if got := WebViewDataDir(); got != "" {
		t.Fatalf("WebViewDataDir = %q, want empty so Wails uses its default", got)
	}
	if IsolationActive() {
		t.Fatalf("IsolationActive must be false by default")
	}
}

func TestIsolationEnabledThroughEnvironment(t *testing.T) {
	root := t.TempDir()
	writeIsolationMarker(t, root)
	t.Setenv(IsolationEnvVar, root)
	resetIsolationCacheForTest()
	t.Cleanup(resetIsolationCacheForTest)

	if !IsolationActive() {
		t.Fatalf("IsolationActive must be true when %s is set", IsolationEnvVar)
	}
	stateRoot := StateRoot()
	if !containsPath(root, stateRoot) {
		t.Fatalf("StateRoot %q is outside the isolation root %q", stateRoot, root)
	}
	if containsPath(root, ProjectRoot()) {
		t.Fatalf("StateRoot must not contain the real project root; isolation root %q vs project %q", root, ProjectRoot())
	}
	if got := WebViewDataDir(); !containsPath(root, got) {
		t.Fatalf("WebViewDataDir %q is outside the isolation root %q", got, root)
	}
	if got := TempDir(); !containsPath(root, got) {
		t.Fatalf("TempDir %q is outside the isolation root %q", got, root)
	}
	if got := InternalConfigDir(); !containsPath(root, got) {
		t.Fatalf("InternalConfigDir %q is outside the isolation root %q", got, root)
	}
	if StatePath("Chats", "example") == filepath.Join(ProjectRoot(), "Chats", "example") {
		t.Fatalf("StatePath must not resolve into the real project while isolated")
	}
}

func TestIsolationInvalidConfigurationFailsLoudly(t *testing.T) {
	t.Setenv(IsolationEnvVar, filepath.Join("not", "absolute"))
	resetIsolationCacheForTest()
	t.Cleanup(resetIsolationCacheForTest)

	if _, err := ResolveIsolation(); err == nil {
		t.Fatalf("expected an error for a relative isolation root")
	}

	defer func() {
		if recover() == nil {
			t.Fatalf("IsolationActive must panic rather than fall back to production")
		}
	}()
	_ = IsolationActive()
}
