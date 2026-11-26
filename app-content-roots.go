package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	apperrors "howl-chat/internal/errors"
	"howl-chat/internal/library"
	"howl-chat/internal/paths"
)

// Content-library readiness and provisioning.
//
// One authority owns the only filesystem provisioning route for a category
// root. Reads never create directories: GetWorlds/GetScenarios report storage
// state instead of mutating it, so a blocked write cannot masquerade as an
// empty library. Provisioning happens once at startup, before the background
// scans, and on an explicit user action — never from a read.
//
// Readiness is reported per library. Reading and modifying are separate
// capabilities: a library that can be read but not written stays visible and is
// reported as read-only, never unavailable. The category rollup exists only so a
// panel can summarise; it keeps the per-library entries authoritative and counts
// the libraries that are not ready, so one healthy library cannot hide a failing
// sibling.

const (
	ContentReady          = "ready"
	ContentReadOnly       = "read-only"
	ContentNotReadable    = "not-readable"
	ContentNotInitialized = "not-initialized"
	ContentUnavailable    = "unavailable"
)

// ContentLibraryState is the authoritative readiness of one content library.
type ContentLibraryState struct {
	Category  string `json:"category"`
	LibraryID string `json:"libraryId"`
	Path      string `json:"path"`
	Readable  bool   `json:"readable"`
	Writable  bool   `json:"writable"`
	SaveRoot  bool   `json:"saveRoot"`
	Status    string `json:"status"`
	Reason    string `json:"reason"`
	Operation string `json:"operation"`
	Retryable bool   `json:"retryable"`
}

// ContentCategoryState is a presentation rollup of the per-library states.
// FailingLibraries is not decoration: it is how a category with three libraries
// reports that one of them is broken while the other two work.
type ContentCategoryState struct {
	Category string `json:"category"`
	Status   string `json:"status"`
	Readable bool   `json:"readable"`
	Writable bool   `json:"writable"`
	// SaveWritable reports whether the category's default save root can be
	// written. Creating goes there, so create is gated on this rather than on
	// Writable: a read-only default root must not block editing items that live
	// in a healthy sibling library, and that sibling must not make create look
	// possible when it would be refused.
	SaveWritable     bool   `json:"saveWritable"`
	Libraries        int    `json:"libraries"`
	FailingLibraries int    `json:"failingLibraries"`
	Reason           string `json:"reason"`
}

type ContentReadiness struct {
	Libraries  []ContentLibraryState  `json:"libraries"`
	Categories []ContentCategoryState `json:"categories"`
}

// ProvisionReport proves what a provisioning attempt actually did, step by
// step. Created directories are removed only when this call created them and
// they are empty; a pre-existing directory is never removed.
type ProvisionReport struct {
	Category        string
	Path            string
	RootCreated     bool
	CategoryCreated bool
	RootRemoved     bool
	CategoryRemoved bool
	FileCreated     bool
	BytesWritten    int
	ReadBack        bool
	Cleanup         bool
	Readable        bool
	Writable        bool
	Step            string
	Operation       string
	Status          string
	Reason          string
	Retryable       bool
	Err             error
}

// probeFile is the small surface the write probe needs, so a test can fail any
// step without ACLs.
type probeFile interface {
	Write([]byte) (int, error)
	Sync() error
	Close() error
}

// contentFSOps is the injectable filesystem seam for the probe.
type contentFSOps struct {
	MkdirAll func(path string, perm os.FileMode) error
	Stat     func(path string) (fs.FileInfo, error)
	ReadDir  func(path string) ([]fs.DirEntry, error)
	OpenFile func(path string) (probeFile, error)
	ReadFile func(path string) ([]byte, error)
	Remove   func(path string) error
}

func defaultContentFSOps() contentFSOps {
	return contentFSOps{
		MkdirAll: os.MkdirAll,
		Stat:     os.Stat,
		ReadDir:  os.ReadDir,
		OpenFile: func(path string) (probeFile, error) {
			return os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_RDWR, 0o600)
		},
		ReadFile: os.ReadFile,
		Remove:   os.Remove,
	}
}

// contentPathBlocked reports a regular file where a directory is required.
type contentPathBlocked struct{ Path string }

func (e contentPathBlocked) Error() string {
	return fmt.Sprintf("%s is a file, not a directory", e.Path)
}

var provisionedContentCategories = []library.Category{library.Worlds, library.Scenarios}

// probeContentDir performs, and reports, the full writability check: root
// creation, category creation, probe-file creation, write, readback, removal and
// cleanup. A successful MkdirAll alone is not treated as proof of writability.
// keepCreated leaves directories this call created in place, which is what
// provisioning needs. A pure probe (no argument) removes them again.
func probeContentDir(ops contentFSOps, category, dir string, keepCreated ...bool) ProvisionReport {
	keep := len(keepCreated) > 0 && keepCreated[0]
	report := ProvisionReport{Category: category, Path: dir}

	root := filepath.Dir(dir)
	rootInfo, rootErr := ops.Stat(root)
	switch {
	case rootErr == nil && !rootInfo.IsDir():
		return failProvision(ops, report, "inspect library root", contentPathBlocked{Path: root}, false)
	case rootErr != nil && !errors.Is(rootErr, fs.ErrNotExist):
		return failProvision(ops, report, "inspect library root", rootErr, false)
	case errors.Is(rootErr, fs.ErrNotExist):
		if err := ops.MkdirAll(root, 0o755); err != nil {
			return failProvision(ops, report, "create library root", err, false)
		}
		report.RootCreated = true
	}

	dirInfo, dirErr := ops.Stat(dir)
	switch {
	case dirErr == nil && !dirInfo.IsDir():
		return failProvision(ops, report, "inspect category directory", contentPathBlocked{Path: dir}, false)
	case dirErr != nil && !errors.Is(dirErr, fs.ErrNotExist):
		return failProvision(ops, report, "inspect category directory", dirErr, false)
	case errors.Is(dirErr, fs.ErrNotExist):
		if err := ops.MkdirAll(dir, 0o755); err != nil {
			return failProvision(ops, report, "create category directory", err, false)
		}
		report.CategoryCreated = true
	}

	readable := false
	if _, err := ops.ReadDir(dir); err == nil {
		readable = true
	}
	report.Readable = readable

	probePath := filepath.Join(dir, fmt.Sprintf(".howl-write-probe-%d-%d.tmp", os.Getpid(), time.Now().UnixNano()))
	file, err := ops.OpenFile(probePath)
	if err != nil {
		return failProvision(ops, report, "create probe file", err, readable)
	}
	report.FileCreated = true

	payload := []byte("howl-write-probe")
	n, writeErr := file.Write(payload)
	if writeErr == nil {
		writeErr = file.Sync()
	}
	closeErr := file.Close()
	if writeErr == nil {
		writeErr = closeErr
	}
	if writeErr != nil {
		_ = ops.Remove(probePath)
		return failProvision(ops, report, "write probe file", writeErr, readable)
	}
	report.BytesWritten = n
	if n != len(payload) {
		_ = ops.Remove(probePath)
		return failProvision(ops, report, "write probe file", io.ErrShortWrite, readable)
	}

	got, err := ops.ReadFile(probePath)
	if err != nil {
		_ = ops.Remove(probePath)
		return failProvision(ops, report, "read back probe file", err, readable)
	}
	if !bytes.Equal(got, payload) {
		_ = ops.Remove(probePath)
		return failProvision(ops, report, "verify probe file", errors.New("the probe file did not read back the bytes that were written"), readable)
	}
	report.ReadBack = true

	if err := ops.Remove(probePath); err != nil {
		return failProvision(ops, report, "remove probe file", err, readable)
	}
	report.Cleanup = true

	// Non-destructive cleanup: only directories this call created, and only when
	// they are empty. Provisioning keeps them; a pure probe removes them.
	if !keep && report.CategoryCreated {
		if entries, err := ops.ReadDir(dir); err == nil && len(entries) == 0 {
			if err := ops.Remove(dir); err == nil {
				report.CategoryRemoved = true
			}
		}
	}
	if !keep && report.RootCreated {
		if entries, err := ops.ReadDir(root); err == nil && len(entries) == 0 {
			if err := ops.Remove(root); err == nil {
				report.RootRemoved = true
			}
		}
	}

	report.Writable = true
	report.Readable = true
	report.Status = ContentReady
	return report
}

func failProvision(ops contentFSOps, report ProvisionReport, step string, err error, readable bool) ProvisionReport {
	report.Step = step
	report.Operation = step
	report.Err = err
	report.Readable = readable
	report.Retryable = true

	var blocked contentPathBlocked
	switch {
	case errors.As(err, &blocked):
		report.Status = ContentUnavailable
		report.Reason = fmt.Sprintf("%s is a file, not a directory, so Howl cannot use it for %s.", blocked.Path, report.Category)
		report.Retryable = false
	case errors.Is(err, fs.ErrPermission):
		if readable {
			report.Status = ContentReadOnly
			report.Reason = fmt.Sprintf("This %s library can be read but not written: Windows denied access to %s. Existing content stays visible; new items cannot be created until the folder is writable.", report.Category, report.Path)
		} else {
			report.Status = ContentUnavailable
			report.Reason = fmt.Sprintf("Windows denied access to %s while trying to %s. The %s library was not created.", report.Path, step, report.Category)
		}
	default:
		report.Status = ContentUnavailable
		report.Reason = fmt.Sprintf("Could not %s in %s: %v", step, report.Path, err)
	}
	return report
}

// provisionDefaultRoot is the single provisioning route for a category root.
func (a *App) provisionDefaultRoot(category library.Category, role library.Role) (ProvisionReport, library.Registration) {
	dir, err := paths.DefaultContentDir(string(category))
	if err != nil {
		report := ProvisionReport{
			Category: string(category), Status: ContentUnavailable, Step: "resolve the default location",
			Operation: "resolve the default location", Retryable: true, Err: err,
			Reason: fmt.Sprintf("Could not resolve the default location for %s: %v", category, err),
		}
		return report, library.Registration{}
	}

	report := probeContentDir(a.contentOpsOrDefault(), string(category), dir, true)
	if report.Err != nil {
		return report, library.Registration{}
	}

	reg := a.libraries()
	if reg == nil {
		report.Status = ContentUnavailable
		report.Reason = "The content library registry is not initialised."
		report.Err = errors.New("content library registry not initialised")
		report.Retryable = false
		return report, library.Registration{}
	}
	registration, err := reg.Add(library.Registration{
		Category:    category,
		Path:        dir,
		Role:        library.NormalizeRole(role),
		Origin:      library.OriginHowl,
		Format:      library.FormatHowl,
		DefaultSave: role&library.RoleSave != 0,
	})
	if err != nil {
		report.Status = ContentUnavailable
		report.Reason = fmt.Sprintf("The %s directory was writable, but it could not be registered: %v", category, err)
		report.Err = err
		return report, library.Registration{}
	}
	a.recordProbe(registration.ID, report)
	return report, registration
}

// initializeContentRoots provisions or probes each category once. It never
// fails a boot: outcomes are recorded as feature state and logged once.
func (a *App) initializeContentRoots() {
	reg := a.libraries()
	if reg == nil {
		return
	}
	ops := a.contentOpsOrDefault()
	for _, category := range provisionedContentCategories {
		views := reg.List(category)
		if len(views) == 0 {
			report, _ := a.provisionDefaultRoot(category, library.RoleRead|library.RoleWrite|library.RoleSave)
			a.recordInit(category, report)
			if report.Err != nil {
				a.logErrorSilent("SYS-CONTENT-LIBRARY", apperrors.Medium, "content.init."+string(category), "app-content-roots.go", report.Err)
			}
			continue
		}
		for _, view := range views {
			if !view.Exists {
				continue
			}
			if !(view.Role.CanWrite() || view.Role.CanSave()) {
				continue
			}
			report := probeContentDir(ops, string(category), view.Path, false)
			a.recordProbe(view.ID, report)
			if report.Err != nil {
				a.logErrorSilent("SYS-CONTENT-LIBRARY", apperrors.Medium, "content.probe."+string(category), "app-content-roots.go", report.Err)
			}
		}
	}
}

// EnsureDefaultContentLibrary is the explicit setup action a panel offers when a
// category is not initialized or unavailable.
func (a *App) EnsureDefaultContentLibrary(category string) (ContentLibraryState, error) {
	cat := library.Category(strings.TrimSpace(category))
	if !cat.Valid() {
		return ContentLibraryState{}, fmt.Errorf("unknown content category %q", category)
	}
	report, _ := a.provisionDefaultRoot(cat, library.RoleRead|library.RoleWrite|library.RoleSave)
	a.recordInit(cat, report)
	if report.Err != nil {
		return synthesizedLibraryState(cat, report), report.Err
	}
	readiness := a.GetContentReadiness()
	best := ContentLibraryState{}
	found := false
	for _, state := range readiness.Libraries {
		if state.Category != string(cat) {
			continue
		}
		if !found || (state.Path == report.Path && best.Path != report.Path) {
			best = state
			found = true
		}
	}
	if found {
		return best, nil
	}
	return synthesizedLibraryState(cat, report), nil
}

// GetContentReadiness reports every library and the derived category summaries.
// It performs no filesystem mutation; probe results come from initialization and
// explicit setup.
func (a *App) GetContentReadiness() ContentReadiness {
	reg := a.libraries()
	return buildContentReadiness(reg, a.probeSnapshot(), a.initSnapshot())
}

func buildContentReadiness(reg *library.Registry, probes, initFails map[string]ProvisionReport) ContentReadiness {
	out := ContentReadiness{Libraries: []ContentLibraryState{}, Categories: []ContentCategoryState{}}
	if reg == nil {
		return out
	}

	present := map[library.Category]bool{}
	for _, view := range reg.List("") {
		probe := ProvisionReport{}
		if stored, ok := probes[view.ID]; ok {
			probe = stored
		}
		state := libraryStateFromView(view, probe)
		// SaveRoot means the designated default destination, not merely a write
		// role: several libraries may be save-capable, but only one receives new
		// documents. Create is gated on that one.
		if view.Role.CanSave() && view.Exists {
			if saveRoot, ok := reg.DefaultSave(view.Category); ok && saveRoot.ID == view.ID {
				state.SaveRoot = true
			}
		}
		out.Libraries = append(out.Libraries, state)
		present[view.Category] = true
	}

	for _, category := range provisionedContentCategories {
		if present[category] {
			continue
		}
		out.Libraries = append(out.Libraries, synthesizedLibraryState(category, initFails[string(category)]))
	}

	for _, category := range orderedContentCategories(out.Libraries) {
		out.Categories = append(out.Categories, rollupContentCategory(category, out.Libraries))
	}
	return out
}

func libraryStateFromView(view library.View, probe ProvisionReport) ContentLibraryState {
	state := ContentLibraryState{
		Category:  string(view.Category),
		LibraryID: view.ID,
		Path:      view.Path,
		Readable:  view.Readable && view.Exists,
		Writable:  view.Writable && view.Exists,
	}
	if probe.Path != "" {
		if probe.Readable {
			state.Readable = true
		}
		if probe.Err != nil {
			// The probe is authoritative about what it actually observed: a
			// readable-but-unwritable library stays readable, while a library the
			// probe could not enumerate is reported unreadable, never "ready".
			state.Readable = probe.Readable
			state.Writable = false
			state.Status = probe.Status
			state.Reason = probe.Reason
			state.Operation = probe.Operation
			state.Retryable = probe.Retryable
			return state
		}
		if probe.Writable {
			state.Writable = view.Role.CanWrite() || view.Role.CanSave()
		}
	}

	switch {
	case !state.Readable:
		state.Status = ContentNotReadable
		state.Reason = fmt.Sprintf("The registered %s library at %s cannot be read. Only cached items are shown.", state.Category, state.Path)
		state.Retryable = true
	case !state.Writable:
		state.Status = ContentReadOnly
		state.Reason = fmt.Sprintf("This %s library can be read but not written: %s does not grant write access. Existing content stays visible; new items cannot be created here.", state.Category, state.Path)
		state.Retryable = true
	default:
		state.Status = ContentReady
	}
	return state
}

func synthesizedLibraryState(category library.Category, failure ProvisionReport) ContentLibraryState {
	state := ContentLibraryState{
		Category:  string(category),
		Path:      failure.Path,
		Status:    ContentNotInitialized,
		Reason:    fmt.Sprintf("No %s library has been set up yet. Set up the default location to start saving.", category),
		Retryable: true,
	}
	if state.Path == "" {
		if dir, err := paths.DefaultContentDir(string(category)); err == nil {
			state.Path = dir
		}
	}
	if failure.Err != nil {
		state.Status = failure.Status
		if state.Status == "" {
			state.Status = ContentUnavailable
		}
		state.Reason = failure.Reason
		state.Operation = failure.Operation
		state.Retryable = failure.Retryable
	}
	return state
}

// rollupContentCategory derives a summary but never overwrites per-library
// truth: the worst library decides the summary status and FailingLibraries counts
// every library that is not ready.
func rollupContentCategory(category string, libraries []ContentLibraryState) ContentCategoryState {
	out := ContentCategoryState{Category: category}
	worst := -1
	firstReason := ""
	for _, state := range libraries {
		if state.Category != category {
			continue
		}
		out.Libraries++
		if state.Readable {
			out.Readable = true
		}
		if state.Writable {
			out.Writable = true
			if state.SaveRoot {
				out.SaveWritable = true
			}
		}
		if state.Status != ContentReady {
			out.FailingLibraries++
			if firstReason == "" {
				firstReason = state.Reason
			}
		}
		if rank := contentStatusRank(state.Status); rank > worst {
			worst = rank
			out.Status = state.Status
		}
	}
	if out.Libraries == 0 {
		out.Status = ContentNotInitialized
	}
	if out.FailingLibraries > 0 {
		if out.Libraries == 1 {
			out.Reason = firstReason
		} else {
			out.Reason = fmt.Sprintf("%d of %d %s libraries need attention. %s", out.FailingLibraries, out.Libraries, category, firstReason)
		}
	}
	return out
}

func contentStatusRank(status string) int {
	switch status {
	case ContentUnavailable:
		return 4
	case ContentNotReadable:
		return 3
	case ContentNotInitialized:
		return 2
	case ContentReadOnly:
		return 1
	case ContentReady:
		return 0
	default:
		return 3
	}
}

func orderedContentCategories(libraries []ContentLibraryState) []string {
	order := []string{}
	seen := map[string]bool{}
	for _, category := range provisionedContentCategories {
		order = append(order, string(category))
		seen[string(category)] = true
	}
	for _, state := range libraries {
		if seen[state.Category] {
			continue
		}
		seen[state.Category] = true
		order = append(order, state.Category)
	}
	return order
}

// requireCreateWritable refuses a create before it reaches the store when the
// default save root cannot be written. The store remains the authority for
// updates; this gate exists so a blocked create names its cause instead of
// surfacing a low-level filesystem error.
func (a *App) requireCreateWritable(category library.Category, label string) error {
	reg := a.libraries()
	if reg == nil {
		return fmt.Errorf("%s: the content library registry is not initialised", label)
	}
	saveRoot, ok := reg.DefaultSave(category)
	if !ok {
		return fmt.Errorf("%s: no %s library is set up yet; set up the default location before saving. Nothing was written", label, category)
	}
	if probe, probed := a.probeFor(saveRoot.ID); probed && !probe.Writable {
		return fmt.Errorf("%s: %s Nothing was written", label, probe.Reason)
	}
	info, err := os.Stat(saveRoot.Path)
	if err != nil || !info.IsDir() {
		return fmt.Errorf("%s: the default %s library at %s is unavailable; nothing was written", label, category, saveRoot.Path)
	}
	return nil
}

func (a *App) libraries() *library.Registry {
	if a == nil || a.registries == nil {
		return nil
	}
	return a.registries.Libraries()
}

func (a *App) contentOpsOrDefault() contentFSOps {
	if a.contentOps != nil {
		return *a.contentOps
	}
	return defaultContentFSOps()
}

func (a *App) recordProbe(libraryID string, report ProvisionReport) {
	if libraryID == "" {
		return
	}
	a.contentMu.Lock()
	defer a.contentMu.Unlock()
	if a.contentProbes == nil {
		a.contentProbes = map[string]ProvisionReport{}
	}
	a.contentProbes[libraryID] = report
}

func (a *App) probeFor(libraryID string) (ProvisionReport, bool) {
	a.contentMu.Lock()
	defer a.contentMu.Unlock()
	report, ok := a.contentProbes[libraryID]
	return report, ok
}

func (a *App) recordInit(category library.Category, report ProvisionReport) {
	a.contentMu.Lock()
	defer a.contentMu.Unlock()
	if a.contentInit == nil {
		a.contentInit = map[string]ProvisionReport{}
	}
	a.contentInit[string(category)] = report
}

func (a *App) probeSnapshot() map[string]ProvisionReport {
	a.contentMu.Lock()
	defer a.contentMu.Unlock()
	out := make(map[string]ProvisionReport, len(a.contentProbes))
	for id, report := range a.contentProbes {
		out[id] = report
	}
	return out
}

func (a *App) initSnapshot() map[string]ProvisionReport {
	a.contentMu.Lock()
	defer a.contentMu.Unlock()
	out := make(map[string]ProvisionReport, len(a.contentInit))
	for category, report := range a.contentInit {
		out[category] = report
	}
	return out
}
