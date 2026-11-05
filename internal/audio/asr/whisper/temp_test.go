package whisper

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"howl-chat/internal/audio/types"
)

func whisperTestModel(t *testing.T) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "ggml-tiny.bin")
	if err := os.WriteFile(path, []byte("model"), 0o644); err != nil {
		t.Fatalf("write model: %v", err)
	}
	return path
}

func whisperTestCLI(t *testing.T) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "whisper-cli")
	if err := os.WriteFile(path, []byte("#!/bin/sh\n"), 0o755); err != nil {
		t.Fatalf("write cli: %v", err)
	}
	return path
}

func newWhisperTestRecognizer(t *testing.T, cli string, limiter *Limiter) *Recognizer {
	t.Helper()
	r := NewWithLimiter(types.ASRWhisperTiny, limiter)
	r.cliSearchPaths = []string{cli}
	r.checkHardware = func() error { return nil }
	return r
}

func writeWhisperWAV(t *testing.T) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "input.wav")
	data := make([]byte, 64)
	copy(data[0:4], "RIFF")
	copy(data[8:12], "WAVE")
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatalf("write wav: %v", err)
	}
	return path
}

func TestWhisperRecognizersOwnDistinctTempDirs(t *testing.T) {
	tempRoot := t.TempDir()
	t.Setenv("TMPDIR", tempRoot)
	t.Setenv("TMP", tempRoot)
	t.Setenv("TEMP", tempRoot)
	model := whisperTestModel(t)
	cli := whisperTestCLI(t)
	cfg := types.ASRConfig{ModelPath: model}

	r1 := newWhisperTestRecognizer(t, cli, NewLimiter(1))
	r2 := newWhisperTestRecognizer(t, cli, NewLimiter(1))
	if err := r1.Initialize(context.Background(), cfg); err != nil {
		t.Fatalf("initialize r1: %v", err)
	}
	if err := r2.Initialize(context.Background(), cfg); err != nil {
		t.Fatalf("initialize r2: %v", err)
	}

	dir1, dir2 := r1.tempDir, r2.tempDir
	if dir1 == "" || dir2 == "" {
		t.Fatalf("temp dirs not set: %q %q", dir1, dir2)
	}
	if dir1 == dir2 {
		t.Fatalf("recognizers share temp dir %q", dir1)
	}
	shared := filepath.Join(tempRoot, "howl-whisper")
	if dir1 == shared || dir2 == shared {
		t.Fatalf("recognizer uses shared temp dir %q", shared)
	}
	if err := os.WriteFile(filepath.Join(dir1, "a.txt"), []byte("a"), 0o644); err != nil {
		t.Fatalf("write r1 artifact: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir2, "b.txt"), []byte("b"), 0o644); err != nil {
		t.Fatalf("write r2 artifact: %v", err)
	}

	if err := r1.Release(); err != nil {
		t.Fatalf("release r1: %v", err)
	}
	if _, err := os.Stat(dir1); !os.IsNotExist(err) {
		t.Fatalf("r1 temp dir still exists after release: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir2, "b.txt")); err != nil {
		t.Fatalf("r2 artifact removed with r1: %v", err)
	}
	if !r2.IsInitialized() {
		t.Fatal("r2 became uninitialized after r1 release")
	}
	if err := r2.Release(); err != nil {
		t.Fatalf("release r2: %v", err)
	}
	if _, err := os.Stat(dir2); !os.IsNotExist(err) {
		t.Fatalf("r2 temp dir still exists after release: %v", err)
	}
	if _, err := os.Stat(shared); !os.IsNotExist(err) {
		t.Fatalf("shared howl-whisper directory was created: %v", err)
	}
}

func TestWhisperInitializeFailureCleansTempDir(t *testing.T) {
	tempRoot := t.TempDir()
	t.Setenv("TMPDIR", tempRoot)
	t.Setenv("TMP", tempRoot)
	t.Setenv("TEMP", tempRoot)
	model := whisperTestModel(t)
	cli := whisperTestCLI(t)
	r := newWhisperTestRecognizer(t, cli, NewLimiter(1))
	r.checkHardware = func() error { return errors.New("hardware check failed") }

	if err := r.Initialize(context.Background(), types.ASRConfig{ModelPath: model}); err == nil {
		t.Fatal("Initialize succeeded despite hardware failure")
	}
	if r.tempDir != "" {
		t.Fatalf("temp dir retained after failed initialize: %q", r.tempDir)
	}
	matches, err := filepath.Glob(filepath.Join(tempRoot, "howl-whisper-*"))
	if err != nil {
		t.Fatalf("glob temp dirs: %v", err)
	}
	if len(matches) != 0 {
		t.Fatalf("failed initialize left temp dirs: %v", matches)
	}
}

func TestWhisperInitializeDoesNotDownloadCLI(t *testing.T) {
	model := whisperTestModel(t)
	r := NewWithLimiter(types.ASRWhisperTiny, NewLimiter(1))
	r.cliSearchPaths = []string{}
	r.checkHardware = func() error { return nil }

	oldWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	work := t.TempDir()
	if err := os.Chdir(work); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	defer func() {
		if err := os.Chdir(oldWD); err != nil {
			t.Errorf("restore cwd: %v", err)
		}
	}()

	err = r.Initialize(context.Background(), types.ASRConfig{ModelPath: model})
	if err == nil {
		t.Fatal("Initialize succeeded without whisper-cli")
	}
	var audioErr *types.AudioError
	if !errors.As(err, &audioErr) {
		t.Fatalf("error type = %T, want *types.AudioError", err)
	}
	if audioErr.Code != types.ErrCodeModelNotFound {
		t.Fatalf("error code = %v, want model-not-found", audioErr.Code)
	}
	if !strings.Contains(err.Error(), "INSTALLATION GUIDE") {
		t.Fatalf("error lacks installation instructions: %v", err)
	}
	if r.tempDir != "" {
		t.Fatalf("temp dir set after failed initialize: %q", r.tempDir)
	}
	if _, statErr := os.Stat(filepath.Join(work, "bin", "whisper")); !os.IsNotExist(statErr) {
		t.Fatalf("initialization created a download directory: %v", statErr)
	}
}
