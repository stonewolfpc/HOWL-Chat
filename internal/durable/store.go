// Package durable provides the single durable-document primitive defined by
// ADR-001 §8.
//
// A document write is: acquire the exclusive lock, write a temporary file,
// flush it, rotate pre-images, commit a manifest describing the intended state,
// replace the live document by rename, then drop pre-images beyond N.
//
// The live path is never opened for writing. Its only mutation is a rename, so a
// concurrent reader observes either the previous document or the new one, never a
// partially written file. That property is what makes this different from
// os.WriteFile, which truncates in place and is observable mid-write.
//
// Reader contract (ADR-001 §8.1, decided in
// TASK-02-READER-SHARING-DECISION-2026-11-22.md): readers open, read and close
// within one call. They must not retain a handle, memory-map a document, request a
// custom share mode, or open a pre-image. A reader that holds the live document
// blocks replacement; the writer retries and then reports NotSaved(rename_denied)
// rather than forcing the write.
package durable

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Outcome is the result classification of a store operation. It is reported, not
// thrown: callers distinguish Busy from NotSaved from Recovered.
type Outcome string

const (
	// Saved means the document was committed.
	Saved Outcome = "Saved"
	// Absent means no document exists yet. Not a failure.
	Absent Outcome = "Absent"
	// Recovered means the live document was unusable and bytes came from a
	// pre-image whose hash matched the manifest.
	Recovered Outcome = "Recovered"
	// NotSaved means the write did not commit and the previous document is intact.
	NotSaved Outcome = "NotSaved"
	// RefusedUnknownSchema means the document is newer than this build supports.
	// It was not rewritten.
	RefusedUnknownSchema Outcome = "RefusedUnknownSchema"
	// Busy means the exclusive lock could not be acquired within the timeout. No
	// bypass exists.
	Busy Outcome = "Busy"
)

// Result reports what happened. Reason is a short machine-readable detail, not a
// user-facing message.
type Result struct {
	Outcome   Outcome
	Reason    string
	From      string // pre-image name, when Recovered
	Schema    int    // schema of the document involved
	Supported int    // highest schema this store accepts
}

// OK reports whether the operation produced usable bytes. It is true for a
// Recovered result, which is usable but is an earlier generation rather than the
// current document; use Current to make that distinction.
func (r Result) OK() bool { return r.Outcome == Saved || r.Outcome == Recovered }

// Current reports whether the bytes returned are the current committed document.
//
// A Recovered result is deliberately *not* current: the live document was
// unusable (or a commit was observed mid-flight) and the bytes came from a
// retained generation, which may be older than the last successful write.
// Callers that must not act on a superseded document use this rather than OK.
func (r Result) Current() bool { return r.Outcome == Saved }

var (
	// ErrBusy is returned when the exclusive lock is not acquired in time.
	ErrBusy = errors.New("durable: store is busy")
	// ErrRenameDenied is returned when an open reader prevented replacement.
	ErrRenameDenied = errors.New("durable: replacement denied by an open reader")
	// ErrCorrupt is returned when the live document and every pre-image fail
	// their recorded hash.
	ErrCorrupt = errors.New("durable: document is corrupt and no valid pre-image exists")
	// ErrSchemaNewer is returned when the document's schema exceeds what is
	// supported.
	ErrSchemaNewer = errors.New("durable: document schema is newer than supported")
)

// Options tune the store. The zero value is not usable; use DefaultOptions.
type Options struct {
	// LockTimeout is how long to wait for the exclusive lock before reporting
	// Busy (ADR-001 §8: default 5000ms).
	LockTimeout time.Duration
	// RenameAttempts is the number of replacement attempts made inside the held
	// lock (ADR-001 §8: default 3).
	RenameAttempts int
	// RenameWait is the pause between replacement attempts (default 50ms).
	RenameWait time.Duration
	// PreImages is how many previous generations are retained (default 2).
	PreImages int
	// SupportedSchema is the highest schema this build accepts. A document with a
	// higher schema is refused, never rewritten.
	SupportedSchema int
}

// DefaultOptions returns the ADR-001 §8 timers with two retained generations.
func DefaultOptions() Options {
	return Options{
		LockTimeout:     5 * time.Second,
		RenameAttempts:  3,
		RenameWait:      50 * time.Millisecond,
		PreImages:       2,
		SupportedSchema: 1,
	}
}

// Store is a DurableStore over one document path.
type Store struct {
	path string
	opts Options
}

// Open returns a store for path. No I/O happens until a method is called.
func Open(path string, opts Options) *Store {
	if opts.LockTimeout <= 0 {
		opts.LockTimeout = 5 * time.Second
	}
	if opts.RenameAttempts <= 0 {
		opts.RenameAttempts = 3
	}
	if opts.RenameWait <= 0 {
		opts.RenameWait = 50 * time.Millisecond
	}
	if opts.PreImages <= 0 {
		opts.PreImages = 2
	}
	return &Store{path: path, opts: opts}
}

// Path is the document this store owns.
func (s *Store) Path() string { return s.path }

// LockPath is the lock file guarding this document. It is a stable name next to
// the document so that every writer contends on the same object.
func (s *Store) LockPath() string { return s.path + ".lock" }

func (s *Store) manifestPath() string { return s.path + ".manifest.json" }
func (s *Store) preImagePath(n int) string {
	return fmt.Sprintf("%s.v%d.bak", s.path, n)
}

// tempPrefix is the temp-name prefix; sweepTemps removes leftovers from writers
// that died before replacement.
func (s *Store) tempPrefix() string {
	return "." + filepath.Base(s.path) + ".tmp-"
}

// testHookBeforeReplace is nil in production. Tests set it to simulate a process
// death between the manifest commit and the replacement, which is otherwise
// impossible to schedule deterministically.
var testHookBeforeReplace func()

type preImageEntry struct {
	Name        string `json:"name"`
	SHA256      string `json:"sha256"`
	Schema      int    `json:"schemaVersion"`
	CommittedAt string `json:"committedAt"`
}

type manifest struct {
	Schema      int             `json:"schemaVersion"`
	SHA256      string          `json:"sha256"`
	Size        int64           `json:"size"`
	CommittedAt string          `json:"committedAt"`
	PreImages   []preImageEntry `json:"preImages"`
}

func hashOf(b []byte) string {
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:])
}

// Save commits doc as the live document.
//
// Every writer takes the exclusive lock; there is no privileged path. On any
// failure the previous document is left intact and the outcome says why.
func (s *Store) Save(doc []byte, schema int) Result {
	if schema > s.opts.SupportedSchema {
		return Result{
			Outcome:   RefusedUnknownSchema,
			Reason:    "schema_newer_than_supported",
			Schema:    schema,
			Supported: s.opts.SupportedSchema,
		}
	}
	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return Result{Outcome: NotSaved, Reason: "mkdir: " + err.Error()}
	}

	lock, err := acquireLock(s.LockPath(), s.opts.LockTimeout)
	if err != nil {
		if errors.Is(err, ErrBusy) {
			return Result{
				Outcome: Busy,
				Reason:  "lock_not_acquired_within_" + s.opts.LockTimeout.String(),
				Schema:  schema,
			}
		}
		return Result{Outcome: NotSaved, Reason: "lock: " + err.Error(), Schema: schema}
	}
	defer lock.release()

	// Under the lock, no other writer is active, so leftovers are from dead
	// processes.
	s.sweepTemps()

	// The manifest that described the previous state, kept so a failed commit can
	// restore it rather than leaving the manifest pointing at a document that was
	// never installed.
	prevManifest, hadPrevManifest := s.readManifestBytes()
	prevSchema := 0
	if hadPrevManifest {
		prevSchema = s.decodeSchema(prevManifest)
	}

	committedAt := time.Now().UTC().Format(time.RFC3339Nano)

	tmpPath, err := s.writeTemp(doc)
	if err != nil {
		return Result{Outcome: NotSaved, Reason: "write_temp: " + err.Error(), Schema: schema}
	}

	entries, staged, err := s.stagePreImages(prevSchema, committedAt)
	if err != nil {
		_ = os.Remove(tmpPath)
		return Result{Outcome: NotSaved, Reason: "pre_image: " + err.Error(), Schema: schema}
	}

	if err := s.writeManifest(manifest{
		Schema:      schema,
		SHA256:      hashOf(doc),
		Size:        int64(len(doc)),
		CommittedAt: committedAt,
		PreImages:   entries,
	}); err != nil {
		s.abortStaged(staged)
		_ = os.Remove(tmpPath)
		if errors.Is(err, ErrRenameDenied) {
			// The manifest describes the previous state and was never replaced.
			return Result{Outcome: NotSaved, Reason: "rename_denied_by_open_reader", Schema: schema}
		}
		return Result{Outcome: NotSaved, Reason: "manifest: " + err.Error(), Schema: schema}
	}

	if testHookBeforeReplace != nil {
		testHookBeforeReplace()
	}

	if err := s.replaceWithRetry(tmpPath, s.path); err != nil {
		_ = os.Remove(tmpPath)
		// The staged generations were never promoted, so the retained files are
		// still exactly what the previous manifest describes. Discard the staging
		// and put that manifest back.
		s.abortStaged(staged)
		s.restoreManifest(prevManifest, hadPrevManifest)
		reason := "rename: " + err.Error()
		if errors.Is(err, ErrRenameDenied) {
			reason = "rename_denied_by_open_reader"
		}
		return Result{Outcome: NotSaved, Reason: reason, Schema: schema}
	}

	// Promote only now that the live document holds the new content. Promoting
	// earlier would consume a generation for a commit that never happened.
	reason := ""
	if err := s.promoteStaged(staged); err != nil {
		reason = "pre_image_promote_partial"
	}

	// Pre-images may be dropped only after the rename succeeded (ADR-001 §8).
	s.trimPreImages()

	return Result{Outcome: Saved, Schema: schema, Reason: reason}
}

// Load returns the committed document.
//
// The read is a single short-lived open/read/close, which is exactly the reader
// contract this store depends on. Integrity is checked against the manifest; a
// document that fails its hash is recovered from the newest pre-image that
// matches, never silently replaced with empty content.
func (s *Store) Load() ([]byte, Result) {
	// A read touches two files (document and manifest). A commit between those
	// two reads would otherwise look like corruption, so the pair is re-read a
	// bounded number of times before any recovery is attempted.
	const pairAttempts = 3

	var data []byte
	var m manifest
	var manifestErr error

	for attempt := 0; attempt < pairAttempts; attempt++ {
		d, err := os.ReadFile(s.path)
		if err != nil {
			if os.IsNotExist(err) {
				return s.loadAbsent()
			}
			return nil, Result{Outcome: NotSaved, Reason: "read: " + err.Error()}
		}
		data = d

		mm, ok, err := s.readManifestStrict()
		if err != nil {
			manifestErr = err
			time.Sleep(2 * time.Millisecond)
			continue
		}
		if !ok {
			// Unmanaged document: nothing recorded to verify it against. The
			// manifest appears on the next Save.
			return data, Result{Outcome: Saved, Reason: "unmanaged_no_manifest"}
		}
		m = mm
		if m.Schema > s.opts.SupportedSchema {
			return data, Result{
				Outcome:   RefusedUnknownSchema,
				Reason:    "schema_newer_than_supported",
				Schema:    m.Schema,
				Supported: s.opts.SupportedSchema,
			}
		}
		if hashOf(data) == m.SHA256 {
			return data, Result{Outcome: Saved, Schema: m.Schema}
		}
		time.Sleep(2 * time.Millisecond)
	}

	if manifestErr != nil {
		return nil, Result{Outcome: NotSaved, Reason: "manifest_unreadable: " + manifestErr.Error()}
	}

	// The document matches no committed manifest and the mismatch persisted.
	// If it is byte-identical to a retained generation, an interrupted commit is
	// being observed: the previous document is intact and usable, so report the
	// rollback rather than claiming corruption.
	if findPreImageHash(m, hashOf(data)) {
		return data, Result{
			Outcome: Recovered,
			Reason:  "interrupted_commit_rolled_back",
			Schema:  m.Schema,
		}
	}

	if doc, from, ok := s.recoverFromPreImages(m); ok {
		return doc, Result{Outcome: Recovered, From: from, Schema: m.Schema, Reason: "live_hash_mismatch"}
	}

	return nil, Result{
		Outcome: NotSaved,
		Reason:  "corrupt_no_valid_preimage",
		Schema:  m.Schema,
	}
}

// loadAbsent handles a missing live document: recover from a pre-image when one
// matches, otherwise report Absent.
func (s *Store) loadAbsent() ([]byte, Result) {
	m, ok, err := s.readManifestStrict()
	if err != nil || !ok {
		return nil, Result{Outcome: Absent, Reason: "no_document"}
	}
	if m.Schema > s.opts.SupportedSchema {
		return nil, Result{
			Outcome:   RefusedUnknownSchema,
			Reason:    "schema_newer_than_supported",
			Schema:    m.Schema,
			Supported: s.opts.SupportedSchema,
		}
	}
	if doc, from, ok := s.recoverFromPreImages(m); ok {
		return doc, Result{Outcome: Recovered, From: from, Schema: m.Schema, Reason: "live_missing"}
	}
	return nil, Result{Outcome: Absent, Reason: "no_document"}
}

// writeTemp writes the candidate document beside the live one and flushes it.
// Flushing is what makes the later rename publish content that reached the
// filesystem rather than a page cache entry. It does not by itself establish
// power-loss durability, which remains unproven.
func (s *Store) writeTemp(doc []byte) (string, error) {
	dir := filepath.Dir(s.path)
	f, err := os.CreateTemp(dir, s.tempPrefix())
	if err != nil {
		return "", err
	}
	name := f.Name()
	cleanup := func() {
		_ = f.Close()
		_ = os.Remove(name)
	}
	if _, err := f.Write(doc); err != nil {
		cleanup()
		return "", err
	}
	if err := f.Sync(); err != nil {
		cleanup()
		return "", err
	}
	if err := f.Close(); err != nil {
		_ = os.Remove(name)
		return "", err
	}
	return name, nil
}

// stagePreImages prepares the previous generations for this commit **without
// touching any retained generation**.
//
// The files are written beside their final names and promoted by promoteStaged
// only after the document rename succeeds. That ordering is what makes a failed
// commit harmless: overwriting v1 before the rename would destroy the newest
// retained generation for a commit that never happened, leaving the restored
// manifest describing files that no longer hold those bytes, and recovery would
// then report corruption on a chain that was intact moments earlier.
func (s *Store) stagePreImages(prevSchema int, committedAt string) ([]preImageEntry, []string, error) {
	type generation struct {
		data        []byte
		schema      int
		committedAt string
	}

	var generations []generation
	if live, err := os.ReadFile(s.path); err == nil {
		generations = append(generations, generation{data: live, schema: prevSchema, committedAt: committedAt})
	}
	for i := 1; i <= s.opts.PreImages; i++ {
		if len(generations) >= s.opts.PreImages {
			break
		}
		data, err := os.ReadFile(s.preImagePath(i))
		if err != nil {
			break
		}
		schema := 0
		if m, ok := s.readManifest(); ok && i-1 < len(m.PreImages) {
			schema = m.PreImages[i-1].Schema
		}
		generations = append(generations, generation{data: data, schema: schema, committedAt: committedAt})
	}

	entries := make([]preImageEntry, 0, len(generations))
	staged := make([]string, 0, len(generations))
	for i, g := range generations {
		name := s.preImagePath(i + 1)
		stagedName := name + stagedSuffix
		if err := os.WriteFile(stagedName, g.data, 0o644); err != nil {
			s.abortStaged(staged)
			return nil, nil, err
		}
		staged = append(staged, stagedName)
		entries = append(entries, preImageEntry{
			Name:        filepath.Base(name),
			SHA256:      hashOf(g.data),
			Schema:      g.schema,
			CommittedAt: g.committedAt,
		})
	}
	return entries, staged, nil
}

// stagedSuffix marks a prepared generation that is not yet a retained one.
const stagedSuffix = ".staged"

// promoteStaged moves prepared generations into their retained names. It runs
// only after the live document holds the new content.
func (s *Store) promoteStaged(staged []string) error {
	var firstErr error
	for _, name := range staged {
		final := strings.TrimSuffix(name, stagedSuffix)
		if err := os.Rename(name, final); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}

// abortStaged discards prepared generations. A commit that did not happen leaves
// the retained chain untouched.
func (s *Store) abortStaged(staged []string) {
	for _, name := range staged {
		_ = os.Remove(name)
	}
}

// trimPreImages drops generations beyond the retention window. Called only after
// a successful replacement.
func (s *Store) trimPreImages() {
	for i := s.opts.PreImages + 1; i <= s.opts.PreImages+4; i++ {
		_ = os.Remove(s.preImagePath(i))
	}
}

// replaceWithRetry installs the temp file as the live document. A reader holding
// the live document blocks this, so the attempts are bounded and the failure is
// reported rather than forced.
func (s *Store) replaceWithRetry(tmpPath, dest string) error {
	var last error
	for attempt := 0; attempt < s.opts.RenameAttempts; attempt++ {
		if err := os.Rename(tmpPath, dest); err == nil {
			return nil
		} else {
			last = err
		}
		if attempt < s.opts.RenameAttempts-1 {
			time.Sleep(s.opts.RenameWait)
		}
	}
	return fmt.Errorf("%w: %v", ErrRenameDenied, last)
}

// sweepTemps removes temporary files left by writers that died before replacing.
// Only called while holding the lock.
func (s *Store) sweepTemps() {
	dir := filepath.Dir(s.path)
	matches, err := filepath.Glob(filepath.Join(dir, s.tempPrefix()+"*"))
	if err == nil {
		for _, m := range matches {
			_ = os.Remove(m)
		}
	}
	// A writer that died between staging and promotion leaves prepared
	// generations behind. They are not retained state and are safe to drop.
	stagedMatches, err := filepath.Glob(filepath.Join(dir, filepath.Base(s.path)+".v*"+stagedSuffix))
	if err == nil {
		for _, m := range stagedMatches {
			_ = os.Remove(m)
		}
	}
}

func (s *Store) readManifest() (manifest, bool) {
	m, ok, err := s.readManifestStrict()
	if err != nil {
		return manifest{}, false
	}
	return m, ok
}

// readManifestStrict distinguishes "no manifest" from "manifest present but
// unreadable", because the first is a normal unmanaged document and the second is
// an error the caller must not mistake for one.
func (s *Store) readManifestStrict() (manifest, bool, error) {
	raw, err := os.ReadFile(s.manifestPath())
	if err != nil {
		if os.IsNotExist(err) {
			return manifest{}, false, nil
		}
		return manifest{}, false, err
	}
	var m manifest
	if err := json.Unmarshal(raw, &m); err != nil {
		return manifest{}, false, err
	}
	return m, true, nil
}

func findPreImageHash(m manifest, hash string) bool {
	for _, e := range m.PreImages {
		if e.SHA256 == hash {
			return true
		}
	}
	return false
}

func (s *Store) readManifestBytes() ([]byte, bool) {
	raw, err := os.ReadFile(s.manifestPath())
	if err != nil {
		return nil, false
	}
	return raw, true
}

func (s *Store) decodeSchema(raw []byte) int {
	var m manifest
	if err := json.Unmarshal(raw, &m); err != nil {
		return 0
	}
	return m.Schema
}

// writeManifest commits the manifest atomically: it is small, and a half-written
// manifest would be read as corruption.
func (s *Store) writeManifest(m manifest) error {
	raw, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}
	dir := filepath.Dir(s.path)
	f, err := os.CreateTemp(dir, "."+filepath.Base(s.manifestPath())+".tmp-")
	if err != nil {
		return err
	}
	name := f.Name()
	if _, err := f.Write(raw); err != nil {
		_ = f.Close()
		_ = os.Remove(name)
		return err
	}
	if err := f.Sync(); err != nil {
		_ = f.Close()
		_ = os.Remove(name)
		return err
	}
	if err := f.Close(); err != nil {
		_ = os.Remove(name)
		return err
	}
	if err := s.replaceWithRetry(name, s.manifestPath()); err != nil {
		_ = os.Remove(name)
		return err
	}
	return nil
}

func (s *Store) restoreManifest(prev []byte, hadPrev bool) {
	if !hadPrev {
		_ = os.Remove(s.manifestPath())
		return
	}
	dir := filepath.Dir(s.path)
	f, err := os.CreateTemp(dir, "."+filepath.Base(s.manifestPath())+".tmp-")
	if err != nil {
		return
	}
	name := f.Name()
	if _, err := f.Write(prev); err != nil {
		_ = f.Close()
		_ = os.Remove(name)
		return
	}
	if err := f.Close(); err != nil {
		_ = os.Remove(name)
		return
	}
	if err := os.Rename(name, s.manifestPath()); err != nil {
		_ = os.Remove(name)
	}
}

// recoverFromPreImages returns the newest pre-image that matches its recorded
// hash. A pre-image that does not match is skipped, not trusted.
func (s *Store) recoverFromPreImages(m manifest) ([]byte, string, bool) {
	dir := filepath.Dir(s.path)
	for _, entry := range m.PreImages {
		data, err := os.ReadFile(filepath.Join(dir, entry.Name))
		if err != nil {
			continue
		}
		if hashOf(data) == entry.SHA256 {
			return data, entry.Name, true
		}
	}
	return nil, "", false
}
