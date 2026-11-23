package diagnostics

import (
	"encoding/json"
	stderrors "errors"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"sort"
	"strings"
	"testing"
	"time"

	apperrors "howl-chat/internal/errors"
)

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func enabledReporter(t *testing.T) (*Owner, *Reporter) {
	t.Helper()
	owner, err := Open(Config{LogDir: t.TempDir(), AppVersion: "test-version"})
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	t.Cleanup(func() { _ = owner.Close() })
	owner.SetEnabled(true)
	return owner, owner.reporter
}

// observe installs the asynchronous result hook. It must be called before the
// upload it observes.
func observe(t *testing.T, r *Reporter) <-chan error {
	t.Helper()
	ch := make(chan error, 1)
	r.onResult = func(err error) { ch <- err }
	return ch
}

func awaitResult(t *testing.T, ch <-chan error) error {
	t.Helper()
	select {
	case err := <-ch:
		return err
	case <-time.After(5 * time.Second):
		t.Fatal("upload outcome was never delivered")
		return nil
	}
}

func keySet(m map[string]json.RawMessage) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func decodeObject(t *testing.T, label string, raw []byte) map[string]json.RawMessage {
	t.Helper()
	var m map[string]json.RawMessage
	if err := json.Unmarshal(raw, &m); err != nil {
		t.Fatalf("%s is not a JSON object: %v (%s)", label, err, raw)
	}
	return m
}

func requireKeys(t *testing.T, label string, m map[string]json.RawMessage, want ...string) {
	t.Helper()
	sort.Strings(want)
	got := keySet(m)
	if strings.Join(got, ",") != strings.Join(want, ",") {
		t.Fatalf("%s keys = %v, want exactly %v", label, got, want)
	}
}

// captured is one request as the endpoint saw it.
type captured struct {
	method string
	path   string
	ctype  string
	marker string
	body   []byte
}

// captureServer receives upload requests and can answer with any status. The
// request is handed to the test through a channel, so the assertion has a
// proper happens-before edge on the captured bytes.
type captureServer struct {
	url      string
	requests chan captured
}

func newCaptureServer(t *testing.T, respond func(w http.ResponseWriter)) *captureServer {
	t.Helper()
	c := &captureServer{requests: make(chan captured, 1)}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		c.requests <- captured{
			method: r.Method,
			path:   r.URL.Path,
			ctype:  r.Header.Get("Content-Type"),
			marker: r.Header.Get("X-Error-Report"),
			body:   body,
		}
		if respond != nil {
			respond(w)
		}
	}))
	t.Cleanup(srv.Close)
	c.url = srv.URL + "/api/report"
	return c
}

func (c *captureServer) await(t *testing.T) captured {
	t.Helper()
	select {
	case got := <-c.requests:
		return got
	case <-time.After(5 * time.Second):
		t.Fatal("the upload request never reached the endpoint")
		return captured{}
	}
}

// ---------------------------------------------------------------------------
// Transmission
// ---------------------------------------------------------------------------

func TestUploadPostsTheDocumentedRequestAndAccepts2xx(t *testing.T) {
	owner, reporter := enabledReporter(t)
	owner.Report(apperrors.New(apperrors.IMGModelLoadFail, apperrors.High, "image.load", "reporter_test.go:1", stderrors.New("img")))
	report, err := owner.BuildReport(50)
	if err != nil {
		t.Fatalf("BuildReport: %v", err)
	}

	capture := newCaptureServer(t, func(w http.ResponseWriter) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"ok":true}`))
	})
	reporter.uploadURL = capture.url
	result := observe(t, reporter)

	// The binding-side contract: Upload returns before the network exchange.
	if err := reporter.Upload(report); err != nil {
		t.Fatalf("Upload returned a synchronous error: %v", err)
	}
	if err := awaitResult(t, result); err != nil {
		t.Fatalf("2xx outcome = %v, want nil", err)
	}

	got := capture.await(t)
	if got.method != http.MethodPost {
		t.Fatalf("method = %q, want POST", got.method)
	}
	if got.path != "/api/report" {
		t.Fatalf("path = %q, want /api/report", got.path)
	}
	if got.ctype != "application/json" {
		t.Fatalf("Content-Type = %q", got.ctype)
	}
	if got.marker != "howl-chat" {
		t.Fatalf("X-Error-Report = %q", got.marker)
	}

	var decoded Report
	if err := json.Unmarshal(got.body, &decoded); err != nil {
		t.Fatalf("request body is not a report: %v (%s)", err, got.body)
	}
	if decoded.ReportID != report.ReportID || decoded.ErrorCount != report.ErrorCount {
		t.Fatalf("body report = %+v, sent %+v", decoded, report)
	}
}

func TestUploadSurfacesANon2xxReply(t *testing.T) {
	owner, reporter := enabledReporter(t)
	owner.Report(apperrors.New(apperrors.CHATSaveFail, apperrors.High, "chat.save", "reporter_test.go:2", stderrors.New("save")))
	report, err := owner.BuildReport(50)
	if err != nil {
		t.Fatalf("BuildReport: %v", err)
	}

	capture := newCaptureServer(t, func(w http.ResponseWriter) {
		w.WriteHeader(http.StatusServiceUnavailable)
		_, _ = w.Write([]byte(`{"error":"maintenance"}`))
	})
	reporter.uploadURL = capture.url
	result := observe(t, reporter)

	if err := reporter.Upload(report); err != nil {
		t.Fatalf("Upload returned a synchronous error for an async failure: %v", err)
	}
	outcome := awaitResult(t, result)
	if outcome == nil {
		t.Fatal("503 reply was reported as a successful upload")
	}
	if !strings.Contains(outcome.Error(), "503") || !strings.Contains(outcome.Error(), "maintenance") {
		t.Fatalf("outcome = %v, want the status and reply body", outcome)
	}

	// The failure must not vanish: it is written to the diagnostics log.
	entries := owner.Entries(20)
	found := false
	for _, e := range entries {
		if e.Code == string(apperrors.NETRequestFail) && e.Operation == "error_report.upload" {
			found = true
			if !strings.Contains(e.Message, "503") {
				t.Fatalf("logged upload failure message = %q, want the status", e.Message)
			}
			if e.Location != "internal/diagnostics/reporter.go:Upload" {
				t.Fatalf("logged upload failure location = %q", e.Location)
			}
		}
	}
	if !found {
		t.Fatalf("no NET-REQ entry for the failed upload: %+v", entries)
	}

	// Reporting still works after an upload failure.
	owner.Report(apperrors.New(apperrors.CHATLoadFail, apperrors.High, "chat.load", "reporter_test.go:3", stderrors.New("load")))
	if got := len(owner.Entries(50)); got != 3 {
		t.Fatalf("entries after failure = %d, want 3 (two reports plus the logged refusal)", got)
	}
}

func TestUploadSurfacesANetworkErrorAndKeepsTheAppAlive(t *testing.T) {
	owner, reporter := enabledReporter(t)
	owner.Report(apperrors.New(apperrors.RESAllocFail, apperrors.High, "resource.allocate", "reporter_test.go:4", stderrors.New("alloc")))
	report, err := owner.BuildReport(50)
	if err != nil {
		t.Fatalf("BuildReport: %v", err)
	}

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	deadURL := "http://" + listener.Addr().String() + "/api/report"
	_ = listener.Close()

	reporter.uploadURL = deadURL
	result := observe(t, reporter)
	if err := reporter.Upload(report); err != nil {
		t.Fatalf("Upload returned a synchronous error for an offline endpoint: %v", err)
	}
	outcome := awaitResult(t, result)
	if outcome == nil {
		t.Fatal("refused connection was reported as success")
	}

	// The app keeps working; the refusal was logged, not swallowed.
	if entries := owner.Entries(20); len(entries) != 2 {
		t.Fatalf("entries = %d, want the report plus one NET-REQ refusal", len(entries))
	}
	owner.Report(apperrors.New(apperrors.CHATLoadFail, apperrors.High, "chat.load", "reporter_test.go:5", stderrors.New("after")))
	if got := len(owner.Entries(50)); got != 3 {
		t.Fatalf("entries after network failure = %d, want 3", got)
	}
}

func TestUploadTimesOutOnAHangingEndpoint(t *testing.T) {
	owner, reporter := enabledReporter(t)
	owner.Report(apperrors.New(apperrors.LLMNetTimeout, apperrors.High, "llm.request", "reporter_test.go:6", stderrors.New("timed out")))
	report, err := owner.BuildReport(50)
	if err != nil {
		t.Fatalf("BuildReport: %v", err)
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Never answer: hold the connection until the client gives up.
		select {
		case <-r.Context().Done():
		case <-time.After(3 * time.Second):
		}
	}))
	t.Cleanup(srv.Close)
	reporter.uploadURL = srv.URL + "/api/report"
	reporter.client = &http.Client{Timeout: 120 * time.Millisecond}

	result := observe(t, reporter)
	start := time.Now()
	if err := reporter.Upload(report); err != nil {
		t.Fatalf("Upload returned a synchronous error: %v", err)
	}
	outcome := awaitResult(t, result)
	elapsed := time.Since(start)

	if outcome == nil {
		t.Fatal("hanging endpoint was reported as success")
	}
	var netErr net.Error
	if !stderrors.As(outcome, &netErr) || !netErr.Timeout() {
		t.Fatalf("outcome = %v, want a timeout", outcome)
	}
	if elapsed > 2*time.Second {
		t.Fatalf("upload took %v, want it bounded by the client timeout", elapsed)
	}

	entries := owner.Entries(20)
	found := false
	for _, e := range entries {
		if e.Code == string(apperrors.NETRequestFail) && e.Operation == "error_report.upload" {
			found = true
		}
	}
	if !found {
		t.Fatalf("timeout was not logged: %+v", entries)
	}
}

func TestUploadDoesNotBlockWhileTheEndpointHangs(t *testing.T) {
	owner, reporter := enabledReporter(t)
	owner.Report(apperrors.New(apperrors.IMGGallerySaveFail, apperrors.High, "image.gallery", "reporter_test.go:7", stderrors.New("gallery")))
	report, err := owner.BuildReport(50)
	if err != nil {
		t.Fatalf("BuildReport: %v", err)
	}

	started := make(chan struct{})
	release := make(chan struct{})
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		close(started)
		<-release
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(srv.Close)
	reporter.uploadURL = srv.URL + "/api/report"

	result := observe(t, reporter)
	begin := time.Now()
	if err := reporter.Upload(report); err != nil {
		t.Fatalf("Upload returned a synchronous error: %v", err)
	}
	elapsed := time.Since(begin)

	select {
	case <-started:
	case <-time.After(2 * time.Second):
		t.Fatal("the request never reached the endpoint")
	}
	if elapsed > 250*time.Millisecond {
		t.Fatalf("Upload blocked for %v while the endpoint hung", elapsed)
	}

	close(release)
	if err := awaitResult(t, result); err != nil {
		t.Fatalf("outcome after release = %v, want nil", err)
	}
}

// ---------------------------------------------------------------------------
// Wire payload / privacy boundary
// ---------------------------------------------------------------------------

func TestSerializedReportCarriesOnlyTheDocumentedFieldsAndNoContextValues(t *testing.T) {
	owner, reporter := enabledReporter(t)

	const (
		secretEmail  = "alice@example.com"
		secretPath   = `C:\Users\alice\Documents\diary.txt`
		secretPrompt = "my private prompt: open the vault"
		secretChat   = "assistant: the secret reply text"
		secretMsg    = "message that must never be uploaded"
	)
	owner.Report(apperrors.New(apperrors.CHATSaveFail, apperrors.High, "chat.save", "reporter_test.go:8", stderrors.New(secretMsg)).
		WithContext("note", secretEmail).
		WithContext("detail", secretPath).
		WithContext("prompt", secretPrompt).
		WithContext("chat", secretChat).
		WithContext("email", secretEmail)) // exact sensitive key: sanitized in the local log too
	owner.Report(apperrors.New(apperrors.CHATLoadFail, apperrors.Medium, "chat.load", "reporter_test.go:9", stderrors.New("plain")))

	report, err := owner.BuildReport(50)
	if err != nil {
		t.Fatalf("BuildReport: %v", err)
	}
	if len(report.Errors) != 2 {
		t.Fatalf("report errors = %d, want 2", len(report.Errors))
	}

	capture := newCaptureServer(t, func(w http.ResponseWriter) { w.WriteHeader(http.StatusOK) })
	reporter.uploadURL = capture.url
	result := observe(t, reporter)
	if err := reporter.Upload(report); err != nil {
		t.Fatalf("Upload: %v", err)
	}
	if err := awaitResult(t, result); err != nil {
		t.Fatalf("outcome = %v", err)
	}

	got := capture.await(t)
	body := string(got.body)
	for _, secret := range []string{secretEmail, secretPath, "diary.txt", secretPrompt, "open the vault", secretChat, "secret reply", secretMsg} {
		if strings.Contains(body, secret) {
			t.Fatalf("serialized report leaked %q: %s", secret, body)
		}
	}
	if strings.Contains(body, `"msg"`) {
		t.Fatalf("serialized report contains a message field: %s", body)
	}
	if strings.Contains(body, `"context":`) {
		t.Fatalf("serialized report contains context values: %s", body)
	}

	top := decodeObject(t, "report", got.body)
	requireKeys(t, "report", top, "error_count", "errors", "generated", "report_id", "system")

	system := decodeObject(t, "system", top["system"])
	requireKeys(t, "system", system, "app_version", "arch", "cpu_count", "go_version", "os")

	var entries []json.RawMessage
	if err := json.Unmarshal(top["errors"], &entries); err != nil {
		t.Fatalf("errors is not an array: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("serialized errors = %d, want 2", len(entries))
	}
	withContext := decodeObject(t, "error[0]", entries[0])
	requireKeys(t, "error[0]", withContext, "code", "context_keys", "location", "operation", "severity", "time")
	plain := decodeObject(t, "error[1]", entries[1])
	requireKeys(t, "error[1]", plain, "code", "location", "operation", "severity", "time")

	var keys []string
	if err := json.Unmarshal(withContext["context_keys"], &keys); err != nil {
		t.Fatalf("context_keys is not an array: %v", err)
	}
	want := []string{"chat", "detail", "email", "note", "prompt"}
	sort.Strings(want)
	if strings.Join(keys, ",") != strings.Join(want, ",") {
		t.Fatalf("context_keys = %v, want the key names only %v", keys, want)
	}
}

func TestReportIsCappedAtFiftyEntries(t *testing.T) {
	owner, _ := enabledReporter(t)
	for i := 0; i < 60; i++ {
		owner.Report(apperrors.New(apperrors.CHATSaveFail, apperrors.High, "chat.save", "reporter_test.go:cap", stderrors.New("x")))
	}

	for _, ask := range []int{50, 1000} {
		report, err := owner.BuildReport(ask)
		if err != nil {
			t.Fatalf("BuildReport(%d): %v", ask, err)
		}
		if len(report.Errors) != maxReportEntries || report.ErrorCount != maxReportEntries {
			t.Fatalf("BuildReport(%d) = %d entries (count %d), want %d", ask, len(report.Errors), report.ErrorCount, maxReportEntries)
		}
	}
}
