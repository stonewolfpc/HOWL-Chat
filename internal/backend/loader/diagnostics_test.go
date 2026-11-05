package loader

import (
	"errors"
	"io"
	"os"
	"os/exec"
	"strings"
	"testing"
)

type failingWriter struct{}

func (failingWriter) Write([]byte) (int, error) {
	return 0, errors.New("the GUI has no usable stderr handle")
}

// TestLoaderHelperProcess is the child used by the process tests below. It only
// does work when explicitly invoked with the environment marker.
func TestLoaderHelperProcess(t *testing.T) {
	if os.Getenv("HOWL_LOADER_HELPER") != "1" {
		return
	}
	if output := os.Getenv("HOWL_HELPER_OUTPUT"); output != "" {
		_, _ = io.WriteString(os.Stderr, output)
	}
	if output := os.Getenv("HOWL_HELPER_STDOUT"); output != "" {
		_, _ = io.WriteString(os.Stdout, output)
	}
	if lines := os.Getenv("HOWL_HELPER_LINES"); lines != "" {
		count := 0
		for i := 0; i < 64; i++ {
			count = count*10 + int(lines[i]-'0')
			if i+1 == len(lines) {
				break
			}
		}
		line := strings.Repeat("x", 511) + "\n"
		for i := 0; i < count; i++ {
			_, _ = io.WriteString(os.Stderr, line)
		}
		_, _ = io.WriteString(os.Stderr, "END-OF-OUTPUT\n")
	}
	os.Exit(3)
}

func runHelper(t *testing.T, capture io.Writer, console io.Writer, env ...string) error {
	t.Helper()
	cmd := exec.Command(os.Args[0], "-test.run=TestLoaderHelperProcess", "--")
	cmd.Env = append(os.Environ(), append([]string{"HOWL_LOADER_HELPER=1"}, env...)...)
	cmd.Stdout = io.MultiWriter(capture, SafeConsole(console))
	cmd.Stderr = io.MultiWriter(capture, SafeConsole(console))
	return cmd.Run()
}

func TestCaptureKeepsABoundedTail(t *testing.T) {
	capture := NewCapture(128)
	chunk := strings.Repeat("a", 100)
	for i := 0; i < 50; i++ {
		if _, err := capture.Write([]byte(chunk)); err != nil {
			t.Fatalf("Write: %v", err)
		}
	}
	if _, err := capture.Write([]byte("LAST-MARKER")); err != nil {
		t.Fatalf("Write: %v", err)
	}

	total, retained, truncated := capture.Stats()
	wantTotal := int64(50*100 + len("LAST-MARKER"))
	if total != wantTotal {
		t.Fatalf("total = %d, want %d", total, wantTotal)
	}
	if retained != 128 {
		t.Fatalf("retained = %d, want 128", retained)
	}
	if !truncated {
		t.Fatal("capture did not report truncation")
	}
	if tail := capture.Tail(0); !strings.HasSuffix(tail, "LAST-MARKER") {
		t.Fatalf("tail %q lost the newest bytes", tail)
	}
}

func TestChildExitIsCapturedWhenTheConsoleIsUnavailable(t *testing.T) {
	capture := NewCapture(DefaultCaptureBytes)
	err := runHelper(t, capture, failingWriter{}, "HOWL_HELPER_OUTPUT=unknown model architecture: vocoder\n")
	if err == nil {
		t.Fatal("helper process unexpectedly succeeded")
	}

	if status := ExitStatus(err); !strings.Contains(status, "3") {
		t.Fatalf("exit status = %q, want the child's code 3", status)
	}
	tail := capture.Tail(0)
	if !strings.Contains(tail, "unknown model architecture") {
		t.Fatalf("capture %q did not retain the child's stderr when the console failed", tail)
	}
}

func TestStdoutAndStderrAreBothCapturedWhenTheConsoleIsUnavailable(t *testing.T) {
	capture := NewCapture(DefaultCaptureBytes)
	err := runHelper(t, capture, failingWriter{},
		"HOWL_HELPER_OUTPUT=stderr-line\n",
		"HOWL_HELPER_STDOUT=stdout-line\n")
	if err == nil {
		t.Fatal("helper process unexpectedly succeeded")
	}
	tail := capture.Tail(0)
	if !strings.Contains(tail, "stderr-line") || !strings.Contains(tail, "stdout-line") {
		t.Fatalf("capture %q did not retain both streams", tail)
	}
}

func TestChildExitIsCapturedWithoutAConsoleWriter(t *testing.T) {
	capture := NewCapture(DefaultCaptureBytes)
	if err := runHelper(t, capture, nil, "HOWL_HELPER_OUTPUT=fatal: bad model\n"); err == nil {
		t.Fatal("helper process unexpectedly succeeded")
	}
	if !strings.Contains(capture.Tail(0), "fatal: bad model") {
		t.Fatalf("capture %q lost output for a nil console", capture.Tail(0))
	}
}

func TestLargeOutputStaysBounded(t *testing.T) {
	const limit = 4096
	capture := NewCapture(limit)
	if err := runHelper(t, capture, failingWriter{}, "HOWL_HELPER_LINES=8192"); err == nil {
		t.Fatal("helper process unexpectedly succeeded")
	}

	total, retained, truncated := capture.Stats()
	if retained > limit {
		t.Fatalf("retained %d bytes, want at most %d", retained, limit)
	}
	if total <= int64(8192*511) {
		t.Fatalf("total = %d, want the full large stream to be counted", total)
	}
	if !truncated {
		t.Fatal("large output was not reported as truncated")
	}
	if !strings.Contains(capture.Tail(0), "END-OF-OUTPUT") {
		t.Fatal("bounded capture lost the end of the stream")
	}
}

func TestFailureRecordCarriesUsefulDetailsAndPersists(t *testing.T) {
	capture := NewCapture(1024)
	_, _ = capture.Write([]byte("llama_model_load: error loading model: unknown architecture\n"))
	total, _, truncated := capture.Stats()

	rec := FailureRecord{
		Stage:           StageStartup,
		ModelPath:       `C:\models\vocoder.gguf`,
		ModelName:       "vocoder.gguf",
		Arch:            "vocoder",
		Projector:       `C:\models\mmproj.gguf`,
		Executable:      `C:\app\runtime\llama\llama-server.exe`,
		Args:            []string{"-m", `C:\models\vocoder.gguf`, "--mmproj", `C:\models\mmproj.gguf`},
		ExitStatus:      "exit status 1",
		Output:          capture.Tail(FailureTailBytes),
		OutputBytes:     total,
		OutputTruncated: truncated,
		Detail:          "llama-server exited before becoming healthy",
	}

	message := rec.UserMessage()
	for _, want := range []string{"vocoder.gguf", "architecture vocoder", "exit status 1", "unknown architecture"} {
		if !strings.Contains(message, want) {
			t.Fatalf("user message %q is missing %q", message, want)
		}
	}

	dir := t.TempDir()
	path, err := WriteFailureRecord(dir, rec)
	if err != nil {
		t.Fatalf("WriteFailureRecord: %v", err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read record: %v", err)
	}
	text := string(data)
	for _, want := range []string{
		"stage: startup",
		"architecture: vocoder",
		"projector:",
		"executable:",
		"exit: exit status 1",
		"output_tail:",
		"unknown architecture",
	} {
		if !strings.Contains(text, want) {
			t.Fatalf("persisted record is missing %q:\n%s", want, text)
		}
	}

	if _, err := WriteFailureRecord(dir, rec); err != nil {
		t.Fatalf("second WriteFailureRecord: %v", err)
	}
	data, err = os.ReadFile(path)
	if err != nil {
		t.Fatalf("read record after append: %v", err)
	}
	if strings.Count(string(data), "stage: startup") != 2 {
		t.Fatalf("records are not appended:\n%s", data)
	}
}

func TestFailureMessageStaysBounded(t *testing.T) {
	rec := FailureRecord{
		Stage:     StageStartup,
		ModelName: "big.gguf",
		Output:    strings.Repeat("noise\n", 5000),
	}
	message := rec.UserMessage()
	if len(message) > 1200 {
		t.Fatalf("user message is %d bytes, want it bounded", len(message))
	}
}
