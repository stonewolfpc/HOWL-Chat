package main

import (
	"context"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"howl-chat/internal/audio"
	audiotypes "howl-chat/internal/audio/types"
	"howl-chat/internal/backend/chat"
)

func testWAVData(samples int) []byte {
	dataSize := samples * 2
	data := make([]byte, 44+dataSize)
	copy(data[0:4], "RIFF")
	binary.LittleEndian.PutUint32(data[4:8], uint32(36+dataSize))
	copy(data[8:12], "WAVE")
	copy(data[12:16], "fmt ")
	binary.LittleEndian.PutUint32(data[16:20], 16)
	binary.LittleEndian.PutUint16(data[20:22], 1)
	binary.LittleEndian.PutUint16(data[22:24], 1)
	binary.LittleEndian.PutUint32(data[24:28], 16000)
	binary.LittleEndian.PutUint32(data[28:32], 32000)
	binary.LittleEndian.PutUint16(data[32:34], 2)
	binary.LittleEndian.PutUint16(data[34:36], 16)
	copy(data[36:40], "data")
	binary.LittleEndian.PutUint32(data[40:44], uint32(dataSize))
	return data
}

func testWAVDataURL(data []byte) string {
	return "data:audio/wav;base64," + base64.StdEncoding.EncodeToString(data)
}

func assertFileGone(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatalf("temporary audio still exists at %s (stat err: %v)", path, err)
	}
}

func TestTranscribeTempAudioCleansUpOnSuccess(t *testing.T) {
	var path string
	got, err := transcribeTempAudio(context.Background(), testWAVData(4), func(_ context.Context, p string) (string, error) {
		path = p
		if _, err := os.Stat(p); err != nil {
			t.Fatalf("temp audio missing during transcription: %v", err)
		}
		if filepath.Ext(p) != ".wav" {
			t.Fatalf("temp audio extension = %q, want .wav", filepath.Ext(p))
		}
		return "hello", nil
	})
	if err != nil {
		t.Fatalf("transcribeTempAudio: %v", err)
	}
	if got != "hello" {
		t.Fatalf("transcript = %q, want hello", got)
	}
	if path == "" {
		t.Fatal("transcription callback never received a path")
	}
	assertFileGone(t, path)
}

func TestTranscribeTempAudioCleansUpOnError(t *testing.T) {
	wantErr := errors.New("recognizer failed")
	var path string
	_, err := transcribeTempAudio(context.Background(), testWAVData(4), func(_ context.Context, p string) (string, error) {
		path = p
		return "", wantErr
	})
	if !errors.Is(err, wantErr) {
		t.Fatalf("error = %v, want %v", err, wantErr)
	}
	assertFileGone(t, path)
}

func TestTranscribeTempAudioCleansUpOnCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	var path string
	_, err := transcribeTempAudio(ctx, testWAVData(4), func(ctx context.Context, p string) (string, error) {
		path = p
		cancel()
		return "", ctx.Err()
	})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("error = %v, want context.Canceled", err)
	}
	assertFileGone(t, path)
}

func TestTranscribeTempAudioUniquePathsAreCleanedUpConcurrently(t *testing.T) {
	const workers = 8
	paths := make([]string, workers)
	errs := make(chan error, workers)
	var wg sync.WaitGroup
	start := make(chan struct{})

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			<-start
			_, err := transcribeTempAudio(context.Background(), testWAVData(4), func(_ context.Context, p string) (string, error) {
				paths[i] = p
				return "ok", nil
			})
			errs <- err
		}(i)
	}
	close(start)
	wg.Wait()
	close(errs)

	for err := range errs {
		if err != nil {
			t.Fatalf("concurrent transcription: %v", err)
		}
	}

	seen := make(map[string]bool, workers)
	for _, path := range paths {
		if path == "" {
			t.Fatal("missing path from concurrent transcription")
		}
		if seen[path] {
			t.Fatalf("two recordings shared temp path %s", path)
		}
		seen[path] = true
		assertFileGone(t, path)
	}
}

func TestTranscribeAudioDataUsesManagedPipeline(t *testing.T) {
	var pathSeenByPipeline string
	pipeline := &appFakePipeline{
		process: func(_ context.Context, msg audiotypes.AudioMessage) (*audiotypes.AudioResponse, error) {
			pathSeenByPipeline = msg.AudioPath
			if _, err := os.Stat(pathSeenByPipeline); err != nil {
				t.Fatalf("pipeline received a missing temp file: %v", err)
			}
			return &audiotypes.AudioResponse{TranscribedText: "from binding"}, nil
		},
	}
	app, manager := newAudioTestApp(t, func(_ context.Context, _ audio.Settings, candidate *audio.Candidate) error {
		*candidate = *appCandidate(pipeline, true, false, nil, nil)
		return nil
	})
	if _, _, err := manager.Apply(context.Background(), audio.Settings{ASREnabled: true}); err != nil {
		t.Fatalf("Apply: %v", err)
	}

	got, err := app.TranscribeAudioData(testWAVDataURL(testWAVData(4)))
	if err != nil {
		t.Fatalf("TranscribeAudioData: %v", err)
	}
	if got != "from binding" {
		t.Fatalf("transcript = %q, want from binding", got)
	}
	if pathSeenByPipeline == "" {
		t.Fatal("pipeline did not receive a temp path")
	}
	assertFileGone(t, pathSeenByPipeline)
}

func TestMicrophoneServiceUsesWAVDataBindingAndComposer(t *testing.T) {
	service, err := os.ReadFile("frontend/dist/js/audio-service.js")
	if err != nil {
		t.Fatalf("read audio service: %v", err)
	}
	text := string(service)
	// Object URLs are legitimate for browser playback (the TTS path turns the
	// GetTTSAudioData bytes into a blob URL). The boundary this test guards is
	// the old capture handoff: a browser object URL passed to the backend as if
	// it were a filesystem path.
	for _, forbidden := range []string{
		"MediaRecorder",
		"saveBlobToTempFile",
		"getElementById('chat-input')",
		"TranscribeAudio(tempPath)",
	} {
		if strings.Contains(text, forbidden) {
			t.Fatalf("audio service still contains %q", forbidden)
		}
	}
	for _, want := range []string{
		"AudioWav.encodeWav16",
		"TranscribeAudioData",
		"getElementById('message-input')",
	} {
		if !strings.Contains(text, want) {
			t.Fatalf("audio service is missing %q", want)
		}
	}

	index, err := os.ReadFile("frontend/dist/index.html")
	if err != nil {
		t.Fatalf("read index: %v", err)
	}
	wavIndex := strings.Index(string(index), "js/audio-wav.js")
	serviceIndex := strings.Index(string(index), "js/audio-service.js")
	if wavIndex < 0 || serviceIndex < 0 {
		t.Fatal("index.html is missing the audio-wav or audio-service script")
	}
	if wavIndex > serviceIndex {
		t.Fatal("audio-wav.js must load before audio-service.js")
	}
}

func TestTranscribedTextPersistsInChatHistory(t *testing.T) {
	pipeline := &appFakePipeline{
		process: func(_ context.Context, _ audiotypes.AudioMessage) (*audiotypes.AudioResponse, error) {
			return &audiotypes.AudioResponse{TranscribedText: "persisted transcript"}, nil
		},
	}
	app, manager := newAudioTestApp(t, func(_ context.Context, _ audio.Settings, candidate *audio.Candidate) error {
		*candidate = *appCandidate(pipeline, true, false, nil, nil)
		return nil
	})
	if _, _, err := manager.Apply(context.Background(), audio.Settings{ASREnabled: true}); err != nil {
		t.Fatalf("Apply: %v", err)
	}

	persistence, err := chat.NewPersistenceManagerWithDir(t.TempDir())
	if err != nil {
		t.Fatalf("NewPersistenceManagerWithDir: %v", err)
	}
	app.conversation = chat.NewActiveConversation(persistence)
	meta, err := app.NewChat("audio-history-test")
	if err != nil {
		t.Fatalf("NewChat: %v", err)
	}

	transcript, err := app.TranscribeAudioData(testWAVDataURL(testWAVData(4)))
	if err != nil {
		t.Fatalf("TranscribeAudioData: %v", err)
	}
	if err := app.AppendChatMessage("user", transcript, "User"); err != nil {
		t.Fatalf("AppendChatMessage: %v", err)
	}
	if err := app.SaveChat(); err != nil {
		t.Fatalf("SaveChat: %v", err)
	}
	loaded, err := app.LoadChat(meta.ID)
	if err != nil {
		t.Fatalf("LoadChat: %v", err)
	}
	for _, entry := range loaded.History {
		if entry.Text == transcript {
			return
		}
	}
	t.Fatalf("history after reload = %+v, want transcript %q", loaded.History, transcript)
}

func TestTranscribeAudioDataRejectsNonWAV(t *testing.T) {
	app := &App{}
	_, err := app.TranscribeAudioData(testWAVDataURL([]byte("not a wav")))
	if err == nil || !strings.Contains(err.Error(), "WAV") {
		t.Fatalf("error = %v, want a WAV validation error", err)
	}
}

func TestDecodeAudioBase64StripsDataURLPrefix(t *testing.T) {
	want := testWAVData(2)
	got, err := decodeAudioBase64(testWAVDataURL(want))
	if err != nil {
		t.Fatalf("decodeAudioBase64: %v", err)
	}
	if string(got) != string(want) {
		t.Fatal("decoded audio differs from input")
	}
}
