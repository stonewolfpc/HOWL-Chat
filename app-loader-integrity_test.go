package main

import (
	"encoding/binary"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"howl-chat/internal/image/registry"
)

// writeLoadTestGGUF writes a minimal GGUF v3 header with metadata so LoadModel
// exercises real GGUF admission rather than a stub classifier.
func writeLoadTestGGUF(t *testing.T, dir, name string, kvs map[string]any) string {
	t.Helper()
	path := filepath.Join(dir, name)
	file, err := os.Create(path)
	if err != nil {
		t.Fatalf("create GGUF: %v", err)
	}
	defer file.Close()

	writeString := func(value string) {
		var size [8]byte
		binary.LittleEndian.PutUint64(size[:], uint64(len(value)))
		if _, err := file.Write(size[:]); err != nil {
			t.Fatalf("write string size: %v", err)
		}
		if _, err := file.WriteString(value); err != nil {
			t.Fatalf("write string: %v", err)
		}
	}

	if err := binary.Write(file, binary.LittleEndian, uint32(0x46554747)); err != nil {
		t.Fatalf("write magic: %v", err)
	}
	if err := binary.Write(file, binary.LittleEndian, uint32(3)); err != nil {
		t.Fatalf("write version: %v", err)
	}
	if err := binary.Write(file, binary.LittleEndian, uint64(0)); err != nil {
		t.Fatalf("write tensor count: %v", err)
	}
	if err := binary.Write(file, binary.LittleEndian, uint64(len(kvs))); err != nil {
		t.Fatalf("write kv count: %v", err)
	}
	for key, value := range kvs {
		writeString(key)
		switch typed := value.(type) {
		case bool:
			if err := binary.Write(file, binary.LittleEndian, uint32(7)); err != nil {
				t.Fatalf("write bool type: %v", err)
			}
			flag := byte(0)
			if typed {
				flag = 1
			}
			if _, err := file.Write([]byte{flag}); err != nil {
				t.Fatalf("write bool: %v", err)
			}
		case uint32:
			if err := binary.Write(file, binary.LittleEndian, uint32(4)); err != nil {
				t.Fatalf("write uint32 type: %v", err)
			}
			if err := binary.Write(file, binary.LittleEndian, typed); err != nil {
				t.Fatalf("write uint32: %v", err)
			}
		case string:
			if err := binary.Write(file, binary.LittleEndian, uint32(8)); err != nil {
				t.Fatalf("write string type: %v", err)
			}
			writeString(typed)
		case []string:
			if err := binary.Write(file, binary.LittleEndian, uint32(9)); err != nil {
				t.Fatalf("write array type: %v", err)
			}
			if err := binary.Write(file, binary.LittleEndian, uint32(8)); err != nil {
				t.Fatalf("write array element type: %v", err)
			}
			if err := binary.Write(file, binary.LittleEndian, uint64(len(typed))); err != nil {
				t.Fatalf("write array count: %v", err)
			}
			for _, item := range typed {
				writeString(item)
			}
		default:
			t.Fatalf("unsupported test metadata type %T", value)
		}
	}
	return path
}

// TestLoaderSlowChild is a live process used to prove that a rejected load does
// not stop the model that is already running.
func TestLoaderSlowChild(t *testing.T) {
	if os.Getenv("HOWL_LOADER_SLOW_CHILD") != "1" {
		return
	}
	time.Sleep(30 * time.Second)
}

func TestLoadModelRejectsVocoderBeforeSpawning(t *testing.T) {
	app := &App{loaderDiagnosticsDir: t.TempDir()}
	path := writeLoadTestGGUF(t, t.TempDir(), "vocoder-LFM2.5-Audio-1.5B-Q4_0.gguf", map[string]any{
		"general.architecture": "this model cannot be used as LLM, use it via --model-vocoder in TTS examples",
	})

	err := app.LoadModel(path)
	if err == nil {
		t.Fatal("LoadModel accepted the vocoder")
	}
	if !strings.Contains(err.Error(), "cannot be used as a chat model") {
		t.Fatalf("error = %v, want the chat-model refusal", err)
	}
	if !strings.Contains(err.Error(), "vocoder") {
		t.Fatalf("error = %v, want it to identify the vocoder", err)
	}
	if app.level0().Processes.Current() != nil {
		t.Fatal("LoadModel started a server for a rejected asset")
	}
}

func TestLoadModelRejectsUnknownArchitectureWithActionableMessage(t *testing.T) {
	app := &App{loaderDiagnosticsDir: t.TempDir()}
	path := writeLoadTestGGUF(t, t.TempDir(), "mystery.gguf", map[string]any{
		"general.architecture": "brandnewarch",
	})

	err := app.LoadModel(path)
	if err == nil {
		t.Fatal("LoadModel accepted an unknown architecture")
	}
	if !strings.Contains(err.Error(), "brandnewarch") {
		t.Fatalf("error = %v, want it to name the architecture", err)
	}
	if !strings.Contains(err.Error(), "not a verified chat architecture") {
		t.Fatalf("error = %v, want the fail-closed explanation", err)
	}
}

func TestRejectedLoadLeavesTheRunningServerAlive(t *testing.T) {
	child := exec.Command(os.Args[0], "-test.run=TestLoaderSlowChild", "--")
	child.Env = append(os.Environ(), "HOWL_LOADER_SLOW_CHILD=1")
	if err := child.Start(); err != nil {
		t.Fatalf("start slow child: %v", err)
	}
	defer func() {
		_ = child.Process.Kill()
		_ = child.Wait()
	}()

	app := &App{loaderDiagnosticsDir: t.TempDir()}
	app.level0().Processes.Adopt(child, nil)
	path := writeLoadTestGGUF(t, t.TempDir(), "vocoder.gguf", map[string]any{
		"general.architecture": "this model cannot be used as LLM, use it via --model-vocoder in TTS examples",
	})

	if err := app.LoadModel(path); err == nil {
		t.Fatal("LoadModel accepted the vocoder")
	}
	if child.ProcessState != nil {
		t.Fatal("a rejected load stopped the already-running llama-server")
	}
	if got := app.level0().Processes.Current(); got != child {
		t.Fatal("a rejected load replaced the running server handle")
	}
}

func TestLLMInfosFiltersComponentsFromTheChatPicker(t *testing.T) {
	dir := t.TempDir()
	chat := writeLoadTestGGUF(t, dir, "chat-Q4_K_M.gguf", map[string]any{
		"general.architecture": "llama",
	})
	vocoder := writeLoadTestGGUF(t, dir, "vocoder-LFM2.5-Audio-1.5B-Q4_0.gguf", map[string]any{
		"general.architecture": "this model cannot be used as LLM, use it via --model-vocoder in TTS examples",
	})
	projector := writeLoadTestGGUF(t, dir, "mmproj-chat-Q8_0.gguf", map[string]any{
		"general.architecture":    "clip",
		"clip.has_vision_encoder": true,
	})

	infos := llmInfos([]registry.Entry{
		{Name: "chat", Path: chat, Type: registry.KindLLM, Format: "gguf"},
		{Name: "vocoder", Path: vocoder, Type: registry.KindLLM, Format: "gguf"},
		{Name: "mmproj-chat", Path: projector, Type: registry.KindLLM, Format: "gguf"},
	})

	if len(infos) != 1 || infos[0].Path != chat {
		t.Fatalf("picker entries = %+v, want only the chat model", infos)
	}
}

func TestLoadModelRejectsTaskIneligibleModelsBeforeSpawning(t *testing.T) {
	tests := []struct {
		name     string
		file     string
		kvs      map[string]any
		wantText string
	}{
		{
			name: "qwen3-asr-0.6b",
			file: "Qwen3-ASR-0.6B-Q8_0.gguf",
			kvs: map[string]any{
				"general.architecture": "qwen3vl",
				"general.tags":         []string{"automatic-speech-recognition"},
			},
			wantText: "speech-recognition",
		},
		{
			name: "qwen3-asr-1.7b",
			file: "Qwen3-ASR-1.7B-Q8_0.gguf",
			kvs: map[string]any{
				"general.architecture": "qwen3vl",
			},
			wantText: "speech-recognition",
		},
		{
			name: "ultravox",
			file: "Ultravox-Llama-3.2-1B-Q4_K_M.gguf",
			kvs: map[string]any{
				"general.architecture": "llama",
				"general.name":         "Llama 3.2 1B Instruct",
			},
			wantText: "audio-specialized",
		},
		{
			name: "voxtral",
			file: "Voxtral-Mini-3B-2507-Q4_K_M.gguf",
			kvs: map[string]any{
				"general.architecture": "llama",
				"general.tags":         []string{"audio-text-to-text"},
			},
			wantText: "audio-specialized",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			app := &App{loaderDiagnosticsDir: t.TempDir()}
			path := writeLoadTestGGUF(t, t.TempDir(), test.file, test.kvs)

			err := app.LoadModel(path)
			if err == nil {
				t.Fatal("LoadModel accepted a task-ineligible model")
			}
			if !strings.Contains(err.Error(), test.wantText) {
				t.Fatalf("error = %v, want %q", err, test.wantText)
			}
			if app.level0().Processes.Current() != nil {
				t.Fatal("LoadModel started a server for a task-ineligible model")
			}
		})
	}
}

func TestLoadModelRejectsWhisperComponentBinBeforeSpawning(t *testing.T) {
	app := &App{loaderDiagnosticsDir: t.TempDir()}
	path := filepath.Join(t.TempDir(), "ggml-base.bin")
	if err := os.WriteFile(path, []byte("legacy whisper component"), 0o644); err != nil {
		t.Fatalf("write legacy bin: %v", err)
	}

	err := app.LoadModel(path)
	if err == nil {
		t.Fatal("LoadModel accepted a legacy .bin as chat")
	}
	if !strings.Contains(err.Error(), "denied for chat by default") {
		t.Fatalf("error = %v, want the legacy default-deny explanation", err)
	}
	if app.level0().Processes.Current() != nil {
		t.Fatal("LoadModel started a server for a legacy .bin")
	}
}

func TestLLMInfosExcludesTaskIneligibleAndLegacyBinEntries(t *testing.T) {
	dir := t.TempDir()
	chat := writeLoadTestGGUF(t, dir, "chat-Q4_K_M.gguf", map[string]any{
		"general.architecture": "llama",
	})
	asr := writeLoadTestGGUF(t, dir, "Qwen3-ASR-1.7B-Q8_0.gguf", map[string]any{
		"general.architecture": "qwen3vl",
		"general.tags":         []string{"automatic-speech-recognition"},
	})
	audio := writeLoadTestGGUF(t, dir, "Ultravox-Llama-3.2-1B-Q4_K_M.gguf", map[string]any{
		"general.architecture": "llama",
	})
	omni := writeLoadTestGGUF(t, dir, "Qwen2.5-Omni-3B-Q4_K_M.gguf", map[string]any{
		"general.architecture": "qwen2vl",
		"general.tags":         []string{"multimodal", "any-to-any"},
	})
	legacy := filepath.Join(dir, "ggml-base.bin")
	if err := os.WriteFile(legacy, []byte("legacy whisper component"), 0o644); err != nil {
		t.Fatalf("write legacy bin: %v", err)
	}

	infos := llmInfos([]registry.Entry{
		{Name: "chat", Path: chat, Type: registry.KindLLM, Format: "gguf"},
		{Name: "asr", Path: asr, Type: registry.KindLLM, Format: "gguf"},
		{Name: "audio", Path: audio, Type: registry.KindLLM, Format: "gguf"},
		{Name: "omni", Path: omni, Type: registry.KindLLM, Format: "gguf"},
		{Name: "whisper", Path: legacy, Type: registry.KindLLM, Format: "bin"},
	})

	if len(infos) != 2 {
		t.Fatalf("picker entries = %+v, want text chat and Omni only", infos)
	}
	got := map[string]bool{}
	for _, info := range infos {
		got[info.Path] = true
	}
	if !got[chat] || !got[omni] {
		t.Fatalf("picker entries = %+v, want chat and Omni", infos)
	}
}

func TestScanDirForTypeDoesNotRegisterLegacyBinForChat(t *testing.T) {
	dir := t.TempDir()
	chat := writeLoadTestGGUF(t, dir, "chat-Q4_K_M.gguf", map[string]any{
		"general.architecture": "llama",
	})
	if err := os.WriteFile(filepath.Join(dir, "ggml-base.bin"), []byte("legacy whisper component"), 0o644); err != nil {
		t.Fatalf("write legacy bin: %v", err)
	}

	entries, _, err := scanDirForType(dir, registry.KindLLM, []string{".gguf"}, nil)
	if err != nil {
		t.Fatalf("scanDirForType: %v", err)
	}
	if len(entries) != 1 || entries[0].Path != chat {
		t.Fatalf("entries = %+v, want only the GGUF chat model", entries)
	}
}

func TestBuildLlamaServerArgsUsesOnlyThePairedProjector(t *testing.T) {
	projector := filepath.Join(`C:\voice`, "mmproj-Qwen3-ASR-1.7B-Q8_0.gguf")
	args := buildLlamaServerArgs("model.gguf", projector, map[string]interface{}{}, false)
	joined := strings.Join(args, " ")

	if strings.Count(joined, "--mmproj") != 1 {
		t.Fatalf("args = %v, want exactly one --mmproj", args)
	}
	index := -1
	for i, arg := range args {
		if arg == "--mmproj" {
			index = i
			break
		}
	}
	if index < 0 || index+1 >= len(args) || args[index+1] != projector {
		t.Fatalf("args = %v, want --mmproj followed by %q", args, projector)
	}
	if !strings.Contains(joined, "--fit off") {
		t.Fatalf("args = %v, want --fit off when fitParams is false", args)
	}

	without := buildLlamaServerArgs("model.gguf", "", map[string]interface{}{}, true)
	if strings.Contains(strings.Join(without, " "), "--mmproj") {
		t.Fatalf("args = %v, want no projector argument", without)
	}
}
