package llama

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"howl-chat/internal/audio/types"
)

func TestTranscribePostsFileAndCleansText(t *testing.T) {
	var sawFile, sawModel bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseMultipartForm(1 << 20); err != nil {
			t.Errorf("parse multipart: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		file, _, err := r.FormFile("file")
		if err == nil {
			sawFile = true
			_, _ = io.Copy(io.Discard, file)
			_ = file.Close()
		}
		sawModel = r.FormValue("model") == "qwen3-asr"
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{"text":"language English<asr_text>Hello from the route."}`)
	}))
	defer server.Close()

	recognizer := New()
	recognizer.baseURL = server.URL
	recognizer.SetInitialized(types.ASRConfig{})

	audioPath := filepath.Join(t.TempDir(), "sample.wav")
	if err := os.WriteFile(audioPath, []byte("RIFF....WAVE"), 0o644); err != nil {
		t.Fatal(err)
	}

	result, err := recognizer.Transcribe(context.Background(), audioPath, nil)
	if err != nil {
		t.Fatalf("Transcribe: %v", err)
	}
	if !sawFile || !sawModel {
		t.Fatalf("request missing file=%t model=%t", sawFile, sawModel)
	}
	if result.Text != "Hello from the route." {
		t.Fatalf("text = %q", result.Text)
	}
	if result.Language != "english" {
		t.Fatalf("language = %q", result.Language)
	}
}

func TestTranscribeHonorsCancellation(t *testing.T) {
	release := make(chan struct{})
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-release:
		case <-r.Context().Done():
		}
	}))
	defer func() {
		close(release)
		server.Close()
	}()

	recognizer := New()
	recognizer.baseURL = server.URL
	recognizer.SetInitialized(types.ASRConfig{})

	audioPath := filepath.Join(t.TempDir(), "sample.wav")
	if err := os.WriteFile(audioPath, []byte("RIFF....WAVE"), 0o644); err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	if _, err := recognizer.Transcribe(ctx, audioPath, nil); err == nil {
		t.Fatal("Transcribe ignored cancellation")
	}
}

func TestTranscribeReportsHTTPFailure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, `{"error":{"message":"model does not support audio"}}`, http.StatusInternalServerError)
	}))
	defer server.Close()

	recognizer := New()
	recognizer.baseURL = server.URL
	recognizer.SetInitialized(types.ASRConfig{})

	audioPath := filepath.Join(t.TempDir(), "sample.wav")
	if err := os.WriteFile(audioPath, []byte("RIFF....WAVE"), 0o644); err != nil {
		t.Fatal(err)
	}

	if _, err := recognizer.Transcribe(context.Background(), audioPath, nil); err == nil {
		t.Fatal("Transcribe accepted a failed HTTP response")
	}
}
