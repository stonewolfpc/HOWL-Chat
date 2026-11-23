package durable

import (
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

func newStore(t *testing.T, opts Options) *Store {
	t.Helper()
	dir := t.TempDir()
	return Open(filepath.Join(dir, "settings.json"), opts)
}

func mustSave(t *testing.T, s *Store, doc string) {
	t.Helper()
	res := s.Save([]byte(doc), 1)
	if res.Outcome != Saved {
		t.Fatalf("Save(%q) = %+v, want Saved", doc, res)
	}
}

func TestSaveThenLoad(t *testing.T) {
	s := newStore(t, DefaultOptions())
	mustSave(t, s, `{"a":1}`)

	got, res := s.Load()
	if res.Outcome != Saved {
		t.Fatalf("Load outcome = %+v, want Saved", res)
	}
	if string(got) != `{"a":1}` {
		t.Fatalf("Load = %q", got)
	}
	if res.Schema != 1 {
		t.Fatalf("schema = %d", res.Schema)
	}
	if _, err := os.Stat(s.manifestPath()); err != nil {
		t.Fatalf("manifest missing after save: %v", err)
	}
}

func TestLoadAbsentWhenNoDocument(t *testing.T) {
	s := newStore(t, DefaultOptions())
	got, res := s.Load()
	if res.Outcome != Absent {
		t.Fatalf("outcome = %+v, want Absent", res)
	}
	if got != nil {
		t.Fatalf("absent document returned bytes %q", got)
	}
}

func TestSaveRotatesRetainedPreImages(t *testing.T) {
	s := newStore(t, DefaultOptions())
	for _, doc := range []string{"d1", "d2", "d3", "d4"} {
		mustSave(t, s, doc)
	}

	if got := readFileOrEmpty(t, s.preImagePath(1)); got != "d3" {
		t.Fatalf("v1 = %q, want d3", got)
	}
	if got := readFileOrEmpty(t, s.preImagePath(2)); got != "d2" {
		t.Fatalf("v2 = %q, want d2", got)
	}
	if _, err := os.Stat(s.preImagePath(3)); !os.IsNotExist(err) {
		t.Fatalf("v3 should not exist beyond the retention window: %v", err)
	}
	if got, _ := s.Load(); string(got) != "d4" {
		t.Fatalf("live = %q, want d4", got)
	}
}

// TestReaderNeverObservesPartialDocument is the core atomicity claim: concurrent
// reads see a complete previous or complete new document, never a mixture.
func TestReaderNeverObservesPartialDocument(t *testing.T) {
	s := newStore(t, DefaultOptions())
	oldDoc := `{"gen":"old","pad":"` + pad() + `"}`
	newDoc := `{"gen":"new","pad":"` + pad() + `"}`
	mustSave(t, s, oldDoc)

	allowed := map[string]bool{oldDoc: true, newDoc: true}
	var mu sync.Mutex
	var bad []string
	stop := make(chan struct{})
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-stop:
				return
			default:
			}
			got, res := s.Load()
			if res.OK() {
				mu.Lock()
				if !allowed[string(got)] {
					bad = append(bad, string(got))
				}
				mu.Unlock()
			}
			time.Sleep(200 * time.Microsecond)
		}
	}()

	for i := 0; i < 20; i++ {
		res := s.Save([]byte(newDoc), 1)
		switch {
		case res.Outcome == Saved:
		case res.Outcome == NotSaved && res.Reason == "rename_denied_by_open_reader":
			// A reader held the document past the retry window. Legitimate:
			// the writer reports it instead of forcing the write.
		default:
			t.Errorf("unexpected save outcome during contention: %+v", res)
		}
	}
	close(stop)
	wg.Wait()

	mu.Lock()
	defer mu.Unlock()
	if len(bad) > 0 {
		t.Fatalf("reader observed %d documents that were neither old nor new; first = %q", len(bad), bad[0])
	}

	// Settle to a known state and confirm the last committed write is what reads.
	deadline := time.Now().Add(2 * time.Second)
	for {
		if res := s.Save([]byte(newDoc), 1); res.Outcome == Saved {
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("could not commit a final write")
		}
	}
	if got, _ := s.Load(); string(got) != newDoc {
		t.Fatalf("final live document = %q, want the committed write", got)
	}
}

func TestSaveReportsBusyWhenLockHeld(t *testing.T) {
	s := newStore(t, Options{
		LockTimeout:     150 * time.Millisecond,
		RenameAttempts:  1,
		RenameWait:      10 * time.Millisecond,
		PreImages:       2,
		SupportedSchema: 1,
	})
	mustSave(t, s, "first")

	held, err := acquireLock(s.LockPath(), time.Second)
	if err != nil {
		t.Fatalf("acquire lock: %v", err)
	}
	defer held.release()

	res := s.Save([]byte("second"), 1)
	if res.Outcome != Busy {
		t.Fatalf("Save under held lock = %+v, want Busy", res)
	}
	if got := readFileOrEmpty(t, s.Path()); got != "first" {
		t.Fatalf("document changed despite Busy: %q", got)
	}
}

// TestSaveReportsRenameDeniedWhenReaderHoldsLive records the measured Windows
// behaviour this design depends on: a held reader blocks replacement.
func TestSaveReportsRenameDeniedWhenReaderHoldsLive(t *testing.T) {
	s := newStore(t, Options{
		LockTimeout:     time.Second,
		RenameAttempts:  2,
		RenameWait:      10 * time.Millisecond,
		PreImages:       2,
		SupportedSchema: 1,
	})
	mustSave(t, s, "previous")

	reader, err := os.Open(s.Path())
	if err != nil {
		t.Fatalf("open reader: %v", err)
	}
	defer reader.Close()

	res := s.Save([]byte("next"), 1)
	if res.Outcome != NotSaved || res.Reason != "rename_denied_by_open_reader" {
		t.Fatalf("Save with held reader = %+v, want NotSaved(rename_denied_by_open_reader)", res)
	}
	if got := readFileOrEmpty(t, s.Path()); got != "previous" {
		t.Fatalf("document = %q, want previous intact", got)
	}
	// The failed commit must not leave the manifest pointing at a document that
	// was never installed.
	if got, lres := s.Load(); string(got) != "previous" || lres.Outcome != Saved {
		t.Fatalf("Load after failed commit = %q %+v, want previous/Saved", got, lres)
	}
}

func TestLoadRecoversFromPreImageWhenLiveCorrupt(t *testing.T) {
	s := newStore(t, DefaultOptions())
	mustSave(t, s, "d1")
	mustSave(t, s, "d2")

	if err := os.WriteFile(s.Path(), []byte("torn"), 0o644); err != nil {
		t.Fatalf("corrupt live: %v", err)
	}

	got, res := s.Load()
	if res.Outcome != Recovered {
		t.Fatalf("outcome = %+v, want Recovered", res)
	}
	if string(got) != "d1" {
		t.Fatalf("recovered = %q, want d1", got)
	}
	if res.From == "" {
		t.Fatal("Recovered without naming the pre-image it came from")
	}
}

func TestLoadRefusesNewerSchema(t *testing.T) {
	s := newStore(t, DefaultOptions())
	mustSave(t, s, "d1")
	writeManifestFor(t, s, 9)

	before := readFileOrEmpty(t, s.Path())
	_, res := s.Load()
	if res.Outcome != RefusedUnknownSchema {
		t.Fatalf("outcome = %+v, want RefusedUnknownSchema", res)
	}
	if res.Schema != 9 || res.Supported != 1 {
		t.Fatalf("schema reporting = %d/%d", res.Schema, res.Supported)
	}
	if after := readFileOrEmpty(t, s.Path()); after != before {
		t.Fatalf("refused document was rewritten: %q -> %q", before, after)
	}
}

func TestSaveRefusesNewerSchema(t *testing.T) {
	s := newStore(t, DefaultOptions())
	res := s.Save([]byte("x"), 9)
	if res.Outcome != RefusedUnknownSchema {
		t.Fatalf("outcome = %+v, want RefusedUnknownSchema", res)
	}
	if _, err := os.Stat(s.Path()); !os.IsNotExist(err) {
		t.Fatalf("document written despite schema refusal: %v", err)
	}
}

// TestCrashBeforeReplaceLeavesPreviousDocument kills a real process between the
// manifest commit and the replacement, then verifies the previous document is
// intact and the leftover temp is swept on the next write.
func TestCrashBeforeReplaceLeavesPreviousDocument(t *testing.T) {
	if os.Getenv("DURABLE_CRASH_CHILD") == "1" {
		t.Skip("helper process; driven by the parent test")
	}
	s := newStore(t, DefaultOptions())
	mustSave(t, s, "previous")

	cmd := exec.Command(os.Args[0], "-test.run=^TestCrashChildDurable$", "-test.v")
	cmd.Env = append(os.Environ(),
		"DURABLE_CRASH_CHILD=1",
		"DURABLE_CRASH_PATH="+s.Path(),
	)
	out, err := cmd.CombinedOutput()
	exitErr, ok := err.(*exec.ExitError)
	if !ok || exitErr.ExitCode() != 9 {
		t.Fatalf("child did not die at the injection point: err=%v out=%s", err, out)
	}

	if got := readFileOrEmpty(t, s.Path()); got != "previous" {
		t.Fatalf("document after crash = %q, want previous intact", got)
	}
	// The crash lands between the manifest commit and the replacement, so the
	// manifest describes a document that was never installed. The previous
	// document must still be readable, and the outcome must say why it is being
	// reported this way rather than claiming a clean save or corruption.
	got, res := s.Load()
	if string(got) != "previous" {
		t.Fatalf("Load after crash = %q, want previous intact", got)
	}
	if res.Outcome != Recovered || res.Reason != "interrupted_commit_rolled_back" {
		t.Fatalf("Load after crash = %+v, want Recovered(interrupted_commit_rolled_back)", res)
	}

	matches, _ := filepath.Glob(filepath.Join(filepath.Dir(s.Path()), s.tempPrefix()+"*"))
	if len(matches) == 0 {
		t.Fatal("expected the crashed writer to leave a temp file to sweep")
	}

	mustSave(t, s, "after-crash")
	matches, _ = filepath.Glob(filepath.Join(filepath.Dir(s.Path()), s.tempPrefix()+"*"))
	if len(matches) != 0 {
		t.Fatalf("leftover temps not swept: %v", matches)
	}
	if got, _ := s.Load(); string(got) != "after-crash" {
		t.Fatalf("live = %q, want after-crash", got)
	}
}

// TestCrashChildDurable is not a test: it is the process that dies mid-commit.
func TestCrashChildDurable(t *testing.T) {
	if os.Getenv("DURABLE_CRASH_CHILD") != "1" {
		t.Skip("only runs as a crash-injection child")
	}
	s := Open(os.Getenv("DURABLE_CRASH_PATH"), DefaultOptions())
	testHookBeforeReplace = func() { os.Exit(9) }
	s.Save([]byte("never-committed"), 1)
	os.Exit(0)
}

// TestBusyAcrossProcesses proves the lock is real mutual exclusion between
// processes, not merely between handles in one process. A second process holds it
// while this one attempts a write.
func TestBusyAcrossProcesses(t *testing.T) {
	if os.Getenv("DURABLE_LOCK_CHILD") == "1" {
		t.Skip("helper process; driven by the parent test")
	}
	s := newStore(t, Options{
		LockTimeout:     250 * time.Millisecond,
		RenameAttempts:  1,
		RenameWait:      10 * time.Millisecond,
		PreImages:       2,
		SupportedSchema: 1,
	})
	marker := s.Path() + ".held"

	cmd := exec.Command(os.Args[0], "-test.run=^TestLockHolderChild$", "-test.v")
	cmd.Env = append(os.Environ(),
		"DURABLE_LOCK_CHILD=1",
		"DURABLE_LOCK_PATH="+s.LockPath(),
		"DURABLE_LOCK_MARKER="+marker,
	)
	var out []byte
	var runErr error
	done := make(chan struct{})
	go func() {
		defer close(done)
		out, runErr = cmd.CombinedOutput()
	}()

	deadline := time.Now().Add(5 * time.Second)
	for {
		if _, err := os.Stat(marker); err == nil {
			break
		}
		if time.Now().After(deadline) {
			<-done // let the child finish before reading its output
			t.Fatalf("child never acquired the lock; err=%v output=%s", runErr, out)
		}
		time.Sleep(10 * time.Millisecond)
	}

	res := s.Save([]byte("blocked"), 1)
	if res.Outcome != Busy {
		t.Fatalf("Save while another process holds the lock = %+v, want Busy", res)
	}

	<-done
	if runErr != nil {
		t.Fatalf("lock holder child failed: %v\n%s", runErr, out)
	}

	// With the other process gone, the same write must now succeed.
	if res := s.Save([]byte("allowed"), 1); res.Outcome != Saved {
		t.Fatalf("Save after lock release = %+v, want Saved", res)
	}
}

// TestLockHolderChild is not a test: it holds the lock so the parent can contend.
func TestLockHolderChild(t *testing.T) {
	if os.Getenv("DURABLE_LOCK_CHILD") != "1" {
		t.Skip("only runs as a lock-holder child")
	}
	lock, err := acquireLock(os.Getenv("DURABLE_LOCK_PATH"), time.Second)
	if err != nil {
		os.Exit(3)
	}
	if err := os.WriteFile(os.Getenv("DURABLE_LOCK_MARKER"), []byte("held"), 0o644); err != nil {
		os.Exit(4)
	}
	time.Sleep(1200 * time.Millisecond)
	lock.release()
	os.Exit(0)
}

func readFileOrEmpty(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return ""
		}
		t.Fatalf("read %s: %v", path, err)
	}
	return string(data)
}

func writeManifestFor(t *testing.T, s *Store, schema int) {
	t.Helper()
	raw := []byte(`{"schemaVersion":` + itoa(schema) + `,"sha256":"0000","size":0,"committedAt":"x"}`)
	if err := os.WriteFile(s.manifestPath(), raw, 0o644); err != nil {
		t.Fatalf("write manifest: %v", err)
	}
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var buf [8]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	return string(buf[i:])
}

func pad() string {
	b := make([]byte, 32*1024)
	for i := range b {
		b[i] = 'x'
	}
	return string(b)
}
