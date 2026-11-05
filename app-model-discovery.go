package main

// app-model-discovery.go
//
// Directory discovery has two jobs, and they are deliberately separate from
// model loading:
//
//  1. The folder locator. The native Windows folder dialog cannot be themed,
//     ignores the app's scaling rules, and hides what a folder actually
//     contains. BrowseModelDirectory backs an in-app locator that lists drives
//     and folders, shows how many model files each one holds, and says which
//     picker a scan would populate.
//
//  2. The "new model folder detected" prompt. A folder is only offered when it
//     sits beside a directory the user already trusts and actually contains
//     model files, so the prompt is a confirmation, not a file-system tour.
//
// Nothing here loads a model, starts inference or reserves VRAM. Scanning only
// reads metadata, classifies it and records the result (see the admission cache
// in app-model-registry.go).
//
// PREFERENCES
//
//   - enabled      (default true)  : offer newly found folders.
//   - never_prompt (default false) : set by "Never show again"; discovery keeps
//     working silently in the background so pickers stay accurate.
//
// Both live in the normal settings map, so they survive restarts and updates and
// are reversible from Settings. Turning prompts back on clears never_prompt.

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"howl-chat/internal/image/registry"
)

const (
	discoverySettingsKey   = "modelDiscovery"
	discoveryRecordFile    = "model-discovery.json"
	discoveryDeferWindow   = 24 * time.Hour
	discoveryCandidateMax  = 3
	discoveryInspectBudget = 800 // files inspected per folder, so a huge tree cannot stall the UI
	discoveryFoldersMax    = 60  // sibling folders examined per parent during detection
)

// DiscoveryCandidate is one folder offered to the user.
type DiscoveryCandidate struct {
	Path       string   `json:"path"`
	Name       string   `json:"name"`
	ModelFiles int      `json:"model_files"`
	Kind       string   `json:"kind"`
	Kinds      []string `json:"kinds"`
}

// ModelDiscoveryState is what the prompt and Settings toggle read.
type ModelDiscoveryState struct {
	Enabled    bool `json:"enabled"`
	Suppressed bool `json:"suppressed"`
	// Prompting is true only when prompts are enabled and not suppressed.
	Prompting bool                `json:"prompting"`
	Candidate *DiscoveryCandidate `json:"candidate,omitempty"`
	Deferred  int                 `json:"deferred"`
}

// DirectoryFolder is one row in the folder locator.
type DirectoryFolder struct {
	Name       string `json:"name"`
	Path       string `json:"path"`
	ModelFiles int    `json:"model_files"`
}

// DirectoryBrowse is one location in the folder locator.
type DirectoryBrowse struct {
	Path       string            `json:"path"`
	Parent     string            `json:"parent"`
	AtRoot     bool              `json:"at_root"`
	Home       string            `json:"home"`
	Drives     []string          `json:"drives"`
	Folders    []DirectoryFolder `json:"folders"`
	ModelFiles int               `json:"model_files"`
	Kind       string            `json:"kind"`
}

// discoveryRecords is per-directory state that is not part of the model
// registry: when a prompt was deferred and when a directory was last scanned.
type discoveryRecords struct {
	Deferred    map[string]time.Time `json:"deferred,omitempty"`
	LastScanned map[string]time.Time `json:"last_scanned,omitempty"`
}

var discoveryMu sync.Mutex

// ---- preferences ------------------------------------------------------------

func (a *App) discoveryPrefMap() map[string]interface{} {
	if a.appSettings == nil {
		return map[string]interface{}{}
	}
	if m, ok := a.appSettings[discoverySettingsKey].(map[string]interface{}); ok {
		return m
	}
	return map[string]interface{}{}
}

func (a *App) discoveryEnabled() bool {
	if v, ok := a.discoveryPrefMap()["enabled"].(bool); ok {
		return v
	}
	return true // beginners get the explanation by default
}

func (a *App) discoverySuppressed() bool {
	v, ok := a.discoveryPrefMap()["never_prompt"].(bool)
	return ok && v
}

func (a *App) writeDiscoveryPrefs(enabled, suppressed bool) error {
	next := make(map[string]interface{}, len(a.appSettings)+1)
	for k, v := range a.appSettings {
		next[k] = v
	}
	next[discoverySettingsKey] = map[string]interface{}{
		"enabled":      enabled,
		"never_prompt": suppressed,
	}
	if a.settingsStore == nil {
		a.appSettings = next
		return nil
	}
	return a.SetSettings(next)
}

// ---- records ----------------------------------------------------------------

// discoveryBaseDir is where per-directory records live. Tests redirect it so
// they never touch the user's real discovery state.
func (a *App) discoveryBaseDir() string {
	if a.discoveryDirOverride != "" {
		return a.discoveryDirOverride
	}
	return a.getSettingsDir()
}

func (a *App) discoveryRecordPath() string {
	return filepath.Join(a.discoveryBaseDir(), discoveryRecordFile)
}

func (a *App) loadDiscoveryRecords() discoveryRecords {
	recs := discoveryRecords{Deferred: map[string]time.Time{}, LastScanned: map[string]time.Time{}}
	data, err := os.ReadFile(a.discoveryRecordPath())
	if err != nil {
		return recs
	}
	if err := json.Unmarshal(data, &recs); err != nil {
		return discoveryRecords{Deferred: map[string]time.Time{}, LastScanned: map[string]time.Time{}}
	}
	if recs.Deferred == nil {
		recs.Deferred = map[string]time.Time{}
	}
	if recs.LastScanned == nil {
		recs.LastScanned = map[string]time.Time{}
	}
	return recs
}

func (a *App) saveDiscoveryRecords(recs discoveryRecords) error {
	if recs.Deferred == nil {
		recs.Deferred = map[string]time.Time{}
	}
	if recs.LastScanned == nil {
		recs.LastScanned = map[string]time.Time{}
	}
	dir := a.discoveryBaseDir()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(recs, "", "  ")
	if err != nil {
		return err
	}
	tmp := a.discoveryRecordPath() + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, a.discoveryRecordPath())
}

// ---- public bindings --------------------------------------------------------

// GetModelDiscoveryState reports the current preference and, when prompting is
// allowed, the best folder that has not been scanned yet.
func (a *App) GetModelDiscoveryState() (ModelDiscoveryState, error) {
	st := ModelDiscoveryState{
		Enabled:    a.discoveryEnabled(),
		Suppressed: a.discoverySuppressed(),
	}
	st.Prompting = st.Enabled && !st.Suppressed
	if !st.Enabled || a.modelRegistry == nil {
		return st, nil
	}
	recs := a.loadDiscoveryRecords()
	now := time.Now().UTC()
	for _, at := range recs.Deferred {
		if now.Sub(at) < discoveryDeferWindow {
			st.Deferred++
		}
	}
	if st.Suppressed {
		return st, nil
	}
	candidates := a.detectModelDirectories()
	if len(candidates) > 0 {
		st.Candidate = &candidates[0]
	}
	return st, nil
}

// SetDirectoryDiscovery turns the prompt on or off. Turning it off keeps
// discovery running silently, so pickers stay accurate without interruptions;
// turning it on clears a previous "Never show again".
func (a *App) SetDirectoryDiscovery(prompt bool) error {
	return a.writeDiscoveryPrefs(true, !prompt)
}

// DeferDiscoveredDirectory postpones the prompt for this folder. It is
// remembered on disk and re-offered after a day rather than on every launch.
func (a *App) DeferDiscoveredDirectory(path string) error {
	abs, err := normalizeDirectory(path)
	if err != nil {
		return err
	}
	discoveryMu.Lock()
	defer discoveryMu.Unlock()
	recs := a.loadDiscoveryRecords()
	recs.Deferred[abs] = time.Now().UTC()
	return a.saveDiscoveryRecords(recs)
}

// ScanDiscoveredDirectory classifies a folder, registers it under the kind its
// contents imply, and populates the matching pickers. It never loads a model.
func (a *App) ScanDiscoveredDirectory(path string) (ModelDirectory, error) {
	if a.modelRegistry == nil {
		return ModelDirectory{}, fmt.Errorf("model registry not initialised")
	}
	abs, err := normalizeDirectory(path)
	if err != nil {
		return ModelDirectory{}, err
	}
	kind, files := inspectModelDirectory(abs)
	if files == 0 || kind == "" {
		return ModelDirectory{}, fmt.Errorf("no model files found in %q", abs)
	}
	// A directory keeps the kind it was first added under; contents decide only
	// for a directory that has never been registered.
	if existing := a.modelRegistry.KindForPath(abs); existing != "" {
		kind = existing
	}
	count, scanErr := a.scanDirectoryIntoRegistry(abs, kind)
	if err := a.modelRegistry.RegisterDir(abs, kind); err != nil {
		return ModelDirectory{}, err
	}

	discoveryMu.Lock()
	recs := a.loadDiscoveryRecords()
	recs.LastScanned[abs] = time.Now().UTC()
	delete(recs.Deferred, abs)
	_ = a.saveDiscoveryRecords(recs)
	discoveryMu.Unlock()

	if scanErr != nil {
		return ModelDirectory{}, scanErr
	}
	return ModelDirectory{Path: abs, Kind: kind, Exists: true, Models: count}, nil
}

// BrowseModelDirectory backs the in-app folder locator. An empty path starts at
// the user's home directory.
func (a *App) BrowseModelDirectory(path string) (DirectoryBrowse, error) {
	home, _ := os.UserHomeDir()
	path = strings.TrimSpace(path)
	if path == "" {
		path = home
	}
	abs, err := normalizeDirectory(path)
	if err != nil {
		return DirectoryBrowse{}, err
	}
	res := DirectoryBrowse{
		Path:   abs,
		Home:   home,
		Drives: driveRoots(),
	}
	parent := filepath.Dir(abs)
	res.Parent = parent
	res.AtRoot = parent == abs || parent == ""

	entries, err := os.ReadDir(abs)
	if err != nil {
		return res, fmt.Errorf("cannot read %q: %w", abs, err)
	}
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		name := e.Name()
		if strings.HasPrefix(name, ".") || strings.EqualFold(name, "$RECYCLE.BIN") ||
			strings.EqualFold(name, "System Volume Information") {
			continue
		}
		folder := filepath.Join(abs, name)
		res.Folders = append(res.Folders, DirectoryFolder{
			Name:       name,
			Path:       folder,
			ModelFiles: countModelFilesShallow(folder),
		})
	}
	sort.Slice(res.Folders, func(i, j int) bool {
		return strings.ToLower(res.Folders[i].Name) < strings.ToLower(res.Folders[j].Name)
	})
	res.Kind, res.ModelFiles = inspectModelDirectory(abs)
	return res, nil
}

// RunSilentModelDiscovery keeps pickers accurate after "Never show again"
// without interrupting the user. It runs in the background at startup.
func (a *App) RunSilentModelDiscovery() {
	if !a.discoveryEnabled() || !a.discoverySuppressed() {
		return
	}
	if a.modelRegistry == nil {
		return
	}
	for _, candidate := range a.detectModelDirectories() {
		if _, err := a.ScanDiscoveredDirectory(candidate.Path); err != nil {
			fmt.Printf("[discovery] silent scan skipped %s: %v\n", candidate.Path, err)
			continue
		}
		startupMark("discovery: silently scanned " + candidate.Path)
	}
}

// ---- detection --------------------------------------------------------------

// detectModelDirectories looks beside directories the user already trusts. It
// never walks a drive: it reads one level of each known parent and inspects the
// folders it finds, so the cost stays flat regardless of library size.
func (a *App) detectModelDirectories() []DiscoveryCandidate {
	recs := a.loadDiscoveryRecords()
	roots := map[string]bool{}
	for _, kind := range []string{registry.KindLLM, registry.KindTTS, registry.KindImage} {
		for _, dir := range a.modelRegistry.Directories(kind) {
			parent := filepath.Dir(dir)
			if parent != "" && parent != dir {
				roots[parent] = true
			}
		}
	}

	now := time.Now().UTC()
	seen := map[string]bool{}
	var out []DiscoveryCandidate
	for root := range roots {
		entries, err := os.ReadDir(root)
		if err != nil {
			continue
		}
		examined := 0
		for _, e := range entries {
			if examined >= discoveryFoldersMax {
				break
			}
			if !e.IsDir() || strings.HasPrefix(e.Name(), ".") {
				continue
			}
			examined++
			path := filepath.Join(root, e.Name())
			key := strings.ToLower(filepath.Clean(path))
			if seen[key] {
				continue
			}
			if a.modelRegistry.KindForPath(path) != "" {
				continue // already scanned
			}
			if at, ok := recs.Deferred[path]; ok && now.Sub(at) < discoveryDeferWindow {
				continue
			}
			kind, files := inspectModelDirectory(path)
			if files == 0 || kind == "" {
				continue
			}
			seen[key] = true
			out = append(out, DiscoveryCandidate{
				Path:       path,
				Name:       e.Name(),
				ModelFiles: files,
				Kind:       kind,
				Kinds:      modelKindsIn(path),
			})
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ModelFiles > out[j].ModelFiles })
	if len(out) > discoveryCandidateMax {
		out = out[:discoveryCandidateMax]
	}
	return out
}

func normalizeDirectory(path string) (string, error) {
	path = strings.TrimSpace(path)
	if path == "" {
		return "", fmt.Errorf("directory path is empty")
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		abs = filepath.Clean(path)
	}
	info, err := os.Stat(abs)
	if err != nil {
		return "", fmt.Errorf("cannot access %q: %w", abs, err)
	}
	if !info.IsDir() {
		return "", fmt.Errorf("%q is not a directory", abs)
	}
	return abs, nil
}

// inspectModelDirectory counts model files up to two levels deep and picks the
// kind with the most evidence. Depth two catches the common "<family>/<quant>"
// layout without turning detection into a disk crawl.
func inspectModelDirectory(dir string) (string, int) {
	counts := map[string]int{}
	inspected := 0
	_ = filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			rel, relErr := filepath.Rel(dir, path)
			if relErr == nil && rel != "." && strings.Count(rel, string(filepath.Separator)) >= 2 {
				return filepath.SkipDir
			}
			if strings.HasPrefix(d.Name(), ".") && path != dir {
				return filepath.SkipDir
			}
			return nil
		}
		if inspected >= discoveryInspectBudget {
			return filepath.SkipAll
		}
		inspected++
		switch strings.ToLower(filepath.Ext(d.Name())) {
		case ".gguf":
			counts[registry.KindLLM]++
		case ".onnx":
			counts[registry.KindTTS]++
		case ".safetensors", ".ckpt", ".sft", ".pt", ".pth":
			counts[registry.KindImage]++
		}
		return nil
	})

	best := ""
	bestCount := 0
	// Ties resolve image first, then chat, then voice: an image library that
	// also ships a text encoder must not be filed as a chat directory.
	for _, kind := range []string{registry.KindImage, registry.KindLLM, registry.KindTTS} {
		if counts[kind] > bestCount {
			best, bestCount = kind, counts[kind]
		}
	}
	return best, bestCount
}

// modelKindsIn names every kind the folder can feed, for the prompt copy.
func modelKindsIn(dir string) []string {
	counts := map[string]int{}
	_ = filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			rel, relErr := filepath.Rel(dir, path)
			if relErr == nil && rel != "." && strings.Count(rel, string(filepath.Separator)) >= 2 {
				return filepath.SkipDir
			}
			return nil
		}
		switch strings.ToLower(filepath.Ext(d.Name())) {
		case ".gguf":
			counts[registry.KindLLM]++
		case ".onnx":
			counts[registry.KindTTS]++
		case ".safetensors", ".ckpt", ".sft", ".pt", ".pth":
			counts[registry.KindImage]++
		}
		return nil
	})
	labels := map[string]string{
		registry.KindLLM:   "chat",
		registry.KindTTS:   "voice",
		registry.KindImage: "image generation",
	}
	var out []string
	for _, kind := range []string{registry.KindLLM, registry.KindTTS, registry.KindImage} {
		if counts[kind] > 0 {
			out = append(out, labels[kind])
		}
	}
	return out
}

// countModelFilesShallow counts model files directly inside a folder, which is
// what the locator shows per row. Cheap on purpose.
func countModelFilesShallow(dir string) int {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return 0
	}
	count := 0
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		switch strings.ToLower(filepath.Ext(e.Name())) {
		case ".gguf", ".onnx", ".safetensors", ".ckpt", ".sft", ".pt", ".pth":
			count++
		}
	}
	return count
}

// driveRoots lists mounted drive letters. Letters that are not mounted are
// skipped, so this is safe on machines with a single drive.
func driveRoots() []string {
	var roots []string
	for c := 'A'; c <= 'Z'; c++ {
		root := string(c) + ":\\"
		if info, err := os.Stat(root); err == nil && info.IsDir() {
			roots = append(roots, root)
		}
	}
	return roots
}
