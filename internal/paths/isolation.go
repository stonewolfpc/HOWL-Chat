package paths

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

// IsolationEnvVar is the single opt-in switch for an isolated native test run.
//
// When it is set to an absolute directory, every application-owned persistent
// location (WebView2 profile, settings, logs, chats, worlds, characters,
// scenarios, images, caches, scratch space) is derived from that directory
// instead. Unset means ordinary behaviour, unchanged.
const IsolationEnvVar = "HOWL_WEBVIEW_DATA_DIR"

// IsolationMarkerName must exist inside the isolation root. Requiring it makes
// it impossible to isolate into a folder that already holds real data: a
// production directory never contains this file, so a mistyped path fails
// instead of quietly adopting a user's real content.
const IsolationMarkerName = ".howl-isolated-root"

// IsolationMarkerContent is the exact marker content accepted.
const IsolationMarkerContent = "howl-isolated-test"

// Isolation is the resolved isolated layout. Every path is a child of Root.
type Isolation struct {
	// Root is the operator-supplied directory that holds the whole test run.
	Root string
	// WebViewDataDir is the WebView2 user-data folder (options.App.Windows.
	// WebviewUserDataPath). It is separate from application state so browser
	// cache never mixes with product data.
	WebViewDataDir string
	// StateRoot replaces the installation root for everything the app writes.
	StateRoot string
	// TempDir replaces os.TempDir for scratch audio/image files.
	TempDir string
}

var (
	isolationOnce   sync.Once
	isolationConfig *Isolation
	isolationErr    error
)

// ResolveIsolation returns the isolation layout when IsolationEnvVar is set.
// A nil layout with a nil error means isolation is off.
func ResolveIsolation() (*Isolation, error) {
	isolationOnce.Do(func() {
		isolationConfig, isolationErr = resolveIsolation()
	})
	return isolationConfig, isolationErr
}

// resetIsolationCacheForTest lets the package's own tests re-resolve after
// changing the environment. The process resolves once in production.
func resetIsolationCacheForTest() {
	isolationOnce = sync.Once{}
	isolationConfig = nil
	isolationErr = nil
}

func resolveIsolation() (*Isolation, error) {
	raw := strings.TrimSpace(os.Getenv(IsolationEnvVar))
	if raw == "" {
		return nil, nil
	}
	return newIsolation(raw, resolve())
}

// newIsolation validates an operator-supplied root. realProject is the
// installation root, which must never contain the isolation root or be
// contained by it; it is a parameter so the rules can be tested without
// depending on where this checkout lives.
func newIsolation(raw, realProject string) (*Isolation, error) {
	if !filepath.IsAbs(raw) {
		return nil, fmt.Errorf("%s must be an absolute path, got %q", IsolationEnvVar, raw)
	}
	root := filepath.Clean(raw)

	if vol := filepath.VolumeName(root); vol != "" {
		rest := strings.TrimPrefix(root, vol)
		if rest == "" || rest == string(filepath.Separator) {
			return nil, fmt.Errorf("%s must not be a drive root (%s)", IsolationEnvVar, root)
		}
	}

	info, err := os.Stat(root)
	if err != nil {
		return nil, fmt.Errorf("%s points at %q, which is not usable: %v; create the directory and add a %s marker containing %q",
			IsolationEnvVar, root, err, IsolationMarkerName, IsolationMarkerContent)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("%s points at %q, which is not a directory", IsolationEnvVar, root)
	}

	if realProject != "" && (containsPath(realProject, root) || containsPath(root, realProject)) {
		return nil, fmt.Errorf("%s must be outside the project tree: %q overlaps %q", IsolationEnvVar, root, realProject)
	}

	markerPath := filepath.Join(root, IsolationMarkerName)
	marker, err := os.ReadFile(markerPath)
	if err != nil {
		return nil, fmt.Errorf("%q has no %s marker (%v); refusing to isolate into a directory that may hold real data",
			root, IsolationMarkerName, err)
	}
	if got := strings.TrimSpace(string(marker)); got != IsolationMarkerContent {
		return nil, fmt.Errorf("marker %q must contain %q, found %q", markerPath, IsolationMarkerContent, got)
	}

	probe, err := os.CreateTemp(root, ".howl-isolation-probe-*")
	if err != nil {
		return nil, fmt.Errorf("%q is not writable: %v", root, err)
	}
	probeName := probe.Name()
	_ = probe.Close()
	_ = os.Remove(probeName)

	cfg := &Isolation{
		Root:           root,
		WebViewDataDir: filepath.Join(root, "webview2"),
		StateRoot:      filepath.Join(root, "data"),
		TempDir:        filepath.Join(root, "temp"),
	}

	derived := map[string]string{
		"webview data": cfg.WebViewDataDir,
		"state":        cfg.StateRoot,
		"temp":         cfg.TempDir,
	}
	for name, p := range derived {
		if !containsPath(root, p) {
			return nil, fmt.Errorf("isolated %s path %q escapes %q", name, p, root)
		}
	}

	for _, dir := range []string{cfg.WebViewDataDir, cfg.StateRoot, cfg.TempDir, filepath.Join(cfg.StateRoot, "config")} {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return nil, fmt.Errorf("create isolated directory %q: %w", dir, err)
		}
	}

	return cfg, nil
}

// containsPath reports whether child is root itself or nested inside it. It
// compares resolved relative components, so a sibling whose name merely shares
// a textual prefix ("C:\root-evil" against "C:\root") is never treated as
// contained.
func containsPath(root, child string) bool {
	if root == "" || child == "" {
		return false
	}
	root = filepath.Clean(root)
	child = filepath.Clean(child)
	if runtime.GOOS == "windows" {
		root = strings.ToLower(root)
		child = strings.ToLower(child)
	}
	rel, err := filepath.Rel(root, child)
	if err != nil {
		return false
	}
	if rel == "." {
		return true
	}
	return rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator))
}

// IsolationActive reports whether an isolated layout is configured. It panics
// when the variable is set but invalid: silently continuing against production
// directories is the exact failure this switch exists to prevent. main.go
// validates first and exits with a clear message, so the panic is a last line
// of defence for other entry points.
func IsolationActive() bool {
	return mustIsolation() != nil
}

// StateRoot is the single root for everything the app persists: settings,
// logs, chats, worlds, characters, scenarios, images, caches and the legacy
// chat store. It is the installation root unless isolation is active.
func StateRoot() string {
	if cfg := mustIsolation(); cfg != nil {
		return cfg.StateRoot
	}
	return ProjectRoot()
}

// StatePath joins path elements onto StateRoot.
func StatePath(elem ...string) string {
	return filepath.Join(append([]string{StateRoot()}, elem...)...)
}

// TempDir is the scratch root: the isolated temp directory when isolation is
// active, os.TempDir otherwise.
func TempDir() string {
	if cfg := mustIsolation(); cfg != nil {
		return cfg.TempDir
	}
	return os.TempDir()
}

// WebViewDataDir is the WebView2 user-data folder to hand to Wails. Empty means
// "use the platform default", which is what a normal run needs.
func WebViewDataDir() string {
	if cfg := mustIsolation(); cfg != nil {
		return cfg.WebViewDataDir
	}
	return ""
}

func mustIsolation() *Isolation {
	cfg, err := ResolveIsolation()
	if err != nil {
		panic("isolation: " + err.Error())
	}
	return cfg
}
