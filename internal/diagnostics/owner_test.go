package diagnostics

import (
	"encoding/json"
	stderrors "errors"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"howl-chat/internal/errors"
)

func today() string { return time.Now().UTC().Format("20060102") }

type portEvent struct {
	name string
	data []interface{}
}

type recordingPort struct {
	mu     sync.Mutex
	events []portEvent
}

func (p *recordingPort) Emit(name string, data ...interface{}) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.events = append(p.events, portEvent{name: name, data: data})
}

func (p *recordingPort) snapshot() []portEvent {
	p.mu.Lock()
	defer p.mu.Unlock()
	return append([]portEvent(nil), p.events...)
}

type panickingPort struct{}

func (panickingPort) Emit(string, ...interface{}) { panic("port exploded") }

func openTestOwner(t *testing.T) (*Owner, string) {
	t.Helper()
	dir := filepath.Join(t.TempDir(), "logs")
	owner, err := Open(Config{LogDir: dir, AppVersion: "test-version"})
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	t.Cleanup(func() { _ = owner.Close() })
	return owner, dir
}

func TestOpenWritesTheSameLogLineShapeAndReadsItBack(t *testing.T) {
	owner, dir := openTestOwner(t)
	owner.Report(errors.New(errors.LLMServerSpawnFail, errors.Critical, "llm.load", "test.go:1", stderrors.New("boom")))

	entries := owner.Entries(10)
	if len(entries) != 1 {
		t.Fatalf("entries = %d, want 1", len(entries))
	}
	if entries[0].Code != string(errors.LLMServerSpawnFail) || entries[0].Severity != "critical" {
		t.Fatalf("entry = %+v", entries[0])
	}
	if entries[0].Message != "boom" {
		t.Fatalf("message = %q, want boom", entries[0].Message)
	}

	files, err := filepath.Glob(filepath.Join(dir, "howl-errors-*.jsonl"))
	if err != nil || len(files) != 1 {
		t.Fatalf("log files = %v, err %v", files, err)
	}
	data, err := os.ReadFile(files[0])
	if err != nil {
		t.Fatal(err)
	}
	var raw map[string]interface{}
	if err := json.Unmarshal([]byte(strings.TrimSpace(string(data))), &raw); err != nil {
		t.Fatalf("log line is not JSON: %v", err)
	}
	for _, key := range []string{"t", "c", "s", "op", "loc", "msg"} {
		if _, ok := raw[key]; !ok {
			t.Fatalf("log line missing key %q: %v", key, raw)
		}
	}
}

func TestOpenReportsAnUnwritableLogDirectory(t *testing.T) {
	file := filepath.Join(t.TempDir(), "not-a-dir")
	if err := os.WriteFile(file, []byte("x"), 0644); err != nil {
		t.Fatal(err)
	}
	if owner, err := Open(Config{LogDir: file}); err == nil || owner != nil {
		t.Fatalf("Open on a file path = (%v, %v), want an error and no owner", owner, err)
	}
}

func TestNilOwnerIsInert(t *testing.T) {
	var owner *Owner
	owner.Attach(&recordingPort{})
	owner.SetEnabled(true)
	if owner.Enabled() {
		t.Fatal("nil owner reported uploads enabled")
	}
	owner.Report(errors.New(errors.SYSPanicRecovered, errors.Critical, "op", "loc", stderrors.New("x")))
	owner.ReportSilent(errors.New(errors.SYSShutdownFail, errors.High, "op", "loc", stderrors.New("x")))
	owner.ReportLLMLoad(stderrors.New("x"), "loc")
	if entries := owner.Entries(10); entries != nil {
		t.Fatalf("nil owner entries = %v", entries)
	}
	if _, err := owner.BuildReport(5); err == nil {
		t.Fatal("nil owner built a report")
	}
	if _, err := owner.SaveReport(&Report{}); err == nil {
		t.Fatal("nil owner saved a report")
	}
	if err := owner.Upload(&Report{}); err == nil {
		t.Fatal("nil owner uploaded a report")
	}
	if err := owner.Close(); err != nil {
		t.Fatalf("nil owner Close: %v", err)
	}
}

func TestReportEmitsOnlyErrorsThatShouldSurface(t *testing.T) {
	owner, _ := openTestOwner(t)
	port := &recordingPort{}
	owner.Attach(port)

	critical := errors.New(errors.LLMServerCrash, errors.Critical, "llm.load", "test.go:2", stderrors.New("crash")).
		WithContext("token", "secret")
	owner.Report(critical)
	owner.Report(errors.New(errors.TTSSynthesizeFail, errors.Medium, "tts.synthesize", "test.go:3", stderrors.New("quiet")))
	owner.ReportSilent(errors.New(errors.CHATSaveFail, errors.High, "chat.save", "test.go:4", stderrors.New("silent")))

	events := port.snapshot()
	if len(events) != 1 {
		t.Fatalf("events = %d, want 1 (only the critical error surfaces)", len(events))
	}
	if events[0].name != "app:error" {
		t.Fatalf("event name = %q", events[0].name)
	}
	payload, ok := events[0].data[0].(string)
	if !ok {
		t.Fatalf("event payload type %T, want string", events[0].data[0])
	}
	var decoded map[string]interface{}
	if err := json.Unmarshal([]byte(payload), &decoded); err != nil {
		t.Fatalf("payload is not JSON: %v", err)
	}
	if decoded["code"] != string(errors.LLMServerCrash) || decoded["severity"] != "critical" {
		t.Fatalf("payload = %v", decoded)
	}
	keys, _ := decoded["context_keys"].([]interface{})
	if len(keys) != 1 || keys[0] != "token" {
		t.Fatalf("context_keys = %v", decoded["context_keys"])
	}
	if _, leaked := decoded["token"]; leaked {
		t.Fatal("payload leaked a context value")
	}

	if entries := owner.Entries(10); len(entries) != 3 {
		t.Fatalf("entries = %d, want 3 (silent and quiet still log)", len(entries))
	}
}

func TestPortFailureDoesNotBreakLogging(t *testing.T) {
	owner, _ := openTestOwner(t)
	owner.Attach(panickingPort{})
	owner.Report(errors.New(errors.SYSPanicRecovered, errors.Critical, "op", "loc", stderrors.New("x")))
	if entries := owner.Entries(10); len(entries) != 1 {
		t.Fatalf("entries = %d, want 1 after a failing port", len(entries))
	}
}

func TestRepeatedReportsAreEachRecordedAndEmitted(t *testing.T) {
	owner, _ := openTestOwner(t)
	port := &recordingPort{}
	owner.Attach(port)
	tracked := errors.New(errors.CHATSaveFail, errors.High, "chat.save", "loc", stderrors.New("again"))
	owner.Report(tracked)
	owner.Report(tracked)
	if entries := owner.Entries(10); len(entries) != 2 {
		t.Fatalf("entries = %d, want 2 (repeats are not deduplicated)", len(entries))
	}
	if events := port.snapshot(); len(events) != 2 {
		t.Fatalf("events = %d, want 2", len(events))
	}
}

func TestCorruptLogLinesAreSkipped(t *testing.T) {
	owner, dir := openTestOwner(t)
	owner.Report(errors.New(errors.LLMServerTimeout, errors.High, "op", "loc", stderrors.New("valid")))
	path := filepath.Join(dir, "howl-errors-"+today()+".jsonl")
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString("this is not json\n{\"t\":\"x\"}\n"); err != nil {
		t.Fatal(err)
	}
	_ = f.Close()

	entries := owner.Entries(10)
	if len(entries) != 2 {
		t.Fatalf("entries = %d, want the valid line plus the partial object", len(entries))
	}
	if entries[0].Code != string(errors.LLMServerTimeout) {
		t.Fatalf("first entry = %+v", entries[0])
	}
}

func TestWritesAfterCloseAreDroppedAndCloseIsIdempotent(t *testing.T) {
	owner, _ := openTestOwner(t)
	owner.Report(errors.New(errors.CHATSaveFail, errors.High, "op", "loc", stderrors.New("before")))
	if err := owner.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}
	owner.Report(errors.New(errors.CHATSaveFail, errors.High, "op", "loc", stderrors.New("after")))
	if entries := owner.Entries(10); len(entries) != 1 {
		t.Fatalf("entries = %d, want only the pre-close entry", len(entries))
	}
	if err := owner.Close(); err != nil {
		t.Fatalf("second Close: %v", err)
	}
}

func TestUploadToggleAndReports(t *testing.T) {
	owner, dir := openTestOwner(t)
	if owner.Enabled() {
		t.Fatal("uploads enabled by default")
	}
	if err := owner.Upload(&Report{}); err == nil || !strings.Contains(err.Error(), "upload disabled") {
		t.Fatalf("Upload while disabled = %v", err)
	}
	if _, err := owner.BuildReport(5); err == nil {
		t.Fatal("empty log produced a report")
	}

	owner.Report(errors.New(errors.IMGModelLoadFail, errors.High, "image.load", "loc", stderrors.New("img")))
	report, err := owner.BuildReport(5)
	if err != nil {
		t.Fatalf("BuildReport: %v", err)
	}
	if report.ErrorCount != 1 || report.System.AppVersion != "test-version" {
		t.Fatalf("report = %+v", report)
	}
	path, err := owner.SaveReport(report)
	if err != nil {
		t.Fatalf("SaveReport: %v", err)
	}
	if filepath.Dir(path) != dir {
		t.Fatalf("report path %q not beside the log %q", path, dir)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var decoded Report
	if err := json.Unmarshal(data, &decoded); err != nil || decoded.ReportID == "" {
		t.Fatalf("saved report unreadable: %v", err)
	}

	owner.SetEnabled(true)
	if !owner.Enabled() {
		t.Fatal("SetEnabled(true) did not stick")
	}
}

func TestConcurrentReportsReadsAndToggles(t *testing.T) {
	owner, _ := openTestOwner(t)
	port := &recordingPort{}
	owner.Attach(port)

	var wg sync.WaitGroup
	for w := 0; w < 8; w++ {
		wg.Add(1)
		go func(w int) {
			defer wg.Done()
			for i := 0; i < 25; i++ {
				owner.Report(errors.New(errors.CHATSaveFail, errors.High, "chat.save", "loc", stderrors.New("x")))
				_ = owner.Entries(5)
				_, _ = owner.BuildReport(5)
				owner.SetEnabled(i%2 == 0)
				_ = owner.Enabled()
			}
		}(w)
	}
	wg.Wait()

	if entries := owner.Entries(1000); len(entries) != 200 {
		t.Fatalf("entries = %d, want 200", len(entries))
	}
	if events := port.snapshot(); len(events) != 200 {
		t.Fatalf("events = %d, want 200", len(events))
	}
}
