package paths

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// This file is the PathService responsibility defined by ADR-001 §4 and contracts
// A and B: root set, mode resolution, category→location, class labelling,
// containment checks and diagnostics.
//
// It is additive. Nothing here changes where the running application reads or
// writes a byte; the legacy functions (ProjectRoot, StateRoot, StatePath,
// InternalConfigDir) keep their existing behaviour until the migration sequence
// in ADR-001 §12 is authorized as its own task.

// Mode is the storage mode. ADR-001 D2 defers portable mode; it is absent by
// decision, not by omission.
type Mode string

const (
	ModeInstalled Mode = "installed"
	ModeIsolated  Mode = "isolated"
)

// RootClass labels what a root may hold. Every path the service returns carries
// exactly one class (ADR-001 A3).
type RootClass string

const (
	ClassInstallation RootClass = "installation"
	ClassAppData      RootClass = "application-data"
	ClassUserContent  RootClass = "user-content"
	ClassExternal     RootClass = "external-resource"
	ClassTemporary    RootClass = "temporary"
)

// Root identifiers. Content categories use "content:<category>".
const (
	RootInstallation    = "installation"
	RootApplicationData = "application-data"
	RootTemporary       = "temporary"
)

// Root is one classified location.
type Root struct {
	ID       string
	Class    RootClass
	Path     string
	Writable bool
}

// RootSet is the resolved layout for the active mode.
type RootSet struct {
	Mode  Mode
	roots map[string]Root
	cats  map[string]Root
}

var (
	// ErrUnknownCategory means the category is not a content category.
	ErrUnknownCategory = errors.New("paths: unknown content category")
	// ErrNotContained means the candidate is not inside the root.
	ErrNotContained = errors.New("paths: path escapes its root")
	// ErrUnavailable means a root could not be resolved.
	ErrUnavailable = errors.New("paths: root unavailable")
)

// RootError is a typed failure that names the root or category and the remedy
// (ADR-001 A6). It never relocates silently.
type RootError struct {
	Root     string
	Category string
	Path     string
	Remedy   string
	Err      error
}

func (e *RootError) Error() string {
	var b strings.Builder
	b.WriteString("paths: ")
	if e.Root != "" {
		fmt.Fprintf(&b, "%s: ", e.Root)
	}
	if e.Category != "" {
		fmt.Fprintf(&b, "category %q: ", e.Category)
	}
	b.WriteString(e.Err.Error())
	if e.Path != "" {
		fmt.Fprintf(&b, " (%s)", e.Path)
	}
	if e.Remedy != "" {
		fmt.Fprintf(&b, "; %s", e.Remedy)
	}
	return b.String()
}

func (e *RootError) Unwrap() error { return e.Err }

// contentCategories lists the categories the service resolves. Images and videos
// are separate roots because they live under Pictures and Videos.
func contentCategories() []string {
	return []string{
		ContentChats,
		ContentCharacters,
		ContentWorlds,
		ContentScenarios,
		ContentLorebooks,
		ContentImages,
		ContentVideos,
	}
}

// Roots resolves the root set for the active mode.
//
// Installed: installation is the project root (read-only resources), application
// data is %AppData%\Howl, temporary is the OS temp directory, and each content
// category resolves through the Windows known-folder API.
//
// Isolated: application data, temporary and content all resolve inside the
// isolation root.
//
// Filesystem effects. In installed mode Roots performs no creation and no
// modification. In isolated mode the existing isolation resolver creates its own
// layout (data, temp, webview2) before returning, so the first resolution does
// create directories. That write belongs to isolation setup, not to this service;
// it is recorded as finding I-02 in the TASK-01 record rather than silently
// corrected here, because changing it would change established isolation
// behaviour.
func Roots() (RootSet, error) {
	iso, err := ResolveIsolation()
	if err != nil {
		return RootSet{}, &RootError{
			Root:   RootApplicationData,
			Remedy: "unset " + IsolationEnvVar + " or point it at a valid isolation root",
			Err:    ErrUnavailable,
		}
	}

	rs := RootSet{
		Mode:  ModeInstalled,
		roots: make(map[string]Root, 3),
		cats:  make(map[string]Root, len(contentCategories())),
	}

	// R1 is a physical classification, never a storage authority.
	rs.roots[RootInstallation] = Root{
		ID:       RootInstallation,
		Class:    ClassInstallation,
		Path:     ProjectRoot(),
		Writable: false,
	}

	if iso != nil {
		rs.Mode = ModeIsolated
		rs.roots[RootApplicationData] = Root{
			ID: RootApplicationData, Class: ClassAppData, Path: iso.StateRoot, Writable: true,
		}
		rs.roots[RootTemporary] = Root{
			ID: RootTemporary, Class: ClassTemporary, Path: iso.TempDir, Writable: true,
		}
	} else {
		base, cfgErr := os.UserConfigDir()
		if cfgErr != nil {
			return RootSet{}, &RootError{
				Root:   RootApplicationData,
				Remedy: "the platform user-configuration directory is unavailable; storage cannot be resolved",
				Err:    ErrUnavailable,
			}
		}
		rs.roots[RootApplicationData] = Root{
			ID: RootApplicationData, Class: ClassAppData, Path: filepath.Join(base, "Howl"), Writable: true,
		}
		rs.roots[RootTemporary] = Root{
			ID: RootTemporary, Class: ClassTemporary, Path: os.TempDir(), Writable: true,
		}
	}

	for _, category := range contentCategories() {
		dir, catErr := DefaultContentDir(category)
		if catErr != nil {
			return RootSet{}, &RootError{
				Category: category,
				Remedy:   "the known folder for this category could not be resolved; the save is refused rather than relocated",
				Err:      ErrUnavailable,
			}
		}
		rs.cats[category] = Root{
			ID:       "content:" + category,
			Class:    ClassUserContent,
			Path:     dir,
			Writable: true,
		}
	}

	return rs, nil
}

// All returns every resolved location, ordered by ID, for diagnostics and tests.
func (rs RootSet) All() []Root {
	out := make([]Root, 0, len(rs.roots)+len(rs.cats))
	for _, r := range rs.roots {
		out = append(out, r)
	}
	for _, r := range rs.cats {
		out = append(out, r)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}

// Get returns a root by identifier.
func (rs RootSet) Get(id string) (Root, bool) {
	r, ok := rs.roots[id]
	return r, ok
}

// Category returns the resolved location for one content category.
func (rs RootSet) Category(category string) (Root, bool) {
	r, ok := rs.cats[category]
	return r, ok
}

// Describe renders the resolved layout deterministically. This is the diagnostic
// surface required by ADR-001 A1/B5: a reader can see mode and roots without
// running the application.
func (rs RootSet) Describe() string {
	var b strings.Builder
	fmt.Fprintf(&b, "mode: %s\n", rs.Mode)
	for _, r := range rs.All() {
		writable := "read-only"
		if r.Writable {
			writable = "writable"
		}
		fmt.Fprintf(&b, "%-20s %-18s %-9s %s\n", r.ID, r.Class, writable, r.Path)
	}
	return b.String()
}

// ExternalRoot classifies a user-selected location — a model directory, a
// library root — as read-only external data (ADR-001 R4).
//
// It confers no write permission. Writing into an external root requires a
// per-operation grant from ExternalDataGuard (ADR-001 B3, contract F), which is
// not part of this package. A caller that treats a writable Root or a
// successful ExternalRoot call as authorization is wrong.
func ExternalRoot(path string) (Root, error) {
	if path == "" {
		return Root{}, &RootError{Path: path, Remedy: "external locations must be supplied by the user", Err: ErrNotContained}
	}
	if !filepath.IsAbs(path) {
		return Root{}, &RootError{
			Path:   path,
			Remedy: "an external location must be absolute; a relative path would be resolved against the working directory",
			Err:    ErrNotContained,
		}
	}
	clean := filepath.Clean(path)
	info, err := os.Stat(clean)
	if err != nil {
		return Root{}, &RootError{Path: clean, Remedy: "the external location does not exist or is not reachable", Err: ErrUnavailable}
	}
	if !info.IsDir() {
		return Root{}, &RootError{Path: clean, Remedy: "the external location is not a directory", Err: ErrUnavailable}
	}
	return Root{
		ID:       "external:" + filepath.Base(clean),
		Class:    ClassExternal,
		Path:     clean,
		Writable: false,
	}, nil
}

// CheckGrant is the containment test a scoped grant performs: the candidate must
// be strictly inside the granted destination root. Equality is refused, because a
// grant authorises one operation's subtree, not the destination root itself.
func CheckGrant(destinationRoot, candidate string) error {
	if err := CheckContained(destinationRoot, candidate); err != nil {
		return err
	}
	rootClean := filepath.Clean(destinationRoot)
	candidateClean := filepath.Clean(candidate)
	if containsPath(rootClean, candidateClean) && candidateClean == rootClean {
		return &RootError{
			Path:   candidateClean,
			Remedy: "a grant covers one operation subtree, not the destination root itself",
			Err:    ErrNotContained,
		}
	}
	return nil
}

// Check reports whether candidate is inside the resolved location for a
// category. Unknown categories are reported as such, never defaulted.
func (rs RootSet) Check(category, candidate string) error {
	root, ok := rs.cats[category]
	if !ok {
		return &RootError{Category: category, Err: ErrUnknownCategory}
	}
	return CheckContained(root.Path, candidate)
}

// CheckContained reports whether candidate is inside root.
//
// Containment is decided on **physical** paths: reparse points are resolved for
// the part of the path that exists, and the not-yet-created remainder is
// re-appended before comparison (ADR-001 B2).
//
// Identity comparison alone is not sufficient, and the junction test proves it:
// walking up lexical ancestors from `root\esc\secret` reaches `root`, whose
// identity matches, even though the path physically lives outside the root. The
// physical path must be resolved first; volume serial and file index then remain
// the basis for the volume-level checks.
//
// This is a pre-flight check, not a TOCTOU guarantee. A write must still use
// handle-relative operations that preserve the validated directory identity
// (ADR-001 §7); otherwise a reparse point created after this call is not
// detected. When a path cannot be resolved, containment fails closed.
func CheckContained(root, candidate string) error {
	if root == "" || candidate == "" {
		return &RootError{Path: candidate, Remedy: "both root and candidate are required", Err: ErrNotContained}
	}
	if !filepath.IsAbs(root) {
		return &RootError{Path: root, Remedy: "containment roots must be absolute", Err: ErrNotContained}
	}
	if !filepath.IsAbs(candidate) {
		return &RootError{
			Path:   candidate,
			Remedy: "a relative path is resolved against the working directory, which is never a storage authority",
			Err:    ErrNotContained,
		}
	}

	rootClean := filepath.Clean(root)
	candidateClean := filepath.Clean(candidate)

	rootPhysical, err := physicalPath(rootClean)
	if err != nil {
		return &RootError{
			Path:   rootClean,
			Remedy: "the root cannot be resolved to a physical path; containment fails closed",
			Err:    ErrNotContained,
		}
	}
	candidatePhysical, err := physicalPath(candidateClean)
	if err != nil {
		return &RootError{
			Path:   candidateClean,
			Remedy: "the candidate cannot be resolved to a physical path; containment fails closed",
			Err:    ErrNotContained,
		}
	}

	if containsPath(rootPhysical, candidatePhysical) {
		return nil
	}
	return &RootError{
		Path:   candidateClean,
		Remedy: "the resolved location is outside its category root",
		Err:    ErrNotContained,
	}
}

// physicalPath resolves reparse points in the portion of path that exists and
// re-appends the remainder that does not exist yet, so containment is decided on
// where the path really is rather than on how it is spelled.
func physicalPath(path string) (string, error) {
	existing := filepath.Clean(path)
	var remainder []string
	for {
		resolved, err := resolvePhysical(existing)
		if err == nil {
			parts := append([]string{resolved}, remainder...)
			return filepath.Join(parts...), nil
		}
		parent := filepath.Dir(existing)
		if parent == existing {
			return "", fmt.Errorf("cannot resolve %q to a physical path", path)
		}
		remainder = append([]string{filepath.Base(existing)}, remainder...)
		existing = parent
	}
}
